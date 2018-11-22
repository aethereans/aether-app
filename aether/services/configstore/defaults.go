// Services > ConfigStore
// This module handles saving and reading values from a config user file.

package configstore

import (
	"time"
)

/*

PORTS quick cheatsheet (assuming all are available)

# Backend
39999 talks to frontends
49999 talks to the Mim network

# Frontend
45001 talks to backend and frontend

# Client
47001 talks to the frontend

*/
// Defaults

// Backend defaults
const (
	defaultPoWBailoutTimeSeconds = 600
	// ^ We'll try to mint the PoW for 10 minutes before giving up. This is mostly for the long tail, so that if people attempt PoW levels that are impossible, they won't be stuck with their machine forever trying to compute it in vain. Anything above 24 is just very hard / impossible (30 secs to 10+ minutes on long tail) on a consumer machine as of mid-2018.
	defaultCacheGenerationIntervalHours            = 6
	defaultCacheDurationHours                      = 6
	defaultPOSTResponseExpiryMinutes               = 540 // 9h
	defaultPOSTResponseIneligibilityMinutes        = 480 // 8h
	defaultConnectionTimeout                       = 60 * time.Second
	defaultTCPConnectTimeout                       = 20 * time.Second
	defaultTLSHandshakeTimeout                     = 10 * time.Second
	defaultPingerPageSize                          = 100
	defaultOnlineAddressFinderPageSize             = 99
	defaultDispatchExclusionExpiryForLiveAddress   = 3 * time.Minute
	defaultDispatchExclusionExpiryForStaticAddress = 72 * time.Hour
	defaultPowStrength                             = 21
	defaultExternalIp                              = "127.0.0.1" // Localhost, if this is still 127.0.0.1 at any point in the future we failed at finding this out.
	defaultExternalIpType                          = 4           // IPv4
	defaultExternalPort                            = 49999
	defaultDbEngine                                = "sqlite" // 'sqlite' or 'mysql'
	defaultDBIp                                    = "127.0.0.1"
	defaultDbPort                                  = 3306
	defaultDbUsername                              = "aether-app-db-access-user"
	defaultDbPassword                              = "exventoveritas"
	defaultNeighbourCount                          = 5
	defaultMaxAddressTableSize                     = 1000
	defaultMaxInboundConns                         = 5
	defaultMaxOutboundConns                        = 1
	defaultMaxDbSizeMb                             = 10000
	defaultVotesMemoryDays                         = 14
	defaultBootstrapAfterOfflineMinutes            = 360
	defaultNodeType                                = 2
	defaultCacheGenerationInterval                 = 10 * time.Minute
)

// Frontend defaults
const (
	defaultFrontendExternalIp                      = "127.0.0.1"
	defaultFrontendExternalIpType                  = 4
	defaultFrontendAPIPort                         = 45001
	defaultBackendAPIAddress                       = "127.0.0.1"
	defaultClientAPIAddress                        = "127.0.0.1"
	defaultClientPort                              = 47001
	defaultThresholdForElectionValidityPercent     = 5
	defaultThresholdForElectionWinPercent          = 51 // If this is 60, any election needs 60% to be considered successful
	defaultBloomFilterSize                         = 10000
	defaultBloomFilterFalsePositiveRatePercent     = 50  // 50%, divide by 100 in use.
	defaultMinimumVoteThresholdForElectionValidity = 100 // Short of 10 votes on any direction, an election is not valid because the size is too small.
	defaultKvStoreRetentionDays                    = 180
	defaultLocalDevBackendDirectory                = "../../../aether-core/aether/backend"
)

// Shared defaults between frontend and backend

const (
	defaultGRPCServiceTimeout        = 60 * time.Second
	defaultBackendAPIPort            = 39999
	defaultMinimumTrustedPoWStrength = 6
	defaultNetworkHeadDays           = 14
	defaultNetworkMemoryDays         = 180
	defaultLocalMemoryDays           = 180
)

// Default entity page sizes

const (
	defaultBoardsPageSize      = 2000  // 0.2x
	defaultThreadsPageSize     = 400   // 1x
	defaultPostsPageSize       = 400   // 1x
	defaultVotesPageSize       = 2000  // 0.2x 2000
	defaultKeysPageSize        = 2000  // 0.2x
	defaultTruststatesPageSize = 3000  // 0.025x
	defaultAddressesPageSize   = 16000 // 0.025x

	defaultBoardIndexesPageSize      = 8000  // 0.025x
	defaultThreadIndexesPageSize     = 16000 // 0.025x
	defaultPostIndexesPageSize       = 12000 // 0.033x
	defaultVoteIndexesPageSize       = 4000  // 0.1x
	defaultKeyIndexesPageSize        = 20000 // 0.02x
	defaultTruststateIndexesPageSize = 15000 // 0.01x
	defaultAddressIndexesPageSize    = 16000 // 0.025x - Address is its own index
	// Every regular page is about 500kb that way.
	// Every index page is about 1mb.

	defaultBoardManifestsPageSize      = 30000
	defaultThreadManifestsPageSize     = 30000
	defaultPostManifestsPageSize       = 30000
	defaultVoteManifestsPageSize       = 30000
	defaultKeyManifestsPageSize        = 30000
	defaultTruststateManifestsPageSize = 30000
	defaultAddressManifestsPageSize    = 30000
	// Manifests are all the same size, so they're all the same.
)

/*----------  Entity-specific defaults  ----------*/

/*----------  Entity versions  ----------*/

const (
	defaultBoardEntityVersion       = 1
	defaultThreadEntityVersion      = 1
	defaultPostEntityVersion        = 1
	defaultVoteEntityVersion        = 1
	defaultKeyEntityVersion         = 1
	defaultTruststateEntityVersion  = 1
	defaultAddressEntityVersion     = 1
	defaultApiResponseEntityVersion = 1
)

/*----------  Other fields  ----------*/

const (
	/*----------  Key  ----------*/
	/*----------  Version: 1  ----------*/
	defaultKeyV1Type = "ed25519"
)
