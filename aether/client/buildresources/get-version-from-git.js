/*
  This script gets the git tag that is the version number, adds the commit hash, and updates the package.json version for the app.

  Heads up - this script (alongside all other scripts ) assumes an UNIX-compatible operating system.
*/

const resolve = require('resolve')
const execSync = require('child_process').execSync
const writePkg = require('write-pkg')
const readPkgUp = require('read-pkg-up')

async function main() {
  let userCmd = process.argv[2]
  if (typeof userCmd === 'undefined' || userCmd === 'useprior') {
    if (userCmd === 'useprior') {
      console.log(process.argv[3])
      return
    }
    /*
            Default path. Compiles the version and updates the package.json.
        */
    const packageJsonPath = resolve.sync('../package.json')
    let packageJsonPkg = await readPkgUp()
    let pj = packageJsonPkg.pkg
    let r = compileFullVersion()
    pj.version = r
    await writePkg(packageJsonPath, pj)
    return
  }
  if (userCmd === 'print') {
    /*
            Generates and outputs the full version (v+build) only to stdout.
        */
    console.log(compileFullVersion())
    return
  }
  if (userCmd === 'print-version-only') {
    /*
            Generates and outputs the version only (without build) to stdout.
        */
    console.log(getBaseVersion())
    return
  }
  if (userCmd === 'print-build-only') {
    /*
            Generates and outputs the build only (without version) to stdout.
        */
    console.log(getBaseBuildNumber())
    return
  }
  console.log('The command you gave was not understood. You gave:', userCmd)
}

function getBaseVersion() {
  return execSync('printf "%s" `git describe --abbrev=0 --tags`').toString()
}

function getBaseBuildNumber() {
  let humanTimestamp = generateTimestamp()
  let commitHash = execSync('printf "%s" `git rev-parse --short HEAD`')
  let dirty = gitIsDirty()
  let str = humanTimestamp + '.' + commitHash
  if (dirty) {
    str = str + '.d'
  }
  return str
}

function compileFullVersion() {
  let version = getBaseVersion()
  let str = getBaseBuildNumber()
  return version + '+' + str
}

function gitIsDirty() {
  /*
        Check if we have any uncommitted changes at the moment of compile.
    */
  let isDirty = execSync(
    `git diff-index --quiet HEAD -- ':!*package-lock.json' ':!*package.json' ':!*buildresources/get-version-from-git.js' ':!../support/getaether-website' || echo 'dirty'`
  )
  return isDirty.toString()
}

function generateTimestamp() {
  var now = new Date()
  // 201809251643
  let year = now
    .getFullYear()
    .toString()
    .substr(-2)
  let month = now.getMonth() + 1
  if (month < 10) {
    month = '0' + month
  }
  let day = now.getDate()
  if (day < 10) {
    day = '0' + day
  }
  let hour = now.getHours()
  if (hour < 10) {
    hour = '0' + hour
  }
  let min = now.getMinutes()
  if (min < 10) {
    min = '0' + min
  }
  let str = year + '' + month + '' + day + '' + hour + '' + min
  return str
}

main()
