// API > Structs
// This file provides the struct definitions for the protocol. This is what should be arriving from the network, and what should be sent over to other nodes.

package api

import (
	"database/sql/driver"
	// "fmt"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"

	"golang.org/x/crypto/ed25519"

	// "github.com/davecgh/go-spew/spew"
	"aether-core/aether/services/globals"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"
)

// Structs for the entity types. There are 7 types. Board, Thread, Post, Vote, Key, Address, Truststate.

// Low-level types

// type Fingerprint [64]byte // 64 char ASCII
type Fingerprint string // 64 char ASCII
type Nonce string       // max 64 char ASCII
type Timestamp int64    // UNIX Timestamp
// type ProofOfWork [1024]byte
type ProofOfWork string // temp
// type Signature [512]byte
type Signature string // temp
type Location string

func (t Timestamp) Humanise() string {
	if t != 0 {
		return fmt.Sprintf("%s (%d)", time.Unix(int64(t), 0).Format(time.Stamp), t)
	} else {
		return fmt.Sprint("Blank")
	}
}

func (f Fingerprint) Value() (driver.Value, error) {
	return string(f), nil
}

func (f *Fingerprint) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*f = Fingerprint(stringVal)
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *Timestamp) Scan(value interface{}) error {
	numVal := value.(int64)
	*t = Timestamp(numVal)
	return nil
}

func (p ProofOfWork) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *ProofOfWork) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*p = ProofOfWork(stringVal)
	return nil
}

func (s Signature) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Signature) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*s = Signature(stringVal)
	return nil
}

func (l Location) Value() (driver.Value, error) {
	return string(l), nil
}

func (l *Location) Scan(value interface{}) error {
	stringVal := string(value.([]uint8))
	*l = Location(stringVal)
	return nil
}

// Basic properties

type ProvableFieldSet struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	Creation    Timestamp   `json:"creation"`
	ProofOfWork ProofOfWork `json:"proof_of_work"`
	Signature   Signature   `json:"signature"`
	Verified    bool        `json:"-"`
}

type UpdateableFieldSet struct { // Common set of properties for all objects that are updateable.
	LastUpdate        Timestamp   `json:"last_update"`
	UpdateProofOfWork ProofOfWork `json:"update_proof_of_work"`
	UpdateSignature   Signature   `json:"update_signature"`
}

// Subentities

type BoardOwner struct {
	KeyFingerprint Fingerprint `json:"key_fingerprint"` // Fingerprint of the key the ownership is associated to.
	Expiry         Timestamp   `json:"expiry"`          // When the ownership expires.
	Level          uint8       `json:"level"`           // mod(1)
}

type Subprotocol struct {
	Name              string   `json:"name"` //2-16 chars
	VersionMajor      uint8    `json:"version_major"`
	VersionMinor      uint16   `json:"version_minor"`
	SupportedEntities []string `json:"supported_entities"`
}

type Protocol struct {
	VersionMajor uint8         `json:"version_major"`
	VersionMinor uint16        `json:"version_minor"`
	Subprotocols []Subprotocol `json:"subprotocols"`
}

type Client struct {
	VersionMajor uint8  `json:"version_major"`
	VersionMinor uint16 `json:"version_minor"`
	VersionPatch uint16 `json:"version_patch"`
	ClientName   string `json:"name"` // Max 128
}

// Entities

type Board struct { // Mutables: BoardOwners, Description, Meta
	ProvableFieldSet
	Name           string       `json:"name"`         // Max 255 char unicode
	BoardOwners    []BoardOwner `json:"board_owners"` // max 128 owners
	Description    string       `json:"description"`  // Max 65535 char unicode
	Owner          Fingerprint  `json:"owner"`
	OwnerPublicKey string       `json:"owner_publickey"`
	EntityVersion  int          `json:"entity_version"`
	Language       string       `json:"language"`
	Meta           string       `json:"meta"` // This is the dynamic JSON field
	RealmId        Fingerprint  `json:"realm_id"`
	EncrContent    string       `json:"encrcontent"`
	UpdateableFieldSet
}

