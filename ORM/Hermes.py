"""
This is the deepest API level which interfaces with SQLAlchemy directly. It's Called Hermes.
This is a synchronous API that only serves the frontend. Network events are handled in Mercury.

"""

import time, random, ujson
from datetime import date, timedelta
from sqlalchemy import func, or_
from twisted.internet import threads

from globals import basedir, refreshBackendValuesGatheredFromJson
import globals
from ORM.models import *
from ORM import Demeter
from DecisionEngine import eventLoop

session = Session()

def timeit(method):

    def timed(*args, **kw):
        ts = time.time() *1000
        result = method(*args, **kw)
        te = time.time() *1000

        print '%r (%r, %r) %2.2f ms' % \
              (method.__name__, args, kw, te-ts)
        return result

    return timed

# The first block is fetchers. These pertain to the read events.

def fetchSinglePost(fingerprint):
    """
        This is used for single precise shots. For example, you want to use this to fetch the
        subject post itself, because fetchPosts, when given a subject post, will return its
        constituents. This will fetch the post and only the post fingerprinted.

        This uses all() instead of one even when it is certain it will return one row,
        because using one() raises noResultFound exception, and I don't want that. Finding no
        results is a perfectly okay answer, it is not an exception. ... Well, shit, it turns out
    """

    post = session.query(Post).filter(Post.PostFingerprint == fingerprint).all()
    return post

def fetchUppermostTopics():
    posts = session.query(Post).filter(Post.ParentPostFingerprint == None).all()
    # At this point I don't know if their counts are zero or not. In the frontend, I'm filtering the topic:
    # a) If not locally created, and b) If has zero children.
    return posts

def fetchDirectDescendantPosts(fingerprint):
    """
    This returns the direct descendants of the fingerprint provided.
    """

    directDescendantPosts = session.query(Post).filter(Post.ParentPostFingerprint == fingerprint)\
        .order_by(Post.CreationDate.desc()).all()

    return directDescendantPosts

def fetchDirectDescendantPostsWithTimespan(fingerprint, daysToSubtract):
    """
    This function takes an amount of days and it returns the direct descendants produced within
    that timespan.
    """
    today = datetime.utcnow()
    targetTime = today - timedelta(days = daysToSubtract)
    #import ipdb; ipdb.set_trace()
    posts = session.query(Post)\
        .filter(Post.ParentPostFingerprint == fingerprint,
            Post.CreationDate >= targetTime).all()
    return posts

def fetchDescendantPosts(fingerprint, depth):
    """
    This puppy takes a depth argument. Zero means it will only return the fingerprint asked.
    1,2,3 as it goes.. This can replace almost all of code her.e
    """
    postsArray = []
    def resolveDescendancy(fp, array, depth):
        if fetchDirectDescendantPosts(fp) is not [] and depth > 0:
            depth = depth - 1
            array.append(fetchSinglePost(fp))
            for post in fetchDirectDescendantPosts(fp):
                resolveDescendancy(post.PostFingerprint, array, depth)
        else:
            array.append(fetchSinglePost(fp))

    resolveDescendancy(fingerprint, postsArray, depth)
    return postsArray

def fetchAllDescendantPosts(fingerprint):
    """
        This returns all descendants of a post. If you apply this to a topic, you'll get subjects
         and posts. If you apply this onto a subject, you'll get the entire subject tree. Don't
         apply this on a topic. This returns an array of posts, not a SQLAlchemy collection.

         In the future I should have this function also take a depth argument, so it won't return
          a massive output. [done, look at fetchDescendantPosts]
    """
    postsArray = []

    def resolveDescendancy(fp, array):
        if fetchDirectDescendantPosts(fp) is not []:
            for post in fetchDirectDescendantPosts(fp):
                array.append(fetchSinglePost(post.PostFingerprint))
                resolveDescendancy(post.PostFingerprint, array)
        else:
            array.append(fetchSinglePost(fp))

    resolveDescendancy(fingerprint, postsArray)

    return postsArray

def countDirectDescendantPosts(fingerprint):
    counter = session.query(Post).filter(Post.ParentPostFingerprint == fingerprint)\
        .count()
    return counter

def countTopics():
    counter = session.query(Post).filter(Post.ParentPostFingerprint == None)\
        .count()
    return counter

def countSubjects():
    topics = session.query(Post).filter(Post.ParentPostFingerprint == None)\
        .all()
    finalCount = 0
    for t in topics:
        finalCount += session.query(Post).filter(Post.ParentPostFingerprint == t.PostFingerprint).count()
    return finalCount

