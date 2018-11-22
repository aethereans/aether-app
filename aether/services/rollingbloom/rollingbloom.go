// Services > RollingBloom
// This service provides a bloom filter that is time limited, in that you can ask 'have we seen this thing in the last x months?'

// The way this works is that you give this bloom filter a duration, and it makes it so that you can query for anything within that duration. Let's say it's 6 months. At any point in time, when you ask this filter, it will tell you whether what you're testing for is in the list, and how many total entries are in the bloom filter in total.

// What it cannot do is that you cannot tell this guy to give you the count of last 3 months when the filter was for 6 months. In exchange to that, the count this has is pessimistic and exact, in that it only adds to the bloom if we know for sure that it has not been put into the filter in the given timeframe. This means it might actually not add some data because it returned 'might be, might not be', so the count will always be lower than the real count. The diference between the real count and the count the filter reports increases as the bloom filter approaches its capacity, so you should have the filter about 10x the size you think you need. This is a bloom filter calculator that can give you an idea on how to set your values. https://hur.st/bloomfilter/ Mind that this count is for one filter, and rollingbloom creates duration / granularity. So if your duration is 30 days and granularity is 10 days, it'll be 3 bloom filters (+1 for the immediate future, so 4). That means if you set your capacity is 10000 at 50% false positive rate, your rolling bloom capacity is 30000 at 50% assuming a perfectly even distribution.

// This way of implementing a rolling bloom trades being able to query a specific time range within the bloom for exactitude on the lowest boundary. It might not count everything (but it will count 99%+ for most of its capacity), but it is guaranteed to not double count anything twice. This is useful for elections, which are conviction checks.

// There is another way to implement this which I'll probably also do, and that way, the bloom filter is capable of answering the queries on a limited, specific time range, not just its whole duration. That way, we always add it to the latest bloom without checking, but it doesn't have a count that increments. When the count query comes in, you union the bloom filters in the given time range, and do an estimation on the bloom filter to calculate the statistically most likely count number. This means the count number will be inexact, but you'll be able to get a count for any duration range within your covered total time range. That's more useful for population counts, which are existence checks.

/*

"Lies, damned lies, and statistics." — Benjamin Disraeli

*/

package rollingbloom

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"github.com/willf/bloom"
	"sync"
	"time"
)

type constituentBloom struct {
	StartTimestamp int64
	EndTimestamp   int64
	Count          int
	Bloom          bloom.BloomFilter
}

func (cb *constituentBloom) AddString(str string) {
	cb.Bloom.AddString(str)
	cb.Count++
}

func (cb *constituentBloom) TestString(str string) bool {
	return cb.Bloom.TestString(str)
}

func newConstituentBloom(durationDays, maxSize uint, start int64) constituentBloom {
	return constituentBloom{
		StartTimestamp: start,
		EndTimestamp:   time.Unix(start, 0).Add(time.Duration(int(durationDays)) * time.Hour * 24).Unix(),
		Bloom:          *bloom.NewWithEstimates(maxSize, float64(globals.FrontendConfig.GetBloomFilterFalsePositiveRatePercent())/100),
	}
}

type RollingBloom struct {
	lock              sync.Mutex
	ConstituentBlooms []constituentBloom
	MaxDurationDays   uint
	GranularityDays   uint
	MaxSize           uint
	lastMaintainRun   int64
}

func (r RollingBloom) String() string {
	return "" // Disable printing of internals.
}

// AddString adds a value to the bloom (to the most recent constituent).
func (r *RollingBloom) AddString(str string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.maintain()
	cb := r.getCurrentConstituentBloom()
	if r.teststr(str) {
		return
		// This check is very important, because we are keeping counts in the bloom filter. These counts increment at every add, so it's very important that every add is actually for something that is not in the bloom.
	}
	cb.AddString(str)
}

// TestString checks whether the value is possibly within *any* of the constituent blooms. This is an important detail. When you're adding something, you should check first, because if you add without checking, you're going to increment the count without actually adding anything to the bloom.
func (r *RollingBloom) TestString(str string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.maintain()
	return r.teststr(str)
}

