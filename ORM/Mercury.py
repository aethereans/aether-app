"""
    This is the deepest layer of the Networking stack. It handles the atomic actions.

    Everything in this API returns deferreds.

    All of the functions given are run in separate threads to not block the reactor. This also sometimes
    causes write failures due to locks and double commits. Beware.
"""
from __future__ import print_function
import ujson
from time import sleep
from datetime import timedelta

from twisted.internet import threads
from termcolor import cprint
from sqlalchemy import exc

from ORM.models import *
from globals import maximumRetryCount, retryWaitTime


"""
    Handshake Methods.
"""

def checkIfNodeExists(nodeId):
    def threadFunction():
        s = Session()
        existence= False
        if s.query(Node).filter(Node.NodeId == nodeId).count():
            existence = True
        s.close()
        return existence

    return threads.deferToThread(threadFunction)

def updateAndGetNode(nodeId, ip, port): # This is used by handshake, so this is used for direct
# connections. So this marks nodes as last connected.
    def threadFunction():
        cprint('CONNECTED TO: A KNOWN NODE. UPDATING RECORD IN DB.', 'cyan', 'on_blue')
        s = Session()
        node = s.query(Node).filter(Node.NodeId == nodeId).one()
        node.LastConnectedIP = ip
        node.LastConnectedPort = port
        node.LastConnectedDate = datetime.utcnow()
        nodeAsDict = node.asDict()
        s.add(node)
        s.commit()
        s.close()
        return nodeAsDict

    return threads.deferToThread(threadFunction)

def createNode(nodeId, ip, port):
    def threadFunction():
        cprint('CONNECTED TO: A NEW NODE. ADDING TO DB.', 'cyan', 'on_blue')
        s = Session()
        node = Node(NodeId=nodeId, LastConnectedIP=ip, LastConnectedDate=datetime.utcnow(),
                    LastConnectedPort=port)
        nodeAsDict = node.asDict()
        s.add(node)
        s.commit()
        s.close()
        return nodeAsDict

    return threads.deferToThread(threadFunction)

"""
    Outbound methods.
    Outbound means the connections this node is initiating. So these are mostly methods to
    insert stuff that is received from the remote. They write to the database.
"""

def insertHeaders(reply): #Reply: PositiveHeaders, NeutralHeaders, NegativeHeaders, TopicHeaders
    cprint('RECEIVED: HEADERS. \n%s' %(reply), 'cyan', 'on_blue')
    print('\n')
    def threadFunction():
        s = Session()
        # Checks all headers to see if they exist. If they do not, creates.
        del reply['TotalNumberOfPackets']
        del reply['CurrentPacketNo']
        for key in reply:
            for h in reply[key]:
                h = ujson.loads(h)
                retryCount = 0
                while True:
                    try:
                        if retryCount >= maximumRetryCount:
                            cprint('PERMFAILURE: Write of HEADER failed for %s times. '
                                   'Aborting write for '
                                   %(retryCount) + h['PostFingerprint'], 'yellow', 'on_red',
                                   attrs=['bold', 'dark', 'underline'])
                            break
                        # Core Part below
                        if not s.query(PostHeader)\
                        .filter(PostHeader.PostFingerprint == h['PostFingerprint'])\
                        .count():
                            header = PostHeader(PostFingerprint=h['PostFingerprint'],
                                    ParentPostFingerprint=h['ParentPostFingerprint'],
                                    Language=h['Language'])
                            s.add(header)
                        # Core Part ends
                        try:
                            cprint('PROCESSED: HEADER, CREATING NEW: \n%s' %(header.PostFingerprint), 'cyan', 'on_blue')
                        except:
                            cprint('PROCESSED: HEADER, ALREADY EXISTS: \n%s' %(h['PostFingerprint']), 'cyan',
                                   'on_blue')
                        print('\n')
                        break
                    except exc.OperationalError:
                        retryCount += 1
                        cprint('FAILURE: Write of HEADER %s. Will retry in %s seconds.'
                               'This is %s. try.'
                               %(h['PostFingerprint'], retryWaitTime, retryCount), 'yellow')

                        s.rollback()
                        sleep(retryWaitTime)
                        continue
        s.commit()
        s.close()
        return reply

    return threads.deferToThread(threadFunction)

