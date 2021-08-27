"use strict";
/*
This is the data store for the client.

This data store does not hold any persistent data, nor does it cache it. The point of this is to hold the instance data. The frontend is the actual caching and compile logic that regenerates the data to be used as needed.
*/
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
Object.defineProperty(exports, "__esModule", { value: true });
var isDev = require('electron-is-dev');
if (isDev) {
    var Vue = require('../../../node_modules/vue/dist/vue.js'); // Production
}
else {
    var Vue = require('../../../node_modules/vue/dist/vue.min.js'); // Production
}
var Vuex = require('../../../node_modules/vuex').default;
Vue.use(Vuex);
var ipc = require('../../../node_modules/electron-better-ipc').ipcMain;
var fe = require('../services/feapiconsumer/feapiconsumer');
var globalMethods = require('../services/globals/methods');
var dataLoaders = require('./dataloaders').default;
var statusLights = require('./statuslights').default;
var contentRelations = require('./contentrelations');
var crumbs = require('./crumbs');
var dataLoaderPlugin = function (store) {
    store.watch(
    // When the returned result changes,
    function (state) {
        return state.route.params;
    }, 
    // Run this callback
    function (newValue, oldValue) {
        console.log(store.state.route.name);
        fe.SendInflightsPruneRequest(function (resp) {
            console.log(resp);
        });
        /*
          Send the location to metrics
        */
        var metrics = require('../services/metrics/metrics')();
        metrics.SendRaw('App navigate', {
            'A-App-Location': store.state.route.name,
        });
        store.dispatch('registerNextMoveToHistoryCounter');
        // Set history state for the back / forward buttons.
        // store.dispatch('setHistoryHasForward', routerHistory.hasForward())
        // store.dispatch('setHistoryHasPrevious', routerHistory.hasPrevious())
        // Set the last load timestamp. We determine by looking at this the animation state some progress bar elements are in. (i.e. If it's a refresh or a return to a page, and we have a progress bar at 33% we don't want to animate from 0 to 33, just start it snapped there.)
        store.dispatch('setLastPageLoadTimestamp');
        // First, check if we should refresh.
        if (oldValue === newValue && !store.state.frontendHasUpdates) {
            // if the values are the same, and frontend has no updates, bail.
            return;
        }
        var routeParams = newValue;
        if (store.state.route.name === 'Board' ||
            store.state.route.name === 'Board>ThreadsNewList' ||
            store.state.route.name === 'Board>Elections') {
            store.dispatch('loadBoardScopeData', routeParams.boardfp);
            store.dispatch('setLastSeenForBoard', { fp: routeParams.boardfp });
            return;
        }
        if (store.state.route.name.includes('Onboard')) {
            /*
              If the route name includes the word 'onboard' (i.e. onboard1, onboard2 ...) we will add a check that makes it so that if the onboarding is complete, we redirect to popular.
            */
            if (store.state.onboardCompleteStatusArrived &&
                store.state.onboardCompleteStatus) {
                var router = require('../renderermain').router;
                router.push('/popular');
                return;
            }
        }
        if (store.state.route.name === 'Board>NewThread' ||
            store.state.route.name === 'Board>BoardInfo') {
            store.dispatch('loadBoardScopeData', routeParams.boardfp);
            return;
        }
        if (store.state.route.name === 'Board>Reports') {
            store.dispatch('setCurrentBoardReportsArrived', false);
            store.dispatch('loadBoardScopeData', routeParams.boardfp);
            store.dispatch('loadBoardReports', routeParams.boardfp);
            // store.dispatch('setCurrentBoardReportsArrived', true)
            return;
        }
        if (store.state.route.name === 'Board>ModActivity') {
            console.log('modactivity triggered');
            store.dispatch('setCurrentBoardModActionsArrived', false);
            store.dispatch('loadBoardScopeData', routeParams.boardfp);
            store.dispatch('loadBoardModActions', routeParams.boardfp);
            // store.dispatch('setCurrentBoardReportsArrived', true)
            return;
        }
        if (store.state.route.name === 'Thread') {
            store.dispatch('loadThreadScopeData', {
                boardfp: routeParams.boardfp,
                threadfp: routeParams.threadfp,
            });
            return;
        }
        if (store.state.route.name === 'Global' ||
            store.state.route.name === 'Global>Subbed') {
            store.dispatch('loadGlobalScopeData');
            return;
        }
        if (store.state.route.name === 'User' ||
            store.state.route.name === 'User>Boards' ||
            store.state.route.name === 'User>Threads' ||
            store.state.route.name === 'User>Posts' ||
            store.state.route.name === 'User>Notifications') {
            store.dispatch('loadUserScopeData', {
                fp: routeParams.userfp,
                userreq: true,
                boardsreq: false,
                threadsreq: false,
                postsreq: false,
            });
            return;
        }
        // If none of the special cases, just trigger an update for breadcrumbs.
        store.dispatch('updateBreadcrumbs');
    });
};
var actions = __assign({ 
    /*----------  Refreshers  ----------*/
    /*
      These are smaller, less encompassing versions of the loaders, and they're meant to be used after the principal payload is brought in.
    */
    refreshCurrentBoardAndThreads: function (context, _a) {
        var boardfp = _a.boardfp, sortByNew = _a.sortByNew;
        fe.GetBoardAndThreads(boardfp, sortByNew, function (resp) {
            actions.pruneInflights();
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList);
            context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true);
        });
    },
    pruneInflights: function () {
        fe.SendInflightsPruneRequest(function () { });
    },
    refreshCurrentThreadAndPosts: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        fe.GetThreadAndPosts(boardfp, threadfp, function (resp) {
            actions.pruneInflights();
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_THREAD', resp.thread);
            context.commit('SET_CURRENT_THREADS_POSTS', resp.postsList);
            // context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', true)
        });
    },
    /*----------  Refreshers END  ----------*/
    // within any of those, context.state is how you access state above.
    setSidebarState: function (context, sidebarOpen) {
        context.commit('SET_SIDEBAR_STATE', sidebarOpen);
    },
    setAmbientBoards: function (context, ambientBoards) {
        context.commit('SET_AMBIENT_BOARDS', ambientBoards);
    },
    setAmbientStatus: function (context, ambientStatus) {
        context.commit('SET_AMBIENT_STATUS', ambientStatus);
        // If any of the items in the ambient status is a board that has just been created, add it to subscribed board.
        for (var _i = 0, _a = context.state.ambientStatus.inflights.boardsList; _i < _a.length; _i++) {
            var val = _a[_i];
            if (val.status.eventtype === 'CREATE' &&
                val.status.completionpercent === 100) {
                actions.subToBoard(context, {
                    fp: val.entity.provable.fingerprint,
                    notify: true,
                });
            }
        }
        // At every refresh, render these dot states.
        actions.setDotStates(context, context.state.ambientStatus);
    },
    setAmbientLocalUserEntity: function (context, ambientLocalUserEntityPayload) {
        context.commit('SET_AMBIENT_LOCAL_USER_ENTITY', ambientLocalUserEntityPayload);
    },
    setCurrentBoardFp: function (context, fp) {
        var sortByNew = false;
        if (context.state.route.name === 'Board>ThreadsNewList') {
            sortByNew = true;
        }
        context.dispatch('setCurrentBoardAndThreads', {
            boardfp: fp,
            sortByNew: sortByNew,
        });
        context.commit('SET_CURRENT_BOARD_FP', fp);
        // // current board fp is the same as what we asked for, but FE has updates.
        // if (context.state.frontendHasUpdates) {
        //   context.dispatch('setCurrentBoardAndThreads', fp)
        //   context.commit('SET_CURRENT_BOARD_FP', fp)
        // }
    },
    setCurrentThreadFp: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        context.dispatch('setCurrentThreadAndPosts', {
            boardfp: boardfp,
            threadfp: threadfp,
        });
        context.commit('SET_CURRENT_THREAD_FP', threadfp);
        // Same scope, but fe has updated.
        return;
        // current thread fp is the same as what we asked for, but FE has updates.
        // if (context.state.frontendHasUpdates) {
        //   context.dispatch('setCurrentThreadAndPosts',
        //     { boardfp: boardfp, threadfp: threadfp })
        //   context.commit('SET_CURRENT_THREAD_FP', threadfp)
        // }
    },
    setCurrentBoardAndThreads: function (context, _a) {
        var boardfp = _a.boardfp, sortByNew = _a.sortByNew;
        if (context.state.currentBoardFp === boardfp) {
            fe.GetBoardAndThreads(boardfp, sortByNew, function (resp) {
                context.commit('SET_CURRENT_BOARD', resp.board);
                context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList);
                context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true);
                context.dispatch('updateBreadcrumbs');
            });
            return;
            // If we're already here, update board but without false/true current board load complete flash.
        }
        context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', false);
        fe.GetBoardAndThreads(boardfp, sortByNew, function (resp) {
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList);
            context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true);
            context.dispatch('updateBreadcrumbs');
            // ^ This has to be here because otherwise the BC compute process runs before the data is ready, resulting in empty breadcrumbs.
        });
    },
    setCurrentThreadAndPosts: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        // if (context.state.currentThreadFp === threadfp) {
        //   context.dispatch('updateBreadcrumbs')
        //   return
        // }
        // ^ Nope, you're trying to be way too smart. The user might have updated the thread and if that's the case trying to do that will prevent that update from being visible. I tried doing it like above in boards where you load but without making it invisible, and what that does it creates a flash of old content. Not great. This is the best.
        context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', false);
        fe.GetThreadAndPosts(boardfp, threadfp, function (resp) {
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_THREAD', resp.thread);
            context.commit('SET_CURRENT_THREADS_POSTS', resp.postsList);
            context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', true);
            context.dispatch('updateBreadcrumbs');
        });
    },
    setLastPageLoadTimestamp: function (context) {
        context.commit('SET_LAST_PAGE_LOAD_TIMESTAMP', globalMethods.NowUnix());
    },
    /*----------  Views insertion  ----------*/
    setHomeView: function (context, threads) {
        context.commit('SET_HOME_VIEW', threads);
    },
    setPopularView: function (context, threads) {
        context.commit('SET_POPULAR_VIEW', threads);
    },
    setNewView: function (context, newView) {
        context.commit('SET_NEW_VIEW', newView);
    },
    setNotifications: function (context, response) {
        var payload = {
            notifications: response.notificationsList,
            unseenNotificationsPresent: false,
        };
        /*
          dev.6 change: Notifications light stays on so long as there are notifications present. This is simpler logic, and makes sure you don't forget about some unread stuff. Side effect: the light won't go out until you click on all of these. Future feature: add a 'mark read' button, and 'mark all as read' button.
        */
        for (var _i = 0, _a = payload.notifications; _i < _a.length; _i++) {
            var notification = _a[_i];
            if (!notification.read) {
                payload.unseenNotificationsPresent = true;
                break;
            }
        }
        // dev.5- code below.
        // for (let val of payload.notifications) {
        //   if (val.creationtimestamp > response.lastseen) {
        //     payload.unseenNotificationsPresent = true
        //     break
        //   }
        // }
        context.commit('SET_NOTIFICATIONS', payload);
    },
    setOnboardCompleteStatus: function (context, ocs) {
        if (ocs === false) {
            var router = require('../renderermain').router;
            router.push('/onboard');
        }
        context.commit('SET_ONBOARD_COMPLETE_STATUS', ocs);
    },
    setModModeEnabledStatus: function (context, modModeEnabled) {
        context.commit('SET_MOD_MODE_ENABLED_STATUS', modModeEnabled);
    },
    setExternalContentAutoloadDisabledStatus: function (context, externalContentAutoloadDisabled) {
        context.commit('SET_EXTERNAL_CONTENT_AUTOLOAD_DISABLED_STATUS', externalContentAutoloadDisabled);
    },
    /*----------  History state  ----------*/
    registerNextActionIsHistoryMoveForward: function (context) {
        context.commit('REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_FORWARD');
    },
    registerNextActionIsHistoryMoveBack: function (context) {
        context.commit('REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_BACK');
    },
    registerNextMoveToHistoryCounter: function (context) {
        context.commit('REGISTER_NEXT_MOVE_TO_HISTORY_COUNTER');
    },
    setFirstRunAfterUpdateState: function (context, fraus) {
        context.commit('SET_FIRST_RUN_AFTER_UPDATE_STATE', fraus);
    },
    saveDraft: function (context, draftContainer) {
        context.commit('SAVE_DRAFT', draftContainer);
    } }, statusLights, dataLoaders, crumbs.crumbActions, contentRelations.actions, { 
    /*----------  Search  ----------*/
    setSearchResult: function (context, response) {
        context.commit('SAVE_SEARCH_RESULT', response);
    } });
