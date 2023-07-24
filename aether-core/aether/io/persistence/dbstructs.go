// Persistence > Structs
// This file contains the struct definitions of the database objects.

package persistence

import (
	"aether-core/aether/io/api"
	// "aether-core/aether/services/compress"
	"aether-core/aether/services/fingerprinting"
	"aether-core/aether/services/logging"
	"encoding/csv"
	"encoding/json"
	"fmt"

	// "github.com/jmoiron/sqlx/types"
	"strings"
	"time"
)

// Basic properties

type DbUpdateable struct {
	LastUpdate        api.Timestamp   `db:"LastUpdate"`
	UpdateProofOfWork api.ProofOfWork `db:"UpdateProofOfWork"`
	UpdateSignature   api.Signature   `db:"UpdateSignature"`
}

type DbProvable struct {
	Creation    api.Timestamp   `db:"Creation"`
	ProofOfWork api.ProofOfWork `db:"ProofOfWork"`
	Signature   api.Signature   `db:"Signature"`
}

// Subentities

type DbBoardOwner struct {
	BoardFingerprint      api.Fingerprint `db:"BoardFingerprint"`
	KeyFingerprint        api.Fingerprint `db:"KeyFingerprint"`
	Expiry                api.Timestamp   `db:"Expiry"`
	Level                 uint8           `db:"Level"`
	ParentBoardCreation   api.Timestamp   `db:"ParentBoardCreation"`
	ParentBoardLastUpdate api.Timestamp   `db:"ParentBoardLastUpdate"`
}

type DbSubprotocol struct {
	Fingerprint       api.Fingerprint `db:"Fingerprint"`
	Name              string          `db:"Name"`
	VersionMajor      uint8           `db:"VersionMajor"`
	VersionMinor      uint16          `db:"VersionMinor"`
	SupportedEntities string          `db:"SupportedEntities"`
}

// Junction Table Entities

type DbAddressSubprotocol struct {
	AddressLocation        api.Location    `db:"AddressLocation"`
	AddressSublocation     api.Location    `db:"AddressSublocation"`
	AddressPort            uint16          `db:"AddressPort"`
	SubprotocolFingerprint api.Fingerprint `db:"SubprotocolFingerprint"`
}

// Entities

type DbBoard struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Name           string          `db:"Name"`
	Owner          api.Fingerprint `db:"Owner"`
	OwnerPublicKey string          `db:"OwnerPublicKey"`
	Description    string          `db:"Description"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Language       string          `db:"Language"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbThread struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Board          api.Fingerprint `db:"Board"`
	Name           string          `db:"Name"`
	Body           string          `db:"Body"`
	Link           string          `db:"Link"`
	Owner          api.Fingerprint `db:"Owner"`
	OwnerPublicKey string          `db:"OwnerPublicKey"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbPost struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Board          api.Fingerprint `db:"Board"`
	Thread         api.Fingerprint `db:"Thread"`
	Parent         api.Fingerprint `db:"Parent"`
	Body           string          `db:"Body"`
	Owner          api.Fingerprint `db:"Owner"`
	OwnerPublicKey string          `db:"OwnerPublicKey"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbVote struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Board          api.Fingerprint `db:"Board"`
	Thread         api.Fingerprint `db:"Thread"`
	Target         api.Fingerprint `db:"Target"`
	Owner          api.Fingerprint `db:"Owner"`
	OwnerPublicKey string          `db:"OwnerPublicKey"`
	TypeClass      int             `db:"TypeClass"`
	Type           int             `db:"Type"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbKey struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Type           string          `db:"Type"`
	PublicKey      string          `db:"PublicKey"`
	Expiry         api.Timestamp   `db:"Expiry"`
	Name           string          `db:"Name"`
	Info           string          `db:"Info"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbTruststate struct {
	Fingerprint    api.Fingerprint `db:"Fingerprint"`
	Target         api.Fingerprint `db:"Target"`
	Owner          api.Fingerprint `db:"Owner"`
	OwnerPublicKey string          `db:"OwnerPublicKey"`
	TypeClass      int             `db:"TypeClass"`
	Type           int             `db:"Type"`
	Domain         api.Fingerprint `db:"Domain"`
	Expiry         api.Timestamp   `db:"Expiry"`
	LocalArrival   api.Timestamp   `db:"LocalArrival"`
	LastReferenced api.Timestamp   `db:"LastReferenced"`
	EntityVersion  int             `db:"EntityVersion"`
	Meta           string          `db:"Meta"`
	RealmId        api.Fingerprint `db:"RealmId"`
	EncrContent    string          `db:"EncrContent"`
	DbProvable
	DbUpdateable
}

