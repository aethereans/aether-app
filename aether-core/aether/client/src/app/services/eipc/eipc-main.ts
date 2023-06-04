export {}

let globals = require('../globals/globals')
let fesupervisor = require('../fesupervisor/fesupervisor')
let ipc = require('../../../../node_modules/electron-better-ipc')

/*----------  Main receivers  ----------*/
// i.e. main does something at the request of renderer

ipc.answerRenderer('GetFrontendReady', function (): boolean {
  return globals.FrontendReady
})

ipc.answerRenderer('SetFrontendReady', function (ready: boolean) {
  globals.FrontendReady = ready
})

ipc.answerRenderer('SetRendererReady', function (ready: boolean) {
  globals.RendererReady = ready
})

ipc.answerRenderer('GetFrontendAPIPort', function (): number {
  return globals.FrontendAPIPort
})

ipc.answerRenderer('SetFrontendAPIPort', function (port: number) {
  globals.FrontendAPIPort = port
})

ipc.answerRenderer('GetFrontendClientConnInitialised', function (): boolean {
  return globals.FrontendClientConnInitialised
})

ipc.answerRenderer(
  'SetFrontendClientConnInitialised',
  function (initialised: boolean) {
    globals.FrontendClientConnInitialised = initialised
  }
)

ipc.answerRenderer('GetClientAPIServerPort', function (): number {
  return globals.ClientAPIServerPort
})

ipc.answerRenderer('SetClientAPIServerPort', function (port: number): boolean {
  console.log('ipc client api server port: ', port)
  globals.ClientAPIServerPort = port
  return fesupervisor.StartFrontendDaemon(globals.ClientAPIServerPort)
})

// module.exports = ipc
