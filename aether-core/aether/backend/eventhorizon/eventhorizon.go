// Backend > Event Horizon
// This library handles the impermanence of all things.

/*
  The birds have vanished down the sky.
  Now the last cloud drains away.

  We sit together, the mountain and me,
  until only the mountain remains.

  — Li Bai 李白
*/

package eventhorizon

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Timestamp int64

const (
	Day = Timestamp(86400) // UNIX timestamp format. Otherwise Go's internal format isn't seconds, it's smaller.
)

func delete(ts Timestamp, entityType string) {
	tableName := ""
	switch entityType {
	case "boards":
		tableName = "Boards"
	case "threads":
		tableName = "Threads"
	case "posts":
		tableName = "Posts"
	case "votes":
		tableName = "Votes"
	case "keys":
		tableName = "Keys"
	case "truststates":
		tableName = "Truststates"
	case "addresses":
		tableName = "Addresses"
	default:
		return
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE LastReferenced < ?", tableName)
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		tx.Rollback()
		logging.Logf(1, "We couldn't begin the deletion process, transaction open failed. Error: %v", err)
		return
	}
	tx.Exec(query, ts)
	tx.Commit()
}

// func CnvToCutoffDays(days int) Timestamp {
// 	return Timestamp(time.Now().Add(-(time.Duration(days) * time.Hour * time.Duration(24))).Unix())
// }

func max(ts1 Timestamp, ts2 Timestamp) Timestamp {
	if ts1 > ts2 {
		return ts1
	}
	return ts2
}

func deleteUpToLocalMemory() {
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := Timestamp(toolbox.CnvToCutoffDays(lmD))
	// vmD := globals.BackendConfig.GetVotesMemoryDays()
	// vmCutoff := Timestamp(toolbox.CnvToCutoffDays(vmD))
	delete(lmCutoff, "boards")
	delete(lmCutoff, "threads")
	delete(lmCutoff, "posts")
	delete(lmCutoff, "keys")
	delete(lmCutoff, "truststates")
	delete(lmCutoff, "addresses")
	// These are the special ones
	delete(lmCutoff, "votes")
	// delete(vmCutoff, "votes")
	// Disables votes memory days impl. We'll reactivate if the need surfaces.
}

func deleteUpToEH(eventhorizon Timestamp) {
	delete(eventhorizon, "boards")
	delete(eventhorizon, "threads")
	delete(eventhorizon, "posts")
	delete(eventhorizon, "keys")
	delete(eventhorizon, "truststates")
	// Addresses is limited to 1000 items and it has its own cycling logic. No need to delete based on event horizon, it will likely yield not many items. The LM cutoff deletion (deleteUpToLocalMemory) does that for us.
	// delete(eventhorizon, "addresses")
}

func getDbSize() int {
	switch globals.BackendConfig.GetDbEngine() {
	case "mysql":
		query := `
      SELECT
      SUM(
        ROUND(
          ((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024 ), 2)
        )
      AS "SIZE IN MB"
      FROM INFORMATION_SCHEMA.TABLES
      WHERE TABLE_SCHEMA = "AetherDB"
    `
		var size int
		err := globals.DbInstance.Get(&size, query)
		if err != nil {
			logging.LogCrash("The attempt to read the MySQL database size failed.")
		}
		return size
	case "sqlite":
		dbLoc := filepath.Join(globals.BackendConfig.GetSQLiteDBLocation(), "AetherDB.db")
		fi, _ := os.Stat(dbLoc)
		// get the size
		// size := fi.Size() / 1048576 // Assuming 1Mb = 1048576 bytes (binary, not decimal)
		size := fi.Size() / 1000000
		return int(size)
	default:
		logging.LogCrash(fmt.Sprintf("This database type is not supported: %s", globals.BackendConfig.GetDbEngine()))
	}
	return -1 // Should never happen
}

