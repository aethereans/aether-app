/**
 * @fileoverview gRPC-Web generated client stub for feapi
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js')

var mimapi_mimapi_pb = require('../mimapi/mimapi_pb.js')
const proto = {};
proto.feapi = require('./feapi_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.feapi.FrontendAPIClient =
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

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.feapi.FrontendAPIPromiseClient =
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

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.SetClientAPIServerPortRequest,
 *   !proto.feapi.SetClientAPIServerPortResponse>}
 */
const methodDescriptor_FrontendAPI_SetClientAPIServerPort = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SetClientAPIServerPort',
  grpc.web.MethodType.UNARY,
  proto.feapi.SetClientAPIServerPortRequest,
  proto.feapi.SetClientAPIServerPortResponse,
  /**
   * @param {!proto.feapi.SetClientAPIServerPortRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SetClientAPIServerPortResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.SetClientAPIServerPortRequest,
 *   !proto.feapi.SetClientAPIServerPortResponse>}
 */
const methodInfo_FrontendAPI_SetClientAPIServerPort = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.SetClientAPIServerPortResponse,
  /**
   * @param {!proto.feapi.SetClientAPIServerPortRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SetClientAPIServerPortResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.SetClientAPIServerPortRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.SetClientAPIServerPortResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.SetClientAPIServerPortResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.setClientAPIServerPort =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SetClientAPIServerPort',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetClientAPIServerPort,
      callback);
};


/**
 * @param {!proto.feapi.SetClientAPIServerPortRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.SetClientAPIServerPortResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.setClientAPIServerPort =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SetClientAPIServerPort',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetClientAPIServerPort);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.ThreadAndPostsRequest,
 *   !proto.feapi.ThreadAndPostsResponse>}
 */
const methodDescriptor_FrontendAPI_GetThreadAndPosts = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/GetThreadAndPosts',
  grpc.web.MethodType.UNARY,
  proto.feapi.ThreadAndPostsRequest,
  proto.feapi.ThreadAndPostsResponse,
  /**
   * @param {!proto.feapi.ThreadAndPostsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ThreadAndPostsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.ThreadAndPostsRequest,
 *   !proto.feapi.ThreadAndPostsResponse>}
 */
const methodInfo_FrontendAPI_GetThreadAndPosts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.ThreadAndPostsResponse,
  /**
   * @param {!proto.feapi.ThreadAndPostsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ThreadAndPostsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.ThreadAndPostsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.ThreadAndPostsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.ThreadAndPostsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.getThreadAndPosts =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/GetThreadAndPosts',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetThreadAndPosts,
      callback);
};


/**
 * @param {!proto.feapi.ThreadAndPostsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.ThreadAndPostsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.getThreadAndPosts =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/GetThreadAndPosts',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetThreadAndPosts);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BoardAndThreadsRequest,
 *   !proto.feapi.BoardAndThreadsResponse>}
 */
const methodDescriptor_FrontendAPI_GetBoardAndThreads = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/GetBoardAndThreads',
  grpc.web.MethodType.UNARY,
  proto.feapi.BoardAndThreadsRequest,
  proto.feapi.BoardAndThreadsResponse,
  /**
   * @param {!proto.feapi.BoardAndThreadsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardAndThreadsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BoardAndThreadsRequest,
 *   !proto.feapi.BoardAndThreadsResponse>}
 */
const methodInfo_FrontendAPI_GetBoardAndThreads = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BoardAndThreadsResponse,
  /**
   * @param {!proto.feapi.BoardAndThreadsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardAndThreadsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BoardAndThreadsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BoardAndThreadsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BoardAndThreadsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.getBoardAndThreads =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/GetBoardAndThreads',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetBoardAndThreads,
      callback);
};


/**
 * @param {!proto.feapi.BoardAndThreadsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BoardAndThreadsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.getBoardAndThreads =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/GetBoardAndThreads',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetBoardAndThreads);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.AllBoardsRequest,
 *   !proto.feapi.AllBoardsResponse>}
 */