def insertVotes(reply, node):
    def threadFunction():
        s = Session()
        # At this point I have all the headers in place. I'll go through headers and add them
        # appropriate votes according to the list they are in.
        for key in reply:
            if key != 'TopicHeaders':
                direction = 0
                if key == 'PositiveHeaders': direction = 1
                elif key == 'NeutralHeaders': direction = 0
                elif key == 'NegativeHeaders': direction = -1
                for h in reply[key]:
                    h = ujson.loads(h)
                    retryCount = 0
                    while True:
                        try:
                            if retryCount >= maximumRetryCount:
                                cprint('PERMFAILURE: Write of VOTE failed for %s times. '
                                       'Aborting write for vote '
                                       %(retryCount) + h['PostFingerprint'], 'yellow', 'on_red',
                                       attrs=['bold', 'dark', 'underline'])
                                break

                            # Core Part below
                            try:
                                s.query(PostHeader)\
                                .filter(PostHeader.PostFingerprint == h['PostFingerprint']).one() # If has multiple records.
                            except:
                                while s.query(PostHeader)\
                                .filter(PostHeader.PostFingerprint == h['PostFingerprint']).count() > 1:
                                    s.delete(s.query(PostHeader)\
                                    .filter(PostHeader.PostFingerprint == h['PostFingerprint']).first())
                                    s.commit()

                            header = s.query(PostHeader)\
                            .filter(PostHeader.PostFingerprint == h['PostFingerprint']).one()
                            header.Dirty = True # Adding a vote to a header or changing an existing vote marks header
                                                # as dirty.
                            s.add(header)
                            try:
                                s.query(Node).filter(Node.NodeId == node['NodeId']).one()
                            except:
                                while s.query(Node).filter(Node.NodeId == node['NodeId']).count() > 1:
                                    nodeToDelete = s.query(Post).filter(Post.PostFingerprint == h.PostFingerprint).first()
                                    for vote in nodeToDelete.Votes:
                                        s.delete(vote)
                                    s.delete(nodeToDelete)
                                    s.commit()
                            n = s.query(Node).filter(Node.NodeId == node['NodeId']).one()

                            if not s.query(Vote)\
                            .filter(Vote.postheader_id==header.ID)\
                            .filter(Vote.node_id==n.ID)\
                            .count():
                                # If the vote does not exist:
                                v = Vote(Direction=direction, PostHeader=header, Node=n)
                                cprint('PROCESSED: VOTE, CREATING NEW: \n%s for %s'
                                       %(v.Direction, header.PostFingerprint),
                                       'cyan', 'on_blue')
                                s.add(v)
                            else:
                                # If vote exists:
                                try:
                                    v = s.query(Vote)\
                                        .filter(Vote.postheader_id==header.ID)\
                                        .filter(Vote.node_id==n.ID)\
                                        .one()
                                        # This is where it crashes if remote sends two votes.  What should I do in case of
                                        # two votes? I should probably silently abort the process and log it to the console.

                                    v.Direction = direction
                                    cprint('PROCESSED: VOTE, ALREADY EXISTS: \n%s for %s'
                                           %(v.Direction, v.PostHeader.PostFingerprint), 'cyan', 'on_blue')
                                    s.add(v)
                                except Exception, ex:
                                    cprint('EXCEPTION: Write of VOTE failed because: %s\n' %(ex) +
                                        'If this says multiple rows found, remote tried to send two votes for same item.'
                                        'Only the 1st is committed to the database.'
                                        'Failed Write Vote Direction: %s, Vote for header: %s' %(direction, header.PostFingerprint),
                                        'yellow', 'on_red', attrs=['bold', 'dark', 'underline'])

                            print('\n')
                            # Core Part ends

                            break
                        except exc.OperationalError:
                            retryCount += 1
                            cprint('FAILURE: Write of VOTE for %s. Will retry in %s seconds.'
                                   'This is %s. try.'
                                   %(h['PostFingerprint'], retryWaitTime, retryCount), 'yellow')
                            s.rollback()
                            sleep(retryWaitTime)
                            continue
        s.commit()
        s.close()
        return reply

    return threads.deferToThread(threadFunction)

def listNeededPosts(reply):
    def threadFunction():
        s = Session()
        neededPosts = []
        # Here, I will look at the replied list, and check whether all of them have their
        # corresponding posts. If any of them does not, I'll add it to the needed posts list.
        # I will also not ask for downvoted posts.
        for key in reply:
            if key == 'NegativeHeaders':
                continue
            for h in reply[key]:
                h = ujson.loads(h)
                if not s.query(Post)\
                .filter(Post.PostFingerprint == h['PostFingerprint']).count():
                    neededPosts.append(h['PostFingerprint'])
        return neededPosts

    return threads.deferToThread(threadFunction)

