"""
This is the deepest API level which interfaces with SQLAlchemy directly. It's Called Hermes.
This is a synchronous API that only serves the frontend. Network events are handled in Mercury.

Remember, this API runs in a different process. It does not have access to GlobalCommitter instance!
It can only fire an AMP command to let the parent process know.

"""

import time, random, ujson
from datetime import date, timedelta
from sqlalchemy import func, or_
from twisted.internet import threads

from globals import basedir, refreshBackendValuesGatheredFromJson, profiledir
import globals
from ORM.models import *
from twisted.python import reflect

import sys

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

class Hermes(object):
    def __init__(self, protInstance):
        self.protInstance = protInstance

    def fetchSinglePost(self, fingerprint):
        """
            This is used for single precise shots. For example, you want to use this to fetch the
            subject post itself, because fetchPosts, when given a subject post, will return its
            constituents. This will fetch the post and only the post fingerprinted.

            This uses all() instead of one even when it is certain it will return one row,
            because using one() raises noResultFound exception, and I don't want that. Finding no
            results is a perfectly okay answer, it is not an exception. ... Well, shit, it turns out
        """
        session = Session()
        post = session.query(Post).filter(Post.PostFingerprint == fingerprint).all()
        session.close()
        return post

    def fetchUppermostTopics(self):
        session = Session()
        posts = session.query(Post).filter(Post.ParentPostFingerprint == None).all()
        # At this point I don't know if their counts are zero or not. In the frontend, I'm filtering the topic:
        # a) If not locally created, and b) If has zero children.
        session.close()
        return posts

    def fetchDirectDescendantPosts(self, fingerprint):
        """
        This returns the direct descendants of the fingerprint provided.
        """
        session = Session()
        directDescendantPosts = session.query(Post).filter(Post.ParentPostFingerprint == fingerprint)\
            .order_by(Post.CreationDate.desc()).all()
        session.close()

        return directDescendantPosts

    def fetchDirectDescendantPostsWithTimespan(self, fingerprint, daysToSubtract):
        """
        This function takes an amount of days and it returns the direct descendants produced within
        that timespan.
        """
        today = datetime.utcnow()
        targetTime = today - timedelta(days = daysToSubtract)
        #import ipdb; ipdb.set_trace()
        session = Session()
        posts = session.query(Post)\
            .filter(Post.ParentPostFingerprint == fingerprint,
                Post.CreationDate >= targetTime).all()
        session.close()
        return posts

    def fetchDescendantPosts(self, fingerprint, depth):
        """
        This puppy takes a depth argument. Zero means it will only return the fingerprint asked.
        1,2,3 as it goes.. This can replace almost all of code her.e
        """
        postsArray = []
        def resolveDescendancy(self, fp, array, depth):
            if self.fetchDirectDescendantPosts(fp) is not [] and depth > 0:
                depth = depth - 1
                array.append(self.fetchSinglePost(fp))
                for post in self.fetchDirectDescendantPosts(fp):
                    resolveDescendancy(self, post.PostFingerprint, array, depth)
            else:
                array.append(self.fetchSinglePost(fp))

        resolveDescendancy(self, fingerprint, postsArray, depth)
        return postsArray

    def fetchAllDescendantPosts(self, fingerprint):
        """
            This returns all descendants of a post. If you apply this to a topic, you'll get subjects
             and posts. If you apply this onto a subject, you'll get the entire subject tree. Don't
             apply this on a topic. This returns an array of posts, not a SQLAlchemy collection.

             In the future I should have this function also take a depth argument, so it won't return
              a massive output. [done, look at fetchDescendantPosts]
        """
        postsArray = []

        def resolveDescendancy(self, fp, array):
            if self.fetchDirectDescendantPosts(fp) is not []:
                for post in self.fetchDirectDescendantPosts(fp):
                    array.append(self.fetchSinglePost(post.PostFingerprint))
                    resolveDescendancy(self, post.PostFingerprint, array)
            else:
                array.append(self.fetchSinglePost(fp))

        resolveDescendancy(self, fingerprint, postsArray)

        return postsArray

    def countDirectDescendantPosts(self, fingerprint):
        session = Session()
        counter = session.query(Post).filter(Post.ParentPostFingerprint == fingerprint)\
            .count()
        session.close()
        return counter

    def countTopics(self):
        session = Session()
        counter = session.query(Post).filter(Post.ParentPostFingerprint == None)\
            .count()
        session.close()
        return counter

    @timeit
    # This counts all subjects in all topics. It performs VERY poorly! Do not use unless absolutely needed!
    def countSubjects(self):
        # session = Session()
        # topics = session.query(Post).filter(Post.ParentPostFingerprint == None)\
        #     .all()
        # finalCount = 0
        # for t in topics:
        #     finalCount += session.query(Post).filter(Post.ParentPostFingerprint == t.PostFingerprint).count()
        # session.close()
        # return finalCount

        # A better implementation maybe? All subjects (where parent is a topic). A subject = there is an owner
        # but not a body.

        session = Session()
        subjectCount = session.query(Post)\
                .filter(or_(Post.OwnerUsername != None,
                            Post.OwnerFingerprint != None),
                        Post.Subject != None,
                        Post.Body == '').count()
        session.close()
        return subjectCount



    @timeit
    def getSubjects(self, topicFingerprint, daysToSubtract):
        # This gets all subjects in a topic with appropriate external data attached (first post and comment count.)
        today = datetime.utcnow()
        targetTime = today - timedelta(days = daysToSubtract)
        #import ipdb; ipdb.set_trace()
        session = Session()
        subjects = session.query(Post)\
            .filter(Post.ParentPostFingerprint == topicFingerprint,
                Post.CreationDate >= targetTime)\
            .order_by(Post.RankScore.desc())\
            .order_by(Post.CreationDate.desc())\
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
        session.close()
        return subjectsArray
        #return posts

    @timeit
    def getHomeScreen(self, numberOfTopics, numberOfSubjects):
        # This returns the home screen. The number of topics shown and number of subjects inside each topic.
        result = []
        topiks = globals.selectedTopics
        session = Session()
        topics = session.query(Post).filter(Post.ParentPostFingerprint == None)\
            .order_by(Post.ReplyCount.desc())\
            .filter(Post.PostFingerprint.in_(globals.selectedTopics))\
            .limit(numberOfTopics).all()

        popularSubjects = session.query(Post)\
            .filter(or_(Post.OwnerUsername != None,
                        Post.OwnerFingerprint != None),
                    Post.Subject != None,
                    Post.Body == '')\
            .order_by(Post.RankScore.desc())\
            .order_by(Post.CreationDate.desc()).limit(numberOfSubjects).all()
            # TODO: Body doesn't do None checking, check '' instead. Get this consistent
            # TODO this is the actual place where the stuff gets chosen. Implement frontend algorithm here.

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
                .order_by(Post.RankScore.desc())\
                .order_by(Post.CreationDate.desc()).limit(numberOfSubjects).all()

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
        session.close()
        return result



    def multiCountDirectDescendantPostsWithTimespan(self, daysToSubtract, fingerprintsAsDict):
        today = datetime.utcnow()
        session = Session()
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

        session.close()
        return fingerprintsAsDict

    def countAllDescendantPosts(self, fingerprint):
        """
        Converting all preliminary checks before diving to calculate (such as calling
        fetchDirectDescendants to know if there are anything to be added to the array could be made
        more performant by changing those calls to counts.
        """
        counter = 0

        def resolveDescendancy(self, fp, counter):
            if self.countDirectDescendantPosts(fp) != 0:
                for post in self.fetchDirectDescendantPosts(fp):
                    counter += self.countDirectDescendantPosts(fp)
                    return resolveDescendancy(self, post.PostFingerprint, counter)
            return counter
        counter = resolveDescendancy(self, fingerprint, counter)
        return counter

    def getTopmostComment(self, fingerprint):
        """
        This function returns the first and highest upvoted comment in a subject.
        """
        session = Session()
        post = session.query(Post).\
            filter(Post.ParentPostFingerprint ==fingerprint)\
            .order_by(Post.UpvoteCount.desc())\
            .first()
        session.close()
        return post

    def getUser(self, fingerprint):
        """
            User in question is returned.
            Mind that this is only available for fingerprinted, i.e. registered users.
        """
        session = Session()
        user = session.query(User).filter(User.Fingerprint == fingerprint).all()[0]
        session.close()
        return user

    def getUnregisteredUserPosts(self, name):
        """
            This returns posts created by the unregistered user requested.
        """
        session = Session()
        posts = session.query(Post).filter(Post.OwnerUsername == name).order_by(Post.CreationDate.desc()).all()
        session.close()
        return posts


    def readUserProfile(self):
        try:
            f = open(profiledir+'UserProfile/UserProfile.json', 'rb')
        except:
            f = open(profiledir+'UserProfile/UserProfile.json', 'wb')
            f.close()
            f = open(profiledir+'UserProfile/UserProfile.json', 'rb')
        jsonAsText = f.read()
        f.close()
        return jsonAsText

    def getSavedPosts(self):
        session = Session()
        posts = session.query(Post).filter(Post.Saved == True, Post.Subject == u'')\
            .order_by(Post.CreationDate.desc()).all()
        session.close()
        return posts

    def countConnectedNodes(self):
        # The current definition of a connected node I am using right now is a node that has updated or
        # confirmed its IP address in the last 30 minutes. So, lastseen should be below 30 mins. As
        # this definition changes, I need to update this function.
        targetTime = datetime.utcnow() - timedelta(minutes = 30)
        session = Session()
        count = session.query(Node).filter(Node.LastConnectedDate > targetTime).count()
        session.close()
        return count

    def getLocallyCreatedPosts(self):
        session = Session()
        posts = session.query(Post).filter(Post.LocallyCreated == True).order_by(Post.CreationDate.desc()).all()
        session.close()
        return posts

    def getReplies(self):
        session = Session()
        replies = session.query(Post).filter(Post.IsReply == True).order_by(Post.CreationDate.desc()).all()
        session.close()
        return replies

    def countReplies(self):
        session = Session()
        count = session.query(Post).filter(Post.IsReply == True).count()
        session.close()
        return count

    def getLastConnectionTime(self):
        session = Session()
        lastConnTime = session.query(Node).order_by(Node.LastConnectedDate.desc()).first().LastConnectedDate
        session.close()
        return lastConnTime

    def __resolveAncestry(self, fingerprint):
        session = Session()
        post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
        if post.OwnerUsername == '' and post.OwnerFingerprint == '':
            # If topic
            session.close()
            return False
        elif post.Subject != '':
            # If subject
            fingerprint = post.PostFingerprint
            session.close()
            return fingerprint
        else:
            #If post
            return self.__resolveAncestry(post.ParentPostFingerprint)

    def getParentSubjectOfGivenPost(self, fingerprint):
        subjectFingerprint = self.__resolveAncestry(fingerprint)
        if subjectFingerprint is not False:
            session = Session()
            parentSubject = session.query(Post).filter(Post.PostFingerprint == subjectFingerprint).all()
            session.close
            return parentSubject


    # These pertain to the write events.

    def createTopic(self, topicName):
        session = Session()
        ephemeralConnectionId = int(random.random()*10**10)
        topic = Post(Subject= topicName, LocallyCreated=True, EphemeralConnectionId=ephemeralConnectionId)
        topicHeader = PostHeader(PostFingerprint=topic.PostFingerprint)
        session.add(topic)
        session.add(topicHeader)
        session.commit()
        print "This topic has been added to the database"
        fingerprint = topic.PostFingerprint
        session.close()
        print('calling commit on protinstance')
        self.protInstance.commit()
        return fingerprint

    def createPost(self, postSubject, postText, parentFingerprint, ownerUsername, postLanguage):
        session = Session()
        post = Post(Subject=postSubject, Body=postText, OwnerUsername=ownerUsername,
                    ParentPostFingerprint=parentFingerprint, LocallyCreated=True, Language=postLanguage,
                    Upvoted=True, UpvoteCount=1)


        postHeader = PostHeader(PostFingerprint=post.PostFingerprint,
                                ParentPostFingerprint=post.ParentPostFingerprint, Language=postLanguage
                                )

        localnode = session.query(Node).filter(Node.LastConnectedIP=='LOCAL').one()
        vote = Vote(Direction=1,NodeId=localnode.NodeId, TargetPostFingerprint=post.PostFingerprint)

        session.add(post)
        session.add(postHeader)
        session.add(vote)
        session.commit()
        print "The post has been added to the database."
        if post.Subject == '':
            self.protInstance.commit(post.PostFingerprint) # This tells the dude I have committed a new post.
            # Only commit if it is a post. Since there is no way to create a subject without a post, this should
            # prevent double signals.
        fingerprint = post.PostFingerprint
        session.close()

        return fingerprint

    def createSignedPost(self, postSubject, postText, parentFingerprint, ownerFingerprint):
        session = Session()
        post = Post(Subject=postSubject, Body=postText, OwnerFingerprint=ownerFingerprint,
                    ownerUsername = getUser(ownerFingerprint).Username,
                    ParentPostFingerprint=parentFingerprint)
        session.add(post)
        session.commit()
        print "The signed post has been added to the database."
        fingerprint = post.PostFingerprint
        session.close()
        return fingerprint

    def writeUserProfile(self, userProfileInJSON):
        # This is the autocommitter.
        print('Autocommit Fired. Writing new user profile.')
        f = open(profiledir+'UserProfile/UserProfile.json', 'wb')
        f.write(userProfileInJSON.encode('utf8'))
        f.close()
        refreshBackendValuesGatheredFromJson(userProfileInJSON)
        return True

    def votePost(self, fingerprint, voteDirection):
        # Vote Direction is either 2, 1, 0 or -1. 2 means remove prior votes. (user clicks on the
        # button again to remove his / her vote)
        session = Session()
        post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
        localnode = session.query(Node).filter(Node.LastConnectedIP=='LOCAL').one()
        post.Dirty = True
        if voteDirection == 2: # If this is a remove vote call:

            if post.Upvoted == True:
                # And it's currently upvoted, remove effects,
                post.Upvoted = False
                post.UpvoteCount -= 1
                post.LastVoteDate = datetime.utcnow()
                session.add(post)
                # And the vote item itself.
                voteUp = session.query(Vote).filter(Vote.Direction == 1,
                                                    Vote.TargetPostFingerprint == fingerprint,
                                                    Vote.NodeId == localnode.NodeId)\
                                                    .one()
                session.delete(voteUp)
                session.commit()
                session.close()
                return True

            if post.Downvoted == True:
                # And it's currently downvoted, remove effects,
                post.Downvoted = False
                post.DownvoteCount -= 1
                post.LastVoteDate = datetime.utcnow()
                session.add(post)
                # And remove the item itself.
                voteDown = session.query(Vote).filter(Vote.Direction == 1,
                                                    Vote.TargetPostFingerprint == fingerprint,
                                                    Vote.NodeId == localnode.NodeId)\
                                                    .one()
                session.delete(voteDown)
                session.commit()
                session.close()
                return True

        else:
            # If this is not a remove vote call:
            if not session.query(Vote).filter(Vote.TargetPostFingerprint == fingerprint, Vote.NodeId == localnode.NodeId).all():
                # If empty (If first vote on this post by this user)
                vote = Vote(Direction=voteDirection,NodeId=localnode.NodeId, TargetPostFingerprint=fingerprint)
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
                vote = session.query(Vote).filter(Vote.TargetPostFingerprint == fingerprint, Vote.NodeId == localnode.NodeId)\
                    .one()
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
            session.close()
            return True

    def savePost(self, fingerprint):
        session = Session()
        post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()
        if post.Saved == True:
            post.Saved = False
        else:
            post.Saved = True
        session.add(post)
        session.commit()
        session.close()
        return True

    def markAllRepliesAsRead(self):
        session = Session()
        replies = session.query(Post).filter(Post.IsReply == True).all()
        for r in replies:
            r.IsReply = False
            session.add(r)

        session.commit()
        session.close()
        return True

    def markAllSavedsAsNotSaved(self):
        session = Session()
        saveds = session.query(Post).filter(Post.Saved == True).all()
        for s in saveds:
            s.Saved = False
            session.add(s)

        session.commit()
        session.close()
        return True

    def sendConnectWithIPSignal(self, ip, port):
        self.protInstance.connectWithIP(ip, port)