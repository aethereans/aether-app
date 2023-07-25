// API > Boundscheck
// This file provides a boundary checker for all entities. This is a mandatory test for all entities before verification.

package api

import (
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Sizes for mathematical constants
// const (
// 	// Why not just ^(uint8(0)? Because we want these things untyped for now.
// 	toolbox.MaxInt64  = 9223372036854775807
// 	toolbox.MaxUint8  = 255
// 	toolbox.MaxUint16 = 65535
// 	toolbox.MaxUint32 = 4294967295
// )

// const (
// 	toolbox.MaxUint8        = 1<<8 - 1
// 	toolbox.MaxUint16       = 1<<16 - 1
// 	toolbox.MaxUint32 uint  = 1<<32 - 1
// 	toolbox.MaxInt32        = 1<<31 - 1
// 	toolbox.MaxInt64  int64 = 1<<63 - 1
// )

// Sizes for objects fields.
/*
  These values represent only what this specific node accepts as valid when it is receiving data from the outside.

  **There is a reason for these values being not configurable with either the JSON configuration or via command line flags.**

  Do not change these values without taking extreme caution. If you change these values on your own (i.e. not as a part of software update) the other nodes will interpret your posts as malicious and they will ban your node.
*/

// Lengths for objects:

// // V1

const (
	// Versionless, we start from 1. This is more pain than others to change, so try to not move the MAX_ENTITYVERSION between app versions if possible.
	MIN_ENTITYVERSION       = 1
	MAX_ENTITYVERSION int64 = toolbox.MaxInt64

	// Common
	MIN_META_V1 = 0
	MAX_META_V1 = 20480

	MIN_ENCR_CONTENT_V1 = 0
	MAX_ENCR_CONTENT_V1 = 131070

	MIN_PUBLICKEY_V1 = 32
	MAX_PUBLICKEY_V1 = 128

	// Board
	MIN_BOARD_NAME_V1 = 2
	MAX_BOARD_NAME_V1 = 128

	MIN_BOARD_DESCRIPTION_V1 = 0
	MAX_BOARD_DESCRIPTION_V1 = 65535

	MIN_BOARD_BOARDOWNERS_V1 = 0
	MAX_BOARD_BOARDOWNERS_V1 = 128

	MIN_BOARD_LANGUAGE_V1 = 0 // 3 char ISO 639-3 codes in lowercase
	// ^ 0 Because in the absence of language data, or when unrecognised, we assume Common Tongue.
	MAX_BOARD_LANGUAGE_V1 = 3

	// Thread
	MIN_THREAD_NAME_V1 = 2
	MAX_THREAD_NAME_V1 = 255

	MIN_THREAD_BODY_V1 = 0
	MAX_THREAD_BODY_V1 = 65535

	MIN_THREAD_LINK_V1 = 0
	MAX_THREAD_LINK_V1 = 4096

	// Post
	MIN_POST_BODY_V1 = 1
	MAX_POST_BODY_V1 = 65535

	// Vote
	MIN_VOTE_TYPE_V1 = 1
	MAX_VOTE_TYPE_V1 = toolbox.MaxUint16

	// Key
	MIN_KEY_TYPE_V1 = 1
	MAX_KEY_TYPE_V1 = 32

	MIN_KEY_NAME_V1 = 1
	MAX_KEY_NAME_V1 = 32

	MIN_KEY_INFO_V1 = 0
	MAX_KEY_INFO_V1 = 65535

	// Truststate

	MIN_TRUSTSTATE_TYPE_V1 = MIN_VOTE_TYPE_V1
	MAX_TRUSTSTATE_TYPE_V1 = MAX_VOTE_TYPE_V1

	// Address

	MIN_ADDRESS_LOCATIONTYPE_V1 = 0
	MAX_ADDRESS_LOCATIONTYPE_V1 = toolbox.MaxUint8

	MIN_ADDRESS_PORT_V1 = 0
	MAX_ADDRESS_PORT_V1 = toolbox.MaxUint16

	MIN_ADDRESS_TYPE_V1 = 0
	MAX_ADDRESS_TYPE_V1 = toolbox.MaxUint8

	MIN_ADDRESS_CLIENT_NAME_V1 = 2
	MAX_ADDRESS_CLIENT_NAME_V1 = 128

	MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_V1 = 1
	MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_V1 = 32

	MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1 = 2
	MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1 = 16

	MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1 = 1
	MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1 = 128

	MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1 = 1
	MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1 = 32

	// ApiResponse
	MIN_APIRESPONSE_ENTITY_NAME_V1_0 = MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1
	MAX_APIRESPONSE_ENTITY_NAME_V1_0 = MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1

	MIN_APIRESPONSE_ENDPOINT_NAME_V1_0 = MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1
	MAX_APIRESPONSE_ENDPOINT_NAME_V1_0 = MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1

	MIN_APIRESPONSE_PAGINATION_PAGES_V1_0 = 0
	MAX_APIRESPONSE_PAGINATION_PAGES_V1_0 = toolbox.MaxInt64

	MIN_APIRESPONSE_FILTER_V1_0 = 0
	MAX_APIRESPONSE_FILTER_V1_0 = 3

	MIN_APIRESPONSE_FILTER_TYPE_V1_0 = 0
	MAX_APIRESPONSE_FILTER_TYPE_V1_0 = toolbox.MaxUint16

	MIN_APIRESPONSE_FILTER_VALUES_FINGERPRINT_V1_0 = 1
	MAX_APIRESPONSE_FILTER_VALUES_FINGERPRINT_V1_0 = toolbox.MaxUint16

	MIN_APIRESPONSE_FILTER_VALUES_EMBED_V1_0 = 1
	MAX_APIRESPONSE_FILTER_VALUES_EMBED_V1_0 = toolbox.MaxUint16

	MIN_APIRESPONSE_FILTER_VALUES_EMBED_NAME_V1_0 = MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1
	MAX_APIRESPONSE_FILTER_VALUES_EMBED_NAME_V1_0 = MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1

	MIN_APIRESPONSE_FILTER_VALUES_TIMESTAMP_V1_0 = 2
	MAX_APIRESPONSE_FILTER_VALUES_TIMESTAMP_V1_0 = 2

	MIN_APIRESPONSE_CACHING_CACHEURL_V1_0 = 0
	MAX_APIRESPONSE_CACHING_CACHEURL_V1_0 = 128 // 64 char sha256 hash + some additions like POST response timestamp, etc.

	MIN_APIRESPONSE_RESULTCACHE_V1_0 = 0
	MAX_APIRESPONSE_RESULTCACHE_V1_0 = toolbox.MaxUint16

	MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0 = 0
	MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0 = toolbox.MaxInt32

	MIN_APIRESONSE_RESPONSEBODY_MANIFEST_ENTITY_V1_0 = 0
	MAX_APIRESONSE_RESPONSEBODY_MANIFEST_ENTITY_V1_0 = 50000

	// Indexes

	MIN_INDEX_PAGENUMBER_V1 = 0
	MAX_INDEX_PAGENUMBER_V1 = MAX_APIRESPONSE_PAGINATION_PAGES_V1_0
)

// Low level

func stringBC(item string, minLen, maxLen int) bool {
	if !(utf8.ValidString(item) && len(item) >= minLen && len(item) <= maxLen) {
		fmt.Printf("STRING CHECK FAIL. Item: %s, Min: %d, Max: %d\n", item, minLen, maxLen)
	}
	return utf8.ValidString(item) &&
		len(item) >= minLen &&
		len(item) <= maxLen
}

func stringSliceBC(item []string, minSliceLen, maxSliceLen, minItemLen, maxItemLen int) bool {
	sliceValid := intBC(int64(len(item)), int64(minSliceLen), int64(maxSliceLen))
	componentsValid := true
	if sliceValid {
		for _, val := range item {
			if !stringBC(val, minItemLen, maxItemLen) {
				componentsValid = false
				break
			}
		}
	}
	return sliceValid && componentsValid
}

// Arguably, we don't need this because the integers come with their own maxes and it will fail much before this if it fails. It's still good to have an explicit safeguard though.
func intBC(item int64, min, max int64) bool {
	if !(item >= min && item <= max) {
		fmt.Printf("INT CHECK FAIL: %d, min: %d, max: %d\n", item, min, max)
	}
	return item >= min &&
		item <= max
}

func boolBC(item bool, expected bool) bool {
	// if item != expected {
	// 	fmt.Printf("BOOL CHECK FAIL. Expected: %b, Received: %b\n", expected, item)
	// }
	if item != expected {
		return false
	}
	return true
}

// Low - mid level

func fingerprintBC(item Fingerprint) bool {
	// if !stringBC(string(item), 0, 64) {
	// 	fmt.Printf("FINGERPRINT FAIL: %s\n", item)
	// }
	return stringBC(string(item), 0, 64)
}

func nonceBC(item Nonce) bool {
	// Min 0 because nonce is only checked when there is a post request.
	return stringBC(string(item), 0, 64)
}
func fingerprintSliceBC(item *[]Fingerprint, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !fingerprintBC((*item)[key]) {
			return false
		}
	}
	return true
}
func timestampBC(item Timestamp) bool {
	// if !intBC(int64(item), 0, toolbox.MaxInt64) {
	// 	fmt.Printf("TIMESTAMP FAIL: %s\n", item)
	// }
	return intBC(int64(item), 0, toolbox.MaxInt64)
	// March 1, 2018 12AM UTC
}
func timestampSliceBC(item *[]Timestamp, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !timestampBC((*item)[key]) {
			return false
		}
	}
	return true
}
func proofOfWorkBC(item ProofOfWork) bool {
	// if !stringBC(string(item), 0, 512) {
	// 	fmt.Printf("POW FAIL: %s\n", item)
	// }
	return stringBC(string(item), 0, 512)
}
func signatureBC(item Signature) bool {
	// if !stringBC(string(item), 0, 512) {
	// 	fmt.Printf("SIG FAIL: %s\n", item)
	// }
	return stringBC(string(item), 0, 512)
}
func locationBC(item Location) bool {
	// if !stringBC(string(item), 0, 512) {
	// 	fmt.Printf("LOC FAIL: %s\n", item)
	// }
	return stringBC(string(item), 0, 512)
}
func publicKeyBC(item string, owner Fingerprint) bool {
	// if !stringBC(item, MIN_PUBLICKEY_V1, MAX_PUBLICKEY_V1) {
	// 	fmt.Printf("PK FAIL: %s\n", item)
	// }
	if len(item) == 0 && len(owner) == 0 {
		return true
	}
	return stringBC(item, MIN_PUBLICKEY_V1, MAX_PUBLICKEY_V1)
}
func nodePublicKeyBC(item string) bool {
	// if !stringBC(item, MIN_PUBLICKEY_V1, MAX_PUBLICKEY_V1) {
	// 	fmt.Printf("NODE PK FAIL: %s\n", item)
	// }
	return stringBC(item, MIN_PUBLICKEY_V1, MAX_PUBLICKEY_V1)
}

