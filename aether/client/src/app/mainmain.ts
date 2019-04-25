// Electron main.
/*
    This file is the main execution path of Electron. It starts up Electron, and loads our HTML main html file. In this HTML file is our app contained, its JS code, etc.

    In other words, anything here runs as an OS-privileged executable a la Node. Anything that takes root from the main HTML file runs effectively as a web page.
*/

var globals = require('./services/globals/globals') // Register globals
var metrics = require('./services/metrics/metrics')(true, false)

require('./services/eipc/eipc-main') // Register IPC events
var ipc = require('../../node_modules/electron-better-ipc') // Register IPC caller
const elc = require('electron')
// const starters = require('./starters')
// const feapiconsumer = require('./services/feapiconsumer/feapiconsumer')
const minimatch = require('../../node_modules/minimatch')
const treekill = require('tree-kill')

// var ipc = require('../../node_modules/electron-better-ipc')
// const fesupervisor = require('./services/fesupervisor/fesupervisor')

// Enable live reload. This should be disabled in production. TODO
const path = require('path')
const maindir = path.dirname(__dirname)

const isDev = require('electron-is-dev')

// const unhandled = require('../../node_modules/electron-unhandled')
// unhandled()

require('electron-context-menu')({
  // prepend: (params, browserWindow) => [{
  // prepend: () => [{
  //   label: 'Rainbow',
  //   // Only show it when right-clicking images
  //   visible: true
  // }]
  showCopyLink: false,
})

/*===================================
=            Auto update            =
===================================*/

const autoUpdater = require('../../node_modules/electron-updater')
autoUpdater.autoUpdater.requestHeaders = { authorization: '' }

autoUpdater.autoUpdater.on('update-downloaded', function() {
  // ev: any, info: any
  ipc.callRenderer(win, 'NewUpdateReady', true)
})

ipc.answerRenderer('RestartToUpdateTheApp', function() {
  autoUpdater.autoUpdater.quitAndInstall()
  elc.app.quit()
})

ipc.answerRenderer('AskNewUpdateReady', function() {
  checkSoftwareUpdate()
})

function checkSoftwareUpdate() {
  autoUpdater.autoUpdater.checkForUpdates()
}

function startAutoUpdateCheck() {
  checkSoftwareUpdate()
  return setInterval(checkSoftwareUpdate, 3600000) // Check every hour
}

/*=====  End of Auto update  ======*/

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win: any
let tray: any = null
var DOM_READY: boolean = false

// Fix for enabling Windows notifications.
if (process.platform === 'win32') {
  elc.app.setAppUserModelId(process.execPath)
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
elc.app.on('ready', main)

// Quit when all windows are closed.
elc.app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  // if (process.platform !== 'darwin' && process.platform !== 'win32') {
  //   elc.app.quit()
  // }
})

/*==========================================================
=            Retaining window size and location            =
==========================================================*/

const Store = require('electron-store')
const store = new Store()

function getWinBounds() {
  let w = store.get('win_width')
  let h = store.get('win_height')
  if (!w || !h) {
    w = 1200
    h = 740
  }
  return { w: w, h: h }
}

// Call this one on close
function saveWinBounds() {
  let b = win.getBounds()
  store.set('win_width', b.width)
  store.set('win_height', b.height)
}

/*=====  End of Retaining window size and location  ======*/

/*----------  Open / close  ----------*/
// These are shared methods. If there are any changes to these events, they should be implemented here.

let openAppWindow = function() {
  if (win === null) {
    main()
  }
  if (win.isMinimized()) {
    win.restore()
  }
  // In the case it's hidden
  win.show()
  win.focus()
  metrics.SendRaw('App window opened')
}

let closeAppWindow = function() {
  // win = null
  win.hide()
  metrics.SendRaw('App window closed')
}

elc.app.on('activate', () => {
  openAppWindow()
})

elc.app.on('before-quit', function(e: any) {
  contextMenuTemplate[0].label = 'Shutting down...'
  saveWinBounds()
  setTimeout(function() {
    elc.app.exit()
  }, 3000) // Hard limit - if it doesn't shut down in 3 seconds, we force kill it.
  const contextMenu = elc.Menu.buildFromTemplate(contextMenuTemplate)
  tray.setContextMenu(contextMenu)
  // console.log('before-quit: electron is quitting. ')
  globals.AppIsShuttingDown = true
  if (win !== null) {
    win.close()
  }
  e.preventDefault()
  globals.FrontendDaemon.on('exit', function() {
    console.log('Frontend has exited.')
    elc.app.exit()
  })
  treekill(globals.FrontendDaemon.pid)
  metrics.SendRaw('App shut down')
})
// elc.app.on('will-quit', function() {
//   console.log('will-quit: electron is quitting. ')
// })