type DbAddress struct {
	Location             api.Location  `db:"Location"`
	Sublocation          api.Location  `db:"Sublocation"`
	Port                 uint16        `db:"Port"`
	LocationType         uint8         `db:"IPType"`
	Type                 uint8         `db:"AddressType"`
	LastSuccessfulPing   api.Timestamp `db:"LastSuccessfulPing"`
	LastSuccessfulSync   api.Timestamp `db:"LastSuccessfulSync"`
	ProtocolVersionMajor uint8         `db:"ProtocolVersionMajor"`
	ProtocolVersionMinor uint16        `db:"ProtocolVersionMinor"`
	ClientVersionMajor   uint8         `db:"ClientVersionMajor"`
	ClientVersionMinor   uint16        `db:"ClientVersionMinor"`
	ClientVersionPatch   uint16        `db:"ClientVersionPatch"`
	ClientName           string        `db:"ClientName"`
	LocalArrival         api.Timestamp `db:"LocalArrival"`
	// LastReferenced       api.Timestamp   `db:"LastReferenced"`
	EntityVersion int             `db:"EntityVersion"`
	RealmId       api.Fingerprint `db:"RealmId"`
}

// Non-communicating entities
type DbNode struct {
	Fingerprint            api.Fingerprint `db:"Fingerprint"`
	BoardsLastCheckin      api.Timestamp   `db:"BoardsLastCheckin"`
	ThreadsLastCheckin     api.Timestamp   `db:"ThreadsLastCheckin"`
	PostsLastCheckin       api.Timestamp   `db:"PostsLastCheckin"`
	VotesLastCheckin       api.Timestamp   `db:"VotesLastCheckin"`
	KeysLastCheckin        api.Timestamp   `db:"KeysLastCheckin"`
	TruststatesLastCheckin api.Timestamp   `db:"TruststatesLastCheckin"`
	AddressesLastCheckin   api.Timestamp   `db:"AddressesLastCheckin"`
}

// Return types of APIToDB. This is necessary because some API objects, when converted to their DB form, return more than one DB object.

type BoardPack struct {
	Board       DbBoard
	BoardOwners []DbBoardOwner
}

type AddressPack struct {
	Address      DbAddress
	Subprotocols []DbSubprotocol
	Junctions    []DbAddressSubprotocol
}

