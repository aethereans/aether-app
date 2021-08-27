// Services > Global Constants

// var ipc = require('../../../../node_modules/electron-better-ipc').ipcMain
// const fesupervisor = require('../fesupervisor/fesupervisor')

interface Globals {
  FrontendReady: boolean
  FrontendAPIPort: number
  FrontendClientConnInitialised: boolean
  ClientAPIServerPort: number
  FrontendDaemonStarted: boolean
  RendererReady: boolean
  FrontendDaemon: any
  AppIsShuttingDown: boolean
}

let glob: Globals = {
  FrontendReady: false,
  FrontendAPIPort: 0,
  FrontendClientConnInitialised: false,
  ClientAPIServerPort: 0,
  FrontendDaemonStarted: false,
  RendererReady: false,
  FrontendDaemon: {},
  AppIsShuttingDown: false,
}

module.exports = glob