const methodDescriptor_FrontendAPI_GetAllBoards = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/GetAllBoards',
  grpc.web.MethodType.UNARY,
  proto.feapi.AllBoardsRequest,
  proto.feapi.AllBoardsResponse,
  /**
   * @param {!proto.feapi.AllBoardsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.AllBoardsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.AllBoardsRequest,
 *   !proto.feapi.AllBoardsResponse>}
 */
const methodInfo_FrontendAPI_GetAllBoards = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.AllBoardsResponse,
  /**
   * @param {!proto.feapi.AllBoardsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.AllBoardsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.AllBoardsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.AllBoardsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.AllBoardsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.getAllBoards =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/GetAllBoards',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetAllBoards,
      callback);
};


/**
 * @param {!proto.feapi.AllBoardsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.AllBoardsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.getAllBoards =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/GetAllBoards',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetAllBoards);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BoardSignalRequest,
 *   !proto.feapi.BoardSignalResponse>}
 */
const methodDescriptor_FrontendAPI_SetBoardSignal = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SetBoardSignal',
  grpc.web.MethodType.UNARY,
  proto.feapi.BoardSignalRequest,
  proto.feapi.BoardSignalResponse,
  /**
   * @param {!proto.feapi.BoardSignalRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardSignalResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BoardSignalRequest,
 *   !proto.feapi.BoardSignalResponse>}
 */
const methodInfo_FrontendAPI_SetBoardSignal = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BoardSignalResponse,
  /**
   * @param {!proto.feapi.BoardSignalRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardSignalResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BoardSignalRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BoardSignalResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BoardSignalResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.setBoardSignal =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SetBoardSignal',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetBoardSignal,
      callback);
};


/**
 * @param {!proto.feapi.BoardSignalRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BoardSignalResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.setBoardSignal =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SetBoardSignal',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetBoardSignal);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.UserAndGraphRequest,
 *   !proto.feapi.UserAndGraphResponse>}
 */
const methodDescriptor_FrontendAPI_GetUserAndGraph = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/GetUserAndGraph',
  grpc.web.MethodType.UNARY,
  proto.feapi.UserAndGraphRequest,
  proto.feapi.UserAndGraphResponse,
  /**
   * @param {!proto.feapi.UserAndGraphRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.UserAndGraphResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.UserAndGraphRequest,
 *   !proto.feapi.UserAndGraphResponse>}
 */
const methodInfo_FrontendAPI_GetUserAndGraph = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.UserAndGraphResponse,
  /**
   * @param {!proto.feapi.UserAndGraphRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.UserAndGraphResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.UserAndGraphRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.UserAndGraphResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.UserAndGraphResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.getUserAndGraph =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/GetUserAndGraph',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetUserAndGraph,
      callback);
};


/**
 * @param {!proto.feapi.UserAndGraphRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.UserAndGraphResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.getUserAndGraph =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/GetUserAndGraph',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetUserAndGraph);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.ContentEventPayload,
 *   !proto.feapi.ContentEventResponse>}
 */
const methodDescriptor_FrontendAPI_SendContentEvent = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendContentEvent',
  grpc.web.MethodType.UNARY,
  proto.feapi.ContentEventPayload,
  proto.feapi.ContentEventResponse,
  /**
   * @param {!proto.feapi.ContentEventPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ContentEventResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.ContentEventPayload,
 *   !proto.feapi.ContentEventResponse>}
 */
const methodInfo_FrontendAPI_SendContentEvent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.ContentEventResponse,
  /**
   * @param {!proto.feapi.ContentEventPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ContentEventResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.ContentEventPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.ContentEventResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.ContentEventResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendContentEvent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendContentEvent',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendContentEvent,
      callback);
};


/**
 * @param {!proto.feapi.ContentEventPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.ContentEventResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendContentEvent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendContentEvent',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendContentEvent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.SignalEventPayload,
 *   !proto.feapi.SignalEventResponse>}
 */
