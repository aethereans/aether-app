"""
    Demeter is the goddess of maintenance / sustenance.
    This is the interface slow-cycle maintenance / sustenance loop (Persephone) uses for maintenance
    / sustenance calls such as dirty checking for updates.

    This is a deferred API. Everything returns deferred and runs in threads. This API is used mostly by eventLoop.
"""

from __future__ import print_function
from twisted.internet import threads
from sqlalchemy.orm import joinedload
from datetime import datetime, timedelta
from termcolor import cprint
from sqlalchemy.orm import exc
import miniupnpc, random
import copy
from twisted.internet import defer
from sqlalchemy.orm import exc
from math import log
import calendar,time

from ORM.models import *
from globals import aetherListeningPort
import time
from InputOutput import interprocessAPI
import globals

if not globals.debugEnabled:
    def print(*a):
        pass
    def cprint(*a):
        pass


def timeit(method):

    def timed(*args, **kw):
        ts = time.time() *1000
        result = method(*args, **kw)
        te = time.time() *1000

        print('%r (%r, %r) %2.2f ms' % \
              (method.__name__, args, kw, te-ts))
        return result

    return timed



def checkUPNPStatus(delay):

    def threadFunction():
        u = miniupnpc.UPnP()
        u.discoverdelay = delay #milliseconds
        def mapUPNP():
            print('Checking UPNP Status... with %d second search...' %(delay/1000))
            numberOfFoundDevices = u.discover()
            print('I have found %s devices.' %numberOfFoundDevices)
            if numberOfFoundDevices > 1:
                print('Multiple UPNP devices found. Don\'t know what to do here. Adding mapping to only what selectigd() gives to me.')
                pass
            if numberOfFoundDevices == 0 :
                print('There are no UPNP devices in the network or this computer is directly connected to the Internet.')
                return False
            i = 0
            mappings = []
            try:
                u.selectigd()
                print('This machine has the IP of %s' %(u.externalipaddress()))
                while True:
                    p = u.getgenericportmapping(i)
                    if p==None:
                        break
                    else:
                        mappings.append(p)
                    i+= 1
                for m in mappings:
                    if m[3][:6] == 'Aether' and \
                            m[2][0] == u.lanaddr and \
                            m[2][1] == aetherListeningPort:
                        print('I have Aether already mapped to this computer at port %s.' % aetherListeningPort)
                        return True

                    elif m[3][:6] == 'Aether' and \
                            m[2][0] == u.lanaddr and \
                            m[2][1] != aetherListeningPort:
                        print('I have Aether already mapped to this computer at port %s, but it is not the port Aether is currently using.  Removing the old port and adding a mapping to %s.' %(m[2][1], aetherListeningPort))
                        u.deleteportmapping(m[2][1], 'TCP') # Delete the port mapping. It will add again in the end.

                import random
                nameString = 'Aether.' + str(random.randint(10000000,99999999))
                u.addportmapping(aetherListeningPort, 'TCP', u.lanaddr, aetherListeningPort, nameString, '')
                print('I have mapped TCP port %s on the router to the same port of this local machine.' % aetherListeningPort)
                return True

            except Exception, ex:
                print(ex)
                return False
        return mapUPNP()
        # I have no idea how long an UPNP mapping lasts.
        # Should I be removing the mapping request?
        # What happens two machines that are in the same network are using Aether? What's the chance of
        # both being assigned to the same port? I need to randomise the port.

    return threads.deferToThread(threadFunction)

