// Backend > ResponseGenerator > CacheGenerator
// This file provides a set of functions that relate to the generation of pregenerated caches in certain intervals.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	// "aether-core/aether/services/configstore"
	"aether-core/aether/backend/feapiconsumer"
	"aether-core/aether/services/extverify"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/randomhashgen"
	"aether-core/aether/services/toolbox"
	"encoding/json"
	"errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CacheResponse is the internal procesing structure for generating caches to be saved to the disk.
type CacheResponse struct {
	cacheName     string
	start         api.Timestamp
	end           api.Timestamp
	entityPages   *[]api.Response
	indexPages    *[]api.Response
	manifestPages *[]api.Response
	counts        *[]api.EntityCount
}

// GatherCacheData responds to a cache generation request. This returns an Api.Response entity with entities, entity indexes, and the cache link that needs to be inserted into the index of the endpoint.
// This has no filters.
func GatherCacheData(etype string, start api.Timestamp, end api.Timestamp) (CacheResponse, error) {
	var cacheRespStruct CacheResponse
	switch etype {
	case "boards", "threads", "posts", "votes", "keys", "truststates":
		localData, dbError := persistence.Read(etype, []api.Fingerprint{}, []string{}, start, end, true, nil)
		if dbError != nil {
			return cacheRespStruct, errors.New(fmt.Sprintf("This cache generation request caused an error in the local database while trying to respond to this request. Error: %#v\n", dbError))
		}
		// if len(localData.Boards) == 0 &&
		// 	len(localData.Threads) == 0 &&
		// 	len(localData.Posts) == 0 &&
		// 	len(localData.Votes) == 0 &&
		// 	len(localData.Keys) == 0 &&
		// 	len(localData.Truststates) == 0 {
		/*
		   There's no data in this result. But the cache generation should continue. Why?

		   1) This cache generation process is guarded by the 'is this node tracking head?' guard. So this part of the code does not need to care about accidentally generating blank caches.

		   2) Consider the case that the most recent data in the network is actually, genuinely three days old. Had we stopped cache generation when empty, the caches for those two blank days would NEVER be generated, but ALWAYS attempted. So every hit of the cache generation cycle would turn out to be an attempt for the cache generation of those two days.

		   How do I know? Because that's exactly what happened and this text is the bug fix.
		*/
		// 	logging.Log(2, fmt.Sprintf("The result for this cache is empty. Entity type: %s, Start: %d, End: %d", etype, start, end))
		// }
		/*
			Disabling above - we don't really care whether caches are empty. We just generate them all the same, for now.
		*/
		// if len(localData.Truststates) > 0 {
		// 	logging.Logf(1, "We've found some truststates: %v", localData.Truststates)
		// }
		entityPages := splitEntitiesToPages(&localData)
		cacheRespStruct.entityPages = entityPages
		// create indexes
		indexes := createUnbakedIndexes(entityPages)
		indexPages := splitEntitiesToPages(indexes)
		cacheRespStruct.indexPages = indexPages
		// fmt.Println("length of index pages")
		// fmt.Println(len(*cacheRespStruct.indexPages))
		// create manifests
		manifest := createUnbakedManifests(entityPages)
		manifestPages := splitManifestToPages(manifest)
		cacheRespStruct.manifestPages = manifestPages
		// fmt.Println("length of manifest pages")
		// fmt.Println(len(*cacheRespStruct.manifestPages))
		// count entities
		entityCounts := countEntities(&localData)
		cacheRespStruct.counts = entityCounts
		cn, err := randomhashgen.GenerateInsecureRandomHash()
		if err != nil {
			return cacheRespStruct, errors.New(fmt.Sprintf("There was an error in the cache generation request serving. Error: %#v\n", err))
		}
		cacheRespStruct.cacheName = cn
		cacheRespStruct.start = start
		cacheRespStruct.end = end

	case "addresses":
		addresses, dbError := persistence.ReadAddresses("", "", 0, start, end, 0, 0, 0, "timerange_all") // Cache generation only generates caches for addresses that this computer has personally connected to.
		if dbError != nil {
			return cacheRespStruct, errors.New(fmt.Sprintf("This cache generation request caused an error in the local database while trying to respond to this request. Error: %#v\n", dbError))
		}
		addresses = *sanitiseOutboundAddresses(&addresses)
		if len(addresses) == 0 {
			/*
			   There's no data in this result. But the cache generation should continue. Why?

			   1) This cache generation process is guarded by the 'is this node tracking head?' guard. So this part of the code does not need to care about accidentally generating blank caches.

			   2) Consider the case that the most recent data in the network is actually, genuinely three days old. Had we stopped cache generation when empty, the caches for those two blank days would NEVER be generated, but ALWAYS attempted. So every hit of the cache generation cycle would turn out to be an attempt for the cache generation of those two days.

			   How do I know? Because that's exactly what happened and this text is the bug fix.
			*/
			logging.Log(2, fmt.Sprintf("The result for this cache is empty. Entity type: %s", etype))
		}
		cacheRespStruct.start = start
		cacheRespStruct.end = end
		var localData api.Response
		localData.Addresses = addresses
		entityPages := splitEntitiesToPages(&localData)
		cacheRespStruct.entityPages = entityPages
		cn, err := randomhashgen.GenerateInsecureRandomHash()
		if err != nil {
			return cacheRespStruct, errors.New(fmt.Sprintf("There was an error in the cache generation request serving. Error: %#v\n", err))
		}
		cacheRespStruct.cacheName = cn
		// count entities
		entityCounts := countEntities(&localData)
		cacheRespStruct.counts = entityCounts
	default:
		return cacheRespStruct, errors.New(fmt.Sprintf("The requested entity type is unknown to the cache generator. Entity type: %s", etype))
	}
	return cacheRespStruct, nil
}