const methodDescriptor_FrontendAPI_SendSignalEvent = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendSignalEvent',
  grpc.web.MethodType.UNARY,
  proto.feapi.SignalEventPayload,
  proto.feapi.SignalEventResponse,
  /**
   * @param {!proto.feapi.SignalEventPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SignalEventResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.SignalEventPayload,
 *   !proto.feapi.SignalEventResponse>}
 */
const methodInfo_FrontendAPI_SendSignalEvent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.SignalEventResponse,
  /**
   * @param {!proto.feapi.SignalEventPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SignalEventResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.SignalEventPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.SignalEventResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.SignalEventResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendSignalEvent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendSignalEvent',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendSignalEvent,
      callback);
};


/**
 * @param {!proto.feapi.SignalEventPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.SignalEventResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendSignalEvent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendSignalEvent',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendSignalEvent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.UncompiledEntityByKeyRequest,
 *   !proto.feapi.UncompiledEntityByKeyResponse>}
 */
const methodDescriptor_FrontendAPI_GetUncompiledEntityByKey = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/GetUncompiledEntityByKey',
  grpc.web.MethodType.UNARY,
  proto.feapi.UncompiledEntityByKeyRequest,
  proto.feapi.UncompiledEntityByKeyResponse,
  /**
   * @param {!proto.feapi.UncompiledEntityByKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.UncompiledEntityByKeyResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.UncompiledEntityByKeyRequest,
 *   !proto.feapi.UncompiledEntityByKeyResponse>}
 */
const methodInfo_FrontendAPI_GetUncompiledEntityByKey = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.UncompiledEntityByKeyResponse,
  /**
   * @param {!proto.feapi.UncompiledEntityByKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.UncompiledEntityByKeyResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.UncompiledEntityByKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.UncompiledEntityByKeyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.UncompiledEntityByKeyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.getUncompiledEntityByKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/GetUncompiledEntityByKey',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetUncompiledEntityByKey,
      callback);
};


/**
 * @param {!proto.feapi.UncompiledEntityByKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.UncompiledEntityByKeyResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.getUncompiledEntityByKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/GetUncompiledEntityByKey',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_GetUncompiledEntityByKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.InflightsPruneRequest,
 *   !proto.feapi.InflightsPruneResponse>}
 */
const methodDescriptor_FrontendAPI_SendInflightsPruneRequest = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendInflightsPruneRequest',
  grpc.web.MethodType.UNARY,
  proto.feapi.InflightsPruneRequest,
  proto.feapi.InflightsPruneResponse,
  /**
   * @param {!proto.feapi.InflightsPruneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.InflightsPruneResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.InflightsPruneRequest,
 *   !proto.feapi.InflightsPruneResponse>}
 */
const methodInfo_FrontendAPI_SendInflightsPruneRequest = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.InflightsPruneResponse,
  /**
   * @param {!proto.feapi.InflightsPruneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.InflightsPruneResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.InflightsPruneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.InflightsPruneResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.InflightsPruneResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendInflightsPruneRequest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendInflightsPruneRequest',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendInflightsPruneRequest,
      callback);
};


/**
 * @param {!proto.feapi.InflightsPruneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.InflightsPruneResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendInflightsPruneRequest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendInflightsPruneRequest',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendInflightsPruneRequest);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.AmbientStatusRequest,
 *   !proto.feapi.AmbientStatusResponse>}
 */
const methodDescriptor_FrontendAPI_RequestAmbientStatus = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestAmbientStatus',
  grpc.web.MethodType.UNARY,
  proto.feapi.AmbientStatusRequest,
  proto.feapi.AmbientStatusResponse,
  /**
   * @param {!proto.feapi.AmbientStatusRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.AmbientStatusResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.AmbientStatusRequest,
 *   !proto.feapi.AmbientStatusResponse>}
 */