// APItoDB translates structs of API objects into structs of DB objects. It also checks whether the objects are verified or not, and if not, it prevents them entry into the database layer.
func APItoDB(object interface{}, ts time.Time) (interface{}, error) {
	now := ts.Unix()
	switch obj := object.(type) {
	// obj: typed API object.
	case api.Board:
		// Corner case: board owners
		if !obj.GetVerified() {
			return BoardPack{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbBoard
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Name = obj.Name
		dbObj.Owner = obj.Owner
		dbObj.OwnerPublicKey = obj.OwnerPublicKey
		dbObj.Description = obj.Description
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Language = obj.Language
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		var DbBoardOwners []DbBoardOwner
		for _, val := range obj.BoardOwners {
			var bo DbBoardOwner
			bo.BoardFingerprint = obj.Fingerprint
			bo.KeyFingerprint = val.KeyFingerprint
			bo.Expiry = val.Expiry
			bo.Level = val.Level
			bo.ParentBoardCreation = obj.Creation
			bo.ParentBoardLastUpdate = obj.LastUpdate
			DbBoardOwners = append(DbBoardOwners, bo)
		}
		var result BoardPack
		result.Board = dbObj
		result.BoardOwners = DbBoardOwners
		return result, nil

	case api.Thread:
		if !obj.GetVerified() {
			return DbThread{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbThread
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Name = obj.Name
		dbObj.Body = obj.Body
		dbObj.Link = obj.Link
		dbObj.Owner = obj.Owner
		dbObj.OwnerPublicKey = obj.OwnerPublicKey
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		return dbObj, nil

	case api.Post:
		if !obj.GetVerified() {
			return DbPost{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		// If a post is its own parent, this is a maliciously crafted post and we do not accept this in. Verification system also checks for this.
		if obj.Parent == obj.Fingerprint {
			logging.Log(1, fmt.Sprintf("We've received a post whose parent is itself, attempting to cross into the database. This might happen normally if you have disabled the verification, but if the verification is enabled in this node, investigate. Post: %#v", obj))
			return DbPost{}, fmt.Errorf("This post's parent is its own fingerprint, therefore attempting to insert this to the database would create infinite recursion up to recursion limit. This does not break the app, but it freezes it for a couple minutes until the recursion limit is reached. This post is thus denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbPost
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Thread = obj.Thread
		dbObj.Parent = obj.Parent
		dbObj.Body = obj.Body
		dbObj.Owner = obj.Owner
		dbObj.OwnerPublicKey = obj.OwnerPublicKey
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		return dbObj, nil

	case api.Vote:
		if !obj.GetVerified() {
			return DbVote{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbVote
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Board = obj.Board
		dbObj.Thread = obj.Thread
		dbObj.Target = obj.Target
		dbObj.Owner = obj.Owner
		dbObj.OwnerPublicKey = obj.OwnerPublicKey
		dbObj.TypeClass = obj.TypeClass
		dbObj.Type = obj.Type
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		return dbObj, nil

	case api.Key:
		if !obj.GetVerified() {
			return DbKey{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbKey
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Type = obj.Type
		dbObj.PublicKey = obj.Key
		dbObj.Expiry = obj.Expiry
		dbObj.Name = obj.Name
		dbObj.Info = obj.Info
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		return dbObj, nil

	case api.Truststate:
		if !obj.GetVerified() {
			return DbTruststate{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbTruststate
		dbObj.Fingerprint = obj.Fingerprint
		dbObj.Target = obj.Target
		dbObj.Owner = obj.Owner
		dbObj.OwnerPublicKey = obj.OwnerPublicKey
		dbObj.TypeClass = obj.TypeClass
		dbObj.Type = obj.Type
		dbObj.Expiry = obj.Expiry
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.Meta = obj.Meta
		dbObj.RealmId = obj.RealmId
		dbObj.EncrContent = obj.EncrContent
		// Provable set
		dbObj.Creation = obj.Creation
		dbObj.ProofOfWork = obj.ProofOfWork
		dbObj.Signature = obj.Signature
		// Updateable set
		dbObj.LastUpdate = obj.LastUpdate
		dbObj.UpdateProofOfWork = obj.UpdateProofOfWork
		dbObj.UpdateSignature = obj.UpdateSignature
		dbObj.Domain = obj.Domain
		return dbObj, nil

	case api.Address:
		if !obj.GetVerified() {
			return AddressPack{}, fmt.Errorf("This Api entity failed verification (or the verification hasn't been run on it), thus is denied conversion to the Db entity. Entity %#v", obj)
		}
		var dbObj DbAddress
		dbObj.Location = obj.Location
		dbObj.Sublocation = obj.Sublocation
		dbObj.LocationType = obj.LocationType
		dbObj.Port = obj.Port
		dbObj.Type = obj.Type
		dbObj.LastSuccessfulPing = obj.LastSuccessfulPing
		dbObj.LastSuccessfulSync = obj.LastSuccessfulSync
		dbObj.ProtocolVersionMajor = obj.Protocol.VersionMajor
		dbObj.ProtocolVersionMinor = obj.Protocol.VersionMinor
		dbObj.ClientVersionMajor = obj.Client.VersionMajor
		dbObj.ClientVersionMinor = obj.Client.VersionMinor
		dbObj.ClientVersionPatch = obj.Client.VersionPatch
		dbObj.ClientName = obj.Client.ClientName
		// now := time.Now().Unix()
		dbObj.LocalArrival = api.Timestamp(now)
		// dbObj.LastReferenced = api.Timestamp(now)
		dbObj.EntityVersion = obj.EntityVersion
		dbObj.RealmId = obj.RealmId
		var ap AddressPack
		ap.Address = dbObj
		// Loop over subprotocols and insert to pack.
		var subprotocols []DbSubprotocol
		for _, val := range obj.Protocol.Subprotocols {
			var s DbSubprotocol
			s.Name = val.Name
			s.VersionMajor = val.VersionMajor
			s.VersionMinor = val.VersionMinor
			// fmt.Printf("%#v", val.SupportedEntities)
			// Convert protocol entities in subprotocol into comma separated list
			parsedStr, err := parseStringSliceToCommaSeparatedString(val.SupportedEntities, "addresses")
			if err != nil {
				logging.Log(1, fmt.Sprintf("Subprotocol %s has an error in its supported entities list. Supported Entities List: %#v Error: %s. This field will be saved as empty in the database.", s.Name, val.SupportedEntities, err))
				return ap, err
			}
			s.SupportedEntities = parsedStr
			// Create Fingerprint for the entity. Mind that this is an internal FP useful for local access purposes, and this is not communicated externally.
			res, _ := json.Marshal(s)
			fp := fingerprinting.Create(string(res))
			s.Fingerprint = api.Fingerprint(fp)
			subprotocols = append(subprotocols, s)
			jItem := generateAdrSprotJunctionItem(dbObj, s)
			ap.Junctions = append(ap.Junctions, jItem)
		}
		ap.Subprotocols = subprotocols
		return ap, nil

	default:
		return nil, fmt.Errorf("APItoDB only takes API (not DB) objects. Your object: %#v\n", obj)
	}
}

// DBtoAPI converts a DB type to an API type. This is useful when this object needs to communicated to the outside world, such as showing it to the user, or to send it over the wire. This merges certain objects (Boards will have their BoardOwners appended) and removes some information (local arrival time).
func DBtoAPI(object interface{}) (interface{}, error) {
	switch obj := object.(type) {
	// obj: typed DB object.
	case DbBoard:
		// Corner case: has to query BoardOwners, too.
		var apiObj api.Board
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Name = obj.Name
		apiObj.Owner = obj.Owner
		apiObj.OwnerPublicKey = obj.OwnerPublicKey
		apiObj.Description = string(obj.Description)
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Language = obj.Language
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		// Pull the board owners for this board from database.
		dbBoardOwners, err := ReadDBBoardOwners(obj.Fingerprint, "")
		if err != nil {
			// This should always crash, it means the local remote lost / corrupted data as network always provides sub-entities and main entity together.
			logging.LogCrash(err)
		}
		for _, dbBoardOwner := range dbBoardOwners {
			var apiBoardOwner api.BoardOwner
			apiBoardOwner.KeyFingerprint = dbBoardOwner.KeyFingerprint
			apiBoardOwner.Expiry = dbBoardOwner.Expiry
			apiBoardOwner.Level = dbBoardOwner.Level
			apiObj.BoardOwners = append(apiObj.BoardOwners, apiBoardOwner)
		}
		return apiObj, nil

	case DbThread:
		var apiObj api.Thread
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Name = obj.Name
		apiObj.Body = string(obj.Body)
		apiObj.Link = obj.Link
		apiObj.Owner = obj.Owner
		apiObj.OwnerPublicKey = obj.OwnerPublicKey
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		return apiObj, nil

	case DbPost:
		var apiObj api.Post
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Thread = obj.Thread
		apiObj.Parent = obj.Parent
		apiObj.Body = string(obj.Body)
		apiObj.Owner = obj.Owner
		apiObj.OwnerPublicKey = obj.OwnerPublicKey
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		return apiObj, nil

	case DbVote:
		var apiObj api.Vote
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Board = obj.Board
		apiObj.Thread = obj.Thread
		apiObj.Target = obj.Target
		apiObj.Owner = obj.Owner
		apiObj.OwnerPublicKey = obj.OwnerPublicKey
		apiObj.TypeClass = obj.TypeClass
		apiObj.Type = obj.Type
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		return apiObj, nil

	case DbKey:
		var apiObj api.Key
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Type = obj.Type
		apiObj.Key = obj.PublicKey
		apiObj.Expiry = obj.Expiry
		apiObj.Name = obj.Name
		apiObj.Info = string(obj.Info)
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		return apiObj, nil

	case DbTruststate:
		// Corner case, comma separated fingerprint parse
		var apiObj api.Truststate
		apiObj.Fingerprint = obj.Fingerprint
		apiObj.Target = obj.Target
		apiObj.Owner = obj.Owner
		apiObj.OwnerPublicKey = obj.OwnerPublicKey
		apiObj.TypeClass = obj.TypeClass
		apiObj.Type = obj.Type
		apiObj.Expiry = obj.Expiry
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.Meta = obj.Meta
		apiObj.RealmId = obj.RealmId
		apiObj.EncrContent = obj.EncrContent
		// Provable set
		apiObj.Creation = obj.Creation
		apiObj.ProofOfWork = obj.ProofOfWork
		apiObj.Signature = obj.Signature
		// Updateable set
		apiObj.LastUpdate = obj.LastUpdate
		apiObj.UpdateProofOfWork = obj.UpdateProofOfWork
		apiObj.UpdateSignature = obj.UpdateSignature
		apiObj.Domain = obj.Domain
		return apiObj, nil

	case DbAddress:
		// Corner case
		var apiObj api.Address
		apiObj.Location = obj.Location
		apiObj.Sublocation = obj.Sublocation
		apiObj.LocationType = obj.LocationType
		apiObj.Port = obj.Port
		apiObj.Type = obj.Type
		apiObj.LastSuccessfulPing = obj.LastSuccessfulPing
		apiObj.LastSuccessfulSync = obj.LastSuccessfulSync
		apiObj.Protocol.VersionMajor = obj.ProtocolVersionMajor
		apiObj.Protocol.VersionMinor = obj.ProtocolVersionMinor
		apiObj.Client.VersionMajor = obj.ClientVersionMajor
		apiObj.Client.VersionMinor = obj.ClientVersionMinor
		apiObj.Client.VersionPatch = obj.ClientVersionPatch
		apiObj.Client.ClientName = obj.ClientName
		apiObj.EntityVersion = obj.EntityVersion
		apiObj.RealmId = obj.RealmId
		dbSubprotocols, err := ReadDBSubprotocols(obj.Location, obj.Sublocation, obj.Port)
		if err != nil {
			// This should always crash, it means the local remote lost / corrupted data as network always provides sub-entities and main entity together.
			logging.LogCrash(err)
		}
		// Convert dbSubprotocols to api.Subprotocols
		var apiSubprotocols []api.Subprotocol
		for _, dbSubprot := range dbSubprotocols {
			var apiSp api.Subprotocol
			apiSp.Name = dbSubprot.Name
			apiSp.VersionMajor = dbSubprot.VersionMajor
			apiSp.VersionMinor = dbSubprot.VersionMinor
			parsedStrSlice, err2 := parseCommaSeparatedStringToStringSlice(dbSubprot.SupportedEntities, 64, 100)
			if err2 != nil {
				return apiSp, err2
			}
			apiSp.SupportedEntities = parsedStrSlice
			apiSubprotocols = append(apiSubprotocols, apiSp)
		}
		apiObj.Protocol.Subprotocols = apiSubprotocols
		return apiObj, nil

	case DbBoardOwner:
		return nil, fmt.Errorf("This object cannot be queried on its own. Try querying the parent Board object. Your object: %#v\n", obj)
	default:
		return nil, fmt.Errorf("DBtoAPI only takes DB (not API) objects. Your object: %#v\n", obj)
	}
}

// parseCommaSeparatedStringToStringSlice converts a single, comma separated string into a group of strings.
func parseCommaSeparatedStringToStringSlice(str string, maxLen int, maxCount int) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(str))
	result, err := reader.ReadAll()
	if err != nil {
		logging.LogCrash(err)
	}
	// Trim and check sanity.
	var cleaned []string
	var err2 error
	if result != nil { // If this isn't present, result[0] ends up undefined.
		for i, val := range result[0] {
			if i >= maxCount {
				// Too many items for the field. It's here, because otherwise somebody could send a 100 billion field string and brick the remote. We're still printing it, but not parsing it any more.
				return nil, fmt.Errorf("The string provided has too many items. String: %#v\n", val)
			}
			if len(val) < maxLen && len(val) > 0 {
				cleaned = append(cleaned, strings.Trim(val, " "))
			} else if len(val) < maxLen { // The length of the string is zero.
				err2 = fmt.Errorf("This string is empty. String: %#v\n", val)
			} else { // The string is longer than max length allowed for field.
				err2 = fmt.Errorf("This string is too long for this field. String: %#v\n", val)
			}
		}
		err3 := checkDuplicatesInStringSlice(cleaned)
		if err3 != nil {
			return cleaned, err3
		}
	}
	return cleaned, err2
}

// parseStringSliceToCommaSeparatedString converts a slice of strings into one single comma separated string.
func parseStringSliceToCommaSeparatedString(strs []string, entityType string) (string, error) {
	maxLen := 0
	maxCount := 0
	if entityType == "addresses" {
		maxLen = api.MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_NAME_V1
		maxCount = api.MAX_ADDRESS_PROTOCOL_SUBPROTOCOL_SUPPORTEDENTITIES_V1
	}
	var err error
	var finalStr string
	if len(strs) > maxCount {
		return finalStr, fmt.Errorf("The string slice provided has too many items. String slice: %#v\n", strs)
	} else {
		// We know it has an acceptable amount of elements. Let's check for duplicates.
		err := checkDuplicatesInStringSlice(strs)
		if err != nil {
			return finalStr, err
		}
	}
	for i, str := range strs {
		// Check sanity
		if len(str) < maxLen && len(str) > 0 {
			// Convert to comma separated string.
			// If this is the first item, no comma at the beginning.
			if i == 0 {
				finalStr = str
			} else {
				finalStr = fmt.Sprint(finalStr, ",", str)
			}
		} else if len(str) < maxLen { // The length of the string is zero.
			err = fmt.Errorf("This string is empty. String: %#v\n", str)
		} else { // The string is longer than max length allowed for field.
			err = fmt.Errorf("This string is too long for this field. String: %#v\n", str)
		}
	}
	return finalStr, err
}

func checkDuplicatesInStringSlice(strs []string) error {
	mappy := make(map[string]int)
	for _, val := range strs {
		mappy[val]++
	}
	for str, occurrenceCount := range mappy {
		if occurrenceCount > 1 {
			return fmt.Errorf("This list includes items that are duplicates. Duplicate item: %#v\n", str)
		}
	}
	return nil
}

// func convertFingerprintSliceToStringSlice(fps []api.Fingerprint) []string {
// 	var strSlice []string
// 	for _, val := range fps {
// 		strSlice = append(strSlice, string(val))
// 	}
// 	return strSlice
// }

// func convertStringSliceToFingerprintSlice(strs []string) []api.Fingerprint {
// 	var fpSlice []api.Fingerprint
// 	for _, val := range strs {
// 		fpSlice = append(fpSlice, api.Fingerprint(val))
// 	}
// 	return fpSlice
// }
