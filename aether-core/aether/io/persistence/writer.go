// Persistence > Writer
// This file collects all of the functions that write to the database. UI uses this for insertions, as well as the Fetcher.

package persistence

import (
	"aether-core/aether/backend/feapiconsumer"
	"aether-core/aether/io/api"
	"fmt"

	// _ "github.com/mattn/go-sqlite3"
	// "aether-core/aether/backend/metrics"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"errors"

	"github.com/fatih/color"

	// "github.com/jmoiron/sqlx/types"
	// "github.com/davecgh/go-spew/spew"
	// "runtime"
	// "github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// Node is a non-communicating entity that holds the LastCheckin timestamps of each of the entities provided in the remote node. There is no way to send this data over to somebody, this is entirely local. There is also no batch processing because there is no situation in which you would need to insert multiple nodes at the same time (since you won't be connecting to multiple nodes simultaneously)

func InsertNode(n DbNode) error {
	err := insertNode(n)
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
			time.Sleep(10 * time.Second)
			logging.Log(1, "Retrying the previously failed InsertNode transaction.")
			err2 := insertNode(n)
			if err2 != nil {
				if strings.Contains(err2.Error(), "Database was locked") {
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // err2 != nil, but error isn't "database was locked"
					return errors.New(fmt.Sprintf("InsertNode encountered an error. Error: %s", err2))
				}
			} else { // If the reattempted transaction succeeds
				logging.Log(1,
					"The retry attempt of the failed transaction succeeded.")
			}
		} else { // err != nil, but error isn't "database was locked"
			logging.Log(1, err)
			return errors.New(fmt.Sprintf("InsertNode encountered an error. Error: %s", err))
		}
	}
	return nil
}
func insertNode(n DbNode) error {
	// fmt.Println(toolbox.DumpStack())
	// fmt.Println("Node to be inserted:")
	// spew.Dump(n)
	// fmt.Printf("%#v\n", n)
	if api.Fingerprint(globals.BackendConfig.GetNodeId()) == n.Fingerprint {
		return errors.New(fmt.Sprintf("The node ID that was attempted to be inserted is the SAME AS the local node's ID. This could be an attempted attack. Node ID of the remote: %s", n.Fingerprint))
	}
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		return err
	}
	_, err2 := tx.NamedExec(nodeInsert, n)
	if err2 != nil {
		return err2
	}
	err3 := tx.Commit()
	if err3 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("InsertNode encountered an error when trying to commit to the database. Error is: %s", err3))
		if strings.Contains(err3.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err3
	}

	nodeAsMap := make(map[string]string)
	nodeAsMap["Fingerprint"] = string(n.Fingerprint)
	nodeAsMap["BoardsLastCheckin"] = strconv.Itoa(int(n.BoardsLastCheckin))
	nodeAsMap["ThreadsLastCheckin"] = strconv.Itoa(int(n.ThreadsLastCheckin))
	nodeAsMap["PostsLastCheckin"] = strconv.Itoa(int(n.PostsLastCheckin))
	nodeAsMap["VotesLastCheckin"] = strconv.Itoa(int(n.VotesLastCheckin))
	nodeAsMap["KeysLastCheckin"] = strconv.Itoa(int(n.KeysLastCheckin))
	nodeAsMap["TruststatesLastCheckin"] = strconv.Itoa(int(n.TruststatesLastCheckin))
	nodeAsMap["AddressesLastCheckin"] = strconv.Itoa(int(n.AddressesLastCheckin))
	// metrics.CollateMetrics("NodeInsertionsSinceLastMetricsDbg", nodeAsMap)
	// client, conn := metrics.StartConnection()
	// defer conn.Close()
	// metrics.SendMetrics(client)
	return nil
}

