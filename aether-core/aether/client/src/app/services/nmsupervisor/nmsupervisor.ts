// Services > NameMinter Supervisor
// This service handles the interface between the GUI app and the name minter.

export {}

const { exec } = require('child_process')

// let os = require('os')
// let path = require('path')

let MintNewUniqueUsername = function (
  requestedUsername: string,
  targetKeyFp: string,
  expiryTimestamp: number,
  password: string,
  callback: any
) {
  let execString = `go run ../../support/nameminter/main.go mint`
  execString += ` --reqname="`
  execString += requestedUsername
  execString += `" --targetkeyfp="`
  execString += targetKeyFp
  execString += `" --expiry="`
  execString += expiryTimestamp
  execString += `" --password="`
  execString += password
  execString += `"`

  exec(execString, function (e: any, stdout: any) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}

let FetchAlreadyMintedPendingUsernames = function (callback: any) {
  let execString = `go run ../../support/nameminter/main.go batchdeliver`
  exec(execString, function (e: any, stdout: any) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}

let MarkUsernamesAsDelivered = function (
  deliveredUsernames: any,
  callback: any
) {
  let execString = `go run ../../support/nameminter/main.go markdelivered`
  execString += ` --deliveredfps='`
  execString += JSON.stringify(deliveredUsernames)
  execString += `'`
  exec(execString, function (e: any, stdout: any) {
    // , stderr: any
    if (e instanceof Error) {
      callback(e.message)
      return
    }
    callback(stdout)
  })
}

module.exports = {
  MintNewUniqueUsername,
  FetchAlreadyMintedPendingUsernames,
  MarkUsernamesAsDelivered,
}