// Mid level

func provableBC(item *ProvableFieldSet) bool {
	return fingerprintBC(item.Fingerprint) &&
		timestampBC(item.Creation) &&
		proofOfWorkBC(item.ProofOfWork) &&
		signatureBC(item.Signature)
}

func updateableBC(item *UpdateableFieldSet) bool {
	return timestampBC(item.LastUpdate) &&
		proofOfWorkBC(item.UpdateProofOfWork) &&
		signatureBC(item.UpdateSignature)
}

// // Subentities

func boardOwnerBC(item *BoardOwner) bool {
	return fingerprintBC(item.KeyFingerprint) &&
		timestampBC(item.Expiry) &&
		intBC(int64(item.Level), 0, toolbox.MaxUint8)
}

func boardOwnerSliceBC(item *[]BoardOwner) bool {
	sliceValid := intBC(int64(len(*item)), MIN_BOARD_BOARDOWNERS_V1, MAX_BOARD_BOARDOWNERS_V1)
	boardOwnersValid := true
	if sliceValid {
		for key, _ := range *item {
			if !boardOwnerBC(&(*item)[key]) {
				boardOwnersValid = false
				break
			}
		}
	}
	return sliceValid && boardOwnersValid
}

func subprotocolBC(item *Subprotocol) bool {
	return stringBC(item.Name, MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1) &&
		intBC(int64(item.VersionMajor), 1, toolbox.MaxUint8) &&
		intBC(int64(item.VersionMinor), 0, toolbox.MaxUint16) &&
		stringSliceBC(item.SupportedEntities,
			MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1,
			MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1)
}