func updateEndpointIndex(cacheIndex *api.ApiResponse, cacheData *CacheResponse) {
	// Save the cache link into the index.
	c := createResultCacheBlockForIndex(cacheData)
	cacheIndex.Results = append(cacheIndex.Results, c)
	cacheIndex.Timestamp = api.Timestamp(int64(time.Now().Unix()))
	cacheIndex.Caching.Pregenerated = true
}

func createResultCacheBlockForIndex(cacheData *CacheResponse) api.ResultCache {
	// Save the cache link into the index.
	var c api.ResultCache
	c.ResponseUrl = fmt.Sprintf("cache_%s", cacheData.cacheName)
	c.StartsFrom = cacheData.start
	c.EndsAt = cacheData.end
	return c
}

/*
	This deletion logic does three things:
	1) Delete caches that have gone past the network head threshold
	2) Delete caches that we do no longer carry in the index payload. (i.e. consolidated caches' non-consolidated old versions.) We do this trailing caches by one hour, to allow anyone that is still downloading from the cache to finish downloading.
*/
func deleteTooOldCaches(etype string, cacheIndex *api.ApiResponse) {
	logging.Logf(2, "DeleteTooOldCaches starts to run for Entity type: %v", etype)
	entityCacheDir, err := generateEndpointDir(etype)
	if err != nil {
		logging.Logf(1, "deleteTooOldCaches errored out when trying to generate the entity cache directory for this entity type. Err: %v", err)
		return
	}
	threshold := api.Timestamp(time.Now().Add(
		-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
	oldestCacheEnd := api.Timestamp(time.Now().Unix())
	for _, cache := range cacheIndex.Results {
		if oldestCacheEnd > cache.EndsAt {
			oldestCacheEnd = cache.EndsAt
		}
	}
	if threshold > oldestCacheEnd {
		// We have more caches than needed. We need to delete some starting from the oldest.
		logging.Log(2, fmt.Sprintf("We have caches for a longer duration of time than we need. (The oldest cache.EndsAt is %d, the threshold is %d) Caches will be purged starting from the oldest. Purge is starting.", oldestCacheEnd, threshold))
		oldCaches := []api.ResultCache{}
		stillValidCaches := []api.ResultCache{}
		for _, cache := range cacheIndex.Results {
			if cache.EndsAt < threshold {
				// This cache is too old.
				oldCaches = append(oldCaches, cache)
				continue
			}
			stillValidCaches = append(stillValidCaches, cache)
		}
		cacheIndex.Results = stillValidCaches
		for _, cache := range oldCaches {
			// Figure out the location of the cache and nuke it.
			location := filepath.Join(entityCacheDir, cache.ResponseUrl)
			os.RemoveAll(location)
		}
		logging.Log(1, fmt.Sprintf("Old cache purging is complete. We've deleted these expired caches from both index and from the local file system: %#v", oldCaches))
	}
	// Now, we want to take a look at the folder structure, and see which ones need to go.
	validsMap := make(map[string]bool)
	for k, _ := range cacheIndex.Results {
		validsMap[cacheIndex.Results[k].ResponseUrl] = true
	}
	folders, err := ioutil.ReadDir(entityCacheDir)
	if err != nil {
		logging.Logf(1, "deleteTooOldCaches had an error trying to read the cache directory. Dir: %v, Error: %v", entityCacheDir, err)
		return
	}
	for k, _ := range folders {
		cacheName := folders[k].Name()
		if cacheName == "index.json" {
			// This folder not only has cache folders. Avoid that one.
			continue
		}
		cachePath := filepath.Join(entityCacheDir, cacheName)
		deletionMarker := filepath.Join(cachePath, "delete_in_next_cycle.txt")
		// First, check if there is a deletion marker. If so, delete.
		deletionMarkerExists, err := fileExists(deletionMarker)
		if err != nil {
			logging.Logf(1, "deleteTooOldCaches encountered an error while trying to check for the existence of the deletion marker. Err: %v", err)
			continue
		}
		if deletionMarkerExists {
			logging.Logf(2, "We found a deletion marker on cache %v, so we're deleting it.", cacheName)
			os.RemoveAll(cachePath)
			continue
		}
		// Second, check if it's still in circulation. If so, don't delete.
		logging.Logf(2, "For etype: %v, this is the cache name we're checking the valids map for: %v", etype, cacheName)
		if validsMap[cacheName] {
			logging.Logf(2, "Cache %v exists in the valids map.", cacheName)
			// If in the map = still in use. Don't delete.
			continue
		}
		logging.Logf(2, "%v %v does not exist in the valids map and is queued for deletion. ", etype, cacheName)
		// If it made it to here, It is not in use. Insert the deletion marker, so we can delete it in the next cycle.
		markCacheForDeletion(cachePath)
	}
	logging.Logf(1, "DeleteTooOldCaches is done for Entity type: %v", etype)
	// logging.Logf(1, "Expired cache deletion is complete.")
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func markCacheForDeletion(cachePath string) {
	marker := filepath.Join(cachePath, "delete_in_next_cycle.txt")
	err := ioutil.WriteFile(marker, []byte{}, 0644)
	if err != nil {
		logging.Logf(1, "markCacheForDeletion encountered an error. Err: %v", err)
	}
	logging.Logf(2, "Delete in next cycle marker inserted to cachepath: %v", cachePath)
}

func generateEndpointDir(etype string) (string, error) {
	protv := globals.BackendConfig.GetProtURLVersion()
	var ecd string
	if etype == "boards" ||
		etype == "threads" ||
		etype == "posts" ||
		etype == "votes" ||
		etype == "keys" ||
		etype == "truststates" {
		ecd = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/c0/", etype)
	} else if etype == "addresses" {
		ecd = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/", etype)
	} else {
		return ecd, errors.New(fmt.Sprintf("Unknown response type: %s", etype))
	}
	return ecd, nil
}

func generateEndpointRelativePath(etype string) (string, error) {
	protv := globals.BackendConfig.GetProtURLVersion()
	var ecd string
	if etype == "boards" ||
		etype == "threads" ||
		etype == "posts" ||
		etype == "votes" ||
		etype == "keys" ||
		etype == "truststates" {
		ecd = filepath.Join(protv, "c0", etype)
	} else if etype == "addresses" {
		ecd = filepath.Join(protv, etype)
	} else {
		return ecd, errors.New(fmt.Sprintf("Unknown response type: %s", etype))
	}
	return ecd, nil
}

// CreateNewCache creates the cache for the given entity type for the given time range.
func CreateNewCache(etype string, start api.Timestamp, end api.Timestamp, allPriorCachesGeneratedSoFarAreEmpty bool) (bool, error) {
	// - Pull the data from the DB
	// - Look at the cache folder. If there is a cache folder and an index there, save the cache and add to index.
	// - If there is no cache present there, create the index and add it as the first entry.
	// fmt.Printf("CreateNewCache was asked to generate a cache for the resp type %#v that ended at the timestamp: %#v\n", etype, end)
	cacheData, err := GatherCacheData(etype, start, end)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err))
	}
	// if (*cacheData.entityPages)[0].Empty() && allPriorCachesGeneratedSoFarAreEmpty {
	// 	// fmt.Printf("This cache and all prior caches generated so far were empty, skipping generation of this cache. Entity type: %s, Start: %d, End: %d\n", etype, start, end)
	// 	return true, nil
	// }
	/*
		^ This is an optimisation that we should do eventually, but for now, it's easier to reason about when they don't actually prevent generation.
	*/
	ePagesApiresp := convertResponsesToApiResponses(cacheData.entityPages)
	iPagesApiresp := convertResponsesToApiResponses(cacheData.indexPages)
	mPagesApiresp := convertResponsesToApiResponses(cacheData.manifestPages)
	startAsString := strconv.FormatInt(int64(cacheData.start), 10)
	endAsString := strconv.FormatInt(int64(cacheData.end), 10)
	filter := api.Filter{Type: "timestamp", Values: []string{startAsString, endAsString}}
	generateContainer(ePagesApiresp, iPagesApiresp, mPagesApiresp, cacheData.counts, &[]api.Filter{filter}, cacheData.cacheName, false, etype, api.Timestamp(cacheData.start))
	// Generate endpoint index.
	epd, err := generateEndpointDir(etype)
	if err != nil {
		return false, err
	}
	toolbox.CreatePath(epd)
	var endpointIndex api.ApiResponse
	// Look for the index.json in it. If it doesn't exist, create.
	// Heads up: we're reading and parsing our own caches.
	endpointIndexAsJson, err3 := ioutil.ReadFile(fmt.Sprint(epd, "/index.json"))
	if err3 != nil && strings.Contains(err3.Error(), "no such file or directory") {
		// The index.json of this cache likely doesn't exist. Create one.
		endpointIndex.Prefill()
		endpointIndex.Entity = etype
		endpointIndex.Endpoint = etype
	} else if err3 != nil {
		// The index is corrupted. The user knowingly modified it or filesystem did, or some other process did.
		//FUTURE: We should regenerate this cache, maybe. But if the user (or a process running as user) modified this cache, we have no guarantee that it will not do that again in the future, so regenerating it might just be a waste of resources.
		return false, errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err3))
	} else {
		// err3 is nil
		json.Unmarshal(endpointIndexAsJson, &endpointIndex)
	}
	// If the file exists, go through with regular processing.
	updateEndpointIndex(&endpointIndex, &cacheData)
	deleteTooOldCaches(etype, &endpointIndex)
	signingErr := endpointIndex.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if signingErr != nil {
		return false, errors.New(fmt.Sprintf("This entity index failed to be page-signed. Error: %#v Page: %#v\n", signingErr, endpointIndex))
	}
	json, err4 := endpointIndex.ToJSON()
	if err4 != nil {
		return false, err
	}
	saveFileToDisk(json, epd, "index.json")
	return false, nil
}

