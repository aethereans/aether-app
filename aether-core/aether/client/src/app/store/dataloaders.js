'use strict'
// Store > Data Loaders
Object.defineProperty(exports, '__esModule', { value: true })
var fe = require('../services/feapiconsumer/feapiconsumer')
var dataLoaders = {
  loadBoardScopeData: function (context, boardfp) {
    context.dispatch('setCurrentBoardFp', boardfp)
  },
  loadThreadScopeData: function (context, _a) {
    var boardfp = _a.boardfp,
      threadfp = _a.threadfp
    context.dispatch('setCurrentThreadFp', {
      boardfp: boardfp,
      threadfp: threadfp,
    })
  },
  loadGlobalScopeData: function (context) {
    context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', false)
    context.dispatch('updateBreadcrumbs')
    fe.GetAllBoards(function (resp) {
      console.log('received the all boards payload from fe')
      context.commit('SET_ALL_BOARDS', resp)
      context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', true)
    })
  },
  loadUserScopeData: function (context, _a) {
    var fp = _a.fp,
      userreq = _a.userreq,
      boardsreq = _a.boardsreq,
      threadsreq = _a.threadsreq,
      postsreq = _a.postsreq
    // context.commit('SET_CURRENT_USER_LOAD_COMPLETE', false)
    // ^ Disabling this prevents flashing between tab switches in user view. Let's make sure that this has no unintended side effects, if we see nothing, we can remove it.
    fe.GetUserAndGraph(
      fp,
      userreq,
      boardsreq,
      threadsreq,
      postsreq,
      function (resp) {
        console.log('Received user scope data')
        // We need to set in the query values we asked, so that the mutation will know what to override, and what not to.
        // resp.userentityrequested = userreq
        // resp.userboardsrequested = boardsreq
        // resp.userthreadsrequested = threadsreq
        // resp.userpostsrequested = postsreq
        context.commit('SET_USER_SCOPE_DATA', resp)
        context.commit('SET_CURRENT_USER_LOAD_COMPLETE', true)
        context.dispatch('updateBreadcrumbs')
      }
    )
  },
  /*----------  User reports loading  ----------*/
  loadBoardReports: function (context, boardfp) {
    fe.RequestBoardReports(boardfp, function (resp) {
      context.commit('SET_CURRENT_BOARD_REPORTS', resp.reportstabentriesList)
    })
  },
  setCurrentBoardReportsArrived: function (context, arrived) {
    context.commit('SET_CURRENT_BOARD_REPORTS_ARRIVED', arrived)
  },
  loadBoardModActions: function (context, boardfp) {
    console.log('load board mod actions called')
    fe.RequestBoardModActions(boardfp, function (resp) {
      context.commit(
        'SET_CURRENT_BOARD_MODACTIONS',
        resp.modactionstabentriesList
      )
    })
  },
  setCurrentBoardModActionsArrived: function (context, arrived) {
    context.commit('SET_CURRENT_BOARD_MODACTIONS_ARRIVED', arrived)
  },
}
exports.default = dataLoaders
//# sourceMappingURL=dataloaders.js.map