func AddrTrustedInsert(a *[]api.Address) error {
	if globals.BackendTransientConfig.ShutdownInitiated {
		return nil
	}
	tx, dbErr := globals.DbInstance.Beginx()
	if dbErr != nil {
		logging.LogCrash(dbErr)
	}
	for key, _ := range *a {
		(*a)[key].SetVerified(true)
		aPkIface, err := APItoDB((*a)[key], time.Now())
		if err != nil {
			// return errors.Wrap(err, "AddrTrustedInsert encountered an error in using APItoDB.")
			logging.Logf(1, "AddrTrustedInsert encountered an error in using APItoDB. Error: %#v", err)
			continue
		}
		addrPack := aPkIface.(AddressPack)
		err2 := enforceNoEmptyIdentityFields(addrPack)
		if err2 != nil {
			logging.Log(1, fmt.Sprintf("AddrTrustedInsert encountered an error in checking identity fields. Error: %#v", err2))
			continue
			// return errors.Wrap(err2, "AddrTrustedInsert encountered an error in checking identity fields")
		}
		err3 := enforceNoEmptyTrustedAddressRequiredFields(addrPack)
		if err3 != nil {
			// return errors.Wrap(err3, "AddrTrustedInsert encountered an error in checking required fields.")
			logging.Log(1, fmt.Sprintf("AddrTrustedInsert encountered an error in checking required fields. Error: %#v", err3))
			continue
		}
		_, err5 := tx.NamedExec(getSQLCommands("dbAddressUpdate")[0], addrPack.Address)
		if err5 != nil {
			logging.Log(1, err5)
		}
		if len(addrPack.Subprotocols) > 0 {
			for _, sp := range addrPack.Subprotocols {
				_, err6 := tx.NamedExec(getSQLCommands("dbSubprotocol")[0], sp)
				if err6 != nil {
					logging.Log(1, err6)
				}
			}
		}
		if len(addrPack.Junctions) > 0 {
			for _, jn := range addrPack.Junctions {
				_, err7 := tx.NamedExec(getSQLCommands("dbAddressSubprotocol")[0], jn)
				if err7 != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	err8 := tx.Commit()
	if err8 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("AddrTrustedInsert encountered an error when trying to commit to the database. Error is: %s", err8))
		return errors.New(fmt.Sprintf("AddrTrustedInsert encountered an error when trying to commit to the database.", err8))
	}
	return nil
}

// InsertOrUpdateAddresses is the multi-entry of the core function InsertOrUpdateAddress. This is the only public API, and it should be used exclusively, because this is where we have the connection retry logic that we need.
func InsertOrUpdateAddresses(a *[]api.Address) []error {
	var errs []error

	for key, _ := range *a {
		(*a)[key].SetVerified(true)
		valid, err := (*a)[key].CheckBounds()
		if err != nil {
			errs = append(errs, err)
		}
		if !valid {
			return errs
		}
	}
	err := AddrTrustedInsert(a)
	if err != nil {
		return []error{err}
	}
	return []error{}
}

func enforceNoEmptyTrustedAddressRequiredFields(obj AddressPack) error {
	if obj.Address.LocationType == 0 || obj.Address.LastSuccessfulPing == 0 || obj.Address.ProtocolVersionMajor == 0 || obj.Address.ClientVersionMajor == 0 || obj.Address.ClientName == "" || obj.Address.EntityVersion == 0 || len(obj.Subprotocols) < 1 {
		return errors.New(
			fmt.Sprintf(
				"This address has some required fields empty (One or more of: LocationType, LastSuccessfulPing, ProtocolVersionMajor, Subprotocols, ClientVersionMajor, ClientName, EntityVersion). Address: %#v\n", obj))
	}
	for _, subprot := range obj.Subprotocols {
		if subprot.Fingerprint == "" || subprot.Name == "" || subprot.VersionMajor == 0 || subprot.SupportedEntities == "" {
			return errors.New(
				fmt.Sprintf(
					"This address' subprotocol has some required fields empty (One or more of: Fingerprint, Name, VersionMajor, SupportedEntities). Address: %#v\n Subprotocol: %#v\n", obj, subprot))
		}
	}
	return nil
}

// insertOrUpdateAddress is the ONLY way to update an address in the database. Be very careful with this, careless use of this function can result in entry of untrusted data from the remotes into the local database. The only legitimate use of this is to put in the details of nodes that this local machine has personally connected to.
func insertOrUpdateAddress(a api.Address) error {
	// Because this address pack is constructed by the local machine based on the inbound TCP connection, we are marking this as trusted.
	a.SetVerified(true)
	addressPackAsInterface, err := APItoDB(a, time.Now())
	if err != nil {
		return errors.New(fmt.Sprint(
			"Error raised from APItoDB function used in Batch insert. Error: ", err))
	}
	addressPack := addressPackAsInterface.(AddressPack)

	err2 := enforceNoEmptyIdentityFields(addressPack)
	if err2 != nil {
		// If this unit does have empty identity fields, we pass on adding it to the database.
		return err2
	}
	err7 := enforceNoEmptyTrustedAddressRequiredFields(addressPack)
	if err7 != nil {
		// If this unit does have empty identity fields, we pass on adding it to the database.

		return errors.New(fmt.Sprintf("InsertOrUpdateAddress encountered an error. Error: %s", err7))
	}
	dbAddress := DbAddress{}
	var dbSubprotocols []DbSubprotocol

	var dbJunctionItems []DbAddressSubprotocol
	// Junction table.
	dbAddress = addressPack.Address // We only have one address.
	for _, dbSubprot := range addressPack.Subprotocols {
		dbSubprotocols = append(dbSubprotocols, dbSubprot)
		jItem := generateAdrSprotJunctionItem(addressPack.Address, dbSubprot)
		dbJunctionItems = append(dbJunctionItems, jItem)
	}
	tx, err3 := globals.DbInstance.Beginx()
	if err3 != nil {
		logging.Log(1, err3)
	}
	_, err4 := tx.NamedExec(getSQLCommands("dbAddressUpdate")[0], dbAddress)
	if err4 != nil {
		logging.Log(1, err4)
	}
	if len(dbSubprotocols) > 0 {
		for _, dbSubprotocol := range dbSubprotocols {
			_, err5 := tx.NamedExec(getSQLCommands("dbSubprotocol")[0], dbSubprotocol)
			if err5 != nil {
				logging.Log(1, err5)
			}
		}
	}
	if len(dbJunctionItems) > 0 {
		for _, dbJunctionItem := range dbJunctionItems {
			_, err5 := tx.NamedExec(getSQLCommands("dbAddressSubprotocol")[0], dbJunctionItem)
			if err5 != nil {
				logging.Log(1, err5)
			}
		}
	}
	err6 := tx.Commit()
	if err6 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("InsertOrUpdateAddress encountered an error when trying to commit to the database. Error is: %s", err6))
		if strings.Contains(err6.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err6
	}
	return nil
}

