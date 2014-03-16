from interprocessAPI import *
from ampoule import child
from ORM import Hermes
from globals import basedir

# This is the code responding to the protocol in the child.
class GUIProcessProtocol(child.AMPChild): # This is the dude which instantiates the child process.
    # I can override connection lost method and prevent the reactor from shutting down, keeping
    # the process alive.

    def killApp(self):
        print('sending kill signal')
        return self.callRemote(killApp)

    def commit(self, PostFingerprint):
        print('commit is called to parent')
        return self.callRemote(commit, PostFingerprint=PostFingerprint)

    def connectWithIP(self, IP, Port):
        return self.callRemote(connectWithIP, IP=IP, Port=int(Port))

    # These are methods arriving from Main thread to GUI thread (this thread.)

    @thereAreReplies.responder
    def respondToThereAreReplies(self):
        # Do stuff here.
        print('I received a NEW REPLIES signal')
        # Here I need to check the reply count and post that reply count to the main page by changing the scope element.
        replyCount = self.Hermes.countReplies()
        from PyQt5 import QtWidgets
        self.trayIcon.showMessage('New Messages', 'You have new messages.', QtWidgets.QSystemTrayIcon.Information)
        self.trayIcon.lightUpIcon()
        # This doesn't seem to be working. It might be that it's running attached to PyCharm, or it might legitimately be broken. TEST
        jsString = \
            ("rootScope = angular.element(document.getElementById('root-body')).scope();"
             "rootScope.totalReplyCount = %s;"
             "rootScope.$apply();" % replyCount
            )
        self.JSContext(jsString)
        return {}


    # def connectionLost(self, reason):
    #     #do nothing?
    #     pass
    #     # amp.AMP.connectionLost(self, reason)
    #     # from twisted.internet import reactor
    #     # try:
    #     #     reactor.stop()
    #     # except error.ReactorNotRunning:
    #     #     # woa, this means that something bad happened,
    #     #     # most probably we received a SIGINT. Now this is only
    #     #     # a problem when you use Ctrl+C to stop the main process
    #     #     # because it would send the SIGINT to child processes too.
    #     #     # In all other cases receiving a SIGINT here would be an
    #     #     # error condition and correctly restarted. maybe we should
    #     #     # use sigprocmask?
    #     #     pass
    #     # if not self.shutdown:
    #     #     # if the shutdown wasn't explicit we presume that it's an
    #     #     # error condition and thus we return a -1 error returncode.
    #     #     import os
    #     #     os._exit(-1)

childBootstrap = """\
# Future imports
from __future__ import print_function
import globals

from twisted.python import log
from twisted.python.logfile import DailyLogFile
log.startLogging(DailyLogFile.fromFullPath(globals.basedir+'Logs/network.log'))

# Aether modules
from GUI.guiElements import *
from globals import enableWebkitInspector, basedir

from ORM import Hermes

app = Aether()

#Python imports
import os

if globals.FROZEN:
    del sys.modules['twisted.internet.reactor']

import qt5reactor
qt5reactor.install()
from twisted.internet import reactor


# Aether modules that need to load after the reactor
from InputOutput.Charon import Charon

baseurl = os.getcwd()+'/'

if len(sys.argv) > 1 and sys.argv[1] == '-openatlogin':
    globals.appStartedAtBoot = True

# if globals.PLATFORM == 'OSX':
#     import  objc
#     import AppKit
#
#     if globals.FROZEN:
#         class NSObjCApp(AppKit.AppKit.NSObject):
#             @objc.signature('B@:#B')
#             def applicationShouldHandleReopen_hasVisibleWindows_(self, nsAppObject, flag):
#                 app.onClickOnDock()
#                 return True
#     else:
#         class NSObjCApp(AppKit.NSObject):
#             @objc.signature('B@:#B')
#             def applicationShouldHandleReopen_hasVisibleWindows_(self, nsAppObject, flag):
#                 app.onClickOnDock()
#                 return True
#
#     cls = objc.lookUpClass('NSApplication')
#     appInstance = cls.sharedApplication() # I'm doing some real crazy runtime shit there.
#     ta = NSObjCApp.alloc().init()
#     appInstance.setDelegate_(ta)

if globals.PLATFORM == 'LNX':
    app.setStyle("Fusion")
    # This is a fix to 'CRITICAL: GTK_IS_WIDGET (widget)' failed' bug on Debian / Ubuntu.
    # Visually it doesn't change anything, it's just an explicit declaration to help Unity.

def main(reactorString, ampChildPath):
#def main():


    from twisted.internet import reactor, stdio
    from twisted.python import reflect, runtime
    ampChild = reflect.namedAny(ampChildPath)
    protInstance = ampChild(*sys.argv[1:-2]) # This is how you reach the prot from here.


    from twisted.internet import reactor # this is actually used, ignore the warning.
    hermes = Hermes.Hermes(protInstance)
    charon = Charon(hermes)
    view = AetherMainWindow(charon, reactor, baseurl, protInstance)
    trayIcon = SystemTrayIcon(basedir, app, view)
    trayIcon.protInstance = protInstance
    protInstance.trayIcon = trayIcon
    ef = ModifierEventFilter()
    app.view = view
    ef.view = view
    app.trayIcon = trayIcon
    view.trayIcon = trayIcon
    ef.trayIcon = trayIcon
    charon.trayIcon = trayIcon
    trayIcon.webView = view
    app.installEventFilter(ef)
    trayIcon.show()

    #FIXME before package
    enableWebkitInspector = False
    if enableWebkitInspector is True:
        from PyQt5.QtWebKit import QWebSettings
        QWebSettings.globalSettings().setAttribute(QWebSettings.DeveloperExtrasEnabled, True)
        inspect = QWebInspector()
        inspect.resize(1450, 300)
        inspect.move(0,0)
        inspect.setPage(view.webView.page())
        view.setContextMenuPolicy(Qt.DefaultContextMenu)
        inspect.show()

    splash.finish(view)
    if not globals.appStartedAtBoot:
        view.show()
        trayIcon.toggleVisibilityMenuItem.setText('Hide Aether')



    def sigTest():
        print('hello this is a before shutdown signal')

    reactor.addSystemEventTrigger('before', 'shutdown', sigTest)

    if runtime.platform.isWindows():
        stdio.StandardIO(protInstance)
    else:
        stdio.StandardIO(protInstance, 3, 4)
    enter = getattr(ampChild, '__enter__', None)
    if enter is not None:
        enter()
    try:
        reactor.run()
    except:
        if enter is not None:
            info = sys.exc_info()
            if not ampChild.__exit__(*info):
                raise
        else:
            raise
    else:
        if enter is not None:
            ampChild.__exit__(None, None, None)

if __name__ == "__main__":

    pixmap = QtGui.QPixmap(basedir+'Assets/splash.png')
    splash = QSplashScreen(pixmap, Qt.WindowStaysOnTopHint)

    if not globals.appStartedAtBoot:
        splash.show()
    main(sys.argv[-2], sys.argv[-1])

    if reactor.threadpool is not None:
        reactor.threadpool.stop()
        reactor.stop()
        sys.exit()

"""