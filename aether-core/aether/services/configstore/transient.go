// Services > ConfigStore
// This module handles saving and reading values from a config user file.

package configstore

import (
	pb "aether-core/aether/backend/metrics/proto"
	"aether-core/aether/protos/clapi"
	"aether-core/aether/services/nonces"
	"sync"
	// "time"
)

// TRANSIENT CONFIG

// These are the items that are set in runtime, and do not change until the application closes. This is different from the application state in the way that they're set-once for the runtime.

// These do not have getters and setters.

var Btc BackendTransientConfig
var Ftc FrontendTransientConfig

// Backend

// Default entity versions for this version of the app. This is not user adjustable.

type entityVersions struct {
	Board       int
	Thread      int
	Post        int
	Vote        int
	Key         int
	Truststate  int
	Address     int
	ApiResponse int
}

/*
The way we open a reverse connection is that we open a conn to the local server, and we open a conn to the remote server, and after sending the TCPMim message to request reverse open, we pipe one conn to another.

							C1															C2
LOCAL SERVER <--> LOCAL END <PIPE> LOCAL END <--> REMOTE SERVER
^ Local Remote    ^ Local Local    ^ Local Local  ^ Remote Remote
*/
type reverseConnData struct {
	C1LocalLocalAddr string
	C1LocalLocalPort uint16
}

/*
#### NONCOMMITTED ITEMS

## PermConfigReadOnly
When enabled, this prevents anything from saved into the config. This value itself is NOT saved into the config, so when the application restarts, this value is reset to false. This is useful in the case that you provide flags to the executable, but you don't want the values in the flags to be permanently saved into the config file. Any flags being provided into the executable will set this to true, therefore any runs with flags will effectively treat the config as read-only.

## AppIdentifier
This is the name of the app as registered to the operating system. This is useful to have here, because what we can do is we can vary this number in the swarm testing (petridish) and each of these nodes will act like a network in a single local machine, each with their own databases and different config files.

## OrgIdentifier
Same as above, but it's probably best to keep it under the same org name just to keep the local machine clean.

## PrintToStdout
This is useful because the logging things the normal kind does not pass the output to the swarm test orchestrator. This flag being enabled routes the logs to stdout so that the orchestrator can show it.

## MetricsDebugMode
This being enabled temporarily makes this node send much more detailed metrics more frequently, so that network connectivity issues can be debugged. This is a transient config on purpose, so that this cannot be enabled permanently. If a frontend connects to a backend with debug mode enabled, it has to show a warning to its user that says this backend node has debugging enabled, and only connect if the user agrees. Mind that the backend doesn't have to be truthful about whether it has the debug mode on. Having this mode on does not immediately compromise the frontend's privacy / identity, but the longer the frontend stays on that backend and the more actions a user commits, the higher the likelihood.

## ExternalPortVerified
Whether the port that was in the config was actually checked to be free and clear. This is important because we'll check once before the server starts to run, and when it starts, that port will no longer be available, and will start to return 'not available'. That will make all subsequent checks fail and that will trigger the port to be moved to a port that is free - but not bound to any server, since the server is bound to the old port, and that in fact is the reason the checks return false.

## SwarmNodeId
This is the number that this specific node will route to the main swarm orchestrator when it's reporting logs. Make sure that the App identifier (Usually in the format of "Aether-N") matches this number N, or it can be confusing.

## ShutdownInitiated
This is set when the shutdown of the backend service is initiated. The processes that take a long time to return should be checking this value periodically, and if it is set, they should stop whatever they're doing and do a graceful shutdown.

## LameduckInitiated
In this mode, the node starts to decline any inbound and outbound requests as well as reverse connections. This is effectively a way to prepare the node for an orderly shutdown. It will finish writing whatever db write it needs to do, and whenever that is done, shutdown can be initiated safely.

## StopAddressScannerCycle
This is the channel to send the message to when you want to stop the address scanner repeated task.

## StopUPNPCycle
This is the channel to send the message to when you want to stop the UPNP mapper repeated task.

## StopCacheGenerationCycle
This is the channel to send the message to when you want to stop the cache generator repeated task.

## AddressesScannerActive
This is the mutex that gets activated when the address scanner is active, so that it cannot be triggered twice at the same time.

## SyncActive
There is currently an ongoing sync happening.

## CurrentMetricsPage
This is the current metrics struct that we are building to send to the metrics server, if enabled.

## ConfigMutex
This is the mutex that prevents configuration from being written from multiple places.

## FingerprintCheckEnabled
Determines whether the entities coming over from the wire are fingerprint-checked for integrity.

## SignatureCheckEnabled
Determines whether the entities coming over from the wire are signature-checked for ownership.

## ProofOfWorkCheckEnabled
Determines whether the entities coming over from the wire are PoW-checked for anti-spam.

## PageSignatureCheckEnabled
Determines whether the pages (entity containers) coming over from the wire are signature-checked for integrity.

## EntityVersions
These are the versions of the entities that we can issue in this version of the app. Mind that this is for issuance, not for acceptance - we should still accept older versions gracefully.

# POSTResponseRepo
This is the repository that we keep our post responses in, so that they can be reused. This resets at every restart.

# NeighboursList
Our list of neighbours that we are checking in with at given intervals.

# Bouncer
Bouncer controls the inbound and outbound connections. This is the library that starts to refuse connections if the node gets too busy.

# ReverseConnData
This is the place we use to save the data so that we know an inbound connection is a reverse-opened one.

# Nonces
This is the library that keeps track of nonces for us.

# NewContentCommitted
We flip this flag to true whenever the users (frontends) of this backend send us new content. It triggers a reverse connection request at the first opportunity, so that the new content can spread to the network as soon as possible. It then flips itself to false after a reverse connection.

# MinimumTrustedPoWStrength
This is the PoW strength we ask for trusted entities coming from a CA that this node has explicitly chosen to trust.

# TetheredToFrontend
This backend was started by a frontend, therefore we should be sending status data to it. If this is not true, the backend will act as a standalone entity.

# AllowLocalhostRemotes
This allows 127.0.0.1 (localhost) as a valid remote address. This is crucial in swarm testing, in that in a swarm test all remotes will be on localhost.
*/

