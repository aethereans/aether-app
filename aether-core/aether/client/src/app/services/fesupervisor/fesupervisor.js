'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
var spawn = require('child_process').spawn
var isDev = require('electron-is-dev')
var globals = require('../globals/globals')
var clientAPIServerIP = '127.0.0.1'
var os = require('os')
var path = require('path')
var StartFrontendDaemon = function (clientAPIServerPort) {
  if (globals.FrontendDaemonStarted) {
    console.log('frontend daemon already running. skipping the start.')
    return false
  }
  globals.FrontendDaemonStarted = true
  // This is where we start the frontend binary.
  console.log('Frontend daemon starting')
  var child = {}
  if (!isDev) {
    // In production
    var fepath = path.join(
      __dirname,
      '../../../../../app',
      'aether-frontend-' + getOsAndArch()
    )
    console.log('fepath')
    console.log(fepath)
    child = spawn(
      fepath,
      [
        'run',
        '--isdev=false',
        '--clientip=' + clientAPIServerIP,
        '--clientport=' + clientAPIServerPort,
      ],
      {
        // env: {}, // no access to environment, enabled this in prod to make sure that the app can run w/out depending on anything
        detached: false,
      }
    )
  } else {
    // In development
    var compilerTags = ''
    /*
          {{ COMPILE INSTRUCTIONS }}
          To run in extvenabled in development, you need to comment out the line below
        */
    compilerTags = 'extvenabled'
    // ^^^^^ This line
    // todo
    child = spawn(
      'go',
      [
        'run',
        '-tags',
        compilerTags,
        '../frontend/main.go',
        'run',
        '--isdev=true',
        '--clientip=' + clientAPIServerIP,
        '--clientport=' + clientAPIServerPort,
      ],
      {
        // , '--logginglevel=1'
        // env: {}, // no access to environment, enabled this in prod to make sure that the app can run w/out depending on anything
        detached: false,
      }
    )
    // console.log(child)
  }
  globals.FrontendDaemon = child
  /*
      What's below needs to within this function, because this needs to all be set whenever the FE daemon is assigned. If you bring it outside, it attempts to set it at the very beginning, when it is actually undefined.
    */
  // child.unref() // Unreference = means it can continue running even when client shuts down. todo: figure out how to make best use of this, we want the frontend to shut down but maybe not the backend? do we want client to have code that searches for an existing fe?
  globals.FrontendDaemon.on('exit', function (code, signal) {
    if (globals.AppIsShuttingDown) {
      // THis is not a crash, it is a shut down. This does not apply.
      return
    }
    globals.FrontendDaemonStarted = false
    console.log(
      'Frontend process exited with ' +
        ('code ' + code + ' and signal ' + signal)
    )
    console.log('We will reattempt to start the frontend daemon in 10 seconds.')
    setTimeout(function () {
      console.log('Attempting to restart the frontend now.')
      console.log(globals.ClientAPIServerPort)
      console.log(globals)
      StartFrontendDaemon(globals.ClientAPIServerPort)
    }, 10000)
  })
  globals.FrontendDaemon.stdout.on('data', function (data) {
    console.log('' + data)
  })
  globals.FrontendDaemon.stderr.on('data', function (data) {
    console.error('' + data)
  })
  return true
}
function getOsAndArch() {
  var opSys = os.platform()
  var arch = os.arch()
  // Valid OS values are: win32, linux, darwin
  // Mapped to            win,   linux, mac
  // (win returns win32 even on 64 bit on node)
  if (opSys === 'win32') {
    opSys = 'win'
  }
  if (opSys === 'darwin') {
    opSys = 'mac'
  }
  // Valid archs are: x64, ia32, arm64, arm
  // Mapped to:       x64, ia32, arm64, arm32
  if (arch === 'arm') {
    arch = 'arm32'
  }
  return opSys + '-' + arch
}
module.exports = {
  StartFrontendDaemon: StartFrontendDaemon,
}
//# sourceMappingURL=fesupervisor.js.map