const methodInfo_FrontendAPI_RequestAmbientStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.AmbientStatusResponse,
  /**
   * @param {!proto.feapi.AmbientStatusRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.AmbientStatusResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.AmbientStatusRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.AmbientStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.AmbientStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestAmbientStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestAmbientStatus',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestAmbientStatus,
      callback);
};


/**
 * @param {!proto.feapi.AmbientStatusRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.AmbientStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestAmbientStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestAmbientStatus',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestAmbientStatus);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.HomeViewRequest,
 *   !proto.feapi.HomeViewResponse>}
 */
const methodDescriptor_FrontendAPI_RequestHomeView = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestHomeView',
  grpc.web.MethodType.UNARY,
  proto.feapi.HomeViewRequest,
  proto.feapi.HomeViewResponse,
  /**
   * @param {!proto.feapi.HomeViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.HomeViewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.HomeViewRequest,
 *   !proto.feapi.HomeViewResponse>}
 */
const methodInfo_FrontendAPI_RequestHomeView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.HomeViewResponse,
  /**
   * @param {!proto.feapi.HomeViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.HomeViewResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.HomeViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.HomeViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.HomeViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestHomeView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestHomeView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestHomeView,
      callback);
};


/**
 * @param {!proto.feapi.HomeViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.HomeViewResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestHomeView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestHomeView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestHomeView);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.PopularViewRequest,
 *   !proto.feapi.PopularViewResponse>}
 */
const methodDescriptor_FrontendAPI_RequestPopularView = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestPopularView',
  grpc.web.MethodType.UNARY,
  proto.feapi.PopularViewRequest,
  proto.feapi.PopularViewResponse,
  /**
   * @param {!proto.feapi.PopularViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.PopularViewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.PopularViewRequest,
 *   !proto.feapi.PopularViewResponse>}
 */
const methodInfo_FrontendAPI_RequestPopularView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.PopularViewResponse,
  /**
   * @param {!proto.feapi.PopularViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.PopularViewResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.PopularViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.PopularViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.PopularViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestPopularView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestPopularView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestPopularView,
      callback);
};


/**
 * @param {!proto.feapi.PopularViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.PopularViewResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestPopularView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestPopularView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestPopularView);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.NewViewRequest,
 *   !proto.feapi.NewViewResponse>}
 */
const methodDescriptor_FrontendAPI_RequestNewView = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestNewView',
  grpc.web.MethodType.UNARY,
  proto.feapi.NewViewRequest,
  proto.feapi.NewViewResponse,
  /**
   * @param {!proto.feapi.NewViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NewViewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.NewViewRequest,
 *   !proto.feapi.NewViewResponse>}
 */
const methodInfo_FrontendAPI_RequestNewView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.NewViewResponse,
  /**
   * @param {!proto.feapi.NewViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NewViewResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.NewViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.NewViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.NewViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestNewView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestNewView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestNewView,
      callback);
};


/**
 * @param {!proto.feapi.NewViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.NewViewResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestNewView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestNewView',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestNewView);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.NotificationsRequest,
 *   !proto.feapi.NotificationsResponse>}
 */
const methodDescriptor_FrontendAPI_RequestNotifications = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestNotifications',
  grpc.web.MethodType.UNARY,
  proto.feapi.NotificationsRequest,
  proto.feapi.NotificationsResponse,
  /**
   * @param {!proto.feapi.NotificationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NotificationsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.NotificationsRequest,
 *   !proto.feapi.NotificationsResponse>}
 */
const methodInfo_FrontendAPI_RequestNotifications = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.NotificationsResponse,
  /**
   * @param {!proto.feapi.NotificationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NotificationsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.NotificationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.NotificationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.NotificationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestNotifications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestNotifications',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestNotifications,
      callback);
};


/**
 * @param {!proto.feapi.NotificationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.NotificationsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestNotifications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestNotifications',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestNotifications);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.NotificationsSignalPayload,
 *   !proto.feapi.NotificationsSignalResponse>}
 */
