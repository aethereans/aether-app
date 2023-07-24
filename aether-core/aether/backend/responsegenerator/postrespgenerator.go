// Backend > ResponseGenerator > PostRespGenerator
// This file provides a set of functions that relate to the generation of non-pregenerated post responses based on the network requests, and their reuse.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/backend/metrics"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/randomhashgen"

	// "encoding/json"
	"errors"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	"strconv"
	// "strings"
	"time"
)

// generatePostFaceResponse is where we generate the actual response to be sent over POST (i.e. not the data container that it points to.)
func generatePostFaceResponse(
	resultPages *[]api.ApiResponse,
	entityCounts *[]api.EntityCount,
	filters *[]api.Filter,
	dirname string,
	reusedPostResponses *[]configstore.POSTResponseEntry,
) *api.ApiResponse {
	resp := api.ApiResponse{}
	resp.Prefill()
	// resp := GeneratePrefilledApiResponse()
	resp.Endpoint = "post_response"
	// This is a slice of filters, but in effect, we are only applying one filter, so it makes sense to return one filter as well. This being a slice is a future expansion point.
	resp.Filters = *filters
	resp.Caching.EntityCounts = *entityCounts
	// Add the chain first as the prior.
	/*
	   Heads up, this means even your response is only one page and you're providing results directly, it can still have []resultcaches because of reused post responses. the remote should look at and process both.
	*/
	chain := generateResultCachesFromPostRespChain(*reusedPostResponses)
	resp.Results = append(resp.Results, chain...)
	// From this point on, our chain is added. These refer to the container being generated right now, or the page generated right now.
	var b int64
	var e int64
	if (*filters)[0].Type == "timestamp" {
		b, _ = strconv.ParseInt((*filters)[0].Values[0], 10, 64)
		e, _ = strconv.ParseInt((*filters)[0].Values[1], 10, 64)
	}
	start := api.Timestamp(b)
	end := api.Timestamp(e)
	if start > end {
		start, end = end, start
	}
	if len(*resultPages) > 1 { // multiple page response: create a container, save it to disk, and return a resultcache[] in the face response.
		// For the time being, we're only generating ONE result cache for POST responses.
		var c api.ResultCache
		// What's below looks like a regeneration of filters for a second time. It is not. We need this to add these values to the resultcache.
		c = constructResultCache(start, end, "")
		foldername := fmt.Sprint("post_", dirname)
		c.ResponseUrl = foldername
		resp.Results = append(resp.Results, c)
		// fmt.Println("results are more than one page")
		// fmt.Println("This is the foldername for the response that was generated based on this request")
		// fmt.Println(fmt.Sprint("post_", dirname))
	} else { // Single page response: send results directly in the face response. (mind that there can still be resultcaches from the reused post responses.)
		// fmt.Println("results are single page")
		entityType := findEntityTypeInApiResponse((*resultPages)[0], "")
		resp.Pagination.Pages = 0 // These start to count from 0
		resp.Pagination.CurrentPage = 0
		resp.Entity = entityType
		resp.ResponseBody = (*resultPages)[0].ResponseBody
	}
	resultsTimeRange := calculateResultTimeRange(resp.Results)
	resp.StartsFrom = resultsTimeRange.Start
	// resp.EndsAt = resultsTimeRange.End
	resp.EndsAt = end // we do not want to return the result time range end - if you had data from t1 to t3 but you were making a call at t5 and no data for t3-t5 range, the timestamp should be t5, not t3.

	// fmt.Println("Outbound POST ApiResponse Results:")
	// fmt.Printf("%#v\n", resp.Results)
	// fmt.Println("Outbound POST ApiResponse ResponseBody lengths:")
	// totalinbound := len(resp.ResponseBody.Boards) + len(resp.ResponseBody.Threads) + len(resp.ResponseBody.Posts) + len(resp.ResponseBody.Votes) + len(resp.ResponseBody.Keys) + len(resp.ResponseBody.Truststates) + len(resp.ResponseBody.Addresses)
	// fmt.Printf("B: %#v\n", len(resp.ResponseBody.Boards))
	// fmt.Printf("T: %#v\n", len(resp.ResponseBody.Threads))
	// fmt.Printf("P: %#v\n", len(resp.ResponseBody.Posts))
	// fmt.Printf("V: %#v\n", len(resp.ResponseBody.Votes))
	// fmt.Printf("K: %#v\n", len(resp.ResponseBody.Keys))
	// fmt.Printf("TS: %#v\n", len(resp.ResponseBody.Truststates))
	// fmt.Printf("A: %#v\n", len(resp.ResponseBody.Addresses))
	// if len(resp.Results) > 0 && totalinbound > 0 {
	// 	fmt.Println("BINGO")
	// }
	// if len(resp.Results) > 1 {
	// 	fmt.Println("TWO RESULTS")
	// }
	logging.LogObj(2, "ResultCaches returning:", resp.Results)
	return &resp
}