func subprotocolSliceBC(item *[]Subprotocol) bool {
	sliceValid := intBC(int64(len(*item)), MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_V1)
	subprotocolsValid := true
	if sliceValid {
		for key, _ := range *item {
			if !subprotocolBC(&(*item)[key]) {
				subprotocolsValid = false
				break
			}
		}
	}
	return sliceValid && subprotocolsValid
}

func protocolBC(item *Protocol) bool {
	// Optional, might not be present in inbound address.
	if item.VersionMajor == 0 &&
		item.VersionMinor == 0 &&
		len(item.Subprotocols) == 0 {
		return true
	}
	return intBC(int64(item.VersionMajor), 1, toolbox.MaxUint8) &&
		intBC(int64(item.VersionMinor), 0, toolbox.MaxUint16) &&
		subprotocolSliceBC(&item.Subprotocols)
}

func clientBC(item *Client) bool {
	// Optional, might not be present in inbound address.
	if item.ClientName == "" &&
		item.VersionMajor == 0 &&
		item.VersionMinor == 0 &&
		item.VersionPatch == 0 {
		return true
	}
	return intBC(int64(item.VersionMajor), 1, toolbox.MaxUint8) &&
		intBC(int64(item.VersionMinor), 0, toolbox.MaxUint16) &&
		intBC(int64(item.VersionPatch), 0, toolbox.MaxUint16) &&
		stringBC(item.ClientName, MIN_ADDRESS_CLIENT_NAME_V1, MAX_ADDRESS_CLIENT_NAME_V1)
}