// CreateNewCache creates the cache for the given entity type for the given time range.
func CreateNewCacheV2(etype string, start, end api.Timestamp) (cachegenSkipped bool, resultCache api.ResultCache, err error) {
	// logging.Logf((1, "CreateNewCache was asked to generate a cache for the resp type %#v that ended at the timestamp: %#v\n", etype, end)
	cacheData, err := GatherCacheData(etype, start, end)
	if err != nil {
		return false, api.ResultCache{}, errors.New(fmt.Sprintf("Cache creation process encountered an error. Error: %s", err))
	}
	ePagesApiresp := convertResponsesToApiResponses(cacheData.entityPages)
	iPagesApiresp := convertResponsesToApiResponses(cacheData.indexPages)
	mPagesApiresp := convertResponsesToApiResponses(cacheData.manifestPages)
	startAsString := strconv.FormatInt(int64(cacheData.start), 10)
	endAsString := strconv.FormatInt(int64(cacheData.end), 10)
	filter := api.Filter{Type: "timestamp", Values: []string{startAsString, endAsString}}
	generateContainer(ePagesApiresp, iPagesApiresp, mPagesApiresp, cacheData.counts, &[]api.Filter{filter}, cacheData.cacheName, false, etype, api.Timestamp(cacheData.start))
	rc := createResultCacheBlockForIndex(&cacheData)
	return false, rc, nil
}

