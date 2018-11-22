// Service > Hashimage

// This service converts a hash into a visually recognisable artifact.

export { }

let blockies = require('./blockies-modified') // the modified version uses a custom colour mapper that uses our colours in globals.scss.
let hqx = require('./hqx-modified').hqx // I made this use Math instead of Window.Math to make it work.

function generate(hash: string, isUser: boolean): string {
  if (typeof hash === 'undefined') {
    return ""
  }
  let blockiesConf: any = {
    seed: hash,
    size: 10,
    scale: 2,
    spotcolor: -1,
  }
  if (isUser) {
    blockiesConf.spotcolor = '#343D46' //a-grey-200
  }
  let canvas = blockies.create(blockiesConf)
  let scaledCanvas = hqx(canvas, 4)
  return scaledCanvas.toDataURL()
}

module.exports = generate