var mutations = __assign({ SET_SIDEBAR_STATE: function (state, sidebarOpen) {
        state.sidebarOpen = sidebarOpen;
    },
    SET_AMBIENT_BOARDS: function (state, ambientBoards) {
        ambientBoards.sort(function (a, b) {
            return b.usercount - a.usercount;
        });
        state.ambientBoards = ambientBoards;
        state.ambientBoardsArrived = true;
    },
    SET_AMBIENT_STATUS: function (state, ambientStatus) {
        if (!globalMethods.IsUndefined(ambientStatus.frontendambientstatus)) {
            state.ambientStatus.frontendambientstatus =
                ambientStatus.frontendambientstatus;
        }
        if (!globalMethods.IsUndefined(ambientStatus.backendambientstatus)) {
            state.ambientStatus.backendambientstatus =
                ambientStatus.backendambientstatus;
        }
        if (!globalMethods.IsUndefined(ambientStatus.inflights)) {
            state.ambientStatus.inflights = ambientStatus.inflights;
        }
    },
    SET_AMBIENT_LOCAL_USER_ENTITY: function (state, payload) {
        state.localUserArrived = true;
        // ^ Always set to true since some message arrived.
        state.localUserExists = payload.localuserexists;
        state.localUser = payload.localuserentity;
    },
    SET_CURRENT_BOARD_FP: function (state, fp) {
        state.currentBoardFp = fp;
    },
    SET_CURRENT_THREAD_FP: function (state, fp) {
        state.currentThreadFp = fp;
    },
    SET_CURRENT_BOARD: function (state, board) {
        state.currentBoard = board;
    },
    SET_CURRENT_THREAD: function (state, thread) {
        state.currentThread = thread;
    },
    SET_CURRENT_BOARDS_THREADS: function (state, threads) {
        state.currentBoardsThreads = threads;
    },
    SET_CURRENT_THREADS_POSTS: function (state, posts) {
        state.currentThreadsPosts = posts;
    },
    SET_ALL_BOARDS: function (state, boards) {
        boards.sort(function (a, b) {
            return b.usercount - a.usercount;
        });
        state.allBoards = boards;
    },
    SET_USER_SCOPE_DATA: function (state, resp) {
        if (resp.userentityrequested) {
            state.currentUserEntity = resp.user;
        }
        if (resp.boardsrequested) {
            state.currentUserBoards = resp.boards;
        }
        if (resp.threadsrequested) {
            state.currentUserThreads = resp.threads;
        }
        if (resp.postsrequested) {
            state.currentUserPosts = resp.posts;
        }
    } }, crumbs.crumbMutations, contentRelations.mutations, { 
    /*----------  Loader mutations that mark a pull done  ----------*/
    /*
      These are important because when these are complete and there is no data, we know that we should show a 404. These only apply to singular entities, not lists, so effectively thread view, board view, user view.
    */
    SET_ALL_BOARDS_LOAD_COMPLETE: function (state, loadComplete) {
        state.allBoardsLoadComplete = loadComplete;
    },
    SET_CURRENT_BOARD_LOAD_COMPLETE: function (state, loadComplete) {
        state.currentBoardLoadComplete = loadComplete;
    },
    SET_CURRENT_THREAD_LOAD_COMPLETE: function (state, loadComplete) {
        state.currentThreadLoadComplete = loadComplete;
    },
    SET_CURRENT_USER_LOAD_COMPLETE: function (state, loadComplete) {
        state.currentUserLoadComplete = loadComplete;
    },
    SET_LAST_PAGE_LOAD_TIMESTAMP: function (state, lastPageLoadTimestamp) {
        state.lastPageLoadTimestamp = lastPageLoadTimestamp;
    },
    SET_DOT_STATES: function (state, dotStates) {
        state.dotStates = dotStates;
    },
    /*----------  Views insertion mutations  ----------*/
    SET_HOME_VIEW: function (state, threads) {
        state.homeViewThreads = threads;
        state.homeViewArrived = true;
    },
    SET_POPULAR_VIEW: function (state, threads) {
        state.popularViewThreads = threads;
        state.popularViewArrived = true;
    },
    SET_NEW_VIEW: function (state, newView) {
        state.newView = newView;
        state.newViewArrived = true;
    },
    SET_NOTIFICATIONS: function (state, payload) {
        state.notifications = payload.notifications;
        state.unseenNotificationsPresent = payload.unseenNotificationsPresent;
        state.notificationsArrived = true;
        /*----------  Handle OS notifications  ----------*/
        var unreads = [];
        var notification = {};
        var localUserFp = '';
        if (state.localUserExists && state.localUserArrived) {
            localUserFp = state.localUser.fingerprint;
        }
        if (localUserFp.length === 0) {
            return;
            // If there's no local user or isn't yet present, we raise no notifications.
        }
        for (var _i = 0, _a = state.notifications; _i < _a.length; _i++) {
            var val = _a[_i];
            if (!val.read) {
                // Did we raise this OS notification before? Check the transient list
                var vAsStr = JSON.stringify(val);
                for (var _b = 0, _c = state.alreadyRaisedOSNotifications; _b < _c.length; _b++) {
                    var v = _c[_b];
                    if (v === vAsStr) {
                        return;
                    }
                }
                unreads.push(val);
                state.alreadyRaisedOSNotifications.push(vAsStr);
            }
        }
        if (unreads.length > 1) {
            notification = new Notification('New notifications', {
                body: 'You have ' + unreads.length + ' unread notifications on Aether.',
            });
            notification.onclick = function () {
                var router = require('../renderermain').router;
                router.push('/user/' + localUserFp + '/notifications');
                ipc.callMain('FocusAndShow');
            };
        }
        if (unreads.length === 1) {
            // Add the name of the user to the notification
            var user = '';
            if (typeof unreads[0].responsepostsusersList[0] !== 'undefined' &&
                typeof unreads[0].responsepostsusersList[0].username !== 'undefined') {
                user = '@' + unreads[0].responsepostsusersList[0].username;
            }
            notification = new Notification('New notification', {
                body: user + ' ' + unreads[0].text,
            });
            if (user.length === 0) {
                /*
                  This might happen if the user's name has not yet arrived at the point of notification creation.
                */
                notification = new Notification('New notification', {
                    body: 'You have one unread notification on Aether.',
                });
            }
            notification.onclick = function () {
                var router = require('../renderermain').router;
                router.push('/user/' + localUserFp + '/notifications');
                ipc.callMain('FocusAndShow');
            };
        }
        /*---------- END Handle OS notifications  ----------*/
    },
    SET_ONBOARD_COMPLETE_STATUS: function (state, ocs) {
        state.onboardCompleteStatus = ocs;
        state.onboardCompleteStatusArrived = true;
    },
    SET_MOD_MODE_ENABLED_STATUS: function (state, modModeEnabled) {
        state.modModeEnabled = modModeEnabled;
        state.modModeEnabledArrived = true;
    },
    SET_EXTERNAL_CONTENT_AUTOLOAD_DISABLED_STATUS: function (state, externalContentAutoloadDisabled) {
        // If there's a change, apply the new whitelist.
        if (state.externalContentAutoloadDisabled != externalContentAutoloadDisabled) {
            // There *is* a change. Apply the change.
            if (externalContentAutoloadDisabled) {
                ipc.callMain('DisableExternalResourceAutoLoad');
            }
            else {
                ipc.callMain('EnableExternalResourceAutoLoad');
            }
        }
        state.externalContentAutoloadDisabled = externalContentAutoloadDisabled;
        state.externalContentAutoloadDisabledArrived = true;
    },
    REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_FORWARD: function (state) {
        state.historyNextActionType = 'HISTORY_BUTTON_MOVE_FORWARD';
    },
    REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_BACK: function (state) {
        state.historyNextActionType = 'HISTORY_BUTTON_MOVE_BACK';
    },
    REGISTER_NEXT_MOVE_TO_HISTORY_COUNTER: function (state) {
        // Regular nav
        if (state.historyNextActionType.length === 0) {
            // We advance current history caret and set max to the same. If you go a few pages via the back button and now click something, then the forward stack is gone.
            state.historyCurrentCaret++;
            state.historyMaxCaret = state.historyCurrentCaret;
            return;
        }
        // History back button
        if (state.historyNextActionType === 'HISTORY_BUTTON_MOVE_BACK') {
            // Only currentHistory caret moves back, max stays the same
            state.historyNextActionType = '';
            state.historyCurrentCaret > 0
                ? state.historyCurrentCaret--
                : (state.historyCurrentCaret = 0);
            return;
        }
        if (state.historyNextActionType === 'HISTORY_BUTTON_MOVE_FORWARD') {
            // Only currentHistory caret moves back, max stays the same
            state.historyNextActionType = '';
            state.historyCurrentCaret < state.historyMaxCaret
                ? state.historyCurrentCaret++
                : (state.historyCurrentCaret = state.historyMaxCaret);
        }
    },
    /*===============================
    =            Reports            =
    ===============================*/
    SET_CURRENT_BOARD_REPORTS: function (state, boardReports) {
        state.currentBoardsReports = boardReports;
        state.currentBoardsReportsArrived = true;
    },
    SET_CURRENT_BOARD_REPORTS_ARRIVED: function (state, arrived) {
        state.currentBoardsReportsArrived = arrived;
    },
    /*=====  End of Reports  ======*/
    /*===================================
    =            Mod actions            =
    ===================================*/
    SET_CURRENT_BOARD_MODACTIONS: function (state, boardModActions) {
        state.currentBoardsModActions = boardModActions;
        state.currentBoardsModActionsArrived = true;
    },
    SET_CURRENT_BOARD_MODACTIONS_ARRIVED: function (state, arrived) {
        state.currentBoardsModActionsArrived = arrived;
    },
    /*=====  End of Mod actions  ======*/
    SET_FIRST_RUN_AFTER_UPDATE_STATE: function (state, fraus) {
        state.firstRunAfterUpdate = fraus;
        state.firstRunAfterUpdateStatusArrived = true;
    },
    SAVE_DRAFT: function (state, draft) {
        // drafts: parentfp > contenttype > fieldtype: draftcontent
        /*
          We start to construct the tree from the bottom.
        */
        // Content type, there is only one. We set the content type map to the drafts map on parentfp key.
        var parentFpDrafts = new Map();
        if (state.drafts.has(draft.parentFp)) {
            parentFpDrafts = state.drafts.get(draft.parentFp);
        }
        parentFpDrafts.set(draft.contentType, draft.fields);
        // And finally, set the parentfp key on the drafts object with our updated map.
        state.drafts.set(draft.parentFp, parentFpDrafts);
    },
    SAVE_SEARCH_RESULT: function (state, searchResult) {
        // state.console.log('result received: ')
        console.log(searchResult);
        if (searchResult.searchtype === 'Board') {
            state.boardsSearchResult = searchResult.boardsList;
        }
        if (searchResult.searchtype === 'Content') {
            state.threadsSearchResult = searchResult.threadsList;
            state.postsSearchResult = searchResult.postsList;
        }
        if (searchResult.searchtype === 'User') {
            state.usersSearchResult = searchResult.usersList;
        }
    } });
