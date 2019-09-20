package beapiconsumer

import (
	"aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"sort"
	// "fmt"
	"time"
)

/*================================================
=            Internal cache structure            =
================================================*/

// This is not publicly accessible directly.

type backendAPICache struct {
	Boards                []*mimapi.Board
	Threads               []*mimapi.Thread
	Posts                 []*mimapi.Post
	Votes                 []*mimapi.Vote
	Keys                  []*mimapi.Key
	Truststates           []*mimapi.Truststate
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
	start := time.Now()
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
	elapsed := time.Since(start)
	logging.Logf(1, "Cache is prefilled in %s and ready. Boards: %v, Threads, %v, Posts: %v, Votes: %v, Keys: %v, Truststates: %v", elapsed, len(cache.Boards), len(cache.Threads), len(cache.Posts), len(cache.Votes), len(cache.Keys), len(cache.Truststates))
	return cache.PrefillEndTimestamp // the refresh cycle should use this as now, so that the time range matches. Otherwise, there might be gaps.
}

func ReleaseCache() {
	cache = backendAPICache{}
}

// DetermineObservableUniverse looks at the oncoming objects on this batch, and determines what objects can possibly be affected. This is useful, because if we don't know this, we need to run the refresh on the whole frontend database, which takes a long time. This, for example, finds out which boards have received new content or signals and adds them to the observable universe, so that they can be refreshed. For boards where nothing has happened, we just don't refresh, and save a lot of work.
func DetermineObservableUniverse() map[string]map[string]bool {
	ou := make(map[string]map[string]bool)
	// Determine observable universe for: boards
	boardFps := make(map[string]bool)
	for k, _ := range cache.Boards {
		boardFps[cache.Boards[k].GetProvable().GetFingerprint()] = true
	}
	for k, _ := range cache.Threads {
		boardFps[cache.Threads[k].GetBoard()] = true
	}
	for k, _ := range cache.Posts {
		boardFps[cache.Posts[k].GetBoard()] = true
	}
	for k, _ := range cache.Votes {
		boardFps[cache.Votes[k].GetBoard()] = true
	}
	// What's not included: keys that create the boards, and truststates that point to those keys. (I.e. if a board owner gets a 'member' badge in orange it won't automatically trigger a wholesale board update.) This is mostly for efficiency reasons, since boards can have an arbitrary number of board owners and elected mods.
	ou["Boards"] = boardFps
	// Determine observable universe for: threads
	// threadFps := make(map[string]bool)
	// for k, _ := range cache.Threads {
	// 	threadFps[cache.Threads[k].GetProvable().GetFingerprint()] = true
	// }
	// for k, _ := range cache.Posts {
	// 	threadFps[cache.Posts[k].GetThread()] = true
	// }
	// for k, _ := range cache.Votes {
	// 	threadFps[cache.Votes[k].GetThread()] = true
	// }
	// ou["Threads"] = threadFps
	// // Same as above re: not included
	// postFps := make(map[string]bool)
	// for k, _ := range cache.Posts {
	// 	postFps[cache.Posts[k].GetProvable().GetFingerprint()] = true
	// }
	// for k, _ := range cache.Posts {
	// 	postFps[cache.Posts[k].GetParent()] = true
	// }
	// for k, _ := range cache.Votes {
	// 	postFps[cache.Votes[k].GetTarget()] = true
	// }
	// ou["Posts"] = postFps
	// Determine observable universe for: votes
	// FUTURE: if we need this, implement it here
	// Determine observable universe for: keys
	keyFps := make(map[string]bool)
	for k, _ := range cache.Keys {
		keyFps[cache.Keys[k].GetProvable().GetFingerprint()] = true
	}
	for k, _ := range cache.Truststates {
		keyFps[cache.Truststates[k].GetTarget()] = true
	}
	ou["Keys"] = keyFps
	// Determine observable universe for: truststates
	// FUTURE: if we need this, implement it here
	return ou
}

/*================================
=            New feed            =
================================*/
// New feed is the feed of content that appears on 'new' on linear time order. New feed is composed of posts and threads, but nothing else.

type NewFeedItem struct {
	EntityType  string
	BoardFp     string
	ThreadFp    string
	ParentFp    string
	TargetFp    string
	Fingerprint string
	UserFp      string
	Creation    int64
}

// Make it satisfy Searchable interface
func (o *NewFeedItem) GetEntityType() string {
	return o.EntityType
}
func (o *NewFeedItem) GetBoardFp() string {
	return o.BoardFp
}
func (o *NewFeedItem) GetThreadFp() string {
	return o.ThreadFp
}
func (o *NewFeedItem) GetParentFp() string {
	return o.ParentFp
}
func (o *NewFeedItem) GetFingerprint() string {
	return o.Fingerprint
}
func (o *NewFeedItem) GetUserFp() string {
	return o.UserFp
}

