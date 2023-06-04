'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
var globals = require('../globals/globals')
var fesupervisor = require('../fesupervisor/fesupervisor')
var ipc = require('../../../../node_modules/electron-better-ipc')
/*----------  Main receivers  ----------*/
// i.e. main does something at the request of renderer
ipc.answerRenderer('GetFrontendReady', function () {
  return globals.FrontendReady
})
ipc.answerRenderer('SetFrontendReady', function (ready) {
  globals.FrontendReady = ready
})
ipc.answerRenderer('SetRendererReady', function (ready) {
  globals.RendererReady = ready
})
ipc.answerRenderer('GetFrontendAPIPort', function () {
  return globals.FrontendAPIPort
})
ipc.answerRenderer('SetFrontendAPIPort', function (port) {
  globals.FrontendAPIPort = port
})
ipc.answerRenderer('GetFrontendClientConnInitialised', function () {
  return globals.FrontendClientConnInitialised
})
ipc.answerRenderer('SetFrontendClientConnInitialised', function (initialised) {
  globals.FrontendClientConnInitialised = initialised
})
ipc.answerRenderer('GetClientAPIServerPort', function () {
  return globals.ClientAPIServerPort
})
ipc.answerRenderer('SetClientAPIServerPort', function (port) {
  console.log('ipc client api server port: ', port)
  globals.ClientAPIServerPort = port
  return fesupervisor.StartFrontendDaemon(globals.ClientAPIServerPort)
})
// module.exports = ipc
//# sourceMappingURL=eipc-main.js.map