@timeit
def getSubjects(topicFingerprint, daysToSubtract):
    # This gets all subjects in a topic with appropriate external data attached (first post and comment count.)
    today = datetime.utcnow()
    targetTime = today - timedelta(days = daysToSubtract)
    #import ipdb; ipdb.set_trace()
    subjects = session.query(Post)\
        .filter(Post.ParentPostFingerprint == topicFingerprint,
            Post.CreationDate >= targetTime)\
        .order_by(Post.UpvoteCount.desc())\
        .all()

    subjectsArray = []
    for s in subjects:
        try:
            firstPost = session.query(Post).\
            filter(Post.ParentPostFingerprint ==s.PostFingerprint)\
            .order_by(Post.UpvoteCount.desc())\
            .first()
            s.FirstPostBody = firstPost.Body[:410] # The current shown amount is 400.
            s.FirstPostOwner = firstPost.OwnerUsername
        except:
            s.FirstPostBody = ''
            s.FirstPostOwner = ''
        subjectsArray.append(s)
    return subjectsArray
    #return posts

@timeit
def getHomeScreen(numberOfTopics, numberOfSubjects):
    # This returns the home screen. The number of topics shown and number of subjects inside each topic.
    result = []
    topiks = globals.selectedTopics
    print('selectedtopics as seen by homescreen get: ', globals.selectedTopics)
    topics = session.query(Post).filter(Post.ParentPostFingerprint == None)\
        .order_by(Post.ReplyCount.desc())\
        .filter(Post.PostFingerprint.in_(globals.selectedTopics))\
        .limit(numberOfTopics).all()

    popularSubjects = session.query(Post)\
        .filter(or_(Post.OwnerUsername != None,
                    Post.OwnerFingerprint != None),
                Post.Subject != None,
                Post.Body == '')\
        .order_by(Post.UpvoteCount.desc()).limit(numberOfSubjects).all()
        # TODO: Body doesn't do None checking, check '' instead. Get this consistent

    popularSubjectsAsJsons = []

    for subject in popularSubjects:
        subjectAsJson = ujson.dumps(subject.asDict(), ensure_ascii=False)
        popularSubjectsAsJsons.append(subjectAsJson)

    popularSubjectsDict = {
        'TopicName': 'Most Popular',
        'Subjects': popularSubjectsAsJsons
    }
    result.append(popularSubjectsDict)
    for topic in topics:
        subjects = session.query(Post).filter(Post.ParentPostFingerprint == topic.PostFingerprint)\
            .order_by(Post.UpvoteCount.desc()).limit(numberOfSubjects).all()

        subjectsAsJsons = []

        for subject in subjects:
            subjectAsJson = ujson.dumps(subject.asDict())
            subjectsAsJsons.append(subjectAsJson)

        topicDict = {
            'TopicName': topic.Subject,
            'Subjects': subjectsAsJsons
        }
        result.append(topicDict)
    # subject: has owner, has subject, but no body
    return result



def multiCountDirectDescendantPostsWithTimespan(daysToSubtract, fingerprintsAsDict):
    today = datetime.utcnow()
    if daysToSubtract != 0:
        targetTime = today - timedelta(days = daysToSubtract)
        for fingerprint in fingerprintsAsDict:
            counter = session.query(Post).filter(
                Post.ParentPostFingerprint == fingerprint,
                Post.CreationDate >= targetTime
            ).count()
            fingerprintsAsDict[fingerprint] = counter

            # Iterating over a dict, fingerprint gives me the key, and fingerprintAsDict[fingerprint]
            #  gives me the value of the key.
    else:
        for fingerprint in fingerprintsAsDict:
            counter = session.query(Post).filter(
                Post.ParentPostFingerprint == fingerprint
            ).count()
            fingerprintsAsDict[fingerprint] = counter


    return fingerprintsAsDict

def countAllDescendantPosts(fingerprint):
    """
    Converting all preliminary checks before diving to calculate (such as calling
    fetchDirectDescendants to know if there are anything to be added to the array could be made
    more performant by changing those calls to counts.
    """
    counter = 0

    def resolveDescendancy(fp, counter):
        if countDirectDescendantPosts(fp) != 0:
            for post in fetchDirectDescendantPosts(fp):
                counter += countDirectDescendantPosts(fp)
                return resolveDescendancy(post.PostFingerprint, counter)
        return counter
    counter = resolveDescendancy(fingerprint, counter)
    return counter

def getTopmostComment(fingerprint):
    """
    This function returns the first and highest upvoted comment in a subject.
    """
    post = session.query(Post).\
        filter(Post.ParentPostFingerprint ==fingerprint)\
        .order_by(Post.UpvoteCount.desc())\
        .first()
    return post

