'use strict'
// Starters for the client.
Object.defineProperty(exports, '__esModule', { value: true })
var clapiserver = require('./services/clapiserver/clapiserver')
var spawn = require('child_process').spawn
var clientAPIServerPort
var clientAPIServerIP = '127.0.0.1'
/**
  This is a documentation.
*/
function StartGRPCServer() {
  console.log('Starting the Client API Server.')
  clientAPIServerPort = clapiserver.StartClientAPIServer()
}
exports.StartGRPCServer = StartGRPCServer
function StartFE() {
  // This is where we start the frontend binary.
  console.log('Frontend daemon starting')
  var child = spawn(
    'go',
    [
      'run',
      '../frontend/main.go',
      'run',
      '--logginglevel=1',
      '--clientip=' + clientAPIServerIP,
      '--clientport=' + clientAPIServerPort,
    ],
    {
      // env: {}, // no access to environment, enabled this in prod to make sure that the app can run w/out depending on anything
      detached: true,
    }
  )
  // child.unref() // Unreference = means it can continue running even when client shuts down. todo: figure out how to make best use of this, we want the frontend to shut down but maybe not the backend? do we want client to have code that searches for an existing fe?
  child.on('exit', function (code, signal) {
    console.log(
      'Frontend process exited with ' +
        ('code ' + code + ' and signal ' + signal)
    )
  })
  child.stdout.on('data', function (data) {
    console.log('' + data)
  })
  child.stderr.on('data', function (data) {
    console.error('' + data)
  })
}
exports.StartFE = StartFE
//# sourceMappingURL=starters.js.map