type Thread struct { // Mutables: Body, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Name           string      `json:"name"`
	Body           string      `json:"body"`
	Link           string      `json:"link"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Post struct { // Mutables: Body, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Thread         Fingerprint `json:"thread"`
	Parent         Fingerprint `json:"parent"`
	Body           string      `json:"body"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Vote struct { // Mutables: Type, Meta
	ProvableFieldSet
	Board          Fingerprint `json:"board"`
	Thread         Fingerprint `json:"thread"`
	Target         Fingerprint `json:"target"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	TypeClass      int         `json:"typeclass"`
	Type           int         `json:"type"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Key struct { // Mutables: Expiry, Info, Meta
	ProvableFieldSet
	Type          string      `json:"type"`
	Key           string      `json:"key"`
	Expiry        Timestamp   `json:"expiry"`
	Name          string      `json:"name"`
	Info          string      `json:"info"`
	EntityVersion int         `json:"entity_version"`
	Meta          string      `json:"meta"`
	RealmId       Fingerprint `json:"realm_id"`
	EncrContent   string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Truststate struct { // Mutables: Type, Expiry, Meta (Domain is immutable)
	ProvableFieldSet
	Target         Fingerprint `json:"target"`
	Owner          Fingerprint `json:"owner"`
	OwnerPublicKey string      `json:"owner_publickey"`
	TypeClass      int         `json:"typeclass"`
	Type           int         `json:"type"`
	Domain         Fingerprint `json:"domain"`
	Expiry         Timestamp   `json:"expiry"`
	EntityVersion  int         `json:"entity_version"`
	Meta           string      `json:"meta"`
	RealmId        Fingerprint `json:"realm_id"`
	EncrContent    string      `json:"encrcontent"`
	UpdateableFieldSet
}

type Address struct { // Mutables: None
	Location           Location    `json:"location"`
	Sublocation        Location    `json:"sublocation"`
	LocationType       uint8       `json:"location_type"`
	Port               uint16      `json:"port"`
	Type               uint8       `json:"type"`
	LastSuccessfulPing Timestamp   `json:"-,omitempty"`
	LastSuccessfulSync Timestamp   `json:"-,omitempty"`
	Protocol           Protocol    `json:"protocol"`
	Client             Client      `json:"client"`
	EntityVersion      int         `json:"entity_version"`
	RealmId            Fingerprint `json:"realm_id"`
	Verified           bool        `json:"-,omitempty"` // This is normally part of the provable field set, but address is not provable, so provided here separately.
}

// Index Form Entities: These are index forms of the entities above.

type BoardIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Owner         Fingerprint `json:",omitempty"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type ThreadIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Owner         Fingerprint `json:",omitempty"`
	Board         Fingerprint `json:"board"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type PostIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Owner         Fingerprint `json:",omitempty"`
	Board         Fingerprint `json:"board"`
	Thread        Fingerprint `json:"thread"`
	Parent        Fingerprint `json:",omitempty"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type VoteIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Owner         Fingerprint `json:",omitempty"`
	Board         Fingerprint `json:"board"`
	Thread        Fingerprint `json:"thread"`
	Target        Fingerprint `json:"target"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type AddressIndex Address

type KeyIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

type TruststateIndex struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Owner         Fingerprint `json:",omitempty"`
	Target        Fingerprint `json:"target"`
	Creation      Timestamp   `json:"creation"`
	LastUpdate    Timestamp   `json:"last_update"`
	EntityVersion int         `json:"entity_version"`
	PageNumber    int         `json:"page_number"`
}

// Index interfaces

type ProvableIndex interface {
	GetFingerprint() Fingerprint
	GetEntityType() string
	GetCreation() Timestamp
	GetLastUpdate() Timestamp
	GetLastModified() Timestamp
	GetOwner() Fingerprint
	IsIndex() bool // We don't want other inter things that aren't indexes to satisfy this interface.
}

// Fingerprint accessors

func (entity *BoardIndex) GetFingerprint() Fingerprint      { return entity.Fingerprint }
func (entity *ThreadIndex) GetFingerprint() Fingerprint     { return entity.Fingerprint }
func (entity *PostIndex) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *VoteIndex) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *KeyIndex) GetFingerprint() Fingerprint        { return entity.Fingerprint }
func (entity *TruststateIndex) GetFingerprint() Fingerprint { return entity.Fingerprint }

// LastUpdate accessors

func (entity *BoardIndex) GetLastUpdate() Timestamp      { return entity.LastUpdate }
func (entity *ThreadIndex) GetLastUpdate() Timestamp     { return entity.LastUpdate }
func (entity *PostIndex) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *VoteIndex) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *KeyIndex) GetLastUpdate() Timestamp        { return entity.LastUpdate }
func (entity *TruststateIndex) GetLastUpdate() Timestamp { return entity.LastUpdate }

// Creation accessors

func (entity *BoardIndex) GetCreation() Timestamp      { return entity.Creation }
func (entity *ThreadIndex) GetCreation() Timestamp     { return entity.Creation }
func (entity *PostIndex) GetCreation() Timestamp       { return entity.Creation }
func (entity *VoteIndex) GetCreation() Timestamp       { return entity.Creation }
func (entity *KeyIndex) GetCreation() Timestamp        { return entity.Creation }
func (entity *TruststateIndex) GetCreation() Timestamp { return entity.Creation }

// EntityType accessors

func (entity *BoardIndex) GetEntityType() string      { return "board" }
func (entity *ThreadIndex) GetEntityType() string     { return "thread" }
func (entity *PostIndex) GetEntityType() string       { return "post" }
func (entity *VoteIndex) GetEntityType() string       { return "vote" }
func (entity *KeyIndex) GetEntityType() string        { return "key" }
func (entity *TruststateIndex) GetEntityType() string { return "truststate" }
func (entity *AddressIndex) GetEntityType() string    { return "address" }

// LastModified accessors (LM: the larger of creation / lastupdate)