# FIXME before package: timeout
def getNodesToConnect(amount=10, timeout=0):
    def threadFunction():
        s = Session()
        connecteds = s.query(Node).order_by(Node.LastConnectedDate.desc())\
            .filter(Node.LastConnectedDate < (datetime.utcnow() - timedelta(minutes = timeout)))\
            .filter(Node.LastConnectedIP != 'LOCAL')\
            .limit(amount*4)\
            .all()
        randConnectedIndexes = [int(amount*4*random.random()) for i in xrange(amount/2)] # get random x items from 0-25 range.
        finalConnectedArray = []
        if len(connecteds) < amount*4:
            finalConnectedArray = connecteds[:amount]
        else:
            for rNumber in randConnectedIndexes:
                finalConnectedArray.append(connecteds[rNumber])
        # This prevents connection to nodes connected to in the last [timeout] minutes.

        # Okay, retrieveds.. I need to increase the accuracy of retrieved selection.
        # retrieveds = s.query(Node).order_by(Node.LastRetrievedDate.desc())\
        #     .filter(Node.LastConnectedIP == None)\
        #     .limit(amount*4).all()
        # randRetrievedIndexes = [int(amount*4*random.random()) for i in xrange(amount/2)] # get random x items from 0-25 range.
        # finalRetrievedArray = []
        # if len(retrieveds) < amount*4:
        #     finalRetrievedArray = retrieveds[:amount]
        # else:
        #     for rNumber in randRetrievedIndexes:
        #         finalRetrievedArray.append(retrieveds[rNumber])
        retrieveds = s.query(Node).order_by(Node.LastRetrievedDate.desc())\
            .filter(Node.LastConnectedIP == None)\
            .limit(amount/2).all()
        finalRetrievedArray = []
        for n in retrieveds:
            finalRetrievedArray.append(n)
        s.close()
        # if (len(connecteds) + len(retrieveds) < amount):
        #     print('All Connected Nodes: %s' %connecteds)
        #     print('All Retrieved Nodes: %s' %retrieveds)
        #     return connecteds + retrieveds # If we don't have enough people, don't filter.
        # else:
        connectedIdList = []
        for n in finalConnectedArray:
            connectedIdList.append(n.NodeId)
        retrievedIdList = []
        for n in finalRetrievedArray:
            retrievedIdList.append(n.NodeId)

        print('Randomly Selected Connected Nodes: %s' %connectedIdList)
        print('Randomly Selected Retrieved Nodes: %s' %retrievedIdList)
        return finalConnectedArray + finalRetrievedArray

    return threads.deferToThread(threadFunction)