const methodDescriptor_FrontendAPI_SetNotificationsSignal = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SetNotificationsSignal',
  grpc.web.MethodType.UNARY,
  proto.feapi.NotificationsSignalPayload,
  proto.feapi.NotificationsSignalResponse,
  /**
   * @param {!proto.feapi.NotificationsSignalPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NotificationsSignalResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.NotificationsSignalPayload,
 *   !proto.feapi.NotificationsSignalResponse>}
 */
const methodInfo_FrontendAPI_SetNotificationsSignal = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.NotificationsSignalResponse,
  /**
   * @param {!proto.feapi.NotificationsSignalPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.NotificationsSignalResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.NotificationsSignalPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.NotificationsSignalResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.NotificationsSignalResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.setNotificationsSignal =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SetNotificationsSignal',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetNotificationsSignal,
      callback);
};


/**
 * @param {!proto.feapi.NotificationsSignalPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.NotificationsSignalResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.setNotificationsSignal =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SetNotificationsSignal',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetNotificationsSignal);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.OnboardCompleteRequest,
 *   !proto.feapi.OnboardCompleteResponse>}
 */
const methodDescriptor_FrontendAPI_SetOnboardComplete = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SetOnboardComplete',
  grpc.web.MethodType.UNARY,
  proto.feapi.OnboardCompleteRequest,
  proto.feapi.OnboardCompleteResponse,
  /**
   * @param {!proto.feapi.OnboardCompleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.OnboardCompleteResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.OnboardCompleteRequest,
 *   !proto.feapi.OnboardCompleteResponse>}
 */
const methodInfo_FrontendAPI_SetOnboardComplete = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.OnboardCompleteResponse,
  /**
   * @param {!proto.feapi.OnboardCompleteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.OnboardCompleteResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.OnboardCompleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.OnboardCompleteResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.OnboardCompleteResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.setOnboardComplete =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SetOnboardComplete',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetOnboardComplete,
      callback);
};


/**
 * @param {!proto.feapi.OnboardCompleteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.OnboardCompleteResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.setOnboardComplete =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SetOnboardComplete',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SetOnboardComplete);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.SendAddressPayload,
 *   !proto.feapi.SendAddressResponse>}
 */
const methodDescriptor_FrontendAPI_SendAddress = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendAddress',
  grpc.web.MethodType.UNARY,
  proto.feapi.SendAddressPayload,
  proto.feapi.SendAddressResponse,
  /**
   * @param {!proto.feapi.SendAddressPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SendAddressResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.SendAddressPayload,
 *   !proto.feapi.SendAddressResponse>}
 */
const methodInfo_FrontendAPI_SendAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.SendAddressResponse,
  /**
   * @param {!proto.feapi.SendAddressPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SendAddressResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.SendAddressPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.SendAddressResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.SendAddressResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendAddress',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendAddress,
      callback);
};


/**
 * @param {!proto.feapi.SendAddressPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.SendAddressResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendAddress',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.FEConfigChangesPayload,
 *   !proto.feapi.FEConfigChangesResponse>}
 */
const methodDescriptor_FrontendAPI_SendFEConfigChanges = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendFEConfigChanges',
  grpc.web.MethodType.UNARY,
  proto.feapi.FEConfigChangesPayload,
  proto.feapi.FEConfigChangesResponse,
  /**
   * @param {!proto.feapi.FEConfigChangesPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.FEConfigChangesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.FEConfigChangesPayload,
 *   !proto.feapi.FEConfigChangesResponse>}
 */
const methodInfo_FrontendAPI_SendFEConfigChanges = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.FEConfigChangesResponse,
  /**
   * @param {!proto.feapi.FEConfigChangesPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.FEConfigChangesResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.FEConfigChangesPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.FEConfigChangesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.FEConfigChangesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendFEConfigChanges =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendFEConfigChanges',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendFEConfigChanges,
      callback);
};


/**
 * @param {!proto.feapi.FEConfigChangesPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.FEConfigChangesResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendFEConfigChanges =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendFEConfigChanges',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendFEConfigChanges);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BoardReportsRequest,
 *   !proto.feapi.BoardReportsResponse>}
 */