func glmIndex(entity ProvableIndex) Timestamp {
	if entity.GetCreation() > entity.GetLastUpdate() {
		return entity.GetCreation()
	} else {
		return entity.GetLastUpdate()
	}
}

func (entity *BoardIndex) GetLastModified() Timestamp      { return glmIndex(entity) }
func (entity *ThreadIndex) GetLastModified() Timestamp     { return glmIndex(entity) }
func (entity *PostIndex) GetLastModified() Timestamp       { return glmIndex(entity) }
func (entity *VoteIndex) GetLastModified() Timestamp       { return glmIndex(entity) }
func (entity *KeyIndex) GetLastModified() Timestamp        { return glmIndex(entity) }
func (entity *TruststateIndex) GetLastModified() Timestamp { return glmIndex(entity) }

// IsIndex accessors

func (entity *BoardIndex) IsIndex() bool      { return true }
func (entity *ThreadIndex) IsIndex() bool     { return true }
func (entity *PostIndex) IsIndex() bool       { return true }
func (entity *VoteIndex) IsIndex() bool       { return true }
func (entity *KeyIndex) IsIndex() bool        { return true }
func (entity *TruststateIndex) IsIndex() bool { return true }

// GetOwner accessors

func (entity *BoardIndex) GetOwner() Fingerprint      { return entity.Owner }
func (entity *ThreadIndex) GetOwner() Fingerprint     { return entity.Owner }
func (entity *PostIndex) GetOwner() Fingerprint       { return entity.Owner }
func (entity *VoteIndex) GetOwner() Fingerprint       { return entity.Owner }
func (entity *KeyIndex) GetOwner() Fingerprint        { return entity.Fingerprint }
func (entity *TruststateIndex) GetOwner() Fingerprint { return entity.Owner }

// Response types

type Pagination struct {
	Pages       uint64 `json:"pages"`
	CurrentPage uint64 `json:"current_page"`
}

type Caching struct {
	Pregenerated    bool          `json:"pregenerated"`
	CurrentCacheUrl string        `json:"current_cache_url"`
	EntityCounts    []EntityCount `json:"entity_counts"`
}

