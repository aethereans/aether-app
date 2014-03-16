from __future__ import print_function
from ORM import Demeter
from InputOutput.aetherProtocol import connectWithNode, connectWithIP
import globals

if not globals.debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass

firstConnectTriggered = False

def marduk(factoryInstance, committerInstance):
    if globals.appIsPaused: # If the app is explicitly paused, don't create active network events.
        return

    d3 = Demeter.getNodesToConnect(globals.maxOutboundCount, globals.cooldown)
    print('max outbound count:', globals.maxOutboundCount, 'cooldown:', globals.cooldown)
    def connectToNodes(nodes, openConnsList):
        if len(nodes) == 0:
            # If this is a bootstrap.
            global firstConnectTriggered
            if not firstConnectTriggered and (globals.newborn or globals.nuked or globals.resetted):
                firstConnectTriggered = True
                print('We\'re bootstrapping!')
                connectWithIP('151.236.11.192', 39994) # austria #FIXME 39994
            else:
                print('There are no nodes to connect to!')
                print('Number of open connections:', len(openConnsList), 'Commit in Progress: ', committerInstance.commitInProgress)
                if not len(openConnsList) and not committerInstance.commitInProgress:
                    connectWithIP('151.236.11.192', 39994) # austria #FIXME 39994

        for n in nodes:
             print('I\'m attempting to connect to node %s at %s:%s'
                   %(n.NodeId,
                     n.LastConnectedIP if n.LastConnectedIP is not None else n.LastRetrievedIP,
                     n.LastConnectedPort if n.LastConnectedPort is not None else n.LastRetrievedPort))
             connectWithNode(n.asDict())
    d3.addCallback(connectToNodes, factoryInstance.openConnections)