func filterBC(item *Filter) bool {
	// This is an optional field. If fully empty, it returns true.
	if item.Type == "" && len(item.Values) == 0 {
		return true
	}
	allowed := (item.Type == "fingerprint" || item.Type == "embed" || item.Type == "timestamp")
	if !allowed {
		return false
	}
	valid := false
	if item.Type == "fingerprint" {
		valid = stringSliceBC(item.Values, 0, 64, // Fingerprint, but in string form
			MIN_APIRESPONSE_FILTER_VALUES_FINGERPRINT_V1_0, MAX_APIRESPONSE_FILTER_VALUES_FINGERPRINT_V1_0)
	} else if item.Type == "embed" {
		valid = stringSliceBC(item.Values,
			MIN_APIRESPONSE_FILTER_VALUES_EMBED_V1_0, MAX_APIRESPONSE_FILTER_VALUES_EMBED_V1_0,
			MIN_APIRESPONSE_FILTER_VALUES_EMBED_NAME_V1_0, MAX_APIRESPONSE_FILTER_VALUES_EMBED_NAME_V1_0)
	} else if item.Type == "timestamp" {
		var tss []Timestamp

		for _, val := range item.Values {
			ts, _ := strconv.ParseInt(val, 10, 64)
			tss = append(tss, Timestamp(ts))
		}
		valid = timestampSliceBC(&tss,
			MIN_APIRESPONSE_FILTER_VALUES_TIMESTAMP_V1_0,
			MAX_APIRESPONSE_FILTER_VALUES_TIMESTAMP_V1_0)
	}
	return valid
}

func filterSliceBC(item *[]Filter, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !filterBC(&(*item)[key]) {
			return false
		}
	}
	return true
}

func resultCacheBC(item *ResultCache) bool {
	// Optional
	if item.ResponseUrl == "" && item.StartsFrom == 0 && item.EndsAt == 0 { // take a look
		return true
	}
	valid := stringBC(item.ResponseUrl,
		MIN_APIRESPONSE_CACHING_CACHEURL_V1_0,
		MAX_APIRESPONSE_CACHING_CACHEURL_V1_0) &&
		timestampBC(item.StartsFrom) &&
		timestampBC(item.EndsAt)
	sane := item.StartsFrom < item.EndsAt
	return valid && sane
}

func resultCacheSliceBC(item *[]ResultCache, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !resultCacheBC(&(*item)[key]) {
			return false
		}
	}
	return true
}

func pageManifestEntityBC(item *PageManifestEntity) bool {
	return fingerprintBC(item.Fingerprint) &&
		timestampBC(item.LastUpdate)
}

