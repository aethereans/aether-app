// Store > Data Loaders

// These actions are the high-level loaders that correspond roughly to page contexts we have.

export {}
var fe = require('../services/feapiconsumer/feapiconsumer')

let dataLoaders = {
  loadBoardScopeData(context: any, boardfp: string) {
    context.dispatch('setCurrentBoardFp', boardfp)
  },

  loadThreadScopeData(
    context: any,
    { boardfp, threadfp }: { boardfp: string; threadfp: string }
  ) {
    context.dispatch('setCurrentThreadFp', {
      boardfp: boardfp,
      threadfp: threadfp,
    })
  },

  loadGlobalScopeData(context: any) {
    context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', false)
    context.dispatch('updateBreadcrumbs')
    fe.GetAllBoards(function (resp: any) {
      console.log('received the all boards payload from fe')
      context.commit('SET_ALL_BOARDS', resp)
      context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', true)
    })
  },

  loadUserScopeData(
    context: any,
    {
      fp,
      userreq,
      boardsreq,
      threadsreq,
      postsreq,
    }: {
      fp: string
      userreq: boolean
      boardsreq: boolean
      threadsreq: boolean
      postsreq: boolean
    }
  ) {
    // context.commit('SET_CURRENT_USER_LOAD_COMPLETE', false)
    // ^ Disabling this prevents flashing between tab switches in user view. Let's make sure that this has no unintended side effects, if we see nothing, we can remove it.
    fe.GetUserAndGraph(
      fp,
      userreq,
      boardsreq,
      threadsreq,
      postsreq,
      function (resp: any) {
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
  loadBoardReports(context: any, boardfp: string) {
    fe.RequestBoardReports(boardfp, function (resp: any) {
      context.commit('SET_CURRENT_BOARD_REPORTS', resp.reportstabentriesList)
    })
  },
  setCurrentBoardReportsArrived(context: any, arrived: boolean) {
    context.commit('SET_CURRENT_BOARD_REPORTS_ARRIVED', arrived)
  },
  loadBoardModActions(context: any, boardfp: string) {
    console.log('load board mod actions called')
    fe.RequestBoardModActions(boardfp, function (resp: any) {
      context.commit(
        'SET_CURRENT_BOARD_MODACTIONS',
        resp.modactionstabentriesList
      )
    })
  },
  setCurrentBoardModActionsArrived(context: any, arrived: boolean) {
    context.commit('SET_CURRENT_BOARD_MODACTIONS_ARRIVED', arrived)
  },
}

export default dataLoaders