func GenerateNewFeed(firstEverGeneration bool) []NewFeedItem {
	newFeed := []NewFeedItem{}
	sfwlist := make(map[string]bool)
	for _, v := range globals.FrontendConfig.ContentRelations.SFWList.Boards {
		sfwlist[v] = true
	}
	for _, v := range cache.Threads {
		if !(sfwlist[v.GetBoard()]) {
			continue
		}
		nfi := NewFeedItem{
			EntityType:  "Thread",
			BoardFp:     v.GetBoard(),
			ThreadFp:    "",
			ParentFp:    "",
			Fingerprint: v.GetProvable().GetFingerprint(),
			UserFp:      v.GetOwner(),
			Creation:    v.GetProvable().GetCreation(),
		}
		newFeed = append(newFeed, nfi)
	}
	for _, v := range cache.Posts {
		if !(sfwlist[v.GetBoard()]) {
			continue
		}
		nfi := NewFeedItem{
			EntityType:  "Post",
			BoardFp:     v.GetBoard(),
			ThreadFp:    v.GetThread(),
			ParentFp:    v.GetParent(),
			Fingerprint: v.GetProvable().GetFingerprint(),
			UserFp:      v.GetOwner(),
			Creation:    v.GetProvable().GetCreation(),
		}
		newFeed = append(newFeed, nfi)
	}
	if !firstEverGeneration {
		for _, v := range cache.Votes {
			if !(sfwlist[v.GetBoard()]) {
				continue
			}
			nfi := NewFeedItem{
				EntityType:  "Vote",
				BoardFp:     v.GetBoard(),
				ThreadFp:    v.GetThread(),
				ParentFp:    "",
				TargetFp:    v.GetTarget(),
				Fingerprint: v.GetProvable().GetFingerprint(),
				UserFp:      v.GetOwner(),
				Creation:    v.GetProvable().GetCreation(),
			}
			newFeed = append(newFeed, nfi)
		}
	}
	sort.Slice(newFeed, func(i, j int) bool {
		// If less or more, fairly straightforward. Sort by that.
		if newFeed[i].Creation > newFeed[j].Creation {
			return true
		}
		return false
	})
	return newFeed
}

/*=====  End of New feed  ======*/

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

func queryBoardsCache(q cacheQuery) []*mimapi.Board {
	var result []*mimapi.Board
	for k, _ := range cache.Boards {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Boards[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		result = append(result, cache.Boards[k])
	}
	return result
}

func queryThreadsCache(q cacheQuery) []*mimapi.Thread {
	var result []*mimapi.Thread
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

func queryPostsCache(q cacheQuery) []*mimapi.Post {
	var result []*mimapi.Post
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

func queryVotesCache(q cacheQuery, parentType string) []*mimapi.Vote {
	var result []*mimapi.Vote
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

func queryKeysCache(q cacheQuery) []*mimapi.Key {
	var result []*mimapi.Key
	for k, _ := range cache.Keys {
		// Fingerprints range checker
		if !fingerprintsFilter(cache.Keys[k].GetProvable().GetFingerprint(), q.Fingerprints) {
			continue
		}
		result = append(result, cache.Keys[k])
	}
	return result
}

func queryTruststatesCache(q cacheQuery) []*mimapi.Truststate {
	var result []*mimapi.Truststate
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
	// return true // debug
	// logging.Logf(1, "Requested Start: %v, CacheStart: %v, RequestedEnd: %v,  CacheEnd: %v", start, cache.PrefillStartTimestamp, end, cache.PrefillEndTimestamp)
	if bypassCache {
		logging.Logf(1, "Cannot be served from the cache: bypasscache enabled.")
		return false
	}
	if !cache.Prefilled {
		logging.Logf(1, "Cannot be served from the cache: cache not prefilled.")
		return false
	}
	if start < cache.PrefillStartTimestamp {
		logging.Logf(1, "Cannot be served from the cache: start ts before the prefill start ts. Prefill start: %#v, Asked start: %#v, Difference in seconds: %#v", cache.PrefillStartTimestamp, start, cache.PrefillStartTimestamp-start)
		// if start != 0 {
		// 	logging.LogCrash("yo")
		// }
		return false
	}
	if end > cache.PrefillEndTimestamp {
		logging.Logf(1, "Cannot be served from the cache: end ts after the prefill end ts. Prefill end: %#v, Asked end: %#v, Difference in seconds: %#v", cache.PrefillEndTimestamp, end, cache.PrefillEndTimestamp-end)
		return false
	}
	return true
}
