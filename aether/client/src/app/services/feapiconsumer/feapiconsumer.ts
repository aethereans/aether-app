// Services > Frontend API Consumer

// This service is the API through which the client accesses data in the frontend.

export { } // This says this file is a module, not a script.

// Imports
const grpc = require('grpc')
// const resolve = require('path').resolve
var ipc = require('../../../../node_modules/electron-better-ipc')

// Consts
// const proto = grpc.load({
//   file: 'feapi/feapi.proto',
//   root: resolve(__dirname, '../protos')
// }).feapi

console.log('feapi consumer init')

var pmessages = require('../../../../../protos/feapi/feapi_pb.js')
// var feobjmessages = require('../../../../../protos/feobjects/feobjects_pb.js');
var proto = require('../../../../../protos/feapi/feapi_grpc_pb')

let feAPIConsumer: any
let Initialised: boolean
let initInProgress: boolean = false
let clientApiServerPortIsSet: boolean = false

function timeout(ms: any) {
  return new Promise(function(resolve) {
    return setTimeout(resolve, ms)
  })
}
// ^ Promisified wait, so that it won't actually block like while .. {} does. Useful with async/await.

async function checkPortSet() {
  if (clientApiServerPortIsSet === false) {
    await timeout(25)
    await checkPortSet()
  }
}

