package beapiconsumer

import (
	pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	// "fmt"
	"time"
)

/*================================================
=            Internal cache structure            =
================================================*/

// This is not publicly accessible directly.

type backendAPICache struct {
	Boards                []*pbstructs.Board
	Threads               []*pbstructs.Thread
	Posts                 []*pbstructs.Post
	Votes                 []*pbstructs.Vote
	Keys                  []*pbstructs.Key
	Truststates           []*pbstructs.Truststate
	Prefilled             bool
	PrefillStartTimestamp int64
	PrefillEndTimestamp   int64
}

var cache backendAPICache

/*==============================================
=            Public API for caching            =
==============================================*/

// This is used by the refresher to prime the cache before a refresh starts.

func PrefillCache(lastRef int64) int64 {
	cache.PrefillStartTimestamp = lastRef
	cache.PrefillEndTimestamp = time.Now().Unix()
	cache.Boards = GetBoards(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{},
		true, true)
	cache.Threads = GetThreads(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{}, "",
		true, true)
	cache.Posts = GetPosts(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{}, "", "",
		true, true)
	cache.Votes = GetVotes(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{}, "", "", "", -1, -1, false,
		true, true)
	cache.Keys = GetKeys(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{},
		true, true)
	cache.Truststates = GetTruststates(
		cache.PrefillStartTimestamp,
		cache.PrefillEndTimestamp, []string{}, -1, -1, "", "",
		true, true)
	cache.Prefilled = true
	logging.Logf(1, "Cache is prefilled and ready. Boards: %v, Threads, %v, Posts: %v, Votes: %v, Keys: %v, Truststates: %v", len(cache.Boards), len(cache.Threads), len(cache.Posts), len(cache.Votes), len(cache.Keys), len(cache.Truststates))
	return cache.PrefillEndTimestamp // the refresh cycle should use this as now, so that the time range matches. Otherwise, there might be gaps.
}

func ReleaseCache() {
	cache = backendAPICache{}
}

/*=============================================
=            Cache Query structure            =
=============================================*/

// This is used by the beapiconsumer methods to ask cache data.

type cacheQuery struct {
	// ^ Why do we have this when we explicitly prefill the cache at every refresh? Well, we can actually stop refreshing some stuff and when the time comes, the range they might need to refresh might fall out of the range of the cache. And they might actually end up requesting this manual refresh while the automatic refresh is in place, meaning there would be a cache active.
	Fingerprints []string
	// Boards
	// Threads
	Thread_Board string
	// Posts
	Post_Board  string
	Post_Thread string
	Post_Parent string
	// Votes
	Vote_Board         string
	Vote_Thread        string
	Vote_Target        string
	Vote_TypeClass     int
	Vote_Type          int
	Vote_NoDescendants bool
	// Keys
	// Truststates
	Truststate_Target    string
	Truststate_Domain    string
	Truststate_TypeClass int
	Truststate_Type      int
}

/*==================================================
=            Methods to pull from cache            =
==================================================*/

