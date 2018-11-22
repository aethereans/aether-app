// Store > Content Relations

// This file handles the content relations (similar to the similarly named file on the frontend > services > configstore > contentrelations).

/*
  Generally speaking, our way of working is this. When there is a change that requires a state change, we apply that state change to the store on the client side, *and* send a signal to the frontend of that change. the change is applied immediately, but if the frontend returns an error or a timeout, the change is reverted.
*/

export { }
var fe = require('../services/feapiconsumer/feapiconsumer')

let actions = {
  subToBoard(context: any, { fp, notify }: { fp: string, notify: boolean }) {
    let now = Math.floor(Date.now() / 1000)
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      'fp': fp,
      'subbed': true,
      'notify': notify,
      'lastseen': now,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, notify, now, false, function(resp: any) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  unsubFromBoard(context: any, { fp }: { fp: string }) {
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      'fp': fp,
      'subbed': false,
      'notify': false
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, false, false, 0, false, function(resp: any) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  silenceBoard(context: any, { fp }: { fp: string }) {
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      'fp': fp,
      'notify': false
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, false, 0, false, function(resp: any) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  unsilenceBoard(context: any, { fp }: { fp: string }) {
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      'fp': fp,
      'notify': true
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, true, true, 0, false, function(resp: any) {
      if (resp.committed) {
        // console.log("Commit Successful!")
      }
    })
  },
  setLastSeenForBoard(context: any, { fp }: { fp: string }) {
    console.log('last seen for board fingerprint: ', fp)
    let now = Math.floor(Date.now() / 1000)
    // Set the local client state
    context.commit('SAVE_BOARD_SIGNAL', {
      'fp': fp,
      'lastseen': now,
    })
    // then send the appropriate message to the frontend
    fe.SetBoardSignal(fp, false, false, now, true, function(resp: any) {
      if (resp.committed) {
        console.log("Last Seen Commit Successful!")
      }
    })
  }
}

let mutations = {
  SAVE_BOARD_SIGNAL(state: any, payload: any) {
    // Find board in the store and change the state appropriately.
    for (var i: number = 0; i < state.allBoards.length; i++) {
      (function(i: number) {
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
      }
      )(i)
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
  actions,
  mutations
}