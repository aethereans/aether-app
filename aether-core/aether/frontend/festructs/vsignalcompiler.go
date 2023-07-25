// Frontend > FEStructs > Vote Signal Compiler
// This library is tasked with compiling high level vote based signals into usable blocks that the refresher loops can consume.

package festructs

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"

	"github.com/willf/bloom"
)

/////////////////////////
// Compiled Vote Signals
/////////////////////////

// CompiledATD is the upvotes and downvotes of a single post. This is generated from multiple ATDs (bloomed aggregate)
type CompiledATD struct {
	TargetFingerprint string `storm:"id"`
	UpvotesCount      int
	UpvotesBloom      bloom.BloomFilter
	DownvotesCount    int
	DownvotesBloom    bloom.BloomFilter
	// This one has self vote marker and the other ones don't, because since we aggregate it into the bloom, we can't really know with certainty that the user has voted on it. We could test the bloom filter against the local user public key, but that would give us 'maybe in it', not 'in it'. 'Maybe' is not good enough when it comes to user's own data.
	SelfVoted         bool
	SelfVoteDirection int
	SelfFingerprint   string
	SelfCreation      int64
	SelfLastUpdate    int64
	LastRefreshed     int64
}

func NewCATD(targetfp string, nowts int64) *CompiledATD {
	return &CompiledATD{
		TargetFingerprint: targetfp,
		UpvotesBloom:      *bloom.NewWithEstimates(uint(globals.FrontendConfig.GetBloomFilterSize()), float64(globals.FrontendConfig.GetBloomFilterFalsePositiveRatePercent())/100),
		DownvotesBloom:    *bloom.NewWithEstimates(uint(globals.FrontendConfig.GetBloomFilterSize()), float64(globals.FrontendConfig.GetBloomFilterFalsePositiveRatePercent())/100),
		LastRefreshed:     nowts,
	}
}

func (c *CompiledATD) Insert(atd AddsToDiscussionSignal) {
	if c.TargetFingerprint != atd.TargetFingerprint {
		logging.Logf(1, "You tried to apply a different entity's ATD to this CATD. CATD's targetfp: %v, ATD's target fp: %v", c.TargetFingerprint, atd.TargetFingerprint)
		return
	}
	/*
	   Check both bloom filters.

	   If none matches, all is well, insert it in.

	   If it's present in one, and it's the one the direction matches, this might be a double entry (or bloom filter being probabilistic). Don't do anything.

	   If it's present in one, and it's not the direction that is in the one being added, the user has changed opinion. Drop one point from the opposite, add one point to the new one.

	   If both matches, it might be the bloom filter being probabilistic, or the user has changed opinion more than once. In that case, we don't know which, and we can't do anything. You only get one shot at changing your vote.

	   Our bloom filters are tuned so that they will expect about 10k items as their max capacity. That makes every bloom filter cost us 1kb. It's a little expensive, but it's much cheaper than holding every single up/downvote in kvstore.

	   That means, though, the vote accrual will slow down as it gets closer to 10k, as the false positive rates increase, and stop at somewhere around 10k. Eventually, we'll have to raise this ceiling to 100k or something similar, based on what we see in the network.
	*/
	// Check whether ATD's own fingerprint (not its target fingerprint) is in the bloom filter.

	inUpvotesBloom := c.UpvotesBloom.TestString(atd.SourceFingerprint)
	inDownvotesBloom := c.DownvotesBloom.TestString(atd.SourceFingerprint)
	logging.Logf(2, "PRE-MODIFY upvotes and downvotes counts: U: %v, D: %v ATD Target: %#v", c.UpvotesCount, c.DownvotesCount, c.TargetFingerprint)
	logging.Logf(2, "In upvotes bloom: %v, In downvotes bloom: %v, Fingerprint: %v", inUpvotesBloom, inDownvotesBloom, c.TargetFingerprint)
	////////////////
	// Both matches.
	////////////////
	if inUpvotesBloom && inDownvotesBloom {
		return // Can't do much here
	}
	// Only one bloom matches.
	if inUpvotesBloom {
		if atd.Type == Signal_Upvote {
			logging.Logf(2, "We already counted this upvote.")
			return // Already added
		}
		// Upvote to downvote flip.
		c.UpvotesCount--
		c.DownvotesCount++
		c.DownvotesBloom.AddString(atd.SourceFingerprint)
	}
	if inDownvotesBloom {
		if atd.Type == Signal_Downvote {
			logging.Logf(2, "We already counted this downvote.")
			return // Already added
		}
		// Downvote to upvote flip.
		c.DownvotesCount--
		c.UpvotesCount++
		c.UpvotesBloom.AddString(atd.SourceFingerprint)
	}
	// None matches.
	if !inUpvotesBloom && !inDownvotesBloom {
		if atd.Type == Signal_Upvote {
			c.UpvotesBloom.AddString(atd.SourceFingerprint)
			c.UpvotesCount++
		}
		if atd.Type == Signal_Downvote {
			c.DownvotesBloom.AddString(atd.SourceFingerprint)
			c.DownvotesCount++
		}
	}
	// If this is a self vote, set the flag for it.
	// if atd.Self && !c.SelfVoted {
	// 	c.SelfVoted = true
	// 	c.SelfVoteDirection = atd.Type
	// 	c.SelfFingerprint = atd.BaseVoteSignal.Fingerprint
	// 	c.SelfCreation = atd.BaseVoteSignal.Creation
	// 	c.SelfLastUpdate = atd.BaseVoteSignal.LastUpdate
	// }
	// ^ This used to block updates at the !c.SelfVoted part. prevented that.
	if atd.Self {
		c.SelfVoted = true
		c.SelfVoteDirection = atd.Type
		c.SelfFingerprint = atd.BaseVoteSignal.Fingerprint
		c.SelfCreation = atd.BaseVoteSignal.Creation
		c.SelfLastUpdate = atd.BaseVoteSignal.LastUpdate
	}
	if c.LastRefreshed < atd.LastRefreshed {
		c.LastRefreshed = atd.LastRefreshed
	}
	// Finally, if they fell below 0, make them zero. This can happen because bloom filters are probabilistic and can return 'maybe', in which case we decrement / increment, but not all of those maybes are actual positives.
	if c.UpvotesCount < 0 {
		c.UpvotesCount = 0
	}
	if c.DownvotesCount < 0 {
		c.DownvotesCount = 0
	}
	logging.Logf(2, "           Upvotes and downvotes counts: U: %v, D: %v ATD Target: %#v", c.UpvotesCount, c.DownvotesCount, c.TargetFingerprint)
}

