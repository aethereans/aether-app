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
from globals import basedir, aetherListeningPort, protocolVersion
import globals

class AetherProtocol(amp.AMP):
    def __init__(self):
        self.statusCheckerLoop = LoopingCall(self.connectionStatusCheck)
        amp.AMP.__init__(self)
        self.connectedNode = None
        self.headerTransferStarted = False
        self.nodeTransferStarted = False
        self.totalNumberOfHeaderPackets = 0
        self.currentHeaderPacketCount = 0
        self.totalNumberOfNodePackets =0
        self.currentNodePacketCount = 0
        self.connectionStartTime = datetime.utcnow()
        self.localConnectionFinishTime = datetime.utcnow()
        self.headersDone = False
        self.nodesDone = False
        self.syncDone = False
        self.remoteIsOkToClose = False
        self.localIsOkToClose = False
        self.ephemeralConnectionId = int(random.random()*10**10)

    incomingPostCount = 0
    # So this is the number of posts that this connection needs to wait for. Length of needed posts sets this, and every
    # successful incoming post decrements by 1. every incoming post also calls a check on this number, and when it hits
    # 0, connection closes.

    def connectionStatusCheck(self):
        print('header transfer started: ', self.headerTransferStarted)
        print('node transfer started: ', self.nodeTransferStarted)
        print('currentHeaderPacketCount: ', self.currentHeaderPacketCount)
        print('total header packet count: ', self.totalNumberOfHeaderPackets)
        print('current node packet count: ', self.currentNodePacketCount)
        print('total node packet count: ', self.totalNumberOfNodePackets)

        deltaTime = (datetime.utcnow() - self.connectionStartTime)
        if deltaTime > timedelta(minutes=globals.maximumAllowedConnectionTimespan):
            cprint('MAXIMUM ALLOWED CONNECTION TIMESPAN EXCEEDED. Closing the connection.', 'white', 'on_red', attrs=['bold'])
            self.factory.currentConnectionCount -= 1
            self.transport.loseConnection()

        deltaSinceLocalFinishedConnection = (datetime.utcnow() - self.localConnectionFinishTime)
        if deltaSinceLocalFinishedConnection > timedelta(seconds=30):
            cprint('Maximum wait time after finish is exceeded. Closing the connection.', 'white', 'on_red', attrs=['bold'])
            self.factory.currentConnectionCount -= 1
            self.transport.loseConnection()

        if (self.headerTransferStarted and (self.currentHeaderPacketCount == self.totalNumberOfHeaderPackets or
        self.currentHeaderPacketCount +1 == self.totalNumberOfHeaderPackets)) and \
            self.nodeTransferStarted and (self.currentNodePacketCount == self.totalNumberOfNodePackets or
        self.currentNodePacketCount + 1 == self.totalNumberOfNodePackets):
            cprint('ALL HEADER PACKAGES and their POSTS HAVE ARRIVED!', 'red', 'on_white')
            cprint('ALL NODE PACKAGES HAVE ARRIVED!', 'red', 'on_white')
            d = self.syncTimestamps() # I ASK for a timestamp. ## OKAY. THIS SHOULD HAPPEN AT THE VERY END, AT CONNECTION CLOSURE.
            def setFlag(ignored):
                self.syncDone = True
            d.addCallback(setFlag)
            d.addCallback(self.requestClosureOfAetherConnection)

        print('remote is ok: ', self.remoteIsOkToClose)
        print('local is ok: ', self.localIsOkToClose)
        print('sync is done: ', self.syncDone)

        if self.remoteIsOkToClose and self.localIsOkToClose and self.syncDone:
            print('Both sides are OK with closing the connection.')
            self.transport.loseConnection()
            # TODO: Implement Aether Keepalive.

    def connectionMade(self):
        self.statusCheckerLoop.start(3) # FIXME: This should be 10 in normal cases. maybe not..
        self.factory.currentConnectionCount += 1
        self.connectionStartTime = datetime.utcnow()

    def connectionLost(self, reason):
        print('Connection lost. Reason: %s' %reason)
        self.factory.currentConnectionCount -= 1
        self.statusCheckerLoop.stop()
        Demeter.incrementAncestryCommentCount(self.ephemeralConnectionId) # this increments all parency of new posts' comment count.
        try:
            cprint('REMOTE %s HAS SUCCESSFULLY CLOSED THE CONNECTION at UTC %s'
                   %(self.connectedNode['NodeId'], datetime.utcnow()), 'grey', 'on_white', attrs=['bold'])
        except:
            cprint('SOMETHING WENT HORRIBLY WRONG. THE CONNECTION FAILED BEFORE THE HANDSHAKE. ONE OF THE PEERS MIGHT HAVE '
                   'ATTEMPTED AN ILLEGAL REQUEST SUCH AS A NON-ENCRYPTED CONNECTION OR DIDN\'T FOLLOW THE PROTOCOL. \nEXITING.'
                , 'white', 'on_red', attrs=['bold'])

    """
        Errbacks used in protocol.
    """

    def closeConnection(self, exception, side, methodName):
        self.transport.loseConnection() # I think twisted auto closes when things go wrong.
        exceptionText = \
            ('The sync process got borked. '
            'Which happened at %s side of this machine in %s method call. Closing the connection.'
            %(side, methodName))
        print(exceptionText)
        return exception
    """
        This is the outbound protocol. These methods send requests.
    """
    def initiateHandshake(self):
        cprint('ASKING FOR: HANDSHAKE to %s:%s'
               %(self.transport.getPeer().host, self.transport.getPeer().port),
               'cyan', 'on_blue', attrs=['bold'])
        def replyArrived(reply):
            cprint('RECEIVED: HANDSHAKE REPLY. \n%s' %(reply), 'cyan', 'on_blue')
            print('\n')
            ip = self.transport.getPeer().host
            d = Mercury.checkIfNodeExists(reply['NodeId'])

            def processNode(nodeExists):
                if nodeExists:
                    return Mercury.updateAndGetNode(reply['NodeId'], ip, reply['ListeningPort'])
                    # This returns deferred
                else:
                    return Mercury.createNode(reply['NodeId'], ip, reply['ListeningPort'])
            d.addCallback(processNode)
            d.addCallback(self.setGlobalConnectedNode) # This is at one level up, directly below protocol class.
            d.addCallback(self.requestHeaders)
            d.addCallback(self.requestNodes)

        self.callRemote(networkAPI.Handshake,
                        NodeId=self.factory.localNodeId,
                        ListeningPort=aetherListeningPort,
                        ProtocolVersion=protocolVersion)\
        .addCallback(replyArrived)\
        .addErrback(self.closeConnection, 'OUTBOUND', self.initiateHandshake.__name__)

    def requestHeaders(self, node):
        cprint('ASKING FOR: HEADERS with timestamp %s'
        %(node['LastSyncTimestamp'] if node['LastSyncTimestamp'] != None
          else 'NO-TIMESTAMP'), 'cyan', 'on_blue', attrs=['bold'])
        print('\n')
        cprint('OLD REMOTE TIMESTAMP - SENT TO REMOTE: %s' %(node['LastSyncTimestamp']
        if node['LastSyncTimestamp'] != None else 'NO-TIMESTAMP'), 'cyan', 'on_blue', attrs=['bold'])
        print('\n')

        def replyArrived(reply):
            return node

        # hmm.. here I need to figure out a way to get replyArrived to wait till all replies come. okay, this
        # is a little more complex than what I thought..

        langs = globals.userLanguages
        return self.callRemote(networkAPI.RequestHeaders, LastSyncTimestamp =
            ujson.dumps(node['LastSyncTimestamp']), Languages = langs)\
            .addCallback(replyArrived)\
            .addErrback(self.closeConnection, 'OUTBOUND', self.requestHeaders.__name__)

    def requestPost(self, fingerprint):
        cprint('ASKING FOR: POST with fingerprint %s' %(fingerprint), 'cyan', 'on_blue', attrs=['bold'])

        def replyArrived(reply):
            return Mercury.insertPost(reply, self.ephemeralConnectionId)

        return self.callRemote(networkAPI.RequestPost, PostFingerprint=fingerprint)\
            .addCallback(replyArrived)\
            .addErrback(self.closeConnection, 'OUTBOUND', self.requestPost.__name__)

    def requestNodes(self, node):
        cprint('ASKING FOR: NODES with timestamp %s'
            %(node['LastSyncTimestamp'] if node['LastSyncTimestamp'] != None
              else 'NO-TIMESTAMP'), 'cyan', 'on_blue', attrs=['bold'])



        self.callRemote(networkAPI.RequestNodes, LastSyncTimestamp =
            unicode(toUnixTimestamp(node['LastSyncTimestamp'].utctimetuple())
            if node['LastSyncTimestamp'] != None else 'null'))\
            .addErrback(self.closeConnection, 'OUTBOUND', self.requestNodes.__name__)
        return node

    def syncTimestamps(self):
        cprint('NEW LOCAL TIMESTAMP - SENT TO THE REMOTE: %s' %datetime.utcnow(),
                   'cyan', 'on_blue', attrs=['bold'])
        print('\n')

        return self.callRemote(networkAPI.SyncTimestamps,
                NewSyncTimestamp=unicode(toUnixTimestamp(datetime.utcnow().utctimetuple())))\
                .addErrback(self.closeConnection, 'OUTBOUND', self.syncTimestamps.__name__)

    def requestClosureOfAetherConnection(self, *args):
        cprint('TO REMOTE: OK TO DISCONNECT WHEN DONE.', 'white', 'on_green')
        self.localIsOkToClose = True
        self.localConnectionFinishTime = datetime.utcnow()
        self.callRemote(networkAPI.KillConnection)\
            .addErrback(self.closeConnection, 'OUTBOUND', self.requestClosureOfAetherConnection.__name__)

    def setGlobalConnectedNode(self, node):
        self.connectedNode = node
        return self.connectedNode

    """
        This is inbound. These methods respond to requests that arrive from the network.
    """

    @networkAPI.Handshake.responder
    def respondToHandshake(self, NodeId, ListeningPort, ProtocolVersion):
        ip = self.transport.getPeer().host
        d = Mercury.checkIfNodeExists(NodeId)

        def processNode(nodeExists):
            if nodeExists:
                return  Mercury.updateAndGetNode(NodeId, ip, ListeningPort) # This returns deferred
            else:
                return Mercury.createNode(NodeId, ip, ListeningPort)

        d.addCallback(processNode)

        def callReverseSync(node):
            # I'm going to ask for nodes in parallel with already existing header call.
            d = self.requestHeaders(node)
            d.addCallback(self.requestNodes)
            return node

        d.addCallback(callReverseSync)\
         .addErrback(self.closeConnection, 'INBOUND', callReverseSync.__name__)

        d.addCallback(self.setGlobalConnectedNode) # This is at one level up, directly below protocol class.
        reply = {'NodeId': self.factory.localNodeId,
                 'ListeningPort': aetherListeningPort,
                 'ProtocolVersion': protocolVersion }
        cprint('FROM REMOTE: HANDSHAKE REQUEST: from %s:%s'
                %(ip, ListeningPort), 'white', 'on_yellow', attrs=['bold'])
        cprint('ANSWER: %s' %(reply), 'white', 'on_yellow')
        print('\n')
        return reply

    @networkAPI.RequestHeaders.responder
    def respondWithHeaders(self, LastSyncTimestamp, Languages):
        if LastSyncTimestamp != 'null':
            LastSyncTimestamp = datetime.utcfromtimestamp(float(LastSyncTimestamp))# + timedelta(minutes = 1)
        else:
            LastSyncTimestamp = 'NO-TIMESTAMP'
        cprint('FROM REMOTE: OLD LOCAL TIMESTAMP: %s' %(LastSyncTimestamp
               if LastSyncTimestamp != 'NO-TIMESTAMP' else LastSyncTimestamp), 'white', 'on_yellow', attrs=['bold'])
        print('\n')
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
            print('the reply includes %s headers' %len(dbReply))
            print('total packet numbers: %s' %totalPacketNumber)
            return dbReply

        d.addCallback(calculatePacketCounts, globals.headerPacketCount)

        def bucketise(dbReply, bucketSize, CPNo):

            def process(slicedDbReply, currentPacketNumber):
                # here I need to fire the bucket given.

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
                    # I ALSO CONVERT TO JSON HERE.

                def send(reply): # here I send the amp packet produced.
                    global totalPacketNumber
                    cprint('THIS IS PACKAGE NUMBERED: %d' %CPNo, 'white', 'on_red')
                    reply['TotalNumberOfPackets'] = totalPacketNumber
                    reply['CurrentPacketNo'] = currentPacketNumber
                    print('reply: ',reply)
                    # here I need to fire receiveheaders methods onto the remote for each of the packets.
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

        d.addErrback(self.closeConnection, 'INBOUND', Mercury.getHeaders.__name__)
        return d

    @networkAPI.ReceiveHeaders.responder
    def receiveHeaderPacket(self, PositiveHeaders, NeutralHeaders, NegativeHeaders,
                            TopicHeaders, TotalNumberOfPackets, CurrentPacketNo):
        cprint('RECEIVE HEADER PACKET RESPONDER WAS CALLED.', 'white', 'on_red')

        # This is just a proof of concept, so I'm mushing all the answers together to create the reply replyArrived
        # (below) expects.
        self.headerTransferStarted = True
        # Here, I assign the totalnumber to the protocol and start counting from that.
        self.totalNumberOfHeaderPackets = TotalNumberOfPackets
        print('total number of header packets: ',self.totalNumberOfHeaderPackets)
        print('current packet count: ', self.currentHeaderPacketCount)
        reply = {
            'PositiveHeaders': PositiveHeaders,
            'NeutralHeaders': NeutralHeaders,
            'NegativeHeaders': NegativeHeaders,
            'TopicHeaders': TopicHeaders,
            'TotalNumberOfPackets': TotalNumberOfPackets,
            'CurrentPacketNo': CurrentPacketNo
        }
        node = self.connectedNode

        # I have two tasks, the first is to insert appropriate votes,
        # the second is to decide on which posts I will request.
        d = Mercury.insertHeaders(reply)
        d.addCallback(Mercury.insertVotes, node)
        d.addCallback(Mercury.listNeededPosts) # this will return an array of needed posts.

        def decrementAndCheckPostCounter(*args):
            setattr(self, 'incomingPostCountForHeaderPacket'+str(CurrentPacketNo),
                    getattr(self, 'incomingPostCountForHeaderPacket'+str(CurrentPacketNo)) - 1)  if \
                getattr(self, 'incomingPostCountForHeaderPacket'+str(CurrentPacketNo)) != 0 else 0
            # Here is how I figure out the connection is completed.
            if getattr(self, 'incomingPostCountForHeaderPacket'+str(CurrentPacketNo)) == 0:
                print('incoming posts of header %s is completed'%CurrentPacketNo)
                # This only enters after the posts of the header packet given complete transmitting.
                # This is the final point of process of a header and this marks this header packet done.
                self.currentHeaderPacketCount += 1

        def firePostRequests(neededPosts):
            setattr(self, 'incomingPostCountForHeaderPacket'+str(CurrentPacketNo), len(neededPosts))
            for fp in neededPosts:
                d2 = self.requestPost(fp)
                d2.addCallback(decrementAndCheckPostCounter)
            if len(neededPosts) == 0:
                decrementAndCheckPostCounter()

        d.addCallback(firePostRequests)
        return node

    @networkAPI.RequestPost.responder
    def respondWithPost(self, PostFingerprint):
        d = Mercury.getPost(PostFingerprint)
        d.addCallback(Mercury.toJson)

        def processAndReturn(postAsJson):
            reply = { 'Post': postAsJson }
            return reply

        d.addCallback(processAndReturn)\
         .addErrback(self.closeConnection, 'INBOUND', self.respondWithPost.__name__)
        return d

    @networkAPI.RequestNodes.responder
    def respondWithNodes(self, LastSyncTimestamp):
        if LastSyncTimestamp != 'null':
            LastSyncTimestamp = datetime.utcfromtimestamp(float(LastSyncTimestamp))# + timedelta(minutes = 1)
        else:
            LastSyncTimestamp = 'NO-TIMESTAMP'
        d = Mercury.getNodes(LastSyncTimestamp)
        d.addCallback(Mercury.processNodes)

        cprint('FROM REMOTE: NODES REQUEST: with sync timestamp %s' %(LastSyncTimestamp
                   if LastSyncTimestamp != 'NO-TIMESTAMP' else LastSyncTimestamp), 'white', 'on_yellow', attrs=['bold'])
        print('\n')

        currentPacketNumber = 1
        totalPacketNumber = 1

        def calculatePacketCounts(dbReply, bucketSize):
            global totalPacketNumber
            totalPacketNumber = len(dbReply)/bucketSize

            if len(dbReply)%bucketSize > 0:
                totalPacketNumber+= 1
            if len(dbReply) == 0:
                totalPacketNumber = 1
            print('the reply includes %s headers' %len(dbReply))
            print('total packet numbers: %s' %totalPacketNumber)
            return dbReply

        d.addCallback(calculatePacketCounts, globals.nodePacketCount)

        def bucketise(dbReply, bucketSize ,CPNo):
            # Here I have the entire reply from the db. I need divide that reply into buckets, and fire each
            # bucket to the remote.

            def process(slicedDbReply, currentPacketNumber):
                global totalPacketNumber
                nodesJson = []
                for n in slicedDbReply:
                    nJson = ujson.dumps(n)
                    nodesJson.append(nJson)
                cprint('THIS IS NODE ANSWER PACKAGE NUMBERED: %d' %currentPacketNumber, 'white', 'on_red')

                print('reply: ',{'Nodes':nodesJson,
                                 'TotalNumberOfPackets':totalPacketNumber,
                                 'CurrentPacketNo':currentPacketNumber})

                # here I need to fire receiveheaders methods onto the remote for each of the packets.
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

        d.addErrback(self.closeConnection, 'INBOUND', self.respondWithNodes.__name__)
        return d

    @networkAPI.ReceiveNodes.responder
    def receiveNodePacket(self, Nodes, TotalNumberOfPackets, CurrentPacketNo):
        cprint('RECEIVE NODE PACKET RESPONDER WAS CALLED.', 'white', 'on_yellow')
        self.totalNumberOfNodePackets = TotalNumberOfPackets

        print('total number of NODE packets: ',self.totalNumberOfNodePackets)
        print('current NODE packet count: ', self.currentNodePacketCount)
        self.nodeTransferStarted = True

        d= Mercury.insertNodes(Nodes, self.factory.localNodeId)
        def incrementNodePacketCount(ignored):
            self.currentNodePacketCount += 1
        d.addCallback(incrementNodePacketCount)
        return {}

    @networkAPI.SyncTimestamps.responder
    def receiveNewSyncTimestamp(self, NewSyncTimestamp):
        newSyncTimestamp = NewSyncTimestamp
        cprint('FROM REMOTE: NEW REMOTE TIMESTAMP: %s' %newSyncTimestamp, 'white', 'on_yellow', attrs=['bold'])
        print('\n')
        Mercury.insertNewSyncTimestamp(self.connectedNode, newSyncTimestamp)
        return {}

    @networkAPI.KillConnection.responder
    def receiveKillConnection(self):
        cprint('FROM REMOTE: OK TO DISCONNECT WHEN DONE.', 'white', 'on_red', attrs=['bold'])
        self.remoteIsOkToClose = True
        return {}