func generateAdrSprotJunctionItem(addr DbAddress, sprot DbSubprotocol) DbAddressSubprotocol {
	var adrSprot DbAddressSubprotocol
	adrSprot.AddressLocation = addr.Location
	adrSprot.AddressSublocation = addr.Sublocation
	adrSprot.AddressPort = addr.Port
	adrSprot.SubprotocolFingerprint = sprot.Fingerprint
	return adrSprot
}

// This is where we capture DB errors like 'DB is locked' and take action, such as retrying.
func BatchInsert(apiObjects []interface{}) (InsertMetrics, error) {
	var im InsertMetrics
	var err error
	im, err = batchInsert(&apiObjects)
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
			time.Sleep(10 * time.Second)
			logging.Log(1, "Retrying the previously failed BatchInsert transaction.")
			var err2 error
			im, err2 = batchInsert(&apiObjects)
			if err2 != nil {
				if strings.Contains(err.Error(), "Database was locked") {
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // Error is not db locked
					return im, err2
				}
			} else { // If the reattempted transaction succeeds
				logging.Log(1, "The retry attempt of the failed transaction succeeded.")
			}
		} else { // Error is not db locked
			return im, err
		}
	}
	return im, nil
}

type batchBucket struct {
	DbBoards      []DbBoard
	DbThreads     []DbThread
	DbPosts       []DbPost
	DbVotes       []DbVote
	DbKeys        []DbKey
	DbTruststates []DbTruststate
	DbAddresses   []DbAddress
	// Sub objects
	// // Parent: Board
	DbBoardOwners         []DbBoardOwner
	DbBoardOwnerDeletions []DbBoardOwner
	// // Parent: Address
	// dbSubprotocols := []DbSubprotocol{}
	// WHY? Because this is untrusted address entry, and subprotocol info coming from the environment is not committed, alongside many other parts of the address data.
}

// InsertMetrics collects how many entities in each type that ended up as candidates for DB insert. These are only candidates, and we cannot know whether an entity will make into the database - because the logic in SQL filters out duplicate entries, or entries with updates older than we have.
type InsertMetrics struct {
	BoardsReceived                 int
	BoardsDBCommitTime             float64
	BoardOwnerDBCommitTime         float64
	BoardOwnerDeletionDBCommitTime float64
	ThreadsReceived                int
	ThreadsDBCommitTime            float64
	PostsReceived                  int
	PostsDBCommitTime              float64
	VotesReceived                  int
	VotesDBCommitTime              float64
	KeysReceived                   int
	KeysDBCommitTime               float64
	TruststatesReceived            int
	TruststatesDBCommitTime        float64
	AddressesReceived              int
	AddressesDBCommitTime          float64
	MultipleInsertDBCommitTime     float64
	TimeElapsedSeconds             int
}

// Add adds an InsertMetrics object into another.
func (im *InsertMetrics) Add(im2 InsertMetrics) {
	im.BoardsReceived = im.BoardsReceived + im2.BoardsReceived
	im.BoardsDBCommitTime = im.BoardsDBCommitTime + im2.BoardsDBCommitTime
	im.BoardOwnerDBCommitTime = im.BoardOwnerDBCommitTime + im2.BoardOwnerDBCommitTime
	im.BoardOwnerDeletionDBCommitTime = im.BoardOwnerDeletionDBCommitTime + im2.BoardOwnerDeletionDBCommitTime
	im.ThreadsReceived = im.ThreadsReceived + im2.ThreadsReceived
	im.ThreadsDBCommitTime = im.ThreadsDBCommitTime + im2.ThreadsDBCommitTime
	im.PostsReceived = im.PostsReceived + im2.PostsReceived
	im.PostsDBCommitTime = im.PostsDBCommitTime + im2.PostsDBCommitTime
	im.VotesReceived = im.VotesReceived + im2.VotesReceived
	im.VotesDBCommitTime = im.VotesDBCommitTime + im2.VotesDBCommitTime
	im.KeysReceived = im.KeysReceived + im2.KeysReceived
	im.KeysDBCommitTime = im.KeysDBCommitTime + im2.KeysDBCommitTime
	im.TruststatesReceived = im.TruststatesReceived + im2.TruststatesReceived
	im.TruststatesDBCommitTime = im.TruststatesDBCommitTime + im2.TruststatesDBCommitTime
	im.AddressesReceived = im.AddressesReceived + im2.AddressesReceived
	im.AddressesDBCommitTime = im.AddressesDBCommitTime + im2.AddressesDBCommitTime
	im.MultipleInsertDBCommitTime = im.MultipleInsertDBCommitTime + im2.MultipleInsertDBCommitTime
	im.TimeElapsedSeconds = im.TimeElapsedSeconds + im2.TimeElapsedSeconds
}