/*
  Methods related to cache days table generation. Cache Days Table is a table of days with beginning and end timestamps that we feed into the cache generator to generate caches for those days.

  We then feed this cache generation table into our cache generator, and it creates the appropriate folder structure for us.
*/

// readEndpointIndex reads the cache index of the requested endpoint from the local drive. This is then used for finding the end timestamp of the last cache generated.
func readEndpointIndex(etype string) (api.ApiResponse, error) {
	protv := globals.BackendConfig.GetProtURLVersion()
	var cacheDir string
	if etype == "boards" || etype == "threads" || etype == "posts" || etype == "votes" || etype == "keys" || etype == "truststates" {
		cacheDir = globals.BackendConfig.GetCachesDirectory() + "/" + protv + "/c0/" + etype
	} else if etype == "addresses" {
		cacheDir = globals.BackendConfig.GetCachesDirectory() + "/" + protv + "/" + etype
	}
	cacheIndex := cacheDir + "/index.json"
	dat, err := ioutil.ReadFile(cacheIndex)
	if err != nil {
		return api.ApiResponse{}, err
	}
	var apiresp api.ApiResponse
	err2 := json.Unmarshal([]byte(dat), &apiresp)
	if err2 != nil {
		logging.Log(2, fmt.Sprintf(fmt.Sprintf(
			"The JSON That was the cache index for the entity type is malformed. Entity type: %s, JSON: %s", etype, string([]byte(dat)))))
		// Delete the whole index folder and return 0 to generate new caches.
		os.RemoveAll(cacheDir)
		return api.ApiResponse{}, errors.New("no such file or directory")
	}
	return apiresp, nil
}

