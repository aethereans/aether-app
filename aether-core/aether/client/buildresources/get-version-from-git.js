const execSync = require('child_process').execSync

function main() {
  let version = (getBaseVersion() + '+' + getBaseBuildNumber()).toLowerCase()
  execSync(
    'npm version --allow-same-version --git-tag-version false ' + version
  )
  console.log(version)
}

function gitIsDirty() {
  /*
          Check if we have any uncommitted changes at the moment of compile.
      */
  let isDirty = ""
  
  let isDirty = execSync(
    `git diff-index --quiet HEAD -- . ':!*package-lock.json' ':!*package.json' ':!*buildresources/get-version-from-git.js' ':!../support/getaether-website' || echo 'dirty'`
  )
  
  return isDirty.toString()
}

function getBaseVersion() {
  if(process.platform == "win32") {
    return execSync('git describe --abbrev=0 --tags').toString()
  } else return execSync('printf "%s" `git describe --abbrev=0 --tags`').toString()
}

function getBaseBuildNumber() {
  let humanTimestamp = generateTimestamp()
  let commitHash = ""
  if(process.platform == "win32") {
    let commitHash = execSync('git rev-parse --short HEAD')
  } else {
    let commitHash = execSync('printf "%s" `git rev-parse --short HEAD`')
  }
  let dirty = gitIsDirty()
  let str = humanTimestamp + '.' + commitHash
  if (dirty) {
    str = str + '.d'
  }
  return str
}

function generateTimestamp() {
  var now = new Date()
  let year = now.getFullYear().toString().substr(-2)
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