// BatchInsert insert a set of objects in a batch as a transaction.
func batchInsert(apiObjectsPtr *[]interface{}) (InsertMetrics, error) {
	/*----------  Tell the frontend that we're doing something  ----------*/
	feapiconsumer.BackendAmbientStatus.DatabaseStatus = "Inserting..."
	feapiconsumer.SendBackendAmbientStatus()
	/*----------  Done  ----------*/

	apiObjects := *apiObjectsPtr
	logging.Log(2, "Batch insert starting.")
	defer logging.Log(2, "Batch insert is complete.")
	numberOfObjectsCommitted := len(apiObjects)
	logging.Log(2, fmt.Sprintf("%v objects are being committed.", numberOfObjectsCommitted))
	start := time.Now()
	insertTimestamp := time.Now() // This is used so that all entities inserted in this insert will have same LocalArrival, LastReferenced, etc. This makes our inserts atomic, single instants in time. This is to prevent the case where another node connects to you with a first sync timestamp acquired while you were inserting from another node.
	bb := batchBucket{}
	// For each API object, convert to DB object and add to transaction.
	for _, apiObject := range apiObjects {
		// apiObject: API type, dbObj: DB type.
		dbo, err := APItoDB(apiObject, insertTimestamp) // does not hit DB
		if err != nil {
			logging.Log(1,
				fmt.Sprint("Error raised from APItoDB function used in Batch insert. Error: ", err))
			if strings.Contains(err.Error(), "APItoDB only takes API (not DB) objects") { // If this triggered, this would be a programming error, in which case, we want to know.
				return InsertMetrics{}, errors.New(fmt.Sprint(
					"Error raised from APItoDB function used in Batch insert. Error: ", err))
			}
			// return InsertMetrics{}, errors.New(fmt.Sprint(
			// 	"Error raised from APItoDB function used in Batch insert. Error: ", err))
		}
		err2 := enforceNoEmptyIdentityFields(dbo) // does not hit DB
		if err2 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(2, err2)
			continue
		}
		err3 := enforceNoEmptyRequiredFields(dbo) // does not hit DB
		if err3 != nil {
			// If this unit does have empty identity fields, we pass on adding it to the database.
			logging.Log(2, err3)
			continue
		}
		switch dbObject := dbo.(type) {
		case BoardPack:
			bb.DbBoards = append(bb.DbBoards, dbObject.Board)
			// spew.Dump(dbObject.BoardOwners)
			for _, boardOwner := range dbObject.BoardOwners {
				bb.DbBoardOwners = append(bb.DbBoardOwners, boardOwner)
			}
		case DbThread:
			bb.DbThreads = append(bb.DbThreads, dbObject)
		case DbPost:
			bb.DbPosts = append(bb.DbPosts, dbObject)
		case DbVote:
			bb.DbVotes = append(bb.DbVotes, dbObject)
		case AddressPack:
			// In case of address, we strip out everything except the primary keys. This is because we cannot trust the data that is coming from the network. We just add the primary key set, and the local node will take care of directly connecting to these nodes and getting the details.

			// The other types of address inputs are not affected by this because they use InsertOrUpdateAddress, not this batch insert. If you're batch inserting addresses, it's by definition third party data.

			// This also means that we will actually be not using the Subprotocols data, as that would be untrusted data.

			dbObject.Address.LocationType = 0       // IPv4 or 6
			dbObject.Address.Type = 0               // 2 = live, 255 = static
			dbObject.Address.LastSuccessfulPing = 0 // We cannot trust someone else's lsp timestamp
			dbObject.Address.LastSuccessfulSync = 0 // We cannot trust someone else's lsc timestamp
			dbObject.Address.ProtocolVersionMajor = 0
			dbObject.Address.ProtocolVersionMinor = 0
			dbObject.Address.ClientVersionMajor = 0
			dbObject.Address.ClientVersionMinor = 0
			dbObject.Address.ClientVersionPatch = 0
			dbObject.Address.ClientName = ""
			bb.DbAddresses = append(bb.DbAddresses, dbObject.Address)
		case DbKey:
			bb.DbKeys = append(bb.DbKeys, dbObject)
		case DbTruststate:
			bb.DbTruststates = append(bb.DbTruststates, dbObject)
		default:
			return InsertMetrics{}, errors.New(
				fmt.Sprintf(
					"This object type is something batch insert does not understand. Your object: %#v\n", dbObject))
		}
	}
	im := InsertMetrics{}
	err := insert(&bb, &im)
	if err != nil {
		return InsertMetrics{}, err
	}
	im.BoardsReceived = len(bb.DbBoards)
	im.ThreadsReceived = len(bb.DbThreads)
	im.PostsReceived = len(bb.DbPosts)
	im.VotesReceived = len(bb.DbVotes)
	im.KeysReceived = len(bb.DbKeys)
	im.TruststatesReceived = len(bb.DbTruststates)
	im.AddressesReceived = len(bb.DbAddresses)
	elapsed := time.Since(start)
	im.TimeElapsedSeconds = int(elapsed.Seconds())
	clr := color.New(color.FgCyan)
	logging.Log(2, clr.Sprintf("It took %v to insert %v objects. %s", elapsed.Round(time.Millisecond), numberOfObjectsCommitted, generateInsertLog(&bb)))
	committedToDb := len(bb.DbBoards) + len(bb.DbThreads) + len(bb.DbPosts) + len(bb.DbVotes) + len(bb.DbKeys) + +len(bb.DbTruststates) + len(bb.DbAddresses)
	if (committedToDb != numberOfObjectsCommitted) && numberOfObjectsCommitted == 1 {
		clr2 := color.New(color.FgRed)
		logging.Log(1, clr2.Sprintf("There is a discrepancy between the number of entities in the inbound package, and those that end up being committed. Inbound entities count: %d, Committed to DB: %d", numberOfObjectsCommitted, committedToDb))
		logging.Log(1, clr.Sprintf("Inbound entities: %#v", apiObjects))
	}
	// if len(apiObjects) > 0 {
	// 	metrics.CollateMetrics("ArrivedEntitiesSinceLastMetricsDbg", apiObjects)
	// 	client, conn := metrics.StartConnection()
	// 	defer conn.Close()
	// 	metrics.SendMetrics(client)
	// }

	/*----------  Send sync metrics to frontend  ----------*/
	feapiconsumer.BackendAmbientStatus.DatabaseStatus = "Available"
	feapiconsumer.BackendAmbientStatus.LastInsertDurationSeconds = int32(elapsed.Seconds())
	feapiconsumer.BackendAmbientStatus.LastDbInsertTimestamp = time.Now().Unix()
	feapiconsumer.BackendAmbientStatus.DbSizeMb = int64(globals.GetDbSize())
	feapiconsumer.BackendAmbientStatus.MaxDbSizeMb = int64(globals.BackendConfig.GetMaxDbSizeMb())
	feapiconsumer.SendBackendAmbientStatus()
	/*----------  And all done!  ----------*/

	return im, nil
}