// Make the app a single-instance app.
// This is actually enforced by Mac OS, but unfortunately Windows isn't that easy - you have to manually implement it.
var shouldQuitDueToAnotherInstanceBeingOpen = elc.app.makeSingleInstance(
  function(argv: any) {
    if (process.platform === 'win32') {
      /*
      This is the windows-specific implementation of the open-url in the will-finish-launching event.

      This code is available in two different places. This one handles deep-linking after the app is open.
    */
      // Keep only command line / deep linked arguments
      if (typeof argv.slice(1)[0] !== 'undefined') {
        linkToLoadAtBoot = argv.slice(1)[0].substring(8)
      }
      if (linkToLoadAtBoot.length > 0) {
        ipc.callRenderer(win, 'RouteTo', linkToLoadAtBoot)
      }
      openAppWindow()
    }
    if (win) {
      if (win.isMinimized()) {
        win.restore()
      }
      win.focus()
    }
  }
)

if (shouldQuitDueToAnotherInstanceBeingOpen) {
  elc.app.exit()
}

function EstablishExternalResourceAutoLoadFiltering() {
  // This is the list of allowed URLs that can be auto-loaded in electron. (This does not prevent links that open in your browser, just ones that fetch data within the app. You can link anywhere, but only links from the whitelist below will have a chance to auto-load.)

  /*
    Why does this list even exist? Shouldn't people be able to link to everywhere?

    You *can* link to everywhere, this list is just for auto-loading previews. Why does this matter? Because when an asset is auto-loaded, the entity on the other end (i.e. the sites below) will see your IP address as a normal user. That means, if there was no whitelist and all links were allowed to auto-load, then anybody could link to a site they control, and by listening to hits from IP addresses, it could figure out which IP addresses are using Aether. It wouldn't be able to figure out who is who, but the fact that IP is using Aether would be revealed.

    If you'd like to make things a little tighter in exchange to not being able to preview, replace this list with an empty one, and all auto-loads will be blocked.
  */

  let autoloadEnabledWhitelist = [
    'https://*.getaether.net/**',
    'https://*.imgur.com/**',
    'https://imgur.com/**',
    'https://gfycat.com/**',
    'https://*.gfycat.com/**',
    'https://giphy.com/**',
    'https://*.giphy.com/**',
    'https://*.mixpanel.com/**',
    'https://*.mxpnl.com/**',
    'https://*.coinbase.com/**',
    'file://**', // So that we can load the local client app itself
    'chrome-devtools://**',
    'chrome-extension://**', // for vue devtools
  ]

  let autoloadDisabledWhitelist = [
    'https://*.getaether.net/**',
    'https://*.mixpanel.com/**',
    'https://*.mxpnl.com/**',
    'file://**',
    'chrome-devtools://**',
    'chrome-extension://**',
  ]

  // This list should be editable. (TODO)
  let whitelist = autoloadEnabledWhitelist

  ipc.answerRenderer('DisableExternalResourceAutoLoad', function() {
    // Only the URLs required for correct app functionality.
    whitelist = autoloadDisabledWhitelist
    return true
  })

  ipc.answerRenderer('EnableExternalResourceAutoLoad', function() {
    // Only the URLs required for correct app functionality.
    whitelist = autoloadEnabledWhitelist
    return true
  })

  // Allow any auto-load request that's in the whitelist. Deny autoload requests to all other domains.
  elc.session.defaultSession.webRequest.onBeforeRequest(function(
    details: any,
    cb: any
  ) {
    // console.log(details.url) // Uncomment this to see all attempted outbound network requests from the client app.
    for (let i = 0; i < whitelist.length; i++) {
      if (minimatch(details.url, whitelist[i], { matchBase: true })) {
        cb({ cancel: false })
        return
      }
    }
    cb({ cancel: true })
  })
}

