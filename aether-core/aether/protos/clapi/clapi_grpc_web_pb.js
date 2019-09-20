/**
 * @fileoverview gRPC-Web generated client stub for clapi
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js')

var mimapi_mimapi_pb = require('../mimapi/mimapi_pb.js')
const proto = {};
proto.clapi = require('./clapi_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.clapi.ClientAPIClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.clapi.ClientAPIPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.FEReadyRequest,
 *   !proto.clapi.FEReadyResponse>}
 */
const methodInfo_ClientAPI_FrontendReady = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.FEReadyResponse,
  /** @param {!proto.clapi.FEReadyRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.FEReadyResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.FEReadyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.FEReadyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.FEReadyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.frontendReady =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/FrontendReady',
      request,
      metadata || {},
      methodInfo_ClientAPI_FrontendReady,
      callback);
};


/**
 * @param {!proto.clapi.FEReadyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.FEReadyResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.frontendReady =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/FrontendReady',
      request,
      metadata || {},
      methodInfo_ClientAPI_FrontendReady);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.AmbientsRequest,
 *   !proto.clapi.AmbientsResponse>}
 */
const methodInfo_ClientAPI_DeliverAmbients = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.AmbientsResponse,
  /** @param {!proto.clapi.AmbientsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.AmbientsResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.AmbientsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.AmbientsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.AmbientsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.deliverAmbients =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/DeliverAmbients',
      request,
      metadata || {},
      methodInfo_ClientAPI_DeliverAmbients,
      callback);
};


/**
 * @param {!proto.clapi.AmbientsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.AmbientsResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.deliverAmbients =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/DeliverAmbients',
      request,
      metadata || {},
      methodInfo_ClientAPI_DeliverAmbients);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.AmbientStatusPayload,
 *   !proto.clapi.AmbientStatusResponse>}
 */
const methodInfo_ClientAPI_SendAmbientStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.AmbientStatusResponse,
  /** @param {!proto.clapi.AmbientStatusPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.AmbientStatusResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.AmbientStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.AmbientStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.AmbientStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendAmbientStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendAmbientStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendAmbientStatus,
      callback);
};


/**
 * @param {!proto.clapi.AmbientStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.AmbientStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendAmbientStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendAmbientStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendAmbientStatus);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.AmbientLocalUserEntityPayload,
 *   !proto.clapi.AmbientLocalUserEntityResponse>}
 */
const methodInfo_ClientAPI_SendAmbientLocalUserEntity = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.AmbientLocalUserEntityResponse,
  /** @param {!proto.clapi.AmbientLocalUserEntityPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.AmbientLocalUserEntityResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.AmbientLocalUserEntityPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.AmbientLocalUserEntityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.AmbientLocalUserEntityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendAmbientLocalUserEntity =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendAmbientLocalUserEntity',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendAmbientLocalUserEntity,
      callback);
};


/**
 * @param {!proto.clapi.AmbientLocalUserEntityPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.AmbientLocalUserEntityResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendAmbientLocalUserEntity =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendAmbientLocalUserEntity',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendAmbientLocalUserEntity);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.HomeViewPayload,
 *   !proto.clapi.HomeViewResponse>}
 */
const methodInfo_ClientAPI_SendHomeView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.HomeViewResponse,
  /** @param {!proto.clapi.HomeViewPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.HomeViewResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.HomeViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.HomeViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.HomeViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendHomeView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendHomeView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendHomeView,
      callback);
};


/**
 * @param {!proto.clapi.HomeViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.HomeViewResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendHomeView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendHomeView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendHomeView);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.PopularViewPayload,
 *   !proto.clapi.PopularViewResponse>}
 */
const methodInfo_ClientAPI_SendPopularView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.PopularViewResponse,
  /** @param {!proto.clapi.PopularViewPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.PopularViewResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.PopularViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.PopularViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.PopularViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendPopularView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendPopularView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendPopularView,
      callback);
};


/**
 * @param {!proto.clapi.PopularViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.PopularViewResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendPopularView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendPopularView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendPopularView);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.NewViewPayload,
 *   !proto.clapi.NewViewResponse>}
 */