func generateInsertLog(bb *batchBucket) string {
	str := "Type:"
	if len(bb.DbBoards) > 0 {
		str = str + fmt.Sprintf(" %d Boards", len(bb.DbBoards))
	}
	if len(bb.DbThreads) > 0 {
		str = str + fmt.Sprintf(" %d Threads", len(bb.DbThreads))
	}
	if len(bb.DbPosts) > 0 {
		str = str + fmt.Sprintf(" %d Posts", len(bb.DbPosts))
	}
	if len(bb.DbVotes) > 0 {
		str = str + fmt.Sprintf(" %d Votes", len(bb.DbVotes))
	}
	if len(bb.DbKeys) > 0 {
		str = str + fmt.Sprintf(" %d Keys", len(bb.DbKeys))
	}
	if len(bb.DbTruststates) > 0 {
		str = str + fmt.Sprintf(" %d Truststates", len(bb.DbTruststates))
	}
	if len(bb.DbAddresses) > 0 {
		str = str + fmt.Sprintf(" %d Untrusted Addresses", len(bb.DbAddresses))
	}
	// Disabled because they don't count as individual entities.
	// if len(bb.DbBoardOwners) > 0 {
	// 	str = str + fmt.Sprintf(" %d Board Owners", len(bb.DbBoardOwners))
	// }
	if len(bb.DbBoards) == 0 &&
		len(bb.DbThreads) == 0 &&
		len(bb.DbPosts) == 0 &&
		len(bb.DbVotes) == 0 &&
		len(bb.DbKeys) == 0 &&
		len(bb.DbTruststates) == 0 &&
		len(bb.DbAddresses) == 0 {
		str = str + " Nothing."
	} else {
		str = str + ". Nothing else."
	}
	return str
}