func saveEndpointIndex(epi api.ApiResponse, etype string) error {
	signingErr := epi.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if signingErr != nil {
		logging.Logf(1, "saveEndpointIndex could not sign this entity index: failed to be page-signed. Error: %#v Page: %#v\n", signingErr, epi)
		return errors.New(fmt.Sprintf("This entity index failed to be page-signed. Error: %#v Page: %#v\n", signingErr, epi))
	}
	epiJson, err := epi.ToJSON()
	if err != nil {
		logging.Logf(1, "saveEndpointIndex could not convert the endpoint index to JSON. Err: %v", err)
		return errors.New(fmt.Sprintf("saveEndpointIndex could not convert the endpoint index to JSON. Err: %v", err))
	}
	endpointDir, err := generateEndpointDir(etype)
	if err != nil {
		logging.Logf(1, "saveEndpointIndex could not generate the endpoint path. Err: %v", err)
		return errors.New(fmt.Sprintf("saveEndpointIndex could not generate the endpoint path. Err: %v", err))
	}
	saveFileToDisk(epiJson, endpointDir, "index.json")
	return nil
}

// determineLastCacheEnd figures out when was the last cache for this entity type was generated. For each entity, we need to look at the last cache that is generated by the entity and find its end timestamp.
func determineLastCacheEnd(etype string) api.Timestamp {
	cacheIndex, err := readEndpointIndex(etype)
	if err != nil {
		// logging.LogCrash(err)
		// FUTURE: Add tampered caches gating
		if strings.Contains(err.Error(), "no such file or directory") {
			logging.Log(2, fmt.Sprintf("The cache for this entity type does not exist yet. We'll be generating this from scratch. Entity type: %#v", etype))
			// var blankTs api.Timestamp
			// return blankTs
		} else {
			logging.Logf(1, "determineLastCacheEnd errored out. Error: %v", err)
		}
	}
	// Identify the most recent end timestamp
	var mostRecentExtantCacheEndTs api.Timestamp
	for _, cache := range cacheIndex.Results {
		if cache.EndsAt > mostRecentExtantCacheEndTs {
			mostRecentExtantCacheEndTs = cache.EndsAt
		}
	}
	// If the most recent extant cache end is lesser than our network head, make it network head (so that we don't have to generate caches starting from 1970s)
	networkHeadThreshold := api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
	if mostRecentExtantCacheEndTs < networkHeadThreshold {
		mostRecentExtantCacheEndTs = networkHeadThreshold
	}
	// fmt.Printf("Most recent extant cache end for type %#v is %#v.\n", etype, mostRecentExtantCacheEndTs)
	return mostRecentExtantCacheEndTs
}

// generateRequestedCachesTable determines how many caches we need to generate, and at which intervals they need to start and end.
func generateRequestedCachesTable(mostRecentExtantCacheEndTs api.Timestamp) []api.ResultCache {
	// Split the difference of most recent cache end and now into 24H slices.
	now := api.Timestamp(time.Now().Unix())
	var cachesTable []api.ResultCache
	currentEndTs := mostRecentExtantCacheEndTs
	// So long as the current end + a day is lesser than timestamp of now, iterate
	for currentEndTs < now {
		// fmt.Println("current end ts smaller than now")
		newEnd := api.Timestamp(time.Unix(int64(currentEndTs), 0).Add(time.Duration(globals.BackendConfig.GetCacheDurationHours()) * time.Hour).Unix())
		cache := api.ResultCache{
			StartsFrom: currentEndTs,
			EndsAt:     newEnd,
		}
		currentEndTs = newEnd
		cachesTable = append(cachesTable, cache)
	}
	if len(cachesTable) == 0 {
		logging.LogCrash(fmt.Sprintf("Cache table length turned out to be zero. Your Time block size and past blocks to check values are invalid. If you have not changed them, please delete the backend configuration file and restart the application."))
	}
	/*
		  After this table generation is done, check the last cache bracket (start>end). If the time difference of its start and now() is less than 12 hours, delete the last bracket, and set the n-1th cache bracket's end timestamp to now.
			   e.g. (assuming the caches are generated every day)
			   IF:
			    Day -3  Day -2  Day -1     Now
			   |-------|-------|-------|----=--|
			   (more than half day's worth data in last)

			   DO:
			    Day -3  Day -2  Day -1     Now
			   |-------|-------|-------|----=|
			   (move the end to now)

			   IF:
			    Day -3  Day -2  Day -1   Now
			   |-------|-------|-------|--=----|
			   (less than 12 hours of data in last)

			   DO:
			    Day -3  Day -2  Day -1  Now
			   |-------|-------|---------=|
			   (remove the last one and move the end of n-1 to now)
	*/
	lastDTItemEndTs := cachesTable[len(cachesTable)-1].EndsAt
	halfCacheIntoFuture := api.Timestamp(time.Now().Add(time.Duration(globals.BackendConfig.GetCacheDurationHours()) / 2 * time.Hour).Unix())
	if halfCacheIntoFuture < lastDTItemEndTs && len(cachesTable) > 1 {
		// The last cache covers less than 12 hours (i.e. it captures more than 12 hours of not-happened-yet) AND it's not just one item in the cache table (in which case, don't do anything.)
		// Chop the last item off.
		cachesTable = cachesTable[:len(cachesTable)-1]
	}
	// Make the last item of the cache table come up to now, not to future.
	cachesTable[len(cachesTable)-1].EndsAt = api.Timestamp(time.Now().Unix())
	// fmt.Printf("Caches table length is: %#v\n", len(cachesTable))
	// fmt.Printf("Caches table is: %#v\n", cachesTable)
	return cachesTable
}

