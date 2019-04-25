// GENERATED CODE -- DO NOT EDIT!

// Original file comments:
// Frontend API server Protobufs
//
'use strict';
var grpc = require('grpc');
var feapi_feapi_pb = require('../feapi/feapi_pb.js');
var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js');
var mimapi_structprotos_pb = require('../mimapi/structprotos_pb.js');

function serialize_feapi_AllBoardsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.AllBoardsRequest)) {
    throw new Error('Expected argument of type feapi.AllBoardsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AllBoardsRequest(buffer_arg) {
  return feapi_feapi_pb.AllBoardsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_AllBoardsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.AllBoardsResponse)) {
    throw new Error('Expected argument of type feapi.AllBoardsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AllBoardsResponse(buffer_arg) {
  return feapi_feapi_pb.AllBoardsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_AmbientStatusRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.AmbientStatusRequest)) {
    throw new Error('Expected argument of type feapi.AmbientStatusRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AmbientStatusRequest(buffer_arg) {
  return feapi_feapi_pb.AmbientStatusRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_AmbientStatusResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.AmbientStatusResponse)) {
    throw new Error('Expected argument of type feapi.AmbientStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_AmbientStatusResponse(buffer_arg) {
  return feapi_feapi_pb.AmbientStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BEReadyRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BEReadyRequest)) {
    throw new Error('Expected argument of type feapi.BEReadyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BEReadyRequest(buffer_arg) {
  return feapi_feapi_pb.BEReadyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BEReadyResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BEReadyResponse)) {
    throw new Error('Expected argument of type feapi.BEReadyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BEReadyResponse(buffer_arg) {
  return feapi_feapi_pb.BEReadyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BackendAmbientStatusPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.BackendAmbientStatusPayload)) {
    throw new Error('Expected argument of type feapi.BackendAmbientStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BackendAmbientStatusPayload(buffer_arg) {
  return feapi_feapi_pb.BackendAmbientStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BackendAmbientStatusResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BackendAmbientStatusResponse)) {
    throw new Error('Expected argument of type feapi.BackendAmbientStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BackendAmbientStatusResponse(buffer_arg) {
  return feapi_feapi_pb.BackendAmbientStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardAndThreadsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardAndThreadsRequest)) {
    throw new Error('Expected argument of type feapi.BoardAndThreadsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardAndThreadsRequest(buffer_arg) {
  return feapi_feapi_pb.BoardAndThreadsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardAndThreadsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardAndThreadsResponse)) {
    throw new Error('Expected argument of type feapi.BoardAndThreadsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardAndThreadsResponse(buffer_arg) {
  return feapi_feapi_pb.BoardAndThreadsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardReportsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardReportsRequest)) {
    throw new Error('Expected argument of type feapi.BoardReportsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardReportsRequest(buffer_arg) {
  return feapi_feapi_pb.BoardReportsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardReportsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardReportsResponse)) {
    throw new Error('Expected argument of type feapi.BoardReportsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardReportsResponse(buffer_arg) {
  return feapi_feapi_pb.BoardReportsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardSignalRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardSignalRequest)) {
    throw new Error('Expected argument of type feapi.BoardSignalRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardSignalRequest(buffer_arg) {
  return feapi_feapi_pb.BoardSignalRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_BoardSignalResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.BoardSignalResponse)) {
    throw new Error('Expected argument of type feapi.BoardSignalResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_BoardSignalResponse(buffer_arg) {
  return feapi_feapi_pb.BoardSignalResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ClientVersionPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.ClientVersionPayload)) {
    throw new Error('Expected argument of type feapi.ClientVersionPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ClientVersionPayload(buffer_arg) {
  return feapi_feapi_pb.ClientVersionPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ClientVersionResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.ClientVersionResponse)) {
    throw new Error('Expected argument of type feapi.ClientVersionResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ClientVersionResponse(buffer_arg) {
  return feapi_feapi_pb.ClientVersionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ContentEventPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.ContentEventPayload)) {
    throw new Error('Expected argument of type feapi.ContentEventPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ContentEventPayload(buffer_arg) {
  return feapi_feapi_pb.ContentEventPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ContentEventResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.ContentEventResponse)) {
    throw new Error('Expected argument of type feapi.ContentEventResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ContentEventResponse(buffer_arg) {
  return feapi_feapi_pb.ContentEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_FEConfigChangesPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.FEConfigChangesPayload)) {
    throw new Error('Expected argument of type feapi.FEConfigChangesPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_FEConfigChangesPayload(buffer_arg) {
  return feapi_feapi_pb.FEConfigChangesPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_FEConfigChangesResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.FEConfigChangesResponse)) {
    throw new Error('Expected argument of type feapi.FEConfigChangesResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_FEConfigChangesResponse(buffer_arg) {
  return feapi_feapi_pb.FEConfigChangesResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_HomeViewRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.HomeViewRequest)) {
    throw new Error('Expected argument of type feapi.HomeViewRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_HomeViewRequest(buffer_arg) {
  return feapi_feapi_pb.HomeViewRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_HomeViewResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.HomeViewResponse)) {
    throw new Error('Expected argument of type feapi.HomeViewResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_HomeViewResponse(buffer_arg) {
  return feapi_feapi_pb.HomeViewResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_InflightsPruneRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.InflightsPruneRequest)) {
    throw new Error('Expected argument of type feapi.InflightsPruneRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_InflightsPruneRequest(buffer_arg) {
  return feapi_feapi_pb.InflightsPruneRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_InflightsPruneResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.InflightsPruneResponse)) {
    throw new Error('Expected argument of type feapi.InflightsPruneResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_InflightsPruneResponse(buffer_arg) {
  return feapi_feapi_pb.InflightsPruneResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_NotificationsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.NotificationsRequest)) {
    throw new Error('Expected argument of type feapi.NotificationsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_NotificationsRequest(buffer_arg) {
  return feapi_feapi_pb.NotificationsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_NotificationsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.NotificationsResponse)) {
    throw new Error('Expected argument of type feapi.NotificationsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_NotificationsResponse(buffer_arg) {
  return feapi_feapi_pb.NotificationsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_NotificationsSignalPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.NotificationsSignalPayload)) {
    throw new Error('Expected argument of type feapi.NotificationsSignalPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_NotificationsSignalPayload(buffer_arg) {
  return feapi_feapi_pb.NotificationsSignalPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_NotificationsSignalResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.NotificationsSignalResponse)) {
    throw new Error('Expected argument of type feapi.NotificationsSignalResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_NotificationsSignalResponse(buffer_arg) {
  return feapi_feapi_pb.NotificationsSignalResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_OnboardCompleteRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.OnboardCompleteRequest)) {
    throw new Error('Expected argument of type feapi.OnboardCompleteRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_OnboardCompleteRequest(buffer_arg) {
  return feapi_feapi_pb.OnboardCompleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_OnboardCompleteResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.OnboardCompleteResponse)) {
    throw new Error('Expected argument of type feapi.OnboardCompleteResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_OnboardCompleteResponse(buffer_arg) {
  return feapi_feapi_pb.OnboardCompleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_PopularViewRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.PopularViewRequest)) {
    throw new Error('Expected argument of type feapi.PopularViewRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_PopularViewRequest(buffer_arg) {
  return feapi_feapi_pb.PopularViewRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_PopularViewResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.PopularViewResponse)) {
    throw new Error('Expected argument of type feapi.PopularViewResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_PopularViewResponse(buffer_arg) {
  return feapi_feapi_pb.PopularViewResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SearchRequestPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.SearchRequestPayload)) {
    throw new Error('Expected argument of type feapi.SearchRequestPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SearchRequestPayload(buffer_arg) {
  return feapi_feapi_pb.SearchRequestPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SearchRequestResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SearchRequestResponse)) {
    throw new Error('Expected argument of type feapi.SearchRequestResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SearchRequestResponse(buffer_arg) {
  return feapi_feapi_pb.SearchRequestResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SendAddressPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.SendAddressPayload)) {
    throw new Error('Expected argument of type feapi.SendAddressPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SendAddressPayload(buffer_arg) {
  return feapi_feapi_pb.SendAddressPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SendAddressResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SendAddressResponse)) {
    throw new Error('Expected argument of type feapi.SendAddressResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SendAddressResponse(buffer_arg) {
  return feapi_feapi_pb.SendAddressResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SendMintedUsernamesPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.SendMintedUsernamesPayload)) {
    throw new Error('Expected argument of type feapi.SendMintedUsernamesPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SendMintedUsernamesPayload(buffer_arg) {
  return feapi_feapi_pb.SendMintedUsernamesPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SendMintedUsernamesResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SendMintedUsernamesResponse)) {
    throw new Error('Expected argument of type feapi.SendMintedUsernamesResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SendMintedUsernamesResponse(buffer_arg) {
  return feapi_feapi_pb.SendMintedUsernamesResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SetClientAPIServerPortRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.SetClientAPIServerPortRequest)) {
    throw new Error('Expected argument of type feapi.SetClientAPIServerPortRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SetClientAPIServerPortRequest(buffer_arg) {
  return feapi_feapi_pb.SetClientAPIServerPortRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SetClientAPIServerPortResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SetClientAPIServerPortResponse)) {
    throw new Error('Expected argument of type feapi.SetClientAPIServerPortResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SetClientAPIServerPortResponse(buffer_arg) {
  return feapi_feapi_pb.SetClientAPIServerPortResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SignalEventPayload(arg) {
  if (!(arg instanceof feapi_feapi_pb.SignalEventPayload)) {
    throw new Error('Expected argument of type feapi.SignalEventPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SignalEventPayload(buffer_arg) {
  return feapi_feapi_pb.SignalEventPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_SignalEventResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.SignalEventResponse)) {
    throw new Error('Expected argument of type feapi.SignalEventResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_SignalEventResponse(buffer_arg) {
  return feapi_feapi_pb.SignalEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ThreadAndPostsRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.ThreadAndPostsRequest)) {
    throw new Error('Expected argument of type feapi.ThreadAndPostsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ThreadAndPostsRequest(buffer_arg) {
  return feapi_feapi_pb.ThreadAndPostsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_ThreadAndPostsResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.ThreadAndPostsResponse)) {
    throw new Error('Expected argument of type feapi.ThreadAndPostsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_ThreadAndPostsResponse(buffer_arg) {
  return feapi_feapi_pb.ThreadAndPostsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UncompiledEntityByKeyRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.UncompiledEntityByKeyRequest)) {
    throw new Error('Expected argument of type feapi.UncompiledEntityByKeyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UncompiledEntityByKeyRequest(buffer_arg) {
  return feapi_feapi_pb.UncompiledEntityByKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UncompiledEntityByKeyResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.UncompiledEntityByKeyResponse)) {
    throw new Error('Expected argument of type feapi.UncompiledEntityByKeyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UncompiledEntityByKeyResponse(buffer_arg) {
  return feapi_feapi_pb.UncompiledEntityByKeyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UserAndGraphRequest(arg) {
  if (!(arg instanceof feapi_feapi_pb.UserAndGraphRequest)) {
    throw new Error('Expected argument of type feapi.UserAndGraphRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UserAndGraphRequest(buffer_arg) {
  return feapi_feapi_pb.UserAndGraphRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_feapi_UserAndGraphResponse(arg) {
  if (!(arg instanceof feapi_feapi_pb.UserAndGraphResponse)) {
    throw new Error('Expected argument of type feapi.UserAndGraphResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_feapi_UserAndGraphResponse(buffer_arg) {
  return feapi_feapi_pb.UserAndGraphResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// These "Set", "Get" verbs are written from the viewpoint of the consumer of this api.
var FrontendAPIService = exports.FrontendAPIService = {
  // ----------  Methods used by client  ----------
  setClientAPIServerPort: {
    path: '/feapi.FrontendAPI/SetClientAPIServerPort',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SetClientAPIServerPortRequest,
    responseType: feapi_feapi_pb.SetClientAPIServerPortResponse,
    requestSerialize: serialize_feapi_SetClientAPIServerPortRequest,
    requestDeserialize: deserialize_feapi_SetClientAPIServerPortRequest,
    responseSerialize: serialize_feapi_SetClientAPIServerPortResponse,
    responseDeserialize: deserialize_feapi_SetClientAPIServerPortResponse,
  },
  getThreadAndPosts: {
    path: '/feapi.FrontendAPI/GetThreadAndPosts',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.ThreadAndPostsRequest,
    responseType: feapi_feapi_pb.ThreadAndPostsResponse,
    requestSerialize: serialize_feapi_ThreadAndPostsRequest,
    requestDeserialize: deserialize_feapi_ThreadAndPostsRequest,
    responseSerialize: serialize_feapi_ThreadAndPostsResponse,
    responseDeserialize: deserialize_feapi_ThreadAndPostsResponse,
  },
  getBoardAndThreads: {
    path: '/feapi.FrontendAPI/GetBoardAndThreads',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BoardAndThreadsRequest,
    responseType: feapi_feapi_pb.BoardAndThreadsResponse,
    requestSerialize: serialize_feapi_BoardAndThreadsRequest,
    requestDeserialize: deserialize_feapi_BoardAndThreadsRequest,
    responseSerialize: serialize_feapi_BoardAndThreadsResponse,
    responseDeserialize: deserialize_feapi_BoardAndThreadsResponse,
  },
  getAllBoards: {
    path: '/feapi.FrontendAPI/GetAllBoards',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.AllBoardsRequest,
    responseType: feapi_feapi_pb.AllBoardsResponse,
    requestSerialize: serialize_feapi_AllBoardsRequest,
    requestDeserialize: deserialize_feapi_AllBoardsRequest,
    responseSerialize: serialize_feapi_AllBoardsResponse,
    responseDeserialize: deserialize_feapi_AllBoardsResponse,
  },
  setBoardSignal: {
    path: '/feapi.FrontendAPI/SetBoardSignal',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BoardSignalRequest,
    responseType: feapi_feapi_pb.BoardSignalResponse,
    requestSerialize: serialize_feapi_BoardSignalRequest,
    requestDeserialize: deserialize_feapi_BoardSignalRequest,
    responseSerialize: serialize_feapi_BoardSignalResponse,
    responseDeserialize: deserialize_feapi_BoardSignalResponse,
  },
  getUserAndGraph: {
    path: '/feapi.FrontendAPI/GetUserAndGraph',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.UserAndGraphRequest,
    responseType: feapi_feapi_pb.UserAndGraphResponse,
    requestSerialize: serialize_feapi_UserAndGraphRequest,
    requestDeserialize: deserialize_feapi_UserAndGraphRequest,
    responseSerialize: serialize_feapi_UserAndGraphResponse,
    responseDeserialize: deserialize_feapi_UserAndGraphResponse,
  },
  sendContentEvent: {
    path: '/feapi.FrontendAPI/SendContentEvent',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.ContentEventPayload,
    responseType: feapi_feapi_pb.ContentEventResponse,
    requestSerialize: serialize_feapi_ContentEventPayload,
    requestDeserialize: deserialize_feapi_ContentEventPayload,
    responseSerialize: serialize_feapi_ContentEventResponse,
    responseDeserialize: deserialize_feapi_ContentEventResponse,
  },
  sendSignalEvent: {
    path: '/feapi.FrontendAPI/SendSignalEvent',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SignalEventPayload,
    responseType: feapi_feapi_pb.SignalEventResponse,
    requestSerialize: serialize_feapi_SignalEventPayload,
    requestDeserialize: deserialize_feapi_SignalEventPayload,
    responseSerialize: serialize_feapi_SignalEventResponse,
    responseDeserialize: deserialize_feapi_SignalEventResponse,
  },
  getUncompiledEntityByKey: {
    path: '/feapi.FrontendAPI/GetUncompiledEntityByKey',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.UncompiledEntityByKeyRequest,
    responseType: feapi_feapi_pb.UncompiledEntityByKeyResponse,
    requestSerialize: serialize_feapi_UncompiledEntityByKeyRequest,
    requestDeserialize: deserialize_feapi_UncompiledEntityByKeyRequest,
    responseSerialize: serialize_feapi_UncompiledEntityByKeyResponse,
    responseDeserialize: deserialize_feapi_UncompiledEntityByKeyResponse,
  },
  sendInflightsPruneRequest: {
    path: '/feapi.FrontendAPI/SendInflightsPruneRequest',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.InflightsPruneRequest,
    responseType: feapi_feapi_pb.InflightsPruneResponse,
    requestSerialize: serialize_feapi_InflightsPruneRequest,
    requestDeserialize: deserialize_feapi_InflightsPruneRequest,
    responseSerialize: serialize_feapi_InflightsPruneResponse,
    responseDeserialize: deserialize_feapi_InflightsPruneResponse,
  },
  requestAmbientStatus: {
    path: '/feapi.FrontendAPI/RequestAmbientStatus',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.AmbientStatusRequest,
    responseType: feapi_feapi_pb.AmbientStatusResponse,
    requestSerialize: serialize_feapi_AmbientStatusRequest,
    requestDeserialize: deserialize_feapi_AmbientStatusRequest,
    responseSerialize: serialize_feapi_AmbientStatusResponse,
    responseDeserialize: deserialize_feapi_AmbientStatusResponse,
  },
  // ^ Client requests ambient status to be sent in.
  requestHomeView: {
    path: '/feapi.FrontendAPI/RequestHomeView',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.HomeViewRequest,
    responseType: feapi_feapi_pb.HomeViewResponse,
    requestSerialize: serialize_feapi_HomeViewRequest,
    requestDeserialize: deserialize_feapi_HomeViewRequest,
    responseSerialize: serialize_feapi_HomeViewResponse,
    responseDeserialize: deserialize_feapi_HomeViewResponse,
  },
  requestPopularView: {
    path: '/feapi.FrontendAPI/RequestPopularView',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.PopularViewRequest,
    responseType: feapi_feapi_pb.PopularViewResponse,
    requestSerialize: serialize_feapi_PopularViewRequest,
    requestDeserialize: deserialize_feapi_PopularViewRequest,
    responseSerialize: serialize_feapi_PopularViewResponse,
    responseDeserialize: deserialize_feapi_PopularViewResponse,
  },
  requestNotifications: {
    path: '/feapi.FrontendAPI/RequestNotifications',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.NotificationsRequest,
    responseType: feapi_feapi_pb.NotificationsResponse,
    requestSerialize: serialize_feapi_NotificationsRequest,
    requestDeserialize: deserialize_feapi_NotificationsRequest,
    responseSerialize: serialize_feapi_NotificationsResponse,
    responseDeserialize: deserialize_feapi_NotificationsResponse,
  },
  setNotificationsSignal: {
    path: '/feapi.FrontendAPI/SetNotificationsSignal',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.NotificationsSignalPayload,
    responseType: feapi_feapi_pb.NotificationsSignalResponse,
    requestSerialize: serialize_feapi_NotificationsSignalPayload,
    requestDeserialize: deserialize_feapi_NotificationsSignalPayload,
    responseSerialize: serialize_feapi_NotificationsSignalResponse,
    responseDeserialize: deserialize_feapi_NotificationsSignalResponse,
  },
  // rpc RequestOnboardCompleteStatus(OnboardCompleteStatusRequest) returns (OnboardCompleteStatusResponse) {}
  setOnboardComplete: {
    path: '/feapi.FrontendAPI/SetOnboardComplete',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.OnboardCompleteRequest,
    responseType: feapi_feapi_pb.OnboardCompleteResponse,
    requestSerialize: serialize_feapi_OnboardCompleteRequest,
    requestDeserialize: deserialize_feapi_OnboardCompleteRequest,
    responseSerialize: serialize_feapi_OnboardCompleteResponse,
    responseDeserialize: deserialize_feapi_OnboardCompleteResponse,
  },
  sendAddress: {
    path: '/feapi.FrontendAPI/SendAddress',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SendAddressPayload,
    responseType: feapi_feapi_pb.SendAddressResponse,
    requestSerialize: serialize_feapi_SendAddressPayload,
    requestDeserialize: deserialize_feapi_SendAddressPayload,
    responseSerialize: serialize_feapi_SendAddressResponse,
    responseDeserialize: deserialize_feapi_SendAddressResponse,
  },
  sendFEConfigChanges: {
    path: '/feapi.FrontendAPI/SendFEConfigChanges',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.FEConfigChangesPayload,
    responseType: feapi_feapi_pb.FEConfigChangesResponse,
    requestSerialize: serialize_feapi_FEConfigChangesPayload,
    requestDeserialize: deserialize_feapi_FEConfigChangesPayload,
    responseSerialize: serialize_feapi_FEConfigChangesResponse,
    responseDeserialize: deserialize_feapi_FEConfigChangesResponse,
  },
  requestBoardReports: {
    path: '/feapi.FrontendAPI/RequestBoardReports',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BoardReportsRequest,
    responseType: feapi_feapi_pb.BoardReportsResponse,
    requestSerialize: serialize_feapi_BoardReportsRequest,
    requestDeserialize: deserialize_feapi_BoardReportsRequest,
    responseSerialize: serialize_feapi_BoardReportsResponse,
    responseDeserialize: deserialize_feapi_BoardReportsResponse,
  },
  sendMintedUsernames: {
    path: '/feapi.FrontendAPI/SendMintedUsernames',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SendMintedUsernamesPayload,
    responseType: feapi_feapi_pb.SendMintedUsernamesResponse,
    requestSerialize: serialize_feapi_SendMintedUsernamesPayload,
    requestDeserialize: deserialize_feapi_SendMintedUsernamesPayload,
    responseSerialize: serialize_feapi_SendMintedUsernamesResponse,
    responseDeserialize: deserialize_feapi_SendMintedUsernamesResponse,
  },
  sendClientVersion: {
    path: '/feapi.FrontendAPI/SendClientVersion',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.ClientVersionPayload,
    responseType: feapi_feapi_pb.ClientVersionResponse,
    requestSerialize: serialize_feapi_ClientVersionPayload,
    requestDeserialize: deserialize_feapi_ClientVersionPayload,
    responseSerialize: serialize_feapi_ClientVersionResponse,
    responseDeserialize: deserialize_feapi_ClientVersionResponse,
  },
  sendSearchRequest: {
    path: '/feapi.FrontendAPI/SendSearchRequest',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.SearchRequestPayload,
    responseType: feapi_feapi_pb.SearchRequestResponse,
    requestSerialize: serialize_feapi_SearchRequestPayload,
    requestDeserialize: deserialize_feapi_SearchRequestPayload,
    responseSerialize: serialize_feapi_SearchRequestResponse,
    responseDeserialize: deserialize_feapi_SearchRequestResponse,
  },
  // ----------  Methods used by backend  ----------
  backendReady: {
    path: '/feapi.FrontendAPI/BackendReady',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BEReadyRequest,
    responseType: feapi_feapi_pb.BEReadyResponse,
    requestSerialize: serialize_feapi_BEReadyRequest,
    requestDeserialize: deserialize_feapi_BEReadyRequest,
    responseSerialize: serialize_feapi_BEReadyResponse,
    responseDeserialize: deserialize_feapi_BEReadyResponse,
  },
  sendBackendAmbientStatus: {
    path: '/feapi.FrontendAPI/SendBackendAmbientStatus',
    requestStream: false,
    responseStream: false,
    requestType: feapi_feapi_pb.BackendAmbientStatusPayload,
    responseType: feapi_feapi_pb.BackendAmbientStatusResponse,
    requestSerialize: serialize_feapi_BackendAmbientStatusPayload,
    requestDeserialize: deserialize_feapi_BackendAmbientStatusPayload,
    responseSerialize: serialize_feapi_BackendAmbientStatusResponse,
    responseDeserialize: deserialize_feapi_BackendAmbientStatusResponse,
  },
};

exports.FrontendAPIClient = grpc.makeGenericClientConstructor(FrontendAPIService);
