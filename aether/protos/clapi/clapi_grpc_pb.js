// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var clapi_clapi_pb = require('../clapi/clapi_pb.js');
var feobjects_feobjects_pb = require('../feobjects/feobjects_pb.js');
var mimapi_structprotos_pb = require('../mimapi/structprotos_pb.js');

function serialize_clapi_AmbientLocalUserEntityPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientLocalUserEntityPayload)) {
    throw new Error('Expected argument of type clapi.AmbientLocalUserEntityPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientLocalUserEntityPayload(buffer_arg) {
  return clapi_clapi_pb.AmbientLocalUserEntityPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientLocalUserEntityResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientLocalUserEntityResponse)) {
    throw new Error('Expected argument of type clapi.AmbientLocalUserEntityResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientLocalUserEntityResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientLocalUserEntityResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientStatusPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientStatusPayload)) {
    throw new Error('Expected argument of type clapi.AmbientStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientStatusPayload(buffer_arg) {
  return clapi_clapi_pb.AmbientStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientStatusResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientStatusResponse)) {
    throw new Error('Expected argument of type clapi.AmbientStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientStatusResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientsRequest(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientsRequest)) {
    throw new Error('Expected argument of type clapi.AmbientsRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientsRequest(buffer_arg) {
  return clapi_clapi_pb.AmbientsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_AmbientsResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.AmbientsResponse)) {
    throw new Error('Expected argument of type clapi.AmbientsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_AmbientsResponse(buffer_arg) {
  return clapi_clapi_pb.AmbientsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_ExternalContentAutoloadDisabledStatusPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.ExternalContentAutoloadDisabledStatusPayload)) {
    throw new Error('Expected argument of type clapi.ExternalContentAutoloadDisabledStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_ExternalContentAutoloadDisabledStatusPayload(buffer_arg) {
  return clapi_clapi_pb.ExternalContentAutoloadDisabledStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_ExternalContentAutoloadDisabledStatusResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.ExternalContentAutoloadDisabledStatusResponse)) {
    throw new Error('Expected argument of type clapi.ExternalContentAutoloadDisabledStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_ExternalContentAutoloadDisabledStatusResponse(buffer_arg) {
  return clapi_clapi_pb.ExternalContentAutoloadDisabledStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_FEReadyRequest(arg) {
  if (!(arg instanceof clapi_clapi_pb.FEReadyRequest)) {
    throw new Error('Expected argument of type clapi.FEReadyRequest');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_FEReadyRequest(buffer_arg) {
  return clapi_clapi_pb.FEReadyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_FEReadyResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.FEReadyResponse)) {
    throw new Error('Expected argument of type clapi.FEReadyResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_FEReadyResponse(buffer_arg) {
  return clapi_clapi_pb.FEReadyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_HomeViewPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.HomeViewPayload)) {
    throw new Error('Expected argument of type clapi.HomeViewPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_HomeViewPayload(buffer_arg) {
  return clapi_clapi_pb.HomeViewPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_HomeViewResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.HomeViewResponse)) {
    throw new Error('Expected argument of type clapi.HomeViewResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_HomeViewResponse(buffer_arg) {
  return clapi_clapi_pb.HomeViewResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_ModModeEnabledStatusPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.ModModeEnabledStatusPayload)) {
    throw new Error('Expected argument of type clapi.ModModeEnabledStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_ModModeEnabledStatusPayload(buffer_arg) {
  return clapi_clapi_pb.ModModeEnabledStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_ModModeEnabledStatusResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.ModModeEnabledStatusResponse)) {
    throw new Error('Expected argument of type clapi.ModModeEnabledStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_ModModeEnabledStatusResponse(buffer_arg) {
  return clapi_clapi_pb.ModModeEnabledStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_NotificationsPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.NotificationsPayload)) {
    throw new Error('Expected argument of type clapi.NotificationsPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_NotificationsPayload(buffer_arg) {
  return clapi_clapi_pb.NotificationsPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_NotificationsResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.NotificationsResponse)) {
    throw new Error('Expected argument of type clapi.NotificationsResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_NotificationsResponse(buffer_arg) {
  return clapi_clapi_pb.NotificationsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_OnboardCompleteStatusPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.OnboardCompleteStatusPayload)) {
    throw new Error('Expected argument of type clapi.OnboardCompleteStatusPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_OnboardCompleteStatusPayload(buffer_arg) {
  return clapi_clapi_pb.OnboardCompleteStatusPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_OnboardCompleteStatusResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.OnboardCompleteStatusResponse)) {
    throw new Error('Expected argument of type clapi.OnboardCompleteStatusResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_OnboardCompleteStatusResponse(buffer_arg) {
  return clapi_clapi_pb.OnboardCompleteStatusResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_PopularViewPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.PopularViewPayload)) {
    throw new Error('Expected argument of type clapi.PopularViewPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_PopularViewPayload(buffer_arg) {
  return clapi_clapi_pb.PopularViewPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_PopularViewResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.PopularViewResponse)) {
    throw new Error('Expected argument of type clapi.PopularViewResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_PopularViewResponse(buffer_arg) {
  return clapi_clapi_pb.PopularViewResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_SearchResultPayload(arg) {
  if (!(arg instanceof clapi_clapi_pb.SearchResultPayload)) {
    throw new Error('Expected argument of type clapi.SearchResultPayload');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_SearchResultPayload(buffer_arg) {
  return clapi_clapi_pb.SearchResultPayload.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_clapi_SearchResultResponse(arg) {
  if (!(arg instanceof clapi_clapi_pb.SearchResultResponse)) {
    throw new Error('Expected argument of type clapi.SearchResultResponse');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_clapi_SearchResultResponse(buffer_arg) {
  return clapi_clapi_pb.SearchResultResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// These "Set", "Get" verbs are written from the viewpoint of the consumer of this api.
var ClientAPIService = exports.ClientAPIService = {
  // ----------  Methods used by frontend  ----------
  frontendReady: {
    path: '/clapi.ClientAPI/FrontendReady',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.FEReadyRequest,
    responseType: clapi_clapi_pb.FEReadyResponse,
    requestSerialize: serialize_clapi_FEReadyRequest,
    requestDeserialize: deserialize_clapi_FEReadyRequest,
    responseSerialize: serialize_clapi_FEReadyResponse,
    responseDeserialize: deserialize_clapi_FEReadyResponse,
  },
  deliverAmbients: {
    path: '/clapi.ClientAPI/DeliverAmbients',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientsRequest,
    responseType: clapi_clapi_pb.AmbientsResponse,
    requestSerialize: serialize_clapi_AmbientsRequest,
    requestDeserialize: deserialize_clapi_AmbientsRequest,
    responseSerialize: serialize_clapi_AmbientsResponse,
    responseDeserialize: deserialize_clapi_AmbientsResponse,
  },
  // ^ Ambient *entities*. This is poorly named, we should make it more specific
  sendAmbientStatus: {
    path: '/clapi.ClientAPI/SendAmbientStatus',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientStatusPayload,
    responseType: clapi_clapi_pb.AmbientStatusResponse,
    requestSerialize: serialize_clapi_AmbientStatusPayload,
    requestDeserialize: deserialize_clapi_AmbientStatusPayload,
    responseSerialize: serialize_clapi_AmbientStatusResponse,
    responseDeserialize: deserialize_clapi_AmbientStatusResponse,
  },
  // ^ Ambient *status* (fe, be status)
  sendAmbientLocalUserEntity: {
    path: '/clapi.ClientAPI/SendAmbientLocalUserEntity',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.AmbientLocalUserEntityPayload,
    responseType: clapi_clapi_pb.AmbientLocalUserEntityResponse,
    requestSerialize: serialize_clapi_AmbientLocalUserEntityPayload,
    requestDeserialize: deserialize_clapi_AmbientLocalUserEntityPayload,
    responseSerialize: serialize_clapi_AmbientLocalUserEntityResponse,
    responseDeserialize: deserialize_clapi_AmbientLocalUserEntityResponse,
  },
  sendHomeView: {
    path: '/clapi.ClientAPI/SendHomeView',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.HomeViewPayload,
    responseType: clapi_clapi_pb.HomeViewResponse,
    requestSerialize: serialize_clapi_HomeViewPayload,
    requestDeserialize: deserialize_clapi_HomeViewPayload,
    responseSerialize: serialize_clapi_HomeViewResponse,
    responseDeserialize: deserialize_clapi_HomeViewResponse,
  },
  sendPopularView: {
    path: '/clapi.ClientAPI/SendPopularView',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.PopularViewPayload,
    responseType: clapi_clapi_pb.PopularViewResponse,
    requestSerialize: serialize_clapi_PopularViewPayload,
    requestDeserialize: deserialize_clapi_PopularViewPayload,
    responseSerialize: serialize_clapi_PopularViewResponse,
    responseDeserialize: deserialize_clapi_PopularViewResponse,
  },
  sendNotifications: {
    path: '/clapi.ClientAPI/SendNotifications',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.NotificationsPayload,
    responseType: clapi_clapi_pb.NotificationsResponse,
    requestSerialize: serialize_clapi_NotificationsPayload,
    requestDeserialize: deserialize_clapi_NotificationsPayload,
    responseSerialize: serialize_clapi_NotificationsResponse,
    responseDeserialize: deserialize_clapi_NotificationsResponse,
  },
  sendOnboardCompleteStatus: {
    path: '/clapi.ClientAPI/SendOnboardCompleteStatus',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.OnboardCompleteStatusPayload,
    responseType: clapi_clapi_pb.OnboardCompleteStatusResponse,
    requestSerialize: serialize_clapi_OnboardCompleteStatusPayload,
    requestDeserialize: deserialize_clapi_OnboardCompleteStatusPayload,
    responseSerialize: serialize_clapi_OnboardCompleteStatusResponse,
    responseDeserialize: deserialize_clapi_OnboardCompleteStatusResponse,
  },
  sendModModeEnabledStatus: {
    path: '/clapi.ClientAPI/SendModModeEnabledStatus',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.ModModeEnabledStatusPayload,
    responseType: clapi_clapi_pb.ModModeEnabledStatusResponse,
    requestSerialize: serialize_clapi_ModModeEnabledStatusPayload,
    requestDeserialize: deserialize_clapi_ModModeEnabledStatusPayload,
    responseSerialize: serialize_clapi_ModModeEnabledStatusResponse,
    responseDeserialize: deserialize_clapi_ModModeEnabledStatusResponse,
  },
  sendExternalContentAutoloadDisabledStatus: {
    path: '/clapi.ClientAPI/SendExternalContentAutoloadDisabledStatus',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.ExternalContentAutoloadDisabledStatusPayload,
    responseType: clapi_clapi_pb.ExternalContentAutoloadDisabledStatusResponse,
    requestSerialize: serialize_clapi_ExternalContentAutoloadDisabledStatusPayload,
    requestDeserialize: deserialize_clapi_ExternalContentAutoloadDisabledStatusPayload,
    responseSerialize: serialize_clapi_ExternalContentAutoloadDisabledStatusResponse,
    responseDeserialize: deserialize_clapi_ExternalContentAutoloadDisabledStatusResponse,
  },
  sendSearchResult: {
    path: '/clapi.ClientAPI/SendSearchResult',
    requestStream: false,
    responseStream: false,
    requestType: clapi_clapi_pb.SearchResultPayload,
    responseType: clapi_clapi_pb.SearchResultResponse,
    requestSerialize: serialize_clapi_SearchResultPayload,
    requestDeserialize: deserialize_clapi_SearchResultPayload,
    responseSerialize: serialize_clapi_SearchResultResponse,
    responseDeserialize: deserialize_clapi_SearchResultResponse,
  },
};

exports.ClientAPIClient = grpc.makeGenericClientConstructor(ClientAPIService);