func pageManifestEntitySliceBC(item *[]PageManifestEntity, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !pageManifestEntityBC(&(*item)[key]) {
			return false
		}
	}
	return true
}

func pageManifestBC(item *PageManifest) bool {
	return intBC(int64(item.Page),
		MIN_APIRESPONSE_PAGINATION_PAGES_V1_0,
		MAX_APIRESPONSE_PAGINATION_PAGES_V1_0) &&
		pageManifestEntitySliceBC(&item.Entities, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_ENTITY_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_ENTITY_V1_0)
}

// func pageManifestBC(item *PageManifest) bool {
// 	return intBC(int64(item.Page),
// 		MIN_APIRESPONSE_PAGINATION_PAGES_V1_0,
// 		MAX_APIRESPONSE_PAGINATION_PAGES_V1_0) &&
// 		fingerprintSliceBC(&item.Fingerprints, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_FINGERPRINTS_V1_0,
// 			MIN_APIRESONSE_RESPONSEBODY_MANIFEST_FINGERPRINTS_V1_0)
// }

func pageManifestSliceBC(item *[]PageManifest, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !pageManifestBC(&(*item)[key]) {
			return false
		}
	}
	return true
}

func entityCountBC(item *EntityCount) bool {
	return stringBC(item.Protocol, MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_NAME_V1) &&
		stringBC(item.Name,
			MIN_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1) &&
		intBC(int64(item.Count), 0, toolbox.MaxInt64)
}

func entityCountSliceBC(item *[]EntityCount, minLen, maxLen int) bool {
	sliceValid := intBC(int64(len(*item)), int64(minLen), int64(maxLen))
	if !sliceValid {
		return false
	}
	for key, _ := range *item {
		if !entityCountBC(&(*item)[key]) {
			return false
		}
	}
	return true
}

// Version-dependent internal API.

// Order for these: identity sets (Provable / Updateable), body fields are second, and slices are last (since they cost the most.) We want to bail as cheaply as possible.
func checkBoardBounds_V1(item *Board) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		stringBC(item.Name, MIN_BOARD_NAME_V1, MAX_BOARD_NAME_V1) &&
		stringBC(item.Description, MIN_BOARD_DESCRIPTION_V1, MAX_BOARD_DESCRIPTION_V1) &&
		fingerprintBC(item.Owner) &&
		publicKeyBC(item.OwnerPublicKey, item.Owner) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Language, MIN_BOARD_LANGUAGE_V1, MAX_BOARD_LANGUAGE_V1) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1) &&
		boardOwnerSliceBC(&item.BoardOwners)
}

// func checkBoardBounds_V1(item *Board) bool {
// 	return provableBC(&item.ProvableFieldSet)
// }

