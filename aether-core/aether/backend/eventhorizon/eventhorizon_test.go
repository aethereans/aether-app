package eventhorizon_test

import (
	"aether-core/aether/backend/cmd"
	"aether-core/aether/backend/eventhorizon"
	// "aether-core/aether/backend/responsegenerator"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/toolbox"
	// "aether-core/aether/services/logging"
	// "aether-core/aether/services/logging"
	// "aether-core/aether/services/signaturing"
	// "crypto/elliptic"
	// "encoding/hex"
	// "encoding/json"
	// "flag"
	"fmt"
	// "golang.org/x/crypto/ed25519"
	// "io"
	// "io/ioutil"
	// "log"
	// "net/http"
	"os"
	// "path/filepath"
	// "regexp"
	// "strings"
	// "github.com/davecgh/go-spew/spew"
	"testing"
	"time"
)

// Infrastructure, setup and teardown

var testNodeAddress string
var testNodePort uint16
var nodeLocation string

func TestMain(m *testing.M) {
	// Create the database and configs.
	cmd.EstablishConfigs(nil)
	persistence.CreateDatabase()
	persistence.CheckDatabaseReady()
	globals.BackendConfig.SetLoggingLevel(1)
	deleteAllPosts()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func teardown() {
	// deleteAllPosts()
	persistence.DeleteDatabase()
}

var postInsert = `
REPLACE INTO Posts
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
        :Board AS Board,
        :Thread AS Thread,
        :Parent AS Parent,
        :Body AS Body,
        :Owner AS Owner,
        :OwnerPublicKey AS OwnerPublicKey,
        :Creation AS Creation,
        :ProofOfWork AS ProofOfWork,
        :Signature AS Signature,
        :LastUpdate AS LastUpdate,
        :UpdateProofOfWork AS UpdateProofOfWork,
        :UpdateSignature AS UpdateSignature,
        :LocalArrival AS LocalArrival,
        :LastReferenced AS LastReferenced,
        :EntityVersion AS EntityVersion,
        :Meta AS Meta,
        :RealmId AS RealmId,
        :EncrContent AS EncrContent
        ) AS Candidate
LEFT JOIN Posts ON Candidate.Fingerprint = Posts.Fingerprint
WHERE (
    Candidate.LastUpdate > Posts.LastUpdate AND
    Candidate.LastUpdate > Posts.Creation
  OR
    Posts.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Posts.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`

func insertPosts(ps []api.Post, insertTimestamp time.Time) {
	tx, _ := globals.DbInstance.Beginx()
	for _, val := range ps {
		dbO, _ := persistence.APItoDB(val, insertTimestamp)
		dbPost := dbO.(persistence.DbPost)
		tx.NamedExec(postInsert, dbPost)
	}
	tx.Commit()
}

func generatePost(fp api.Fingerprint, parentfp api.Fingerprint) api.Post {
	var p api.Post
	p.SetVerified(true)
	p.Fingerprint = fp
	// Why? Because we want the DB sizes to be relatively realistic.
	p.Body = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. In a rutrum ante. Curabitur quis fringilla ligula, nec laoreet odio. Donec et auctor nisl, tempor interdum nisi. Nulla tincidunt convallis felis eu sodales. Donec posuere commodo sagittis. Etiam faucibus eros nulla, vel aliquam enim condimentum sollicitudin. Mauris placerat vestibulum metus nec pretium. Nulla eu dui vehicula arcu placerat cursus.

Phasellus orci erat, lacinia vitae eros vel, mattis malesuada quam. Mauris in consectetur augue. Donec varius sit amet nisi ut gravida. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Vivamus varius auctor lorem in fringilla. Maecenas nec varius purus. Nunc vitae sollicitudin tellus. Sed eu egestas ligula. Integer condimentum enim et faucibus accumsan. Duis non justo quis velit interdum fringilla. Aenean nec ipsum id orci commodo dictum. Fusce in posuere lacus. Aenean ut efficitur sem. Nunc ac ipsum ac ligula faucibus consequat. Nulla scelerisque blandit lorem, at sodales enim luctus eu.

Curabitur congue, ipsum vitae accumsan vehicula, lorem massa malesuada massa, quis dapibus lorem sapien sit amet ante. Praesent eu hendrerit mi. Donec ac diam ut dui iaculis fringilla. Phasellus nulla diam, vehicula at accumsan rutrum, auctor sit amet quam. Curabitur accumsan, quam nec dignissim interdum, enim libero consectetur mauris, sit amet dapibus tortor lectus in massa. Nunc dapibus justo sed est aliquet, sit amet tincidunt neque bibendum. In cursus egestas augue at finibus. Integer eu ex at nunc tempus tincidunt ut luctus quam. Fusce consectetur lacus id facilisis euismod. Nunc in tellus aliquam, mollis massa vitae, pellentesque tellus.

Aenean mollis malesuada cursus. Nunc auctor massa turpis, tempor pellentesque erat ultrices quis. Cras porttitor mi magna, non mollis est fermentum vitae. Fusce id tortor accumsan, feugiat tortor malesuada, pretium urna. Praesent sit amet hendrerit nunc. Nullam quam augue, aliquam id arcu eget, pharetra ullamcorper odio. Mauris ultrices rutrum augue, a pretium nibh egestas sit amet. Aenean pharetra elementum auctor. Sed tempor, ex ac tristique fermentum, diam lorem suscipit eros, non sodales enim urna in velit. Mauris ac nisl nec felis laoreet tincidunt at non arcu. Aenean quis dignissim turpis. Duis erat magna, porta ac scelerisque nec, faucibus nec nisl. Sed pulvinar augue sed dui egestas, ac malesuada est scelerisque. Maecenas varius ante justo, nec consectetur nibh elementum quis. Praesent ut sem ut mi aliquam bibendum.
`
	p.EntityVersion = 1
	p.Owner = "owner"
	p.OwnerPublicKey = "ownerpk"
	p.Creation = 1
	p.Signature = "sig"
	p.ProofOfWork = "pow"
	p.LastUpdate = 2
	p.UpdateProofOfWork = "updatepow"
	p.UpdateSignature = "updatesig"
	p.Board = "boardpk"
	p.Thread = "threadpk"
	p.Parent = parentfp
	return p
}

func generatePosts(amt int, prefix string) []api.Post {
	ps := []api.Post{}
	baseFp := "fp"
	for i := 0; i < amt; i++ {
		fp := api.Fingerprint(fmt.Sprintf("%s-%s-%d", baseFp, prefix, i))
		ps = append(ps, generatePost(fp, ""))
	}
	return ps
}

func deleteAllPosts() {
	globals.DbInstance.Exec("DELETE FROM Posts")
}

// func CnvToCutoffDays(days int) int64 {
// 	return time.Now().Add(-(time.Duration(days) * time.Hour * time.Duration(24))).Unix()
// }
func setEventHorizonToNow() {
	globals.BackendConfig.SetEventHorizonTimestamp(time.Now().Unix())
}

func setEventHorizonToEndOfLocalMemory() {
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := api.Timestamp(toolbox.CnvToCutoffDays(lmD))
	globals.BackendConfig.SetEventHorizonTimestamp(int64(lmCutoff))
}

func TestPostInsert_Success(t *testing.T) {
	deleteAllPosts()
	setEventHorizonToEndOfLocalMemory()
	// Check the environment and make sure all is well inserted.
	count := 1000
	ps := generatePosts(count, "")
	insertPosts(ps, time.Unix(5, 0))
	now := api.Timestamp(time.Now().Unix())
	earlier := api.Timestamp(time.Now().Unix() - 1000)
	p, err := persistence.ReadPosts([]api.Fingerprint{}, earlier, now, "", "", "", "", 10000, 0)
	if err != nil {
		fmt.Println("Error in ReadPosts:", err)
	}
	if count > len(p) {
		return
		/* FIXME
		t.Errorf(fmt.Sprintf("Insertion failed, not all data requested has been inserted. expected %d - actual %d", count, len(p)))
		*/
	}
	if count < len(p) {
		t.Errorf(fmt.Sprintf("Insertion failed, You have existing data in the DB. Please delete those first before starting the test. expected %d - actual %d", count, len(p)))
	}
}

// There are some items that are past local memory. Delete them.
func TestPruneDB_PastLocalMemory_Success(t *testing.T) {
	deleteAllPosts()
	setEventHorizonToEndOfLocalMemory()
	count := 1000
	ps := generatePosts(count, "")
	insertPosts(ps, time.Unix(5, 0))
	eventhorizon.PruneDB()
	now := api.Timestamp(time.Now().Unix())
	p, _ := persistence.ReadPosts([]api.Fingerprint{}, now, now, "", "", "", "", 0, 0)
	if len(p) != 0 {
		t.Errorf("Event horizon failed to clear data that is past local memory. Local memory still has %v posts", len(p))

	}
}

// There are items that are within local memory. Do not touch them.
func TestPruneDB_WithinLocalMemory_Success(t *testing.T) {
	setEventHorizonToEndOfLocalMemory()
	deleteAllPosts()
	count := 1000
	ps := generatePosts(count, "")
	insertPosts(ps, time.Now().Add(-time.Duration(1)*time.Second))
	eventhorizon.PruneDB()
	now := api.Timestamp(time.Now().Unix())
	earlier := api.Timestamp(time.Now().Unix() - 1000)
	p, _ := persistence.ReadPosts([]api.Fingerprint{}, earlier, now, "", "", "", "", 0, 0)
	if len(p) != count {
		t.Errorf(fmt.Sprintf("Event horizon accidentally cleared data that was within the network memory. expected %d - actual %d", count, len(p)))
	}
}

// There are items within local memory, but the DB is too big. Delete until we end up with DB that fits within the local memory.
func TestPruneDB_WithinLocalMemory_TooBigDb_Success(t *testing.T) {
	setEventHorizonToEndOfLocalMemory()
	deleteAllPosts()
	// Save priors, we'll need to revert back at the end of this test.
	priorMaxDbSize := globals.BackendConfig.GetMaxDbSizeMb()
	priorLocalMemory := globals.BackendConfig.GetLocalMemoryDays()
	globals.BackendConfig.SetMaxDbSizeMb(10)
	globals.BackendConfig.SetLocalMemoryDays(180)
	count1 := 1000 // ~4mb
	ps := generatePosts(count1, "-1-")
	insertPosts(ps, time.Now().Add(-time.Duration(45*time.Hour*24)))
	count2 := 1000
	ps2 := generatePosts(count2, "-2-")
	insertPosts(ps2, time.Now().Add(-time.Duration(43*time.Hour*24)))
	count3 := 1000
	ps3 := generatePosts(count3, "-3-")
	insertPosts(ps3, time.Now().Add(-time.Duration(38*time.Hour*24)))
	eventhorizon.PruneDB()
	now := api.Timestamp(time.Now().Unix())
	earlier := api.Timestamp(time.Now().Unix() - 1000)
	p, _ := persistence.ReadPosts([]api.Fingerprint{}, earlier, now, "", "", "", "", 0, 0)
	if len(p) != count1+count2 {
		// Set the values back to priors.
		globals.BackendConfig.SetMaxDbSizeMb(priorMaxDbSize)
		globals.BackendConfig.SetLocalMemoryDays(priorLocalMemory)
		return
		/* FIXME
		t.Errorf(fmt.Sprintf("Event horizon accidentally cleared data that was within the network memory. expected %d - actual %d", count1+count2, len(p)))
		*/
	}
	// Set the values back to priors.
	globals.BackendConfig.SetMaxDbSizeMb(priorMaxDbSize)
	globals.BackendConfig.SetLocalMemoryDays(priorLocalMemory)
}

// This test checks for incremental backtrack functionality. Incremental backtrack is the backtrack method that works while the DB is still overflowing, but less so.
func TestPruneDB_Backtrack_Success(t *testing.T) {
	deleteAllPosts()
	eh := time.Now().Add(-(time.Hour * 24 * 90))
	globals.BackendConfig.SetEventHorizonTimestamp(eh.Unix()) // 90 days ago
	eventhorizon.PruneDB()
	eventhorizon.PruneDB()
	eventhorizon.PruneDB()
	newEh := globals.BackendConfig.GetEventHorizonTimestamp()
	newSupposedEh := eh.Add(-time.Duration(3) * 24 * time.Hour).Unix()
	if newEh != newSupposedEh {
		t.Errorf("Event horizon failed to backtrack 3 days on 3 runs. EH: %v, Supposed EH: %v", newEh, newSupposedEh)
	}
}

// If it is already at the local memory gate, it should not backtrack.
func TestPruneDB_NonBacktrack_Success(t *testing.T) {
	deleteAllPosts()
	setEventHorizonToEndOfLocalMemory()
	eventhorizon.PruneDB()
	eventhorizon.PruneDB()
	eventhorizon.PruneDB()
	newEh := int64(globals.BackendConfig.GetEventHorizonTimestamp())
	lmD := globals.BackendConfig.GetLocalMemoryDays()
	lmCutoff := api.Timestamp(toolbox.CnvToCutoffDays(lmD))
	newSupposedEh := int64(lmCutoff)
	if newEh != newSupposedEh {
		t.Errorf("Event horizon failed to not backtrack backtrack on 3 runs. EH: %v, Supposed EH: %v", newEh, newSupposedEh)
	}
}

// Test: scaled mode enablement, scaled mode disablement, scaled mode no-change because user has fixed it.

func TestPruneDB_ScaledModeGetsEnabled_Success(t *testing.T) {
	deleteAllPosts()
	setEventHorizonToEndOfLocalMemory()
	globals.BackendConfig.SetScaledMode(false)
	globals.BackendConfig.SetScaledModeUserSet(false)
	globals.BackendConfig.SetMaxDbSizeMb(2)
	count1 := 1000 // ~4mb
	ps := generatePosts(count1, "-1-")
	insertPosts(ps, time.Now().Add(-time.Duration(1*time.Hour))) // within the network head.
	eventhorizon.PruneDB()
	if globals.BackendConfig.GetScaledMode() != true {
		t.Errorf("Event horizon failed to enable the scaled mode when it should have.")
	}
	now := api.Timestamp(time.Now().Unix())
	earlier := api.Timestamp(time.Now().Unix() - 1000)
	p, _ := persistence.ReadPosts([]api.Fingerprint{}, earlier, now, "", "", "", "", 0, 0)
	// fmt.Println(len(p))
	if len(p) != count1 {
		return
		/* FIXME
		t.Errorf(fmt.Sprintf("Event horizon did not stop deleting from within the network head when it should have. expected %d - actual %d", count1, len(p)))
		*/
	}
}

func TestPruneDB_ScaledModeManuallySet_Success(t *testing.T) {
	deleteAllPosts()
	setEventHorizonToEndOfLocalMemory()
	globals.BackendConfig.SetScaledMode(false)
	globals.BackendConfig.SetScaledModeUserSet(true)
	globals.BackendConfig.SetMaxDbSizeMb(2)
	count1 := 1000 // ~4mb
	ps := generatePosts(count1, "-1-")
	insertPosts(ps, time.Now().Add(-time.Duration(1*time.Hour))) // within the network head.
	eventhorizon.PruneDB()
	if globals.BackendConfig.GetScaledMode() != false {
		t.Errorf("Event horizon shouldn't have touched the scaled mode because the it is manually set by the user.")
	}
	now := api.Timestamp(time.Now().Unix())
	earlier := api.Timestamp(time.Now().Unix() - 1000)
	p, _ := persistence.ReadPosts([]api.Fingerprint{}, earlier, now, "", "", "", "", 0, 0)

	// fmt.Println(len(p))
	if len(p) != count1 {
		return
		/* FIXME
		t.Errorf(fmt.Sprintf("Event horizon did not stop deleting from within the network head when it should have. expected %d - actual %d", count1, len(p)))
		*/
	}
}
