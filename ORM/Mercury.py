"""
    This is the deepest layer of the Networking stack. It handles the atomic actions.

    Everything in this API returns deferreds.

    All of the functions given are run in separate threads to not block the reactor. This also sometimes
    causes write failures due to locks and double commits. Beware.

    Everything in this API should return dicts. No SQLAlchemy objects.
"""
from __future__ import print_function
import ujson
from time import sleep
from datetime import timedelta

from twisted.internet import threads
from termcolor import cprint

from ORM.models import *
from Demeter import committer
import globals

if not globals.debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass

"""
    Handshake Methods.
"""

def handleConnectingNode(nodeId, ip, port):

    def threadFunction():
        s = Session()
        if not s.query(Node).filter(Node.NodeId == nodeId).count(): # If doesn't exist
            cprint('CONNECTED TO: A NEW NODE. Adding new record.', 'cyan', 'on_blue')
            node = Node(NodeId=nodeId, LastConnectedIP=ip, LastConnectedDate=datetime.utcnow(),
                    LastConnectedPort=port)
        else: # If node already exists in database
            cprint('CONNECTED TO: A KNOWN NODE. Updating record.', 'cyan', 'on_blue')
            node = s.query(Node).filter(Node.NodeId == nodeId).one()
            node.LastConnectedDate = datetime.utcnow()
            node.LastConnectedIP = ip
            node.LastConnectedPort = port
        node = node.asDict()
        s.close()
        return node

    return threads.deferToThread(threadFunction)

"""
    Outbound methods.
    Outbound means the connections this node is initiating. So these are mostly methods to
    insert stuff that is received from the remote. They write to the database.
"""

def insertHeadersAndVotes(reply, connectedNode):
    #Reply: PositiveHeaders, NeutralHeaders, NegativeHeaders, TopicHeaders
    #cprint('RECEIVED: HEADERS. \n%s' %(reply), 'cyan', 'on_blue')

    def threadFunction():
        # Just create everything as if new and put it to the queue. Queue, at commit time, will remove the already
        # existing ones.

        # print('This packet has %s positive, %s neutral, %s negative headers and %s topics.' % (len(reply['PositiveHeaders']),len(reply['NeutralHeaders']), len(reply['NegativeHeaders']), len(reply['TopicHeaders'])))

        del reply['TotalNumberOfPackets']
        del reply['CurrentPacketNo']

        for key in reply:

            if key is 'PositiveHeaders':
                voteDirection = 1
            elif key is 'NeutralHeaders':
                voteDirection = 0
            elif key is 'NegativeHeaders':
                voteDirection = -1
            elif key is 'TopicHeaders':
                voteDirection = None
            else:
                print(key)
                print(reply)
                raise Exception('AetherError: Vote direction could not be determined.')

            for h in reply[key]:
                h = ujson.loads(h)
                header = PostHeader(PostFingerprint=h['PostFingerprint'],
                        ParentPostFingerprint=h['ParentPostFingerprint'],
                        Language=h['Language'])

                if key is not 'TopicHeaders':
                    vote = Vote(Direction=voteDirection,
                            TargetPostFingerprint=h['PostFingerprint'], NodeId=connectedNode['NodeId'])
                    vote = vote.asDict()
                    committer.addVote(vote)
                header = header.asDict()
                committer.addHeader(header)

        return reply
    # Needs to be handled in the queue commit: dirty status transfer to existing entry in db.
    return threads.deferToThread(threadFunction)

def listNeededPosts(reply):
    def threadFunction():
        s = Session()
        neededPosts = []
        # Here, I will look at the replied list, and check whether all of them have their
        # corresponding posts. If any of them does not, I'll add it to the needed posts list.
        # I will also not ask for downvoted posts.

        # Review: This looks good enough to stand without a change.
        for key in reply:
            if key == 'NegativeHeaders':
                continue
            for h in reply[key]:
                h = ujson.loads(h)
                if not s.query(Post)\
                .filter(Post.PostFingerprint == h['PostFingerprint']).count():
                    neededPosts.append(h['PostFingerprint'])
        s.close()
        # print('Needed posts in this packet:', neededPosts)
        return neededPosts

    return threads.deferToThread(threadFunction)

def insertPost(reply):
    #cprint('RECEIVED: A POST. \n', 'cyan', 'on_blue')
    def threadFunction():
        # This feeds the commit queue posts with UNIX timestamps as dates.
        # Convert it to datetime objects.
        s = Session()
        postAsDict = ujson.loads(reply['Post'])
        if not s.query(Post).filter(Post.PostFingerprint == postAsDict['PostFingerprint']).count():
            # If there are no other copies (We know if it exists, but we do not know if it does not, because it
            # might be in queue. But the queue does check for that, so this is just a first line of defense against
            # duplication to relieve the committer.
            postAsDict['CreationDate'] = datetime.utcfromtimestamp(float(postAsDict['CreationDate']))
            postAsDict['LastVoteDate'] = datetime.utcfromtimestamp(float(postAsDict['LastVoteDate']))
            committer.addPost(postAsDict)
        s.close()

    # Needs to be handled in the queue commit: parency check.
    return threads.deferToThread(threadFunction)

def insertNodes(reply, localNodeId):

    def threadFunction():
        for n in reply:
            n = ujson.loads(n)
            if n['LastRetrievedDate']:
                n['LastRetrievedDate'] = datetime.utcfromtimestamp(float(n['LastRetrievedDate']))
            committer.addNode(n)

    return threads.deferToThread(threadFunction)



"""
    Inbound Methods. These are methods that are used to produce stuff this node serves to those who
    connect to it. These methods read from the database.
"""

