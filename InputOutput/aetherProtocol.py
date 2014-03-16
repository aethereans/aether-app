"""
    This is the implementation of the protocol defined in networkAPI. The first part is inbound, the second part is
    outbound.
"""
from __future__ import print_function
import ujson
from calendar import timegm as toUnixTimestamp
from datetime import datetime, timedelta
import cPickle as pickle
import random

from twisted.protocols import amp
from twisted.internet.protocol import ClientFactory
from twisted.internet.task import LoopingCall
from twisted.internet.endpoints import TCP4ClientEndpoint, SSL4ClientEndpoint
from termcolor import cprint
from twisted.internet import reactor

from InputOutput import networkAPI
from ORM import Mercury, Demeter
from globals import basedir, aetherListeningPort, protocolVersion, profiledir, debugEnabled
import globals
from ORM.Demeter import committer

if not globals.debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass


class AetherProtocol(amp.AMP):
    def __init__(self):
        amp.AMP.__init__(self)
        self.remoteFinished = False
        self.connectedNode = None
        # Sync timestamp received from remote.
        self.remoteSyncTimestamp = None
        # State Machine
        self.connectionState = 'NEWBORN' # Possible values: NEWBORN, HANDSHAKE, HEADER, POST, NODE, CLOSED
        self.lastArrival = datetime.utcnow() # Last time anything happened on the wire.
        # When in HEADER
        self.expectedHeaderPackets = 0 # Ex: 46
        self.arrivedHeaderPackets = 0 # Ex: 1
        self.lastHeaderPacketArrival = None
        self.neededPosts = [] # Gathering place for all the posts that will be requested in POST step.
        # When in POST
        self.expectedPosts = 0
        self.arrivedPosts = 0
        self.lastPostArrival = None
        # When in NODE
        self.expectedNodePackets = 0
        self.arrivedNodePackets = 0
        self.lastNodePacketArrival = None
        # When connection is stuck and pushed to next state manually.
        # The timestamp needs to ne not written, because if it is,
        # the unarrived posts will be marked as successful in the
        # remote and those will never be sent again. Invalidating
        # allows the next sync of that node to start from scratch.
        self.remoteSyncTimestampInvalidated = False

    #def connectionMade(self):

    def connectionLost(self, reason):
        # This errors out if connection is lost before a node is determined. Prevent that.
        if self.connectedNode:
            # This gets called after self.transport.loseConnection. At this point protocol isn't available anymore.
            # So you can't call any methods that are connecting to the remote, the gate is closed.
            if not self.remoteSyncTimestampInvalidated:
                print('L<-R: Timestamp at connection closure: ', self.remoteSyncTimestamp, ' N:', self.connectedNode['NodeId'])
                self.connectedNode['LastSyncTimestamp'] = self.remoteSyncTimestamp
            else:
                print('L<-R: Timestamp at closure, but invalidated: ', self.remoteSyncTimestamp, ' N:', self.connectedNode['NodeId'])

            # If there is no sync timestamp, it's probably a failed connection, one way or another. Don't commit the update.
            if self.connectedNode['LastSyncTimestamp']:
                committer.addNode(self.connectedNode, True)
            self.connectionState = 'CLOSED'
            # try:
            #     self.commitQueue.commit()
            # except AttributeError:
            #     # If the connection closes before late init which creates commitQueue, this can happen. It's OK.
            #     pass
            # This increments all parency of new posts' comment count.
            #Demeter.incrementAncestryCommentCount(self.ephemeralConnectionId)
            print('L=/=R: Connection Closed. Reason: ', reason)

    """
        Errbacks used in protocol.
    """
    def abortConnection(self, exception, side, methodName):
        # Escape from the connecting without sending any data waiting to be sent.
        self.transport.abortConnection()
        self.connectionState = 'CLOSED'
        exceptionText = \
            ('Aborting the sync. Possibly happened in: ', side, ' Method: ', methodName)
        print(exceptionText)
        print(exception)

    """
        This is the outbound (L->R) protocol. These methods send commands to remote.
    """
    def initiateHandshake(self):
        print('L->R: Handshake request to ', self.transport.getPeer().host)

        def replyArrived(reply):
            self.connectionState = 'HANDSHAKE'
            ip = self.transport.getPeer().host
            d = Mercury.handleConnectingNode(reply['NodeId'], ip, reply['ListeningPort']) # returns a node as dict.
            # This sets protocol globals that only become available after node is ID'd.
            d.addCallback(self.postHandshakeInitialization)
            d.addCallback(self.requestHeaders)
            print('L<-R: Handshake reply. N: ', reply['NodeId'])

        self.callRemote(networkAPI.Handshake, NodeId=self.factory.localNodeId, ListeningPort=aetherListeningPort,
                        ProtocolVersion=protocolVersion)\
        .addCallback(replyArrived)\
        .addErrback(self.abortConnection, 'L->R', self.initiateHandshake.__name__)

    def requestHeaders(self, *args):
        # If there is no self.connectedNode at this point, the connection is duplicate. This is handled here because
        # handling it in late init doesn't stop this being called as the next callback in chain.
        if not self.connectedNode:
            self.transport.loseConnection
            return
        # Why no replyArrived? Because I don't need whatever this returns. It returns empty.
        print('L->R: Headers request, Timestamp: ', self.connectedNode['LastSyncTimestamp'] if self.connectedNode['LastSyncTimestamp'] != None else 'Initial connection', ' N: ', self.connectedNode['NodeId'])
        langs = globals.userLanguages # FIXME
        #langs = ['English']
        return self.callRemote(networkAPI.RequestHeaders, LastSyncTimestamp =
                ujson.dumps(self.connectedNode['LastSyncTimestamp']), Languages = langs)\
                .addErrback(self.abortConnection, 'L->R', self.requestHeaders.__name__)

    def requestPost(self, fingerprint):
        print('L->R: Post request, Fingerprint: ', fingerprint, ' N:', self.connectedNode['NodeId'])

        def replyArrived(reply):
            print('L<-R: Post reply, Fingerprint: ', fingerprint, ' N:', self.connectedNode['NodeId'])
            # Things that happen at every arrival.
            self.lastArrival = datetime.utcnow()
            self.lastPostArrival = datetime.utcnow()
            self.arrivedPosts += 1
            self.connectionState = 'POST'
            d = Mercury.insertPost(reply)
            # Things that happen only after all posts arrive.
            if self.expectedPosts == self.arrivedPosts:
                self.advanceToNextStateFromPost()
            return d

        return self.callRemote(networkAPI.RequestPost, PostFingerprint=fingerprint)\
            .addCallback(replyArrived)\
            .addErrback(self.abortConnection, 'L->R', self.requestPost.__name__)

    def advanceToNextStateFromPost(self):
        self.requestNodes()

    def requestNodes(self):
        print('L->R: Nodes request, Timestamp: ', self.connectedNode['LastSyncTimestamp'] if self.connectedNode['LastSyncTimestamp'] != None else 'Initial connection', ' N:', self.connectedNode['NodeId'])
        self.callRemote(networkAPI.RequestNodes, LastSyncTimestamp =
            unicode(toUnixTimestamp(self.connectedNode['LastSyncTimestamp'].utctimetuple())
            if self.connectedNode['LastSyncTimestamp'] != None else 'null'))\
            .addErrback(self.abortConnection, 'L->R', self.requestNodes.__name__)
        return

    def sendTimestamp(self):
        timestamp = unicode(toUnixTimestamp(datetime.utcnow().utctimetuple()))
        print('L->R: New Timestamp: : ', timestamp, ' N:', self.connectedNode['NodeId'])
        return self.callRemote(networkAPI.SyncTimestamps,  NewSyncTimestamp=timestamp)

    def postHandshakeInitialization(self, node):
        # If there are any other protocols not in CLOSED state
        # that are currently connected to this node, kill this.
        for prot in self.factory.openConnections:
            # Ensure we're looping within prots which have a connectedNode and is not closed.
            if prot.connectionState != 'CLOSED' \
            and prot.connectedNode \
            and prot.connectedNode['NodeId'] == node['NodeId']:
                cprint('I HAVE THE SAME CONNECTION OPEN', 'red')
                return
                # Right. actually returning here doesn't stop anything.. the next callback still fires.
        self.connectedNode = node

    """
        This is inbound (L<-R). These methods respond to commands that arrive from the remote.
    """

    @networkAPI.Handshake.responder
    def respondToHandshake(self, NodeId, ListeningPort, ProtocolVersion):
        print('L<-R: Handshake request from ', self.transport.getPeer().host, ' N:', NodeId)
        self.connectionState = 'HANDSHAKE'
        ip = self.transport.getPeer().host
        if len(self.factory.openConnections) > globals.maxInboundCount:
            # cprint('Connection refused because I\'m too busy.', 'red')
            self.transport.loseConnection()
        else:
            d = Mercury.handleConnectingNode(NodeId, ip, ListeningPort) # returns a node as dict.
            d.addCallback(self.postHandshakeInitialization)
            # Reverse sync here. I'm asking the remote about his contents.
            d.addCallback(self.requestHeaders)\
             .addErrback(self.abortConnection, 'L<-R', self.requestHeaders.__name__)
        reply = {'NodeId': self.factory.localNodeId,
                 'ListeningPort': aetherListeningPort,
                 'ProtocolVersion': protocolVersion }
        print('L->R: Handshake reply sent.', ' N:', NodeId)
        return reply

    @networkAPI.RequestHeaders.responder
    def respondWithHeaders(self, LastSyncTimestamp, Languages):
        #Check if handshake actually occurred. If not, abort.
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'respondWithHeaders')
            return {}
        if LastSyncTimestamp != 'null':
            LastSyncTimestamp = datetime.utcfromtimestamp(float(LastSyncTimestamp))
        else:
            LastSyncTimestamp = 'NO-TIMESTAMP'
        print('L<-R: Headers request, Timestamp: ', LastSyncTimestamp if LastSyncTimestamp != 'NO-TIMESTAMP' else 'Initial Connection', ' N:', self.connectedNode['NodeId'])
        d = Mercury.getHeaders(LastSyncTimestamp, Languages)
        currentPacketNumber = 1
        totalPacketNumber = 1 # start from one.

        def calculatePacketCounts(dbReply, bucketSize):
            global totalPacketNumber
            totalPacketNumber = len(dbReply)/bucketSize
            if len(dbReply)%bucketSize > 0:
                totalPacketNumber+= 1
            if len(dbReply) == 0:
                totalPacketNumber = 1
            return dbReply

        d.addCallback(calculatePacketCounts, globals.headerPacketCount)

        def bucketise(dbReply, bucketSize, CPNo):

            def process(slicedDbReply, currentPacketNumber):
                # Here I need to fire the bucket given.

                def construct(bucket):
                    positiveHeadersArray = []
                    neutralHeadersArray = []
                    negativeHeadersArray = []
                    topicHeadersArray = []
                    for h in bucket:
                        if h['DIRECTION'] == 'POSITIVE':
                            del h['DIRECTION']
                            hJson = ujson.dumps(h)
                            positiveHeadersArray.append(hJson)
                        elif h['DIRECTION'] == 'NEUTRAL':
                            del h['DIRECTION']
                            hJson = ujson.dumps(h)
                            neutralHeadersArray.append(hJson)
                        elif h['DIRECTION'] == 'NEGATIVE':
                            del h['DIRECTION']
                            hJson = ujson.dumps(h)
                            negativeHeadersArray.append(hJson)
                        elif h['DIRECTION'] == 'TOPIC':
                            del h['DIRECTION']
                            hJson = ujson.dumps(h)
                            topicHeadersArray.append(hJson)
                    return {
                            'PositiveHeaders': positiveHeadersArray,
                            'NeutralHeaders': neutralHeadersArray,
                            'NegativeHeaders': negativeHeadersArray,
                            'TopicHeaders': topicHeadersArray
                    }
                    # Here i need to get the bucket, produce a valid amp packet with the data given and return that.
                    # I also convert to JSON here.

                def send(reply): # Here I send the amp packet produced.
                    global totalPacketNumber
                    print('L->R: Header reply, Package number: ', CPNo, ' N:', self.connectedNode['NodeId'])
                    reply['TotalNumberOfPackets'] = totalPacketNumber
                    reply['CurrentPacketNo'] = currentPacketNumber
                    self.callRemote(networkAPI.ReceiveHeaders,
                                    PositiveHeaders = reply['PositiveHeaders'],
                                    NeutralHeaders = reply['NeutralHeaders'],
                                    NegativeHeaders = reply['NegativeHeaders'],
                                    TopicHeaders = reply['TopicHeaders'],
                                    TotalNumberOfPackets = reply['TotalNumberOfPackets'],
                                    CurrentPacketNo = reply['CurrentPacketNo']
                                    )
                    return {}
                return send(construct(slicedDbReply))

            if len(dbReply) > bucketSize:
                bucketise(dbReply[bucketSize:], bucketSize, CPNo + 1) # next batch.
                return process(dbReply[0:bucketSize], CPNo) # current batch.
            else:
                return process(dbReply[0:bucketSize], CPNo) # current batch.

        d.addCallback(bucketise, globals.headerPacketCount, currentPacketNumber)
        d.addErrback(self.abortConnection, 'L<-R', Mercury.getHeaders.__name__)
        #return d # bad local return?
        return {}

    @networkAPI.ReceiveHeaders.responder
    def receiveHeaderPacket(self, PositiveHeaders, NeutralHeaders, NegativeHeaders,
                            TopicHeaders, TotalNumberOfPackets, CurrentPacketNo):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'receiveHeaderPacket')
        print('L<-R: Header reply, Package number: ', CurrentPacketNo, ' N:', self.connectedNode['NodeId'])
        # Define what just arrived from the network.
        reply = {
            'PositiveHeaders': PositiveHeaders,
            'NeutralHeaders': NeutralHeaders,
            'NegativeHeaders': NegativeHeaders,
            'TopicHeaders': TopicHeaders,
            'TotalNumberOfPackets': TotalNumberOfPackets,
            'CurrentPacketNo': CurrentPacketNo
        }

        # Set the global state machine information.
        self.lastArrival = datetime.utcnow()
        self.connectionState = 'HEADER'

        # Set the local data about this state's status.
        self.expectedHeaderPackets = reply['TotalNumberOfPackets']



        # Things that need to happen with every packet.
        d = Mercury.insertHeadersAndVotes(reply, self.connectedNode)
        d.addCallback(Mercury.listNeededPosts) # Returns: an array of needed posts' fingerprints.

        def saveNeededPostsList(currentPacketNeededPosts):
            #print('current packet needed posts: ', currentPacketNeededPosts)
            # for postFingerprint in currentPacketNeededPosts:
            #     if postFingerprint not in self.neededPosts:
            #         self.neededPosts.append(postFingerprint)
            #         self.expectedPosts += 1
            # I don't need to do this check, because I am doing this kind of check at the commit phase already.
            self.neededPosts += currentPacketNeededPosts
            self.expectedPosts += len(currentPacketNeededPosts)

            # Set the local data about this state's status.
            self.arrivedHeaderPackets += 1
            self.lastHeaderPacketArrival = datetime.utcnow()
            # I moved them here from one outer scope, because when the remote has a million posts, listneededposts
            # take much longer than header arrival. as checkifdone is now async, checkifdone actually triggers
            # much later than header arrival, unlike intended use of the last one triggering it, almost all triggers.

        d.addCallback(saveNeededPostsList)

        # Things that need to happen only after all header posts have arrived.
        def checkIfDone(*a):
            if self.expectedHeaderPackets == self.arrivedHeaderPackets:
                self.advanceToNextStateFromHeader()
        # This is a callback because before I did that it was triggering before everything was actually appended.
        d.addCallback(checkIfDone)
        return {}

    def advanceToNextStateFromHeader(self):
        for postFingerprint in self.neededPosts:
                self.requestPost(postFingerprint)
        if len(self.neededPosts) == 0:
            self.advanceToNextStateFromPost()

    # For RequestPost, the reply is the post itself, so there is no ReceivePost, there is no need to.
    @networkAPI.RequestPost.responder
    def respondWithPost(self, PostFingerprint):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'respondWithPost')
        print('L->R: Post reply. Post: ', PostFingerprint, ' N:', self.connectedNode['NodeId'])
        d = Mercury.getPost(PostFingerprint)
        d.addCallback(Mercury.toJson)

        def processAndReturn(postAsJson):
            reply = { 'Post': postAsJson }
            return reply

        d.addCallback(processAndReturn)\
         .addErrback(self.abortConnection, 'L<-R', self.respondWithPost.__name__)
        return d

    @networkAPI.RequestNodes.responder
    def respondWithNodes(self, LastSyncTimestamp):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'respondWithNodes')
        if LastSyncTimestamp != 'null':
            LastSyncTimestamp = datetime.utcfromtimestamp(float(LastSyncTimestamp))
        else:
            LastSyncTimestamp = 'NO-TIMESTAMP'
        print('L<-R: Nodes request, Timestamp: ', LastSyncTimestamp if LastSyncTimestamp != 'NO-TIMESTAMP' else 'Initial Connection', ' N:', self.connectedNode['NodeId'])
        d = Mercury.getNodes(LastSyncTimestamp)
        d.addCallback(Mercury.processNodes)
        currentPacketNumber = 1
        totalPacketNumber = 1

        def calculatePacketCounts(dbReply, bucketSize):
            global totalPacketNumber
            totalPacketNumber = len(dbReply)/bucketSize

            if len(dbReply)%bucketSize > 0:
                totalPacketNumber+= 1
            if len(dbReply) == 0:
                totalPacketNumber = 1
            return dbReply

        d.addCallback(calculatePacketCounts, globals.nodePacketCount)

        def bucketise(dbReply, bucketSize ,CPNo):
            # Here I have the entire reply from the db. I need divide that reply into buckets, and fire each
            # bucket to the remote.

            def process(slicedDbReply, currentPacketNumber):
                print('L->R: Nodes reply, Package number: ', currentPacketNumber, ' N:', self.connectedNode['NodeId'])
                global totalPacketNumber
                nodesJson = []
                for n in slicedDbReply:
                    nJson = ujson.dumps(n)
                    nodesJson.append(nJson)
                self.callRemote(networkAPI.ReceiveNodes,
                                Nodes = nodesJson,
                                TotalNumberOfPackets = totalPacketNumber,
                                CurrentPacketNo = currentPacketNumber
                                )
                return {}

            if len(dbReply) > bucketSize:
                bucketise(dbReply[bucketSize:], bucketSize, CPNo + 1) # next batch.
                return process(dbReply[0:bucketSize], CPNo) # current batch.
            else:
                return process(dbReply[0:bucketSize], CPNo) # current batch.

        d.addCallback(bucketise, globals.nodePacketCount, currentPacketNumber)
        d.addErrback(self.abortConnection, 'L<-R', self.respondWithNodes.__name__)
        return d

    @networkAPI.ReceiveNodes.responder
    def receiveNodePacket(self, Nodes, TotalNumberOfPackets, CurrentPacketNo):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'receiveNodePacket')
        print('L<-R: Nodes reply, Package number: ', CurrentPacketNo, ' N:', self.connectedNode['NodeId'])
        # Set the global state machine information.
        self.lastArrival = datetime.utcnow()
        self.connectionState = 'NODE'

        # Set the local data about this state's status.
        self.expectedNodePackets = TotalNumberOfPackets
        self.arrivedNodePackets += 1

        # Things that need to happen at every node packet arrival.
        Mercury.insertNodes(Nodes, self.factory.localNodeId)

        # Things that need to happen only after all node packets arrive.
        if self.expectedNodePackets == self.arrivedNodePackets:
            self.advanceToNextStateFromNode()
        return {}

    def advanceToNextStateFromNode(self):
        # Save everything gathered so far into the database.
            self.sendTimestamp()
            # And we're done. Check if remote is done already.
            if self.remoteFinished:
                # If so, insert the timestamp if not invalidated.

                # And kill the connection.
                self.connectionState = 'CLOSED'
                print('L=/=R: I have received', self.arrivedHeaderPackets, 'headers,', self.arrivedPosts, 'posts, and', self.arrivedNodePackets, 'nodes.')
                self.transport.loseConnection()
            else:
                # Tell remote that this computer is done. It will close the connection when done.
                print('L->R: Finished', ' N:', self.connectedNode['NodeId'])
                self.callRemote(networkAPI.KillConnection)


    @networkAPI.SyncTimestamps.responder
    def receiveTimestamp(self, NewSyncTimestamp):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'receiveTimestamp')
        print('L<-R: New Timestamp: ', NewSyncTimestamp, ' N:', self.connectedNode['NodeId'])
        self.remoteSyncTimestamp = datetime.utcfromtimestamp(float(NewSyncTimestamp))
        return {}

    @networkAPI.KillConnection.responder
    def receiveKillConnection(self):
        if not self.connectedNode:
            self.abortConnection(Exception('AetherError: Illegal request from remote.'), 'L<-R', 'receiveKillConnection')
        print('L<-R: Finished', ' N:', self.connectedNode['NodeId'])
        self.remoteFinished = True
        return {}


