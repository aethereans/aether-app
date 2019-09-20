export {}

// let ipc = require('../../../../node_modules/electron-better-ipc')

// /*----------  Renderer receivers  ----------*/
// // i.e. renderer does something at the request of renderer

// ipc.answerMain('RouteTo', function(route: string) {
//   console.log('route is: ', route)
//   var router = require('../../renderermain')
//   router.push(route)
// })

/*
  The only way to get access to the router here appears to be importing the main. That's a big no - I'd rather have this mapping on the main instead.
*/