func insert(batchBucket *batchBucket, im *InsertMetrics) error {
	if globals.BackendTransientConfig.ShutdownInitiated {
		return nil
	}
	// We have our final list of entries. Add these objects to DB and let DB deal with what is a new addition and what is an update.
	// (Hot code path.) Start transaction.
	start := time.Now()
	bb := *batchBucket
	var insertType []string

	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.Log(1, err)
	}
	if len(bb.DbBoards) > 0 {
		etype := "dbBoard"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbBoard := range bb.DbBoards {
				_, err := tx.NamedExec(cmd, dbBoard)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbThreads) > 0 {
		etype := "dbThread"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbThread := range bb.DbThreads {
				_, err := tx.NamedExec(cmd, dbThread)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbPosts) > 0 {
		etype := "dbPost"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbPost := range bb.DbPosts {
				_, err := tx.NamedExec(cmd, dbPost)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbVotes) > 0 {
		etype := "dbVote"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbVote := range bb.DbVotes {
				_, err := tx.NamedExec(cmd, dbVote)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbKeys) > 0 {
		etype := "dbKey"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbKey := range bb.DbKeys {
				_, err := tx.NamedExec(cmd, dbKey)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbTruststates) > 0 {
		etype := "dbTruststate"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbTruststate := range bb.DbTruststates {
				_, err := tx.NamedExec(cmd, dbTruststate)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}
	if len(bb.DbAddresses) > 0 {
		etype := "dbAddress"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbAddress := range bb.DbAddresses {
				_, err := tx.NamedExec(cmd, dbAddress)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
		// fmt.Println("Db address prune hits.")
		cmds := getSQLCommands("dbAddressPrune")
		_, err := tx.Exec(cmds[0], globals.BackendConfig.GetMaxAddressTableSize())
		if err != nil {
			logging.Log(1, err)
		}
	}
	if len(bb.DbBoardOwners) > 0 {
		etype := "dbBoardOwner"
		insertType = append(insertType, etype)
		for _, cmd := range getSQLCommands(etype) {
			for _, dbBoardOwner := range bb.DbBoardOwners {
				_, err := tx.NamedExec(cmd, dbBoardOwner)
				if err != nil {
					logging.Log(1, err)
				}
			}
		}
	}

	err2 := tx.Commit()
	if err2 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("BatchInsert encountered an error when trying to commit to the database. Error is: %s", err2))
		if strings.Contains(err.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err2
	}
	elapsed := time.Since(start)
	if len(insertType) == 1 ||
		(len(insertType) == 2 && toolbox.IndexOf("dbBoard", insertType) != -1 && toolbox.IndexOf("dbBoardOwner", insertType) != -1) { // If this is a multiple insert, I won't save the time it takes, because we don't know which part takes the most time.
		if insertType[0] == "dbBoard" {
			im.BoardsDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbThread" {
			im.ThreadsDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbPost" {
			im.PostsDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbVote" {
			im.VotesDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbKey" {
			im.KeysDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbTruststate" {
			im.TruststatesDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		} else if insertType[0] == "dbAddress" {
			im.AddressesDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		}
		// else if insertType[0] == "dbBoardOwner" {
		// 	im.BoardOwnerDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
		// }
	} else if len(insertType) > 1 {
		fmt.Println("this insert has multiple types")
		fmt.Println(insertType)
		// Multiple insert - save it to multiple insert time.
		im.MultipleInsertDBCommitTime = toolbox.Round(elapsed.Seconds(), 0.1)
	}
	return nil
}

// getSQLCommands determines the order of execution of these commands based on the order they're appended here. The fundamental rule is that all of these are gated, and gate only works in the case the object has not already inserted, so the object's actual insertion always comes last.
func getSQLCommands(dbType string) []string {
	var sqlstrs []string
	if dbType == "dbBoard" {
		sqlstrs =
			append(sqlstrs, boardInsert_BoardsBoardOwners_DeletePriors)
		sqlstrs =
			append(sqlstrs, boardInsert_BoardsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, boardInsert)
	} else if dbType == "dbThread" {
		sqlstrs =
			append(sqlstrs, threadInsert_ThreadsBoard_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, threadInsert_ThreadsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, threadInsert_ThreadsBoardsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, threadInsert)
	} else if dbType == "dbPost" {
		sqlstrs =
			append(sqlstrs, postInsert_PostsBoard_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsBoardsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsThread_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsThreadsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsPosts_Recursive_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert_PostsPostsKeys_Recursive_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, postInsert)
	} else if dbType == "dbVote" {
		sqlstrs =
			append(sqlstrs, voteInsert_VotesKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, voteInsert)
	} else if dbType == "dbKey" {
		sqlstrs =
			append(sqlstrs, keyInsert)
	} else if dbType == "dbTruststate" {
		sqlstrs =
			append(sqlstrs, truststateInsert_TruststatesKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, truststateInsert_TruststatesTargetKey_LastReferencedUpdate)
		sqlstrs =
			append(sqlstrs, truststateInsert)
	} else if dbType == "dbAddress" { // untrusted address
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstrs =
				append(sqlstrs, addressInsertMySQL)
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstrs =
				append(sqlstrs, addressInsertSQLite)
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	} else if dbType == "dbAddressUpdate" { // trusted address
		sqlstrs = append(sqlstrs, addressUpdateInsert)
	} else if dbType == "dbBoardOwner" {
		sqlstrs = append(sqlstrs, boardOwnerInsert)
	} else if dbType == "dbSubprotocol" {
		sqlstrs = append(sqlstrs, subprotocolInsert)
	} else if dbType == "dbAddressSubprotocol" {
		if globals.BackendConfig.GetDbEngine() == "mysql" {
			sqlstrs = append(sqlstrs, addressSubprotocolInsertMySQL)
		} else if globals.BackendConfig.GetDbEngine() == "sqlite" {
			sqlstrs = append(sqlstrs, addressSubprotocolInsertSQLite)
		} else {
			logging.LogCrash(fmt.Sprintf("Db Engine type not recognised."))
		}
	} else if dbType == "dbAddressPrune" {
		sqlstrs = append(sqlstrs, addressPrune)
	}
	return sqlstrs
}

// enforceNoEmptyIdentityFields enforces that nothing will enter the database without having proper identity columns. For most objects this is Fingerprint(s), for some, like address, it's a combination of multiple fields.
func enforceNoEmptyIdentityFields(object interface{}) error {
	switch obj := object.(type) {
	case BoardPack:
		if obj.Board.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This board has an empty primary key. BoardPack: %#v\n", obj))
		}
		for _, bo := range obj.BoardOwners {
			if bo.BoardFingerprint == "" || bo.KeyFingerprint == "" {
				return errors.New(
					fmt.Sprintf(
						"This board owner has one or more empty primary key(s). BoardPack: %#v\n", obj))
			}
		}
	case DbThread:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This thread has an empty primary key. Thread: %#v\n", obj))
		}
	case DbPost:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This post has an empty primary key. Post: %#v\n", obj))
		}
	case DbVote:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This vote has an empty primary key. Vote: %#v\n", obj))
		}

	case AddressPack:
		if obj.Address.Location == "" || obj.Address.Port == 0 {
			return errors.New(
				fmt.Sprintf(
					"This address has one or more empty primary key(s). Address: %#v\n", obj))
		}
		if !globals.BackendTransientConfig.AllowLocalhostRemotes {
			if obj.Address.Location == "127.0.0.1" {
				return errors.New(fmt.Sprintf("This address declares a localhost (127.0.0.1) address. Address: %#v", obj))
			}
		}
		/*
			Decided to allow URLs in — and nodes will keep connecting to URL addresses as normal.

			The problem with this is that you save the URL and make outbound requests and it will work.

			But when that node comes back at you and asks for sync, what it will give you is not its DNS domain name, but its IP address.

			Now, it could have given you the URL, but you can't trust that - it might be someone else's domain that an attacker could be trying to DDoS.

			Therefore, you will have the record of the URL, and the IP of the remote in two different address records, which will cause the remote see traffic from you more frequently. Considering that if a node has an URL, it's probably pretty established, so it's a beneficial thing that the sync happens more often.
		*/

		// if obj.Address.LocationType == 3 {
		// 	return errors.New(
		// 		fmt.Sprintf(
		// 			"This address is an URL. URL address locations are useful for initiating outgoing connections, but they are not allowed to enter the database as is — only when the IP address of the remote is given. Address: %#v\n", obj))
		// }
		// if toolbox.IsIPv6String(string(obj.Address.Location)) {
		// 	return errors.New(fmt.Sprintf("This address declares an IPv6 address. This is not supported, since large parts of the world are IPv4, and an IPv6 address is permanently unconnectable to them. An Ipv6 address holder will almost certainly have an IPv4, but an IPv4 address is not guaranteed to have IPv6 connectivity. This remote should use an Ipv4 address. Address: %#v", obj))
		// }
	case DbKey:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This key has an empty primary key. Key: %#v\n", obj))
		}
	case DbTruststate:
		if obj.Fingerprint == "" {
			return errors.New(
				fmt.Sprintf(
					"This trust state has an empty primary key. Truststate: %#v\n", obj))
		}
	}
	return nil
}

