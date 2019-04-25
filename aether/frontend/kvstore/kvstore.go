// Frontend > KeyValueStore
// This package handles the reads from and writes to the key value store the frontend keeps fresh based on the data received from the backend. This KV store is the main source of truth for the client.

package kvstore

import (
	"aether-core/aether/frontend/search"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"github.com/asdine/storm"
	"os"
	"path/filepath"
	"strings"
	// "strconv"
)

func OpenKVStore() {
	kvdir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	kvloc := filepath.Join(kvdir, "KVStore.kv")
	/*
		Check if a search index exists. If not, we try to delete the existing KV store and start again.
	*/
	if !search.IndexExists() && kvStoreExists() {
		/*
			We want this KVStore deletion to only happen if the KVStore is present, but index is not. In that case, we silence the notifications raise, because there is no point in reshowing the prior notifications. But in the case both KvStore and index is missing, that is likely a whole new node with a JSON config imported, so we actually *do* want to show notifications in that case. Checking for KVStore existence guards for that case.

			We also have our notifications silencer gate trigger only on notification fetches with n>0 items, because the notification fetches are not in our control, and sometimes the fetch happens faster than we can actually ready the notifications. In this case, the one-time notification silencer gate is exhausted on the empty notification fetch and the actual second fetch with the whole payload comes through unsilenced.

			To prevent that, we added a n>0 condition on the extinguishment. But what that means is that it would also silence the first notification of every new user if it arrives on, since every user starts with a no kvstore and no index. That would make their first notification marked read automatically, not great. That's why this gate is kicked off on the absence of index, AND presence of kvstore. It spares new users from having their first notification extinguished.

			Gotcha: the notifications visibility order will be mangled, because we are showing notifications by generation time, not the last updated time of the underlying entity. This is because we don't want a missing item that you ended up receiving a week late because the guy went offline immediately to appear 50 lines down, it needs to appear at the top so you get a chance to respond even if late.
		*/
		deleteKVStore()
		globals.FrontendTransientConfig.SilenceNotificationsOnce = true
	}
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

func deleteKVStore() {
	kvdir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	kvloc := filepath.Join(kvdir, "KVStore.kv")
	toolbox.DeleteFromDisk(kvloc)
}

func kvStoreExists() bool {
	kvdir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	kvloc := filepath.Join(kvdir, "KVStore.kv")
	if _, err := os.Stat(kvloc); !os.IsNotExist(err) {
		return true
	}
	return false
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
