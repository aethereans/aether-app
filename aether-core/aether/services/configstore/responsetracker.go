// Services > ConfigStore > ResponseTracker
// This file handles the maintenance and control of the POST response records that are kept in the ephemeral working memory. Keeping records of this is important because when responding to POST responses, we can chain these responses, and not have to spend more resources regenerating them every time we need to provide it.

/*
What's the logic there?

The idea is that we keep a list of post responses that we have generated that did not have any filters except for time. That means these responses are actually full dumps of the database at that given time for that type of object.

We have caches, which are generated every x ticks. Let's say this is 60 minutes. Let's think of the worst case. 59 minutes after you generate a cache, people are getting up to speed until the last 59 minutes via the caches, but every response ends up with you providing 59 minutes of data directly from the DB also. This is a little bit of a waste. These 59 minutes of data that you generated for the last 20 people are exactly the same, or nearly so (depending on which exact second).

So, instead of doing this, what we do is we actually go in and generate a response when somebody asks. And when a second guy comes in, the second guy will receive a response that points to the first response, and the delta from the end of the first response to now.

So, wouldn't this end up with some guy receiving 500 responses in a list? No. Because the only way you can get into this chain is that you generate a multi-page response, which triggers a save to disk. Unless your response is saved to disk, you are not getting to this list. So if the response is small enough that it fits one single page, then an entry is not created here. Effectively the page size becomes a gate for the minimum timespan covered of the responses here.

Also, our responses are deleted after some time, and rendered unusable after that time divided by 2. Why? Because if those were the same, we would be able to provide a response instants before it was deleted, which would effectively be unavailable. We still need to provide some time for the remote to be able to download it.

This does not affect more specialised filters that are applied - if a response is generated as a result of a post request asking for a specific entity, or had a specific embed, then it is ineligible for this processing, and those requests are served normally.
*/

package configstore

