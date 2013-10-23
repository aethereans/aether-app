from __future__ import print_function


import sys, webbrowser
from PyQt5.QtCore import *
from PyQt5.QtWebKitWidgets import *
from PyQt5.QtWidgets import QApplication
from twisted.internet.task import LoopingCall
from twisted.internet.endpoints import SSL4ServerEndpoint

from datetime import datetime, timedelta


app = QApplication(sys.argv)

from globals import baseURL, enableWebkitInspector, \
    aetherListeningPort, newborn, resetted, nuked, basedir

import globals

if globals.FROZEN:
    del sys.modules['twisted.internet.reactor']

import qt5reactor
qt5reactor.install()
from twisted.internet import reactor
from twisted.web.client import getPage

from InputOutput import aetherProtocol
from DecisionEngine import eventLoop
from InputOutput.Charon import Charon
from ORM import Demeter

from PyQt5.QtWidgets import QSystemTrayIcon, QWidget, QMenu, QSplashScreen
from PyQt5 import QtGui
from PyQt5 import QtPrintSupport
from PyQt5 import QtCore

if len(sys.argv) > 1 and sys.argv[1] == '-openatlogin':
    globals.appStartedAtBoot = True

modifierKeypressed = False
lastModifierKeypressDatetime = ''

class ModifierEventFilter(QObject):
    def eventFilter(self, receiver, event):
        if (event.type() == QEvent.KeyPress):
            if (event.modifiers() == QtCore.Qt.ControlModifier and event.key() == QtCore.Qt.Key_W):
                view.hide()
                toggleVisibilityMenuItem.setText('Show Aether')

            elif (event.key() == QtCore.Qt.Key_Control): # Control is Cmd in Mac.
                global modifierKeypressed
                modifierKeypressed = True
                global lastModifierKeypressDatetime
                lastModifierKeypressDatetime = datetime.now()
                return True
            else:
                #Call Base Class Method to Continue Normal Event Processing
                return super(ModifierEventFilter,self).eventFilter(receiver, event)
        else:
            #Call Base Class Method to Continue Normal Event Processing
            return super(ModifierEventFilter,self).eventFilter(receiver, event)

ef = ModifierEventFilter()
app.installEventFilter(ef)
app.setAttribute(QtCore.Qt.AA_UseHighDpiPixmaps)


class aetherMainWindow(QWebView):
    def __init__(self, reactor):
        super(aetherMainWindow, self).__init__()
        self.reactor = reactor
        self.initUI()

    def initUI(self):
        self.resize(1148, 680)
        self.move(126,300)
        self.setContextMenuPolicy(Qt.NoContextMenu)
        self.load(QUrl('file:///'+baseURL+'GUI/WebKitApp/index.html'))
        self.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)

        def linkClick(url):
            webbrowser.open(str(url.toString()))

        self.linkClicked.connect(linkClick)

    def closeEvent(self, e):
        if modifierKeypressed:
            delta = datetime.now() - lastModifierKeypressDatetime
            if delta.microseconds/1000 < 300: # Microseconds / 1000 = milliseconds.
            # This is CMD+Q. I have to handle this because Alt+F4 and Cmd+Q are special and cannot be caught by Qt.
            # See the reference here: http://qt-project.org/faq/answer/how_can_i_catch_altf4_in_my_qt_application
                if self.reactor.threadpool is not None:
                    self.reactor.threadpool.stop()
                self.close()
                self.reactor.stop()
                app.quit()
                sys.exit()
        else:
            self.hide()
            toggleVisibilityMenuItem.setText('Show Aether')

