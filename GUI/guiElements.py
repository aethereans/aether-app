from PyQt5.QtWidgets import QSystemTrayIcon, QWidget, QMenu, QSplashScreen
from PyQt5 import QtGui
from PyQt5 import QtPrintSupport
from PyQt5 import QtCore
import sys, webbrowser
from PyQt5.QtCore import *
from PyQt5.QtWebKitWidgets import *
from PyQt5.QtWidgets import QApplication, QMainWindow, QFileDialog
import globals
from datetime import datetime

class Aether(QApplication):
    def __init__(self):
        QApplication.__init__(self, sys.argv)
        QApplication.setAttribute(QtCore.Qt.AA_UseHighDpiPixmaps)
        QApplication.setQuitOnLastWindowClosed(False)

    def onClickOnDock(self): # This only gets called on Mac.
        ##print('dock clicked')
        self.view.show()
        self.view.raise_() # because the app needs to fire up in the bg to be brought to front.
        globals.raiseAndFocusApp()
        self.trayIcon.toggleVisibilityMenuItem.setText('Hide Aether')



class ModifierEventFilter(QObject):
    def eventFilter(self, receiver, event):
        if (event.type() == QEvent.KeyPress):
            if (event.modifiers() == QtCore.Qt.ControlModifier and event.key() == QtCore.Qt.Key_W):
                self.view.hide()
                self.trayIcon.toggleVisibilityMenuItem.setText('Show Aether')
                return True

            elif (event.key() == QtCore.Qt.Key_Control): # Control is Cmd in Mac.
                global lastModifierKeypressDatetime
                lastModifierKeypressDatetime = datetime.now()
                return True

            elif (event.key() == QtCore.Qt.Key_Escape): # Control is Cmd in Mac.
                return True # Esc key does not move. Return true stops the propagation there.
            else:
                #Call Base Class Method to Continue Normal Event Processing
                return super(ModifierEventFilter,self).eventFilter(receiver, event)
        #elif (event.type() == QEvent.ApplicationActivate):
        #    ##print('app activate fired')
        #    view.show()
        #    toggleVisibilityMenuItem.setText('Hide Aether')
        #
        #    return True
        else:
            #Call Base Class Method to Continue Normal Event Processing
            return super(ModifierEventFilter,self).eventFilter(receiver, event)

class AetherMainWindow(QMainWindow):
    def __init__(self, charon, reactor, baseurl, protInstance):
        super(AetherMainWindow, self).__init__()
        self.resize(1148, 680)
        self.reactor = reactor
        webView = AetherWebView(reactor, baseurl)
        webView.page().mainFrame().addToJavaScriptWindowObject("Charon", charon)
        webView.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)
        self.setContextMenuPolicy(Qt.NoContextMenu)
        self.setCentralWidget(webView)
        self.webView = webView
        self.protInstance = protInstance
        self.protInstance.JSContext = webView.page().mainFrame().evaluateJavaScript
        self.JSContext = webView.page().mainFrame().evaluateJavaScript
        self.protInstance.Hermes = charon.Hermes
        #self.setWindowFlags(QtCore.Qt.WindowMinimizeButtonHint)

        # from PyQt5.QtWebKit import QWebSettings
        # QWebSettings.globalSettings().setAttribute(QWebSettings.DeveloperExtrasEnabled, True)
        # self.inspector = QWebInspector()
        # self.inspector.resize(1450, 300)
        # self.inspector.move(0,0)
        # self.inspector.setPage(self.webView.page())
        # self.setContextMenuPolicy(Qt.DefaultContextMenu)


    def hideEvent(self, QHideEvent):
        self.trayIcon.toggleVisibilityMenuItem.setText('Show Aether')

class AetherWebView(QWebView):
    def __init__(self, reactor, baseurl):
        super(AetherWebView, self).__init__()
        self.reactor = reactor
        self.load(QUrl('file:///' + baseurl + 'GUI/WebKitApp/index.html'))
        self.page().setLinkDelegationPolicy(QWebPage.DelegateAllLinks)

        def linkClick(url):
            webbrowser.open(str(url.toString()))

        self.linkClicked.connect(linkClick)

