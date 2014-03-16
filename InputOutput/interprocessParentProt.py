#from twisted.protocols import amp
from __future__ import print_function
from interprocessAPI import *
from twisted.internet import reactor
from ORM.Demeter import committer, connectToLastConnected
import globals

if not globals.debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass

class MainProcessProtocol(amp.AMP):
    def __init__(self):
        amp.AMP.__init__(self)
        print('Main Process Protocol initialized')
        pass

    @commit.responder
    def guiCommitted(self, PostFingerprint):
        #committer.commit()
        # TODO: Do other stuff required when user adds something.
        # I need to find the last n dudes I connected at, and ignore the cooldown and connect back to them
        # to serve my newest shit.
        # Okay, get the connectToNode method
        committer.newPostsToIncrement.append(PostFingerprint)
        connectToLastConnected(10)
        print('I received a commit signal from child process.')
        return {}

    @killApp.responder
    def killAppResponder(self):
        print('I received a kill signal. KTHXBAI')
        d = self.processPool.stop()
        def stopR(*a):
            reactor.stop()
        d.addCallback(stopR)
        return {}

    @connectWithIP.responder
    def respondToConnectButton(self, IP, Port):
        print('I received a connect request to %s:%s' %(IP, Port))
        return {}