def getUser(fingerprint):
    """
        User in question is returned.
        Mind that this is only available for fingerprinted, i.e. registered users.
    """
    user = session.query(User).filter(User.Fingerprint == fingerprint).all()[0]
    return user

def getUnregisteredUserPosts(name):
    """
        This returns posts created by the unregistered user requested.
    """
    return session.query(Post).filter(Post.OwnerUsername == name).order_by(Post.CreationDate.desc()).all()


def readUserProfile():
    try:
        f = open(basedir+'UserProfile/UserProfile.json', 'rb')
    except:
        f = open(basedir+'UserProfile/UserProfile.json', 'wb')
        f.close()
        f = open(basedir+'UserProfile/UserProfile.json', 'rb')
    jsonAsText = f.read()
    f.close()
    return jsonAsText

def getSavedPosts():
    posts = session.query(Post).filter(Post.Saved == True, Post.Subject == u'')\
        .order_by(Post.CreationDate.desc()).all()
    return posts

def countConnectedNodes():
    # The current definition of a connected node I am using right now is a node that has updated or
    # confirmed its IP address in the last 30 minutes. So, lastseen should be below 30 mins. As
    # this definition changes, I need to update this function.
    targetTime = datetime.utcnow() - timedelta(minutes = 30)
    count = session.query(Node).filter(Node.LastConnectedDate > targetTime).count()
    return count

def getLocallyCreatedPosts():
    posts = session.query(Post).filter(Post.LocallyCreated == True).order_by(Post.CreationDate.desc()).all()
    return posts

def getReplies():
    return session.query(Post).filter(Post.IsReply == True).order_by(Post.CreationDate.desc()).all()

def countReplies():
    return session.query(Post).filter(Post.IsReply == True).count()

def getParentSubjectOfGivenPost(fingerprint):
    subjectFingerprint = __resolveAncestry(fingerprint)
    if subjectFingerprint is not False:
        return session.query(Post).filter(Post.PostFingerprint == subjectFingerprint).all()

def getLastConnectionTime():
    return session.query(Node).order_by(Node.LastConnectedDate.desc()).first().LastConnectedDate

def __resolveAncestry(fingerprint):
    post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
    if post.OwnerUsername == '' and post.OwnerFingerprint == '':
        # If topic
        return False
    elif post.Subject != '':
        # If subject
        return post.PostFingerprint
    else:
        #If post
        return __resolveAncestry(post.ParentPostFingerprint)



# These pertain to the write events.

def createTopic(topicName):
    ephemeralConnectionId = int(random.random()*10**10)
    topic = Post(Subject= topicName, LocallyCreated=True, EphemeralConnectionId=ephemeralConnectionId)
    topicHeader = PostHeader(PostFingerprint=topic.PostFingerprint, Dirty = False)
    session.add(topic)
    session.add(topicHeader)
    session.commit()
    print "This topic has been added to the database"
    Demeter.incrementAncestryCommentCount(ephemeralConnectionId)
    return topic.PostFingerprint

def createPost(postSubject, postText, parentFingerprint, ownerUsername, postLanguage):
    ephemeralConnectionId = int(random.random()*10**10)
    post = Post(Subject=postSubject, Body=postText, OwnerUsername=ownerUsername,
                ParentPostFingerprint=parentFingerprint, LocallyCreated=True, Language=postLanguage,
                EphemeralConnectionId=ephemeralConnectionId, Upvoted=True, UpvoteCount=1)


    parent = session.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
    postHeader = PostHeader(PostFingerprint=post.PostFingerprint,
                            ParentPostFingerprint=post.ParentPostFingerprint, Language=postLanguage
                            )

    localnode = session.query(Node).filter(Node.LastConnectedIP=='LOCAL').one()
    vote = Vote(Direction=1,Node=localnode, PostHeader=postHeader)

    session.add(post)
    session.add(postHeader)
    session.add(vote)
    session.commit()
    print "The post has been added to the database."
    Demeter.incrementAncestryCommentCount(ephemeralConnectionId)
    if post.Subject == '':
        eventLoop.marduk() # Only call if post, because otherwise it results in a double call. And there is no way
        # to create a subject without a post anyway.
    return post.PostFingerprint

def createSignedPost(postSubject, postText, parentFingerprint, ownerFingerprint):

    post = Post(Subject=postSubject, Body=postText, OwnerFingerprint=ownerFingerprint,
                ownerUsername = getUser(ownerFingerprint).Username,
                ParentPostFingerprint=parentFingerprint)
    session.add(post)
    session.commit()
    print "The signed post has been added to the database."
    return post.PostFingerprint