type BackendTransientConfig struct {
	ConfigMutex                sync.Mutex
	PermConfigReadOnly         bool
	AppIdentifier              string
	OrgIdentifier              string
	PrintToStdout              bool
	MetricsDebugMode           bool
	ExternalPortVerified       bool
	SwarmNodeId                int
	ShutdownInitiated          bool
	LameduckInitiated          bool
	StopNeighbourhoodCycle     chan bool
	StopInboundConnectionCycle chan bool
	StopExplorerCycle          chan bool
	StopAddressScannerCycle    chan bool
	StopNetworkScanCycle       chan bool
	StopUPNPCycle              chan bool
	StopCacheGenerationCycle   chan bool
	StopBadlistRefreshCycle    chan bool
	AddressesScannerActive     sync.Mutex
	SyncActive                 sync.Mutex
	CurrentMetricsPage         pb.Metrics
	FingerprintCheckEnabled    bool
	SignatureCheckEnabled      bool
	ProofOfWorkCheckEnabled    bool
	PageSignatureCheckEnabled  bool
	EntityVersions             entityVersions
	POSTResponseRepo           POSTResponseRepo // empty at start, empty at every app start
	NeighboursList             NeighboursList
	Bouncer                    Bouncer
	TLSEnabled                 bool
	ReverseConnData            reverseConnData
	Nonces                     nonces.RemotesNonces
	NewContentCommitted        bool
	BackendAPIPortVerified     bool
	MinimumTrustedPoWStrength  int
	TetheredToFrontend         bool
	AllowLocalhostRemotes      bool
}

// Set transient backend config defaults. Only need to set defaults that are not the type default.

// Mind that if you somehow manage to call something before SetDefaults is called, it will return its zero value without warning. This transient config does not have a Initialised gate that we can check, because adding that gate would have us convert everything in this place to getters / setters. We might do that in the future, but the point of BTC/FTC is that these are the things where the default value of the thing is the empty value of that variable.

// The problem here is that the default value of the field being empty value of that variable type and configs that need to be transient don't exactly match. So we will probably eventually move to a get/set model where it checks for init.

