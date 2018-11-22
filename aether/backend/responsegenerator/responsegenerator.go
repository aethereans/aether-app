// Backend > ResponseGenerator
// This file provides a set of functions that take a database response, and convert it into a set of paginated (or nonpaginated) results.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	// "aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "aether-core/aether/services/randomhashgen"
	// "aether-core/aether/services/toolbox"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	// "os"
	"strconv"
	// "strings"
	"time"
)

type FilterSet struct {
	Fingerprints []api.Fingerprint
	TimeStart    api.Timestamp
	TimeEnd      api.Timestamp
	Embeds       []string
}

func processFilters(req *api.ApiResponse) FilterSet {
	var fs FilterSet
	for _, filter := range req.Filters {
		// Fingerprint
		if filter.Type == "fingerprint" {
			for _, fp := range filter.Values {
				fs.Fingerprints = append(fs.Fingerprints, api.Fingerprint(fp))
			}
		}
		// Embeds
		if filter.Type == "embed" {
			for _, embed := range filter.Values {
				fs.Embeds = append(fs.Embeds, embed)
			}
		}
		// If a time filter is given, timeStart is either the timestamp provided by the remote if it's larger than the end date of the last cache, or the end timestamp of the last cache.
		// In essence, we do not provide anything that is already cached from the live server.
		if filter.Type == "timestamp" {
			// now := int64(time.Now().Unix())
			start, _ := strconv.ParseInt(filter.Values[0], 10, 64)
			end, _ := strconv.ParseInt(filter.Values[1], 10, 64)

			// If there is a value given (not 0), that is, the timerange filter is active.
			// The sanitisation of these ranges are done in the DB level, so this is just intake.
			if start > 0 || end > 0 {
				fs.TimeStart = api.Timestamp(start)
				fs.TimeEnd = api.Timestamp(end)
			}
			if fs.TimeStart == 0 {
				networkHeadThreshold := api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
				fs.TimeStart = networkHeadThreshold
			}
			if fs.TimeEnd == 0 {
				fs.TimeEnd = api.Timestamp(time.Now().Unix())
			}

		}
	}
	return fs
}

// convertResponsesToApiResponses - why is this plural and not singular? Because it needs to know the number of pages to insert to every page and it can only do that if it has the whole [].
func convertResponsesToApiResponses(r *[]api.Response) *[]api.ApiResponse {
	var responses []api.ApiResponse
	if r == nil {
		return &responses
	}
	for i, _ := range *r {
		resp := api.ApiResponse{}
		resp.Prefill()
		// resp := GeneratePrefilledApiResponse()
		resp.ResponseBody.Boards = (*r)[i].Boards
		resp.ResponseBody.Threads = (*r)[i].Threads
		resp.ResponseBody.Posts = (*r)[i].Posts
		resp.ResponseBody.Votes = (*r)[i].Votes
		resp.ResponseBody.Addresses = (*r)[i].Addresses
		resp.ResponseBody.Keys = (*r)[i].Keys
		resp.ResponseBody.Truststates = (*r)[i].Truststates
		// Indexes
		resp.ResponseBody.BoardIndexes = (*r)[i].BoardIndexes
		resp.ResponseBody.ThreadIndexes = (*r)[i].ThreadIndexes
		resp.ResponseBody.PostIndexes = (*r)[i].PostIndexes
		resp.ResponseBody.VoteIndexes = (*r)[i].VoteIndexes
		resp.ResponseBody.AddressIndexes = (*r)[i].AddressIndexes
		resp.ResponseBody.KeyIndexes = (*r)[i].KeyIndexes
		resp.ResponseBody.TruststateIndexes = (*r)[i].TruststateIndexes
		// Manifests
		resp.ResponseBody.BoardManifests = (*r)[i].BoardManifests
		resp.ResponseBody.ThreadManifests = (*r)[i].ThreadManifests
		resp.ResponseBody.PostManifests = (*r)[i].PostManifests
		resp.ResponseBody.VoteManifests = (*r)[i].VoteManifests
		resp.ResponseBody.AddressManifests = (*r)[i].AddressManifests
		resp.ResponseBody.KeyManifests = (*r)[i].KeyManifests
		resp.ResponseBody.TruststateManifests = (*r)[i].TruststateManifests

		resp.Pagination.Pages = uint64(len(*r) - 1) // pagination starts from 0
		resp.Pagination.CurrentPage = uint64(i)
		// entityType := findEntityTypeInApiResponse(*resp)
		// resp.Entity = entityType
		responses = append(responses, resp)
	}
	return &responses
}

// findEntityTypeInApiResponse looks at the response and determines the type based on its ingredients. If the response is empty, it returns the response type from the outside context.
func findEntityTypeInApiResponse(resp api.ApiResponse, respType string) string {
	if len(resp.ResponseBody.Boards) > 0 || len(resp.ResponseBody.BoardIndexes) > 0 || len(resp.ResponseBody.BoardManifests) > 0 {
		return "boards"
	}
	if len(resp.ResponseBody.Threads) > 0 || len(resp.ResponseBody.ThreadIndexes) > 0 || len(resp.ResponseBody.ThreadManifests) > 0 {
		return "threads"
	}
	if len(resp.ResponseBody.Posts) > 0 || len(resp.ResponseBody.PostIndexes) > 0 || len(resp.ResponseBody.PostManifests) > 0 {
		return "posts"
	}
	if len(resp.ResponseBody.Votes) > 0 || len(resp.ResponseBody.VoteIndexes) > 0 || len(resp.ResponseBody.VoteManifests) > 0 {
		return "votes"
	}
	if len(resp.ResponseBody.Addresses) > 0 || len(resp.ResponseBody.AddressIndexes) > 0 || len(resp.ResponseBody.AddressManifests) > 0 {
		return "addresses"
	}
	if len(resp.ResponseBody.Keys) > 0 || len(resp.ResponseBody.KeyIndexes) > 0 || len(resp.ResponseBody.KeyManifests) > 0 {
		return "keys"
	}
	if len(resp.ResponseBody.Truststates) > 0 || len(resp.ResponseBody.TruststateIndexes) > 0 || len(resp.ResponseBody.TruststateManifests) > 0 {
		return "truststates"
	}
	return respType
}