func checkThreadBounds_V1(item *Thread) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		fingerprintBC(item.Board) &&
		stringBC(item.Name, MIN_THREAD_NAME_V1, MAX_THREAD_NAME_V1) &&
		stringBC(item.Body, MIN_THREAD_BODY_V1, MAX_THREAD_BODY_V1) &&
		stringBC(item.Link, MIN_THREAD_LINK_V1, MAX_THREAD_LINK_V1) &&
		fingerprintBC(item.Owner) &&
		publicKeyBC(item.OwnerPublicKey, item.Owner) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1)
}
func checkPostBounds_V1(item *Post) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		fingerprintBC(item.Board) &&
		fingerprintBC(item.Thread) &&
		fingerprintBC(item.Parent) &&
		stringBC(item.Body, MIN_POST_BODY_V1, MAX_POST_BODY_V1) &&
		fingerprintBC(item.Owner) &&
		publicKeyBC(item.OwnerPublicKey, item.Owner) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1)
}
func checkVoteBounds_V1(item *Vote) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		fingerprintBC(item.Board) &&
		fingerprintBC(item.Thread) &&
		fingerprintBC(item.Target) &&
		fingerprintBC(item.Owner) &&
		publicKeyBC(item.OwnerPublicKey, item.Owner) &&
		intBC(int64(item.TypeClass), MIN_VOTE_TYPE_V1, MAX_VOTE_TYPE_V1) &&
		intBC(int64(item.Type), MIN_VOTE_TYPE_V1, MAX_VOTE_TYPE_V1) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1)
}
func checkKeyBounds_V1(item *Key) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		stringBC(item.Type, MIN_KEY_TYPE_V1, MAX_KEY_TYPE_V1) &&
		timestampBC(item.Expiry) &&
		stringBC(item.Name, MIN_KEY_NAME_V1, MAX_KEY_NAME_V1) &&
		stringBC(item.Info, MIN_KEY_INFO_V1, MAX_KEY_INFO_V1) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1)
}
func checkTruststateBounds_V1(item *Truststate) bool {
	return provableBC(&item.ProvableFieldSet) &&
		updateableBC(&item.UpdateableFieldSet) &&
		fingerprintBC(item.Target) &&
		fingerprintBC(item.Owner) &&
		publicKeyBC(item.OwnerPublicKey, item.Owner) &&
		intBC(int64(item.TypeClass), MIN_TRUSTSTATE_TYPE_V1, MAX_TRUSTSTATE_TYPE_V1) &&
		intBC(int64(item.Type), MIN_TRUSTSTATE_TYPE_V1, MAX_TRUSTSTATE_TYPE_V1) &&
		fingerprintBC(item.Domain) &&
		timestampBC(item.Expiry) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		stringBC(item.Meta, MIN_META_V1, MAX_META_V1) &&
		fingerprintBC(item.RealmId) &&
		stringBC(item.EncrContent, MIN_ENCR_CONTENT_V1, MAX_ENCR_CONTENT_V1)
}
func checkAddressBounds_V1(item *Address) bool {
	return locationBC(item.Location) &&
		locationBC(item.Sublocation) &&
		intBC(int64(item.LocationType), MIN_ADDRESS_LOCATIONTYPE_V1, MAX_ADDRESS_LOCATIONTYPE_V1) &&
		intBC(int64(item.Port), MIN_ADDRESS_PORT_V1, MAX_ADDRESS_PORT_V1) &&
		intBC(int64(item.Type), MIN_ADDRESS_TYPE_V1, MAX_ADDRESS_TYPE_V1) &&
		timestampBC(item.LastSuccessfulPing) &&
		timestampBC(item.LastSuccessfulSync) &&
		clientBC(&item.Client) &&
		intBC(int64(item.EntityVersion), MIN_ENTITYVERSION, MAX_ENTITYVERSION) &&
		fingerprintBC(item.RealmId)
}