def insertPost(reply, ephemeralConnectionId):
    cprint('RECEIVED: A POST. \n', 'cyan', 'on_blue')
    print('\n')
    def threadFunction():
        s = Session()
        #print(postAsDict['PostFingerprint'])
        #print(postAsDict['CreationDate'])
        postAsDict = ujson.loads(reply['Post'])
        retryCount = 0
        while True:
            try:
                if retryCount >= maximumRetryCount:
                    cprint('PERMANENT FAILURE: Write of POST failed for %s times. '
                           'Aborting write for post '
                           %(retryCount) + postAsDict['PostFingerprint'], 'red', 'on_red',
                           attrs=['bold', 'dark', 'underline'])
                    break
                # Core Part below
                if not s.query(Post).filter(Post.PostFingerprint == postAsDict['PostFingerprint']).count():
                    post = Post(PostFingerprint=postAsDict['PostFingerprint'],
                                Subject=postAsDict['Subject'],
                                Body=postAsDict['Body'],
                                OwnerFingerprint=postAsDict['OwnerFingerprint'],
                                OwnerUsername=postAsDict['OwnerUsername'],
                                CreationDate=datetime.utcfromtimestamp(float(postAsDict['CreationDate'])),
                                ParentPostFingerprint=postAsDict['ParentPostFingerprint'],
                                ProtocolVersion=postAsDict['ProtocolVersion'],
                                Language=postAsDict['Language'],
                                EphemeralConnectionId=ephemeralConnectionId,
                                Neutral=True, # These two will be set by persephone anyway, but setting these here
                                NeutralCount=1 # ensures if a sync happens before persephone has a chance, they won't be
                                # blackholed because of the sync timestamp.
                                )
                    #if not post.ParentPostFingerprint == None: # This is for server.
                    #   post.Neutral = True
                    try:
                        # Because we might not have the parent at all, in case of a topic arriving.
                        parent = s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
                        if (parent.LocallyCreated == True):
                            post.IsReply = True

                    except Exception, e:
                        print(e)
                    s.add(post)
                    s.commit()
                try:
                    cprint('PROCESSED: POST, CREATING NEW: \n%s' %(post.PostFingerprint), 'cyan', 'on_blue')
                except:
                    cprint('PROCESSED: POST, ALREADY EXISTS: \n%s' %(postAsDict['PostFingerprint']), 'cyan',
                               'on_blue')
                print('\n')
                s.close()
                # Core part ends
                break
            except exc.OperationalError:
                retryCount += 1
                cprint('TEMPORARY FAILURE: Write of POST %s. Will retry in %s seconds.'
                       'This is %s. try.'
                       %(postAsDict['PostFingerprint'], retryWaitTime, retryCount), 'yellow')
                s.rollback()
                sleep(retryWaitTime)
                continue
    return threads.deferToThread(threadFunction)

def insertNodes(reply, localNodeId):
    cprint('RECEIVED: NODES. \n%s' %(reply), 'cyan', 'on_blue')
    print('\n')
    def threadFunction():
        s = Session()
        for n in reply:
            retryCount = 0
            while True:
                try:
                    if retryCount >= maximumRetryCount:
                        cprint('PERMANENT FAILURE: Write of NODE failed for %s times. '
                               'Aborting write for node '
                               %(retryCount) + n['NodeId'], 'yellow', 'on_red',
                               attrs=['bold', 'dark', 'underline'])
                        break
                    # Core Part below
                    n = ujson.loads(n)

                    if s.query(Node).filter(Node.NodeId==n['NodeId']).count(): # if node exists
                        node = s.query(Node).filter(Node.NodeId == n['NodeId']).one()
                        if node.NodeId != localNodeId:
                            node.LastRetrievedIP = \
                                n['LastConnectedIP'] if n['LastConnectedIP'] != None else n['LastRetrievedIP']
                            node.LastRetrievedPort = \
                                n['LastConnectedPort'] if n['LastConnectedPort'] != None else n['LastRetrievedPort']
                            node.LastRetrievedDate = datetime.utcfromtimestamp(n['LastConnectedDate']
                                if n['LastConnectedDate'] != None else n['LastRetrievedDate'])
                            cprint('PROCESSED: NODE, ALREADY EXISTS: \n%s' %(node.NodeId), 'cyan', 'on_blue')
                        else:
                            cprint('PROCESSED: NODE, SKIPPING UPDATE, This is the local node.: \n%s' %(node.NodeId),
                                   'white', 'on_green')
                    else:
                        node = Node(
                        LastRetrievedIP =
                                    n['LastConnectedIP'] if n['LastConnectedIP'] != None else n['LastRetrievedIP'],
                        LastRetrievedPort =
                                    n['LastConnectedPort'] if n['LastConnectedPort'] != None else n['LastRetrievedPort'],
                        LastRetrievedDate =
                                    datetime.utcfromtimestamp(float(n['LastConnectedDate'])
                                    if n['LastConnectedDate'] != None else float(n['LastRetrievedDate'])),
                        NodeId = n['NodeId']
                                    )
                        cprint('PROCESSED: NODE, CREATING NEW: \n%s' %(node.NodeId), 'cyan', 'on_blue')
                    print('\n')
                    s.add(node)
                    # Core Part ends
                    break
                except exc.OperationalError:
                    retryCount += 1
                    cprint('TEMPORARY FAILURE: Write of NODE for %s. Will retry in %s seconds.'
                           'This is %s. try.'
                           %(n['NodeId'], retryWaitTime, retryCount), 'yellow')
                    s.rollback()
                    sleep(retryWaitTime)
                    continue

        s.commit()
        s.close()
    return threads.deferToThread(threadFunction)