import (
	"aether-core/aether/services/toolbox"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

// Basic types

type Timestamp int64

func (t Timestamp) String() string {
	return time.Unix(int64(t), 0).String()
}

type EntityCount struct {
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type POSTResponseEntry struct {
	ResponseUrl  string
	StartsFrom   Timestamp
	EndsAt       Timestamp
	Creation     Timestamp
	EntityCounts []EntityCount
}

// In case we need to add some more fields in the future.
type POSTResponseRepo struct {
	lock      sync.Mutex
	Responses []POSTResponseEntry
}

// Internal helper functions

// deleteFromDisk deletes the POST response from the directory. This only triggers when the item is also removed from the Responses repo.
func deleteFromDisk(url string) {
	// postDir := fmt.Sprintf("%s/", bc.GetProtURLVersion(), "/responses/%s", bc.GetCachesDirectory(), url)
	postDir := filepath.Join(bc.GetCachesDirectory(), bc.GetProtURLVersion(), "responses", url)
	toolbox.DeleteFromDisk(postDir)
}

// collectCounts assumes there is only one type of entity - this function is internal to this library, do not use is elsewhere as that assumption won't hold anywhere else.
func collectCounts(pres []POSTResponseEntry) EntityCount {
	ec := EntityCount{}
	for i, _ := range pres {
		if i == 0 {
			// [0] below because they will only have one. (Guaranteed in getLongestUniqueTimespan)
			ec.Name = pres[i].EntityCounts[0].Name
			ec.Protocol = pres[i].EntityCounts[0].Protocol
		}
		ec.Count = ec.Count + pres[i].EntityCounts[0].Count
	}
	return ec
}

// Basic entry actions

// Eligible checks whether this post response entry is eligible for inclusion into a chain.
func (e *POSTResponseEntry) eligible() bool {
	expiryDur := time.Duration(bc.GetPOSTResponseExpiryMinutes()) * time.Minute
	ineligibilityDur := time.Duration(bc.GetPOSTResponseIneligibilityMinutes()) * time.Minute
	eligibleDur := expiryDur - ineligibilityDur
	lastEligibleMoment := time.Now().Add(-eligibleDur)
	if e.Creation > Timestamp(lastEligibleMoment.Unix()) {
		return true
	} else {
		return false
	}
}

// Basic repo actions

func (r *POSTResponseRepo) indexOf(url string) int {
	for i, _ := range r.Responses {
		if r.Responses[i].ResponseUrl == url {
			return i
		}
	}
	return -1
}

func (r *POSTResponseRepo) deleteFromRepo(i int) {
	newResp := append(r.Responses[0:i], r.Responses[i+1:len(r.Responses)]...)
	r.Responses = newResp
}

func (r *POSTResponseRepo) Add(url string, startsFrom Timestamp, endsAt Timestamp, creation Timestamp, entityCounts *[]EntityCount) {
	r.lock.Lock()
	defer r.lock.Unlock()
	pre := POSTResponseEntry{ResponseUrl: url, StartsFrom: startsFrom, EndsAt: endsAt, Creation: creation}
	pre.EntityCounts = *entityCounts
	r.Responses = append(r.Responses, pre)
	// fmt.Printf("Addded %#v to the response tracker repo. The repo is now: %#v", pre, r.Responses)
}

func (r *POSTResponseRepo) Remove(url string) {
	index := r.indexOf(url)
	if index != -1 { // if exists on the list
		r.deleteFromRepo(index)
	}
	deleteFromDisk(url)
}

// Flush deletes all post response entries and their representations on the disk up to the given timestamp. We count reverse because deleting from a list you're iterating on forwards will make the index number shift.
func (r *POSTResponseRepo) flush(cutoff Timestamp) {
	for i := len(r.Responses) - 1; i >= 0; i-- {
		// fmt.Println(i)
		if r.Responses[i].EndsAt < cutoff {
			r.deleteFromRepo(i)
		}
	}
}

// getLongestUniqueTimespan returns the POST response with the longest timespan starting from the t given. This t is usually the end of the last item of the chain, or if at the beginning, the beginning of the time range requested.
// If we sort by shortest timespan and then do this, given that their unique lengths are equal, the shorter timespanned response will be selected, because the only way to override the selection would be to be to have longer duration, not equal duration. But what's the performance improvement in both of them having the exat same unique duration? extremely unlikely benefit for constant, ongoing cost of sorting them in reverse by duration.
func (r *POSTResponseRepo) getLongestUniqueTimespan(t Timestamp, entityName string) *POSTResponseEntry {
	index := -1
	longestDuration := Timestamp(0) // This starts from 0, which means we do not allow for the chain to start after the given timestamp. It has to start before.
	for i, _ := range r.Responses {
		only1Ec := len(r.Responses[i].EntityCounts) == 1
		entityNameMatches := r.Responses[i].EntityCounts[0].Name == entityName
		hasEnoughTime := r.Responses[i].eligible()
		if !only1Ec || !entityNameMatches || !hasEnoughTime {
			// fmt.Printf("This post response was ineligible: Why: only1Ec: %#v, entityNameMatches: %#v, hasEnoughTime: %#v, Resp:%#v\n", only1Ec, entityNameMatches, hasEnoughTime, r.Responses[i])
			continue
		}
		// fmt.Printf("This post response is ELIGIBLE: %#v\n", r.Responses[i])
		uniqlength := r.Responses[i].EndsAt - t
		uniqlengthLongerThanLongestDur := uniqlength > longestDuration
		startsFromBeforeCutoff := t > r.Responses[i].StartsFrom
		// fmt.Printf("uniq length longer than longest duration: %#v, uniqlength: %#v, starts from before cutoff: %#v cutoff: %#v \n", uniqlengthLongerThanLongestDur, uniqlength, startsFromBeforeCutoff, t)
		if uniqlengthLongerThanLongestDur && startsFromBeforeCutoff {
			index = i
			longestDuration = r.Responses[i].EndsAt - t
		}
	}
	if index != -1 {
		return &r.Responses[index]
	} else {
		return &POSTResponseEntry{}
	}
}

// GetPostResponseChain constructs a POST response chain from generated post responses from prior that starts from before the given timeline and ends after the given end, with as few responses as possible.
func (r *POSTResponseRepo) GetPostResponseChain(start Timestamp, end Timestamp, entityName string) (*[]POSTResponseEntry, Timestamp, Timestamp, EntityCount) {
	r.lock.Lock()
	defer r.lock.Unlock()
	// fmt.Printf("This is all post responses available for chaining: %#v\n", r)
	// return &[]POSTResponseEntry{}, Timestamp(0), Timestamp(0), EntityCount{}
	var chain []POSTResponseEntry

	linkStart := start
	firstLinkStartsFrom := Timestamp(0)
	lastLinkEndsAt := Timestamp(0)
	ename := toolbox.Singular(entityName)
	// if end == 0 {
	// 	end = Timestamp(time.Now().Unix())
	// }
	for {
		chainLink := r.getLongestUniqueTimespan(linkStart, ename)
		if len((*chainLink).ResponseUrl) > 0 { // if it actually exists
			if linkStart == start {
				// This was the first ever link found, map its start to the first link start
				firstLinkStartsFrom = chainLink.StartsFrom
			}
			linkStart = chainLink.EndsAt
			chain = append(chain, *chainLink)
			lastLinkEndsAt = chainLink.EndsAt // any link can be last link if the next iteration doesn't work out.
			if chainLink.EndsAt > end {       // This chainLink went beyond our end, let's finalise the chain here.
				break
			}
		} else { // if it doesn't exist, this was the end of the link
			break
		}
	}
	count := collectCounts(chain)
	// fmt.Printf("This is the chain count we've found. Chain: %v, Count: %v\n", chain, count)
	return &chain, firstLinkStartsFrom, lastLinkEndsAt, count
}

func (r *POSTResponseRepo) Maintain() {
	r.lock.Lock()
	defer r.lock.Unlock()
	expiryDur := time.Duration(bc.GetPOSTResponseExpiryMinutes()) * time.Minute
	expiryThreshold := Timestamp(time.Now().Add(-expiryDur).Unix())
	r.flush(expiryThreshold)
}

func (r *POSTResponseRepo) DeleteAllFromDisk() {
	postDir := fmt.Sprintf("%s/", bc.GetProtURLVersion(), "/responses", bc.GetCachesDirectory())
	toolbox.DeleteFromDisk(postDir)
}