// enforceNoEmptyRequiredFields enforces that nothing will enter the database without having proper required columns. What columns are required depends on the type of the entity. See documentation for details.
func enforceNoEmptyRequiredFields(object interface{}) error {
	// FUTURE: This needs to be able to also defend against unicode replacement char or unicode rune error characters, as well as fields that are somehow only composed of spaces. That defence *should* be in place, but only real-world testing can prove it. There have been occurrences in the past where people tried to get past this by editing their own local database. The local machine assumes zero trust, everything that is coming in needs to be fully checked for sanity.
	switch obj := object.(type) {
	case BoardPack:
		if obj.Board.Name == "" ||
			obj.Board.Creation == 0 ||
			obj.Board.EntityVersion == 0 ||
			len(obj.Board.Owner) == 0 ||
			len(obj.Board.OwnerPublicKey) == 0 {
			return errors.New(
				fmt.Sprintf(
					"This board has some required fields empty (One or more of: Name, Creation, PoW, EntityVersion, Owner, OwnerPublicKey). BoardPack: %#v\n", obj))
		}
		if powEnabled() && obj.Board.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This board has the PoW field empty. BoardPack: %#v\n", obj))
		}
		for _, bo := range obj.BoardOwners {
			if bo.Level == 0 {
				return errors.New(
					fmt.Sprintf(
						"This boardowner has some required fields empty (One or more of: Level). BoardPack: %#v\n", obj))
			}
		}
	case DbThread:
		if obj.Board == "" ||
			obj.Name == "" ||
			obj.Creation == 0 ||
			obj.EntityVersion == 0 ||
			len(obj.Owner) == 0 ||
			len(obj.OwnerPublicKey) == 0 {
			return errors.New(
				fmt.Sprintf(
					"This thread has some required fields empty (One or more of: Board, Name, Creation, PoW, EntityVersion, Owner, OwnerPublicKey). Thread: %#v\n", obj))
		}
		if powEnabled() && obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This thread has the PoW field empty. Thread: %#v\n", obj))
		}
	case DbPost:
		if obj.Board == "" ||
			obj.Thread == "" ||
			obj.Parent == "" ||
			string(obj.Body) == "" ||
			obj.Creation == 0 ||
			obj.EntityVersion == 0 ||
			len(obj.Owner) == 0 ||
			len(obj.OwnerPublicKey) == 0 {
			return errors.New(
				fmt.Sprintf(
					"This post has some required fields empty (One or more of: Board, Thread, Parent, Body, Creation, PoW, EntityVersion, Owner, OwnerPublicKey). Post: %#v\n", obj))
		}
		if powEnabled() && obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This post has the PoW field empty. Post: %#v\n", obj))
		}
	case DbVote:
		if obj.Board == "" ||
			obj.Thread == "" ||
			obj.Target == "" ||
			obj.Owner == "" ||
			obj.Type == 0 ||
			obj.Creation == 0 ||
			obj.EntityVersion == 0 ||
			obj.Signature == "" ||
			len(obj.Owner) == 0 ||
			len(obj.OwnerPublicKey) == 0 {
			return errors.New(
				fmt.Sprintf(
					"This vote has some required fields empty (One or more of: Board, Thread, Target, Owner, Type, Creation, Signature, PoW, EntityVersion, Owner, OwnerPublicKey). Vote: %#v\n", obj))
		}
		if powEnabled() && obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This vote has the PoW field empty. Vote: %#v\n", obj))
		}
	case AddressPack:
		/*
			Why only these? Address is special. When address traverses over the network, it is mostly emptied out, because the information it contains is untrustable - a remote might be maliciously replacing those fields to get the network to do its bidding.
			The only address entry that is trustable is gained first-person, that is, this node connects to the node on the address, and that direct connection can update this address entity with real, first-party data.
		*/

		if obj.Address.Location == "" || obj.Address.Port == 0 || obj.Address.EntityVersion == 0 {
			return errors.New(
				fmt.Sprintf(
					"This address has some required fields empty (One or more of: Location, Port, EntityVersion. Address: %#v", obj))
		}

	case DbKey:
		if obj.Type == "" ||
			obj.PublicKey == "" ||
			obj.Creation == 0 ||
			obj.EntityVersion == 0 ||
			obj.Signature == "" {
			return errors.New(
				fmt.Sprintf(
					"This key has some required fields empty (One or more of: Type, PublicKey, Creation, PoW, Signature, EntityVersion) Key: %#v\n", obj))
		}
		if powEnabled() && obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This key has the PoW field empty. Key: %#v\n", obj))
		}
	case DbTruststate:
		if obj.Target == "" ||
			obj.Owner == "" ||
			obj.Type == 0 ||
			obj.Creation == 0 ||
			obj.EntityVersion == 0 ||
			obj.Signature == "" ||
			len(obj.Owner) == 0 ||
			len(obj.OwnerPublicKey) == 0 {
			return errors.New(
				fmt.Sprintf(
					"This trust state has some required fields empty (One or more of: Target, Owner, Type, Creation, PoW, Signature, EntityVersion, Owner, OwnerPublicKey). Truststate: %#v\n", obj))
		}
		if powEnabled() && obj.ProofOfWork == "" {
			return errors.New(
				fmt.Sprintf(
					"This truststate has the PoW field empty. Truststate: %#v\n", obj))
		}
	}
	return nil
}

/*=================================
=            Utilities            =
=================================*/

func powEnabled() bool {
	return globals.BackendTransientConfig.ProofOfWorkCheckEnabled
}

/*=====  End of Utilities  ======*/