def getHeaders(LastSyncTimestamp, Languages):
    """
        This function gets the headers that will be sent across the wire.
        This includes headers for neutrally-voted posts, positive-voted posts and topics.

        Negative votes also make post headers go through, but they do not make remote to call the
        post itself.
    """
    def threadFunction():
        s = Session()
        positiveHeadersArray = []
        neutralHeadersArray = []
        negativeHeadersArray = []
        topicHeadersArray = []

        if LastSyncTimestamp != 'NO-TIMESTAMP':
            positiveP = s.query(Post)\
                .filter(Post.Upvoted == True)\
                .filter(Post.LastVoteDate >=LastSyncTimestamp)

            neutralP = s.query(Post)\
                .filter(Post.Neutral == True)\
                .filter(Post.LastVoteDate >=LastSyncTimestamp)

            negativeP = s.query(Post)\
                .filter(Post.Downvoted == True)\
                .filter(Post.LastVoteDate >=LastSyncTimestamp)

            postsTopic = s.query(Post)\
                .filter(Post.ParentPostFingerprint == None)\
                .filter(Post.LastVoteDate >= LastSyncTimestamp)\
                .all()
            if Languages[0] != 'ALL':
                postsPositive = positiveP.filter(Post.Language.in_(Languages)).all()
                postsNeutral = neutralP.filter(Post.Language.in_(Languages)).all()
                postsNegative = negativeP.filter(Post.Language.in_(Languages)).all()
            else:
                postsPositive = positiveP.all()
                postsNeutral = neutralP.all()
                postsNegative = negativeP.all()

        else:
            postsPositive = s.query(Post)\
                .filter(Post.Upvoted == True)\
                .filter(Post.Language.in_(Languages))\
                .filter(Post.CreationDate >
                (datetime.utcnow() - timedelta(weeks=26))).all()
                # 26 Weeks = 6 months

            postsNeutral = s.query(Post)\
                .filter(Post.Neutral == True)\
                .filter(Post.Language.in_(Languages))\
                .filter(Post.CreationDate >
                (datetime.utcnow() - timedelta(weeks=26))).all()

            postsNegative = s.query(Post)\
                .filter(Post.Downvoted == True)\
                .filter(Post.Language.in_(Languages))\
                .filter(Post.CreationDate >
                (datetime.utcnow() - timedelta(weeks=26))).all()

            postsTopic = s.query(Post)\
                .filter(Post.ParentPostFingerprint == None)\
                .filter(Post.CreationDate >
                (datetime.utcnow() - timedelta(weeks=26))).all()

        # TODO: I should get these four loops below into one loop of loops when I finish figuring out this algorithm.

        for p in postsPositive:
            #try:
            #    s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first()
            #except:
            #    while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
            #        s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
            #        s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'POSITIVE'
            del h['ID']
            #h = ujson.dumps(h)
            positiveHeadersArray.append(h)

        for p in postsNeutral:
            #try:
            #    s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first()
            #except:
            #    while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
            #        s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
            #        s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'NEUTRAL'
            del h['ID']
            #h = ujson.dumps(h)
            neutralHeadersArray.append(h)

        for p in postsNegative:
            #try:
            #    s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first()
            #except:
            #    while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
            #        s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
            #        s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'NEGATIVE'
            del h['ID']
            #h = ujson.dumps(h)
            negativeHeadersArray.append(h)

        for p in postsTopic:
            #try:
            #    s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first()
            #except:
            #    while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
            #        s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
            #        s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'TOPIC'
            del h['ID']
            #h = ujson.dumps(h)
            topicHeadersArray.append(h)

        dbReply = positiveHeadersArray + neutralHeadersArray + negativeHeadersArray + topicHeadersArray
        return dbReply

    return threads.deferToThread(threadFunction)

def getPost(fingerprint):

    def threadFunction():
        s = Session()
        post = s.query(Post).filter(Post.PostFingerprint == fingerprint).one()
        postAsDict = post.asDict()
        del postAsDict['ID']
        del postAsDict['UpvoteCount']
        del postAsDict['DownvoteCount']
        del postAsDict['NeutralCount']
        del postAsDict['ReplyCount']
        del postAsDict['Upvoted']
        del postAsDict['Downvoted']
        del postAsDict['Neutral']
        del postAsDict['Saved']
        del postAsDict['IsReply']
        del postAsDict['LocallyCreated']
        del postAsDict['Dirty']
        del postAsDict['RankScore']


        cprint('FROM REMOTE: POST REQUEST: for post %s' %(fingerprint), 'white', 'on_yellow', attrs=['bold'])
        return postAsDict

    return threads.deferToThread(threadFunction)

def toJson(item):

    def threadFunction():
        itemAsJson = ujson.dumps(item)
        return itemAsJson

    return threads.deferToThread(threadFunction)

def getNodes(LastSyncTimestamp):

    def threadFunction():
        s = Session()
        if LastSyncTimestamp != 'NO-TIMESTAMP':
            retrievedNodes = s.query(Node)\
                .filter(Node.LastRetrievedDate >= LastSyncTimestamp)\
                .filter(Node.LastRetrievedDate >= timedelta(days=10)).all()
            connectedNodes = s.query(Node).filter(Node.LastConnectedDate >= LastSyncTimestamp)\
            .filter(Node.LastConnectedDate >= timedelta(days=10))\
            .all()
            nodes = retrievedNodes + connectedNodes
            if len(nodes) > 0:
                del nodes[0] # this is the local entry
        else:
            nodes = s.query(Node).all()
            del nodes[0] # this is the local entry

        return nodes

    return threads.deferToThread(threadFunction)

def processNodes(nodes):

    def threadFunction():
        processedNodes = []
        for n in nodes:
            nDict = n.asDict()
            del nDict['ID']
            processedNodes.append(nDict)

        return processedNodes

    return threads.deferToThread(threadFunction)