const methodDescriptor_FrontendAPI_RequestBoardReports = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestBoardReports',
  grpc.web.MethodType.UNARY,
  proto.feapi.BoardReportsRequest,
  proto.feapi.BoardReportsResponse,
  /**
   * @param {!proto.feapi.BoardReportsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardReportsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BoardReportsRequest,
 *   !proto.feapi.BoardReportsResponse>}
 */
const methodInfo_FrontendAPI_RequestBoardReports = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BoardReportsResponse,
  /**
   * @param {!proto.feapi.BoardReportsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardReportsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BoardReportsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BoardReportsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BoardReportsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestBoardReports =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestBoardReports',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestBoardReports,
      callback);
};


/**
 * @param {!proto.feapi.BoardReportsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BoardReportsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestBoardReports =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestBoardReports',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestBoardReports);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BoardModActionsRequest,
 *   !proto.feapi.BoardModActionsResponse>}
 */
const methodDescriptor_FrontendAPI_RequestBoardModActions = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/RequestBoardModActions',
  grpc.web.MethodType.UNARY,
  proto.feapi.BoardModActionsRequest,
  proto.feapi.BoardModActionsResponse,
  /**
   * @param {!proto.feapi.BoardModActionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardModActionsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BoardModActionsRequest,
 *   !proto.feapi.BoardModActionsResponse>}
 */
const methodInfo_FrontendAPI_RequestBoardModActions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BoardModActionsResponse,
  /**
   * @param {!proto.feapi.BoardModActionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BoardModActionsResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BoardModActionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BoardModActionsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BoardModActionsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.requestBoardModActions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestBoardModActions',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestBoardModActions,
      callback);
};


/**
 * @param {!proto.feapi.BoardModActionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BoardModActionsResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.requestBoardModActions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/RequestBoardModActions',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_RequestBoardModActions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.SendMintedUsernamesPayload,
 *   !proto.feapi.SendMintedUsernamesResponse>}
 */
const methodDescriptor_FrontendAPI_SendMintedUsernames = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendMintedUsernames',
  grpc.web.MethodType.UNARY,
  proto.feapi.SendMintedUsernamesPayload,
  proto.feapi.SendMintedUsernamesResponse,
  /**
   * @param {!proto.feapi.SendMintedUsernamesPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SendMintedUsernamesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.SendMintedUsernamesPayload,
 *   !proto.feapi.SendMintedUsernamesResponse>}
 */
const methodInfo_FrontendAPI_SendMintedUsernames = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.SendMintedUsernamesResponse,
  /**
   * @param {!proto.feapi.SendMintedUsernamesPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SendMintedUsernamesResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.SendMintedUsernamesPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.SendMintedUsernamesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.SendMintedUsernamesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendMintedUsernames =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendMintedUsernames',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendMintedUsernames,
      callback);
};


/**
 * @param {!proto.feapi.SendMintedUsernamesPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.SendMintedUsernamesResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendMintedUsernames =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendMintedUsernames',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendMintedUsernames);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.ClientVersionPayload,
 *   !proto.feapi.ClientVersionResponse>}
 */
const methodDescriptor_FrontendAPI_SendClientVersion = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendClientVersion',
  grpc.web.MethodType.UNARY,
  proto.feapi.ClientVersionPayload,
  proto.feapi.ClientVersionResponse,
  /**
   * @param {!proto.feapi.ClientVersionPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ClientVersionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.ClientVersionPayload,
 *   !proto.feapi.ClientVersionResponse>}
 */
