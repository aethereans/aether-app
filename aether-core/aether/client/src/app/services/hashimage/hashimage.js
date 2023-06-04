'use strict'
// Service > Hashimage
Object.defineProperty(exports, '__esModule', { value: true })
var blockies = require('./blockies-modified') // the modified version uses a custom colour mapper that uses our colours in globals.scss.
var hqx = require('./hqx-modified').hqx // I made this use Math instead of Window.Math to make it work.
function generate(hash, isUser) {
  if (typeof hash === 'undefined') {
    return ''
  }
  var blockiesConf = {
    seed: hash,
    size: 10,
    scale: 2,
    spotcolor: -1,
  }
  if (isUser) {
    blockiesConf.spotcolor = '#343D46' //a-grey-200
  }
  var canvas = blockies.create(blockiesConf)
  var scaledCanvas = hqx(canvas, 4)
  return scaledCanvas.toDataURL()
}
module.exports = generate
//# sourceMappingURL=hashimage.js.map