type EntityCount struct {
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type Filter struct { // Timestamp filter or embeds, or fingerprint
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type ResultCache struct { // These are caches shown in the index endpoint of a particular entity.
	ResponseUrl string    `json:"response_url"`
	StartsFrom  Timestamp `json:"starts_from"`
	EndsAt      Timestamp `json:"ends_at"`
}

type Answer struct { // Bodies of API Endpoint responses from remote. This will be filled and unused field will be omitted.
	Boards      []Board      `json:"boards,omitempty"`
	Threads     []Thread     `json:"threads,omitempty"`
	Posts       []Post       `json:"posts,omitempty"`
	Votes       []Vote       `json:"votes,omitempty"`
	Keys        []Key        `json:"keys,omitempty"`
	Truststates []Truststate `json:"truststates,omitempty"`
	Addresses   []Address    `json:"addresses,omitempty"`

	BoardIndexes      []BoardIndex      `json:"boards_index,omitempty"`
	ThreadIndexes     []ThreadIndex     `json:"threads_index,omitempty"`
	PostIndexes       []PostIndex       `json:"posts_index,omitempty"`
	VoteIndexes       []VoteIndex       `json:"votes_index,omitempty"`
	KeyIndexes        []KeyIndex        `json:"keys_index,omitempty"`
	TruststateIndexes []TruststateIndex `json:"truststates_index,omitempty"`
	AddressIndexes    []AddressIndex    `json:"addresses_index,omitempty"`

	BoardManifests      []PageManifest `json:"boards_manifest,omitempty"`
	ThreadManifests     []PageManifest `json:"threads_manifest,omitempty"`
	PostManifests       []PageManifest `json:"posts_manifest,omitempty"`
	VoteManifests       []PageManifest `json:"votes_manifest,omitempty"`
	KeyManifests        []PageManifest `json:"keys_manifest,omitempty"`
	TruststateManifests []PageManifest `json:"truststates_manifest,omitempty"`
	AddressManifests    []PageManifest `json:"addresses_manifest,omitempty"`
}

// Manifest type
type PageManifest struct {
	Page     uint64               `json:"page_number"`
	Entities []PageManifestEntity `json:"entities"`
}

type PageManifestEntity struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	LastUpdate  Timestamp   `json:"last_update"`
}

// ApiResponse is the blueprint of all requests and responses. This is the 'external' communication structure backend uses to talk to other backends. Ideally, this should have been called ApiPayload, since api requests are also somewhat confusingly of the type ApiResponse
type ApiResponse struct {
	NodeId        Fingerprint   `json:"-"` // Generated and used at the ApiResponse signature verification, from the NodePublicKey. It doesn't transmit in or out, only generated on the fly. This blocks both inbound and outbound.
	NodePublicKey string        `json:"node_public_key,omitempty"`
	Signature     Signature     `json:"page_signature,omitempty"`
	ProofOfWork   ProofOfWork   `json:"proof_of_work"`
	Nonce         Nonce         `json:"nonce,omitempty"`
	EntityVersion int           `json:"entity_version,omitempty"`
	Address       Address       `json:"address,omitempty"`
	Entity        string        `json:"entity,omitempty"`
	Endpoint      string        `json:"endpoint,omitempty"`
	Filters       []Filter      `json:"filters,omitempty"`
	Timestamp     Timestamp     `json:"timestamp,omitempty"`
	StartsFrom    Timestamp     `json:"starts_from,omitempty"`
	EndsAt        Timestamp     `json:"ends_at,omitempty"`
	Pagination    Pagination    `json:"pagination,omitempty"`
	Caching       Caching       `json:"caching,omitempty"`
	Results       []ResultCache `json:"results,omitempty"`  // Pages
	ResponseBody  Answer        `json:"response,omitempty"` // Entities, Full size or Index versions.
}

// GetProvables gets all provables in an ApiResponse.
func (r *ApiResponse) GetProvables() *[]Provable {
	var p []Provable

	for key := range r.ResponseBody.Boards {
		p = append(p, Provable(&r.ResponseBody.Boards[key]))
	}
	for key := range r.ResponseBody.Threads {
		p = append(p, Provable(&r.ResponseBody.Threads[key]))
	}
	for key := range r.ResponseBody.Posts {
		p = append(p, Provable(&r.ResponseBody.Posts[key]))
	}
	for key := range r.ResponseBody.Votes {
		p = append(p, Provable(&r.ResponseBody.Votes[key]))
	}
	for key := range r.ResponseBody.Keys {
		p = append(p, Provable(&r.ResponseBody.Keys[key]))
	}
	for key := range r.ResponseBody.Truststates {
		p = append(p, Provable(&r.ResponseBody.Truststates[key]))
	}
	return &p
}

// Dump dumps the apiresponse in JSON format to a predetermined location on disk for inspection.
func (r *ApiResponse) Dump() error {
	fileContents, err := json.Marshal(r)
	if err != nil {
		return errors.New(fmt.Sprint(
			"This ApiResponse failed to convert to JSON. Error: %#v, ApiResponse: %#v", err, r))
	}
	path := globals.BackendConfig.GetCachesDirectory() + "/dumps/"
	toolbox.CreatePath(path)
	filename := fmt.Sprint("ApiRespDump-", time.Now().Unix(), ".json")
	// fmt.Printf("%s%s%s\n", path, "/dumps/", filename)
	err2 := ioutil.WriteFile(fmt.Sprint(path, filename), fileContents, 0755)
	if err2 != nil {
		logging.LogCrash(err2)
	}
	return nil
}

// Verify verifies all items and flags them appropriately in a response.
func (r *ApiResponse) Verify() []error {
	var errs []error

	// First of all, run boundary check on the apiresponse. This does NOT check for the .Address field, neither does it check the entities contained within. We'll do those afterwards.
	boundsOk, err := r.CheckBounds()
	if len(r.ResponseBody.PostIndexes) > 0 && boundsOk {
		// fmt.Println("this api response has post indexes in it and it verified.")
	} else if len(r.ResponseBody.PostIndexes) > 0 {
		fmt.Println("this api response has post indexes in it and it did not verify.")
	}
	if err != nil {
		return []error{fmt.Errorf("This ApiResponse failed the boundary check for its general structure (not for its contents -- it didn't come to that.) Error: %#v, ApiResponse: %#v", err, r)}
	}
	if !boundsOk {
		// logging.LogCrash("yo")
		return []error{fmt.Errorf("This ApiResponse failed the boundary check for its general structure (not for its contents -- it didn't come to that.) ApiResponse: %#v", r)}
	}
	remoteAddrOk, err := r.Address.CheckBounds()
	if err != nil {
		return []error{err}
	}
	if !remoteAddrOk {
		return []error{fmt.Errorf("This ApiResponse's remote Address failed the boundary check. ApiResponse.Address: %#v", r.Address)}
	}
	// This is all the verification we need for addresses - just a bounds check. It does not go into the more involved Verify() flow.
	for key := range r.ResponseBody.Addresses { // this is a concrete type..
		err := Verify(&r.ResponseBody.Addresses[key])
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	provables := r.GetProvables()
	for _, e := range *provables { // provable is an interface, so pointer..
		err := Verify(e)
		if err != nil && !strings.Contains(err.Error(), "This entity is in a badlist") {
			/*
				We do not count badlist errors as verification errors for the purposes of cutting the connection. The other types of errors will still count for malformed objects threshold, though.
			*/
			errs = append(errs, err)
			continue
		}
	}
	for _, err := range errs {
		logging.Log(1, err)
	}
	return errs
}

func (r *ApiResponse) ToJSON() ([]byte, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return result, errors.New(fmt.Sprint(
			"This ApiResponse failed to convert to JSON. Error: %#v, ApiResponse: %#v", err, r))
	}
	return result, nil
}

func (r *ApiResponse) Prefill() {
	subprotsAsShims := globals.BackendConfig.GetServingSubprotocols()
	var subprotsSupported []Subprotocol

	for _, val := range subprotsAsShims {
		subprotsSupported = append(subprotsSupported, Subprotocol(val))
	}
	r.NodePublicKey = globals.BackendConfig.GetMarshaledBackendPublicKey()
	addr := Address{}
	addr.LocationType = globals.BackendConfig.GetExternalIpType()
	addr.Type = globals.BackendConfig.GetNodeType()
	if !globals.BackendConfig.GetRenderNonconnectible() {
		addr.Port = uint16(globals.BackendConfig.GetExternalPort())
	}
	addr.Protocol.VersionMajor = globals.BackendConfig.GetProtocolVersionMajor()
	addr.Protocol.VersionMinor = globals.BackendConfig.GetProtocolVersionMinor()
	addr.Protocol.Subprotocols = subprotsSupported
	addr.Client.VersionMajor = globals.BackendConfig.GetClientVersionMajor()
	addr.Client.VersionMinor = globals.BackendConfig.GetClientVersionMinor()
	addr.Client.VersionPatch = globals.BackendConfig.GetClientVersionPatch()
	addr.Client.ClientName = globals.BackendConfig.GetClientName()
	addr.EntityVersion = globals.BackendTransientConfig.EntityVersions.Address
	r.EntityVersion = globals.BackendTransientConfig.EntityVersions.ApiResponse
	r.Address = addr
	r.CreateNonce()
	r.Timestamp = Timestamp(time.Now().Unix())
}

// // Interfaces

type Fingerprintable interface {
	GetFingerprint() Fingerprint // Field accessor
	CreateFingerprint() error
	VerifyFingerprint() bool
}

type PoWAble interface {
	GetProofOfWork() ProofOfWork // Field accessor
	CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error
	VerifyPoW(pubKey string) (bool, error)
}

type Signable interface {
	GetSignature() Signature   // Field accessor
	GetOwnerPublicKey() string // Field accessor
	CreateSignature(keyPair *ed25519.PrivateKey) error
	VerifySignature(pubKey string) (bool, error)
}

type BoundsCheckable interface {
	CheckBounds() (bool, error)
}

type Verifiable interface {
	Fingerprintable
	PoWAble
	Signable
	BoundsCheckable
	SetVerified(bool)
	GetVerified() bool
	VerifyEntitlements() bool
	NotInBadlist() bool
}

type Encryptable interface {
	GetEncrContent() string
}

type Shardable interface {
	GetRealmId() Fingerprint
}

type Provable interface {
	Verifiable
	Shardable
	Encryptable
	GetOwner() Fingerprint
	GetLastUpdate() Timestamp
	GetCreation() Timestamp
	GetEntityType() string
	GetLastModified() Timestamp
}

type Updateable interface {
	GetUpdateProofOfWork() ProofOfWork // Field accessor
	GetUpdateSignature() Signature     // Field accessor
	CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error
	CreateUpdateSignature(keyPair *ed25519.PrivateKey) error
}

type Versionable interface {
	GetVersion() int
}

// Accessor methods. These methods allow access to fields from the interfaces. The reason why we need these is that interfaces cannot take struct fields, so I have to create these accessor methods to let them be accessible over interfaces.

// Version accessors

func (entity *Board) GetVersion() int       { return entity.EntityVersion }
func (entity *Thread) GetVersion() int      { return entity.EntityVersion }
func (entity *Post) GetVersion() int        { return entity.EntityVersion }
func (entity *Vote) GetVersion() int        { return entity.EntityVersion }
func (entity *Key) GetVersion() int         { return entity.EntityVersion }
func (entity *Truststate) GetVersion() int  { return entity.EntityVersion }
func (entity *Address) GetVersion() int     { return entity.EntityVersion }
func (entity *ApiResponse) GetVersion() int { return entity.EntityVersion }

// Fingerprint accessors

func (entity *Board) GetFingerprint() Fingerprint      { return entity.Fingerprint }
func (entity *Thread) GetFingerprint() Fingerprint     { return entity.Fingerprint }
func (entity *Post) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Vote) GetFingerprint() Fingerprint       { return entity.Fingerprint }
func (entity *Key) GetFingerprint() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetFingerprint() Fingerprint { return entity.Fingerprint }