class AetherProtocolFactory(ClientFactory):
    # Anything here is globally accessible to ALL protocols instantiated.
    protocol = AetherProtocol
    with open(basedir + 'UserProfile/backendSettings.dat', 'rb') as f2:
        localNodeId = pickle.load(f2)

    maximumConnectionCount = 50
    # When maximum connection count is hit, I don't know how to say it to the remote.
    currentConnectionCount = 0

    def clientConnectionFailed(self, connector, reason):
        print('CONNECTION FAILED because: %s' %reason.getErrorMessage())


"""
    The API used for connecting to the nodes, from other parts of the application.
    The two methods below are usually all you need.
"""
from twisted.internet.ssl import ClientContextFactory
def connectWithNode(node):
    ip = node['LastConnectedIP'] if node['LastConnectedIP'] != None else node['LastRetrievedIP']
    port = node['LastConnectedPort'] if node['LastConnectedPort'] != None else node['LastRetrievedPort']
    endpoint = SSL4ClientEndpoint(reactor, ip , port, globals.AetherClientContextFactory())
    endpoint.connect(AetherProtocolFactory())\
        .addCallback(lambda p: p.initiateHandshake())\
        .addErrback(printError,
                'Connection failed because node  %s at the address %s:%s is not responding to or actively '
                'refusing the connection request. It can be offline or permanently dead.'
                %(node['NodeId'], ip, port))

# Errbacks used in API above.

def printError(exception, definition):
    print('\n'
          'NETWORK FAILURE: '
          '%s\n'
          'TRACEBACK: %s'
          '\n' %(definition, exception))


def connectWithIP(IP, Port):
    cprint('CONNECT WITH IP IS CALLED FOR THE ADDRESS %s:%s' %(IP, Port), 'white', 'on_red')
    endpoint = SSL4ClientEndpoint(reactor, IP , Port, globals.AetherClientContextFactory())
    endpoint.connect(AetherProtocolFactory())\
        .addCallback(lambda p: p.initiateHandshake())\
        .addErrback(printError,
            'Connection failed because node at the address %s:%s is not responding to or actively '
            'refusing the connection request. It can be offline or permanently dead.'
            %(IP, Port))