// bakeFinalPOSTApiResponse looks at the resultpages. If there is one, it is directly provided as is. If there is more, the results are committed into the file system, and a cachelink page is provided instead.
func bakeFinalPOSTApiResponse(
	resultPages *[]api.ApiResponse,
	indexPages *[]api.ApiResponse,
	manifestPages *[]api.ApiResponse,
	entityCounts *[]api.EntityCount,
	filters *[]api.Filter,
	dirname string,
	mergedEntityCounts *[]api.EntityCount,
	reusedPostResponses *[]configstore.POSTResponseEntry,
	dbReadStartLoc api.Timestamp, // This is needed in the case of container generation.
) (*api.ApiResponse, error) {
	if len(*resultPages) > 1 {
		logging.Logf(2, "This result is still more than one page after the chain addition. We are generating the container %s", dirname)
		generateContainer(resultPages, indexPages, manifestPages, entityCounts, filters, dirname, true, "", dbReadStartLoc) // this will save to disk, doesn't return anything.
		// Generate container needs to generate its own entity
		resp := generatePostFaceResponse(resultPages, mergedEntityCounts, filters, dirname, reusedPostResponses)
		return resp, nil
	} else if len(*resultPages) == 1 {
		logging.Logf(2, "This result is one page after the chain addition. We are not generating a container")
		logging.Logf(2, "This result that is being sent as single page has these items: \nB: %v, T: %v, P: %v, V: %v, K: %v, TS: %v, A: %v",
			len((*resultPages)[0].ResponseBody.Boards),
			len((*resultPages)[0].ResponseBody.Threads),
			len((*resultPages)[0].ResponseBody.Posts),
			len((*resultPages)[0].ResponseBody.Votes),
			len((*resultPages)[0].ResponseBody.Keys),
			len((*resultPages)[0].ResponseBody.Truststates),
			len((*resultPages)[0].ResponseBody.Addresses))
		// Has no container, will be directly served across.
		resp := generatePostFaceResponse(resultPages, mergedEntityCounts, filters, dirname, reusedPostResponses)
		logging.Logf(2, "GeneratePostFaceResponse for the single page result returned these: \nB: %v, T: %v, P: %v, V: %v, K: %v, TS: %v, A: %v",
			len(resp.ResponseBody.Boards),
			len(resp.ResponseBody.Threads),
			len(resp.ResponseBody.Posts),
			len(resp.ResponseBody.Votes),
			len(resp.ResponseBody.Keys),
			len(resp.ResponseBody.Truststates),
			len(resp.ResponseBody.Addresses))
		return resp, nil
	} else {
		logging.Logf(1, "This post request produced both no results and no resulting apiResponses. []ApiResponse: %#v", *resultPages)
		// logging.LogCrash(fmt.Sprintf("This post request produced both no results and no resulting apiResponses. []ApiResponse: %#v", *resultPages))
		return &api.ApiResponse{}, nil
	}
}

