'use strict'
// Store > Content Relations
Object.defineProperty(exports, '__esModule', { value: true })
var fe = require('../services/feapiconsumer/feapiconsumer')
var actions = {
  subToBoard: function (context, _a) {
    var fp = _a.fp,
      notify = _a.notify
    var now = Math.floor(Date.now() / 1000)
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      fp: fp,
      subbed: true,
      notify: notify,
      lastseen: now,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, notify, now, false, function (resp) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  unsubFromBoard: function (context, _a) {
    var fp = _a.fp
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      fp: fp,
      subbed: false,
      notify: false,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, false, false, 0, false, function (resp) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  silenceBoard: function (context, _a) {
    var fp = _a.fp
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      fp: fp,
      notify: false,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, false, 0, false, function (resp) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  unsilenceBoard: function (context, _a) {
    var fp = _a.fp
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      fp: fp,
      notify: true,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, true, 0, false, function (resp) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  setLastSeenForBoard: function (context, _a) {
    var fp = _a.fp
    console.log('last seen for board fingerprint: ', fp)
    var now = Math.floor(Date.now() / 1000)
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      fp: fp,
      lastseen: now,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, false, false, now, true, function (resp) {
      if (resp.committed) {
        console.log('Last Seen Commit Successful!')
      }
    })
  },
}
var mutations = {
  SAVE_BOARD_SIGNAL: function (state, payload) {
    // Find board in the store and change the state appropriately.
    for (var i = 0; i < state.allBoards.length; i++) {
      ;(function (i) {
        if (state.allBoards[i].fingerprint === payload.fp) {
          if (typeof payload.subbed !== 'undefined') {
            state.allBoards[i].subscribed = payload.subbed
          }
          if (typeof payload.notify !== 'undefined') {
            state.allBoards[i].notify = payload.notify
          }
          // if (typeof payload.lastseen !== 'undefined') {
          //   state.allBoards[i].LastSeen = payload.lastseen
          // }
        }
      })(i)
    }
    // Take a look at the currentBoard and if it's the same fp as the one we're subbing to, change its state as well.
    if (state.currentBoard.fingerprint === payload.fp) {
      if (typeof payload.subbed !== 'undefined') {
        state.currentBoard.subscribed = payload.subbed
      }
      if (typeof payload.notify !== 'undefined') {
        state.currentBoard.notify = payload.notify
      }
    }
  },
}
module.exports = {
  actions: actions,
  mutations: mutations,
}
//# sourceMappingURL=contentrelations.js.map
