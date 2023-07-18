'use strict'
// Services > NameMinter Supervisor
// This service handles the interface between the GUI app and the name minter.
Object.defineProperty(exports, '__esModule', { value: true })
var exec = require('child_process').exec
// let os = require('os')
// let path = require('path')
var MintNewUniqueUsername = function (
  requestedUsername,
  targetKeyFp,
  expiryTimestamp,
  password,
  callback
) {
  var execString = 'go run ../../support/nameminter/main.go mint'
  execString += ' --reqname="'
  execString += requestedUsername
  execString += '" --targetkeyfp="'
  execString += targetKeyFp
  execString += '" --expiry="'
  execString += expiryTimestamp
  execString += '" --password="'
  execString += password
  execString += '"'
  exec(execString, function (e, stdout) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}
var FetchAlreadyMintedPendingUsernames = function (callback) {
  var execString = 'go run ../../support/nameminter/main.go batchdeliver'
  exec(execString, function (e, stdout) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}
var MarkUsernamesAsDelivered = function (deliveredUsernames, callback) {
  var execString = 'go run ../../support/nameminter/main.go markdelivered'
  execString += " --deliveredfps='"
  execString += JSON.stringify(deliveredUsernames)
  execString += "'"
  exec(execString, function (e, stdout) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}
module.exports = {
  MintNewUniqueUsername: MintNewUniqueUsername,
  FetchAlreadyMintedPendingUsernames: FetchAlreadyMintedPendingUsernames,
  MarkUsernamesAsDelivered: MarkUsernamesAsDelivered,
}
//# sourceMappingURL=nmsupervisor.js.map