// LastUpdate accessors

func (entity *Board) GetLastUpdate() Timestamp      { return entity.LastUpdate }
func (entity *Thread) GetLastUpdate() Timestamp     { return entity.LastUpdate }
func (entity *Post) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *Vote) GetLastUpdate() Timestamp       { return entity.LastUpdate }
func (entity *Key) GetLastUpdate() Timestamp        { return entity.LastUpdate }
func (entity *Truststate) GetLastUpdate() Timestamp { return entity.LastUpdate }

// Creation accessors

func (entity *Board) GetCreation() Timestamp      { return entity.Creation }
func (entity *Thread) GetCreation() Timestamp     { return entity.Creation }
func (entity *Post) GetCreation() Timestamp       { return entity.Creation }
func (entity *Vote) GetCreation() Timestamp       { return entity.Creation }
func (entity *Key) GetCreation() Timestamp        { return entity.Creation }
func (entity *Truststate) GetCreation() Timestamp { return entity.Creation }

// EntityType accessors

func (entity *Board) GetEntityType() string      { return "board" }
func (entity *Thread) GetEntityType() string     { return "thread" }
func (entity *Post) GetEntityType() string       { return "post" }
func (entity *Vote) GetEntityType() string       { return "vote" }
func (entity *Key) GetEntityType() string        { return "key" }
func (entity *Truststate) GetEntityType() string { return "truststate" }
func (entity *Address) GetEntityType() string    { return "address" }

