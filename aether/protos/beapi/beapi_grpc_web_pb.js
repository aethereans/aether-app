/**
 * @fileoverview gRPC-Web generated client stub for beapi
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var mimapi_structprotos_pb = require('../mimapi/structprotos_pb.js')
const proto = {};
proto.beapi = require('./beapi_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.beapi.BackendAPIClient =
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
proto.beapi.BackendAPIPromiseClient =
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
 *   !proto.beapi.AccessRequest,
 *   !proto.beapi.AccessResponse>}
 */
const methodInfo_BackendAPI_RequestBackendAccess = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.AccessResponse,
  /** @param {!proto.beapi.AccessRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.AccessResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.AccessRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.AccessResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.AccessResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.requestBackendAccess =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/RequestBackendAccess',
      request,
      metadata || {},
      methodInfo_BackendAPI_RequestBackendAccess,
      callback);
};


/**
 * @param {!proto.beapi.AccessRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.AccessResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.requestBackendAccess =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/RequestBackendAccess',
      request,
      metadata || {},
      methodInfo_BackendAPI_RequestBackendAccess);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.BoardsRequest,
 *   !proto.beapi.BoardsResponse>}
 */
const methodInfo_BackendAPI_GetBoards = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.BoardsResponse,
  /** @param {!proto.beapi.BoardsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.BoardsResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.BoardsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.BoardsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.BoardsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getBoards =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetBoards',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetBoards,
      callback);
};


/**
 * @param {!proto.beapi.BoardsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.BoardsResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getBoards =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetBoards',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetBoards);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.ThreadsRequest,
 *   !proto.beapi.ThreadsResponse>}
 */
const methodInfo_BackendAPI_GetThreads = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.ThreadsResponse,
  /** @param {!proto.beapi.ThreadsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.ThreadsResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.ThreadsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.ThreadsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.ThreadsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getThreads =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetThreads',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetThreads,
      callback);
};


/**
 * @param {!proto.beapi.ThreadsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.ThreadsResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getThreads =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetThreads',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetThreads);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.PostsRequest,
 *   !proto.beapi.PostsResponse>}
 */
const methodInfo_BackendAPI_GetPosts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.PostsResponse,
  /** @param {!proto.beapi.PostsRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.PostsResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.PostsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.PostsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.PostsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getPosts =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetPosts',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetPosts,
      callback);
};


/**
 * @param {!proto.beapi.PostsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.PostsResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getPosts =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetPosts',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetPosts);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.VotesRequest,
 *   !proto.beapi.VotesResponse>}
 */
const methodInfo_BackendAPI_GetVotes = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.VotesResponse,
  /** @param {!proto.beapi.VotesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.VotesResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.VotesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.VotesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.VotesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getVotes =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetVotes',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetVotes,
      callback);
};


/**
 * @param {!proto.beapi.VotesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.VotesResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getVotes =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetVotes',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetVotes);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.KeysRequest,
 *   !proto.beapi.KeysResponse>}
 */
const methodInfo_BackendAPI_GetKeys = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.KeysResponse,
  /** @param {!proto.beapi.KeysRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.KeysResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.KeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.KeysResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.KeysResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getKeys =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetKeys',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetKeys,
      callback);
};


/**
 * @param {!proto.beapi.KeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.KeysResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getKeys =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetKeys',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetKeys);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.TruststatesRequest,
 *   !proto.beapi.TruststatesResponse>}
 */
const methodInfo_BackendAPI_GetTruststates = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.TruststatesResponse,
  /** @param {!proto.beapi.TruststatesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.TruststatesResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.TruststatesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.TruststatesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.TruststatesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getTruststates =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetTruststates',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetTruststates,
      callback);
};


/**
 * @param {!proto.beapi.TruststatesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.TruststatesResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getTruststates =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetTruststates',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetTruststates);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.BoardThreadsCountRequest,
 *   !proto.beapi.BoardThreadsCountResponse>}
 */
const methodInfo_BackendAPI_GetBoardThreadsCount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.BoardThreadsCountResponse,
  /** @param {!proto.beapi.BoardThreadsCountRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.BoardThreadsCountResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.BoardThreadsCountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.BoardThreadsCountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.BoardThreadsCountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getBoardThreadsCount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetBoardThreadsCount',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetBoardThreadsCount,
      callback);
};


/**
 * @param {!proto.beapi.BoardThreadsCountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.BoardThreadsCountResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getBoardThreadsCount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetBoardThreadsCount',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetBoardThreadsCount);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.ThreadPostsCountRequest,
 *   !proto.beapi.ThreadPostsCountResponse>}
 */
const methodInfo_BackendAPI_GetThreadPostsCount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.ThreadPostsCountResponse,
  /** @param {!proto.beapi.ThreadPostsCountRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.ThreadPostsCountResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.ThreadPostsCountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.ThreadPostsCountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.ThreadPostsCountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.getThreadPostsCount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/GetThreadPostsCount',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetThreadPostsCount,
      callback);
};


/**
 * @param {!proto.beapi.ThreadPostsCountRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.ThreadPostsCountResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.getThreadPostsCount =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/GetThreadPostsCount',
      request,
      metadata || {},
      methodInfo_BackendAPI_GetThreadPostsCount);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.MintedContentPayload,
 *   !proto.beapi.MintedContentResponse>}
 */
const methodInfo_BackendAPI_SendMintedContent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.MintedContentResponse,
  /** @param {!proto.beapi.MintedContentPayload} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.MintedContentResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.MintedContentPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.MintedContentResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.MintedContentResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.sendMintedContent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/SendMintedContent',
      request,
      metadata || {},
      methodInfo_BackendAPI_SendMintedContent,
      callback);
};


/**
 * @param {!proto.beapi.MintedContentPayload} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.MintedContentResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.sendMintedContent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/SendMintedContent',
      request,
      metadata || {},
      methodInfo_BackendAPI_SendMintedContent);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.beapi.ConnectToRemoteRequest,
 *   !proto.beapi.ConnectToRemoteResponse>}
 */
const methodInfo_BackendAPI_SendConnectToRemoteRequest = new grpc.web.AbstractClientBase.MethodInfo(
  proto.beapi.ConnectToRemoteResponse,
  /** @param {!proto.beapi.ConnectToRemoteRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.beapi.ConnectToRemoteResponse.deserializeBinary
);


/**
 * @param {!proto.beapi.ConnectToRemoteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.beapi.ConnectToRemoteResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.beapi.ConnectToRemoteResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.beapi.BackendAPIClient.prototype.sendConnectToRemoteRequest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/beapi.BackendAPI/SendConnectToRemoteRequest',
      request,
      metadata || {},
      methodInfo_BackendAPI_SendConnectToRemoteRequest,
      callback);
};


/**
 * @param {!proto.beapi.ConnectToRemoteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.beapi.ConnectToRemoteResponse>}
 *     A native promise that resolves to the response
 */
proto.beapi.BackendAPIPromiseClient.prototype.sendConnectToRemoteRequest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/beapi.BackendAPI/SendConnectToRemoteRequest',
      request,
      metadata || {},
      methodInfo_BackendAPI_SendConnectToRemoteRequest);
};


module.exports = proto.beapi;