// GeneratePOSTResponse creates a response that is directly returned to a custom request by the remote.
func GeneratePOSTResponse(respType string, req api.ApiResponse) ([]byte, error) {
	metrics.SendDbState()
	var resp api.ApiResponse
	resp.Prefill()
	// Look at filterset to figure out what is being requested
	logging.Logf(3, "Filters received raw: %#v", req.Filters)
	filterset := processFilters(&req)
	logging.Logf(3, "Filters processed: %#v", filterset)
	filter := reconstructFilters(filterset)
	logging.Logf(3, "Filters reconstructed: %#v", filter)
	filters := []api.Filter{filter}
	// Create a random SHA256 hash as folder name to use in the case the response has more than one page.
	dirname, err := randomhashgen.GenerateInsecureRandomHash()
	if err != nil {
		logging.Log(1, err)
	}
	logging.Logf(1, "We got a %v POST request with the filters: %#v", respType, filterset)
	switch respType {
	case "node":
		// r := GeneratePrefilledApiResponse()
		// resp = *r
		resp.Endpoint = "node"
		resp.Entity = "node"
	case "boards", "threads", "posts", "votes", "keys", "truststates":
		// Check our post response repo to check if there are any suitable post responses that we can reuse.
		start := configstore.Timestamp(filterset.TimeStart)
		end := configstore.Timestamp(filterset.TimeEnd)
		chain, _, chainEnd, chainCount := globals.BackendTransientConfig.POSTResponseRepo.GetPostResponseChain(start, end, respType)
		dbReadStartLoc := api.Timestamp(0)
		if len(*chain) == 0 {
			dbReadStartLoc = filterset.TimeStart
			logging.Logf(2, "Chain count is zero, therefore dbReadStartLoc is filterset.Timestart, which is %#v", dbReadStartLoc)
		} else {
			dbReadStartLoc = api.Timestamp(chainEnd)
			logging.Logf(2, "Chain count is NOT zero, therefore dbReadStartLoc is chainEnd, which is %#v", dbReadStartLoc)
		}
		// test end
		logging.Logf(2, "Chain: %#v, Start: %v, End: %v Chain Count: %#v Time: %s", chain, start, chainEnd, chainCount, time.Now())
		logging.Logf(2, "These are the values being fed to the persistence.Read. RespType: %s, filterset.Fingerprints: %v, filterset.Embeds: %v, dbReadStartLoc: %v, filterset.TimeEnd: %v", respType, filterset.Fingerprints, filterset.Embeds, dbReadStartLoc, filterset.TimeEnd)
		localData, dbError := persistence.Read(respType, filterset.Fingerprints, filterset.Embeds, dbReadStartLoc, filterset.TimeEnd, false, nil)

		if dbError != nil {
			return []byte{}, errors.New(fmt.Sprintf("The query coming from the remote caused an error in the local database while trying to respond to this request. Error: %#v\n, Request: %#v\n", dbError, req))
		}
		// Generate main data & count the entities resulting. This will go to all three of the response entity pages themselves, the index and the manifest pages.
		pages := splitEntitiesToPages(&localData)
		pagesAsApiResponses := convertResponsesToApiResponses(pages)
		// entityCounts is the count of the cache we're just generating. entities, indexes and manifests should still use it.
		entityCounts := countEntities(&localData)
		// mergedEntityCounts should be used ONLY by the post face response.
		mergedEntityCounts := mergeCounts(entityCounts, chainCount)
		// FUTURE: If there's an entity count bug that surfaces, this is the first place to look. As well as its Address counterpart below.

		/*
		   Future optimisation: we might want to avoid generating indexes and manifests in the case we know the response is going to be only 1 page, thus directly served. That said, if the response is 1 page, then the effort to generate the indexes and the manifest is negligible, so it doesn't matter that much. It's a tradeoff for code clarity, helps us punt the decision whether to paginate until the last moment possible.
		*/
		// Generate indexes
		indexes := createUnbakedIndexes(pages)
		indexPages := splitEntitiesToPages(indexes)
		indexApiResponse := convertResponsesToApiResponses(indexPages)
		for key, _ := range *indexApiResponse {
			(*indexApiResponse)[key].Endpoint = fmt.Sprintf("%s_index_post", (*indexApiResponse)[key].Entity)
		}
		// Generate manifest
		manifest := createUnbakedManifests(pages)
		manifestPages := splitManifestToPages(manifest)
		manifestApiResponse := convertResponsesToApiResponses(manifestPages)
		for key, _ := range *manifestApiResponse {
			(*manifestApiResponse)[key].Endpoint = "manifest_post"
		}
		// bakeFinalPOSTApiResponse wraps the data up and assigns proper metadata. It does not pull any further data in.
		finalResponse, err := bakeFinalPOSTApiResponse(pagesAsApiResponses, indexApiResponse, manifestApiResponse, entityCounts, &filters, dirname, &mergedEntityCounts, chain, dbReadStartLoc)
		// fmt.Printf("%#v", finalResponse)
		if err != nil {
			return []byte{}, errors.New(fmt.Sprintf("An error was encountered while trying to finalise the API response. Error: %#v\n, Request: %#v\n", err, req))
		}
		resp = *finalResponse
	case "addresses": // Addresses can't do address search by loc/subloc/port. Only time search is available, since adresses don't have fingerprints defined.
		/*
		   An addresses POST response returns results within the time boundary that has been seen online first-person by the remote. It does not communicate addresses that the remote has not connected to.
		*/
		// Check our post response repo to check if there are any suitable post responses that we can reuse.
		start := configstore.Timestamp(filterset.TimeStart)
		end := configstore.Timestamp(filterset.TimeEnd)
		chain, _, chainEnd, chainCount := globals.BackendTransientConfig.POSTResponseRepo.GetPostResponseChain(start, end, respType)
		// fmt.Printf("requested a chain for %#v, this is what we got: %#v\n", respType, chain)
		dbReadStartLoc := api.Timestamp(0)
		if len(*chain) == 0 {
			dbReadStartLoc = filterset.TimeStart
		} else {
			dbReadStartLoc = api.Timestamp(chainEnd)
		}
		// Containers will always return max 100 addresses. (80 live nodes, 10 bootstrap nodes and 10 static nodes max)
		addresses, dbError := persistence.ReadAddresses("", "", 0, api.Timestamp(dbReadStartLoc), filterset.TimeEnd, 100, 0, 0, "container_generate")
		addresses = *sanitiseOutboundAddresses(&addresses)
		var localData api.Response
		localData.Addresses = addresses
		if dbError != nil {
			return []byte{}, errors.New(fmt.Sprintf("The query coming from the remote caused an error in the local database while trying to respond to this request. Error: %#v\n, Request: %#v\n", dbError, req))
		}
		// Addresses do not have indexes or manifests, so this is much simpler.
		pages := splitEntitiesToPages(&localData)
		pagesAsApiResponses := convertResponsesToApiResponses(pages)
		entityCounts := countEntities(&localData)
		// mergedEntityCounts should be used ONLY by the post face response.
		mergedEntityCounts := mergeCounts(entityCounts, chainCount)
		finalResponse, err := bakeFinalPOSTApiResponse(pagesAsApiResponses, nil, nil, entityCounts, &filters, dirname, &mergedEntityCounts, chain, dbReadStartLoc)
		if err != nil {
			return []byte{}, errors.New(fmt.Sprintf("An error was encountered while trying to finalise the API response. Error: %#v\n, Request: %#v\n", err, req))
		}
		resp = *finalResponse
		resp.Endpoint = "entity"
	}
	// Build the response itself
	resp.Entity = respType
	resp.Timestamp = api.Timestamp(time.Now().Unix())
	signingErr := resp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if signingErr != nil {
		return []byte{}, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to be page-signed. Error: %#v Response Body: %#v\n", signingErr, resp))
	}
	// Construct the query, and run an index to determine how many entries we have for the filter.
	// jsonResp, err := ConvertApiResponseToJson(&resp)
	jsonResp, err := resp.ToJSON()
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("The response that was prepared to respond to this query failed to convert to JSON. Error: %#v\n, Request Body: %#v\n", err, req))
	}
	return jsonResp, nil
}