def writeUserProfile(userProfileInJSON):
    # This is the autocommitter.
    print('Autocommit Fired. Writing new user profile.')
    f = open(basedir+'UserProfile/UserProfile.json', 'wb')
    f.write(userProfileInJSON.encode('utf8'))
    f.close()
    refreshBackendValuesGatheredFromJson(userProfileInJSON)
    return True

def votePost(fingerprint, voteDirection):
    # Vote Direction is either 2, 1, 0 or -1. 2 means remove prior votes. (user clicks on the
    # button again to remove his / her vote)
    postheader = session.query(PostHeader).filter(PostHeader.PostFingerprint == fingerprint).one()
    post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
    localnode = session.query(Node).filter(Node.LastConnectedIP=='LOCAL').one()

    if voteDirection == 2: # If this is a remove vote call:

        if post.Upvoted == True:
            # And it's currently upvoted, remove effects,
            post.Upvoted = False
            post.UpvoteCount -= 1
            post.LastVoteDate = datetime.utcnow()
            session.add(post)
            # And the vote item itself.
            voteUp = session.query(Vote).filter(Vote.Direction == 1,
                                                Vote.PostHeader == postheader).one()
            session.delete(voteUp)
            session.commit()
            return True

        if post.Downvoted == True:
            # And it's currently downvoted, remove effects,
            post.Downvoted = False
            post.DownvoteCount -= 1
            post.LastVoteDate = datetime.utcnow()
            session.add(post)
            # And remove the item itself.
            voteDown = session.query(Vote).filter(Vote.Direction == -1,
                                                  Vote.PostHeader == postheader).one()
            session.delete(voteDown)
            session.commit()
            return True

    else:
        # If this is not a remove vote call:
        if not session.query(Vote).filter(Vote.PostHeader == postheader, Vote.Node == localnode).all():
            # If empty (If first vote on this post by this user)
            vote = Vote(Direction=voteDirection,Node=localnode, PostHeader=postheader)
            if voteDirection == 1:
                # If this is empty and first vote is positive
                vote.Direction = 1
                post.Upvoted = True
                post.UpvoteCount += 1
                post.LastVoteDate = datetime.utcnow()
                if post.Neutral == True and post.NeutralCount > 0: post.NeutralCount -= 1
                # If it is a neutral, it is not anymore and I should remove one from the counter. Why zero checking?
                # because marduk loops in long intervals and there can be a position where the post switches from
                # neutral to positive, and marduk hasn't incremented the neutral count yet. The entire reason I'm doing
                # this is that this allows me to keep NeutralCount consistent with Neutral flag at all times even before
                # Marduk hits.
                post.Neutral = False # in reverse action, neutral is set by marduk anyway,
                # so it'll eventually regain neutral status if needed.

            if voteDirection == -1:
                # If negative
                vote.Direction = -1
                post.Downvoted = True
                post.DownvoteCount += 1
                post.LastVoteDate = datetime.utcnow()
                if post.Neutral == True and post.NeutralCount > 0: post.NeutralCount -= 1
                post.Neutral = False
        else:
            # If this is not the first vote:
            vote = session.query(Vote).filter(Vote.PostHeader == postheader, Vote.Node == localnode).first()
            if voteDirection != vote.Direction:
                # If the user is actually changing vote direction:
                if voteDirection == 1:
                    # And if the new vote is positive:
                    vote.Direction = 1
                    post.Upvoted = True
                    post.UpvoteCount += 1
                    post.LastVoteDate = datetime.utcnow()
                    # Because we can be sure if the user is changing the vote direction and the
                    # new vote is positive, the old vote was negative.
                    post.Downvoted = False
                    post.DownvoteCount -= 1

                if voteDirection == -1:
                    # If new vote is negative:
                    vote.Direction = -1
                    post.Downvoted = True
                    post.DownvoteCount += 1
                    post.UpvoteCount -= 1
                    post.LastVoteDate = datetime.utcnow()
                    # If the vote is changing and the new one is -, the old was +.
                    post.Upvoted = False
        session.add(vote)
        session.add(post)
        session.commit()
        return True

def savePost(fingerprint):
    post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
    if post.Saved == True:
        post.Saved = False
    else:
        post.Saved = True
    session.add(post)
    session.commit()
    return True

def markAllRepliesAsRead():
    replies = session.query(Post).filter(Post.IsReply == True).all()
    for r in replies:
        r.IsReply = False
        session.add(r)

    session.commit()
    return True