// Compiled Follows Guidelines
type CompiledFG struct {
	TargetFingerprint string
	FGs               []FollowsGuidelinesSignal
	SelfReported      bool
	LastRefreshed     int64
}

func NewCFG(targetfp string, nowts int64) *CompiledFG {
	return &CompiledFG{
		TargetFingerprint: targetfp,
		LastRefreshed:     nowts,
	}
}

func (c *CompiledFG) Insert(sg FollowsGuidelinesSignal) {
	// Check if it's the right bucket first. If not, bail
	if c.TargetFingerprint != sg.BaseVoteSignal.TargetFingerprint {
		return
	}
	// If it already exists in the list, overwrite existing
	for k := range c.FGs {
		if c.FGs[k].Fingerprint != sg.Fingerprint {
			continue
		}
		// If it exists in the list, and is newer,
		if c.FGs[k].BaseVoteSignal.LastRefreshed >= sg.BaseVoteSignal.LastRefreshed {
			// If older than what we have, bail
			return
		}
		c.FGs[k] = sg // overwrite existing if the new one is newer
		if sg.BaseVoteSignal.Self {
			c.SelfReported = true
		}
		if c.LastRefreshed < sg.LastRefreshed {
			c.LastRefreshed = sg.LastRefreshed
		}
		return
	}
	// It's not extant. We'll be adding it if this is the right bucket
	c.FGs = append(c.FGs, sg)
	if sg.BaseVoteSignal.Self {
		c.SelfReported = true
	}
	if c.LastRefreshed < sg.LastRefreshed {
		c.LastRefreshed = sg.LastRefreshed
	}
}

// Compiled Mod Actions

type CompiledMA struct {
	TargetFingerprint string
	MAs               []ModActionsSignal
	SelfModBlocked    bool
	SelfModApproved   bool
	LastRefreshed     int64
}

func NewCMA(targetfp string, nowts int64) *CompiledMA {
	return &CompiledMA{
		TargetFingerprint: targetfp,
		LastRefreshed:     nowts,
	}
}