class AetherProtocolFactory(ClientFactory):
    def __init__(self):
        self.openConnections = []
        self.overseerLoop = LoopingCall(self.overseer)
        self.overseerLoop.start(5)

    def buildProtocol(self, addr):
        protocol = AetherProtocol()
        protocol.factory = self
        self.openConnections.append(protocol)
        return protocol

    def stopFactory(self):
        for protocol in self.openConnections:
            if protocol.connectionState != 'CLOSED':
                protocol.transport.loseConnection()

    with open(profiledir + 'UserProfile/backendSettings.dat', 'rb') as f2:
        localNodeId = pickle.load(f2)

    # This should be called every five seconds.
    def overseer(self):
        print(len(self.openConnections), 'connections open at', datetime.utcnow())
        for conn in self.openConnections:
            print(conn.connectionState)
        # Reversed: So if I remove something from the list, I won't skip an item.
        for protocol in reversed(self.openConnections):
            # List open connections for debug.
            #print('Currently open connection to: ', protocol.connectedNode['NodeId'])
            # commented out because we removed the guards, so connections without connected nodes (newborn, handshake) also
            # get into this.
            tenSecondsAgo = datetime.utcnow() - timedelta(seconds=10)
            thirtySecondsAgo = datetime.utcnow() - timedelta(seconds=30)
            state = protocol.connectionState
            lastActivity = protocol.lastArrival
            if protocol.connectionState == 'CLOSED':
                # If CLOSED, remove the connection from the list.
                self.openConnections.remove(protocol)
            # If last activity is older than ten seconds.
            if lastActivity < thirtySecondsAgo:
                # After 3 tries of below, if nothing arrives, close the connection.
                protocol.abortConnection(Exception('AetherError: Remote unresponsive'), 'L<-R', 'Factory Overseer')
            if lastActivity < tenSecondsAgo:
                if state == 'HANDSHAKE' or state == 'NEWBORN':
                    protocol.abortConnection(Exception('AetherError: Remote is borked. Stuck in handshake, got killed.'), 'L<-R', 'Factory Overseer')
                    self.openConnections.remove(protocol)
                if state == 'HEADER':
                    # Force advance to the next state.
                    print('Expected header packets: ', protocol.expectedHeaderPackets, 'Arrived header packets: ', protocol.arrivedHeaderPackets)
                    print('Connection stuck. Force advancing from HEADER to POST')
                    # Invalidate the timestamp.
                    protocol.remoteSyncTimestampInvalidated = True
                    # Advance to the next state.
                    protocol.advanceToNextStateFromHeader()
                elif state == 'POST':
                    # Force advance to the next state.
                    print('Expected posts: ', protocol.expectedPosts, 'Arrived posts: ', protocol.arrivedPosts)
                    print('Connection stuck. Force advancing from POST to NODE')
                    # Invalidate the timestamp.
                    protocol.remoteSyncTimestampInvalidated = True
                    # Advance to the next state.
                    protocol.advanceToNextStateFromPost()
                elif state == 'NODE':
                    # Force advance to the next state.
                    print('Expected node packets: ', protocol.expectedNodePackets, 'Arrived node packets: ', protocol.arrivedNodePackets)
                    print('Connection stuck. Force advancing from NODE to CLOSED')
                    # Invalidate the timestamp.
                    #protocol.remoteSyncTimestampInvalidated = True
                    # Advance to the next state.
                    protocol.advanceToNextStateFromNode()

    #def clientConnectionFailed(self, connector, reason):
    #    print('Connection failed, because: %s' %reason.getErrorMessage())
    #    pass