func PruneDB() {
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := Timestamp(toolbox.CnvToCutoffDays(lmD))
	nhD := globals.BackendConfig.GetNetworkHeadDays()
	nhCutoff := Timestamp(toolbox.CnvToCutoffDays(nhD))
	tempeh := Timestamp(globals.BackendConfig.GetEventHorizonTimestamp())
	logging.Logf(2, "DbSize at the beginning of PruneDB: %v", getDbSize())
	logging.Logf(2, "Event horizon at the beginning of PruneDB: %v", time.Unix(int64(tempeh), 0).String())
	deleteUpToLocalMemory()
	if getDbSize() <= globals.BackendConfig.GetMaxDbSizeMb() {
		/*
			Below, we move event horizon one day behind, OR, if one day behind the EH goes out of the range for local memory, the local memory.

			If EH is at LM (there is no storage pressure), this will always equate to local memory. If there is, and it has caused EH to get closer to now, this will bring it back one step. This part is important, because not only we want our EH to get closer to now when there is a pressure, we also want it to move farther into history (up to LM) when the pressure is relieved. This is called backtracking.

			If DB size is below max after we remove up to local memory cutoff, set the EH to one day into the past, or, to the local memory cutoff, whichever is more recent.
		*/
		tempeh = max(tempeh-Day, Timestamp(lmCutoff))
	}
	itercount := 0
	now := Timestamp(time.Now().Unix())
	for getDbSize() > globals.BackendConfig.GetMaxDbSizeMb() {
		// First, a sanity check / guard against infinite loop.
		if tempeh > now {
			logging.Logf(1, "Something went haywire and event horizon is in the future. This might mean that the addresses table + votes table (both of which are not deleted by this loop) might have gone above the max database size threshold, and the network head is set incorrectly. Bailing. Event Horizon: %v, Now: %v", tempeh, now)
			tempeh = lmCutoff
			break
		}
		// If the user hasn't fixed the scaled mode to a setting or another,
		if !globals.BackendConfig.GetScaledModeUserSet() {
			// Check if the EH is in danger of crossing the network head threshold. If so, flip on the scaled mode.
			if tempeh+Day >= nhCutoff {
				globals.BackendConfig.SetScaledMode(true)
				logging.Log(2, "Event horizon crossed the network head. We're stopping deletion and enabling the scaled mode.")
				break // Stop deleting. We do not delete from within the network head.
			} else {
				// This triggers not because the EH moves, but because nhCutoff does with time.
				globals.BackendConfig.SetScaledMode(false)
			}
		} else {
			if tempeh+Day >= nhCutoff {
				break // We do not delete from within the network head. Force-setting the scaled mode off will make DB size grow, it won't eat into the network head.
			}
		}
		// If DB size is still larger than the max after we remove up to local memory cutoff, start deleting iteratively moving one day by day closer to now.
		// One inefficiency, if the EH is exactly at the cutoff when overflow happens, the cutoff deletion will run twice. It's OK to do that - negligible cost, and reduces complexity of logic here.
		// insert nh cutoff here.
		logging.Logf(2, "DbSize at the beginning of this iteration of PruneDB: \n%v", getDbSize())
		logging.Logf(2, "Event horizon at the beginning of this iteration of PruneDB: \n%v", time.Unix(int64(tempeh), 0).String())
		if itercount > 0 {
			// Delete up to EH first. then the next cycle delete up to eh + 1 day
			tempeh = tempeh + Day
		}
		deleteUpToEH(tempeh)
		itercount++

		logging.Logf(2, "DbSize at the end of this iteration of PruneDB: \n%v", getDbSize())
		logging.Logf(2, "Event horizon at the end of this iteration of PruneDB: \n%v", time.Unix(int64(tempeh), 0).String())

	}
	// At the end, save the new EH.
	globals.BackendConfig.SetEventHorizonTimestamp(int64(tempeh))
	logging.Logf(2, "DbSize at the end of PruneDB: %v", getDbSize())
	logging.Logf(2, "Event horizon at the end of PruneDB: %v", time.Unix(int64(tempeh), 0).String())
}