func checkApiResponseBounds_V1_0(item *ApiResponse) bool {
	bodyOk := nodePublicKeyBC(item.NodePublicKey) &&
		signatureBC(item.Signature) &&
		nonceBC(item.Nonce) &&
		stringBC(item.Entity,
			MIN_APIRESPONSE_ENTITY_NAME_V1_0,
			MAX_APIRESPONSE_ENTITY_NAME_V1_0) &&
		stringBC(item.Endpoint,
			MIN_APIRESPONSE_ENDPOINT_NAME_V1_0,
			MAX_APIRESPONSE_ENDPOINT_NAME_V1_0) &&
		timestampBC(item.Timestamp) &&
		timestampBC(item.StartsFrom) &&
		timestampBC(item.EndsAt) &&
		intBC(int64(item.Pagination.Pages),
			MIN_APIRESPONSE_PAGINATION_PAGES_V1_0, MAX_APIRESPONSE_PAGINATION_PAGES_V1_0) &&
		intBC(int64(item.Pagination.CurrentPage), MIN_APIRESPONSE_PAGINATION_PAGES_V1_0, MAX_APIRESPONSE_PAGINATION_PAGES_V1_0) &&
		stringBC(item.Caching.CurrentCacheUrl, MIN_APIRESPONSE_CACHING_CACHEURL_V1_0, MAX_APIRESPONSE_CACHING_CACHEURL_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.BoardManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.ThreadManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.PostManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.VoteManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.KeyManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.TruststateManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		pageManifestSliceBC(&item.ResponseBody.AddressManifests, MIN_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0, MAX_APIRESONSE_RESPONSEBODY_MANIFEST_PAGES_V1_0) &&
		entityCountSliceBC(&item.Caching.EntityCounts, 0, MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_V1*MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1) // 32 subprotocols with 128 entities each is our max.
	if !bodyOk {
		logging.Logf(1, "This ApiResponse failed Boundscheck: %#v", item)
	}
	bodySane := item.Pagination.Pages >= item.Pagination.CurrentPage
	if !(bodyOk && bodySane) {
		return false
	}
	filtersOk := filterSliceBC(&item.Filters,
		MIN_APIRESPONSE_FILTER_V1_0,
		MAX_APIRESPONSE_FILTER_V1_0)
	if !filtersOk {
		return false
	}
	resultCachesOk := resultCacheSliceBC(&item.Results,
		MIN_APIRESPONSE_RESULTCACHE_V1_0,
		MAX_APIRESPONSE_RESULTCACHE_V1_0)
	if !resultCachesOk {
		return false
	}
	return bodyOk
}

func checkBoardIndexBounds_V1(item *BoardIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

func checkThreadIndexBounds_V1(item *ThreadIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		fingerprintBC(item.Board) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

func checkPostIndexBounds_V1(item *PostIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		fingerprintBC(item.Board) &&
		fingerprintBC(item.Thread) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

func checkVoteIndexBounds_V1(item *VoteIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		fingerprintBC(item.Board) &&
		fingerprintBC(item.Thread) &&
		fingerprintBC(item.Target) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

func checkKeyIndexBounds_V1(item *KeyIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

func checkTruststateIndexBounds_V1(item *TruststateIndex) bool {
	return fingerprintBC(item.Fingerprint) &&
		fingerprintBC(item.Target) &&
		timestampBC(item.Creation) &&
		timestampBC(item.LastUpdate) &&
		intBC(int64(item.PageNumber),
			MIN_INDEX_PAGENUMBER_V1, MAX_INDEX_PAGENUMBER_V1)
}

// High level version-independent API.
func (item *Board) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkBoardBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Thread) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkThreadBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Post) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkPostBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Vote) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkVoteBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Key) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkKeyBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Truststate) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkTruststateBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
func (item *Address) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkAddressBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func checkIndexes(item *Answer) (bool, error) {
	for key, _ := range item.BoardIndexes {
		valid, err := item.BoardIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.BoardIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.ThreadIndexes {
		valid, err := item.ThreadIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.ThreadIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.PostIndexes {
		valid, err := item.PostIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.PostIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.VoteIndexes {
		valid, err := item.VoteIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.VoteIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.KeyIndexes {
		valid, err := item.KeyIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.KeyIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.TruststateIndexes {
		valid, err := item.TruststateIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.TruststateIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	for key, _ := range item.AddressIndexes {
		valid, err := item.AddressIndexes[key].CheckBounds()
		if err != nil {
			return false, fmt.Errorf("Check Index encountered a failure. Object: %#v", item.AddressIndexes[key])
		}
		if !valid {
			return false, nil
		}
	}
	return true, nil
}

// CheckBounds of the ApiResponse does NOT check the entities contained within the APIResponse (including the address that it comes from at .Address field), only the ApiResponse's own structure, and index forms.
func (item *ApiResponse) CheckBounds() (bool, error) {
	if item.Address.Protocol.VersionMajor == 1 &&
		item.Address.Protocol.VersionMinor == 0 {
		indexesValid, err := checkIndexes(&item.ResponseBody)
		if err != nil {
			return false, fmt.Errorf("ApiResponse bounds checker encountered an error. Error: %#v", err)
		}
		return checkApiResponseBounds_V1_0(item) && indexesValid, nil
		// return checkApiResponseBounds_V1_0(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

// These are boundary checks for the index forms.

func (item *BoardIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkBoardIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *ThreadIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkThreadIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *PostIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkPostIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *VoteIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkVoteIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *KeyIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkKeyIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *TruststateIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		return checkTruststateIndexBounds_V1(item), nil
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}

func (item *AddressIndex) CheckBounds() (bool, error) {
	if item.EntityVersion == 1 {
		addr := Address(*item)
		return checkAddressBounds_V1(&addr), nil // Not addressIndex - address. Because they're the same.
	} else {
		return false, fmt.Errorf("We do not support this version of this entity for bounds checking. Entity: %#v", item)
	}
}