function EstablishElectronWindow() {
  // If not prod, install Vue devtools.
  if (isDev) {
    require('vue-devtools').install()
  }
  /* Check whether we are started in the hidden state as part of the computer startup process. */
  const loginSettings = elc.app.getLoginItemSettings()
  let hiddenStartAtBoot = loginSettings.wasOpenedAsHidden
  // ^ Only applies to Mac OS
  if (process.platform === 'win32') {
    // ^ Windows version of above
    for (let arg of process.argv) {
      if (arg === '--hidden') {
        hiddenStartAtBoot = true
      }
    }
  }

  // Create the browser window.
  // let dm = elc.screen.getPrimaryDisplay().size
  let bounds = getWinBounds()
  let windowSpec: any = {
    show: false,
    // width: dm.width * 0.8,
    width: bounds.w,
    // height: dm.height * 0.8,
    height: bounds.h,
    title: 'Aether',
    fullscreenWindowTitle: true,
    backgroundColor: '#292b2f',
    disableBlinkFeatures: 'Auxclick', // disable middle click new window
    autoHideMenuBar: true,
    webPreferences: {
      // blinkFeatures: 'OverlayScrollbars'
    },
  }
  if (process.platform === 'win32') {
    windowSpec.frame = false // We have our traffic lights implementation for Win
  }
  if (process.platform === 'darwin') {
    windowSpec.titleBarStyle = 'hiddenInset' // Mac traffic lights
  }
  if (process.platform === 'linux') {
    windowSpec.darkTheme = true // GTK3+ Only
    // Nothing specific for the frame for now.
  }
  win = new elc.BrowserWindow(windowSpec)
  win.once('ready-to-show', function() {
    // We want to show the window only after Electron is done readying itself.
    setTimeout(function() {
      if (!hiddenStartAtBoot) {
        win.show()
      }
    }, 100)
    // Unfortunately, there's a race condition from the Electron side here (I might be making a mistake also, but it is simple enough to reproduce that there is not much space for me to make a mistake). If the setTimeout is 0 or is not present, there's about 1/10 chance the window is painted but completely frozen. Having 100ms seems to make it go away, but it's a little icky, because that basically is my guess. Not great. Hopefully they'll fix this in upcoming Electron 3.
  })
  win.loadFile('index.html')
  if (isDev) {
    // Open the DevTools.
    win.webContents.openDevTools({ mode: 'bottom' })
  }

  // win.webContents.openDevTools({ mode: 'bottom' })

  win.on('close', function(e: any) {
    e.preventDefault()
    // ^ Prevents the app from continuing on with destroying the window element. We need that element.
    closeAppWindow()
    // DOM_READY = false // This is useful when the electron window fully shuts down, not when it's not fully shut down.
  })

  win.webContents.on('dom-ready', function() {
    DOM_READY = true
    // This is needed because the renderer process won't be able to respond to IPC requests before this event happens.

    if (process.platform == 'win32') {
      /*----------  Windows specific deep linker  ----------*/
      /*
        This is the windows-specific implementation of the open-url in the will-finish-launching event.

        This code is available in two different places. This one handles deep-linking from cold boot.
      */
      // Keep only command line / deep linked arguments
      if (typeof process.argv.slice(1)[0] !== 'undefined') {
        linkToLoadAtBoot = process.argv.slice(1)[0].substring(8)
      }
    }
    // Normal open-url event works for Mac and Linux
    if (linkToLoadAtBoot.length > 0) {
      ipc.callRenderer(win, 'RouteTo', linkToLoadAtBoot)
    }
  })
  win.webContents.on('will-navigate', function(e: any, reqUrl: any) {
    e.preventDefault()
    elc.shell.openExternal(reqUrl)
    // return
    // let getHost = function(url: any) { require('url').parse(url).host }
    // let reqHost = getHost(reqUrl)
    // let isExternal = reqHost && reqHost != getHost(win.webContents.getURL())
    // if (isExternal) {
    //   e.preventDefault()
    //   elc.shell.openExternal(reqUrl)
    // }
  })

  win.webContents.on('new-window', function(e: any) {
    e.preventDefault()
  })

  /*----------  Fullscreen state comms to the renderer.  ----------*/
  function sendFullscreenState(isFullscreen: boolean) {
    ipc.callRenderer(win, 'FullscreenState', isFullscreen)
  }
  win.on('enter-full-screen', function() {
    sendFullscreenState(true)
  })
  win.on('leave-full-screen', function() {
    sendFullscreenState(false)
  })
  /*---------- END Fullscreen state comms to the renderer.  ----------*/

  elc.app.on('open-url', function(e: any, url: string) {
    e.preventDefault()
    ipc.callRenderer(win, 'RouteTo', url.substring(8))
    openAppWindow()
  })
}

let linkToLoadAtBoot = ''

elc.app.on('will-finish-launching', function() {
  // Register Aether's aether:// as a standard (http-like) protocol
  elc.protocol.registerStandardSchemes(['aether'])

  elc.app.on('open-url', function(e: any, url: any) {
    e.preventDefault()
    linkToLoadAtBoot = url.substring(8)
  })
})

let openPreferences = function() {
  openAppWindow()
  let rendererReadyChecker = function() {
    if (!(globals.RendererReady && DOM_READY)) {
      return setTimeout(rendererReadyChecker, 100)
    }
    return ipc.callRenderer(win, 'RouteTo', '/settings')
  }
  setTimeout(rendererReadyChecker, 100)
}
let openSupport = function() {
  elc.shell.openExternal('https://meta.getaether.net/c/support')
}

let quitApp = function() {
  elc.app.quit()
}