// insertIntoPOSTResponseTracker checks some conditions to determine whether this post response is eligible for reuse. If that is the case, it will insert it to the queue to render it eligible.
func insertIntoPOSTResponseReuseTracker(resultPage *api.ApiResponse, foldername string, dbReadStartLoc api.Timestamp) {
	pg := *resultPage // we can take a look at only one page because we know all are the same.
	// Rules: A reusable response can only have a time range filter (timestamp) and it can only have one single entity type (represented by entity counts). If those are true, then we can reuse this.
	if len(pg.Caching.EntityCounts) == 1 &&
		len(pg.Filters) == 1 &&
		pg.Filters[0].Type == "timestamp" &&
		len(pg.Filters[0].Values) == 2 {
		s, _ := strconv.Atoi(pg.Filters[0].Values[0])
		e, _ := strconv.Atoi(pg.Filters[0].Values[1])
		start := configstore.Timestamp(s)
		end := configstore.Timestamp(e)
		if start > end {
			start, end = end, start
		}
		creation := configstore.Timestamp(pg.Timestamp)
		counts := convertToConfigStoreEntityCount(pg.Caching.EntityCounts)
		globals.BackendTransientConfig.POSTResponseRepo.Add(foldername, configstore.Timestamp(dbReadStartLoc), end, creation, &counts) // not start, but dbReadStartLoc. Because start is the start of the entire response, not the container we just generated. This is related to the container that we generated, so we use the dbReadStartLoc.
	}
}

func generateResultCachesFromPostRespChain(chain []configstore.POSTResponseEntry) []api.ResultCache {
	var rcachs []api.ResultCache

	for i, _ := range chain {
		rcachs = append(rcachs, constructResultCache(
			api.Timestamp(chain[i].StartsFrom),
			api.Timestamp(chain[i].EndsAt), chain[i].ResponseUrl))
	}
	return rcachs
}