class GlobalCommitter(object):
    """
        This is a singleton that commits changes to the database on certain intervals. To not block addition of further
        elements into the queue while the batch is processing, when a commit fires, it will copy the contents of
        queues, and clear them. Then the commit will proceed with copied instances, while not blocking additions.

        This is the only place SQL definitions in models table should be used. The rest of the application should
        be free of concerns of data persistence and SQLAlchemy layer for writing. Reading is permitted across the board.

        There are many ways Aether writes to the database. Getting them together into one place will be tricky.

        Adding check functions into the addition functions: would that kill the performance? Well, It's not like
        I am doing bulk edits even if I did that in the commit, so yeah, it should be similar.
    """
    def __init__(self):
        self.postQueue = []
        self.voteQueue = []
        self.headerQueue = []
        self.nodeQueue = []

        self.sanitizedPostQueue = []
        self.sanitizedVoteQueue = []
        self.sanitizedHeaderQueue = []
        self.sanitizedNodeQueue = []
        self.lastCommit = datetime.now()
        self.commitInProgress = False

        self.thereAreReplies = False
        self.newPostsToIncrement = []

    def receiveInterprocessProtocolInstance(self, protInstance):
        self.interProt = protInstance

    def incrementReplyCounts(self, fingerprintsArray, session):
        print('Reply count insertion starts.')
        for fingerprint in fingerprintsArray:
            self.resolveReplyCount(fingerprint, session)
        print('Reply count insertion ends.')

    def resolveReplyCount(self, fingerprint, session):
        """
        Given a fingerprint, this method will traverse upwards in the ancestry chain and will +1 the reply count
        of each of the ancestors.
        """

        #### Find the post in question, first from the Q, then from the DB. Standard fare.

        found = False
        for pendingPost in self.sanitizedPostQueue:
            if pendingPost.PostFingerprint == fingerprint:
                found = True
                servedFromPending = True
                post = pendingPost # this is a post item
        if not found:
            if session.query(Post).filter(Post.PostFingerprint == fingerprint).count():
                found = True
                servedFromPending = False
                # If found in database:
                post = session.query(Post).filter(Post.PostFingerprint == fingerprint).one()

        if not found:
            # If does not exist at all.
            raise Exception('The post with fingerprint %s does not exist. End of the ancestry chain.' % fingerprint)
            # When

        #### At this point, I have an item from persistence called 'post' and I know the persistence type.

        assert(isinstance(post, Post))
        post.ReplyCount += 1
        if not servedFromPending:
            self.sanitizedPostQueue.append(post)
        try:
            return self.resolveReplyCount(post.ParentPostFingerprint, session)
        except:
            return True

    def rankPost(self, postToRank):
        # For now we are only ranking subjects. Post ranking will probably have a different algorithm.
        assert(isinstance(postToRank, Post))
        if postToRank.OwnerUsername == None and postToRank.OwnerFingerprint == None: #if topic
            # This will never hit because topics are never dirty nor marked as such.
            return 0.0
        elif postToRank.Subject != None: # if subject
            rawScore = postToRank.UpvoteCount - postToRank.DownvoteCount
            if not rawScore:
                return 0.0 # If raw score is 0, I don't need to calculate, it returns 0.
            order = log(max(abs(rawScore), 1), 2) # I'm using log 2, instead of log 10.
            sign = 1 if rawScore > 0 else -1 if rawScore < 0 else 0
            seconds = calendar.timegm(time.gmtime()) - 1352092272 # This is Mon, 05 Nov 2012 05:11:12 GMT
            return round(order + sign * seconds / float(65535), 7) # Float forces the result to be a float.
        else: #if post
            return 0.0

    @timeit
    def resolveVotesAndFlags(self, session):
        """
            This loop has three functions:
            1) It counts the votes a post has received.
            2) It flags and unflags neutral broadcast.
            3) It figures out whether a post is reply to a local post.

            This is the only part that is concerned with dirties.
            ---

            What can we assume in things that are available in maintenance? They are dirty flagged items, so they EXIST.

            THIS DEPENDS ON POSTS BEING APPROPRIATELY MARKED AS DIRTY.

        """

        print('Vote counting and neutral & reply flagging starts.')
        dirties = []
        # First, check if there are dirties in the database. This should only be the case for stuff the user adds.
        dirtiesFromDb = session.query(Post).filter(Post.Dirty == True).all()
        for dirtyPost in dirtiesFromDb:
            dirties.append(dirtyPost)
        # Then, look through the pending queue to grab dirtied posts.
        for post in self.sanitizedPostQueue:
            if post.Dirty:
                dirties.append(post)

        # DO NOT WRITE TO THE 'post' DIRECTLY. IT'S READ ONLY AND NOT SAVED. Pass it to updatePost for update.

        for post in dirties:
            # This part checks if this is a reply to a post user has written.
            replyFlag = False
            try:
                parent = session.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
            except:
                pass
            else:
                if parent.LocallyCreated:
                    replyFlag = True
                    self.thereAreReplies = True

            # Here I start counting the upvotes, downvotes and else. This is simpler, because it's flat.
            upvoteCount = session.query(Vote).filter(Vote.Direction == 1,
                                                          Vote.TargetPostFingerprint == post.PostFingerprint).count()
            downvoteCount = session.query(Vote).filter(Vote.Direction == -1,
                                                            Vote.TargetPostFingerprint == post.PostFingerprint).count()
            neutralCount = session.query(Vote).filter(Vote.Direction == 0,
                                                           Vote.TargetPostFingerprint == post.PostFingerprint).count()
            neutralFlag = False

            lastVoteDate = None
            # This part below checks for neutral broadcast.
            if upvoteCount + downvoteCount + neutralCount < 10  and post.ParentPostFingerprint is not None:
                # If all votes in total is lesser than 10, and if not a topic.
                # (Does this check for the case all are empty?)
                if not post.Upvoted and not post.Downvoted:
                    # If it's not upvoted or downvoted.
                    neutralFlag = True
                else:
                    # If it is actually upvoted or downvoted.
                    neutralFlag = False
                    #post.LastVoteDate = datetime.utcnow() #FIXME update
                    lastVoteDate = datetime.utcnow() #FIXME FIXED I think (update
            elif upvoteCount + downvoteCount + neutralCount > 10:
                # If the post actually got more than 10 votes.
                neutralFlag = False
                neutralCount -= 1 # Remove your own.

            # This part below calculates the rank of the post.
            postRank = self.rankPost(post)

            # And send the values over to the updater.
            self.updatePost(session, post.PostFingerprint, upvoteCount, downvoteCount, neutralCount,
                                                neutralFlag, replyFlag, postRank, lastVoteDate)

        print('Vote counting and neutral & reply flagging ends.')

    def updatePost(self, session, fingerprint, upvoteCount=None, downvoteCount=None, neutralCount=None,
               neutralFlag=None, replyFlag=False, postRank=None, lastVoteDate=None): # last vote date can be empty.

        assert(isinstance(upvoteCount, int))
        assert(isinstance(downvoteCount, int))
        assert(isinstance(downvoteCount, int))
        assert(isinstance(neutralFlag, bool))
        assert(isinstance(replyFlag, bool))
        assert(isinstance(postRank, float))
        # Last Vote Date can be empty if no new votes arrive.

        # Here I need to decide whether this post is a topic, subject or post.
        # topic = no rank, subject = rank, post = no rank ( for now )
        for pendingPost in self.sanitizedPostQueue:
            if pendingPost.PostFingerprint == fingerprint:
                # If exists in sanitized queue, update values.
                pendingPost.UpvoteCount = upvoteCount
                pendingPost.DownvoteCount = downvoteCount
                pendingPost.NeutralCount = neutralCount
                pendingPost.Neutral = neutralFlag
                pendingPost.IsReply = replyFlag
                if lastVoteDate:
                    pendingPost.LastVoteDate = lastVoteDate
                if pendingPost.Subject: # If this has a subject.
                    pendingPost.RankScore = postRank
                pendingPost.Dirty = False
                return
        if session.query(Post).filter(Post.PostFingerprint==fingerprint).count():
            # If exists in database,
            post = session.query(Post).filter(Post.PostFingerprint==fingerprint).one()
            post.UpvoteCount = upvoteCount
            post.DownvoteCount = downvoteCount
            post.NeutralCount = neutralCount
            post.Neutral = neutralFlag
            post.IsReply = replyFlag
            if lastVoteDate:
                    post.LastVoteDate = lastVoteDate
            if post.Subject: # If this has a subject.
                post.RankScore = postRank
            post.Dirty = False
            self.sanitizedPostQueue.append(post)
            return
        raise Exception('Dirty item with fingerprint %s found neither in queue nor in the database.' % fingerprint)


    def dirtyVoteTargetPosts(self, session):
        """
        This marks target posts of all newly arrived votes dirty.
        This happens before rank assignment. If a corresponding post exists in the post queue, mark as dirty.
        If a corresponding post exists in the database, mark as dirty and add to the queue.
        I don't need to discern whether there is a rank or not. If something points to that post, mark as dirty.
        The recount is idempotent, so there should be no problems.
        """
        for vote in self.sanitizedVoteQueue:
            found = False
            for pendingPost in self.sanitizedPostQueue:
                if pendingPost.PostFingerprint == vote.TargetPostFingerprint:
                    found = True
                    pendingPost.Dirty = True
                    return
            if not found:
                if session.query(Post).filter(Post.PostFingerprint == vote.TargetPostFingerprint).count():
                    found = True
                    # If found in database:
                    post = session.query(Post).filter(Post.PostFingerprint == vote.TargetPostFingerprint).one()
                    post.Dirty = True
                    self.sanitizedPostQueue.append(post)
                    return
            if not found:
                # Do nothing, because when the post is added in the future, it will start as dirty, triggering a
                # count on itself without me doing anything.
                pass



    def commit(self):
        @timeit
        def threadFunction():
            self.commitInProgress = True
            self.lastCommit = datetime.now()
            print('Commit starts')
            session = Session()
            # Create copies of queues.
            postsCopy = copy.deepcopy(self.postQueue)
            votesCopy = copy.deepcopy(self.voteQueue)
            headersCopy = copy.deepcopy(self.headerQueue)
            nodesCopy = copy.deepcopy(self.nodeQueue)
            #print('In this commit: posts, votes, headers, nodes:', postsCopy, votesCopy, headersCopy, nodesCopy)
            # Then empty them so they can continue to accept further input.
            self.postQueue = []
            self.voteQueue = []
            self.headerQueue = []
            self.nodeQueue = []
            # Produce sanitized queues.

            for post in postsCopy:
                self.addSanitizedPost(post)
            for header in headersCopy:
                self.addSanitizedHeader(header)
            for node in nodesCopy:
                self.addSanitizedNode(node)
            for vote in votesCopy:
                self.addSanitizedVote(vote)

            # Modifications to queues.

            self.dirtyVoteTargetPosts(session) # This will modify sanitized post queue, to mark
            # new votes' posts as dirty.
            self.resolveVotesAndFlags(session) # This will modify sanitized post queue, to update the
            # counters on dirtied posts and attach a postrank.
            # This will also move all locally created posts to pending post queue.
            self.incrementReplyCounts(self.newPostsToIncrement, session) # This will modify sanitized post queue to update reply counter.
            # I am also adding myself to this newPostsToIncrement on commit from gui. So, they are checked now.

            # At this point, all edits are finished. Moving to adding them to the session.

            session.add_all(self.sanitizedPostQueue)
            session.add_all(self.sanitizedHeaderQueue)
            session.add_all(self.sanitizedNodeQueue)
            session.add_all(self.sanitizedVoteQueue)


            # Commit all of those into the database.
            session.commit()
            # Declare values.
            # print('Commit completed.', len(self.sanitizedHeaderQueue), 'headers,',
            #       len(self.sanitizedPostQueue),'posts,', len(self.sanitizedVoteQueue), 'votes, and',
            #       len(self.sanitizedNodeQueue), 'nodes are committed.')
            print('Commit completed.', len(self.sanitizedHeaderQueue), 'headers,', len(self.sanitizedPostQueue),'posts,', len(self.sanitizedVoteQueue), 'votes, and', len(self.sanitizedNodeQueue), 'nodes are committed.')
            # And clear the copied values.
            postsCopy = []
            votesCopy = []
            headersCopy = []
            nodesCopy = []
            # And clear the sanitized queues.
            self.sanitizedHeaderQueue = []
            self.sanitizedPostQueue = []
            self.sanitizedVoteQueue = []
            self.sanitizedNodeQueue = []
            # And clear 'new posts to increment'.
            self.newPostsToIncrement = []
            # If we have a reply at hand, send a signal.
            if self.thereAreReplies:
                self.interProt.callRemote(interprocessAPI.thereAreReplies)
                self.thereAreReplies = False
            # Finally, call the maintenance loop.
            session.close()
            self.commitInProgress = False


        return threads.deferToThread(threadFunction)

    def addPost(self, postAsDict):
        assert isinstance(postAsDict, dict)
        self.postQueue.append(postAsDict)

    def addHeader(self, headerAsDict):
        assert isinstance(headerAsDict, dict)
        self.headerQueue.append(headerAsDict)

    def addVote(self, voteAsDict):
        assert isinstance(voteAsDict, dict)
        self.voteQueue.append(voteAsDict)

    def addNode(self, nodeAsDict, firstHand=False):
        assert isinstance(nodeAsDict, dict)
        if firstHand:
            nodeAsDict['FIRSTHAND'] = True
        self.nodeQueue.append(nodeAsDict)


    def addSanitizedHeader(self, headerAsDict):
        """
            If there is no header, create a header. If there is already one, skip.
            Updateable parts:
            : No updateable parts. Cool.
        """
        session = Session()
        found = False
        # Check the queue to see whether it exists.
        for pendingHeader in self.sanitizedHeaderQueue:
            if headerAsDict['PostFingerprint'] == pendingHeader.PostFingerprint:
                found = True
                session.close()
                return
        # Then check the database to see whether it exists.
        if not found:
            if session.query(PostHeader)\
            .filter(PostHeader.PostFingerprint == headerAsDict['PostFingerprint']).count():
                found = True
                session.close()
                return
        # If it exists in neither, convert it to SQLAlchemy object and add to list.
        if not found:
            header = PostHeader(PostFingerprint=headerAsDict['PostFingerprint'],
                    ParentPostFingerprint=headerAsDict['ParentPostFingerprint'],
                    Language=headerAsDict['Language'])
            session.close()
            self.sanitizedHeaderQueue.append(header)
            return


    def addSanitizedVote(self, voteAsDict):
        """
            If there is no vote from that node, add. If there is,
            override fields (So vote direction changes are reflected.)
            Updateable parts:
            : Direction
        """

        assert(isinstance(voteAsDict, dict))
        session = Session()
        found = False
        for pendingVote in self.sanitizedVoteQueue:
            if voteAsDict['TargetPostFingerprint'] == pendingVote.TargetPostFingerprint and voteAsDict['NodeId'] == pendingVote.NodeId:
                # If found in the queue:
                found = True
                print('VOTE found in QUEUE')
                # Change the entry. It is already in the queue.
                if pendingVote.Direction != voteAsDict['Direction']:
                    print('Vote for post %s flipped from %d to %d' % (pendingVote.TargetPostFingerprint, pendingVote.Direction, voteAsDict['Direction']))
                pendingVote.Direction = voteAsDict['Direction']
                session.close()
                return
        if not found:
            if session.query(Vote)\
                .filter(Vote.TargetPostFingerprint == voteAsDict['TargetPostFingerprint'],
                Vote.NodeId == voteAsDict['NodeId']).count():
                # If found in the database:
                print('VOTE found in DATABASE')
                found = True
                # Change the entry and add it to queue.
                vote = session.query(Vote)\
                .filter(Vote.TargetPostFingerprint == voteAsDict['TargetPostFingerprint'],
                Vote.NodeId == voteAsDict['NodeId']).one()
                if vote.Direction != voteAsDict['Direction']:
                    print('Vote for post %s flipped from %d to %d' % (vote.TargetPostFingerprint, vote.Direction, voteAsDict['Direction']))
                vote.Direction = voteAsDict['Direction']
                session.close()
                assert(isinstance(vote, Vote))
                self.sanitizedVoteQueue.append(vote)
                return
        if not found:
            vote = Vote(Direction=voteAsDict['Direction'],
                            TargetPostFingerprint=voteAsDict['TargetPostFingerprint'],
                            NodeId=voteAsDict['NodeId'])
            session.close()
            assert(isinstance(vote, Vote))
            self.sanitizedVoteQueue.append(vote)
            return

    def addSanitizedPost(self, postAsDict):
        """
            If isn't in DB/Queue, add. If it is, don't touch.
        """
        assert(isinstance(postAsDict, dict))
        session = Session()
        found = False
        for pendingPost in self.sanitizedPostQueue:
            if pendingPost.PostFingerprint == postAsDict['PostFingerprint']:
                found = True
                session.close()
                return
        if not found:
            if session.query(Post).filter(Post.PostFingerprint == postAsDict['PostFingerprint']).count():
                found = True
                # If found in database:
                session.close()
                return
        if not found:
            # If does not exist at all.

            # Even if the remote mistakenly leaks the information, the local actively refuses to save any of it.
            if 'ID' in postAsDict: del postAsDict['ID']
            if 'UpvoteCount' in postAsDict: del postAsDict['UpvoteCount']
            if 'DownvoteCount' in postAsDict: del postAsDict['DownvoteCount']
            if 'NeutralCount' in postAsDict: del postAsDict['NeutralCount']
            if 'ReplyCount' in postAsDict: del postAsDict['ReplyCount']
            if 'Upvoted' in postAsDict: del postAsDict['Upvoted']
            if 'Downvoted' in postAsDict: del postAsDict['Downvoted']
            if 'Neutral' in postAsDict: del postAsDict['Neutral']
            if 'Saved' in postAsDict: del postAsDict['Saved']
            if 'IsReply' in postAsDict: del postAsDict['IsReply']
            if 'LocallyCreated' in postAsDict: del postAsDict['LocallyCreated']
            if 'Dirty' in postAsDict: del postAsDict['Dirty']
            if 'RankScore' in postAsDict: del postAsDict['RankScore']
            post = Post()
            for key in postAsDict:
                    if postAsDict[key] is not None and hasattr(post, key):
                        # second part of top: a guard to make sure the input is legit.
                        setattr(post, key, postAsDict[key])

            post.Neutral=True
            post.NeutralCount=1
            session.close()
            assert isinstance(post, Post)
            self.sanitizedPostQueue.append(post)
            self.newPostsToIncrement.append(post.PostFingerprint)
            return


    def addSanitizedNode(self, nodeAsDict):
        """
            If it isn't in the DB/Queue, add. If it is, change.
            Updateable parts:
            if firsthand:
                : last connected ip
                : last connected port
                : timestamp
            else:
                : last retrieved ip
                : last retrieved port

        """
        assert(isinstance(nodeAsDict, dict))
        if not nodeAsDict['NodeId']: # To prevent arrival of nodes with empty nodeid.
            return
        session = Session()
        found = False
        for pendingNode in self.sanitizedNodeQueue:
            if nodeAsDict['NodeId'] == pendingNode.NodeId:
                # If found in Queue:
                found = True
                # I think this is the place where most problems happen. TODO.
                if 'FIRSTHAND' in nodeAsDict:
                    print('NODE UPDATE from QUEUE @ FIRSTHAND. Old timestamp is %s, New timestamp is %s for node %s' % (pendingNode.LastSyncTimestamp, nodeAsDict['LastSyncTimestamp'], nodeAsDict['NodeId']))
                    # If this is the directly connected node (the remote).
                    if 'LastConnectedIP' in nodeAsDict and nodeAsDict['LastConnectedIP'] is not None:
                        pendingNode.LastConnectedIP = nodeAsDict['LastConnectedIP']
                    if 'LastConnectedPort' in nodeAsDict and nodeAsDict['LastConnectedPort'] is not None:
                        pendingNode.LastConnectedPort = nodeAsDict['LastConnectedPort']
                    if 'LastConnectedDate' in nodeAsDict and nodeAsDict['LastConnectedDate'] is not None:
                        pendingNode.LastConnectedDate = nodeAsDict['LastConnectedDate']
                    if 'LastSyncTimestamp' in nodeAsDict and nodeAsDict['LastSyncTimestamp'] is not None:
                        pendingNode.LastSyncTimestamp = nodeAsDict['LastSyncTimestamp']

                else:
                    # If this is a node arrived over the network.
                    if 'LastRetrievedIP' in nodeAsDict and nodeAsDict['LastRetrievedIP'] is not None:
                        pendingNode.LastRetrievedIP = nodeAsDict['LastRetrievedIP']
                    if 'LastRetrievedPort' in nodeAsDict and nodeAsDict['LastRetrievedPort'] is not None:
                        pendingNode.LastRetrievedPort = nodeAsDict['LastRetrievedPort']
                    if 'LastRetrievedDate' in nodeAsDict and nodeAsDict['LastRetrievedDate'] is not None:
                        pendingNode.LastRetrievedDate = nodeAsDict['LastRetrievedDate']
                    session.close()
                    return
        if not found:
            if session.query(Node).filter(Node.NodeId == nodeAsDict['NodeId']).count():
                # If found in Database:
                found = True

                node = session.query(Node).filter(Node.NodeId == nodeAsDict['NodeId']).one()
                if 'FIRSTHAND' in nodeAsDict:
                    print('NODE UPDATE from DATABASE @ FIRSTHAND.  Old timestamp is %s, New timestamp is %s for node %s' % (node.LastSyncTimestamp, nodeAsDict['LastSyncTimestamp'], nodeAsDict['NodeId']))
                    if 'LastConnectedIP' in nodeAsDict and nodeAsDict['LastConnectedIP'] is not None:
                        node.LastConnectedIP = nodeAsDict['LastConnectedIP']
                    if 'LastConnectedPort' in nodeAsDict and nodeAsDict['LastConnectedPort'] is not None:
                        node.LastConnectedPort = nodeAsDict['LastConnectedPort']
                    if 'LastConnectedDate' in nodeAsDict and nodeAsDict['LastConnectedDate'] is not None:
                        node.LastConnectedDate = nodeAsDict['LastConnectedDate']
                    if 'LastSyncTimestamp' in nodeAsDict and nodeAsDict['LastSyncTimestamp'] is not None:
                        node.LastSyncTimestamp = nodeAsDict['LastSyncTimestamp']
                else:
                    if 'LastRetrievedIP' in nodeAsDict and nodeAsDict['LastRetrievedIP'] is not None:
                        node.LastRetrievedIP = nodeAsDict['LastRetrievedIP']
                    if 'LastRetrievedPort' in nodeAsDict and nodeAsDict['LastRetrievedPort'] is not None:
                        node.LastRetrievedPort = nodeAsDict['LastRetrievedPort']
                    if 'LastRetrievedDate' in nodeAsDict and nodeAsDict['LastRetrievedDate'] is not None:
                        node.LastRetrievedDate = nodeAsDict['LastRetrievedDate']
                session.close()
                assert(isinstance(node, Node))
                self.sanitizedNodeQueue.append(node)
                return
        if not found:
            # If not found:
            if 'FIRSTHAND' in nodeAsDict:
                print('NODE CREATE @ FIRSTHAND. New timestamp is %s for node %s' % (nodeAsDict['LastSyncTimestamp'], nodeAsDict['NodeId']))
                node = Node()
                node.NodeId = nodeAsDict['NodeId']
                node.LastConnectedIP = nodeAsDict['LastConnectedIP']
                node.LastConnectedPort = nodeAsDict['LastConnectedPort']
                node.LastConnectedDate = nodeAsDict['LastConnectedDate']
                node.LastSyncTimestamp = nodeAsDict['LastSyncTimestamp']
                print('the timestamp of node about to be committed: %s' % node.LastSyncTimestamp)
            else:
                node = Node()
                node.NodeId = nodeAsDict['NodeId']
                node.LastRetrievedIP = nodeAsDict['LastRetrievedIP']
                node.LastRetrievedPort = nodeAsDict['LastRetrievedPort']
                node.LastRetrievedDate = nodeAsDict['LastRetrievedDate']
            session.close()
            assert(isinstance(node, Node))
            self.sanitizedNodeQueue.append(node)

committer = GlobalCommitter()

def connectToLastConnected(amount=10):
    from InputOutput.aetherProtocol import connectWithNode
    session = Session()
    lastConnectedNodes = session.query(Node).order_by(Node.LastConnectedDate.desc())\
            .filter(Node.LastConnectedIP != 'LOCAL')\
            .limit(amount)\
            .all()

    for node in lastConnectedNodes:
        node = node.asDict()
        connectWithNode(node)
