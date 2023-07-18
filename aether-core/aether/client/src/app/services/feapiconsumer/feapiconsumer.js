'use strict'
// Services > Frontend API Consumer
var __awaiter =
  (this && this.__awaiter) ||
  function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
      function fulfilled(value) {
        try {
          step(generator.next(value))
        } catch (e) {
          reject(e)
        }
      }
      function rejected(value) {
        try {
          step(generator['throw'](value))
        } catch (e) {
          reject(e)
        }
      }
      function step(result) {
        result.done
          ? resolve(result.value)
          : new P(function (resolve) {
              resolve(result.value)
            }).then(fulfilled, rejected)
      }
      step((generator = generator.apply(thisArg, _arguments || [])).next())
    })
  }
var __generator =
  (this && this.__generator) ||
  function (thisArg, body) {
    var _ = {
        label: 0,
        sent: function () {
          if (t[0] & 1) throw t[1]
          return t[1]
        },
        trys: [],
        ops: [],
      },
      f,
      y,
      t,
      g
    return (
      (g = { next: verb(0), throw: verb(1), return: verb(2) }),
      typeof Symbol === 'function' &&
        (g[Symbol.iterator] = function () {
          return this
        }),
      g
    )
    function verb(n) {
      return function (v) {
        return step([n, v])
      }
    }
    function step(op) {
      if (f) throw new TypeError('Generator is already executing.')
      while (_)
        try {
          if (
            ((f = 1),
            y &&
              (t =
                op[0] & 2
                  ? y['return']
                  : op[0]
                  ? y['throw'] || ((t = y['return']) && t.call(y), 0)
                  : y.next) &&
              !(t = t.call(y, op[1])).done)
          )
            return t
          if (((y = 0), t)) op = [op[0] & 2, t.value]
          switch (op[0]) {
            case 0:
            case 1:
              t = op
              break
            case 4:
              _.label++
              return { value: op[1], done: false }
            case 5:
              _.label++
              y = op[1]
              op = [0]
              continue
            case 7:
              op = _.ops.pop()
              _.trys.pop()
              continue
            default:
              if (
                !((t = _.trys), (t = t.length > 0 && t[t.length - 1])) &&
                (op[0] === 6 || op[0] === 2)
              ) {
                _ = 0
                continue
              }
              if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) {
                _.label = op[1]
                break
              }
              if (op[0] === 6 && _.label < t[1]) {
                _.label = t[1]
                t = op
                break
              }
              if (t && _.label < t[2]) {
                _.label = t[2]
                _.ops.push(op)
                break
              }
              if (t[2]) _.ops.pop()
              _.trys.pop()
              continue
          }
          op = body.call(thisArg, _)
        } catch (e) {
          op = [6, e]
          y = 0
        } finally {
          f = t = 0
        }
      if (op[0] & 5) throw op[1]
      return { value: op[0] ? op[1] : void 0, done: true }
    }
  }