// GenerateCachedEndpoint determines how many caches we will need to create for a given entity type, and generates them. This is a separate function because different endpoint can have different last cache ends.
func GenerateCachedEndpointV1(etype string) int64 {
	currentCacheEnd := int64(0)
	// ^ What are are going to set the last cachegen timestamp
	// Read the end of the last cache, or if there are none, start from the beginning.
	lastCacheEndTs := determineLastCacheEnd(etype)
	// If the lastCacheEndTs is younger than globals.BackendConfig.GetCacheGenerationIntervalHours()-1) hours, we do nothing. The cache generator cycle will attempt to create a cache every hour, so this is where we gate how often we create caches.

	// If last cache end is more than 1 hour ago
	cachegenThreshold := api.Timestamp(
		time.Now().Add(-time.Duration(globals.BackendConfig.GetCacheGenerationIntervalHours()-1) * time.Hour).Unix())
	// fmt.Println("Cachegen threshold: ", cachegenThreshold)
	// fmt.Println("Last cache end TS: ", lastCacheEndTs)
	if cachegenThreshold > lastCacheEndTs {
		cachesTable := generateRequestedCachesTable(lastCacheEndTs)
		currentCacheEnd = int64(cachesTable[len(cachesTable)-1].EndsAt)
		// ^ We have caches to generate. The end of our last cache is going to be our last cache generation timestamp.
		allPriorCachesGeneratedSoFarAreEmpty := true
		for _, val := range cachesTable {
			if len(val.ResponseUrl) > 0 {
				continue
				/*
					In V1 cache generation, this ResponseURL is always empty, because the cachesTable provided only starts from the end of the last cache, and ends with now.

					In the V2 cache generation, there is a consolidation pass that merges past caches as they become older and older. However, this pass tries to retain as many caches intact as possible. Therefore there is a chance that some caches survive it intact, especially the ones that hasn't grown to the next cache size yet. (Ex: you have 9 10 minute caches, when the 10th is added, that becomes one 100 minute cache). In that case, we do not regenerate those intact caches.

					If this is the case, this has been grabbed from a prior cache list as a consolidation intact-survivor. We don't generate those caches.
				*/
			}
			empty, err := CreateNewCache(etype, val.StartsFrom, val.EndsAt, allPriorCachesGeneratedSoFarAreEmpty)
			if err != nil {
				logging.Log(2, err)
			}
			if !empty {
				allPriorCachesGeneratedSoFarAreEmpty = false
			}
		}
	} else {
		currentCacheEnd = int64(lastCacheEndTs)
		// ^ We have no caches to generate, so we should reuse the last cache end.
		logging.Log(3, fmt.Sprintf("Last cache that was created for %s was newer than %d hours ago. Please wait until after.", etype, globals.BackendConfig.GetCacheDurationHours()-1))
	}
	return currentCacheEnd
}

/*=============================================================================
=            V2 CACHING FOR EXPONENTIAL GENERATION & CONSOLIDATION            =
=============================================================================*/

func generateRequestedCachesTableV2(etype string, mostRecentExtantCacheEndTs api.Timestamp) []api.ResultCache {
	// Split the difference of most recent cache end and now into 24H slices.
	now := api.Timestamp(time.Now().Unix())
	currentEndTs := mostRecentExtantCacheEndTs
	newtt := NewCacheTimeTable(currentEndTs, now, cacheTimeBlocks)
	// Get the currently present cache index.
	extanttt := []api.ResultCache{}
	cacheIndex, err := readEndpointIndex(etype)
	if err != nil {
		logging.Logf(1, "Read cache index errored out we'll regenerate every cache from scratch. Err: %v", err)
		// todo
	} else {
		extanttt = cacheIndex.Results
	}
	mergedtt := append(extanttt, newtt...)
	consolidatedtt := MakeConsolidatedTimeTable(&mergedtt, cacheTimeBlocks)
	return consolidatedtt
}

