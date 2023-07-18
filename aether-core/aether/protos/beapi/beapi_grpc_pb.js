// GENERATED CODE -- DO NOT EDIT!

// Original file comments:
// BackendAPI Protobufs
//
"use strict";
var grpc = require("grpc");
var beapi_beapi_pb = require("../beapi/beapi_pb.js");
var mimapi_mimapi_pb = require("../mimapi/mimapi_pb.js");

function serialize_beapi_AccessRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.AccessRequest)) {
    throw new Error("Expected argument of type beapi.AccessRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_AccessRequest(buffer_arg) {
  return beapi_beapi_pb.AccessRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_AccessResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.AccessResponse)) {
    throw new Error("Expected argument of type beapi.AccessResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_AccessResponse(buffer_arg) {
  return beapi_beapi_pb.AccessResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_BoardThreadsCountRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.BoardThreadsCountRequest)) {
    throw new Error("Expected argument of type beapi.BoardThreadsCountRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_BoardThreadsCountRequest(buffer_arg) {
  return beapi_beapi_pb.BoardThreadsCountRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_BoardThreadsCountResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.BoardThreadsCountResponse)) {
    throw new Error(
      "Expected argument of type beapi.BoardThreadsCountResponse",
    );
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_BoardThreadsCountResponse(buffer_arg) {
  return beapi_beapi_pb.BoardThreadsCountResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_BoardsRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.BoardsRequest)) {
    throw new Error("Expected argument of type beapi.BoardsRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_BoardsRequest(buffer_arg) {
  return beapi_beapi_pb.BoardsRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_BoardsResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.BoardsResponse)) {
    throw new Error("Expected argument of type beapi.BoardsResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_BoardsResponse(buffer_arg) {
  return beapi_beapi_pb.BoardsResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ConnectToRemoteRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.ConnectToRemoteRequest)) {
    throw new Error("Expected argument of type beapi.ConnectToRemoteRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ConnectToRemoteRequest(buffer_arg) {
  return beapi_beapi_pb.ConnectToRemoteRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ConnectToRemoteResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.ConnectToRemoteResponse)) {
    throw new Error("Expected argument of type beapi.ConnectToRemoteResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ConnectToRemoteResponse(buffer_arg) {
  return beapi_beapi_pb.ConnectToRemoteResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_KeysRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.KeysRequest)) {
    throw new Error("Expected argument of type beapi.KeysRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_KeysRequest(buffer_arg) {
  return beapi_beapi_pb.KeysRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_KeysResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.KeysResponse)) {
    throw new Error("Expected argument of type beapi.KeysResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_KeysResponse(buffer_arg) {
  return beapi_beapi_pb.KeysResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_MintedContentPayload(arg) {
  if (!(arg instanceof beapi_beapi_pb.MintedContentPayload)) {
    throw new Error("Expected argument of type beapi.MintedContentPayload");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_MintedContentPayload(buffer_arg) {
  return beapi_beapi_pb.MintedContentPayload.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_MintedContentResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.MintedContentResponse)) {
    throw new Error("Expected argument of type beapi.MintedContentResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_MintedContentResponse(buffer_arg) {
  return beapi_beapi_pb.MintedContentResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_PostsRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.PostsRequest)) {
    throw new Error("Expected argument of type beapi.PostsRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_PostsRequest(buffer_arg) {
  return beapi_beapi_pb.PostsRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_PostsResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.PostsResponse)) {
    throw new Error("Expected argument of type beapi.PostsResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_PostsResponse(buffer_arg) {
  return beapi_beapi_pb.PostsResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ThreadPostsCountRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.ThreadPostsCountRequest)) {
    throw new Error("Expected argument of type beapi.ThreadPostsCountRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ThreadPostsCountRequest(buffer_arg) {
  return beapi_beapi_pb.ThreadPostsCountRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ThreadPostsCountResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.ThreadPostsCountResponse)) {
    throw new Error("Expected argument of type beapi.ThreadPostsCountResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ThreadPostsCountResponse(buffer_arg) {
  return beapi_beapi_pb.ThreadPostsCountResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ThreadsRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.ThreadsRequest)) {
    throw new Error("Expected argument of type beapi.ThreadsRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ThreadsRequest(buffer_arg) {
  return beapi_beapi_pb.ThreadsRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_ThreadsResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.ThreadsResponse)) {
    throw new Error("Expected argument of type beapi.ThreadsResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_ThreadsResponse(buffer_arg) {
  return beapi_beapi_pb.ThreadsResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_TruststatesRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.TruststatesRequest)) {
    throw new Error("Expected argument of type beapi.TruststatesRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_TruststatesRequest(buffer_arg) {
  return beapi_beapi_pb.TruststatesRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_TruststatesResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.TruststatesResponse)) {
    throw new Error("Expected argument of type beapi.TruststatesResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_TruststatesResponse(buffer_arg) {
  return beapi_beapi_pb.TruststatesResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_VotesRequest(arg) {
  if (!(arg instanceof beapi_beapi_pb.VotesRequest)) {
    throw new Error("Expected argument of type beapi.VotesRequest");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_VotesRequest(buffer_arg) {
  return beapi_beapi_pb.VotesRequest.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

function serialize_beapi_VotesResponse(arg) {
  if (!(arg instanceof beapi_beapi_pb.VotesResponse)) {
    throw new Error("Expected argument of type beapi.VotesResponse");
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_beapi_VotesResponse(buffer_arg) {
  return beapi_beapi_pb.VotesResponse.deserializeBinary(
    new Uint8Array(buffer_arg),
  );
}

var BackendAPIService = (exports.BackendAPIService = {
  // You need to request access first, before anything.
  requestBackendAccess: {
    path: "/beapi.BackendAPI/RequestBackendAccess",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.AccessRequest,
    responseType: beapi_beapi_pb.AccessResponse,
    requestSerialize: serialize_beapi_AccessRequest,
    requestDeserialize: deserialize_beapi_AccessRequest,
    responseSerialize: serialize_beapi_AccessResponse,
    responseDeserialize: deserialize_beapi_AccessResponse,
  },
  //
  // Lower level APIs to get uncompiled Mim objects in the case that there is no specific higher level API is available to you. Look at higher level APIs requesting more specific things first.
  //
  // For example, an uncompiled Mim object is a Truststate, which can mean any of those entities: a TS of publictrust, naming, f451, or elector class types. If you get the truststate itself, you'll have to figure this out on the frontend, which isn't very efficient (but the API is there if you have no other way.)
  //
  // Ideally though, if you end up using these often, it might be time to build a specific higher level API for it.
  getBoards: {
    path: "/beapi.BackendAPI/GetBoards",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.BoardsRequest,
    responseType: beapi_beapi_pb.BoardsResponse,
    requestSerialize: serialize_beapi_BoardsRequest,
    requestDeserialize: deserialize_beapi_BoardsRequest,
    responseSerialize: serialize_beapi_BoardsResponse,
    responseDeserialize: deserialize_beapi_BoardsResponse,
  },
  getThreads: {
    path: "/beapi.BackendAPI/GetThreads",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.ThreadsRequest,
    responseType: beapi_beapi_pb.ThreadsResponse,
    requestSerialize: serialize_beapi_ThreadsRequest,
    requestDeserialize: deserialize_beapi_ThreadsRequest,
    responseSerialize: serialize_beapi_ThreadsResponse,
    responseDeserialize: deserialize_beapi_ThreadsResponse,
  },
  getPosts: {
    path: "/beapi.BackendAPI/GetPosts",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.PostsRequest,
    responseType: beapi_beapi_pb.PostsResponse,
    requestSerialize: serialize_beapi_PostsRequest,
    requestDeserialize: deserialize_beapi_PostsRequest,
    responseSerialize: serialize_beapi_PostsResponse,
    responseDeserialize: deserialize_beapi_PostsResponse,
  },
  getVotes: {
    path: "/beapi.BackendAPI/GetVotes",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.VotesRequest,
    responseType: beapi_beapi_pb.VotesResponse,
    requestSerialize: serialize_beapi_VotesRequest,
    requestDeserialize: deserialize_beapi_VotesRequest,
    responseSerialize: serialize_beapi_VotesResponse,
    responseDeserialize: deserialize_beapi_VotesResponse,
  },
  getKeys: {
    path: "/beapi.BackendAPI/GetKeys",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.KeysRequest,
    responseType: beapi_beapi_pb.KeysResponse,
    requestSerialize: serialize_beapi_KeysRequest,
    requestDeserialize: deserialize_beapi_KeysRequest,
    responseSerialize: serialize_beapi_KeysResponse,
    responseDeserialize: deserialize_beapi_KeysResponse,
  },
  getTruststates: {
    path: "/beapi.BackendAPI/GetTruststates",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.TruststatesRequest,
    responseType: beapi_beapi_pb.TruststatesResponse,
    requestSerialize: serialize_beapi_TruststatesRequest,
    requestDeserialize: deserialize_beapi_TruststatesRequest,
    responseSerialize: serialize_beapi_TruststatesResponse,
    responseDeserialize: deserialize_beapi_TruststatesResponse,
  },
  getBoardThreadsCount: {
    path: "/beapi.BackendAPI/GetBoardThreadsCount",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.BoardThreadsCountRequest,
    responseType: beapi_beapi_pb.BoardThreadsCountResponse,
    requestSerialize: serialize_beapi_BoardThreadsCountRequest,
    requestDeserialize: deserialize_beapi_BoardThreadsCountRequest,
    responseSerialize: serialize_beapi_BoardThreadsCountResponse,
    responseDeserialize: deserialize_beapi_BoardThreadsCountResponse,
  },
  getThreadPostsCount: {
    path: "/beapi.BackendAPI/GetThreadPostsCount",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.ThreadPostsCountRequest,
    responseType: beapi_beapi_pb.ThreadPostsCountResponse,
    requestSerialize: serialize_beapi_ThreadPostsCountRequest,
    requestDeserialize: deserialize_beapi_ThreadPostsCountRequest,
    responseSerialize: serialize_beapi_ThreadPostsCountResponse,
    responseDeserialize: deserialize_beapi_ThreadPostsCountResponse,
  },
  sendMintedContent: {
    path: "/beapi.BackendAPI/SendMintedContent",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.MintedContentPayload,
    responseType: beapi_beapi_pb.MintedContentResponse,
    requestSerialize: serialize_beapi_MintedContentPayload,
    requestDeserialize: deserialize_beapi_MintedContentPayload,
    responseSerialize: serialize_beapi_MintedContentResponse,
    responseDeserialize: deserialize_beapi_MintedContentResponse,
  },
  sendConnectToRemoteRequest: {
    path: "/beapi.BackendAPI/SendConnectToRemoteRequest",
    requestStream: false,
    responseStream: false,
    requestType: beapi_beapi_pb.ConnectToRemoteRequest,
    responseType: beapi_beapi_pb.ConnectToRemoteResponse,
    requestSerialize: serialize_beapi_ConnectToRemoteRequest,
    requestDeserialize: deserialize_beapi_ConnectToRemoteRequest,
    responseSerialize: serialize_beapi_ConnectToRemoteResponse,
    responseDeserialize: deserialize_beapi_ConnectToRemoteResponse,
  },
});

exports.BackendAPIClient = grpc.makeGenericClientConstructor(BackendAPIService);