class SystemTrayIcon(QSystemTrayIcon):
    def __init__(self, basedir, app, parent=None):
        self.basedir = basedir
        if globals.PLATFORM == 'OSX':
            if app.devicePixelRatio() == 2:
                self.icon = QtGui.QIcon(basedir+'Assets/aether-black-tray.svg')
                self.iconActive = QtGui.QIcon(basedir+'Assets/aether-white-tray.svg')
                self.iconHighlight =  QtGui.QIcon(self.basedir+'Assets/aether-blue-tray.svg')
            else:
                self.icon = QtGui.QIcon(basedir+'Assets/aether-black-tray.png')
                self.iconActive = QtGui.QIcon(basedir+'Assets/aether-white-tray.png')
                self.iconHighlight =  QtGui.QIcon(self.basedir+'Assets/aether-blue-tray.png')
        elif globals.PLATFORM == 'LNX':
            self.icon = QtGui.QIcon(basedir+'Assets/aether-white-tray.png')
            self.iconActive = self.icon
            self.iconHighlight = self.icon
        elif globals.PLATFORM == 'WIN':
            self.icon = QtGui.QIcon(basedir+'Assets/aether-black-tray-win.svg')
            self.iconActive = self.icon
            self.iconHighlight = self.icon
        else:
            pass

        QSystemTrayIcon.__init__(self, self.icon, parent)

        self.menu = QMenu(parent)
        if globals.appIsPaused:
            self.menu.addAction('Aether is paused.').setDisabled(True)
        else:
            self.menu.addAction('Aether is connected.').setDisabled(True)
        globalStatusMenuItem = self.menu.actions()[0]

        self.menu.addSeparator() # 1
        self.menu.addAction('You have no replies.').setDisabled(True)
        self.messagesMenuItem = self.menu.actions()[2]
        def goToMessages():
            self.messagesMenuItem.setText('You have no replies.')
            self.messagesMenuItem.setDisabled(True)
            if parent.isHidden():
                parent.show()
                parent.raise_()
            jsString = \
            ("firstFrameScope = angular.element(document.getElementById('first-frame-contents')).scope();"
             "firstFrameScope.repliesButtonClick();"
             "firstFrameScope.$apply();"
            )
            self.webView.JSContext(jsString)
            # reach out to jscontext and
            # Here, I need to call qtwebkit and tell it to open messages.
        self.messagesMenuItem.triggered.connect(goToMessages)
        self.menu.addSeparator() # 3
        if globals.appIsPaused:
            self.menu.addAction('Resume')
        else:
            self.menu.addAction('Pause')
        self.togglePauseMenuItem = self.menu.actions()[4]
        def togglePause():
            if globals.appIsPaused:
                globals.appIsPaused = False
                self.togglePauseMenuItem.setText('Pause')
                globalStatusMenuItem.setText('Aether is connected.')
            else:
                globals.appIsPaused = True
                self.togglePauseMenuItem.setText('Resume')
                self.globalStatusMenuItem.setText('Aether is paused.')
        self.togglePauseMenuItem.triggered.connect(togglePause)

        if not globals.appStartedAtBoot:
            self.menu.addAction('Show Aether')
        else:
            self.menu.addAction('Hide Aether')
        self.toggleVisibilityMenuItem = self.menu.actions()[5]
        def toggleVisibility():
            if parent.isHidden():
                parent.show()
                parent.raise_()
                # if globals.PLATFORM == 'OSX':
                #     globals.raiseAndFocusApp() #FIXME BEFORE RELEASE
                self.toggleVisibilityMenuItem.setText('Hide Aether')

            else:
                parent.hide()
                parent.lower()
                self.toggleVisibilityMenuItem.setText('Show Aether')
        self.toggleVisibilityMenuItem.triggered.connect(toggleVisibility)

        self.menu.addAction('Email the developer')
        self.emailDevMenuItem = self.menu.actions()[6]
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
        self.emailDevMenuItem.triggered.connect(emailDev)

        self.menu.addSeparator() # 5

        self.menu.addAction('Quit')
        self.quitAppMenuItem = self.menu.actions()[8]
        # This is below reactor.run to allow access from other places outside main.
        def quitApp():
            # This is buggy...
            # if parent.reactor.threadpool is not None:
            #     parent.reactor.threadpool.stop()
            # parent.close()
            print('firing kill app on prot instance')
            self.protInstance.killApp()
            # def finishExit():
            #     parent.reactor.stop()
            #     app.quit()
            #     sys.exit()
            # d.addCallback(finishExit)
        self.quitAppMenuItem.triggered.connect(quitApp)

        self.setContextMenu(self.menu)
        self.setIcon(self.icon)

        def changeIconToActiveState():
            self.setIcon(self.iconActive)
        def changeIconToPassiveState():
            self.setIcon(self.icon)
        self.menu.aboutToShow.connect(changeIconToActiveState)
        self.menu.aboutToHide.connect(changeIconToPassiveState)
        if globals.PLATFORM == 'WIN':
            def showOnLeftClick(reason):
                if reason == self.Trigger:
                    toggleVisibility() # I hate that Python doesn't have anonymous functions.
            self.activated.connect(showOnLeftClick)

    def lightUpIcon(self):
        self.setIcon(self.iconHighlight)
        self.messagesMenuItem.setText('New replies available.')
        self.messagesMenuItem.setDisabled(False)
        self.toggleVisibilityMenuItem.setText('Hide Aether')

    def makeIconGoDark(self):
        self.setIcon(self.icon)
        self.messagesMenuItem.setText('You have no replies.')
        self.messagesMenuItem.setDisabled(True)