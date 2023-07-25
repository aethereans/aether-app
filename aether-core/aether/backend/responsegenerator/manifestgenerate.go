// Backend > ResponseGenerator > ManifestGenerate
// This file contains the manifest generation logic from given main entities.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"

	// "aether-core/aether/services/randomhashgen"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	// "time"
)

type unbakedManifestItem struct {
	Fingerprint api.Fingerprint
	LastUpdate  api.Timestamp
	Page        uint64
}

func createUnbakedManifestItem(entity api.Provable, pageNo uint64) unbakedManifestItem {
	return unbakedManifestItem{
		Fingerprint: entity.GetFingerprint(),
		LastUpdate:  entity.GetLastUpdate(),
		Page:        pageNo,
	}
}

type unbakedManifestCarrier struct {
	BoardManifests      []unbakedManifestItem
	ThreadManifests     []unbakedManifestItem
	PostManifests       []unbakedManifestItem
	VoteManifests       []unbakedManifestItem
	KeyManifests        []unbakedManifestItem
	TruststateManifests []unbakedManifestItem
	AddressManifests    []unbakedManifestItem
}

// pgNoExistsInSlice looks whether a certain page number was ever created in this manifest slice. This isn't super efficient, but also not a hot path. Page numbers rarely go above 1000.
func pgNoExistsInSlice(pgNo uint64, slc *[]api.PageManifest) int {
	for i := range *slc {
		if (*slc)[i].Page == pgNo {
			return i
		}
	}
	return -1
}

// addToPageManifestSlice adds a PageManifestEntity item to the slice, creating a new page number and container as needed.
func addToPageManifestSlice(pmans *[]api.PageManifest, item *unbakedManifestItem) {
	manifestPageLoc := pgNoExistsInSlice(item.Page, pmans)
	pme := api.PageManifestEntity{
		Fingerprint: item.Fingerprint,
		LastUpdate:  item.LastUpdate,
	}
	// If this is the first time we're encountering this page no, create container for it and push it into the []
	if manifestPageLoc == -1 {
		pm := api.PageManifest{Page: item.Page}
		pm.Entities = append(pm.Entities, pme)
		*pmans = append(*pmans, pm)
	} else {
		// If we've seen this before, get that and insert
		(*pmans)[manifestPageLoc].Entities = append((*pmans)[manifestPageLoc].Entities, pme)
	}
}

func constructManifestStructure(items *[]unbakedManifestItem) *[]api.PageManifest {
	var pmans []api.PageManifest

	for i := range *items {
		addToPageManifestSlice(&pmans, &(*items)[i])
	}
	return &pmans
}

// createUnbakedManifests returns a unbakedManifestCarrier because manifest is a two-level deep entity. It looks something like page:0 > manifest manifest manifest, page:1 > manifest manifest manifest. If you use the top level page count as the page border, it is useless because they can carry arbitrary amounts of data. Instead, we need to count manifests, split to pages based on manifest counts, and THEN bundle them up to page:0 > ... structure which is the final structure.
func createUnbakedManifests(fullData *[]api.Response) *unbakedManifestCarrier {
	umc := unbakedManifestCarrier{}
	for i := range *fullData {
		for j := range (*fullData)[i].Boards {
			umc.BoardManifests = append(umc.BoardManifests, createUnbakedManifestItem(&(*fullData)[i].Boards[j], uint64(i)))
		}
		for j := range (*fullData)[i].Threads {
			umc.ThreadManifests = append(umc.ThreadManifests, createUnbakedManifestItem(&(*fullData)[i].Threads[j], uint64(i)))
		}
		for j := range (*fullData)[i].Posts {
			umc.PostManifests = append(umc.PostManifests, createUnbakedManifestItem(&(*fullData)[i].Posts[j], uint64(i)))
		}
		for j := range (*fullData)[i].Votes {
			umc.VoteManifests = append(umc.VoteManifests, createUnbakedManifestItem(&(*fullData)[i].Votes[j], uint64(i)))
		}
		for j := range (*fullData)[i].Keys {
			umc.KeyManifests = append(umc.KeyManifests, createUnbakedManifestItem(&(*fullData)[i].Keys[j], uint64(i)))
		}
		for j := range (*fullData)[i].Truststates {
			umc.TruststateManifests = append(umc.TruststateManifests, createUnbakedManifestItem(&(*fullData)[i].Truststates[j], uint64(i)))
		}
		// for j, _ := range (*fullData)[i].Addresses {
		//  umc.AddressManifests = append(umc.AddressManifests, createUnbakedManifestItem(&(*fullData)[i].Addresses[j], uint64(i)))
		// }
	}
	return &umc
}

func bakeManifests(manifestPages *[]api.ApiResponse, entityCounts *[]api.EntityCount, filters *[]api.Filter, foldername string, isPOST bool, respType string, entityType string) {
	if respType == "addresses" {
		return // addresses do not generate manifests.
	}
	protv := globals.BackendConfig.GetProtURLVersion()
	// Create directory
	var manifestdir string
	if isPOST {
		manifestdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/responses/", foldername, "/manifest")
	} else {
		if respType == "addresses" {
			manifestdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/", respType, "/", foldername, "/manifest")
		} else {
			manifestdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/c0/", respType, "/", foldername, "/manifest")
		}
	}
	toolbox.CreatePath(manifestdir)
	for key, val := range *manifestPages {
		// Add counts
		val.Caching.EntityCounts = *entityCounts
		// Add filters
		if filters != nil {
			val.Filters = *filters
		}
		val.Entity = entityType
		val.Endpoint = "manifest"
		if isPOST {
			val.Endpoint = fmt.Sprint(val.Endpoint, "_post")
		}
		// Sign
		signingErr := val.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			logging.Log(1, fmt.Sprintf("This page of the manifest of a multiple-page post response failed to be page-signed. Error: %#v Page: %#v\n", signingErr, val))
		}
		// Convert to JSON
		manifestJsonResp, err := val.ToJSON()
		if err != nil {
			logging.Log(1, fmt.Sprintf("This page of the manifest of a multiple-page post response failed to convert to JSON. Error: %#v\n", err))
		}
		// Save to disk
		filename := fmt.Sprint(key, ".json")
		saveFileToDisk(manifestJsonResp, manifestdir, filename)
	}
}