func GenerateCachedEndpointV2(etype string) int64 {
	logging.Logf(2, "GenerateCachedEndpointV2 starting to run. Endpoint: %v", etype)
	currentCacheEnd := int64(0)
	// ^ What are are going to set the last cachegen timestamp
	// Read the end of the last cache, or if there are none, start from the beginning.
	lastCacheEndTs := determineLastCacheEnd(etype)
	// If the lastCacheEndTs is younger than globals.BackendConfig.GetCacheGenerationIntervalHours()-1) hours, we do nothing. The cache generator cycle will attempt to create a cache every hour, so this is where we gate how often we create caches.
	cachegenThreshold := api.Timestamp(
		time.Now().Add(-time.Duration(globals.BackendConfig.GetCacheGenerationInterval())).Unix())
	if cachegenThreshold <= lastCacheEndTs {
		// No cache needs to be generated, because we generated it too short a time ago.
		currentCacheEnd = int64(lastCacheEndTs)
		return currentCacheEnd
	}
	// We need to generate some caches.
	cachesTable := generateRequestedCachesTableV2(etype, lastCacheEndTs)
	logging.Logf(2, "New caches table: %v", cachesTable)
	currentCacheEnd = int64(cachesTable[len(cachesTable)-1].EndsAt)
	// ^ We have caches to generate. The end of our last cache is going to be our last cache generation timestamp.
	finalCachesTable := []api.ResultCache{}
	for _, val := range cachesTable {
		if len(val.ResponseUrl) > 0 {
			// If this is a cache we are keeping intact, just add it to the final caches table.
			finalCachesTable = append(finalCachesTable, val)
			continue
			/*
				In the V2 cache generation, there is a consolidation pass that merges past caches as they become older and older. However, this pass tries to retain as many caches intact as possible. Therefore there is a chance that some caches survive it intact, especially the ones that hasn't grown to the next cache size yet. (Ex: you have 9 10 minute caches, when the 10th is added, that becomes one 100 minute cache). In that case, we do not regenerate those intact caches.

				If this is the case, this has been grabbed from a prior cache list as a consolidation intact-survivor. We don't generate those caches.
			*/
		}
		_, resultCache, err := CreateNewCacheV2(etype, val.StartsFrom, val.EndsAt)
		finalCachesTable = append(finalCachesTable, resultCache)
		if err != nil {
			logging.Log(2, err)
		}
	}
	// Read or create the endpoint index.
	epd, err := generateEndpointDir(etype)
	if err != nil {
		logging.Logf(1, "GenerateCachedEndpoint encountered an error while trying to generate the path for the endpoint. Err: %v", err)
		return int64(lastCacheEndTs)
	}
	toolbox.CreatePath(epd)
	var endpointIndex api.ApiResponse
	// Look for the index.json in it. If it doesn't exist, create.
	// Heads up: we're reading and parsing our own caches.
	endpointIndexPath := filepath.Join(epd, "index.json")
	endpointIndexAsJson, err3 := ioutil.ReadFile(endpointIndexPath)
	if err3 != nil && strings.Contains(err3.Error(), "no such file or directory") {
		// The index.json of this cache doesn't exist. Create one.
		endpointIndex.Prefill()
		endpointIndex.Entity = etype
		endpointIndex.Endpoint = etype
	} else if err3 != nil {
		// The index is corrupted. The user knowingly modified it or filesystem did, or some other process did.
		logging.Logf(1, "Cache creation process encountered an error while trying to read the endpoint index. Error: %s", err3)
		return int64(lastCacheEndTs)
	} else {
		// err3 == nil
		json.Unmarshal(endpointIndexAsJson, &endpointIndex)
	}
	// Write the caches table in.
	endpointIndex.Results = finalCachesTable
	// Delete too old caches
	deleteTooOldCaches(etype, &endpointIndex)
	// ^ This modifies the endpoint index. that's why it's before signing.
	// Sign the endpoint index after the changes.
	signingErr := endpointIndex.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if signingErr != nil {
		logging.Logf(1, "This entity index failed to be page-signed. Error: %#v Page: %#v\n", signingErr, endpointIndex)
		return int64(lastCacheEndTs)
	}
	// Save the endpoint index to the disk.
	err2 := saveEndpointIndex(endpointIndex, etype)
	if err2 != nil {
		logging.Logf(1, "GenerateCachedEndpointV2 could not save the endpoint index to update the []ResultsCache in there. Err: %v", err2)
		return int64(lastCacheEndTs)
	}
	rpath, err := generateEndpointRelativePath(etype)
	if err != nil {
		logging.Logf(1, "generateRequestedCachesTableV2 could not generate endpoint relative path for invalidation. Err: %v", err)
	}
	relativeIndexPath := filepath.Join(rpath, "index.json")
	extverify.Verifier.Invalidate(relativeIndexPath)
	logging.Logf(1, "GenerateCachedEndpointV2 is done. Endpoint: %v", etype)
	return currentCacheEnd
}