def insertNewSyncTimestamp(node, newSyncTimestamp):
    def threadFunction():
        s = Session()
        #try:
        #    s.query(Node).filter(Node.NodeId == node['NodeId']).one()
        #    # Check if there is only one instance of the same node. This is a little problematic actually.. Votes are
        #    # foreign keyed to the nodes. I think duplicate node only happens when two connections to the same node
        #    # fires at the same moment, adding it twice to the same database.
        #    # I need to retain database integrity, so I need to remove the votes from that vote, which are hard foreign
        #    # keyed to the node in the table.
        #except:
        #    while s.query(Node).filter(Node.NodeId == node['NodeId']).count() > 1:
        #        nodeToDelete = s.query(Node).filter(Node.NodeId == node['NodeId']).first()
        #        for vote in nodeToDelete.Votes:
        #            s.delete(vote)
        #        s.delete(nodeToDelete)
        #        s.commit()
        n = s.query(Node).filter(Node.NodeId == node['NodeId']).one()
        n.LastSyncTimestamp = datetime.utcfromtimestamp(float(newSyncTimestamp))

        s.add(n)
        s.commit()
        s.close()
        return True
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
            try:
                s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            except:
                while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
                    s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
                    s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'POSITIVE'
            del h['ID'], h['Dirty']
            #h = ujson.dumps(h)
            positiveHeadersArray.append(h)

        for p in postsNeutral:
            try:
                s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            except:
                while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
                    s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
                    s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'NEUTRAL'
            del h['ID'], h['Dirty']
            #h = ujson.dumps(h)
            neutralHeadersArray.append(h)

        for p in postsNegative:
            try:
                s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            except:
                while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
                    s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
                    s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'NEGATIVE'
            del h['ID'], h['Dirty']
            #h = ujson.dumps(h)
            negativeHeadersArray.append(h)

        for p in postsTopic:
            try:
                s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            except:
                while s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).count() > 1:
                    s.delete(s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).first())
                    s.commit()
            h = s.query(PostHeader).filter(PostHeader.PostFingerprint == p.PostFingerprint).one()
            h = h.asDict()
            h['DIRECTION'] = 'TOPIC'
            del h['ID'], h['Dirty']
            #h = ujson.dumps(h)
            topicHeadersArray.append(h)

        dbReply = positiveHeadersArray + neutralHeadersArray + negativeHeadersArray + topicHeadersArray
        return dbReply

    return threads.deferToThread(threadFunction)

def getPost(fingerprint):

    def threadFunction():
        s = Session()
        try:
            s.query(Post).filter(Post.PostFingerprint == fingerprint).one()
        except:
            while s.query(Post).filter(Post.PostFingerprint == fingerprint).count() > 1:
                s.delete(s.query(Post).filter(Post.PostFingerprint == fingerprint).first())
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

        cprint('FROM REMOTE: POST REQUEST: for post %s' %(fingerprint), 'white', 'on_yellow', attrs=['bold'])
        print('\n')
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
                .filter(Node.LastRetrievedDate >= timedelta(days=30)).all()
            connectedNodes = s.query(Node).filter(Node.LastConnectedDate >= LastSyncTimestamp)\
            .filter(Node.LastConnectedDate >= timedelta(days=30))\
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
