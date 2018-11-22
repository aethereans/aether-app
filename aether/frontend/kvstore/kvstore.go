// Frontend > KeyValueStore
// This package handles the reads from and writes to the key value store the frontend keeps fresh based on the data received from the backend. This KV store is the main source of truth for the client.

package kvstore

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"github.com/asdine/storm"
	"path/filepath"
	"strings"
	// "strconv"
)

func OpenKVStore() {
	kvdir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	kvloc := filepath.Join(kvdir, "KVStore.kv")
	toolbox.CreatePath(kvdir)
	kv, err := storm.Open(kvloc)
	if err != nil {
		logging.LogCrashf("Frontend KV store could not be opened. Error was: %v", err)
	}
	globals.KvInstance = kv
}

func CloseKVStore() {
	globals.KvInstance.Close()
}

// CheckKVStoreReady checks whether KV store is applying the CRUD operations correctly. If this fails, the application will crash - if the storage isn't running right, there isn't much we can do.
func CheckKVStoreReady() {
	logging.Logf(1, "Starting the Frontend KVStore readiness test.")
	type diag struct {
		ID       int `storm:"id,increment"`
		CRUDTest int
	}
	rnd := toolbox.GetInsecureRand(65535)
	d := diag{ID: 1, CRUDTest: rnd}
	err := globals.KvInstance.Save(&d)
	if err != nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to write. Error: %v", err)
	}
	d2 := diag{}
	err2 := globals.KvInstance.One("ID", 1, &d2)
	if err2 != nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to read. Error: %v", err2)
	}
	if d2.CRUDTest != rnd {
		logging.LogCrashf("KVStore readiness check failed. What we read is not what we put in.")
	}
	rnd2 := toolbox.GetInsecureRand(65535)
	err3 := globals.KvInstance.Update(&diag{ID: 1, CRUDTest: rnd2})
	if err3 != nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to update. Error: %v", err3)
	}
	d3 := diag{}
	err4 := globals.KvInstance.One("ID", 1, &d3)
	if err4 != nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to read after update. Error: %v", err4)
	}
	if d3.CRUDTest != rnd2 {
		logging.LogCrashf("KVStore readiness check failed. What we read after update is not what we put in with the update.")
	}
	err5 := globals.KvInstance.DeleteStruct(&d3)
	if err5 != nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to delete. Error: %v", err5)
	}
	d4 := diag{}
	err6 := globals.KvInstance.One("ID", 1, &d4)
	if err6 != nil && !strings.Contains(err6.Error(), "not found") {
		logging.LogCrashf("KVStore readiness check failed when attempting to read a deleted object. This should have failed with an error not found. We got a different error instead. Error: %v", err6)
	}
	if err6 == nil {
		logging.LogCrashf("KVStore readiness check failed when attempting to read a deleted object. This should have failed with an error not found. We got no error instead.")
	}
	diags := []diag{}
	globals.KvInstance.All(&diags)
	if len(diags) > 0 {
		logging.LogCrashf("KVStore readiness check failed when attempting to read deleted objects. This query should have returned no results, but it has some entries. Entries: %#v", diags)
	}
	logging.Logf(1, "Frontend KVStore is ready. Just verified by inserting and reading data successfully.")
}
