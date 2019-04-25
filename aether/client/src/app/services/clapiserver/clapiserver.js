"use strict";
// Client > ClientAPIServer
// This file is the grpc server we want to use to talk to the frontend.
Object.defineProperty(exports, "__esModule", { value: true });
// Imports
var grpc = require('grpc');
// const resolve = require('path').resolve
// let globals = require('../globals/globals')
var feapiconsumer = require('../feapiconsumer/feapiconsumer');
var ipc = require('../../../../node_modules/electron-better-ipc');
var vuexStore = require('../../store/index').default;
// // Load the proto file
// const proto = grpc.load({
//   file: 'clapi/clapi.proto',
//   root: resolve(__dirname, '../protos')
// }).clapi
var messages = require('../../../../../protos/clapi/clapi_pb.js');
var services = require('../../../../../protos/clapi/clapi_grpc_pb');
/**
 Client-side GRPC server so that the frontend can talk to the client. This is useful at the first start where the Frontend needs to start its own GRPC server and return its address to the client.
 */
function StartClientAPIServer() {
    var server = new grpc.Server();
    server.addService(services.ClientAPIService, {
        frontendReady: FrontendReady,
        deliverAmbients: DeliverAmbients,
        sendAmbientStatus: SendAmbientStatus,
        sendAmbientLocalUserEntity: SendAmbientLocalUserEntity,
        sendHomeView: SendHomeView,
        sendPopularView: SendPopularView,
        sendNotifications: SendNotifications,
        sendOnboardCompleteStatus: SendOnboardCompleteStatus,
        sendModModeEnabledStatus: SendModModeEnabledStatus,
        sendExternalContentAutoloadDisabledStatus: SendExternalContentAutoloadDisabledStatus,
        sendSearchResult: SendSearchResult,
    });
    var boundPort = server.bind('127.0.0.1:0', grpc.ServerCredentials.createInsecure());
    server.start();
    return boundPort;
}
exports.StartClientAPIServer = StartClientAPIServer;
function FrontendReady(req, callback) {
    var r = req.request.toObject();
    console.log('frontend ready at: ', r.address, ':', r.port);
    // globals.FrontendReady = true
    ipc.callMain('SetFrontendReady', true);
    // globals.FrontendAPIPort = req.request.port
    ipc.callMain('SetFrontendAPIPort', r.port);
    feapiconsumer.Initialise();
    var resp = new messages.FEReadyResponse();
    callback(null, resp);
}
function DeliverAmbients(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setAmbientBoards', r.boardsList);
    var resp = new messages.AmbientsResponse();
    callback(null, resp);
}
function SendAmbientStatus(req, callback) {
    var r = req.request.toObject();
    // console.log(r)
    vuexStore.dispatch('setAmbientStatus', r);
    var resp = new messages.AmbientStatusResponse();
    callback(null, resp);
}
function SendAmbientLocalUserEntity(req, callback) {
    var r = req.request.toObject();
    // console.log(r)
    vuexStore.dispatch('setAmbientLocalUserEntity', r);
    var resp = new messages.AmbientLocalUserEntityResponse();
    callback(null, resp);
}
function SendHomeView(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setHomeView', r.threadsList);
    var resp = new messages.HomeViewResponse();
    callback(null, resp);
}
function SendPopularView(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setPopularView', r.threadsList);
    var resp = new messages.PopularViewResponse();
    callback(null, resp);
}
function SendNotifications(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setNotifications', r);
    var resp = new messages.NotificationsResponse();
    callback(null, resp);
}
function SendOnboardCompleteStatus(req, callback) {
    var r = req.request.toObject();
    /*======================================================
    =            CLIENT VERSION / UPGRADE LOGIC            =
    ======================================================*/
    /*
      Q: Why is this even here?
  
      A: OnboardComplete status and Client version both affects which page needs to be shown. Onboardcomplete takes priority, because if the user is not onboarded, then the updates do not really matter - the onboarding has to happen first.
  
      Having this here ensures that the first open after update check will happen after the onboarding check happens. If the onboard hasn't happened yet, there is no point redirecting to the update check.
    */
    /*
      At the end of the init, send in the client version. The FE will respond with the last known version.
      - If the FE response is empty, this is a first-ever boot.
      - If the FE response isn't empty, but they're not the same, this is the first start after a version upgrade.
    */
    var versionAndBuild = require('electron').remote.app.getVersion();
    var firstEverOpen = false;
    var firstOpenAfterSuccessfulUpdate = false;
    feapiconsumer.SendClientVersion(versionAndBuild, function (resp) {
        if (resp.lastknownclientversion.length === 0) {
            console.log('This is the first time this app was ever open.');
            firstEverOpen = true;
        }
        if (resp.lastknownclientversion != versionAndBuild) {
            console.log('First open after successful version update.');
            firstOpenAfterSuccessfulUpdate = true;
        }
        // console.log('NOT a first open after successful version update.')
        if (r.onboardcomplete) {
            /*
            This logic runs only if the onboarding is complete. This means the new version first run bit won't be set if the onboarding isn't done - there is no reason to ever show the banner on changelog which says it has new tricks, because the user hasn't seen the old 'tricks' yet.
            */
            if (firstEverOpen) {
                var metrics_1 = require('../metrics/metrics')();
                metrics_1.SendRaw('App first-ever opened');
                return;
            }
            if (firstOpenAfterSuccessfulUpdate) {
                var metrics_2 = require('../metrics/metrics')();
                metrics_2.SendRaw('App update successful');
                vuexStore.dispatch('setFirstRunAfterUpdateState', true);
                var router = require('../../renderermain').router;
                router.push('/changelog');
            }
            else {
                vuexStore.dispatch('setFirstRunAfterUpdateState', false);
            }
        }
        else {
            vuexStore.dispatch('setFirstRunAfterUpdateState', false);
        }
    });
    /*=====  End of CLIENT VERSION / UPGRADE LOGIC  ======*/
    vuexStore.dispatch('setOnboardCompleteStatus', r.onboardcomplete);
    var resp = new messages.OnboardCompleteStatusResponse();
    callback(null, resp);
}
function SendModModeEnabledStatus(req, callback) {
    var r = req.request.toObject();
    vuexStore.dispatch('setModModeEnabledStatus', r.modmodeenabled);
    var resp = new messages.ModModeEnabledStatusResponse();
    callback(null, resp);
}
function SendExternalContentAutoloadDisabledStatus(req, callback) {
    console.log('external content autoload disabled status arrived.');
    var r = req.request.toObject();
    vuexStore.dispatch('setExternalContentAutoloadDisabledStatus', r.externalcontentautoloaddisabled);
    var resp = new messages.ExternalContentAutoloadDisabledStatusResponse();
    callback(null, resp);
}
function SendSearchResult(req, callback) {
    console.log('Search result arrived.');
    var r = req.request.toObject();
    vuexStore.dispatch('setSearchResult', r);
    var resp = new messages.SearchResultResponse();
    callback(null, resp);
}
//# sourceMappingURL=clapiserver.js.map