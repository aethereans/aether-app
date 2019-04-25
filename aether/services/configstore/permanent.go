// Services > ConfigStore
// This module handles saving and reading values from a config user file.

package configstore

import (
	"aether-core/aether/services/fingerprinting"
	"aether-core/aether/services/signaturing"
	"aether-core/aether/services/toolbox"
	"encoding/json"
	"errors"
	"fmt"
	cdir "github.com/shibukawa/configdir"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Config interface, so that we can actually have methods that take either frontend or backend config.

type Config interface {
	BlankCheck()
	SanityCheck()
}

/*
This package handles any data that gets saved to the user profile. This is important because everything that does not get saved into the database gets saved into this. Also important is this is where we allow multiple users to use the same database.
*/

// 0) UTILITY FUNCTIONS

func invalidDataError(input interface{}) error {
	return errors.New(fmt.Sprintf("An invalid value for this setting was provided by the user / application (in Set) or by the storage backend (in Get). Value provided: %#v", input))
}

// Maximums

const (
	maxPOWBailoutSeconds            = 3600 // 1h
	maxCacheGenerationIntervalHours = 168  // 7 days
	maxCacheDurationHours           = 72   // 3 days
	maxAbsolutePageSize             = 1000000
	maxPOWStrength                  = 63 // Our PoWs are 64 bytes long
	maxLocationSize                 = 2500
)

/*
... (ll. 116-138) Verily at the first Chaos came to be, but next wide-bosomed Earth, the ever-sure foundations of all the deathless ones who hold the peaks of snowy Olympus, and dim Tartarus in the depth of the wide-pathed Earth, and Eros, fairest among the deathless gods, who unnerves the limbs and overcomes the mind and wise counsels of all gods and all men within them. From Chaos came forth Erebus and black Night; but of Night were born Aether and Day, whom she conceived and bare from union in love with Erebus.

... In other words, I had to pick a random constant.

*/

const (
	night = 4386570
)

// 1) BACKEND

// Backend sub-entities

type EntityPageSizes struct {
	Boards      int
	Threads     int
	Posts       int
	Votes       int
	Keys        int
	Truststates int
	Addresses   int

	BoardIndexes      int
	ThreadIndexes     int
	PostIndexes       int
	VoteIndexes       int
	KeyIndexes        int
	TruststateIndexes int
	AddressIndexes    int

	BoardManifests      int
	ThreadManifests     int
	PostManifests       int
	VoteManifests       int
	KeyManifests        int
	TruststateManifests int
	AddressManifests    int
}

type MinimumPoWStrengths struct {
	Board            int
	BoardUpdate      int
	Thread           int
	ThreadUpdate     int
	Post             int
	PostUpdate       int
	Vote             int
	VoteUpdate       int
	Key              int
	KeyUpdate        int
	Truststate       int
	TruststateUpdate int
	ApiResponse      int
}

/*

This is an exact copy of the api.Subprotocol. This is here because we cannot import api here — it creates a circular reference. I've tried splitting API in many ways to avoid this issue, but each of the solutions to do that cause a lot more problems elsewhere since structs defined in the API have methods that reference other libraries, and moving those methods out of the structs mean the code gets a lot messier, etc. In short, unlikely as this sounds, creating a shim here and translating on the fly is the cleanest solution.

https://play.golang.org/p/x8wk4d7oAar <- an example of casting a shim to its actual thing. This could be worth it for the address as well, but address is a multi level entity so it might be not a one shot cast.. or maybe it would. Let's see. Ah yeah it doesn't work.

*/

type SubprotocolShim struct {
	Name              string   `json:"name"`
	VersionMajor      uint8    `json:"version_major"`
	VersionMinor      uint16   `json:"version_minor"`
	SupportedEntities []string `json:"supported_entities"`
}

// CONFIGS

var bc BackendConfig
var fc FrontendConfig

/*
Backend configuration.

## NetworkHeadDays
Days  of data that will be broadcast out in form of caches.

## NetworkMemoryDays
Days of data that will be provided to network upon request.

## LocalMemoryDays
Days of data to be kept before deletion.

## LastCacheGenerationTimestamp
The last time a new cache was generated locally.

## EntityPageSizes
How many entities will be put in a response page in POST responses and caches.

## MinimumPoWStrengths
The minimum number of zeros hashcash algorithm needs to have at the beginning of the PoW to accept it as valid.

## PoWBailoutTimeSeconds
How long does it take before a PoW timestamp is marked unattainable by the local computer. This is to make sure that the app doesn't keep attempting forever for an unattainably strong PoW it attempted to generate.

## CacheGenerationIntervalHours
How often does the node generate a new cache. By default, it generates a new cache every day. This is only used for the standard, static interval cache generation. DEPRECATED

## CacheDurationHours
When a cache is generated, how long does the cache cover? If your network head is 14 days, and your cache duration is 24 hours, you will generate 14 caches. If your cache duration is 6 hours, you'll generate 56 caches. This is only used for the standard, static interval cache generation. DEPRECATED

## ClientVersionMajor
Major version of the client software (Aether). x.0.0

## ClientVersionMinor
Minor version of the client software (Aether). 0.x.0

## ClientVersionPatch
Patch version of the client software (Aether). 0.0.x

## ClientName
Name of the client that this node is part of. (Aether)

## ProtocolVersionMajor
Major version of the Mim protocol that content is served over.

## ProtocolVersionMinor
Minor version of the Mim protocol that content is served over.

## POSTResponseExpiryMinutes
When a remote node makes a request via a POST response, a post response is generated, saved as a temporary file, and the access instructions are sent to a remote node. Remote node has a certain amount of time from this point on to fetch this response, around 30 minutes. After 30 minutes, this response is deleted.

Since we reuse POST responses, this should be always longer than the cache generation duration. If you're generating caches every 6 hours, This should probably be 9 hours. So that reused POST responses won't expire before a cache that covers for that timespan is generated. If so, you might end up having to generate POST responses that cover the entire 6 hours.

## POSTResponseIneligibilityMinutes
If the Post response has less than this many minutes left to expire, it is ineligible to be included in POST response chains.

This should be the same as cache generation duration, or slightly bigger to accommodate delays in cache generation. If you're generating caches every 6 hours, this should be 8 hours.

## ConnectionTimeout
How long the local node tries to attempt to connect to a remote node before deeming it unusable.

## TCPConnectTimeout
How long the local node tries to attempt to establish a TCP connection to a remote node before deeming it unusable.

## TLSHandshakeTimeout
How long the local node tries to attempt to complete a TLS handshake to a remote node before deeming it unusable.

## PingerPageSize
Pinger goes through all available addresses to find out whether they are online or not. This is done to keep a list of nodes that are usually online and in a connectable state. Pinger does this in form of pages (because there are occasionally more addresses available than there are sockets available in the local machine). This number determines how many nodes Pinger will attempt to connect at the same time.

## OnlineAddressFinderPageSize
This page size is slightly different than above. This one is for the local database call. Effectively, it looks at the most recent X addresses in the database to find ones that were active recently, and if that page does not yield enough online addresses, moves to the next page. This is to prevent querying a huge number of addresses.

## DispatchExclusionExpiryForLiveAddress
This is how long we wait until we reconnect to the same live address to look for updates.

## DispatchExclusionExpiryForStaticAddress
This is how long we wait until we reconnect to the same static address to look for updates.

## LoggingLevel
How deeply do we want to keep logs, or if any. 0 is no logs, 1 is medium, 2 is deep logs.

## ExternalIp
The external IP of this machine.

## ExternalIpType
The external IP type of this machine. 4: IPv4, 6: IPv6, 3: URL (in case of static)

## ExternalPort
The external port type of this machine.

## LastStaticAddressConnectionTimestamp
The last time we synced with a static node.

## LastLiveAddressConnectionTimestamp
The last time we synced with a live node.

## ServingSubprotocols
The subprotocols that this machine supports. In this case, c0 and dweb.

## NodeId
The node id of this machine. This is a randomly generated number. It does not have much significance beyond letting remote nodes keep their sync timestamps in check.

## UserDirectory
Where we save the backend , and if this node has a frontend, the frontend profile. This directory is given by the OS.

## CachesDirectory
Where we save the caches. This directory is given by the OS.

## Initialised
Whether the configuration file is properly initialised. If this is false, the initialisation did not complete.

## DbEngine
DbEngine allows the user to choose the database they want to use. SQLite is better for local installations where the app stays running on a desktop machine. It is simple and fast. MySQL is better when there are multiple users on the same backend, and it's a lot more robust against concurrent accesses. The preferred MySQL implementation is MariaDB, but original MySQL should also work.

Important: Do not forget that you have to create a DB called "aetherdb" in your preferred SQL engine with read/write access for the Username you give below.

(I thought of making this an iota and saving the numbers in this slot instead of string, but then that would make other parts of the code harder to read, because a DbEngine named 0 gives no information about what db engine it is, and you'd need to refer to this file to understand. I'd rather be infinitesimally less efficient and require less human RAM to read.)

## DbIP
This is the IP of the SQL server, if not SQLite3. By default, it's 127.0.0.1.

## DbPort
Port of the SQL server, if not SQLite3. By default, it's 3306 (MySQL default port)

## DbUsername
DbUsername is the username of the account that has read/write access to the "aetherdb" database, if not SQLite3. By default it's "aether-app-db-access-user".

## DbPassword
The password of the DB user, if not SQLite3. By default it's "exventoveritas". It's highly recommended that you change this if you're using MySQL. If you don't know what you're using, you're not using MySQL.

## MetricsLevel
## MetricsToken

## BackendKeyPair
Backend key pair is the key for this specific backend by which it signs the pages it creates. This is a combination of both private and public keys.

## AllowUnsignedEntities
If this is set to true, the node accepts posts that are anonymous. (But still with PoW and Fingerprint). This is disabled by default.

## MaxInboundPageSizeKb
Sets the threshold for bailout when a page being downloaded from the remote is too big.

# MaxAddressTableSize
This is how many addresses our database will hold at max.

# NeighbourCount
How many nodes we are interested in keeping in touch with on a rolling basis.

# MaxInboundConns
How many nodes do we allow to be simultaneously connected to this node. This number depends on your bandwidth and CPU resources. Setting this number to zero renders the config invalid (same as most things in config) and it will automatically regenerate from scratch, removing all prior config data.

# MaxOutboundConns
How many outbounds do we allow. Otherwise same as MaxInboundConns.

# MaxPingConns
How many ping (inbound 'hello's) do we allow. Otherwise same as MaxInboundConns.

# MaxDbSizeMb
This is the size that the user has allotted the application to use in the computer. Mind that this is only the database, and it is only the threshold where the event horizon starts to delete. Even when this threshold is not reached, if entities's last references reach the threshold of local memory, they will still be deleted.

# VotesMemoryDays DEPRECATED TODO FUTURE REMOVE
How long will the votes be retained in memory. This is a special case of LocalMemoryDays. We retain the votes much fewer days than the rest of the items because they're much more numerous and much less information dense. That does not mean all voting information will disappear though - when the frontend compiles votes, the compiled vote counts will be retained normally.

# Scaled Mode
This mode is normally disabled, and it is automatically enabled if the event horizon crosses network head threshold. In this mode, the behaviour of the node changes to only collecting and sharing boards that the user(s) on this backend has explicitly requested vs. the normal behaviour where all data on the network is collected for a limited timespan.

# ScaledModeUserSet
If this is enabled, the user has made a decision to keep scaled mode on or off, and we will not be flipping it back and forth based on disk space pressure.

# LastBootstrapAddressConnectionTimestamp
This is the last successful bootstrap timestamp. Every time a bootstrap is completed, this runs. If a node remains offline long enough that a given amount of time passes, bootstrap runs again.

# BootstrapAfterOfflineMinutes
This is how long a node can remain offline before a bootstrap kicks in at next restart. So if this value is 6 hours, if you're offline for more than 6h, this means when you start the app again, it will do a bootstrap. If you set this to -1, this will be disabled.

# SOCKS5ProxyEnabled
If this is true, a proxy is used for all outbound connections. This is an advanced feature, please read below if you are planning to use this, and use this only if you understand what this means. This feature can be used to use Aether over Tor.

IMPORTANT: Mind that proxies are one-way gates, so your requests will be able to go out and get the data your node requests, but your node will not be able to receive and serve any requests coming in from the network. That means the data you create (posts, threads, upvotes ... anything) will NOT be able to leave under normal conditions. If nobody can connect to you, they cannot get the data you create, thus your data won't ever reach the network, but ...

...thankfully, Aether has some mitigations for this case, and if this is the case, your node will attempt to ask other nodes to connect to you by pre-establishing connections from your side and handing the pipe over to them. But the nodes you do this to have *no* obligation to connect to you using the pipe you gave them, and they will likely ignore your node's request if they are under load, for example.

In essence, if you use a proxy, the ability of your node to get your data spread to the network can be anything from effectively unharmed, to completely nonexistent, based on network load.

Aether already anonymises who posts which content: after a piece of content spreads to a few nodes, it is impossible to find its owners' IP address. Therefore privacy is already inbuilt to the system, you do not need to do anything extra for it. This feature is available for those that require extra privacy above this baseline.

The nature of proxy technology means that unless configured specifically for it, a proxy or a VPN will not be able to accept inbound connections into your machine. For that to happen, you need to set up 'remote' or 'reverse' port forwarding (as opposed to *local* port forwarding.) Google it, and make sure you understand the security implications of doing that.

This applies to proxies, Tor, and VPN equally. It's easier to find VPNs that allow you to do a remote port forwarding than SOCKS5 proxies, since it's natively supported with VPNs.

# SOCKS5ProxyAddress
This is the address of the proxy that is going ot be used if the proxy is enabled. If you enable proxies and see no network traffic, it probably means your proxy configuration is broken. In the case proxy doesn't respond, it won't transmit anything. If proxy refuses your connection, the backend will terminate itself so as to not expose you to non-proxied traffic.

# SOCKS5ProxyUsername
# SOCKS5ProxyPassword
Username and password for the SOCKS5 proxy if required.

# NodeType

This value sets the node class. See below for potential values. Currently extant options: 2, 3, 254, 255

## NodeType: 2 (LiveNode)
This is the default setting. This means your node will act as a standard member of the network.

## NodeType: 255 (StaticNode)
This means your node is a static node. Other people will only be able to make you GET requests, not POST requests. This in effect renders you a static web server for the remote. The other nodes will be able to sync with you, but only up to the point of your last cache generation. This saves CPU, because it's only POST requests that need processing and database seeks.

Note: While setting this here and running the backend to serve a static node is possible, one great thing about static nodes is that they don't depend on the backend running to be static nodes. You can export your whole node as a folder that contains only json files, and stick it to any http server, and that server *will* be a static Aether node. So this allows you to just post a copy of your node into any web host without any sort of application running except Apache or Nginx, and it will be a valid static node. If you want to update the data in the stash, export your node again, replace the folder on your web server and you'll have it updated. Pretty simple.

Effects:
- Remote nodes cannot make POST requests (this node's bouncer will start to autodecline all inbound POST requests)

- Remotes will save you as a static node, which means their speed of attempting to sync with you will drop to every 3 hours, if they choose you as a sync target.

- As a side effect of not responding to POST requests, the data you provide will be somewhat stale. If you are generating caches every 6 hours (which is default), that means your data can be stale up to 5h59m. (Normally, this period is covered by POST requests)

- Reverse connections will be disabled. If you are in this node, the content you create will take more time to be distributed to the network (but it won't affect your speed of getting content others create).

## NodeType: 3 (BootstrapNode)
The node declares itself as a bootstrap node. This means the node declares itself to be fairly available and has bandwidth to help out the network by bringing newcomers into the sync state.

Setting this will cause your node to use more CPU and more bandwidth. This should realistically be only set to true when this backend is running on a server somewhere.

It's great for supporting the network, but it's not great as a daily-use node. If you want to help the network, stick this in an AWS (or similar) device and run it as bootstrap.

Effects:

- Other nodes, after their bootstrap, will connect to you less in an attempt to conserve your bandwidth and make it available for newcomers. So if this mode is set, you should not use this node as your regular day-to-day node because the content generated directly on it will take more time than usual to propagate out. (But it won't affect this node's speed of getting content others create.)

- It will disable reverse connection. That means your node should be directly accessible from the Internet, not through a NAT, or so.

- Node will prioritise syncs with other bootstrap nodes over standard nodes.

## NodeType: 254 (StaticBootstrapNode)
This is a combination of a static node and a bootstrap node.

## NodeType: 4 (CANode)
This node is a CA that is principally concerned with serving the CA-specific trust signals that it generates. These are things such as name mappings, or f451 assignments. These nodes have no special software, it's just a self identification so that other nodes can regularly check with them. They're checked in the same loop as bootstrap nodes, with the same caveats, which means if a node switches to this node, the inbounds to that node will drastically drop.

## NodeType: 253 (StaticCANode)
A CA node that does not respond to POST requests. Same as described above.

# BackendAPIPublic
This determines whether the backend API endpoint, which is the GRPC API for frontends(') consumption, is publicly available over the Internet. This is disabled by default. The use case for this would be that you have a personal server that you want to run your backend on 24/7, and your frontend (GUI app) on your local machine connects to that personal server. In this case, you want to make this public, so you can connect to your server over the Internet. (If you do this, make sure that you're the only person that can connect by adding your frontend key into the backend).

# BackendAPIPort
This is the port of the backend API endpoint.

# AdminFrontendAddress
This is the address of the frontend that has spawned the backend. Backend will use this address to reach out to the frontend if the need arises. Can be blank.

# AdminFrontendPublicKey
This is the public key of the frontend that has spawned the backend. This is effectively the 'admin' frontend for the backend (there can be multiple frontends). The admin frontend can run privileged requests, and change the admin fe address to a specific frontend (or to itself).

# GRPCServiceTimeout
How long does a GRPC service attempts to connect before considering the connection unusable.

# ExternalVerifyEnabled
This enables an extension interface for secondary verification of IP addresses. This feature requires the binary to be compiled with the appropriate flags to work.

# SQLiteDBLocation
This sets the location of the SQLite database. This is useful when you want to have your node's Db on a faster hard drive. A plausible use case for this is that you're running your node on a Raspberry Pi, and you want to move the Db to the faster external hard drive connected to it, rather than the SD card the RPi runs on. Mind that you need to shut down your app, move the existing DB, and then change this directory so it can find the moved file, in this order.

# DeclineInboundReverseRequests
This will make your node not listen to inbound reverse connection requests from other nodes. Flipping this to true reduces your connectivity, and there is normally no reason to flip this to true except for debugging.

# PreventOutboundReverseRequests
This will make your node not attempt outbound reverse connection requests. Flipping this to true reduces your connectivity, and there is normally no reason to flip this to true except for debugging.

# CacheGenerationInterval
Cache generation interval for dynamic cache generation with consolidation. This dynamic cache generation has variable size caches, and the caches closer to now are smaller. The benefit of this is that it reduces the time between the last cache's end and now(), which reduces the amount of data your node has to provide with a live POST request. The side effect is that your node uses more CPU, because it has to consolidate caches, which are, in the classic caching style, not touched ever after their creation.

# PrimaryBootstrap
Enabling this will make your bootstrap node be 'primary'. This means it will not share the bootstrapping load that it receives with other bootstrap nodes. If your node is not a bootstrap node, this setting has no effect.

# RenderNonconnectible
Enabling this will make your node not provide its port number to other remote nodes, rendering them unable to access and sync with you. This means the content you post will *not* leave your computer unless your computer is capable of reverse opens. This is good for debugging and simulating cases where the port mapping fails, but this is not a state you want to be in if you can help it.
*/

// Every time you add a new item here, please add getters, setters and to blankcheck method

// And before you think "hm, these would be better if they were private with lowercase letters..." that means you can't export them with JSON. Been there.

// Backend config base
type BackendConfig struct {
	NetworkHeadDays                         uint                // 14
	NetworkMemoryDays                       uint                // 180
	LocalMemoryDays                         uint                // 180
	LastCacheGenerationTimestamp            int64               //
	EntityPageSizes                         EntityPageSizes     //
	MinimumPoWStrengths                     MinimumPoWStrengths //
	PoWBailoutTimeSeconds                   uint                // 30
	CacheGenerationIntervalHours            uint                // 24
	CacheDurationHours                      uint                // 6
	ClientVersionMajor                      uint8               // 2 addr
	ClientVersionMinor                      uint16              // 0 addr
	ClientVersionPatch                      uint16              // 0 addr
	ClientName                              string              // Aether addr
	ProtocolVersionMajor                    uint8               // 1 (This refers to Mim, not subprotocols) addr
	ProtocolVersionMinor                    uint16              // 0 addr
	POSTResponseExpiryMinutes               uint                // 60
	POSTResponseIneligibilityMinutes        uint                // 10
	ConnectionTimeout                       time.Duration
	TCPConnectTimeout                       time.Duration
	TLSHandshakeTimeout                     time.Duration
	PingerPageSize                          uint
	OnlineAddressFinderPageSize             uint
	DispatchExclusionExpiryForLiveAddress   time.Duration
	DispatchExclusionExpiryForStaticAddress time.Duration
	LoggingLevel                            uint
	ExternalIp                              string // addr
	ExternalIpType                          uint8
	ExternalPort                            uint16
	LastStaticAddressConnectionTimestamp    int64
	LastLiveAddressConnectionTimestamp      int64
	ServingSubprotocols                     []SubprotocolShim
	NodeId                                  string
	UserDirectory                           string
	CachesDirectory                         string
	Initialised                             bool // False by default, init to set true
	DbEngine                                string
	DbIp                                    string // Only applies to non-sqlite
	DbPort                                  uint16 // Only applies to non-sqlite
	DbUsername                              string // Only applies to non-sqlite
	DbPassword                              string // Only applies to non-sqlite
	MetricsLevel                            uint8  // 0: no metrics transmitted
	MetricsToken                            string // If metrics level is not zero, metrics token is the anonymous identifier for the metrics server. Resetting this to 0 makes this node behave like a new node as far as metrics go, but if you don't want metrics to be collected, you can set it through the application or set the metrics level to zero in the JSON settings file.
	BackendKeyPair                          string // This is the Aether key, not TLS key
	MarshaledBackendPublicKey               string // This is the Aether key, not TLS key
	AllowUnsignedEntities                   bool
	MaxInboundPageSizeKb                    uint
	NeighbourCount                          uint
	MaxAddressTableSize                     uint
	MaxInboundConns                         uint
	MaxOutboundConns                        uint
	MaxPingConns                            uint
	MaxDbSizeMb                             uint
	VotesMemoryDays                         uint // 14
	EventHorizonTimestamp                   int64
	ScaledMode                              bool
	ScaledModeUserSet                       bool
	LastBootstrapAddressConnectionTimestamp int64
	BootstrapAfterOfflineMinutes            int // 360
	SOCKS5ProxyEnabled                      bool
	SOCKS5ProxyAddress                      string // Format: "127.0.0.1:65535"
	SOCKS5ProxyUsername                     string
	SOCKS5ProxyPassword                     string
	NodeType                                uint8
	BackendAPIPublic                        bool
	BackendAPIPort                          uint16
	AdminFrontendAddress                    string // Format: "127.0.0.1:65535"
	AdminFrontendPublicKey                  string
	GRPCServiceTimeout                      time.Duration
	ExternalVerifyEnabled                   bool
	SQLiteDBLocation                        string
	DeclineInboundReverseRequests           bool
	PreventOutboundReverseRequests          bool
	CacheGenerationInterval                 time.Duration
	PrimaryBootstrap                        bool
	RenderNonconnectible                    bool
}

// GETTERS AND SETTERS

/*
Q: Why do we even have these?

Because some of our types are not directly convertible to JSON, like the public / private key pairs.

Having this kind of set interface allows us to replace storage implementations later in the process without disrupting the rest of the app. The get / setter methods might look simple now, but they have no guarantee of remaining so in the future.

Q: Why the pain of uint as much as possible, then converting to ints?

Because we do not want users to provide negative values and make the application behave unpredictably. Any negative value should make the app not even start at all.
*/

// Init check gate

func (config *BackendConfig) InitCheck() {
	if !config.Initialised {
		log.Fatal(fmt.Sprintf("You've attempted to access a config before it was initialised. Trace: %v\n", toolbox.DumpStack()))
	}
}

// Getters
func (config *BackendConfig) GetLocalMemoryDays() int {
	config.InitCheck()
	if config.LocalMemoryDays < night &&
		config.LocalMemoryDays > 0 {
		return int(config.LocalMemoryDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LocalMemoryDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetNetworkMemoryDays() int {
	config.InitCheck()
	if config.NetworkMemoryDays < night &&
		config.NetworkMemoryDays > 0 {
		return int(config.NetworkMemoryDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NetworkMemoryDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetNetworkHeadDays() int {
	config.InitCheck()
	if config.NetworkHeadDays < night &&
		config.NetworkHeadDays > 0 {
		return int(config.NetworkHeadDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NetworkHeadDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetLastCacheGenerationTimestamp() int64 {
	config.InitCheck()
	if config.LastCacheGenerationTimestamp < toolbox.MaxInt64 { // can be zero
		return config.LastCacheGenerationTimestamp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LastCacheGenerationTimestamp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetEntityPageSizes() EntityPageSizes {
	config.InitCheck()
	if config.EntityPageSizes.Boards < maxAbsolutePageSize &&
		config.EntityPageSizes.Boards > 0 &&
		config.EntityPageSizes.Threads < maxAbsolutePageSize &&
		config.EntityPageSizes.Threads > 0 &&
		config.EntityPageSizes.Posts < maxAbsolutePageSize &&
		config.EntityPageSizes.Posts > 0 &&
		config.EntityPageSizes.Keys < maxAbsolutePageSize &&
		config.EntityPageSizes.Keys > 0 &&
		config.EntityPageSizes.Votes < maxAbsolutePageSize &&
		config.EntityPageSizes.Votes > 0 &&
		config.EntityPageSizes.Truststates < maxAbsolutePageSize &&
		config.EntityPageSizes.Truststates > 0 &&
		config.EntityPageSizes.Addresses < maxAbsolutePageSize &&
		config.EntityPageSizes.Addresses > 0 &&

		config.EntityPageSizes.BoardIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.BoardIndexes > 0 &&
		config.EntityPageSizes.ThreadIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.ThreadIndexes > 0 &&
		config.EntityPageSizes.PostIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.PostIndexes > 0 &&
		config.EntityPageSizes.KeyIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.KeyIndexes > 0 &&
		config.EntityPageSizes.VoteIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.VoteIndexes > 0 &&
		config.EntityPageSizes.TruststateIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.TruststateIndexes > 0 &&
		config.EntityPageSizes.AddressIndexes < maxAbsolutePageSize &&
		config.EntityPageSizes.AddressIndexes > 0 &&

		config.EntityPageSizes.BoardManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.BoardManifests > 0 &&
		config.EntityPageSizes.ThreadManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.ThreadManifests > 0 &&
		config.EntityPageSizes.PostManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.PostManifests > 0 &&
		config.EntityPageSizes.VoteManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.VoteManifests > 0 &&
		config.EntityPageSizes.KeyManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.KeyManifests > 0 &&
		config.EntityPageSizes.TruststateManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.TruststateManifests > 0 &&
		config.EntityPageSizes.AddressManifests < maxAbsolutePageSize &&
		config.EntityPageSizes.AddressManifests > 0 {
		return config.EntityPageSizes
	} else {
		log.Fatal(fmt.Sprintf("%#v", invalidDataError(config.EntityPageSizes)) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return EntityPageSizes{}
}
func (config *BackendConfig) GetMinimumPoWStrengths() MinimumPoWStrengths {
	config.InitCheck()
	if config.MinimumPoWStrengths.Board < maxPOWStrength &&
		config.MinimumPoWStrengths.Board > 0 &&
		config.MinimumPoWStrengths.BoardUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.BoardUpdate > 0 &&
		config.MinimumPoWStrengths.Thread < maxPOWStrength &&
		config.MinimumPoWStrengths.Thread > 0 &&
		config.MinimumPoWStrengths.ThreadUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.ThreadUpdate > 0 &&
		config.MinimumPoWStrengths.Post < maxPOWStrength &&
		config.MinimumPoWStrengths.Post > 0 &&
		config.MinimumPoWStrengths.PostUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.PostUpdate > 0 &&
		config.MinimumPoWStrengths.Vote < maxPOWStrength &&
		config.MinimumPoWStrengths.Vote > 0 &&
		config.MinimumPoWStrengths.VoteUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.VoteUpdate > 0 &&
		config.MinimumPoWStrengths.Key < maxPOWStrength &&
		config.MinimumPoWStrengths.Key > 0 &&
		config.MinimumPoWStrengths.KeyUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.KeyUpdate > 0 &&
		config.MinimumPoWStrengths.Truststate < maxPOWStrength &&
		config.MinimumPoWStrengths.Truststate > 0 &&
		config.MinimumPoWStrengths.TruststateUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.TruststateUpdate > 0 &&
		config.MinimumPoWStrengths.ApiResponse < maxPOWStrength &&
		config.MinimumPoWStrengths.ApiResponse > 0 {
		return config.MinimumPoWStrengths
	} else {
		log.Fatal(fmt.Sprintf("%#v", invalidDataError(config.MinimumPoWStrengths)) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return MinimumPoWStrengths{}
}
func (config *BackendConfig) GetPoWBailoutTimeSeconds() int {
	config.InitCheck()
	if config.PoWBailoutTimeSeconds < maxPOWBailoutSeconds &&
		config.PoWBailoutTimeSeconds > 0 {
		return int(config.PoWBailoutTimeSeconds)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.PoWBailoutTimeSeconds) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetCacheGenerationIntervalHours() int {
	config.InitCheck()
	if config.CacheGenerationIntervalHours < maxCacheGenerationIntervalHours &&
		config.CacheGenerationIntervalHours > 0 {
		return int(config.CacheGenerationIntervalHours)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.CacheGenerationIntervalHours) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetCacheDurationHours() int {
	config.InitCheck()
	if config.CacheDurationHours < maxCacheDurationHours &&
		config.CacheDurationHours > 0 {
		return int(config.CacheDurationHours)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.CacheDurationHours) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetClientVersionMajor() uint8 {
	config.InitCheck()
	if config.ClientVersionMajor < toolbox.MaxUint8 &&
		config.ClientVersionMajor > 0 {
		return config.ClientVersionMajor
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientVersionMajor) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetClientVersionMinor() uint16 {
	config.InitCheck()
	if config.ClientVersionMinor < toolbox.MaxUint16 { // can be zero
		return config.ClientVersionMinor
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientVersionMinor) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetClientVersionPatch() uint16 {
	config.InitCheck()
	if config.ClientVersionPatch < toolbox.MaxUint16 { // can be zero
		return config.ClientVersionPatch
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientVersionPatch) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetClientName() string {
	config.InitCheck()
	if len(config.ClientName) < toolbox.MaxUint8 &&
		len(config.ClientName) > 0 {
		return config.ClientName
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientName) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetProtocolVersionMajor() uint8 {
	config.InitCheck()
	if config.ProtocolVersionMajor < toolbox.MaxUint8 &&
		config.ProtocolVersionMajor > 0 {
		return config.ProtocolVersionMajor
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ProtocolVersionMajor) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetProtocolVersionMinor() uint16 {
	config.InitCheck()
	if config.ProtocolVersionMinor < toolbox.MaxUint16 { // can be zero
		return config.ProtocolVersionMinor
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ProtocolVersionMinor) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetPOSTResponseExpiryMinutes() int {
	config.InitCheck()
	if config.POSTResponseExpiryMinutes < toolbox.MaxInt32 &&
		config.POSTResponseExpiryMinutes > 0 {
		return int(config.POSTResponseExpiryMinutes)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.POSTResponseExpiryMinutes) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetPOSTResponseIneligibilityMinutes() int {
	config.InitCheck()
	if config.POSTResponseIneligibilityMinutes < toolbox.MaxInt32 &&
		config.POSTResponseIneligibilityMinutes > 0 {
		return int(config.POSTResponseIneligibilityMinutes)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.POSTResponseIneligibilityMinutes) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetConnectionTimeout() time.Duration {
	config.InitCheck()
	if config.ConnectionTimeout >= 1*time.Second { // Any value under is probably an attack.
		return config.ConnectionTimeout
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ConnectionTimeout) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}
func (config *BackendConfig) GetTCPConnectTimeout() time.Duration {
	config.InitCheck()
	if config.TCPConnectTimeout >= 1*time.Second { // Any value under is probably an attack.
		return config.TCPConnectTimeout
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.TCPConnectTimeout) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}
func (config *BackendConfig) GetTLSHandshakeTimeout() time.Duration {
	config.InitCheck()
	if config.TLSHandshakeTimeout >= 1*time.Second { // Any value under is probably an attack.
		return config.TLSHandshakeTimeout
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.TLSHandshakeTimeout) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}
func (config *BackendConfig) GetPingerPageSize() int {
	config.InitCheck()
	if config.PingerPageSize < toolbox.MaxInt32 &&
		config.PingerPageSize > 0 {
		return int(config.PingerPageSize)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.PingerPageSize) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetOnlineAddressFinderPageSize() int {
	config.InitCheck()
	if config.OnlineAddressFinderPageSize < toolbox.MaxInt32 &&
		config.OnlineAddressFinderPageSize > 0 {
		return int(config.OnlineAddressFinderPageSize)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.OnlineAddressFinderPageSize) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetDispatchExclusionExpiryForLiveAddress() time.Duration {
	config.InitCheck()
	if config.DispatchExclusionExpiryForLiveAddress >= 1*time.Minute { // Any value under is probably an attack.
		return config.DispatchExclusionExpiryForLiveAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DispatchExclusionExpiryForLiveAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}
func (config *BackendConfig) GetDispatchExclusionExpiryForStaticAddress() time.Duration {
	config.InitCheck()
	if config.DispatchExclusionExpiryForStaticAddress >= 1*time.Minute { // Any value under is probably an attack.
		return config.DispatchExclusionExpiryForStaticAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DispatchExclusionExpiryForStaticAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}
func (config *BackendConfig) GetLoggingLevel() int {
	config.InitCheck()
	if config.LoggingLevel < toolbox.MaxInt32 { // can be zero
		return int(config.LoggingLevel)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LoggingLevel) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetExternalIp() string {
	config.InitCheck()
	if len(config.ExternalIp) < maxLocationSize &&
		len(config.ExternalIp) > 0 {
		return config.ExternalIp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ExternalIp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetLastStaticAddressConnectionTimestamp() int64 {
	config.InitCheck()
	if config.LastStaticAddressConnectionTimestamp < toolbox.MaxInt64 { // can be zero
		return config.LastStaticAddressConnectionTimestamp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LastStaticAddressConnectionTimestamp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetLastLiveAddressConnectionTimestamp() int64 {
	config.InitCheck()
	if config.LastLiveAddressConnectionTimestamp < toolbox.MaxInt64 { // can be zero
		return config.LastLiveAddressConnectionTimestamp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LastLiveAddressConnectionTimestamp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetInitialised() bool {
	return config.Initialised
}
func (config *BackendConfig) GetServingSubprotocols() []SubprotocolShim {
	config.InitCheck()
	for _, val := range config.ServingSubprotocols {
		if len(val.SupportedEntities) == 0 {
			log.Fatal(invalidDataError(fmt.Sprintf("%#v", val.SupportedEntities) + " Trace: " + toolbox.Trace()))
		}
	}
	return config.ServingSubprotocols
}
func (config *BackendConfig) GetExternalIpType() uint8 {
	config.InitCheck()
	if config.ExternalIpType == 6 || config.ExternalIpType == 4 || config.ExternalIpType == 3 { // 6: ipv6, 4: ipv4, 3: URL (useful in static nodes)
		return config.ExternalIpType
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ExternalIpType) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetNodeId() string {
	config.InitCheck()
	if len(config.NodeId) == 64 {
		return config.NodeId
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NodeId) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetExternalPort() uint16 {
	config.InitCheck()
	if config.ExternalPort < toolbox.MaxUint16 && config.ExternalPort > 0 {
		return config.ExternalPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ExternalPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetUserDirectory() string {
	config.InitCheck()
	if len(config.UserDirectory) < toolbox.MaxUint16 &&
		len(config.UserDirectory) > 0 {
		return config.UserDirectory
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.UserDirectory) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetCachesDirectory() string {
	config.InitCheck()
	if len(config.CachesDirectory) < toolbox.MaxUint16 &&
		len(config.CachesDirectory) > 0 {
		return config.CachesDirectory
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.CachesDirectory) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetDbEngine() string {
	config.InitCheck()
	if config.DbEngine == "sqlite" || config.DbEngine == "mysql" {
		return config.DbEngine
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DbEngine) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetDbIp() string {
	config.InitCheck()
	if len(config.DbIp) < maxLocationSize &&
		len(config.DbIp) > 0 {
		return config.DbIp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DbIp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetDbPort() uint16 {
	config.InitCheck()
	if config.DbPort < toolbox.MaxUint16 && config.DbPort > 0 {
		return config.DbPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DbPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *BackendConfig) GetDbUsername() string {
	config.InitCheck()
	if len(config.DbUsername) < toolbox.MaxUint8 &&
		len(config.DbUsername) > 0 {
		return config.DbUsername
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DbUsername) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *BackendConfig) GetDbPassword() string {
	config.InitCheck()
	if len(config.DbPassword) < toolbox.MaxUint8 &&
		len(config.DbPassword) > 0 {
		return config.DbPassword
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DbPassword) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetMetricsLevel() uint8 {
	config.InitCheck()
	if config.MetricsLevel == 0 || config.MetricsLevel == 1 { // 0: no metrics, 1: anonymous metrics
		return config.MetricsLevel
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MetricsLevel) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMetricsToken() string {
	config.InitCheck()
	if len(config.MetricsToken) < 65 &&
		len(config.MetricsToken) >= 0 {
		return config.MetricsToken
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MetricsToken) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetBackendKeyPair() *ed25519.PrivateKey {
	config.InitCheck()
	keyPair, err := signaturing.UnmarshalPrivateKey(config.BackendKeyPair)
	if err != nil {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BackendKeyPair) + " Trace: " + toolbox.Trace() + "Error: " + err.Error()))
	}
	return &keyPair
}

func (config *BackendConfig) GetMarshaledBackendPublicKey() string {
	config.InitCheck()
	return config.MarshaledBackendPublicKey
}

func (config *BackendConfig) GetAllowUnsignedEntities() bool {
	config.InitCheck()
	return config.AllowUnsignedEntities
}

func (config *BackendConfig) GetMaxInboundPageSizeKb() int {
	config.InitCheck()
	if config.MaxInboundPageSizeKb < toolbox.MaxInt32 &&
		config.MaxInboundPageSizeKb > 500 { // can be zero
		return int(config.MaxInboundPageSizeKb)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxInboundPageSizeKb) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetNeighbourCount() int {
	config.InitCheck()
	if config.NeighbourCount < toolbox.MaxInt32 &&
		config.NeighbourCount > 0 {
		return int(config.NeighbourCount)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NeighbourCount) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMaxAddressTableSize() int {
	config.InitCheck()
	if config.MaxAddressTableSize < toolbox.MaxInt32 &&
		config.MaxAddressTableSize > 100 {
		return int(config.MaxAddressTableSize)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxAddressTableSize) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMaxInboundConns() int {
	config.InitCheck()
	if config.MaxInboundConns < toolbox.MaxInt32 &&
		config.MaxInboundConns > 0 {
		return int(config.MaxInboundConns)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxInboundConns) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMaxOutboundConns() int {
	config.InitCheck()
	if config.MaxOutboundConns < toolbox.MaxInt32 &&
		config.MaxOutboundConns > 0 {
		return int(config.MaxOutboundConns)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxOutboundConns) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMaxPingConns() int {
	config.InitCheck()
	if config.MaxPingConns < toolbox.MaxInt32 &&
		config.MaxPingConns > 0 {
		return int(config.MaxPingConns)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxPingConns) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetMaxDbSizeMb() int {
	config.InitCheck()
	if config.MaxDbSizeMb < toolbox.MaxInt32 &&
		config.MaxDbSizeMb > 0 {
		return int(config.MaxDbSizeMb)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MaxDbSizeMb) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetVotesMemoryDays() int {
	config.InitCheck()
	if config.VotesMemoryDays < toolbox.MaxInt32 &&
		config.VotesMemoryDays > 0 {
		return int(config.VotesMemoryDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.VotesMemoryDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetEventHorizonTimestamp() int64 {
	config.InitCheck()
	if config.EventHorizonTimestamp < toolbox.MaxInt64 &&
		config.EventHorizonTimestamp > 0 {
		return config.EventHorizonTimestamp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.EventHorizonTimestamp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetScaledMode() bool {
	config.InitCheck()
	return config.ScaledMode
}

func (config *BackendConfig) GetScaledModeUserSet() bool {
	config.InitCheck()
	return config.ScaledModeUserSet
}

func (config *BackendConfig) GetLastBootstrapAddressConnectionTimestamp() int64 {
	config.InitCheck()
	if config.LastBootstrapAddressConnectionTimestamp < toolbox.MaxInt64 &&
		config.LastBootstrapAddressConnectionTimestamp >= 0 {
		return config.LastBootstrapAddressConnectionTimestamp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LastBootstrapAddressConnectionTimestamp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetBootstrapAfterOfflineMinutes() int {
	config.InitCheck()
	if (config.BootstrapAfterOfflineMinutes < toolbox.MaxInt32 &&
		config.BootstrapAfterOfflineMinutes > 0) ||
		config.BootstrapAfterOfflineMinutes == -1 {
		return int(config.BootstrapAfterOfflineMinutes)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BootstrapAfterOfflineMinutes) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetSOCKS5ProxyEnabled() bool {
	config.InitCheck()
	return config.SOCKS5ProxyEnabled
}

func (config *BackendConfig) GetSOCKS5ProxyAddress() string {
	config.InitCheck()
	if len(config.SOCKS5ProxyAddress) < maxLocationSize &&
		len(config.SOCKS5ProxyAddress) >= 0 {
		return config.SOCKS5ProxyAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.SOCKS5ProxyAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetSOCKS5ProxyUsername() string {
	config.InitCheck()
	if len(config.SOCKS5ProxyUsername) < 1024 &&
		len(config.SOCKS5ProxyUsername) >= 0 {
		return config.SOCKS5ProxyUsername
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.SOCKS5ProxyUsername) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetSOCKS5ProxyPassword() string {
	config.InitCheck()
	if len(config.SOCKS5ProxyPassword) < 1024 &&
		len(config.SOCKS5ProxyPassword) >= 0 {
		return config.SOCKS5ProxyPassword
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.SOCKS5ProxyPassword) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetNodeType() uint8 {
	config.InitCheck()
	if config.NodeType == 2 || config.NodeType == 3 || config.NodeType == 254 || config.NodeType == 255 {
		return config.NodeType
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NodeType) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetBackendAPIPublic() bool {
	config.InitCheck()
	return config.BackendAPIPublic
}

func (config *BackendConfig) GetBackendAPIPort() uint16 {
	config.InitCheck()
	if config.BackendAPIPort < toolbox.MaxUint16 && config.BackendAPIPort > 0 {
		return config.BackendAPIPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BackendAPIPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *BackendConfig) GetAdminFrontendAddress() string {
	config.InitCheck()
	if len(config.AdminFrontendAddress) < maxLocationSize &&
		len(config.AdminFrontendAddress) >= 0 {
		return config.AdminFrontendAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.AdminFrontendAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetAdminFrontendPublicKey() string {
	config.InitCheck()
	return config.AdminFrontendPublicKey
}
func (config *BackendConfig) GetGRPCServiceTimeout() time.Duration {
	config.InitCheck()
	if config.GRPCServiceTimeout >= 1*time.Second { // Any value under is probably an attack.
		return config.GRPCServiceTimeout
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.GRPCServiceTimeout) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}

func (config *BackendConfig) GetExternalVerifyEnabled() bool {
	config.InitCheck()
	return config.ExternalVerifyEnabled
}

func (config *BackendConfig) GetSQLiteDBLocation() string {
	config.InitCheck()
	if len(config.SQLiteDBLocation) < toolbox.MaxUint16 &&
		len(config.SQLiteDBLocation) > 0 {
		return config.SQLiteDBLocation
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.SQLiteDBLocation) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *BackendConfig) GetDeclineInboundReverseRequests() bool {
	config.InitCheck()
	return config.DeclineInboundReverseRequests
}

func (config *BackendConfig) GetPreventOutboundReverseRequests() bool {
	config.InitCheck()
	return config.PreventOutboundReverseRequests
}

func (config *BackendConfig) GetCacheGenerationInterval() time.Duration {
	config.InitCheck()
	if config.CacheGenerationInterval >= 1*time.Minute { // Any value under is probably an attack.
		return config.CacheGenerationInterval
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.CacheGenerationInterval) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}

func (config *BackendConfig) GetPrimaryBootstrap() bool {
	config.InitCheck()
	return config.PrimaryBootstrap
}

func (config *BackendConfig) GetRenderNonconnectible() bool {
	config.InitCheck()
	return config.RenderNonconnectible
}

/*****************************************************************************/

// Setters

func (config *BackendConfig) SetLocalMemoryDays(val int) error {
	if val > 0 {
		config.InitCheck()
		config.LocalMemoryDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetNetworkMemoryDays(val int) error {
	config.InitCheck()
	if val > 0 {
		config.NetworkMemoryDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetNetworkHeadDays(val int) error {
	config.InitCheck()
	if val > 0 {
		config.NetworkHeadDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetLastCacheGenerationTimestamp(val int64) error {
	config.InitCheck()
	if val > 0 {
		config.LastCacheGenerationTimestamp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetEntityPageSizes(val EntityPageSizes) error {
	config.InitCheck()
	if val.Boards < maxAbsolutePageSize &&
		val.Boards > 0 &&
		val.BoardIndexes < maxAbsolutePageSize &&
		val.BoardIndexes > 0 &&
		val.Threads < maxAbsolutePageSize &&
		val.Threads > 0 &&
		val.ThreadIndexes < maxAbsolutePageSize &&
		val.ThreadIndexes > 0 &&
		val.Posts < maxAbsolutePageSize &&
		val.Posts > 0 &&
		val.PostIndexes < maxAbsolutePageSize &&
		val.PostIndexes > 0 &&
		val.Keys < maxAbsolutePageSize &&
		val.Keys > 0 &&
		val.KeyIndexes < maxAbsolutePageSize &&
		val.KeyIndexes > 0 &&
		val.Votes < maxAbsolutePageSize &&
		val.Votes > 0 &&
		val.VoteIndexes < maxAbsolutePageSize &&
		val.VoteIndexes > 0 &&
		val.Truststates < maxAbsolutePageSize &&
		val.Truststates > 0 &&
		val.TruststateIndexes < maxAbsolutePageSize &&
		val.TruststateIndexes > 0 &&

		val.BoardManifests < maxAbsolutePageSize &&
		val.BoardManifests > 0 &&
		val.ThreadManifests < maxAbsolutePageSize &&
		val.ThreadManifests > 0 &&
		val.PostManifests < maxAbsolutePageSize &&
		val.PostManifests > 0 &&
		val.VoteManifests < maxAbsolutePageSize &&
		val.VoteManifests > 0 &&
		val.KeyManifests < maxAbsolutePageSize &&
		val.KeyManifests > 0 &&
		val.TruststateManifests < maxAbsolutePageSize &&
		val.TruststateManifests > 0 {
		config.EntityPageSizes = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetMinimumPoWStrengths(powStr int) error {
	config.InitCheck()
	var mps MinimumPoWStrengths
	if powStr > 6 && powStr < maxPOWStrength {
		mps.Board = powStr
		mps.BoardUpdate = powStr
		mps.Thread = powStr
		mps.ThreadUpdate = powStr
		mps.Post = powStr
		mps.PostUpdate = powStr
		mps.Vote = powStr
		mps.VoteUpdate = powStr
		mps.Key = powStr
		mps.KeyUpdate = powStr
		mps.Truststate = powStr
		mps.TruststateUpdate = powStr
		mps.ApiResponse = powStr - 6
		// ^ ApiResponse PoW strength is a little softer because remotes need it to make POST requests to each other. Temp value, subject to solidification or change in the future.
		config.MinimumPoWStrengths = mps
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		} else {
			return invalidDataError(fmt.Sprintf("%#v", powStr) + " Trace: " + toolbox.Trace())
		}
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetPoWBailoutTimeSeconds(val int) error {
	config.InitCheck()
	if val > 0 {
		config.PoWBailoutTimeSeconds = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetCacheGenerationIntervalHours(val int) error {
	config.InitCheck()
	if val > 0 {
		config.CacheGenerationIntervalHours = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetCacheDurationHours(val int) error {
	config.InitCheck()
	if val > 0 {
		config.CacheDurationHours = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetClientVersionMajor(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint8 {
		config.ClientVersionMajor = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetClientVersionMinor(val int) error {
	config.InitCheck()
	if val >= 0 && val < toolbox.MaxUint16 {
		config.ClientVersionMinor = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetClientVersionPatch(val int) error {
	config.InitCheck()
	if val >= 0 && val < toolbox.MaxUint16 {
		config.ClientVersionPatch = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetClientName(val string) error {
	config.InitCheck()
	if len(val) > 0 {
		config.ClientName = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetProtocolVersionMajor(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint8 {
		config.ProtocolVersionMajor = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetProtocolVersionMinor(val int) error {
	config.InitCheck()
	if val >= 0 && val < toolbox.MaxUint16 {
		config.ProtocolVersionMinor = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetPOSTResponseExpiryMinutes(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.POSTResponseExpiryMinutes = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetPOSTResponseIneligibilityMinutes(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.POSTResponseIneligibilityMinutes = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetConnectionTimeout(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Second { // Any value under is probably an attack.
		config.ConnectionTimeout = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetTCPConnectTimeout(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Second { // Any value under is probably an attack.
		config.TCPConnectTimeout = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetTLSHandshakeTimeout(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Second { // Any value under is probably an attack.
		config.TLSHandshakeTimeout = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetPingerPageSize(val int) error {
	config.InitCheck()
	if val > 0 {
		config.PingerPageSize = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetOnlineAddressFinderPageSize(val int) error {
	config.InitCheck()
	if val > 0 {
		config.OnlineAddressFinderPageSize = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDispatchExclusionExpiryForLiveAddress(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Minute {
		config.DispatchExclusionExpiryForLiveAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDispatchExclusionExpiryForStaticAddress(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Minute { // Any value under is probably an attack.
		config.DispatchExclusionExpiryForStaticAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetLoggingLevel(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.LoggingLevel = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetExternalIp(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.ExternalIp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetLastStaticAddressConnectionTimestamp(val int64) error {
	config.InitCheck()
	if val >= 0 {
		config.LastStaticAddressConnectionTimestamp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetLastLiveAddressConnectionTimestamp(val int64) error {
	config.InitCheck()
	if val >= 0 {
		config.LastLiveAddressConnectionTimestamp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetInitialised(val bool) error {
	// No init check on this one so we can start inserting data.
	config.Initialised = true
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}
func (config *BackendConfig) SetServingSubprotocols(subprotocols []interface{}) error {
	config.InitCheck()
	var castSubprots []SubprotocolShim
	for _, val := range subprotocols {
		item, ok := val.(SubprotocolShim)
		if !ok {
			return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
		}
		castSubprots = append(castSubprots, item)
	}
	config.ServingSubprotocols = castSubprots
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}
func (config *BackendConfig) SetExternalIpType(val int) error {
	config.InitCheck()
	if val == 6 || val == 4 || val == 3 {
		config.ExternalIpType = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetNodeId(val string) error {
	config.InitCheck()
	if len(val) == 64 {
		config.NodeId = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetExternalPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.ExternalPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetUserDirectory(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint16 {
		config.UserDirectory = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetCachesDirectory(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint16 {
		config.CachesDirectory = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDbEngine(val string) error {
	config.InitCheck()
	if val == "mysql" || val == "sqlite" {
		config.DbEngine = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDbIp(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.DbIp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDbPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.DbPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDbUsername(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint8 {
		config.DbUsername = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetDbPassword(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint8 {
		config.DbPassword = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMetricsLevel(val int) error {
	config.InitCheck()
	if val == 0 || val == 1 {
		config.MetricsLevel = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMetricsToken(val string) error {
	config.InitCheck()
	if len(val) >= 0 && len(val) < 65 {
		config.MetricsToken = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetBackendKeyPair(val *ed25519.PrivateKey) error {
	config.InitCheck()
	config.BackendKeyPair = signaturing.MarshalPrivateKey(*val)
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

// The only way to set this is to set backend key pair first.
func (config *BackendConfig) SetMarshaledBackendPublicKey(val *ed25519.PrivateKey) error {
	config.InitCheck()
	marshaledPk := signaturing.MarshalPublicKey(val.Public().(ed25519.PublicKey))
	config.MarshaledBackendPublicKey = marshaledPk
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetAllowUnsignedEntities(val bool) error {
	config.InitCheck()
	config.AllowUnsignedEntities = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetMaxInboundPageSizeKb(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxInboundPageSizeKb = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetNeighbourCount(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.NeighbourCount = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMaxAddressTableSize(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxAddressTableSize = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMaxInboundConns(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxInboundConns = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetMaxOutboundConns(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxOutboundConns = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMaxPingConns(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxPingConns = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetMaxDbSizeMb(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.MaxDbSizeMb = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetVotesMemoryDays(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.VotesMemoryDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *BackendConfig) SetEventHorizonTimestamp(val int64) error {
	config.InitCheck()
	if val > 0 {
		config.EventHorizonTimestamp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetScaledMode(val bool) error {
	config.InitCheck()
	config.ScaledMode = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetScaledModeUserSet(val bool) error {
	config.InitCheck()
	config.ScaledModeUserSet = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetLastBootstrapAddressConnectionTimestamp(val int64) error {
	config.InitCheck()
	if val >= 0 {
		config.LastBootstrapAddressConnectionTimestamp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetBootstrapAfterOfflineMinutes(val int) error {
	config.InitCheck()
	if val > 0 || val == -1 {
		config.BootstrapAfterOfflineMinutes = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetSOCKS5ProxyEnabled(val bool) error {
	config.InitCheck()
	config.SOCKS5ProxyEnabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetSOCKS5ProxyAddress(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.SOCKS5ProxyAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetSOCKS5ProxyUsername(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < 1024 {
		config.SOCKS5ProxyUsername = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetSOCKS5ProxyPassword(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < 1024 {
		config.SOCKS5ProxyPassword = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetNodeType(val int) error {
	config.InitCheck()
	if val == 2 || val == 3 || val == 254 || val == 255 {
		config.NodeType = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetBackendAPIPublic(val bool) error {
	config.InitCheck()
	config.BackendAPIPublic = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetBackendAPIPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.BackendAPIPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetAdminFrontendAddress(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.AdminFrontendAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetAdminFrontendPublicKey(val string) error {
	config.InitCheck()
	config.AdminFrontendPublicKey = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetGRPCServiceTimeout(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Second { // Any value under is probably an attack.
		config.GRPCServiceTimeout = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetExternalVerifyEnabled(val bool) error {
	config.InitCheck()
	config.ExternalVerifyEnabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetSQLiteDBLocation(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint16 {
		config.SQLiteDBLocation = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetDeclineInboundReverseRequests(val bool) error {
	config.InitCheck()
	config.DeclineInboundReverseRequests = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetPreventOutboundReverseRequests(val bool) error {
	config.InitCheck()
	config.PreventOutboundReverseRequests = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetCacheGenerationInterval(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Minute {
		config.CacheGenerationInterval = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *BackendConfig) SetPrimaryBootstrap(val bool) error {
	config.InitCheck()
	config.PrimaryBootstrap = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *BackendConfig) SetRenderNonconnectible(val bool) error {
	config.InitCheck()
	config.RenderNonconnectible = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

/*****************************************************************************/

// BlankCheck looks at all variables and if it finds they're at their zero value, sets the default value for it. This is a guard against a new item being added to the config store as a result of a version update, but it being zero value. If a zero'd value is found, we change it to its default before anything else happens. This also effectively runs at the first pass to set the defaults.

func (config *BackendConfig) BlankCheck() {
	// Init needs to be first so that we can actually start editing these stuff.
	if !config.Initialised {
		config.SetInitialised(true)
	}
	if config.NetworkHeadDays == 0 {
		config.SetNetworkHeadDays(defaultNetworkHeadDays)
	}
	if config.NetworkMemoryDays == 0 {
		config.SetNetworkMemoryDays(defaultNetworkMemoryDays)
	}
	if config.LocalMemoryDays == 0 {
		config.SetLocalMemoryDays(defaultLocalMemoryDays)
	}
	// ::LastCacheGenerationTimestamp: can be zero, no need to blank check.
	if config.MinimumPoWStrengths.Board == 0 ||
		config.MinimumPoWStrengths.BoardUpdate == 0 ||
		config.MinimumPoWStrengths.Thread == 0 ||
		config.MinimumPoWStrengths.ThreadUpdate == 0 ||
		config.MinimumPoWStrengths.Post == 0 ||
		config.MinimumPoWStrengths.PostUpdate == 0 ||
		config.MinimumPoWStrengths.Vote == 0 ||
		config.MinimumPoWStrengths.VoteUpdate == 0 ||
		config.MinimumPoWStrengths.Key == 0 ||
		config.MinimumPoWStrengths.KeyUpdate == 0 ||
		config.MinimumPoWStrengths.Truststate == 0 ||
		config.MinimumPoWStrengths.TruststateUpdate == 0 ||
		config.MinimumPoWStrengths.ApiResponse == 0 {
		config.SetMinimumPoWStrengths(defaultPowStrength)
	}
	if config.EntityPageSizes.Boards == 0 ||
		config.EntityPageSizes.Threads == 0 ||
		config.EntityPageSizes.Posts == 0 ||
		config.EntityPageSizes.Votes == 0 ||
		config.EntityPageSizes.Keys == 0 ||
		config.EntityPageSizes.Truststates == 0 ||
		config.EntityPageSizes.Addresses == 0 ||

		config.EntityPageSizes.BoardIndexes == 0 ||
		config.EntityPageSizes.ThreadIndexes == 0 ||
		config.EntityPageSizes.PostIndexes == 0 ||
		config.EntityPageSizes.VoteIndexes == 0 ||
		config.EntityPageSizes.KeyIndexes == 0 ||
		config.EntityPageSizes.TruststateIndexes == 0 ||
		config.EntityPageSizes.AddressIndexes == 0 ||

		config.EntityPageSizes.BoardManifests == 0 ||
		config.EntityPageSizes.ThreadManifests == 0 ||
		config.EntityPageSizes.PostManifests == 0 ||
		config.EntityPageSizes.VoteManifests == 0 ||
		config.EntityPageSizes.KeyManifests == 0 ||
		config.EntityPageSizes.TruststateManifests == 0 ||
		config.EntityPageSizes.AddressManifests == 0 {
		config.setDefaultEntityPageSizes()
	}
	if config.PoWBailoutTimeSeconds == 0 {
		config.SetPoWBailoutTimeSeconds(defaultPoWBailoutTimeSeconds)
	}
	if config.CacheGenerationIntervalHours == 0 {
		config.SetCacheGenerationIntervalHours(defaultCacheGenerationIntervalHours)
	}
	if config.CacheDurationHours == 0 {
		config.SetCacheDurationHours(defaultCacheDurationHours)
	}
	if config.ClientVersionMajor == 0 {
		config.SetClientVersionMajor(clientVersionMajor)
	}
	if config.ClientVersionMinor != clientVersionMinor {
		config.SetClientVersionMinor(clientVersionMinor)
	}
	if config.ClientVersionPatch != clientVersionPatch {
		config.SetClientVersionPatch(clientVersionPatch)
	}
	if config.ClientName == "" || config.ClientName != clientName {
		config.SetClientName(clientName)
	}
	if config.ProtocolVersionMajor == 0 || config.ProtocolVersionMajor != protocolVersionMajor {
		config.SetProtocolVersionMajor(protocolVersionMajor)
	}
	if config.ProtocolVersionMinor != protocolVersionMinor {
		config.SetProtocolVersionMinor(protocolVersionMinor)
	}
	if config.POSTResponseExpiryMinutes == 0 {
		config.SetPOSTResponseExpiryMinutes(defaultPOSTResponseExpiryMinutes)
	}
	if config.POSTResponseIneligibilityMinutes == 0 {
		config.SetPOSTResponseIneligibilityMinutes(defaultPOSTResponseIneligibilityMinutes)
	}
	if config.ConnectionTimeout == 0 {
		config.SetConnectionTimeout(defaultConnectionTimeout)
	}
	if config.TCPConnectTimeout == 0 {
		config.SetTCPConnectTimeout(defaultTCPConnectTimeout)
	}
	if config.TLSHandshakeTimeout == 0 {
		config.SetTLSHandshakeTimeout(defaultTLSHandshakeTimeout)
	}
	if config.PingerPageSize == 0 {
		config.SetPingerPageSize(defaultPingerPageSize)
	}
	if config.OnlineAddressFinderPageSize == 0 {
		config.SetOnlineAddressFinderPageSize(defaultOnlineAddressFinderPageSize)
	}
	if config.DispatchExclusionExpiryForLiveAddress == 0 {
		config.SetDispatchExclusionExpiryForLiveAddress(defaultDispatchExclusionExpiryForLiveAddress)
	}
	if config.DispatchExclusionExpiryForStaticAddress == 0 {
		config.SetDispatchExclusionExpiryForStaticAddress(defaultDispatchExclusionExpiryForStaticAddress)
	}
	// ::LoggingLevel: can be zero, no need to blank check.
	// if config.LoggingLevel == 0 {
	// 	config.SetLoggingLevel(1)
	// }
	if config.ExternalIp == "" {
		config.SetExternalIp(defaultExternalIp)
	}
	if config.ExternalIpType == 0 {
		config.SetExternalIpType(defaultExternalIpType)
	}
	if config.ExternalPort == 0 {
		config.SetExternalPort(defaultExternalPort)
	}
	// ::LastStaticAddressConnectionTimestamp: can be zero, no need to blank check.
	// ::LastLiveAddressConnectionTimestamp: can be zero, no need to blank check.
	var servingSubprotocolsNeedRegeneration bool
	if len(config.ServingSubprotocols) == 0 {
		servingSubprotocolsNeedRegeneration = true
	} else {
		for _, val := range config.ServingSubprotocols {
			if len(val.SupportedEntities) == 0 {
				servingSubprotocolsNeedRegeneration = true
			}
		}
	}
	if servingSubprotocolsNeedRegeneration {
		c0 := SubprotocolShim{Name: "c0", VersionMajor: 1, VersionMinor: 0, SupportedEntities: []string{"board", "thread", "post", "vote", "key", "truststate"}}
		// dweb := SubprotocolShim{Name: "dweb", VersionMajor: 1, VersionMinor: 0, SupportedEntities: []string{"page"}}
		config.SetServingSubprotocols([]interface{}{c0})
	}
	if len(config.UserDirectory) == 0 {
		config.SetUserDirectory(cdir.New(Btc.OrgIdentifier, Btc.AppIdentifier).QueryFolders(cdir.Global)[0].Path)
	}
	if len(config.CachesDirectory) == 0 {
		config.SetCachesDirectory(cdir.New(Btc.OrgIdentifier, Btc.AppIdentifier).QueryCacheFolder().Path)
	}
	if len(config.DbEngine) == 0 {
		config.SetDbEngine(defaultDbEngine)
	}
	if len(config.DbIp) == 0 {
		config.SetDbIp(defaultDBIp)
	}
	if config.DbPort == 0 {
		config.SetDbPort(defaultDbPort)
	}
	if len(config.DbUsername) == 0 {
		config.SetDbUsername(defaultDbUsername)
	}
	if len(config.DbPassword) == 0 {
		config.SetDbPassword(defaultDbPassword)
	}
	// ::MetricsLevel: can be zero, no need to blank check.
	// ::MetricsToken: can be zero, no need to blank check.
	if len(config.BackendKeyPair) == 0 {
		privKey, _ := signaturing.CreateKeyPair()
		config.SetBackendKeyPair(privKey)
	}
	// This needs to be after Backend key pair generation.
	if len(config.MarshaledBackendPublicKey) == 0 {
		config.SetMarshaledBackendPublicKey(config.GetBackendKeyPair())
	}
	// This needs to be after key pair generation, because it uses the key pair. Node Id is the Fingerprint of the public key of the backend.
	if config.NodeId == "" {
		nodeid := fingerprinting.Create(config.GetMarshaledBackendPublicKey())
		config.SetNodeId(nodeid)
	}
	// ::AllowUnsignedEntities: can be false, no need to blank check.
	if config.MaxInboundPageSizeKb == 0 {
		config.SetMaxInboundPageSizeKb(150000)
	}
	if config.NeighbourCount == 0 {
		config.SetNeighbourCount(defaultNeighbourCount)
	}
	if config.MaxAddressTableSize == 0 {
		config.SetMaxAddressTableSize(defaultMaxAddressTableSize)
	}
	if config.MaxInboundConns == 0 {
		config.SetMaxInboundConns(defaultMaxInboundConns)
	}
	if config.MaxOutboundConns == 0 {
		config.SetMaxOutboundConns(defaultMaxOutboundConns)
	}
	if config.MaxPingConns == 0 {
		config.SetMaxPingConns(defaultMaxPingConns)
	}
	if config.MaxDbSizeMb == 0 {
		config.SetMaxDbSizeMb(defaultMaxDbSizeMb)
	}
	if config.VotesMemoryDays == 0 {
		config.SetVotesMemoryDays(defaultVotesMemoryDays)
	}
	if config.EventHorizonTimestamp == 0 {
		config.ResetEventHorizon()
	}
	// ::ScaledMode: can be false, no need to blank check.
	// ::ScaledModeUserSet: can be false, no need to blank check.
	// ::LastBootstrapAddressConnectionTimestamp: can be 0, no need to blank check.
	if config.BootstrapAfterOfflineMinutes == 0 {
		config.SetBootstrapAfterOfflineMinutes(defaultBootstrapAfterOfflineMinutes)
	}
	// ::SOCKS5ProxyEnabled: can be false, no need to blank check.
	// ::SOCKS5ProxyAddress: can be blank, no need to blank check.
	// ::SOCKS5ProxyUsername: can be blank, no need to blank check.
	// ::SOCKS5ProxyPassword: can be blank, no need to blank check.
	if config.NodeType == 0 {
		config.SetNodeType(defaultNodeType)
	}
	// ::BackendAPIPublic: can be false, no need to blank check.
	if config.BackendAPIPort == 0 {
		config.SetBackendAPIPort(defaultBackendAPIPort)
	}
	// ::AdminFrontendAddress: can be blank, no need to blank check.
	// ::AdminFrontendPublicKey: can be blank, no need to blank check.
	if config.GRPCServiceTimeout == 0 {
		config.SetGRPCServiceTimeout(defaultGRPCServiceTimeout)
	}
	// ::ExternalVerifyEnabled: can be false, no need to blank check.
	if len(config.SQLiteDBLocation) == 0 {
		config.SetSQLiteDBLocation(filepath.Join(config.UserDirectory, "backend"))
	}
	// ::DeclineInboundReverseRequests: can be false, no need to blank check.
	// ::PreventOutboundReverseRequests: can be false, no need to blank check.
	if config.CacheGenerationInterval == 0 {
		config.SetCacheGenerationInterval(defaultCacheGenerationInterval)
	}
	// ::PrimaryBootstrap: can be false, no need to blank check.
	// ::RenderNonconnectible: can be false, no need to blank check.
}

// Resets

func (config *BackendConfig) ResetEventHorizon() {
	localMemCutoff := time.Now().Add(-(time.Duration(config.LocalMemoryDays) * time.Hour * time.Duration(24))).Unix()
	config.SetEventHorizonTimestamp(localMemCutoff)
}

func (config *BackendConfig) ResetLastBootstrapAddressConnectionTimestamp() {
	config.SetLastBootstrapAddressConnectionTimestamp(0)
}

func (config *BackendConfig) ResetLastLiveAddressConnectionTimestamp() {
	config.SetLastLiveAddressConnectionTimestamp(0)
}

func (config *BackendConfig) ResetLastStaticAddressConnectionTimestamp() {
	config.SetLastStaticAddressConnectionTimestamp(0)
}

/*
Backend config sanity check.Everything you add to above, needs to also be added to the sanity check. This runs at the initialisation at the beginning of the program, and it checks that the values actually make sense. Sanity checks also run on gets and sets, but they don't normally run at startup. This function covers that base.
*/
func (config *BackendConfig) SanityCheck() {
	if !config.GetInitialised() {
		log.Fatal("Backend configuration is not initialised. Please initialise it before use.")
	} else {
		// If there is an error, the appropriate getter function will fail and crash the app.
		config.GetLocalMemoryDays()
		config.GetNetworkMemoryDays()
		config.GetNetworkHeadDays()
		config.GetLastCacheGenerationTimestamp()
		config.GetEntityPageSizes()
		config.GetMinimumPoWStrengths()
		config.GetPoWBailoutTimeSeconds()
		config.GetCacheGenerationIntervalHours()
		config.GetCacheDurationHours()
		config.GetClientVersionMajor()
		config.GetClientVersionMinor()
		config.GetClientVersionPatch()
		config.GetClientName()
		config.GetProtocolVersionMajor()
		config.GetProtocolVersionMinor()
		config.GetPOSTResponseExpiryMinutes()
		config.GetPOSTResponseIneligibilityMinutes()
		config.GetConnectionTimeout()
		config.GetTCPConnectTimeout()
		config.GetTLSHandshakeTimeout()
		config.GetPingerPageSize()
		config.GetOnlineAddressFinderPageSize()
		config.GetDispatchExclusionExpiryForLiveAddress()
		config.GetDispatchExclusionExpiryForStaticAddress()
		config.GetLoggingLevel()
		config.GetExternalIp()
		config.GetExternalPort()
		config.GetLastStaticAddressConnectionTimestamp()
		config.GetLastLiveAddressConnectionTimestamp()
		config.GetServingSubprotocols()
		config.GetDbEngine()
		config.GetDbIp()
		config.GetDbPort()
		config.GetDbPassword()
		config.GetMetricsLevel()
		config.GetMetricsToken()
		config.GetBackendKeyPair()
		config.GetMarshaledBackendPublicKey() // location sensitive, needs to happen after getbackendkeypair
		config.GetNodeId()                    // location sensitive, needs to happen after getbackendkeypair
		config.GetMaxInboundPageSizeKb()
		config.GetNeighbourCount()
		config.GetMaxAddressTableSize()
		config.GetMaxInboundConns()
		config.GetMaxOutboundConns()
		config.GetMaxPingConns()
		config.GetMaxDbSizeMb()
		config.GetVotesMemoryDays()
		config.GetEventHorizonTimestamp()
		config.GetLastBootstrapAddressConnectionTimestamp()
		config.GetBootstrapAfterOfflineMinutes()
		config.GetSOCKS5ProxyEnabled()
		config.GetSOCKS5ProxyAddress()
		config.GetSOCKS5ProxyUsername()
		config.GetSOCKS5ProxyPassword()
		config.GetNodeType()
		config.GetBackendAPIPublic()
		config.GetAdminFrontendAddress()
		config.GetAdminFrontendPublicKey()
		config.GetGRPCServiceTimeout()
		config.GetSQLiteDBLocation()
		config.GetDeclineInboundReverseRequests()
		config.GetPreventOutboundReverseRequests()
		config.GetCacheGenerationInterval()
	}
}

/*
Commit saves the file to memory. This is usually called after a Set operation.
*/
func (config *BackendConfig) Commit() error {
	if Btc.PermConfigReadOnly {
		return nil
	}
	Btc.ConfigMutex.Lock()
	defer Btc.ConfigMutex.Unlock()
	confAsByte, err3 := json.MarshalIndent(config, "", "    ")
	if err3 != nil {
		log.Fatal(fmt.Sprintf("JSON marshaler encountered an error while marshaling this config into JSON. Config: %#v, Error: %#v", config, err3))
	}
	configDirs := cdir.New(Btc.OrgIdentifier, Btc.AppIdentifier)
	folders := configDirs.QueryFolders(cdir.Global)
	toolbox.CreatePath(filepath.Join(folders[0].Path, "backend"))
	writeAheadPath := filepath.Join(folders[0].Path, "backend", "backend_config_writeahead.json")
	targetPath := filepath.Join(folders[0].Path, "backend", "backend_config.json")
	err := ioutil.WriteFile(writeAheadPath, confAsByte, 0755)
	if err != nil {
		return err
	}
	err2 := os.Rename(writeAheadPath, targetPath)
	// ^ Rename is atomic in UNIX, should be in Windows, as well. This is useful because if the user's app crashes at the right exact moment while committing a config change, if you don't have a readahead, you can end up with an empty file in the case the process crashes before inserting the bytes but after cleaning out the old ones to ready it for insertion.
	if err2 != nil {
		return err2
	}
	return nil
}

// Cycle commits the whole struct into memory, generating fields in JSON that were newly added.
func (config *BackendConfig) Cycle() error {
	err := config.Commit()
	if err != nil {
		return err
	}
	return nil
}

// The default base size is 1x (The thread size). At the base size, a page gets 100 entries.
func (config *BackendConfig) setDefaultEntityPageSizes() {
	var eps EntityPageSizes
	eps.Boards = defaultBoardsPageSize
	eps.Threads = defaultThreadsPageSize
	eps.Posts = defaultPostsPageSize
	eps.Votes = defaultVotesPageSize
	eps.Keys = defaultKeysPageSize
	eps.Truststates = defaultTruststatesPageSize
	eps.Addresses = defaultAddressesPageSize

	eps.BoardIndexes = defaultBoardIndexesPageSize
	eps.ThreadIndexes = defaultThreadIndexesPageSize
	eps.PostIndexes = defaultPostIndexesPageSize
	eps.VoteIndexes = defaultVoteIndexesPageSize
	eps.KeyIndexes = defaultKeyIndexesPageSize
	eps.TruststateIndexes = defaultTruststateIndexesPageSize
	eps.AddressIndexes = defaultAddressIndexesPageSize

	eps.BoardManifests = defaultBoardManifestsPageSize
	eps.ThreadManifests = defaultThreadManifestsPageSize
	eps.PostManifests = defaultPostManifestsPageSize
	eps.VoteManifests = defaultVoteManifestsPageSize
	eps.KeyManifests = defaultKeyManifestsPageSize
	eps.TruststateManifests = defaultTruststateManifestsPageSize
	eps.AddressManifests = defaultAddressManifestsPageSize
	config.SetEntityPageSizes(eps)
}

// Methods that have no specific direct backing in the config store but calculated from it.

var protv string

// GetProtURLVersion returns the protocol version to be used in the URL scheme. If the major version of mim is 1, the url scheme will be based on v0.
func (config *BackendConfig) GetProtURLVersion() string {
	if len(protv) > 0 {
		return protv
	}
	config.InitCheck()
	switch config.GetProtocolVersionMajor() {
	case 1:
		return "v0"
	default:
		return "v0"
	}
	return "v0"
}

// ===========================================

// 2) FRONTEND

/*
Frontend configuration.

# UserKeyPair
# MarshaledUserPublicKey
This is the key of the user residing in this frontend. Every user has their own frontend, but backends can be shared. Every device is their own frontend.

# FrontendKeyPair
# MarshaledFrontendPublicKey
This is the 'admin' key for either the local backend, or if this frontend is the maintainer of a remote backend, the key used to be authenticated with that backend. If multiple people are maintaining a backend, they can share this key between them without sharing their user accounts (their user keys).

# DehydratedLocalUserKeyEntity
This is the actual Key entity that the localUser owns (not just the PK, but also the actual entity, with its username, creation, etc.) This is a dehydrated (jsonified) copy of the entity. This is useful when you move your frontend to a different backend - it can push it to the new backend if it doesn't have it already.

## PoWBailoutTimeSeconds
How long does it take before a PoW timestamp is marked unattainable by the local computer. This is to make sure that the app doesn't keep attempting forever for an unattainably strong PoW it attempted to generate.
*/

// Frontend config base
type FrontendConfig struct {
	UserKeyPair                             string
	MarshaledUserPublicKey                  string
	Initialised                             bool   // False by default, init to set true
	MetricsLevel                            uint8  // 0: no metrics transmitted
	MetricsToken                            string // If metrics level is not zero, metrics token is the anonymous identifier for the metrics server. Resetting this to 0 makes this node behave like a new node as far as metrics go, but if you don't want metrics to be collected, you can set it through the application or set the metrics level to zero in the JSON settings file.
	LoggingLevel                            uint
	FrontendKeyPair                         string
	MarshaledFrontendPublicKey              string
	ExternalIp                              string
	ExternalIpType                          uint8
	FrontendAPIPort                         uint16
	UserDirectory                           string
	BackendAPIAddress                       string // This is IP - but it can also be a URL. That's why it's called an address and not an IP outright.
	BackendAPIPort                          uint16
	ClientAPIAddress                        string // almost always 127.0.0.1 because frontend and client usually run in the same machine.
	ClientPort                              uint16
	GRPCServiceTimeout                      time.Duration
	UserRelations                           UserRelations    // e.g. Local user's followed, mademod users
	ContentRelations                        ContentRelations // e.g. Local user's subbed boards, threads
	NetworkHeadDays                         uint             // 14
	NetworkMemoryDays                       uint             // 180
	LocalMemoryDays                         uint             // 180
	ThresholdForElectionValidityPercent     uint             // 5%
	ThresholdForElectionWinPercent          uint             // 51%
	BloomFilterSize                         uint             // 10000
	BloomFilterFalsePositiveRatePercent     uint             // 50%
	MinimumVoteThresholdForElectionValidity uint             // 100
	DehydratedLocalUserKeyEntity            string
	MinimumPoWStrengths                     MinimumPoWStrengths //
	PoWBailoutTimeSeconds                   uint                // 30
	OnboardComplete                         bool
	SFWListDisabled                         bool
	ModModeEnabled                          bool
	KvStoreRetentionDays                    uint
	LocalDevBackendEnabled                  bool
	LocalDevBackendDirectory                string
	LastKnownClientVersion                  string
	ExternalContentAutoloadDisabled         bool
}

// Init check gate

func (config *FrontendConfig) InitCheck() {
	if !config.Initialised {
		log.Fatal(fmt.Sprintf("You've attempted to access a config before it was initialised. Trace: %v\n", toolbox.DumpStack()))
	}
}

// Getters and setters

// Getters

// func (config *FrontendConfig) GetUserKeyPair() *ed25519.PrivateKey {
// 	config.InitCheck()
// 	keyPair := ed25519.PrivateKey([]byte(config.UserKeyPair))
// 	return &keyPair
// }

func (config *FrontendConfig) GetUserKeyPair() *ed25519.PrivateKey {
	config.InitCheck()
	keyPair, err := signaturing.UnmarshalPrivateKey(config.UserKeyPair)
	if err != nil {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.UserKeyPair) + " Trace: " + toolbox.Trace() + "Error: " + err.Error()))
	}
	return &keyPair
}

func (config *FrontendConfig) GetMarshaledUserPublicKey() string {
	config.InitCheck()
	return config.MarshaledUserPublicKey
}

func (config *FrontendConfig) GetInitialised() bool {
	return config.Initialised
}

func (config *FrontendConfig) GetMetricsLevel() uint8 {
	config.InitCheck()
	if config.MetricsLevel == 0 || config.MetricsLevel == 1 { // 0: no metrics, 1: anonymous metrics
		return config.MetricsLevel
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MetricsLevel) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetMetricsToken() string {
	config.InitCheck()
	if len(config.MetricsToken) < 65 &&
		len(config.MetricsToken) >= 0 {
		return config.MetricsToken
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MetricsToken) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetLoggingLevel() int {
	config.InitCheck()
	if config.LoggingLevel < toolbox.MaxInt32 { // can be zero
		return int(config.LoggingLevel)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LoggingLevel) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetFrontendKeyPair() *ed25519.PrivateKey {
	config.InitCheck()
	keyPair, err := signaturing.UnmarshalPrivateKey(config.FrontendKeyPair)
	if err != nil {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.GetFrontendKeyPair) + " Trace: " + toolbox.Trace() + "Error: " + err.Error()))
	}
	return &keyPair
}

func (config *FrontendConfig) GetMarshaledFrontendPublicKey() string {
	config.InitCheck()
	return config.MarshaledFrontendPublicKey
}

func (config *FrontendConfig) GetExternalIp() string {
	config.InitCheck()
	if len(config.ExternalIp) < maxLocationSize &&
		len(config.ExternalIp) > 0 {
		return config.ExternalIp
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ExternalIp) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}
func (config *FrontendConfig) GetExternalIpType() uint8 {
	config.InitCheck()
	if config.ExternalIpType == 6 || config.ExternalIpType == 4 || config.ExternalIpType == 3 { // 6: ipv6, 4: ipv4, 3: URL (useful in static nodes)
		return config.ExternalIpType
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ExternalIpType) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *FrontendConfig) GetFrontendAPIPort() uint16 {
	config.InitCheck()
	if config.FrontendAPIPort < toolbox.MaxUint16 && config.FrontendAPIPort > 0 {
		return config.FrontendAPIPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.FrontendAPIPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetUserDirectory() string {
	config.InitCheck()
	if len(config.UserDirectory) < toolbox.MaxUint16 &&
		len(config.UserDirectory) > 0 {
		return config.UserDirectory
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.UserDirectory) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetBackendAPIAddress() string {
	config.InitCheck()
	if len(config.BackendAPIAddress) < maxLocationSize &&
		len(config.BackendAPIAddress) > 0 {
		return config.BackendAPIAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BackendAPIAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetBackendAPIPort() uint16 {
	config.InitCheck()
	if config.BackendAPIPort < toolbox.MaxUint16 && config.BackendAPIPort > 0 {
		return config.BackendAPIPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BackendAPIPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetClientAPIAddress() string {
	config.InitCheck()
	if len(config.ClientAPIAddress) < maxLocationSize &&
		len(config.ClientAPIAddress) > 0 {
		return config.ClientAPIAddress
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientAPIAddress) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetClientPort() uint16 {
	config.InitCheck()
	if config.ClientPort < toolbox.MaxUint16 && config.ClientPort > 0 {
		return config.ClientPort
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ClientPort) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetGRPCServiceTimeout() time.Duration {
	config.InitCheck()
	if config.GRPCServiceTimeout >= 1*time.Second { // Any value under is probably an attack.
		return config.GRPCServiceTimeout
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.GRPCServiceTimeout) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return time.Duration(0)
}

func (config *FrontendConfig) GetLocalMemoryDays() int {
	config.InitCheck()
	if config.LocalMemoryDays < night &&
		config.LocalMemoryDays > 0 {
		return int(config.LocalMemoryDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LocalMemoryDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *FrontendConfig) GetNetworkMemoryDays() int {
	config.InitCheck()
	if config.NetworkMemoryDays < night &&
		config.NetworkMemoryDays > 0 {
		return int(config.NetworkMemoryDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NetworkMemoryDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *FrontendConfig) GetNetworkHeadDays() int {
	config.InitCheck()
	if config.NetworkHeadDays < night &&
		config.NetworkHeadDays > 0 {
		return int(config.NetworkHeadDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.NetworkHeadDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}
func (config *FrontendConfig) GetThresholdForElectionValidityPercent() int {
	config.InitCheck()
	if config.ThresholdForElectionValidityPercent < 100 &&
		config.ThresholdForElectionValidityPercent > 0 {
		return int(config.ThresholdForElectionValidityPercent)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ThresholdForElectionValidityPercent) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetThresholdForElectionWinPercent() int {
	config.InitCheck()
	if config.ThresholdForElectionWinPercent < 100 &&
		config.ThresholdForElectionWinPercent > 50 {
		return int(config.ThresholdForElectionWinPercent)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ThresholdForElectionWinPercent) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetBloomFilterSize() int {
	config.InitCheck()
	if config.BloomFilterSize < 10000000 &&
		config.BloomFilterSize > 100 {
		return int(config.BloomFilterSize)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BloomFilterSize) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetBloomFilterFalsePositiveRatePercent() int {
	config.InitCheck()
	if config.BloomFilterFalsePositiveRatePercent < 100 &&
		config.BloomFilterFalsePositiveRatePercent > 0 {
		return int(config.BloomFilterFalsePositiveRatePercent)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.BloomFilterFalsePositiveRatePercent) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetMinimumVoteThresholdForElectionValidity() int {
	config.InitCheck()
	if config.MinimumVoteThresholdForElectionValidity < 10000000 &&
		config.MinimumVoteThresholdForElectionValidity > 0 {
		return int(config.MinimumVoteThresholdForElectionValidity)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.MinimumVoteThresholdForElectionValidity) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetUserRelations() UserRelations {
	config.InitCheck()
	if config.UserRelations.Initialised {
		return config.UserRelations
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.UserRelations) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return UserRelations{}
}

func (config *FrontendConfig) GetContentRelations() ContentRelations {
	config.InitCheck()
	if config.ContentRelations.Initialised {
		return config.ContentRelations
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.ContentRelations) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ContentRelations{}
}

func (config *FrontendConfig) GetDehydratedLocalUserKeyEntity() string {
	config.InitCheck()
	if uint(len(config.DehydratedLocalUserKeyEntity)) < toolbox.MaxUint32 {
		return config.DehydratedLocalUserKeyEntity
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.DehydratedLocalUserKeyEntity) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetMinimumPoWStrengths() MinimumPoWStrengths {
	config.InitCheck()
	if config.MinimumPoWStrengths.Board < maxPOWStrength &&
		config.MinimumPoWStrengths.Board > 0 &&
		config.MinimumPoWStrengths.BoardUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.BoardUpdate > 0 &&
		config.MinimumPoWStrengths.Thread < maxPOWStrength &&
		config.MinimumPoWStrengths.Thread > 0 &&
		config.MinimumPoWStrengths.ThreadUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.ThreadUpdate > 0 &&
		config.MinimumPoWStrengths.Post < maxPOWStrength &&
		config.MinimumPoWStrengths.Post > 0 &&
		config.MinimumPoWStrengths.PostUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.PostUpdate > 0 &&
		config.MinimumPoWStrengths.Vote < maxPOWStrength &&
		config.MinimumPoWStrengths.Vote > 0 &&
		config.MinimumPoWStrengths.VoteUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.VoteUpdate > 0 &&
		config.MinimumPoWStrengths.Key < maxPOWStrength &&
		config.MinimumPoWStrengths.Key > 0 &&
		config.MinimumPoWStrengths.KeyUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.KeyUpdate > 0 &&
		config.MinimumPoWStrengths.Truststate < maxPOWStrength &&
		config.MinimumPoWStrengths.Truststate > 0 &&
		config.MinimumPoWStrengths.TruststateUpdate < maxPOWStrength &&
		config.MinimumPoWStrengths.TruststateUpdate > 0 &&
		config.MinimumPoWStrengths.ApiResponse < maxPOWStrength &&
		config.MinimumPoWStrengths.ApiResponse > 0 {
		return config.MinimumPoWStrengths
	} else {
		log.Fatal(fmt.Sprintf("%#v", invalidDataError(config.MinimumPoWStrengths)) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return MinimumPoWStrengths{}
}

func (config *FrontendConfig) GetPoWBailoutTimeSeconds() int {
	config.InitCheck()
	if config.PoWBailoutTimeSeconds < maxPOWBailoutSeconds &&
		config.PoWBailoutTimeSeconds > 0 {
		return int(config.PoWBailoutTimeSeconds)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.PoWBailoutTimeSeconds) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetOnboardComplete() bool {
	config.InitCheck()
	return config.OnboardComplete
}

func (config *FrontendConfig) GetSFWListDisabled() bool {
	config.InitCheck()
	return config.SFWListDisabled
}

func (config *FrontendConfig) GetModModeEnabled() bool {
	config.InitCheck()
	return config.ModModeEnabled
}

func (config *FrontendConfig) GetKvStoreRetentionDays() int {
	config.InitCheck()
	if config.KvStoreRetentionDays < night &&
		config.KvStoreRetentionDays > 0 {
		return int(config.KvStoreRetentionDays)
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.KvStoreRetentionDays) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return 0
}

func (config *FrontendConfig) GetLocalDevBackendEnabled() bool {
	config.InitCheck()
	return config.LocalDevBackendEnabled
}

func (config *FrontendConfig) GetLocalDevBackendDirectory() string {
	config.InitCheck()
	if uint(len(config.LocalDevBackendDirectory)) < toolbox.MaxUint32 {
		return config.LocalDevBackendDirectory
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LocalDevBackendDirectory) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetLastKnownClientVersion() string {
	config.InitCheck()
	if uint(len(config.LastKnownClientVersion)) < toolbox.MaxUint32 {
		return config.LastKnownClientVersion
	} else {
		log.Fatal(invalidDataError(fmt.Sprintf("%#v", config.LastKnownClientVersion) + " Trace: " + toolbox.Trace()))
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return ""
}

func (config *FrontendConfig) GetExternalContentAutoloadDisabled() bool {
	config.InitCheck()
	return config.ExternalContentAutoloadDisabled
}

/*****************************************************************************/

// Setters

func (config *FrontendConfig) SetUserKeyPair(val *ed25519.PrivateKey) error {
	config.InitCheck()
	config.UserKeyPair = signaturing.MarshalPrivateKey(*val)
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

// The only way to set this is to set frontend key pair first.
func (config *FrontendConfig) SetMarshaledUserPublicKey(val *ed25519.PrivateKey) error {
	config.InitCheck()
	marshaledPk := signaturing.MarshalPublicKey(val.Public().(ed25519.PublicKey))
	config.MarshaledUserPublicKey = marshaledPk
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetInitialised(val bool) error {
	// No init check on this one, so we can start inserting data.
	config.Initialised = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetMetricsLevel(val int) error {
	config.InitCheck()
	if val == 0 || val == 1 {
		config.MetricsLevel = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetMetricsToken(val string) error {
	config.InitCheck()
	if len(val) >= 0 && len(val) < 65 {
		config.MetricsToken = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetLoggingLevel(val int) error {
	config.InitCheck()
	if val >= 0 {
		config.LoggingLevel = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetFrontendKeyPair(val *ed25519.PrivateKey) error {
	config.InitCheck()
	config.FrontendKeyPair = signaturing.MarshalPrivateKey(*val)
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

// The only way to set this is to set frontend key pair first.
func (config *FrontendConfig) SetMarshaledFrontendPublicKey(val *ed25519.PrivateKey) error {
	config.InitCheck()
	marshaledPk := signaturing.MarshalPublicKey(val.Public().(ed25519.PublicKey))
	config.MarshaledFrontendPublicKey = marshaledPk
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetExternalIp(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.ExternalIp = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetExternalIpType(val int) error {
	config.InitCheck()
	if val == 6 || val == 4 || val == 3 {
		config.ExternalIpType = uint8(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetFrontendAPIPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.FrontendAPIPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetUserDirectory(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < toolbox.MaxUint16 {
		config.UserDirectory = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetBackendAPIAddress(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.BackendAPIAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetBackendAPIPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.BackendAPIPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetClientAPIAddress(val string) error {
	config.InitCheck()
	if len(val) > 0 && len(val) < maxLocationSize {
		config.ClientAPIAddress = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetClientPort(val int) error {
	config.InitCheck()
	if val > 0 && val < toolbox.MaxUint16 {
		config.ClientPort = uint16(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetGRPCServiceTimeout(val time.Duration) error {
	config.InitCheck()
	if val >= 1*time.Second { // Any value under is probably an attack.
		config.GRPCServiceTimeout = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
}

func (config *FrontendConfig) SetLocalMemoryDays(val int) error {
	if val > 0 {
		config.InitCheck()
		config.LocalMemoryDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *FrontendConfig) SetNetworkMemoryDays(val int) error {
	config.InitCheck()
	if val > 0 {
		config.NetworkMemoryDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}
func (config *FrontendConfig) SetNetworkHeadDays(val int) error {
	config.InitCheck()
	if val > 0 {
		config.NetworkHeadDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetThresholdForElectionValidityPercent(val int) error {
	if val > 0 && val < 100 {
		config.InitCheck()
		config.ThresholdForElectionValidityPercent = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetThresholdForElectionWinPercent(val int) error {
	if val > 50 && val < 100 {
		config.InitCheck()
		config.ThresholdForElectionWinPercent = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetBloomFilterSize(val int) error {
	if val > 100 && val < 10000000 {
		config.InitCheck()
		config.BloomFilterSize = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetBloomFilterFalsePositiveRatePercent(val int) error {
	if val > 0 && val < 100 {
		config.InitCheck()
		config.BloomFilterFalsePositiveRatePercent = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetMinimumVoteThresholdForElectionValidity(val int) error {
	if val > 0 && val < 10000000 {
		config.InitCheck()
		config.MinimumVoteThresholdForElectionValidity = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetUserRelations(val UserRelations) error {
	if config.UserRelations.Initialised {
		config.InitCheck()
		config.UserRelations = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetContentRelations(val ContentRelations) error {
	if config.ContentRelations.Initialised {
		config.InitCheck()
		config.ContentRelations = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetDehydratedLocalUserKeyEntity(val string) error {
	config.InitCheck()
	if uint(len(val)) < toolbox.MaxUint32 {
		config.DehydratedLocalUserKeyEntity = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetMinimumPoWStrengths(powStr int) error {
	config.InitCheck()
	var mps MinimumPoWStrengths
	if powStr > 6 && powStr < maxPOWStrength {
		mps.Board = powStr
		mps.BoardUpdate = powStr
		mps.Thread = powStr
		mps.ThreadUpdate = powStr
		mps.Post = powStr
		mps.PostUpdate = powStr
		mps.Vote = powStr
		mps.VoteUpdate = powStr
		mps.Key = powStr
		mps.KeyUpdate = powStr
		mps.Truststate = powStr
		mps.TruststateUpdate = powStr
		mps.ApiResponse = powStr - 3
		// ^ ApiResponse PoW strength is a little softer because remotes need it to make POST requests to each other. Temp value, subject to solidification or change in the future.
		config.MinimumPoWStrengths = mps
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		} else {
			return invalidDataError(fmt.Sprintf("%#v", powStr) + " Trace: " + toolbox.Trace())
		}
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetPoWBailoutTimeSeconds(val int) error {
	config.InitCheck()
	if val > 0 {
		config.PoWBailoutTimeSeconds = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetOnboardComplete(val bool) error {
	config.InitCheck()
	config.OnboardComplete = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetSFWListDisabled(val bool) error {
	config.InitCheck()
	config.SFWListDisabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetModModeEnabled(val bool) error {
	config.InitCheck()
	config.ModModeEnabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetKvStoreRetentionDays(val int) error {
	if val > 0 {
		config.InitCheck()
		config.KvStoreRetentionDays = uint(val)
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetLocalDevBackendEnabled(val bool) error {
	config.InitCheck()
	config.LocalDevBackendEnabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

func (config *FrontendConfig) SetLocalDevBackendDirectory(val string) error {
	config.InitCheck()
	if uint(len(val)) < toolbox.MaxUint32 {
		config.LocalDevBackendDirectory = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetLastKnownClientVersion(val string) error {
	config.InitCheck()
	if uint(len(val)) < toolbox.MaxUint32 {
		config.LastKnownClientVersion = val
		commitErr := config.Commit()
		if commitErr != nil {
			return commitErr
		}
		return nil
	} else {
		return invalidDataError(fmt.Sprintf("%#v", val) + " Trace: " + toolbox.Trace())
	}
	log.Fatal("This should never happen." + toolbox.Trace())
	return nil
}

func (config *FrontendConfig) SetExternalContentAutoloadDisabled(val bool) error {
	config.InitCheck()
	config.ExternalContentAutoloadDisabled = val
	commitErr := config.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}

/*****************************************************************************/

// Frontend config methods

func (config *FrontendConfig) BlankCheck() {
	// SetInitialised needs to be true before we can make any changes.
	if !config.Initialised {
		config.SetInitialised(true)
	}
	if len(config.UserKeyPair) == 0 {
		privKey, _ := signaturing.CreateKeyPair()
		config.SetUserKeyPair(privKey)
	}
	// This needs to be after Frontend user key pair generation.
	if len(config.MarshaledUserPublicKey) == 0 {
		config.SetMarshaledUserPublicKey(config.GetUserKeyPair())
	}
	// ::MetricsLevel: can be zero, no need to blank check.
	// ::MetricsToken: can be zero, no need to blank check.
	if len(config.FrontendKeyPair) == 0 {
		privKey, _ := signaturing.CreateKeyPair()
		config.SetFrontendKeyPair(privKey)
	}
	// This needs to be after SpawnerFE  key pair generation.
	if len(config.MarshaledFrontendPublicKey) == 0 {
		config.SetMarshaledFrontendPublicKey(config.GetFrontendKeyPair())
	}
	if len(config.ExternalIp) == 0 {
		config.SetExternalIp(defaultFrontendExternalIp)
	}
	if config.ExternalIpType == 0 {
		config.SetExternalIpType(defaultFrontendExternalIpType)
	}
	if config.FrontendAPIPort == 0 {
		config.SetFrontendAPIPort(defaultFrontendAPIPort)
	}
	if len(config.UserDirectory) == 0 {
		config.SetUserDirectory(cdir.New(Ftc.OrgIdentifier, Ftc.AppIdentifier).QueryFolders(cdir.Global)[0].Path)
	}
	if len(config.BackendAPIAddress) == 0 {
		config.SetBackendAPIAddress(defaultBackendAPIAddress)
	}
	if config.BackendAPIPort == 0 {
		config.SetBackendAPIPort(defaultBackendAPIPort)
	}
	if len(config.ClientAPIAddress) == 0 {
		config.SetClientAPIAddress(defaultClientAPIAddress)
	}
	if config.ClientPort == 0 {
		config.SetClientPort(defaultClientPort)
	}
	if config.GRPCServiceTimeout == 0 {
		config.SetGRPCServiceTimeout(defaultGRPCServiceTimeout)
	}
	if config.NetworkHeadDays == 0 {
		config.SetNetworkHeadDays(defaultNetworkHeadDays)
	}
	if config.NetworkMemoryDays == 0 {
		config.SetNetworkMemoryDays(defaultNetworkMemoryDays)
	}
	if config.LocalMemoryDays == 0 {
		config.SetLocalMemoryDays(defaultLocalMemoryDays)
	}
	if config.ThresholdForElectionValidityPercent == 0 {
		config.SetThresholdForElectionValidityPercent(defaultThresholdForElectionValidityPercent)
	}
	if config.ThresholdForElectionWinPercent == 0 {
		config.SetThresholdForElectionWinPercent(defaultThresholdForElectionWinPercent)
	}
	if config.BloomFilterSize == 0 {
		config.SetBloomFilterSize(defaultBloomFilterSize)
	}
	if config.BloomFilterFalsePositiveRatePercent == 0 {
		config.SetBloomFilterFalsePositiveRatePercent(defaultBloomFilterFalsePositiveRatePercent)
	}
	if config.MinimumVoteThresholdForElectionValidity == 0 {
		config.SetMinimumVoteThresholdForElectionValidity(defaultMinimumVoteThresholdForElectionValidity)
	}
	if config.UserRelations.Initialised == false {
		config.UserRelations.Init()
		config.SetUserRelations(config.UserRelations)
	}
	if config.ContentRelations.Initialised == false {
		config.ContentRelations.Init()
		config.SetContentRelations(config.ContentRelations)
	}
	// ::DehydratedLocalUserKeyEntity: can be empty, no need to blank check.
	if config.MinimumPoWStrengths.Board == 0 ||
		config.MinimumPoWStrengths.BoardUpdate == 0 ||
		config.MinimumPoWStrengths.Thread == 0 ||
		config.MinimumPoWStrengths.ThreadUpdate == 0 ||
		config.MinimumPoWStrengths.Post == 0 ||
		config.MinimumPoWStrengths.PostUpdate == 0 ||
		config.MinimumPoWStrengths.Vote == 0 ||
		config.MinimumPoWStrengths.VoteUpdate == 0 ||
		config.MinimumPoWStrengths.Key == 0 ||
		config.MinimumPoWStrengths.KeyUpdate == 0 ||
		config.MinimumPoWStrengths.Truststate == 0 ||
		config.MinimumPoWStrengths.TruststateUpdate == 0 ||
		config.MinimumPoWStrengths.ApiResponse == 0 {
		config.SetMinimumPoWStrengths(defaultPowStrength)
	}
	if config.PoWBailoutTimeSeconds == 0 {
		config.SetPoWBailoutTimeSeconds(defaultPoWBailoutTimeSeconds)
	}
	// ::OnboardComplete: can be false, no need to blank check.
	// ::SFWListDisabled: can be false, no need to blank check.
	// ::ModModeEnabled: can be false, no need to blank check.
	// ::LocalDevBackendEnabled: can be false, no need to blank check.
	if config.KvStoreRetentionDays == 0 {
		config.SetKvStoreRetentionDays(defaultKvStoreRetentionDays)
	}
	if len(config.LocalDevBackendDirectory) == 0 {
		config.SetLocalDevBackendDirectory(defaultLocalDevBackendDirectory)
	}
	// ::LastKnownClientVersion: can be false, no need to blank check.
	// ::ExternalContentAutoloadDisabled: can be false, no need to blank check.

}
func (config *FrontendConfig) SanityCheck() {
	if !config.GetInitialised() {
		log.Fatal("Frontend configuration is not initialised. Please initialise it before use.")
	} else {
		config.GetUserKeyPair()
		config.GetMetricsLevel()
		config.GetMetricsToken()
		config.GetFrontendKeyPair()
		config.GetExternalIp()
		config.GetFrontendAPIPort()
		config.GetBackendAPIAddress()
		config.GetBackendAPIPort()
		config.GetClientAPIAddress()
		config.GetClientPort()
		config.GetGRPCServiceTimeout()
		config.GetLocalMemoryDays()
		config.GetNetworkMemoryDays()
		config.GetNetworkHeadDays()
		config.GetThresholdForElectionValidityPercent()
		config.GetThresholdForElectionWinPercent()
		config.GetBloomFilterSize()
		config.GetBloomFilterFalsePositiveRatePercent()
		config.GetMinimumVoteThresholdForElectionValidity()
		config.GetMinimumPoWStrengths()
		config.GetPoWBailoutTimeSeconds()
		config.GetKvStoreRetentionDays()
		config.GetLocalDevBackendDirectory()
	}
}

/*
Commit saves the file to memory. This is usually called after a Set operation.
*/
func (config *FrontendConfig) Commit() error {
	if Ftc.PermConfigReadOnly {
		return nil
	}
	Ftc.ConfigMutex.Lock()
	defer Ftc.ConfigMutex.Unlock()
	confAsByte, err3 := json.MarshalIndent(config, "", "    ")
	if err3 != nil {
		log.Fatal(fmt.Sprintf("JSON marshaler encountered an error while marshaling this config into JSON. Config: %#v, Error: %#v", config, err3))
	}
	configDirs := cdir.New(Ftc.OrgIdentifier, Ftc.AppIdentifier)
	folders := configDirs.QueryFolders(cdir.Global)
	toolbox.CreatePath(filepath.Join(folders[0].Path, "frontend"))
	writeAheadPath := filepath.Join(folders[0].Path, "frontend", "frontend_config_writeahead.json")
	targetPath := filepath.Join(folders[0].Path, "frontend", "frontend_config.json")
	err := ioutil.WriteFile(writeAheadPath, confAsByte, 0755)
	if err != nil {
		return err
	}
	err2 := os.Rename(writeAheadPath, targetPath)
	if err2 != nil {
		return err2
	}
	return nil
}

// Cycle commits the whole struct into memory, generating fields in JSON that were newly added.
func (config *FrontendConfig) Cycle() error {
	err := config.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*****************************************************************************/

// 3) CONFIG METHODS

/*
EstablishBackendConfig establishes the connection with the config file, and makes it available as an object to the rest of the application.
*/
func EstablishBackendConfig() (*BackendConfig, error) {
	// var config BackendConfig
	configDirs := cdir.New(Btc.OrgIdentifier, Btc.AppIdentifier)
	folder := configDirs.QueryFolderContainsFile("backend/backend_config.json")
	if folder != nil {
		configJson, _ := folder.ReadFile("backend/backend_config.json")
		err := json.Unmarshal(configJson, &bc)
		if err != nil || fmt.Sprintf("%#v", string(configJson)) == "\"{}\"" {
			return &bc, errors.New(fmt.Sprintf("Back-end configuration file is corrupted. Please fix the configuration file, or delete it. If deleted a new configuration will be generated with default values. Error: %#v, ConfigJson: %#v", err, string(configJson)))
		}
	}
	// Folder is nil - the configuration file in question does not exist. Ask to create.
	bc.BlankCheck()
	bc.SanityCheck()
	return &bc, nil
}

/*
EstablishFrontendConfig establishes the connection with the config file, and makes it available as an object to the rest of the application.
*/
func EstablishFrontendConfig() (*FrontendConfig, error) {
	// var config FrontendConfig
	configDirs := cdir.New(Ftc.OrgIdentifier, Ftc.AppIdentifier)
	// _ = os.Mkdir(configDirs, mode)
	folder := configDirs.QueryFolderContainsFile("frontend/frontend_config.json")
	if folder != nil {
		configJson, _ := folder.ReadFile("frontend/frontend_config.json")
		err := json.Unmarshal(configJson, &fc)
		if err != nil || fmt.Sprintf("%#v", string(configJson)) == "\"{}\"" {
			return &fc, errors.New(fmt.Sprintf("Front-end configuration file is corrupted. Please fix the configuration file, or delete it. If deleted a new configuration will be generated with default values. Error: %#v, ConfigJson: %#v", err, string(configJson)))
		}
	}
	// Folder is nil - the configuration file in question does not exist. Ask to create.
	fc.BlankCheck()
	fc.SanityCheck()
	return &fc, nil
}
