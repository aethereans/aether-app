// Backend > ResponseGenerator > EntityCount
// This file provides a function that can count entities in a response. The main reason it's moved here is that it's a little verbose.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/configstore"
	// "aether-core/aether/services/logging"
	// "aether-core/aether/services/randomhashgen"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	// "time"
)

func countEntities(r *api.Response) *[]api.EntityCount {
	ecs := []api.EntityCount{}
	if len(r.Boards) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "board",
			Count:    len(r.Boards),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Threads) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "thread",
			Count:    len(r.Threads),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Posts) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "post",
			Count:    len(r.Posts),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Votes) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "vote",
			Count:    len(r.Votes),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Keys) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "key",
			Count:    len(r.Keys),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Truststates) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "truststate",
			Count:    len(r.Truststates),
		}
		ecs = append(ecs, ec)
	}
	if len(r.Addresses) > 0 {
		ec := api.EntityCount{
			Protocol: "c0",
			Name:     "address",
			Count:    len(r.Addresses),
		}
		ecs = append(ecs, ec)
	}
	return &ecs
}

// mergeCounts merges one configstore entity count into a list of api entitycounts.
func mergeCounts(entityCount *[]api.EntityCount, csEntityCount configstore.EntityCount) []api.EntityCount {
	ec := []api.EntityCount{}
	for _, val := range *entityCount {
		ec = append(ec, val)
	}
	if csEntityCount.Count == 0 {
		return ec
	}
	// ec := *entityCount // create a copy, don't manipulate the original
	for i, _ := range ec {
		if ec[i].Name == csEntityCount.Name {
			ec[i].Count = ec[i].Count + csEntityCount.Count
		}
	}
	return ec
}

// convertToConfigStoreEntityCount converts an api entity count slice to a configstore entity count slice.
func convertToConfigStoreEntityCount(apiec []api.EntityCount) []configstore.EntityCount {
	csec := []configstore.EntityCount{}
	for _, val := range apiec {
		csec = append(csec, configstore.EntityCount(val))
	}
	return csec
}