// Count is guaranteed to not double count anything, because when we add things, we check *all* constituent blooms, and if it is present in any, we don't add.
func (r *RollingBloom) Count() int {
	var c int
	for k, _ := range r.ConstituentBlooms {
		c += r.ConstituentBlooms[k].Count
	}
	return c
}

// NewRollingBloom creates a bloom filter that can keep track of a rolling, but limited time. Mind that maxSize is the size of every single constituent bloom - so if your maximum duration is 180 days and your resolution is 14 day blocks with a max size of 10000, you'll have 13~ blocks, each of which can hold 10000. Your total capacity at 50% fail rate and perfect distribution will be 130000. Since this calculation can depend on many things, you should go for a size an order of magnitude larger than you think you'll need.
func NewRollingBloom(maxDurationDays, GranularityDays, maxSize uint) RollingBloom {
	if maxDurationDays < GranularityDays || maxDurationDays == 0 || GranularityDays == 0 {
		logging.LogCrash("You've provided an invalid combination for RollingBloom. This is a programming error.")
	}
	rb := RollingBloom{
		MaxDurationDays: maxDurationDays,
		GranularityDays: GranularityDays,
		MaxSize:         maxSize,
	}
	rb.maintain()
	return rb
}

func (r *RollingBloom) teststr(str string) bool {
	for k, _ := range r.ConstituentBlooms {
		if r.ConstituentBlooms[k].TestString(str) {
			return true
		}
	}
	return false
}

// delete too old blooms and open new ones as needed.
func (r *RollingBloom) maintain() {
	/*
		Below is a bugfix for a prior release that accidentally set granularity to 0. After releasing one version with this, we should remove it afterwards. TODO
	*/
	if r.GranularityDays == 0 {
		r.GranularityDays = 14
	}
	now := time.Now()
	if r.lastMaintainRun > now.Add(-23*time.Hour).Unix() {
		return
		// 23 hours because even in the highest resolution (1 day per bloom) lastMaintainRun cannot possibly prevent maintain() in a case where a most recent bloom that covers now() does not exist.
	}
	r.lastMaintainRun = now.Unix()
	cleanedCBlooms := []constituentBloom{}
	var lastBloomEnd int64
	for k, _ := range r.ConstituentBlooms {
		// Move to the new list only if the new bloom is still valid in duration
		if b := r.ConstituentBlooms[k]; b.EndTimestamp > now.Add(-24*time.Hour*time.Duration(int(r.MaxDurationDays))).Unix() {
			// This bloom is valid.
			if b.EndTimestamp > lastBloomEnd {
				// If this is the most recent time this bloom ends save it.
				lastBloomEnd = b.EndTimestamp
			}
			cleanedCBlooms = append(cleanedCBlooms, b)
		}
	}
	lastNewBloomEnd := lastBloomEnd
	// If last bloom end is older than our max duration, bring it back to our max duration.
	if cutoff := now.Add(-(time.Duration(int(r.MaxDurationDays)) * 24 * time.Hour)).Unix(); lastNewBloomEnd < cutoff {
		lastNewBloomEnd = cutoff
	}
	// Generate bloom buckets until the end reaches beyond one interval into the future. This is to prevent a case where the max duration and resolution are multiples of each other, and the last generated constitutent bloom terminates right at the moment of creation - with the lastMaintainRun gate preventing a new run until a day after. This way, the end date of the bloom filter will at least be one more cycle into the future.
	for now.Add(24*time.Hour*time.Duration(int(r.GranularityDays))).Unix() > lastNewBloomEnd {
		cb := newConstituentBloom(r.GranularityDays, r.MaxSize, lastNewBloomEnd)
		lastNewBloomEnd = cb.EndTimestamp
		r.ConstituentBlooms = append(r.ConstituentBlooms, cb)
	}
}

func (r *RollingBloom) getCurrentConstituentBloom() *constituentBloom {
	now := time.Now().Unix()
	for k, _ := range r.ConstituentBlooms {
		if r.ConstituentBlooms[k].EndTimestamp > now {
			return &r.ConstituentBlooms[k]
		}
	}
	logging.LogCrash("This rolling bloom didn't have a current constituent ready when asked. This is a programming error.")
	return nil
}