class SystemTrayIcon(QSystemTrayIcon):
    def __init__(self, activateCallback, parent=None):
        if app.devicePixelRatio() == 2:
            icon = QtGui.QIcon(basedir+'Assets/aether-black-tray.svg')
            iconActive = QtGui.QIcon(basedir+'Assets/aether-white-tray.svg')
        else:
            icon = QtGui.QIcon(basedir+'Assets/aether-black-tray.png')
            iconActive = QtGui.QIcon(basedir+'Assets/aether-white-tray.png')
        QSystemTrayIcon.__init__(self, icon, parent)
        menu = QMenu(parent)
        if globals.appIsPaused:
            menu.addAction('Aether is paused.').setDisabled(True)
        else:
            menu.addAction('Aether is connected.').setDisabled(True)
        globalStatusMenuItem = menu.actions()[0]

        menu.addSeparator() # 1
        if globals.appIsPaused:
            menu.addAction('Resume')
        else:
            menu.addAction('Pause')
        togglePauseMenuItem = menu.actions()[2]
        def togglePause():
            if globals.appIsPaused:
                globals.appIsPaused = False
                togglePauseMenuItem.setText('Pause')
                globalStatusMenuItem.setText('Aether is connected.')
            else:
                globals.appIsPaused = True
                togglePauseMenuItem.setText('Resume')
                globalStatusMenuItem.setText('Aether is paused.')
        togglePauseMenuItem.triggered.connect(togglePause)

        if not globals.appStartedAtBoot:
            menu.addAction('Show Aether')
        else:
            menu.addAction('Hide Aether')
        global toggleVisibilityMenuItem
        toggleVisibilityMenuItem = menu.actions()[3]
        def toggleVisibility():
            if parent.isHidden():
                parent.show()
                toggleVisibilityMenuItem.setText('Hide Aether')

            else:
                parent.close()
                toggleVisibilityMenuItem.setText('Show Aether')
        toggleVisibilityMenuItem.triggered.connect(toggleVisibility)

        menu.addAction('Email the developer')
        emailDevMenuItem = menu.actions()[4]
        def emailDev():
            mailInitialiser = \
                QUrl('mailto:burak@nehbit.net'
                     '?subject=Feedback for Aether'
                     '&body=<i><br><br>Hello there! Thanks for taking time to give feedback, I really appreciate it. '
                     'If you are having problems, please right click on Aether.app on your Applications folder, '
                     'click Show Package Contents, go to Contents/MacOS/Logs and attach the network.log file there to '
                     'this email. <br><br>'
                     'You can delete this text before sending.'
                     '<br><br>You can find my PGP key here:</i> '
                     'http://pgp.mit.edu:11371/pks/lookup?search=Burak+Nehbit')
            QtGui.QDesktopServices.openUrl(mailInitialiser)
        emailDevMenuItem.triggered.connect(emailDev)

        menu.addSeparator() # 5

        menu.addAction('Quit')
        quitAppMenuItem = menu.actions()[6]
        # This is below reactor.run to allow access from other places outside main.
        def quitApp():
            # This is buggy...
            if parent.reactor.threadpool is not None:
                parent.reactor.threadpool.stop()
            parent.close()
            parent.reactor.stop()
            app.quit()
            sys.exit()
        quitAppMenuItem.triggered.connect(quitApp)

        self.setContextMenu(menu)
        self.setIcon(icon)
        self.activated.connect(lambda r: activateCallback(self, r))
        def changeIconToActiveState():
            self.setIcon(iconActive)
        def changeIconToPassiveState():
            self.setIcon(icon)
        menu.aboutToShow.connect(changeIconToActiveState)
        menu.aboutToHide.connect(changeIconToPassiveState)


def main():
    from time import sleep
    charon = Charon()
    global view
    view = aetherMainWindow(reactor)
    view.page().mainFrame().addToJavaScriptWindowObject("Charon", charon)
    view.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)

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


    def iconClicked(self, reason):
        pass

    trayIcon = SystemTrayIcon(iconClicked, view)
    trayIcon.show()

    d = Demeter.checkUPNPStatus(10000)
    def upnpSanityCheck(returnVal):
        if returnVal is False:
            print('Oops, the network could not map in 10 seconds. Trying a minute.')
            Demeter.checkUPNPStatus(60000)
            # If that also fails, we're in hostile territory.
    d.addCallback(upnpSanityCheck)

    persephone = LoopingCall(eventLoop.persephone)
    persephone.start(10) # this should be 60 under normal circumstances.
    marduk = LoopingCall(eventLoop.marduk)
    marduk.start(60) # 5 minutes normally, which is 300

    logSender = LoopingCall(eventLoop.sendLogs, reactor)
    logSender.start(86400) # send it every day.

    #enableWebkitInspector = True
    if enableWebkitInspector is True:
        from PyQt5.QtWebKit import QWebSettings
        QWebSettings.globalSettings().setAttribute(QWebSettings.DeveloperExtrasEnabled, True)
        inspect = QWebInspector()
        inspect.resize(1450, 300)
        inspect.move(0,0)
        inspect.setPage(view.page())
        view.setContextMenuPolicy(Qt.DefaultContextMenu)
        inspect.show()

    listenerEndpoint = SSL4ServerEndpoint(reactor, aetherListeningPort, globals.AetherContextFactory())
    listenerEndpoint.listen(aetherProtocol.AetherProtocolFactory())

    splash.finish(view)
    if not globals.appStartedAtBoot:
        view.show()
        toggleVisibilityMenuItem.setText('Hide Aether')
    reactor.run()


if __name__ == "__main__":

    if app.devicePixelRatio() == 2:
        splashPicture = QtGui.QPixmap(basedir+'Assets/splash@2x.png')
        sp2 = splashPicture.copy(QtCore.QRect(0,0, 450,450)) # This is a workaround for a bug in Qt 5.1.1

    else:
        sp2 = QtGui.QPixmap(basedir+'Assets/splash.png')

    splash = QSplashScreen(sp2, Qt.WindowStaysOnTopHint)

    if not globals.appStartedAtBoot:
        splash.show()
    main()