Object.defineProperty(exports, '__esModule', { value: true })
// Imports
var grpc = require('grpc')
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
var feAPIConsumer
var Initialised
var initInProgress = false
var clientApiServerPortIsSet = false
function timeout(ms) {
  return new Promise(function (resolve) {
    return setTimeout(resolve, ms)
  })
}
// ^ Promisified wait, so that it won't actually block like while .. {} does. Useful with async/await.
function checkPortSet() {
  return __awaiter(this, void 0, void 0, function () {
    return __generator(this, function (_a) {
      switch (_a.label) {
        case 0:
          if (!(clientApiServerPortIsSet === false)) return [3 /*break*/, 3]
          return [4 /*yield*/, timeout(25)]
        case 1:
          _a.sent()
          return [4 /*yield*/, checkPortSet()]
        case 2:
          _a.sent()
          _a.label = 3
        case 3:
          return [2 /*return*/]
      }
    })
  })
}
var ExportedMethods = {
  Initialise: function () {
    return __awaiter(this, void 0, void 0, function () {
      var feapiport, clapiserverport
      return __generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            if (!initInProgress) return [3 /*break*/, 2]
            // console.log('init is already in progress, waiting until the other one completes')
            return [
              4 /*yield*/,
              checkPortSet(),
              // console.log('init is complete, returning to normal process')
            ]
          case 1:
            // console.log('init is already in progress, waiting until the other one completes')
            _a.sent()
            // console.log('init is complete, returning to normal process')
            return [2 /*return*/]
          case 2:
            initInProgress = true
            console.log('init is called')
            return [4 /*yield*/, ipc.callMain('GetFrontendAPIPort')]
          case 3:
            feapiport = _a.sent()
            feAPIConsumer = new proto.FrontendAPIClient(
              '127.0.0.1:' + feapiport,
              grpc.credentials.createInsecure(),
              {
                'grpc.max_recv_message_length': 2147483647,
                'grpc.max_send_message_length': 2147483647,
              }
            )
            console.log(feAPIConsumer)
            return [4 /*yield*/, ipc.callMain('GetClientAPIServerPort')]
          case 4:
            clapiserverport = _a.sent()
            return [
              4 /*yield*/,
              ExportedMethods.SetClientAPIServerPort(clapiserverport),
            ]
          case 5:
            _a.sent()
            ipc.callMain('SetFrontendClientConnInitialised', true)
            Initialised = true
            initInProgress = false
            return [2 /*return*/]
        }
      })
    })
  },
  GetAllBoards: function (callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              console.log('get all boards is making a call')
              console.log('initstate: ', Initialised)
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              req = new pmessages.AllBoardsRequest()
              feAPIConsumer.getAllBoards(req, function (err, response) {
                if (err) {
                  console.log(err)
                } else {
                  callback(response.toObject().allboardsList)
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SetClientAPIServerPort: function (clientAPIServerPort) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          console.log(
            'clapiserverport mapping is triggered. initstate: ',
            Initialised
          )
          req = new pmessages.SetClientAPIServerPortRequest()
          req.setPort(clientAPIServerPort)
          // console.log(req)
          feAPIConsumer.setClientAPIServerPort(req, function (err, response) {
            if (err) {
              console.log(err)
            } else {
              console.log(response)
              clientApiServerPortIsSet = true
            }
          })
          return [2 /*return*/]
        })
      })
    })
  },
  GetBoardAndThreads: function (boardfp, sortByNew, callback) {
    console.log('get boards and threads received:')
    console.log(boardfp, sortByNew)
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('GetBoardsAndThread triggered.')
              req = new pmessages.BoardAndThreadsRequest()
              req.setBoardfingerprint(boardfp)
              if (sortByNew) {
                req.setSortthreadsbynew(true)
              }
              feAPIConsumer.getBoardAndThreads(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  console.log(resp.toObject())
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  GetThreadAndPosts: function (boardfp, threadfp, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('GetThreadAndPosts triggered.')
              req = new pmessages.ThreadAndPostsRequest()
              req.setBoardfingerprint(boardfp)
              req.setThreadfingerprint(threadfp)
              feAPIConsumer.getThreadAndPosts(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SetBoardSignal: function (
    fp,
    subbed,
    notify,
    lastseen,
    lastSeenOnly,
    callback
  ) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SetBoardSignal triggered.')
              req = new pmessages.BoardSignalRequest()
              req.setFingerprint(fp)
              req.setSubscribed(subbed)
              req.setNotify(notify)
              req.setLastseen(lastseen)
              req.setLastseenonly(lastSeenOnly)
              feAPIConsumer.setBoardSignal(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  GetUserAndGraph: function (
    fp,
    userEntityRequested,
    boardsRequested,
    threadsRequested,
    postsRequested,
    callback
  ) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('GetUserAndGraph triggered.')
              req = new pmessages.UserAndGraphRequest()
              req.setFingerprint(fp)
              req.setUserentityrequested(userEntityRequested)
              req.setUserboardsrequested(boardsRequested)
              req.setUserthreadsrequested(threadsRequested)
              req.setUserpostsrequested(postsRequested)
              console.log(req)
              feAPIConsumer.getUserAndGraph(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  GetUncompiledEntityByKey: function (
    entityType,
    ownerfp,
    boardName,
    keyName,
    limit,
    offset,
    callback
  ) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('GetUncompiledEntityByKey triggered.')
              req = new pmessages.UncompiledEntityByKeyRequest()
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
              if (typeof ownerfp !== 'undefined' && ownerfp.length > 0) {
                req.setOwnerfingerprint(ownerfp)
              }
              if (typeof boardName !== 'undefined' && boardName.length > 0) {
                req.setBoardname(boardName)
              }
              if (typeof keyName !== 'undefined' && keyName.length > 0) {
                req.setKeyname(keyName)
              }
              console.log(req)
              feAPIConsumer.getUncompiledEntityByKey(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SendInflightsPruneRequest: function (callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SendInflightsPruneRequest triggered.')
              req = new pmessages.InflightsPruneRequest()
              feAPIConsumer.sendInflightsPruneRequest(
                req,
                function (err, resp) {
                  if (err) {
                    console.log(err)
                  } else {
                    callback(resp.toObject())
                  }
                }
              )
              return [2 /*return*/]
          }
        })
      })
    })
  },
  RequestAmbientStatus: function (callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('RequestAmbientStatus triggered.')
              req = new pmessages.AmbientStatusRequest()
              feAPIConsumer.requestAmbientStatus(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SetNotificationsSignal: function (seen, fp, markAllAsRead, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SetNotificationsSignal triggered.')
              req = new pmessages.NotificationsSignalPayload()
              req.setSeen(seen)
              req.setReaditemfingerprint(fp)
              req.setMarkallasread(markAllAsRead)
              feAPIConsumer.setNotificationsSignal(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SetOnboardComplete: function (onboardComplete, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SetOnboardComplete triggered.')
              req = new pmessages.OnboardCompleteRequest()
              req.setOnboardcomplete(onboardComplete)
              feAPIConsumer.setOnboardComplete(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  var metrics_1 = require('../metrics/metrics')()
                  metrics_1.SendRaw('Onboarding complete')
                  callback(resp.toObject())
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SendAddress: function (addr, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SendAddress triggered.')
              req = new pmessages.SendAddressPayload()
              req.setAddress(addr)
              try {
                feAPIConsumer.sendAddress(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  RequestBoardReports: function (boardfp, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('RequestBoardReports triggered.')
              req = new pmessages.BoardReportsRequest()
              req.setBoardfingerprint(boardfp)
              try {
                feAPIConsumer.requestBoardReports(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  RequestBoardModActions: function (boardfp, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('RequestBoardModActions triggered.')
              req = new pmessages.BoardModActionsRequest()
              req.setBoardfingerprint(boardfp)
              try {
                feAPIConsumer.requestBoardModActions(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SendMintedUsernames: function (payload, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SendMintedUsernames triggered.')
              req = new pmessages.SendMintedUsernamesPayload()
              req.setMintedusernamesrawjson(payload)
              try {
                feAPIConsumer.sendMintedUsernames(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  SendClientVersion: function (payload, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SendClientVersion triggered.')
              req = new pmessages.ClientVersionPayload()
              req.setCurrentclientversion(payload)
              try {
                feAPIConsumer.sendClientVersion(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  /*----------  FE config changes  ----------*/
  SendModModeEnabledStatus: function (modModeEnabled, callback) {
    var e = new pmessages.FEConfigChangesPayload()
    e.setModmodeenabled(modModeEnabled)
    e.setModmodeenabledisset(true)
    this.SendFEConfigChanges(e, callback)
  },
  SendAlwaysShowNSFWListStatus: function (alwaysShowNSFWList, callback) {
    var e = new pmessages.FEConfigChangesPayload()
    e.setAlwaysshownsfwlist(alwaysShowNSFWList)
    e.setAlwaysshownsfwlistisset(true)
    this.SendFEConfigChanges(e, callback)
  },
  SendExternalContentAutoloadDisabledStatus: function (
    externalContentAutoloadDisabled,
    callback
  ) {
    var e = new pmessages.FEConfigChangesPayload()
    e.setExternalcontentautoloaddisabled(externalContentAutoloadDisabled)
    e.setExternalcontentautoloaddisabledisset(true)
    this.SendFEConfigChanges(e, callback)
  },
  SendSFWListDisabledStatus: function (sfwListDisabled, callback) {
    var e = new pmessages.FEConfigChangesPayload()
    e.setSfwlistdisabled(sfwListDisabled)
    e.setSfwlistdisabledisset(true)
    this.SendFEConfigChanges(e, callback)
  },
  /*----------  FE Config changes delivery base  ----------*/
  SendFEConfigChanges: function (feconfig, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('SendFEConfigChanges triggered.')
              req = new pmessages.FEConfigChangesPayload()
              req = feconfig
              try {
                feAPIConsumer.sendFEConfigChanges(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
  /*----------  Notifications signals  ----------*/
  markSeen: function () {
    this.SetNotificationsSignal(true, '', false, function () {})
  },
  markRead: function (fp) {
    this.SetNotificationsSignal(true, fp, false, function () {})
  },
  markAllAsRead: function (callback) {
    this.SetNotificationsSignal(false, '', true, callback)
  },
  IsInitialised: function () {
    return Initialised
  },
  /*----------  Methods for user signal actions  ----------*/
  /*
      Important thing here. We do have a RETRACT defined, but this is not defined for anything that goes into a bloom filter. Which means, you cannot retract an upvote, but you can downvote to reverse it. The reason why is that upvotes and downvotes (and elects) are aggregated, therefore after they get added to the bloom filter, we only know probabilistically that they're there. We have two bloom filters for each so we can have a +1 and -1, but adding 0 means adding another bloom filter in. Depending on the demand for a retract we can add a third bloom to the implementation to keep tracking of that, but bloom filters are very expensive because they're per-entity, and we have a lot of entities.
  
      This does not apply to non-aggregated signals like reporting to mod, those are kept instact and individual, and they can be retracted.
    */
  Upvote: function (targetfp, priorfp, boardfp, threadfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Upvote', priorfp.length ? 'Edit' : 'Create')
  },
  Downvote: function (targetfp, priorfp, boardfp, threadfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Downvote', priorfp.length ? 'Edit' : 'Create')
  },
  ReportToMod: function (
    targetfp,
    priorfp,
    reason,
    boardfp,
    threadfp,
    callback
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ReportToMod', priorfp.length ? 'Edit' : 'Create')
  },
  // ModDelete instead of ModBlock, to keep it more human-meaningful.
  ModDelete: function (targetfp, priorfp, reason, boardfp, threadfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModDelete', priorfp.length ? 'Edit' : 'Create')
  },
  ModApprove: function (
    targetfp,
    priorfp,
    reason,
    boardfp,
    threadfp,
    callback
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModApprove', priorfp.length ? 'Edit' : 'Create')
  },
  ModIgnore: function (targetfp, priorfp, reason, boardfp, threadfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('ModIgnore', priorfp.length ? 'Edit' : 'Create')
  },
  Follow: function (targetfp, priorfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Follow', priorfp.length ? 'Edit' : 'Create')
  },
  Block: function (targetfp, priorfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Block', priorfp.length ? 'Edit' : 'Create')
  },
  Elect: function (targetfp, priorfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Elect', priorfp.length ? 'Edit' : 'Create')
  },
  Disqualify: function (targetfp, priorfp, callback) {
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
    var metrics = require('../metrics/metrics')()
    metrics.SendSignalEvent('Disqualify', priorfp.length ? 'Edit' : 'Create')
  },
  /*----------  Base signal event action.  ----------*/
  sendSignalEvent: function (
    targetfp,
    priorfp,
    typeclass,
    typ,
    targettype,
    signaltext,
    boardfp,
    threadfp,
    callback
  ) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var now, req, e, localUser
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('Send Signal Event base triggered.')
              now = Math.floor(Date.now() / 1000)
              req = new pmessages.SignalEventPayload()
              e = new pmessages.Event()
              localUser = require('../../store/index').default.state.localUser
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
              feAPIConsumer.sendSignalEvent(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp)
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  /*----------  Methods for content event actions  ----------*/
  /*
      These are things like creating or editing entities that the user has created. If a priorfp is provided, it is an update. If not, it is a create.
    */
  SendBoardContent: function (priorfp, boarddata, callback) {
    this.sendContentEvent(
      priorfp,
      boarddata,
      undefined,
      undefined,
      undefined,
      callback
    )
    var metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Board', priorfp.length ? 'Edit' : 'Create')
  },
  SendThreadContent: function (priorfp, threaddata, callback) {
    this.sendContentEvent(
      priorfp,
      undefined,
      threaddata,
      undefined,
      undefined,
      callback
    )
    var metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Thread', priorfp.length ? 'Edit' : 'Create')
  },
  SendPostContent: function (priorfp, postdata, callback) {
    this.sendContentEvent(
      priorfp,
      undefined,
      undefined,
      postdata,
      undefined,
      callback
    )
    var metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('Post', priorfp.length ? 'Edit' : 'Create')
  },
  SendUserContent: function (priorfp, userdata, callback) {
    this.sendContentEvent(
      priorfp,
      undefined,
      undefined,
      undefined,
      userdata,
      callback
    )
    var metrics = require('../metrics/metrics')()
    metrics.SendContentEvent('User', priorfp.length ? 'Edit' : 'Create')
  },
  /*----------  Base content event action.  ----------*/
  sendContentEvent: function (
    priorfp,
    boarddata,
    threaddata,
    postdata,
    userdata,
    callback
  ) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var now, req, e, localUser, globalMethods
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('Send Content Event base triggered.')
              now = Math.floor(Date.now() / 1000)
              req = new pmessages.ContentEventPayload()
              e = new pmessages.Event()
              localUser = require('../../store/index').default.state.localUser
              globalMethods = require('../globals/methods')
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
              feAPIConsumer.sendContentEvent(req, function (err, resp) {
                if (err) {
                  console.log(err)
                } else {
                  callback(resp)
                }
              })
              return [2 /*return*/]
          }
        })
      })
    })
  },
  /*----------  Search types  ----------*/
  SendCommunitySearchRequest: function (query, callback) {
    this.sendSearchRequest('Board', query, callback)
  },
  SendContentSearchRequest: function (query, callback) {
    this.sendSearchRequest('Content', query, callback)
  },
  SendUserSearchRequest: function (query, callback) {
    this.sendSearchRequest('User', query, callback)
  },
  /*----------  Base search request action  ----------*/
  sendSearchRequest: function (searchType, query, callback) {
    WaitUntilFrontendReady(function () {
      return __awaiter(this, void 0, void 0, function () {
        var req
        return __generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!!Initialised) return [3 /*break*/, 2]
              return [4 /*yield*/, ExportedMethods.Initialise()]
            case 1:
              _a.sent()
              _a.label = 2
            case 2:
              console.log('sendSearchRequest triggered.')
              req = new pmessages.SearchRequestPayload()
              req.setSearchtype(searchType)
              req.setSearchquery(query)
              try {
                feAPIConsumer.sendSearchRequest(req, function (err, resp) {
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
              return [2 /*return*/]
          }
        })
      })
    })
  },
}
module.exports = ExportedMethods
function WaitUntilFrontendReady(cb) {
  function check() {
    return __awaiter(this, void 0, void 0, function () {
      var initialised
      return __generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            return [
              4 /*yield*/,
              ipc.callMain('GetFrontendClientConnInitialised'),
              // console.log(initialised)
            ]
          case 1:
            initialised = _a.sent()
            // console.log(initialised)
            if (!initialised) {
              // console.log("Frontend still not ready, waiting a little more...")
              return [2 /*return*/, setTimeout(check, 333)]
            } else {
              return [2 /*return*/, cb()]
            }
            return [2 /*return*/]
        }
      })
    })
  }
  return check()
}
//# sourceMappingURL=feapiconsumer.js.map
