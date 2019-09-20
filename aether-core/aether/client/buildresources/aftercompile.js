const fs = require('fs')
const path = require('path')

function mkDirByPathSync(targetDir, { isRelativeToScript = false } = {}) {
  const sep = path.sep
  const initDir = path.isAbsolute(targetDir) ? sep : ''
  const baseDir = isRelativeToScript ? __dirname : '.'

  return targetDir.split(sep).reduce((parentDir, childDir) => {
    const curDir = path.resolve(baseDir, parentDir, childDir)
    try {
      fs.mkdirSync(curDir)
    } catch (err) {
      if (err.code === 'EEXIST') {
        // curDir already exists!
        return curDir
      }

      // To avoid `EISDIR` error on Mac and `EACCES`-->`ENOENT` and `EPERM` on Windows.
      if (err.code === 'ENOENT') {
        // Throw the original parentDir error on curDir `ENOENT` failure.
        throw new Error(`EACCES: permission denied, mkdir '${parentDir}'`)
      }

      const caughtErr = ['EACCES', 'EPERM', 'EISDIR'].indexOf(err.code) > -1
      if (!caughtErr || (caughtErr && targetDir === curDir)) {
        throw err // Throw if it's just the last created dir.
      }
    }

    return curDir
  }, initDir)
}

/*
  Heads up, if you run this file separately instead of letting electron builder call it as a hook, you need to add ../ to both target and source directories. Or you should call while your current working directory is the same as the directory package.json is in.
*/
let targetDirectory = '../../../ReleaseArchive'
let sourceDirectory = '../../../BundledReleases'

let serverSourceDirectory = 'MAKE_BINARIES'

/*
  For below, we don't retain unpacked versions, because they aren't named appropriately for us to be able to put them in a folder in a structured fashion.

  We could probably set a structure where we have that data taken from the prior pass if it turns out empty in the current pass, but not only that is brittle, it's also reliant on the fact that our normal bundles start with A and come first (as in, 'A'ether). That's super brittle.
*/
let fileEndings = [
  ['.dmg', 'mac'],
  ['.dmg.blockmap', 'mac'],
  ['.zip', 'mac'], // Apparently needed for auto update. Who knew?
  ['.snap', 'linux'],
  ['.exe', 'win'],
  ['.exe.blockmap', 'win'],
  ['mac', 'delete'],
  ['linux-unpacked', 'delete'],
  ['win-ia32-unpacked', 'delete'],
  ['win-unpacked', 'delete'],
]

let serverFileEndings = [['', 'linux-server']]

function main() {
  if (process.platform === 'linux') {
    /*
            Assuming Linux = running inside docker build env. for cross-compile.
            If you're actually compiling natively in Linux, comment this if clause out.
        */
    targetDirectory = '/ReleaseArchive'
    sourceDirectory = '/BundledReleases'
  }
  fs.readdir(sourceDirectory, function(err, files) {
    for (let val of files) {
      let slc = val.split('+')
      moveFiles(val, slc, fileEndings, sourceDirectory, targetDirectory)
    }
  })
  fs.readdir(serverSourceDirectory, function(err, files) {
    for (let val of files) {
      let slc = val.split('+')
      moveFiles(
        val,
        slc,
        serverFileEndings,
        serverSourceDirectory,
        targetDirectory
      )
    }
  })
}

function moveFiles(unparsedName, parsedName, mapping, sourceDir, targetDir) {
  // if (parsedName.length < 2) {
  //     // No build data - we won't file.
  //     return
  // }
  if (parsedName.length < 2) {
    return
    /*
            If not a final result (i.e. just an untagged binary like aether-backend-mac-x64) without a plus that indicates the build number, we skip - not the finished product we're looking for.
        */
  }
  for (let ending of mapping) {
    if (parsedName[parsedName.length - 1].endsWith(ending[0])) {
      if (ending[1] === 'delete') {
        rimraf(sourceDir + '/' + unparsedName)
        return
      }
      // Set path variables
      let version = parsedName[0].replace('-Setup', '')
      // ^ This is for the Windows version, whose executable is a setup executable, but should be filed appropriately together with the others.
      if (ending[1] === 'linux-server') {
        version = parsedName[0].replace(
          'aether-backend-linux-x64-extverify',
          ''
        )
        version = version.replace('aether-backend-linux-x64', '')
        version = 'Aether' + version
      }
      let build = parsedName[1].substr(0, parsedName[1].lastIndexOf(ending[0]))
      let operatingSystem = ending[1]
      // Set paths
      let sourcePath = sourceDir
      let destPath =
        targetDir + '/' + version + '/' + build + '/' + operatingSystem
      // make paths
      mkDirByPathSync(destPath)
      fs.copyFileSync(
        sourcePath + '/' + unparsedName,
        destPath + '/' + unparsedName
      )
      fs.unlinkSync(sourcePath + '/' + unparsedName)
    }
  }
}

function rimraf(dir_path) {
  if (fs.existsSync(dir_path)) {
    try {
      fs.readdirSync(dir_path).forEach(function(entry) {
        var entry_path = path.join(dir_path, entry)
        if (fs.lstatSync(entry_path).isDirectory()) {
          rimraf(entry_path)
        } else {
          fs.unlinkSync(entry_path)
        }
      })
    } catch (exc) {
      fs.unlinkSync(dir_path)
      return
    }
    fs.rmdirSync(dir_path)
  }
}

if (process.argv[2] === 'run') {
  main()
}

exports.default = main