// LastModified accessors (LM: the larger of creation / lastupdate)

// (get last modified)
func glm(entity Provable) Timestamp {
	if entity.GetCreation() > entity.GetLastUpdate() {
		return entity.GetCreation()
	} else {
		return entity.GetLastUpdate()
	}
}

func (entity *Board) GetLastModified() Timestamp      { return glm(entity) }
func (entity *Thread) GetLastModified() Timestamp     { return glm(entity) }
func (entity *Post) GetLastModified() Timestamp       { return glm(entity) }
func (entity *Vote) GetLastModified() Timestamp       { return glm(entity) }
func (entity *Key) GetLastModified() Timestamp        { return glm(entity) }
func (entity *Truststate) GetLastModified() Timestamp { return glm(entity) }

// Signature accessors

func (entity *Board) GetSignature() Signature      { return entity.Signature }
func (entity *Thread) GetSignature() Signature     { return entity.Signature }
func (entity *Post) GetSignature() Signature       { return entity.Signature }
func (entity *Vote) GetSignature() Signature       { return entity.Signature }
func (entity *Key) GetSignature() Signature        { return entity.Signature }
func (entity *Truststate) GetSignature() Signature { return entity.Signature }

// OwnerPublicKey accessors

func (entity *Board) GetOwnerPublicKey() string  { return entity.OwnerPublicKey }
func (entity *Thread) GetOwnerPublicKey() string { return entity.OwnerPublicKey }
func (entity *Post) GetOwnerPublicKey() string   { return entity.OwnerPublicKey }
func (entity *Vote) GetOwnerPublicKey() string   { return entity.OwnerPublicKey }

// Heads up, this is slightly different in Key below.
func (entity *Key) GetOwnerPublicKey() string        { return entity.Key }
func (entity *Truststate) GetOwnerPublicKey() string { return entity.OwnerPublicKey }

// Verifiable accessors / setters
func (entity *Board) GetVerified() bool      { return entity.Verified }
func (entity *Thread) GetVerified() bool     { return entity.Verified }
func (entity *Post) GetVerified() bool       { return entity.Verified }
func (entity *Vote) GetVerified() bool       { return entity.Verified }
func (entity *Key) GetVerified() bool        { return entity.Verified }
func (entity *Truststate) GetVerified() bool { return entity.Verified }
func (entity *Address) GetVerified() bool    { return entity.Verified }

func (entity *Board) SetVerified(v bool)      { entity.Verified = v }
func (entity *Thread) SetVerified(v bool)     { entity.Verified = v }
func (entity *Post) SetVerified(v bool)       { entity.Verified = v }
func (entity *Vote) SetVerified(v bool)       { entity.Verified = v }
func (entity *Key) SetVerified(v bool)        { entity.Verified = v }
func (entity *Truststate) SetVerified(v bool) { entity.Verified = v }
func (entity *Address) SetVerified(v bool)    { entity.Verified = v }

// UpdateSignature accessors

func (entity *Board) GetUpdateSignature() Signature      { return entity.UpdateSignature }
func (entity *Thread) GetUpdateSignature() Signature     { return entity.UpdateSignature }
func (entity *Post) GetUpdateSignature() Signature       { return entity.UpdateSignature }
func (entity *Vote) GetUpdateSignature() Signature       { return entity.UpdateSignature }
func (entity *Key) GetUpdateSignature() Signature        { return entity.UpdateSignature }
func (entity *Truststate) GetUpdateSignature() Signature { return entity.UpdateSignature }

// ProofOfWork accessors

func (entity *Board) GetProofOfWork() ProofOfWork      { return entity.ProofOfWork }
func (entity *Thread) GetProofOfWork() ProofOfWork     { return entity.ProofOfWork }
func (entity *Post) GetProofOfWork() ProofOfWork       { return entity.ProofOfWork }
func (entity *Vote) GetProofOfWork() ProofOfWork       { return entity.ProofOfWork }
func (entity *Key) GetProofOfWork() ProofOfWork        { return entity.ProofOfWork }
func (entity *Truststate) GetProofOfWork() ProofOfWork { return entity.ProofOfWork }

// UpdateProofOfWork accessors