func (c *CompiledMA) Insert(sg ModActionsSignal) {
	// Check if it's the right bucket first. If not, bail
	if c.TargetFingerprint != sg.BaseVoteSignal.TargetFingerprint {
		return
	}
	// If it already exists in the list, overwrite existing
	for k := range c.MAs {
		if c.MAs[k].Fingerprint != sg.Fingerprint {
			continue
		}
		// If it exists in the list, and is newer,
		if c.MAs[k].BaseVoteSignal.LastRefreshed >= sg.BaseVoteSignal.LastRefreshed {
			// If older than what we have, bail
			return
		}
		c.MAs[k] = sg // overwrite existing if the new one is newer
		if sg.BaseVoteSignal.Self {
			// You can have selfmodblocked and selfmodapproved at the same time. It doesn't make sense, but hey, people don't make sense either.
			if sg.BaseVoteSignal.Type == 1 {
				c.SelfModBlocked = true
			}
			if sg.BaseVoteSignal.Type == 2 {
				c.SelfModApproved = true
			}
		}
		if c.LastRefreshed < sg.LastRefreshed {
			c.LastRefreshed = sg.LastRefreshed
		}
		return
	}
	// It's not extant. We'll be adding it if this is the right bucket
	c.MAs = append(c.MAs, sg)
	if sg.BaseVoteSignal.Self {
		// You can have selfmodblocked and selfmodapproved at the same time. It doesn't make sense, but hey, people don't make sense either.
		if sg.BaseVoteSignal.Type == 1 {
			c.SelfModBlocked = true
		}
		if sg.BaseVoteSignal.Type == 2 {
			c.SelfModApproved = true
		}
	}
	if c.LastRefreshed < sg.LastRefreshed {
		c.LastRefreshed = sg.LastRefreshed
	}
}

/////////////////////////
// Batched Compiled Signals
/////////////////////////

// CATDBatch is a collection of compiled ATDs, each of which is the probabilistic aggregation of upvotes and downvotes for a specific entity.
type CATDBatch []CompiledATD

func (cbatch *CATDBatch) Insert(atds []AddsToDiscussionSignal, nowts int64) {
	for k := range atds {
		// catd := cbatch.FindOrCreate(atds[k].TargetFingerprint)
		var catd *CompiledATD
		i := cbatch.Find(atds[k].TargetFingerprint)
		if i == -1 {
			catd = NewCATD(atds[k].TargetFingerprint, nowts)
			catd.Insert(atds[k])
			*cbatch = append(*cbatch, *catd)
		} else {
			catd = &(*cbatch)[i]
			catd.Insert(atds[k])
		}
	}
}

func (cbatch *CATDBatch) Find(targetfp string) int {
	for k := range *cbatch {
		if targetfp == (*cbatch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

// CFGBatch is a collection of compiled CFGs.
type CFGBatch []CompiledFG

func (batch *CFGBatch) Insert(signals []FollowsGuidelinesSignal, nowts int64) {

	for k := range signals {
		var compiledSignal *CompiledFG
		i := batch.Find(signals[k].TargetFingerprint)
		if i == -1 {
			compiledSignal = NewCFG(signals[k].TargetFingerprint, nowts)
			compiledSignal.Insert(signals[k])
			*batch = append(*batch, *compiledSignal)
		} else {
			compiledSignal = &(*batch)[i]
			compiledSignal.Insert(signals[k])
		}
	}
}

func (batch *CFGBatch) Find(targetfp string) int {
	for k := range *batch {
		if targetfp == (*batch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

// CMABatch is a collection of compiled MAs.
type CMABatch []CompiledMA

func (batch *CMABatch) Insert(signals []ModActionsSignal, nowts int64) {
	for k := range signals {
		var compiledSignal *CompiledMA
		i := batch.Find(signals[k].TargetFingerprint)
		if i == -1 {
			compiledSignal = NewCMA(signals[k].TargetFingerprint, nowts)
			compiledSignal.Insert(signals[k])
			*batch = append(*batch, *compiledSignal)
		} else {
			compiledSignal = &(*batch)[i]
			compiledSignal.Insert(signals[k])
		}
	}
}

func (batch *CMABatch) Find(targetfp string) int {
	for k := range *batch {
		if targetfp == (*batch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}
