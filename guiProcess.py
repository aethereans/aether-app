# Future imports
from __future__ import print_function
import globals

from twisted.python import log
from twisted.python.logfile import DailyLogFile
#log.startLogging(DailyLogFile.fromFullPath(globals.basedir+'Logs/network.log'))

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
    #ampChild = reflect.namedAny(ampChildPath)
    #protInstance = ampChild(*sys.argv[1:-2]) # This is how you reach the prot from here.
    protInstance = ''

    from twisted.internet import reactor # this is actually used, ignore the warning.
    hermes = Hermes.Hermes(protInstance)
    charon = Charon(hermes)
    view = AetherMainWindow(charon, reactor, baseurl, protInstance)
    trayIcon = SystemTrayIcon(basedir, app, view)
    trayIcon.protInstance = protInstance
    #protInstance.trayIcon = trayIcon
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

    if runtime.platform.isWindows():
        stdio.StandardIO(protInstance)
    else:
        stdio.StandardIO(protInstance, 3, 4)
    #enter = getattr(ampChild, '__enter__', None)
    # if enter is not None:
    #     enter()
    # try:
    #     reactor.run()
    # except:
    #     if enter is not None:
    #         info = sys.exc_info()
    #         if not ampChild.__exit__(*info):
    #             raise
    #     else:
    #         raise
    # else:
    #     if enter is not None:
    #         ampChild.__exit__(None, None, None)
    reactor.run()

if __name__ == "__main__":

    pixmap = QtGui.QPixmap(basedir+'Assets/splash.png')
    splash = QSplashScreen(pixmap, Qt.WindowStaysOnTopHint)

    if not globals.appStartedAtBoot:
        splash.show()
    main('', '')

    if reactor.threadpool is not None:
        reactor.threadpool.stop()
        reactor.stop()
        sys.exit()