func (entity *Board) GetUpdateProofOfWork() ProofOfWork      { return entity.UpdateProofOfWork }
func (entity *Thread) GetUpdateProofOfWork() ProofOfWork     { return entity.UpdateProofOfWork }
func (entity *Post) GetUpdateProofOfWork() ProofOfWork       { return entity.UpdateProofOfWork }
func (entity *Vote) GetUpdateProofOfWork() ProofOfWork       { return entity.UpdateProofOfWork }
func (entity *Key) GetUpdateProofOfWork() ProofOfWork        { return entity.UpdateProofOfWork }
func (entity *Truststate) GetUpdateProofOfWork() ProofOfWork { return entity.UpdateProofOfWork }

// Signature accessors

func (entity *Board) GetOwner() Fingerprint  { return entity.Owner }
func (entity *Thread) GetOwner() Fingerprint { return entity.Owner }
func (entity *Post) GetOwner() Fingerprint   { return entity.Owner }
func (entity *Vote) GetOwner() Fingerprint   { return entity.Owner }

// (For below, owner of the entity is itself.)
func (entity *Key) GetOwner() Fingerprint        { return entity.Fingerprint }
func (entity *Truststate) GetOwner() Fingerprint { return entity.Owner }

// RealmId accessors

func (entity *Board) GetRealmId() Fingerprint      { return entity.RealmId }
func (entity *Thread) GetRealmId() Fingerprint     { return entity.RealmId }
func (entity *Post) GetRealmId() Fingerprint       { return entity.RealmId }
func (entity *Vote) GetRealmId() Fingerprint       { return entity.RealmId }
func (entity *Key) GetRealmId() Fingerprint        { return entity.RealmId }
func (entity *Truststate) GetRealmId() Fingerprint { return entity.RealmId }

// EncrContent accessors

func (entity *Board) GetEncrContent() string      { return entity.EncrContent }
func (entity *Thread) GetEncrContent() string     { return entity.EncrContent }
func (entity *Post) GetEncrContent() string       { return entity.EncrContent }
func (entity *Vote) GetEncrContent() string       { return entity.EncrContent }
func (entity *Key) GetEncrContent() string        { return entity.EncrContent }
func (entity *Truststate) GetEncrContent() string { return entity.EncrContent }

// Response styles.

// Response is the interface junction that batch processing functions take and emit. This is the 'internal' communication structure within the backend. It is the big carrier type for the end result of a pull from a remote.
type Response struct {
	Boards      []Board
	Threads     []Thread
	Posts       []Post
	Votes       []Vote
	Keys        []Key
	Addresses   []Address
	Truststates []Truststate

	BoardIndexes      []BoardIndex
	ThreadIndexes     []ThreadIndex
	PostIndexes       []PostIndex
	VoteIndexes       []VoteIndex
	KeyIndexes        []KeyIndex
	AddressIndexes    []AddressIndex
	TruststateIndexes []TruststateIndex

	BoardManifests      []PageManifest
	ThreadManifests     []PageManifest
	PostManifests       []PageManifest
	VoteManifests       []PageManifest
	KeyManifests        []PageManifest
	TruststateManifests []PageManifest
	AddressManifests    []PageManifest

	CacheLinks                []ResultCache
	MostRecentSourceTimestamp Timestamp
}

func (r *Response) Empty() bool {
	return len(r.Boards) == 0 &&
		len(r.Threads) == 0 &&
		len(r.Posts) == 0 &&
		len(r.Votes) == 0 &&
		len(r.Keys) == 0 &&
		len(r.Truststates) == 0 &&
		len(r.Addresses) == 0 &&

		len(r.BoardIndexes) == 0 &&
		len(r.ThreadIndexes) == 0 &&
		len(r.PostIndexes) == 0 &&
		len(r.VoteIndexes) == 0 &&
		len(r.KeyIndexes) == 0 &&
		len(r.TruststateIndexes) == 0 &&
		len(r.AddressIndexes) == 0 &&

		len(r.BoardManifests) == 0 &&
		len(r.ThreadManifests) == 0 &&
		len(r.PostManifests) == 0 &&
		len(r.VoteManifests) == 0 &&
		len(r.KeyManifests) == 0 &&
		len(r.TruststateManifests) == 0 &&
		len(r.AddressManifests) == 0 &&

		len(r.CacheLinks) == 0
}