func saveFileToDisk(fileContents []byte, path string, filename string) {
	ioutil.WriteFile(fmt.Sprint(path, "/", filename), fileContents, 0755)
}

// reconstructFilters reconstructs the filters to record in the response. This also does validation so that it will match what we have on the response itself.
func reconstructFilters(filterset FilterSet) api.Filter {
	filter := api.Filter{}
	if len(filterset.Fingerprints) > 0 {
		filter.Type = "fingerprint"
		for _, val := range filterset.Fingerprints {
			filter.Values = append(filter.Values, string(val))
		}
	} else if len(filterset.Embeds) > 0 {
		filter.Type = "embed"
		filter.Values = filterset.Embeds
	} else {
		// If no filters are given, this is by default a time range search.
		filter.Type = "timestamp"
		begin, end, err := persistence.SanitiseTimeRange(filterset.TimeStart, filterset.TimeEnd, api.Timestamp(time.Now().Unix()), false)
		if err != nil {
			logging.Log(1, fmt.Sprintf("SanitiseTimeRange errored out. Error: %#v", err))
		}
		// The beginning can be 0 in conditions where the remote wants you to just give from the beginning of time, or however much you allow. The allowance is the network head, so we generate the end of the network head here.
		if begin == 0 {
			begin = api.Timestamp(time.Now().Add(-time.Duration(globals.BackendConfig.GetNetworkHeadDays()*24) * time.Hour).Unix())
		}
		if end == 0 {
			end = api.Timestamp(time.Now().Unix())
		}
		filter.Values = append(filter.Values, strconv.FormatInt(int64(begin), 10))
		filter.Values = append(filter.Values, strconv.FormatInt(int64(end), 10))
	}
	return filter
}

type resultTimeRange struct {
	Start api.Timestamp
	End   api.Timestamp
}

func calculateResultTimeRange(res []api.ResultCache) resultTimeRange {
	start := api.Timestamp(0)
	end := api.Timestamp(0)
	for _, val := range res {
		if val.StartsFrom < start || start == 0 {
			start = val.StartsFrom
		}
		if val.EndsAt > end || end == 0 {
			end = val.EndsAt
		}
	}
	return resultTimeRange{start, end}
}

// generateContainer always creates a container from the given data (can be a post response or a cache response) and it saves it to the disk. It does not care about how many pages the result is.
func generateContainer(
	entityPages *[]api.ApiResponse,
	indexPages *[]api.ApiResponse,
	manifestPages *[]api.ApiResponse,
	entityCounts *[]api.EntityCount,
	filters *[]api.Filter,
	dirname string,
	isPOST bool,
	respType string,
	dbReadStartLoc api.Timestamp,
) {
	foldername := ""
	// Gate the filter in such a way that the beginning of the range will be the beginning of the DB read that this container will hold, NOT the beginning of the scan range. Scan range can include the chain with other reused responses, but dbReadStartLoc is the range of the DB read only.
	flt := *filters
	if flt[0].Type == "timestamp" {
		rangeBeg := flt[0].Values[0]
		rangeEnd := flt[0].Values[1]
		if rangeBeg > rangeEnd {
			rangeBeg, rangeEnd = rangeEnd, rangeBeg
		}
		flt[0].Values[0] = strconv.Itoa(int(dbReadStartLoc))
		flt[0].Values[1] = rangeEnd
	}
	if isPOST {
		foldername = fmt.Sprint("post_", dirname)
	} else {
		foldername = fmt.Sprint("cache_", dirname)
	}
	entityType := findEntityTypeInApiResponse((*entityPages)[0], respType)
	// fmt.Println("entityType")
	// fmt.Println(entityType)
	// Create the index and manifest pages.
	if indexPages != nil {
		bakeIndexes(indexPages, entityCounts, &flt, foldername, isPOST, respType, entityType)
	}
	if manifestPages != nil {
		bakeManifests(manifestPages, entityCounts, &flt, foldername, isPOST, respType, entityType)
	}
	// Bake the main entity pages.
	bakeEntityPages(entityPages, entityCounts, &flt, foldername, isPOST, respType, entityType)
}

func constructResultCache(beg api.Timestamp, end api.Timestamp, url string) api.ResultCache {
	if beg > end {
		beg, end = end, beg
	}
	return api.ResultCache{StartsFrom: beg, EndsAt: end, ResponseUrl: url}
}

// sanitiseOutboundAddresses removes untrusted address data from the addresses destined to go out of this node. The remote node will also remove it, but there is no reason to leak information unnecessarily.
func sanitiseOutboundAddresses(addrsPtr *[]api.Address) *[]api.Address {
	addrs := *addrsPtr
	for key, _ := range addrs {
		addrs[key].LocationType = 0
		addrs[key].Type = 0
		addrs[key].LastSuccessfulPing = 0
		addrs[key].LastSuccessfulSync = 0
		addrs[key].Protocol = api.Protocol{}
		addrs[key].Protocol.Subprotocols = []api.Subprotocol{}
		addrs[key].Client = api.Client{}
	}
	return &addrs
}