func queryBoardsCache(q cacheQuery) []*pbstructs.Board {
	var result []*pbstructs.Board
	for k, _ := range cache.Boards {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Boards[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		result = append(result, cache.Boards[k])
	}
	return result
}

func queryThreadsCache(q cacheQuery) []*pbstructs.Thread {
	var result []*pbstructs.Thread
	for k, _ := range cache.Threads {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Threads[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		if !parentFpFilter(cache.Threads[k].GetBoard(), q.Thread_Board) {
			continue
		}
		result = append(result, cache.Threads[k])
	}
	return result
}

func queryPostsCache(q cacheQuery) []*pbstructs.Post {
	var result []*pbstructs.Post
	for k, _ := range cache.Posts {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Posts[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		if !boardFpFilter(cache.Posts[k].GetBoard(), q.Post_Board) {
			continue
		}
		if !threadFpFilter(cache.Posts[k].GetThread(), q.Post_Thread) {
			continue
		}
		if !parentFpFilter(cache.Posts[k].GetParent(), q.Post_Parent) {
			continue
		}
		result = append(result, cache.Posts[k])
	}
	return result
}

func queryVotesCache(q cacheQuery, parentType string) []*pbstructs.Vote {
	var result []*pbstructs.Vote
	for k, _ := range cache.Votes {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Votes[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		if !boardFpFilter(cache.Votes[k].GetBoard(), q.Vote_Board) {
			continue
		}
		if !threadFpFilter(cache.Votes[k].GetThread(), q.Vote_Thread) {
			continue
		}
		if !targetFpFilter(cache.Votes[k].GetTarget(), q.Vote_Target) {
			// if parentType == "thread" {
			// 	logging.Logf(1, "failed target fp filter")
			// }
			continue
		}
		if !typeFilter(cache.Votes[k].GetType(), q.Vote_Type) {
			continue
		}
		if !typeClassFilter(cache.Votes[k].GetTypeClass(), q.Vote_TypeClass) {
			continue
		}
		if q.Vote_NoDescendants {
			if parentType == "board" {
				// Get all things whose parent is a board, but no children of those. (i.e. get all thread entities but not their boards)
				if !noDescendantsFilter(cache.Votes[k].GetTarget(), cache.Votes[k].GetThread(), q.Vote_NoDescendants) {
					continue
				}
				// nodescendants is not valid for "thread" entitytype, since that would be not very useful (it would get all first-level posts, which we have no specific use for)
			}
		}
		result = append(result, cache.Votes[k])
	}
	return result
}

func queryKeysCache(q cacheQuery) []*pbstructs.Key {
	var result []*pbstructs.Key
	for k, _ := range cache.Keys {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Keys[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		result = append(result, cache.Keys[k])
	}
	return result
}

func queryTruststatesCache(q cacheQuery) []*pbstructs.Truststate {
	var result []*pbstructs.Truststate
	for k, _ := range cache.Truststates {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Truststates[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		if !targetFpFilter(cache.Truststates[k].GetTarget(), q.Truststate_Target) {
			continue
		}
		if !domainFpFilter(cache.Truststates[k].GetDomain(), q.Truststate_Domain) {
			continue
		}
		if !typeFilter(cache.Truststates[k].GetType(), q.Truststate_Type) {
			continue
		}
		if !typeClassFilter(cache.Truststates[k].GetTypeClass(), q.Truststate_TypeClass) {
			continue
		}
		result = append(result, cache.Truststates[k])
	}
	return result
}

/*==============================================
=            Filter validity checks            =
==============================================*/

func fingerprintsFilter(fp string, fps []string) bool {
	// Check if filter enabled
	if len(fps) == 0 {
		return true
	}
	if i := toolbox.IndexOf(
		fp, fps); i != -1 {
		return true
	}
	return false
}

func fpBaseFilter(item, requested string, checkEnabled bool) bool {
	if checkEnabled && len(requested) == 0 {
		return true
	}
	return item == requested
}

func parentFpFilter(item, requested string) bool {
	return fpBaseFilter(item, requested, true)
}

func boardFpFilter(item, requested string) bool {
	return fpBaseFilter(item, requested, true)
}

func threadFpFilter(item, requested string) bool {
	return fpBaseFilter(item, requested, true)
}

func targetFpFilter(item, requested string) bool {
	return fpBaseFilter(item, requested, true)
}

func domainFpFilter(item, requested string) bool {
	return fpBaseFilter(item, requested, false)
}

func typeBaseFilter(item int32, requested int) bool {
	if requested == -1 {
		return true
	}
	return item == int32(requested)
}

func typeFilter(item int32, requested int) bool {
	return typeBaseFilter(item, requested)
}

func typeClassFilter(item int32, requested int) bool {
	return typeBaseFilter(item, requested)
}

func noDescendantsFilter(itemTarget, itemParent string, noDescendants bool) bool {
	if noDescendants {
		return itemTarget == itemParent
	}
	return true
}

func servableFromCache(start, end int64, bypassCache bool) bool {
	// logging.Logf(1, "Requested Start: %v, CacheStart: %v, RequestedEnd: %v,  CacheEnd: %v", start, cache.PrefillStartTimestamp, end, cache.PrefillEndTimestamp)
	if bypassCache {
		// logging.Logf(1, "failed due to bypasscache enabled")
		return false
	}
	if !cache.Prefilled {
		// logging.Logf(1, "prefilled failed")
		return false
	}
	if start < cache.PrefillStartTimestamp {
		// logging.Logf(1, "prefill start timestamp failed")
		return false
	}
	if end > cache.PrefillEndTimestamp {
		// logging.Logf(1, "prefill end timestamp failed")
		return false
	}
	return true
}