func (r *Response) Insert(r2 *Response) {
	r.Boards = append(r.Boards, r2.Boards...)
	r.Threads = append(r.Threads, r2.Threads...)
	r.Posts = append(r.Posts, r2.Posts...)
	r.Votes = append(r.Votes, r2.Votes...)
	r.Keys = append(r.Keys, r2.Keys...)
	r.Truststates = append(r.Truststates, r2.Truststates...)
	r.Addresses = append(r.Addresses, r2.Addresses...)

	r.BoardIndexes = append(r.BoardIndexes, r2.BoardIndexes...)
	r.ThreadIndexes = append(r.ThreadIndexes, r2.ThreadIndexes...)
	r.PostIndexes = append(r.PostIndexes, r2.PostIndexes...)
	r.VoteIndexes = append(r.VoteIndexes, r2.VoteIndexes...)
	r.KeyIndexes = append(r.KeyIndexes, r2.KeyIndexes...)
	r.TruststateIndexes = append(r.TruststateIndexes, r2.TruststateIndexes...)
	r.AddressIndexes = append(r.AddressIndexes, r2.AddressIndexes...)

	r.BoardManifests = append(r.BoardManifests, r2.BoardManifests...)
	r.ThreadManifests = append(r.ThreadManifests, r2.ThreadManifests...)
	r.PostManifests = append(r.PostManifests, r2.PostManifests...)
	r.VoteManifests = append(r.VoteManifests, r2.VoteManifests...)
	r.KeyManifests = append(r.KeyManifests, r2.KeyManifests...)
	r.TruststateManifests = append(r.TruststateManifests, r2.TruststateManifests...)
	r.AddressManifests = append(r.AddressManifests, r2.AddressManifests...)

	r.CacheLinks = append(r.CacheLinks, r2.CacheLinks...)

	if r.MostRecentSourceTimestamp < r2.MostRecentSourceTimestamp {
		r.MostRecentSourceTimestamp = r2.MostRecentSourceTimestamp
	} else {
		r.MostRecentSourceTimestamp = r.MostRecentSourceTimestamp
	}
}

func (r *Response) IndexOf(item Provable) int {
	switch entity := item.(type) {
	case *Board:
		for key := range r.Boards {
			if r.Boards[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *Thread:
		for key := range r.Threads {
			if r.Threads[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *Post:
		for key := range r.Posts {
			if r.Posts[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *Vote:
		for key := range r.Votes {
			if r.Votes[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *Key:
		for key := range r.Keys {
			if r.Keys[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *Truststate:
		for key := range r.Truststates {
			if r.Truststates[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	}
	return -1
}

func (r *Response) RemoveByIndex(i int, entityType string) {
	if i == -1 {
		return
	}
	switch entityType {
	case "board":
		if len(r.Boards) > i { // i:5, len(boards): 5, fails. because boards[5] is out of bounds.
			r.Boards = append(r.Boards[0:i], r.Boards[i+1:len(r.Boards)]...)
		}
	case "thread":
		if len(r.Threads) > i {
			r.Threads = append(r.Threads[0:i], r.Threads[i+1:len(r.Threads)]...)
		}
	case "post":
		if len(r.Posts) > i {
			r.Posts = append(r.Posts[0:i], r.Posts[i+1:len(r.Posts)]...)
		}
	case "vote":
		if len(r.Votes) > i {
			r.Votes = append(r.Votes[0:i], r.Votes[i+1:len(r.Votes)]...)
		}
	case "key":
		if len(r.Keys) > i {
			r.Keys = append(r.Keys[0:i], r.Keys[i+1:len(r.Keys)]...)
		}
	case "truststate":
		if len(r.Truststates) > i {
			r.Truststates = append(r.Truststates[0:i], r.Truststates[i+1:len(r.Truststates)]...)
		}
	default:
		logging.LogCrash(fmt.Sprintf("You gave Response.RemoveByIndex an unknown entity type. You gave: %s", entityType))
	}
}

func isInIndexSlice(idx int, idxs []int) bool {
	for key := range idxs {
		if idx == idxs[key] {
			idxs = append(idxs[0:key], idxs[key+1:len(idxs)]...)
			return true
		}
	}
	return false
}

func (r *Response) MassRemoveByIndex(idxs []int, entityType string) {
	if len(idxs) == 0 {
		return
	}
	switch entityType {
	case "board":
		if len(r.Boards) == len(idxs) {
			// There is no way any entities will remain unless indexes are nonexistent.
			r.Boards = []Board{}
			return
		}
		var retained []Board

		for key := range r.Boards {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Boards[key])
			}
		}
		r.Boards = retained
	case "thread":
		if len(r.Threads) == len(idxs) {
			r.Threads = []Thread{}
			return
		}
		var retained []Thread

		for key := range r.Threads {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Threads[key])
			}
		}
		r.Threads = retained
	case "post":
		if len(r.Posts) == len(idxs) {
			r.Posts = []Post{}
			return
		}
		var retained []Post

		for key := range r.Posts {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Posts[key])
			}
		}
		r.Posts = retained
	case "vote":
		if len(r.Votes) == len(idxs) {
			r.Votes = []Vote{}
			return
		}
		var retained []Vote

		for key := range r.Votes {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Votes[key])
			}
		}
		r.Votes = retained
	case "key":
		if len(r.Keys) == len(idxs) {
			r.Keys = []Key{}
			return
		}
		var retained []Key

		for key := range r.Keys {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Keys[key])
			}
		}
		r.Keys = retained
	case "truststate":
		if len(r.Truststates) == len(idxs) {
			r.Truststates = []Truststate{}
			return
		}
		var retained []Truststate

		for key := range r.Truststates {
			if !isInIndexSlice(key, idxs) {
				retained = append(retained, r.Truststates[key])
			}
		}
		r.Truststates = retained
	default:
		logging.LogCrash(fmt.Sprintf("You gave Response.RemoveByIndex an unknown entity type. You gave: %s", entityType))
	}
}