aetherProtocolFactoryInstance = AetherProtocolFactory()


"""
    The API used for connecting to the nodes, from other parts of the application.
    The two methods below are usually all you need.
"""
def connectWithNode(node):
    ip = node['LastConnectedIP'] if node['LastConnectedIP'] != None else node['LastRetrievedIP']
    port = node['LastConnectedPort'] if node['LastConnectedPort'] != None else node['LastRetrievedPort']
    endpoint = SSL4ClientEndpoint(reactor, ip , int(port), globals.aetherClientContextFactoryInstance)
    endpoint.connect(aetherProtocolFactoryInstance)\
        .addCallback(lambda p: p.initiateHandshake())\
        .addErrback(printError,
            'Connection failed because node at the address %s is not responding.'%ip)

def connectWithIP(ip, Port):
    print('Connect with Ip: %s:%s' %(ip, Port))
    endpoint = SSL4ClientEndpoint(reactor, ip , int(Port), globals.AetherClientContextFactory())
    endpoint.connect(aetherProtocolFactoryInstance)\
        .addCallback(lambda p: p.initiateHandshake())\
        .addErrback(printError,
            'Connection failed because node at the address %s is not responding.'%ip)

# Errbacks used in API above.

def printError(exception, definition):
    print('Connection Failed. ', definition, exception)
    pass

