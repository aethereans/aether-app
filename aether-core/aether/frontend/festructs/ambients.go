// FEStructs > Ambients

// This is the file that contains ambients for each of the primary items, boards, threads and users.

package festructs

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"sync"
)

/*----------  Ambient entities  ----------*/
/*
  These entities are pushed into the client as part of the refresh cycle. For example, ambient boards are how we show the boards in the sidebar. Since these entities are pushed in, they need to be very small and only contain the relevant data.

  Ambient threads and users are useful in being able to autocomplete.
*/

type AmbientBoard struct {
	Fingerprint          string `storm:"id"`
	Name                 string
	LastUpdate           int64
	Notify               bool
	LastSeen             int64
	UserCount            int64
	LastNewThreadArrived int64
	// Notify and LastSeen is added using the user profile data after protobuf conversion.
}

func (c *CompiledBoard) ConvertToAmbientBoard(hasNewThreads bool, nowts int64) AmbientBoard {
	ab := AmbientBoard{
		Fingerprint: c.Fingerprint,
		Name:        c.Name,
		LastUpdate:  c.LastUpdate,
		UserCount:   int64(c.UserCount),
	}
	if hasNewThreads {
		ab.LastNewThreadArrived = nowts
	}
	return ab
}

// type AmbientBoardBatch []AmbientBoard
type AmbientBoardBatch struct {
	lock   sync.Mutex
	Boards []AmbientBoard
}

func (b *AmbientBoardBatch) UpdateBatch(abs []AmbientBoard) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for key := range abs {
		if loc := b.Find(abs[key]); loc != -1 {
			// AB already exists, update last updated timestamp
			// Heads up: board name can't change, that's why we don't have it here.
			b.Boards[loc].LastUpdate = abs[key].LastUpdate
			b.Boards[loc].LastNewThreadArrived = abs[key].LastNewThreadArrived
			b.Boards[loc].UserCount = abs[key].UserCount
		} else {
			// AB doesn't exist
			b.Boards = append(b.Boards, abs[key])
		}
	}
}

func (b *AmbientBoardBatch) Find(ab AmbientBoard) int {
	for key := range b.Boards {
		if b.Boards[key].Fingerprint == ab.Fingerprint {
			return key
		}
	}
	return -1
}

func (b *AmbientBoardBatch) Save() {
	logging.Logf(3, "Save happens in AmbientBoardBatch>Save")
	tx, err := globals.KvInstance.Begin(true)
	if err != nil {
		logging.Logf(1, "AmbientBoardBatch save failed. Error: %#v", err)
		return
	}
	defer tx.Rollback()
	for key := range b.Boards {
		err := tx.Save(&b.Boards[key])
		if err != nil {
			logging.Logf(1, "AmbientBoardBatch add board to transaction failed. Error: %#v", err)
			return
		}
	}
	err2 := tx.Commit()
	if err2 != nil {
		logging.Logf(1, "AmbientBoardBatch commit transaction failed. Error: %#v", err)
		return
	}
}

/*----------  Ambient methods  ----------*/

func GetCurrentAmbients() *AmbientBoardBatch {
	var abs []AmbientBoard
	err := globals.KvInstance.All(&abs)
	if err != nil {
		logging.Logf(1, "Existing ambient board retrieval encountered an error. Err: %v", err)
	}
	var filteredAbs []AmbientBoard
	for k := range abs {
		subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(abs[k].Fingerprint)
		if subbed {
			abs[k].Notify = notify
			abs[k].LastSeen = lastseen
			filteredAbs = append(filteredAbs, abs[k])
		}
	}
	absBatch := AmbientBoardBatch{
		Boards: filteredAbs,
		// Boards: abs, // debug
	}
	return &absBatch
}