const methodInfo_ClientAPI_SendNewView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.NewViewResponse,
  /** @param {!proto.clapi.NewViewPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.NewViewResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.NewViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.NewViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.NewViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendNewView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendNewView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendNewView,
      callback);
};


/**
 * @param {!proto.clapi.NewViewPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.NewViewResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendNewView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendNewView',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendNewView);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.NotificationsPayload,
 *   !proto.clapi.NotificationsResponse>}
 */
const methodInfo_ClientAPI_SendNotifications = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.NotificationsResponse,
  /** @param {!proto.clapi.NotificationsPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.NotificationsResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.NotificationsPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.NotificationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.NotificationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendNotifications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendNotifications',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendNotifications,
      callback);
};


/**
 * @param {!proto.clapi.NotificationsPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.NotificationsResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendNotifications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendNotifications',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendNotifications);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.OnboardCompleteStatusPayload,
 *   !proto.clapi.OnboardCompleteStatusResponse>}
 */
const methodInfo_ClientAPI_SendOnboardCompleteStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.OnboardCompleteStatusResponse,
  /** @param {!proto.clapi.OnboardCompleteStatusPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.OnboardCompleteStatusResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.OnboardCompleteStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.OnboardCompleteStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.OnboardCompleteStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendOnboardCompleteStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendOnboardCompleteStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendOnboardCompleteStatus,
      callback);
};


/**
 * @param {!proto.clapi.OnboardCompleteStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.OnboardCompleteStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendOnboardCompleteStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendOnboardCompleteStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendOnboardCompleteStatus);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.ModModeEnabledStatusPayload,
 *   !proto.clapi.ModModeEnabledStatusResponse>}
 */
const methodInfo_ClientAPI_SendModModeEnabledStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.ModModeEnabledStatusResponse,
  /** @param {!proto.clapi.ModModeEnabledStatusPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.ModModeEnabledStatusResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.ModModeEnabledStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.ModModeEnabledStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.ModModeEnabledStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendModModeEnabledStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendModModeEnabledStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendModModeEnabledStatus,
      callback);
};


/**
 * @param {!proto.clapi.ModModeEnabledStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.ModModeEnabledStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendModModeEnabledStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendModModeEnabledStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendModModeEnabledStatus);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.ExternalContentAutoloadDisabledStatusPayload,
 *   !proto.clapi.ExternalContentAutoloadDisabledStatusResponse>}
 */
const methodInfo_ClientAPI_SendExternalContentAutoloadDisabledStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.ExternalContentAutoloadDisabledStatusResponse,
  /** @param {!proto.clapi.ExternalContentAutoloadDisabledStatusPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.ExternalContentAutoloadDisabledStatusResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.ExternalContentAutoloadDisabledStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.ExternalContentAutoloadDisabledStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.ExternalContentAutoloadDisabledStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendExternalContentAutoloadDisabledStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendExternalContentAutoloadDisabledStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendExternalContentAutoloadDisabledStatus,
      callback);
};


/**
 * @param {!proto.clapi.ExternalContentAutoloadDisabledStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.ExternalContentAutoloadDisabledStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendExternalContentAutoloadDisabledStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendExternalContentAutoloadDisabledStatus',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendExternalContentAutoloadDisabledStatus);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.clapi.SearchResultPayload,
 *   !proto.clapi.SearchResultResponse>}
 */
const methodInfo_ClientAPI_SendSearchResult = new grpc.web.AbstractClientBase.MethodInfo(
  proto.clapi.SearchResultResponse,
  /** @param {!proto.clapi.SearchResultPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.clapi.SearchResultResponse.deserializeBinary
);


/**
 * @param {!proto.clapi.SearchResultPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.clapi.SearchResultResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.clapi.SearchResultResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.clapi.ClientAPIClient.prototype.sendSearchResult =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/clapi.ClientAPI/SendSearchResult',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendSearchResult,
      callback);
};


/**
 * @param {!proto.clapi.SearchResultPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.clapi.SearchResultResponse>}
 *     A native promise that resolves to the response
 */
proto.clapi.ClientAPIPromiseClient.prototype.sendSearchResult =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/clapi.ClientAPI/SendSearchResult',
      request,
      metadata || {},
      methodInfo_ClientAPI_SendSearchResult);
};


module.exports = proto.clapi;

