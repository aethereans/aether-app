from __future__ import print_function
from ORM import Demeter
from InputOutput.aetherProtocol import connectWithNode, connectWithIP
from globals import basedir, nodeid, appVersion, protocolVersion, appIsPaused
import globals


def marduk():
    if globals.appIsPaused: # If the app is explicitly paused, don't create active network events.
        return

    d3 = Demeter.getNodesToConnect(10, 0)
    def connectToNodes(nodes):
        if len(nodes) == 0 and (globals.newborn or globals.nuked):
        # If there is only the local node, that means no one to connect to.
            print('There are no nodes to connect to!')
            connectWithIP('151.236.11.192', 39994) # austria
        for n in nodes:
             print('I\'m attempting to connect to node %s at %s:%s'
                   %(n.NodeId,
                     n.LastConnectedIP if n.LastConnectedIP is not None else n.LastRetrievedIP,
                     n.LastConnectedPort if n.LastConnectedPort is not None else n.LastRetrievedPort))
             connectWithNode(n.asDict())
    d3.addCallback(connectToNodes)

def persephone():
    Demeter.updatePostStatus()


def sendLogs(reactor):
    """
        This sends the past days' logs to the log server. Because Twisted names past days logs with an added timestamp,
        I can look for that timestamp to accurately determine files as past. It does not matter how many we have, it
        will iterate over all files and sends the bunch.
    """
    import os
    from twisted.python import failure
    from twisted.web.client import Agent
    from twisted.web.http_headers import Headers
    from twisted.web.client import FileBodyProducer

    for file_ in os.listdir(basedir+'Logs'):
        if file_.startswith('network.log.'):
            file_ = basedir + 'Logs/' + file_
            fileToRead = open(file_, 'rb')
            body = FileBodyProducer(fileToRead)
            agent = Agent(reactor)
            d = agent.request(
                'POST', 'http://151.236.11.192:32891',
                Headers({'User-Agent': ['Aether %d : %d' %(appVersion, protocolVersion)],
                         'Content-Type': ['text/log'],
                         'Node-Id':[nodeid],
                         'Log-Timespan':[os.path.splitext(fileToRead.name)[-1][1:]] # get the last of splitext and remove the dot.
                }), body)

            def resultArrived(ignored, fileToRead):
                print('The logs of the day %s was successfully sent to the log server.' %os.path.splitext(fileToRead.name)[-1][1:])
                return fileToRead
            d.addCallback(resultArrived, fileToRead)
            def failureArrived(failure):
                raise Exception('The log server failed to get the logs.')
            d.addErrback(failureArrived)
            def closeFile(fileToRead):
                fileToRead.close()
                return fileToRead
            d.addCallback(closeFile)

            def deleteFile(fileToRead):
                print('file %s is deleted.' %fileToRead.name)
                os.remove(fileToRead.name)
            d.addCallback(deleteFile)