let contextMenuTemplate = [
  { label: 'Online', enabled: false },
  { type: 'separator' },
  { label: 'Open Aether', click: openAppWindow },
  { type: 'separator' },
  { label: 'Preferences...', click: openPreferences },
  { label: 'Community support', click: openSupport },
  { type: 'separator' },
  { label: 'Quit Aether', click: quitApp },
]

function EstablishTray() {
  if (tray !== null) {
    return
  }
  /*----------  Tray functions  ----------*/

  /*----------  Tray functions END  ----------*/
  let trayIconLocation = ''
  if (process.platform === 'darwin') {
    trayIconLocation = 'ext_dep/images/TrayTemplate.png'
  }
  if (process.platform === 'win32') {
    trayIconLocation = 'ext_dep/images/WindowsTrayIconAlt3.png'
    // trayIconLocation = "ext_dep/images/WindowsTrayIcon.ico"
  }
  if (process.platform === 'linux') {
    trayIconLocation = 'ext_dep/images/LinuxTrayIcon.png'
  }
  tray = new elc.Tray(path.join(__dirname, trayIconLocation))
  tray.setToolTip('Aether')
  const contextMenu = elc.Menu.buildFromTemplate(contextMenuTemplate)
  tray.setContextMenu(contextMenu)
  tray.on('click', () => {
    // On windows, the convention is that when an icon in the tray is clicked, it should spawn the app window.
    if (process.platform === 'win32') {
      openAppWindow()
    }
  })
}

ipc.answerRenderer('QuitApp', function() {
  elc.app.quit()
})

ipc.answerRenderer('FocusAndShow', function() {
  openAppWindow()
  return true
})

let previewBuildStatus = {
  isPreview: false,
  Expiry: 0,
}

// 1540857600 = Oct 30 2018
function EnforcePreviewBuildExpiry() {
  if (!previewBuildStatus.isPreview) {
    return
  }
  // If preview build, check timestamp.
  let d: any = new Date()
  let now: number = Math.floor(d / 1000)
  if (now > previewBuildStatus.Expiry) {
    elc.dialog.showMessageBox(
      {
        type: 'error',
        title: 'Developer build expired',
        message:
          'Hey there! This preview build of Aether has expired. You can get the most recent version of Aether from the meta forum at meta.getaether.net.',
        buttons: ['Quit', 'Get new version'],
      },
      function(respButtonIndex: any) {
        if (respButtonIndex === 1) {
          // The user asked to go to the downloads page.
          elc.shell.openExternal('https://meta.getaether.net')
        }
      }
    )
    // Quit app regardless of what the user chooses.
    elc.app.quit()
  }
}

function ConstructAppMenu() {
  const menu = []
  if (process.platform !== 'linux') {
    // In Linux (i.e. Ubuntu), having this causes two 'Aether' items to show, one of them being the app name, the other being this menu. So this should only be added in the case the OS is not Linux.
    menu.push({
      label: 'Aether',
      submenu: [
        {
          label: 'About Aether',
          click: function() {
            return ipc.callRenderer(win, 'RouteTo', '/about')
          },
        },
        {
          type: 'separator',
        },
        {
          role: 'quit',
        },
      ],
    })
  }

  menu.push({
    label: 'Edit',
    submenu: [
      {
        role: 'undo',
      },
      {
        role: 'redo',
      },
      {
        type: 'separator',
      },
      {
        role: 'cut',
      },
      {
        role: 'copy',
      },
      {
        role: 'paste',
      },
      {
        role: 'selectAll',
      },
    ],
  })
  menu.push({
    label: 'View',
    submenu: [
      {
        role: 'resetZoom',
      },
      {
        role: 'zoomIn',
      },
      {
        role: 'zoomOut',
      },
      {
        role: 'toggleDevTools',
      },
    ],
  })
  if (isDev) {
    menu.push({
      label: 'View',
      submenu: [
        {
          role: 'reload',
        },
        {
          role: 'toggleDevTools',
        },
        {
          type: 'separator',
        },
        {
          role: 'resetZoom',
        },
        {
          role: 'zoomIn',
        },
        {
          role: 'zoomOut',
        },
      ],
    })
  }
  menu.push({
    role: 'windowMenu',
  })
  elc.Menu.setApplicationMenu(elc.Menu.buildFromTemplate(menu))
}

/**
  This is the main() of Electron. It starts the Client GRPC server, and kicks of the frontend and the backend daemons.
*/

function main() {
  startAutoUpdateCheck()
  ConstructAppMenu()
  EstablishExternalResourceAutoLoadFiltering()
  EstablishElectronWindow()
  EstablishTray()
  EnforcePreviewBuildExpiry()
  metrics.SendRaw('App started')
}