let ExportedMethods = {
  async Initialise() {
    // console.log('init is entered')
    if (initInProgress) {
      // console.log('init is already in progress, waiting until the other one completes')
      await checkPortSet()
      // console.log('init is complete, returning to normal process')
      return
    }
    initInProgress = true
    console.log('init is called')
    let feapiport = await ipc.callMain('GetFrontendAPIPort')
    feAPIConsumer = new proto.FrontendAPIClient(
      '127.0.0.1:' + feapiport,
      grpc.credentials.createInsecure()
    )
    console.log(feAPIConsumer)
    let clapiserverport = await ipc.callMain('GetClientAPIServerPort')
    await ExportedMethods.SetClientAPIServerPort(clapiserverport)
    ipc.callMain('SetFrontendClientConnInitialised', true)
    Initialised = true
    initInProgress = false
  },
  GetAllBoards(callback: any) {
    WaitUntilFrontendReady(async function() {
      console.log('get all boards is making a call')
      console.log('initstate: ', Initialised)
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      let req = new pmessages.AllBoardsRequest()
      feAPIConsumer.getAllBoards(req, function(err: any, response: any) {
        if (err) {
          console.log(err)
        } else {
          callback(response.toObject().allboardsList)
        }
      })
    })
  },
  SetClientAPIServerPort(clientAPIServerPort: number) {
    WaitUntilFrontendReady(async function() {
      console.log(
        'clapiserverport mapping is triggered. initstate: ',
        Initialised
      )
      // if (!Initialised) {
      //   await ExportedMethods.Initialise()
      // }
      let req = new pmessages.SetClientAPIServerPortRequest()
      req.setPort(clientAPIServerPort)
      // console.log(req)
      feAPIConsumer.setClientAPIServerPort(req, function(
        err: any,
        response: any
      ) {
        if (err) {
          console.log(err)
        } else {
          console.log(response)
          clientApiServerPortIsSet = true
        }
      })
    })
  },
  GetBoardAndThreads(boardfp: string, sortByNew: boolean, callback: any) {
    console.log('get boards and threads received:')
    console.log(boardfp, sortByNew)
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetBoardsAndThread triggered.')
      let req = new pmessages.BoardAndThreadsRequest()
      req.setBoardfingerprint(boardfp)
      if (sortByNew) {
        req.setSortthreadsbynew(true)
      }
      feAPIConsumer.getBoardAndThreads(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          console.log(resp.toObject())
          callback(resp.toObject())
        }
      })
    })
  },
  GetThreadAndPosts(boardfp: string, threadfp: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetThreadAndPosts triggered.')
      let req = new pmessages.ThreadAndPostsRequest()
      req.setBoardfingerprint(boardfp)
      req.setThreadfingerprint(threadfp)
      feAPIConsumer.getThreadAndPosts(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  SetBoardSignal(
    fp: string,
    subbed: boolean,
    notify: boolean,
    lastseen: number,
    lastSeenOnly: boolean,
    callback: any
  ) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SetBoardSignal triggered.')
      let req = new pmessages.BoardSignalRequest()
      req.setFingerprint(fp)
      req.setSubscribed(subbed)
      req.setNotify(notify)
      req.setLastseen(lastseen)
      req.setLastseenonly(lastSeenOnly)
      feAPIConsumer.setBoardSignal(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  GetUserAndGraph(
    fp: string,
    userEntityRequested: boolean,
    boardsRequested: boolean,
    threadsRequested: boolean,
    postsRequested: boolean,
    callback: any
  ) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetUserAndGraph triggered.')
      let req = new pmessages.UserAndGraphRequest()
      req.setFingerprint(fp)
      req.setUserentityrequested(userEntityRequested)
      req.setUserboardsrequested(boardsRequested)
      req.setUserthreadsrequested(threadsRequested)
      req.setUserpostsrequested(postsRequested)
      console.log(req)
      feAPIConsumer.getUserAndGraph(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  GetUncompiledEntityByKey(
    entityType: string,
    ownerfp: string,
    limit: number,
    offset: number,
    callback: any
  ) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('GetUncompiledEntityByKey triggered.')
      let req = new pmessages.UncompiledEntityByKeyRequest()
      if (entityType === 'Board') {
        req.setEntitytype(pmessages.UncompiledEntityType.BOARD)
      }
      if (entityType === 'Thread') {
        req.setEntitytype(pmessages.UncompiledEntityType.THREAD)
      }
      if (entityType === 'Post') {
        req.setEntitytype(pmessages.UncompiledEntityType.POST)
      }
      if (entityType === 'Vote') {
        req.setEntitytype(pmessages.UncompiledEntityType.VOTE)
      }
      if (entityType === 'Key') {
        req.setEntitytype(pmessages.UncompiledEntityType.KEY)
      }
      if (entityType === 'Truststate') {
        req.setEntitytype(pmessages.UncompiledEntityType.TRUSTSTATE)
      }
      req.setLimit(limit)
      req.setOffset(offset)
      req.setOwnerfingerprint(ownerfp)
      console.log(req)
      feAPIConsumer.getUncompiledEntityByKey(req, function(
        err: any,
        resp: any
      ) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  SendInflightsPruneRequest(callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendInflightsPruneRequest triggered.')
      let req = new pmessages.InflightsPruneRequest()
      feAPIConsumer.sendInflightsPruneRequest(req, function(
        err: any,
        resp: any
      ) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  RequestAmbientStatus(callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('RequestAmbientStatus triggered.')
      let req = new pmessages.AmbientStatusRequest()
      feAPIConsumer.requestAmbientStatus(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },
  SetNotificationsSignal(seen: boolean, fp: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SetNotificationsSignal triggered.')
      let req = new pmessages.NotificationsSignalPayload()
      req.setSeen(seen)
      req.setReaditemfingerprint(fp)
      feAPIConsumer.setNotificationsSignal(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp.toObject())
        }
      })
    })
  },

  SetOnboardComplete(callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SetOnboardComplete triggered.')
      let req = new pmessages.OnboardCompleteRequest()
      req.setOnboardcomplete(true)
      feAPIConsumer.setOnboardComplete(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          const metrics = require('../metrics/metrics')()
          metrics.SendRaw('Onboarding complete')
          callback(resp.toObject())
        }
      })
    })
  },

  SendAddress(addr: any, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendAddress triggered.')
      let req = new pmessages.SendAddressPayload()
      req.setAddress(addr)
      try {
        feAPIConsumer.sendAddress(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },

  RequestBoardReports(boardfp: any, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('RequestBoardReports triggered.')
      let req = new pmessages.BoardReportsRequest()
      req.setBoardfingerprint(boardfp)
      try {
        feAPIConsumer.requestBoardReports(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },

  SendMintedUsernames(payload: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendMintedUsernames triggered.')
      let req = new pmessages.SendMintedUsernamesPayload()
      req.setMintedusernamesrawjson(payload)
      try {
        feAPIConsumer.sendMintedUsernames(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },

  SendClientVersion(payload: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendClientVersion triggered.')
      let req = new pmessages.ClientVersionPayload()
      req.setCurrentclientversion(payload)
      try {
        feAPIConsumer.sendClientVersion(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },

  /*----------  FE config changes  ----------*/

  SendModModeEnabledStatus(modModeEnabled: boolean, callback: any) {
    let e = new pmessages.FEConfigChangesPayload()
    e.setModmodeenabled(modModeEnabled)
    e.setModmodeenabledisset(true)
    this.SendFEConfigChanges(e, callback)
  },

  SendExternalContentAutoloadDisabledStatus(
    externalContentAutoloadDisabled: boolean,
    callback: any
  ) {
    let e = new pmessages.FEConfigChangesPayload()
    e.setExternalcontentautoloaddisabled(externalContentAutoloadDisabled)
    e.setExternalcontentautoloaddisabledisset(true)
    this.SendFEConfigChanges(e, callback)
  },

  /*----------  FE Config changes delivery base  ----------*/

  SendFEConfigChanges(feconfig: any, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('SendFEConfigChanges triggered.')
      let req = new pmessages.FEConfigChangesPayload()
      req = feconfig
      try {
        feAPIConsumer.sendFEConfigChanges(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },

  /*----------  Notifications signals  ----------*/

  markSeen() {
    this.SetNotificationsSignal(true, '', function() { })
  },
  markRead(fp: string) {
    this.SetNotificationsSignal(true, fp, function() { })
  },

  IsInitialised(): boolean {
    return Initialised
  },

  /*----------  Methods for user signal actions  ----------*/

  /*
    Important thing here. We do have a RETRACT defined, but this is not defined for anything that goes into a bloom filter. Which means, you cannot retract an upvote, but you can downvote to reverse it. The reason why is that upvotes and downvotes (and elects) are aggregated, therefore after they get added to the bloom filter, we only know probabilistically that they're there. We have two bloom filters for each so we can have a +1 and -1, but adding 0 means adding another bloom filter in. Depending on the demand for a retract we can add a third bloom to the implementation to keep tracking of that, but bloom filters are very expensive because they're per-entity, and we have a lot of entities.

    This does not apply to non-aggregated signals like reporting to mod, those are kept instact and individual, and they can be retracted.
  */
  Upvote(
    this: any,
    targetfp: string,
    priorfp: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'ADDS_TO_DISCUSSION',
      'UPVOTE',
      'CONTENT',
      '',
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Upvote', priorfp.length ? 'Edit' : 'Create')
  },

  Downvote(
    this: any,
    targetfp: string,
    priorfp: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'ADDS_TO_DISCUSSION',
      'DOWNVOTE',
      'CONTENT',
      '',
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Downvote', priorfp.length ? 'Edit' : 'Create')
  },

  ReportToMod(
    this: any,
    targetfp: string,
    priorfp: string,
    reason: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'FOLLOWS_GUIDELINES',
      'REPORT_TO_MOD',
      'CONTENT',
      reason,
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ReportToMod', priorfp.length ? 'Edit' : 'Create')
  },
  // ModDelete instead of ModBlock, to keep it more human-meaningful.
  ModDelete(
    this: any,
    targetfp: string,
    priorfp: string,
    reason: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'MOD_ACTIONS',
      'MODBLOCK',
      'CONTENT',
      reason,
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModDelete', priorfp.length ? 'Edit' : 'Create')
  },

  ModApprove(
    this: any,
    targetfp: string,
    priorfp: string,
    reason: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'MOD_ACTIONS',
      'MODAPPROVE',
      'CONTENT',
      reason,
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModApprove', priorfp.length ? 'Edit' : 'Create')
  },

  ModIgnore(
    this: any,
    targetfp: string,
    priorfp: string,
    reason: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'MOD_ACTIONS',
      'MODIGNORE',
      'CONTENT',
      reason,
      boardfp,
      threadfp,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModIgnore', priorfp.length ? 'Edit' : 'Create')
  },

  Follow(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'PUBLIC_TRUST',
      'FOLLOW',
      'USER',
      '',
      '',
      '',
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Follow', priorfp.length ? 'Edit' : 'Create')
  },

  Block(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'PUBLIC_TRUST',
      'BLOCK',
      'USER',
      '',
      '',
      '',
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Block', priorfp.length ? 'Edit' : 'Create')
  },

  Elect(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'PUBLIC_ELECT',
      'ELECT',
      'USER',
      '',
      '',
      '',
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Elect', priorfp.length ? 'Edit' : 'Create')
  },
  Disqualify(this: any, targetfp: string, priorfp: string, callback: any) {
    this.sendSignalEvent(
      targetfp,
      priorfp,
      'PUBLIC_ELECT',
      'DISQUALIFY',
      'USER',
      '',
      '',
      '',
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Disqualify', priorfp.length ? 'Edit' : 'Create')
  },

  /*----------  Base signal event action.  ----------*/

  sendSignalEvent(
    targetfp: string,
    priorfp: string,
    typeclass: string,
    typ: string,
    targettype: string,
    signaltext: string,
    boardfp: string,
    threadfp: string,
    callback: any
  ) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('Send Signal Event base triggered.')
      let now = Math.floor(Date.now() / 1000)
      let req = new pmessages.SignalEventPayload()
      let e = new pmessages.Event()
      var localUser = require('../../store/index').default.state.localUser
      // ^ Only import when needed and only the specific part. Because vuexstore is also importing this feapi - we don't want it being imported at the beginning to prevent vuexstore from loading feapi.
      e.setOwnerfingerprint(localUser.fingerprint)
      e.setPriorfingerprint(priorfp)
      e.setEventtype(
        priorfp.length === 0
          ? pmessages.EventType.CREATE
          : pmessages.EventType.UPDATE
      )
      e.setTimestamp(now)
      req.setEvent(e)
      req.setSignaltargettype(pmessages.SignalTargetType[targettype])
      if (targettype === 'CONTENT') {
        req.setTargetboard(boardfp)
        req.setTargetthread(threadfp)
      }
      if (targettype === 'USER') {
        req.setDomain() // todo
      }
      req.setTargetfingerprint(targetfp)
      req.setSignaltypeclass(pmessages.SignalTypeClass[typeclass])
      req.setSignaltext(signaltext)
      console.log('signal type:')
      console.log(pmessages.SignalType[typ])
      req.setSignaltype(pmessages.SignalType[typ])
      feAPIConsumer.sendSignalEvent(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp)
        }
      })
    })
  },

  /*----------  Methods for content event actions  ----------*/

  /*
    These are things like creating or editing entities that the user has created. If a priorfp is provided, it is an update. If not, it is a create.
  */

  SendBoardContent(this: any, priorfp: string, boarddata: any, callback: any) {
    this.sendContentEvent(
      priorfp,
      boarddata,
      undefined,
      undefined,
      undefined,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Board', priorfp.length ? 'Edit' : 'Create')
  },

  SendThreadContent(
    this: any,
    priorfp: string,
    threaddata: any,
    callback: any
  ) {
    this.sendContentEvent(
      priorfp,
      undefined,
      threaddata,
      undefined,
      undefined,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Thread', priorfp.length ? 'Edit' : 'Create')
  },

  SendPostContent(this: any, priorfp: string, postdata: any, callback: any) {
    this.sendContentEvent(
      priorfp,
      undefined,
      undefined,
      postdata,
      undefined,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Post', priorfp.length ? 'Edit' : 'Create')
  },

  SendUserContent(this: any, priorfp: string, userdata: any, callback: any) {
    this.sendContentEvent(
      priorfp,
      undefined,
      undefined,
      undefined,
      userdata,
      callback
    )
    const metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('User', priorfp.length ? 'Edit' : 'Create')
  },

  /*----------  Base content event action.  ----------*/

  sendContentEvent(
    priorfp: string,
    boarddata: any,
    threaddata: any,
    postdata: any,
    userdata: any,
    callback: any
  ) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('Send Content Event base triggered.')
      let now = Math.floor(Date.now() / 1000)
      let req = new pmessages.ContentEventPayload()
      let e = new pmessages.Event()
      var localUser = require('../../store/index').default.state.localUser
      // ^ Only import when needed and only the specific part. Because vuexstore is also importing this feapi - we don't want it being imported at the beginning to prevent vuexstore from loading feapi.
      let globalMethods = require('../globals/methods')
      if (!globalMethods.IsUndefined(localUser)) {
        e.setOwnerfingerprint(localUser.fingerprint)
      }
      e.setPriorfingerprint(priorfp)
      e.setEventtype(
        priorfp.length === 0
          ? pmessages.EventType.CREATE
          : pmessages.EventType.UPDATE
      )
      e.setTimestamp(now)
      req.setEvent(e)
      req.setBoarddata(boarddata)
      req.setThreaddata(threaddata)
      req.setPostdata(postdata)
      req.setKeydata(userdata)
      feAPIConsumer.sendContentEvent(req, function(err: any, resp: any) {
        if (err) {
          console.log(err)
        } else {
          callback(resp)
        }
      })
    })
  },

  /*----------  Search types  ----------*/

  SendCommunitySearchRequest(query: string, callback: any) {
    this.sendSearchRequest('Board', query, callback)
  },

  SendContentSearchRequest(query: string, callback: any) {
    this.sendSearchRequest('Content', query, callback)
  },

  SendUserSearchRequest(query: string, callback: any) {
    this.sendSearchRequest('User', query, callback)
  },

  /*----------  Base search request action  ----------*/

  sendSearchRequest(searchType: string, query: string, callback: any) {
    WaitUntilFrontendReady(async function() {
      if (!Initialised) {
        await ExportedMethods.Initialise()
      }
      console.log('sendSearchRequest triggered.')
      let req = new pmessages.SearchRequestPayload()
      req.setSearchtype(searchType)
      req.setSearchquery(query)
      try {
        feAPIConsumer.sendSearchRequest(req, function(err: any, resp: any) {
          if (err) {
            console.log(err)
          } else {
            callback(resp.toObject())
          }
        })
      } catch (err) {
        // This catches non-grpc errors like assert.
        callback(err)
      }
    })
  },
}
module.exports = ExportedMethods

function WaitUntilFrontendReady(cb: any): any {
  async function check() {
    let initialised = await ipc.callMain('GetFrontendClientConnInitialised')
    // console.log(initialised)
    if (!initialised) {
      // console.log("Frontend still not ready, waiting a little more...")
      return setTimeout(check, 333)
    } else {
      return cb()
    }
  }
  return check()
}
