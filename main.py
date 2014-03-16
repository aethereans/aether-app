from __future__ import print_function
import sys, datetime

from twisted.internet import reactor
from twisted.web.client import getPage
from twisted.internet.task import LoopingCall
from twisted.internet.endpoints import SSL4ServerEndpoint

from InputOutput import aetherProtocol
from DecisionEngine import eventLoop
from ORM import Demeter
from globals import aetherListeningPort, basedir
import globals

if not globals.debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass

if len(sys.argv) > 1 and sys.argv[1] == '-openatlogin':
    globals.appStartedAtBoot = True

if globals.PLATFORM == 'LNX':
    app.setStyle("Fusion")
    # This is a fix to 'CRITICAL: GTK_IS_WIDGET (widget)' failed' bug on Debian / Ubuntu.
    # Visually it doesn't change anything, it's just an explicit declaration to help Unity.





# figure this out the closure issue. it tries to restart the process.
def shutdown():
    print('shutdown is called')
    pp.stop()
    #reactor.stop()
    # d = pp.stop()
    # d.addCallback(reactor.stop)
    # if reactor.threadpool is not None:
    #     def killReactor(*args):
    #         reactor.threadpool.stop()
    #         reactor.stop()
    #     d.addCallback(killReactor)
    # d.addCallback(reactor.stop)

def main():

    from InputOutput import interprocessChildProt
    from ampoule import pool
    from  InputOutput import interprocessParentProt
    from ampoule import main as ampouleMain

    procStarter = ampouleMain.ProcessStarter(bootstrap=interprocessChildProt.childBootstrap)
    global pp
    pp = pool.ProcessPool(interprocessChildProt.GUIProcessProtocol,
                          ampParent=interprocessParentProt.MainProcessProtocol,
                          starter=procStarter,
                          recycleAfter=0,
                          min=1, max=1)

    pp.start()
    pp.ampParent.processPool = pp # Self referential much?

    Demeter.committer.receiveInterprocessProtocolInstance(pp)

    def checkForUpdates():
        # One catch, any result available out of this will only be visible after next boot of the app.
        d = getPage('http://192.30.33.227')
        def processReceivedVersion(reply):

            if int(reply[:3]) > globals.appVersion:
                globals.updateAvailable = True
                print('There is an update available, local version is %d and gathered version is %s.' %(globals.appVersion, reply))
            else:
                globals.updateAvailable = False
                print('There is no update available')
            globals.setUpdateAvailable(globals.updateAvailable)
        d.addCallback(processReceivedVersion)

    checkForUpdates()

    d = Demeter.checkUPNPStatus(2000)

    def maybeCommit():
        thirtySecondsAgo = datetime.datetime.utcnow() - datetime.timedelta(seconds=10) # FIXME make it 30 in live.
        if Demeter.committer.lastCommit < thirtySecondsAgo and Demeter.committer.commitInProgress is False:
            print('Commit loop decided to commit.')
            Demeter.committer.commit()

    persephone = LoopingCall(maybeCommit)
    persephone.start(10) # this should be 60 under normal circumstances.
    marduk = LoopingCall(eventLoop.marduk, aetherProtocol.aetherProtocolFactoryInstance, Demeter.committer)
    marduk.start(60)
    #FIXME#marduk.start(60) # 5 minutes normally, which is 300

    listenerEndpoint = SSL4ServerEndpoint(reactor, aetherListeningPort, globals.AetherContextFactory())
    listenerEndpoint.listen(aetherProtocol.aetherProtocolFactoryInstance)


    # def checksan():
    #     d = pp.callRemote(interprocessChildProt.checkSanity)
    #     d.addCallback(print)
    # def bootstuff():
    #     d = pp.callRemote(interprocessChildProt.bootGUI)
    #     d.addCallback(print)
    # reactor.callLater(2, bootstuff)
    # reactor.callLater(20, checksan)
    #reactor.callLater(5, aetherProtocol.connectWithIP,'151.236.11.192', 39994) #192 ends
    reactor.run()

if __name__ == "__main__":

    main()
    # Below is the hard-won quit code. So THAT was the reason the quit code always came after loop start...
    # if reactor.threadpool is not None:
    #     reactor.threadpool.stop()
    #     reactor.stop()
    #     sys.exit()