func (config *BackendTransientConfig) SetDefaults() {
	config.AppIdentifier = "Aether"
	config.OrgIdentifier = "Air Labs"
	// config.ConfigMutex = &sync.Mutex{}
	config.TLSEnabled = true
	// config.DispatcherExclusions = make(map[*interface{}]time.Time)
	// config.FingerprintCheckEnabled = true
	// config.SignatureCheckEnabled = true
	// config.ProofOfWorkCheckEnabled = true
	// config.PageSignatureCheckEnabled = true
	config.Nonces = nonces.NewRemotesNonces()
	config.MinimumTrustedPoWStrength = defaultMinimumTrustedPoWStrength

	// debug
	config.FingerprintCheckEnabled = false
	config.SignatureCheckEnabled = false
	config.ProofOfWorkCheckEnabled = false
	config.PageSignatureCheckEnabled = false
	// config.TLSEnabled = false

	ev := entityVersions{
		Board:       defaultBoardEntityVersion,
		Thread:      defaultThreadEntityVersion,
		Post:        defaultPostEntityVersion,
		Vote:        defaultVoteEntityVersion,
		Key:         defaultKeyEntityVersion,
		Truststate:  defaultTruststateEntityVersion,
		Address:     defaultAddressEntityVersion,
		ApiResponse: defaultApiResponseEntityVersion,
	}
	config.EntityVersions = ev
}

// Frontend

/*

# RefresherCacheNowTimestamp
This is the synchronised end timestamp for all things related to refreshing the FE. Why is this needed? Because the way FE cache works is that we pull all the necessary data into the frontend first as a whole to minimise the number of back and forths between the frontend and the backend.

That means, though, when that cache is queried, the last piece of data that is available on the cache will be the now() at the time of the beginning of the cache query (t1). That means if you query the cache at t2 and save it as now() of your query of the cache, you'll have the data from t1, but you'll timestamp it as t2, which means the next time you query, you'll start from t1, not t2, missing the data between t1 and t2.

To prevent that, at the beginning of a refresh cycle, the first thing that happens is pull-all-data-from-be event, and that event supplies a cache end timestamp. Anything that happens inside the refresh cycle assumes that timestamp is now() for all intents and purposes.

# SilenceNotificationsOnce

This is triggered when the search index being not present triggers a regeneration of the frontend kv store. Since this regeneration is going to create a lot of unread notifications, the notifications generator listens to this signal, and if this is set, the new notifications are generated as read, once. Whenever all notifications are marked as read, this is flipped back to allow for new unreads.
*/
type FrontendTransientConfig struct {
	ConfigMutex                 sync.Mutex
	PermConfigReadOnly          bool
	MetricsDebugMode            bool
	PrintToStdout               bool
	ShutdownInitiated           bool
	AppIdentifier               string
	OrgIdentifier               string
	FrontendAPIPortVerified     bool
	MinimumTrustedPoWStrength   int
	RefresherCacheNowTimestamp  int64
	CurrentAmbientStatus        clapi.AmbientStatusPayload
	StopRefresherCycle          chan bool
	StopSFWListUpdateCycle      chan bool
	StopNotificationsPruneCycle chan bool
	BackendReady                bool
	DefaultKeyType              string
	EntityVersions              entityVersions
	RefresherMutex              sync.Mutex
	SilenceNotificationsOnce    bool
}

// Set transient frontend config defaults

func (config *FrontendTransientConfig) SetDefaults() {
	config.PermConfigReadOnly = false
	// config.ConfigMutex = &sync.Mutex{}
	config.AppIdentifier = "Aether"
	config.OrgIdentifier = "Air Labs"
	config.MinimumTrustedPoWStrength = defaultMinimumTrustedPoWStrength
	config.DefaultKeyType = defaultKeyV1Type
	ev := entityVersions{
		Board:       defaultBoardEntityVersion,
		Thread:      defaultThreadEntityVersion,
		Post:        defaultPostEntityVersion,
		Vote:        defaultVoteEntityVersion,
		Key:         defaultKeyEntityVersion,
		Truststate:  defaultTruststateEntityVersion,
		Address:     defaultAddressEntityVersion,
		ApiResponse: defaultApiResponseEntityVersion,
	}
	config.EntityVersions = ev
}