/*=====  End of V2 CACHING FOR EXPONENTIAL GENERATION & CONSOLIDATION  ======*/

func GenerateCachedEndpoint(etype string) int64 {
	if globals.BackendConfig.GetNetworkHeadDays() > 28 {
		// If network head is larger than 28 days, we enable the cache consolidation. This reduces the number of individual cache files by consolidating them so that syncs happen faster, but it also costs some CPU, because consolidating caches means caches sometimes need to be regenerated to consolidate.
		return GenerateCachedEndpointV2(etype)
	}
	return GenerateCachedEndpointV1(etype)
}

// GenerateCaches generates all caches for all entities and saves them to disk.
func GenerateCaches() {
	logging.Logf(2, "Cache generation has started.")
	oldestCacheEnd := time.Now().Unix()
	// ^ We will set our oldest cache end as cache generation timestamp.
	feapiconsumer.BackendAmbientStatus.CachingStatus = "Generating caches..."
	feapiconsumer.SendBackendAmbientStatus()
	start := time.Now()
	entityTypes := []string{"boards", "threads", "posts", "votes", "keys", "truststates", "addresses"}
	for _, val := range entityTypes {
		cachedEndpointEndTs := GenerateCachedEndpoint(val)
		if cachedEndpointEndTs < oldestCacheEnd {
			oldestCacheEnd = cachedEndpointEndTs
		}
	}
	// /*================================
	// =            debugger            =
	// ================================*/
	// // We want to see the changes first, before committing the cache changes.
	// // taken from generatecachedendpoint()

	// lastCacheEndTs := determineLastCacheEnd("posts")
	// cacheIndex, err := readEndpointIndex("posts")
	// if err != nil {
	// 	logging.Logf(0, "Read cache index errored out. Err: %v", err)
	// 	return
	// }
	// currCachesTable := cacheIndex.Results
	// newCachesViaV1 := generateRequestedCachesTable(lastCacheEndTs)
	// oldtt := append(currCachesTable, newCachesViaV1...)
	// // generateRequestedCachesTableV2(lastCacheEndTs)
	// newCachesViaV2 := generateRequestedCachesTableV2(lastCacheEndTs)
	// newCachesViaV2 = append(currCachesTable, newCachesViaV2...)
	// newtt := MakeConsolidatedTimeTable(&newCachesViaV2, cacheTimeBlocks)
	// logging.Logf(0, "Old caches table")
	// dbg_printCache(oldtt)
	// logging.Logf(0, "New caches table")
	// dbg_printCache(newtt)
	// /*=====  End of debugger  ======*/

	// We're setting this for the purposes of denying POST requests with a timestamp that is partially or wholly available within our cache bracket. (That is, it's not used to determine where to start generating caches from, we read the actual saved cache for that.)
	globals.BackendConfig.SetLastCacheGenerationTimestamp(oldestCacheEnd)
	elapsed := time.Since(start)
	logging.Logf(2, "Cache generation is complete. It took: %s", elapsed)
	feapiconsumer.BackendAmbientStatus.CachingStatus = "Idle"
	feapiconsumer.BackendAmbientStatus.LastCacheGenerationTimestamp = globals.BackendConfig.GetLastCacheGenerationTimestamp()
	feapiconsumer.BackendAmbientStatus.LastCacheGenerationDurationSeconds = int32(elapsed.Seconds())
	feapiconsumer.SendBackendAmbientStatus()
}

// MaintainCaches maintains the Reusable POST response repository by deleting too old post responses that got superseded by caches, and triggers the cache generation if the timing is ready. If not, GenerateCaches will stop itself, so there is no harm in calling this more frequently than cache duration.
func MaintainCaches() {
	globals.BackendTransientConfig.POSTResponseRepo.Maintain()
	GenerateCaches()
}

// func dbg_printCache(ctt []api.ResultCache) {
// 	for _, v := range ctt {
// 		logging.Logf(0, "Start: %v, End: %v, Δ: %v, RUrl: %v \n", time.Unix(int64(v.StartsFrom), 0), time.Unix(int64(v.EndsAt), 0), time.Duration(int64(v.EndsAt-v.StartsFrom))*time.Second, v.ResponseUrl)
// 	}
// }