/*

registerNextActionIsHistoryMoveForward(context: any) {
  context.commit('REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_FORWARD')
},
registerNextActionIsHistoryMovePrevious(context: any) {
  context.commit('REGISTER_NEXT_ACTION_IS_HISTORY_MOVE_PREVIOUS')
},
registerNextMoveToHistoryCounter(context: any) {
  context.commit('REGISTER_NEXT_MOVE_TO_HISTORY_COUNTER')
},

*/
var st = new Vuex.Store({
    state: {
        /*----------  All boards main  ----------*/
        allBoards: [],
        allBoardsLoadComplete: false,
        /*----------  Current board main  ----------*/
        currentBoard: {},
        currentBoardFp: '',
        currentBoardLoadComplete: false,
        /*----------  Current board sub data  ----------*/
        currentBoardsThreads: [],
        currentBoardsReports: [],
        currentBoardsReportsArrived: false,
        currentBoardsModActions: [],
        currentBoardsModActionsArrived: false,
        /*----------  Current thread main  ----------*/
        currentThread: {},
        currentThreadFp: '',
        currentThreadLoadComplete: false,
        /*----------  Current thread sub data  ----------*/
        currentThreadsPosts: [],
        currentUserBoards: [],
        /*----------  Current user main  ----------*/
        currentUserEntity: {},
        currentUserLoadComplete: false,
        /*----------  Current user sub data  ----------*/
        currentUserPosts: [],
        currentUserThreads: [],
        /*----------  Ambient data pushed in from frontend  ----------*/
        ambientBoards: {},
        ambientBoardsArrived: false,
        /*
          It's quite important that the schema below is available in JS here, not just in protobuf. Because if it is not, it will not be reactive - Vue needs to know beforehand which items you want it to track.
        */
        ambientStatus: {
            backendambientstatus: {
                /*----------  Bootstrap  ----------*/
                bootstrapinprogress: false,
                lastbootstraptimestamp: 0,
                /*----------  Network  ----------*/
                inboundscount15: 0,
                lastinboundconntimestamp: 0,
                lastoutboundconntimestamp: 0,
                lastoutbounddurationseconds: 0,
                outboundscount15: 0,
                localnodeexternalip: '',
                localnodeexternalport: 0,
                upnpstatus: '',
                /*----------  Database  ----------*/
                databasestatus: '',
                dbsizemb: 0,
                lastdbinserttimestamp: 0,
                lastinsertdurationseconds: 0,
                maxdbsizemb: 0,
                /*----------  Caching  ----------*/
                cachingstatus: '',
                lastcachegenerationdurationseconds: 0,
                lastcachegenerationtimestamp: 0,
                backendconfiglocation: '',
            },
            frontendambientstatus: {
                /*----------  Bootstrap  ----------*/
                bootstraprefreshcomplete: false,
                bootstraprefreshinprogress: false,
                /*----------  Refresh  ----------*/
                lastrefreshdurationseconds: 0,
                lastrefreshtimestamp: 0,
                refresherstatus: '',
                frontendconfiglocation: '',
                sfwlistdisabled: false,
            },
            inflights: {
                boardsList: [],
                threadsList: [],
                postsList: [],
                votesList: [],
                keysList: [],
                truststatesList: [],
            },
        },
        // States for the status dots visible at the bottom of the sidebar and in the status page.
        dotStates: {
            /*----------  Main dot statuses  ----------*/
            backendDotState: 'status_section_unknown',
            frontendDotState: 'status_section_unknown',
            /*----------  Sub dot states  ----------*/
            refresherDotState: 'status_subsection_unknown',
            inflightsDotState: 'status_subsection_unknown',
            networkDotState: 'status_subsection_unknown',
            dbDotState: 'status_subsection_unknown',
            cachingDotState: 'status_subsection_unknown',
        },
        /*----------  Local user data  ----------*/
        localUser: {},
        localUserExists: false,
        localUserArrived: false,
        // ^ Did we ever get a payload from FE? Until this is true, you can hide unready parts.
        /*----------  Views payloads  ----------*/
        homeViewThreads: {},
        homeViewArrived: false,
        popularViewThreads: {},
        popularViewArrived: false,
        newView: {},
        newViewArrived: false,
        /*----------  Notifications  ----------*/
        notifications: [],
        notificationsArrived: false,
        unseenNotificationsPresent: false,
        alreadyRaisedOSNotifications: [],
        /*----------  Onboard status  ----------*/
        onboardCompleteStatus: false,
        onboardCompleteStatusArrived: false,
        /*----------  Mod mode enabled status  ----------*/
        modModeEnabled: false,
        modModeEnabledArrived: false,
        /*----------  External content autoload disabled status  ----------*/
        externalContentAutoloadDisabled: false,
        externalContentAutoloadEnabledDisabled: false,
        /*----------  History state  ----------*/
        historyMaxCaret: 0,
        historyCurrentCaret: 0,
        historyNextActionType: '',
        /*----------  App fullscreen state  ----------*/
        appIsFullscreen: false,
        /*----------  Auto update state  ----------*/
        newUpdateReady: false,
        firstRunAfterUpdate: false,
        firstRunAfterUpdateStatusArrived: false,
        /*----------  Metrics  ----------*/
        metricsDisabled: false,
        // ^ Metrics are enabled by default on pre-release builds (as the user is notified of, in the onboarding process of pre-release versions.)
        /*----------  Drafts  ----------*/
        drafts: new Map(),
        /*----------  Search result  ----------*/
        boardsSearchResult: [],
        threadsSearchResult: [],
        postsSearchResult: [],
        usersSearchResult: [],
        /* ----------  Misc  ----------*/
        frontendHasUpdates: true,
        frontendPort: 0,
        route: {},
        sidebarOpen: true,
        breadcrumbs: [],
        lastPageLoadTimestamp: 0,
    },
    actions: actions,
    mutations: mutations,
    plugins: [dataLoaderPlugin],
});
exports.default = st;
/*

Reminder:

changeTestData(context: any) {
  // console.log("yo yo yo ")
  context.commit('editTestData')
}

is the same as:

changeTestData({commit}) {
  // console.log("yo yo yo ")
  commit('editTestData')
}
*/
//# sourceMappingURL=index.js.map