const methodInfo_FrontendAPI_SendClientVersion = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.ClientVersionResponse,
  /**
   * @param {!proto.feapi.ClientVersionPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.ClientVersionResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.ClientVersionPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.ClientVersionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.ClientVersionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendClientVersion =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendClientVersion',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendClientVersion,
      callback);
};


/**
 * @param {!proto.feapi.ClientVersionPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.ClientVersionResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendClientVersion =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendClientVersion',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendClientVersion);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.SearchRequestPayload,
 *   !proto.feapi.SearchRequestResponse>}
 */
const methodDescriptor_FrontendAPI_SendSearchRequest = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendSearchRequest',
  grpc.web.MethodType.UNARY,
  proto.feapi.SearchRequestPayload,
  proto.feapi.SearchRequestResponse,
  /**
   * @param {!proto.feapi.SearchRequestPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SearchRequestResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.SearchRequestPayload,
 *   !proto.feapi.SearchRequestResponse>}
 */
const methodInfo_FrontendAPI_SendSearchRequest = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.SearchRequestResponse,
  /**
   * @param {!proto.feapi.SearchRequestPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.SearchRequestResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.SearchRequestPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.SearchRequestResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.SearchRequestResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendSearchRequest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendSearchRequest',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendSearchRequest,
      callback);
};


/**
 * @param {!proto.feapi.SearchRequestPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.SearchRequestResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendSearchRequest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendSearchRequest',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendSearchRequest);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BEReadyRequest,
 *   !proto.feapi.BEReadyResponse>}
 */
const methodDescriptor_FrontendAPI_BackendReady = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/BackendReady',
  grpc.web.MethodType.UNARY,
  proto.feapi.BEReadyRequest,
  proto.feapi.BEReadyResponse,
  /**
   * @param {!proto.feapi.BEReadyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BEReadyResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BEReadyRequest,
 *   !proto.feapi.BEReadyResponse>}
 */
const methodInfo_FrontendAPI_BackendReady = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BEReadyResponse,
  /**
   * @param {!proto.feapi.BEReadyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BEReadyResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BEReadyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BEReadyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BEReadyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.backendReady =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/BackendReady',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_BackendReady,
      callback);
};


/**
 * @param {!proto.feapi.BEReadyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BEReadyResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.backendReady =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/BackendReady',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_BackendReady);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.feapi.BackendAmbientStatusPayload,
 *   !proto.feapi.BackendAmbientStatusResponse>}
 */
const methodDescriptor_FrontendAPI_SendBackendAmbientStatus = new grpc.web.MethodDescriptor(
  '/feapi.FrontendAPI/SendBackendAmbientStatus',
  grpc.web.MethodType.UNARY,
  proto.feapi.BackendAmbientStatusPayload,
  proto.feapi.BackendAmbientStatusResponse,
  /**
   * @param {!proto.feapi.BackendAmbientStatusPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BackendAmbientStatusResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.feapi.BackendAmbientStatusPayload,
 *   !proto.feapi.BackendAmbientStatusResponse>}
 */
const methodInfo_FrontendAPI_SendBackendAmbientStatus = new grpc.web.AbstractClientBase.MethodInfo(
  proto.feapi.BackendAmbientStatusResponse,
  /**
   * @param {!proto.feapi.BackendAmbientStatusPayload} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.feapi.BackendAmbientStatusResponse.deserializeBinary
);


/**
 * @param {!proto.feapi.BackendAmbientStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.feapi.BackendAmbientStatusResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.feapi.BackendAmbientStatusResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.feapi.FrontendAPIClient.prototype.sendBackendAmbientStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/feapi.FrontendAPI/SendBackendAmbientStatus',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendBackendAmbientStatus,
      callback);
};


/**
 * @param {!proto.feapi.BackendAmbientStatusPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.feapi.BackendAmbientStatusResponse>}
 *     A native promise that resolves to the response
 */
proto.feapi.FrontendAPIPromiseClient.prototype.sendBackendAmbientStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/feapi.FrontendAPI/SendBackendAmbientStatus',
      request,
      metadata || {},
      methodDescriptor_FrontendAPI_SendBackendAmbientStatus);
};


module.exports = proto.feapi;

