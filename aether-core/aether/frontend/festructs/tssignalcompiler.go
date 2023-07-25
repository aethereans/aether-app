// Frontend > FEStructs > Truststate Signal Compiler
// This library is tasked with compiling high level truststate based signals into usable blocks that the refresher loops can consume.

package festructs

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/rollingbloom"
	// "github.com/willf/bloom"
)

/////////////////////////
// Compiled Truststate Signals
/////////////////////////

// Public trust (collection)
type CompiledPT struct {
	TargetFingerprint string
	PTs               []PublicTrustSignal
	LastRefreshed     int64
}

func NewCPT(targetfp string, nowts int64) *CompiledPT {
	return &CompiledPT{
		TargetFingerprint: targetfp,
		LastRefreshed:     nowts,
	}
}

func (c *CompiledPT) Insert(s PublicTrustSignal, nowts int64) {
	// Check if right bucket
	if c.TargetFingerprint != s.TargetFingerprint {
		return
	}
	// It's going in, one way or ther other, so make sure lastrefreshed is updated
	if c.LastRefreshed < s.LastRefreshed {
		c.LastRefreshed = s.LastRefreshed
	}
	// Check if already exists
	/*
		This used to have s.Expiry > nowts as a condition below. It's removed because it can cut off updates to an item that was not expired when it was updated, but expired in the time it took for it to arrive. It should still be ineffectual, but it should be ineffectual in the right state. We should probably prevent propagation of the expired truststates in the backend, though.
	*/
	for k := range c.PTs {
		if c.PTs[k].Fingerprint == s.Fingerprint &&
			s.LastRefreshed > c.PTs[k].LastRefreshed {
			c.PTs[k] = s
			return
		}
	}
	// It doesn't exist - insert it.
	c.PTs = append(c.PTs, s)
}

// CanonicalName (collection)
type CompiledCN struct {
	TargetFingerprint string
	CNs               []CanonicalNameSignal
	LastRefreshed     int64
}

func NewCCN(targetfp string, nowts int64) *CompiledCN {
	return &CompiledCN{
		TargetFingerprint: targetfp,
		LastRefreshed:     nowts,
	}
}

func (c *CompiledCN) Insert(s CanonicalNameSignal, nowts int64) {
	// Check if right bucket
	if c.TargetFingerprint != s.TargetFingerprint {
		return
	}
	// Check if already exists
	for k := range c.CNs {
		if c.CNs[k].Fingerprint != s.Fingerprint {
			continue
		}
		if c.CNs[k].LastRefreshed >= s.LastRefreshed {
			return
		}
		c.CNs[k] = s
		if c.LastRefreshed < s.LastRefreshed {
			c.LastRefreshed = s.LastRefreshed
		}
		return
	}
	// It doesn't exist - insert it.
	c.CNs = append(c.CNs, s)
	if c.LastRefreshed < s.LastRefreshed {
		c.LastRefreshed = s.LastRefreshed
	}
}

// F451 (collection)
type CompiledF451 struct {
	TargetFingerprint string
	F451s             []F451Signal
	LastRefreshed     int64
}

func NewCF451(targetfp string, nowts int64) *CompiledF451 {
	return &CompiledF451{
		TargetFingerprint: targetfp,
		LastRefreshed:     nowts,
	}
}

func (c *CompiledF451) Insert(s F451Signal, nowts int64) {
	// Check if right bucket
	if c.TargetFingerprint != s.TargetFingerprint {
		return
	}
	// Check if already exists
	for k := range c.F451s {
		if c.F451s[k].Fingerprint != s.Fingerprint {
			continue
		}
		if c.F451s[k].LastRefreshed >= s.LastRefreshed {
			return
		}
		c.F451s[k] = s
		if c.LastRefreshed < s.LastRefreshed {
			c.LastRefreshed = s.LastRefreshed
		}
		return
	}
	// It doesn't exist - insert it.
	c.F451s = append(c.F451s, s)
	if c.LastRefreshed < s.LastRefreshed {
		c.LastRefreshed = s.LastRefreshed
	}
}

// PublicElect (bloomed aggregate)
type CompiledPE struct {
	TargetFingerprint       string
	ElectsCount             int
	ElectsBloom             rollingbloom.RollingBloom
	DisqualifiesCount       int
	DisqualifiesBloom       rollingbloom.RollingBloom
	SelfPubliclyVoted       bool
	SelfPublicVoteDirection int
	SelfFingerprint         string
	SelfCreation            int64
	SelfLastUpdate          int64
	LastRefreshed           int64
}

func NewCPE(targetfp string, nowts int64) *CompiledPE {
	return &CompiledPE{
		TargetFingerprint: targetfp,
		ElectsBloom: rollingbloom.NewRollingBloom(
			uint(globals.FrontendConfig.GetNetworkMemoryDays()),
			uint(globals.FrontendConfig.GetNetworkHeadDays()),
			uint(globals.FrontendConfig.GetBloomFilterSize())),
		DisqualifiesBloom: rollingbloom.NewRollingBloom(
			uint(globals.FrontendConfig.GetNetworkMemoryDays()),
			uint(globals.FrontendConfig.GetNetworkHeadDays()),
			uint(globals.FrontendConfig.GetBloomFilterSize())),
		LastRefreshed: nowts,
	}
}

func (c *CompiledPE) Insert(s PublicElectSignal, nowts int64) {
	if c.TargetFingerprint != s.TargetFingerprint {
		logging.Logf(1, "You tried to apply a different entity's Public Elect signal to this CPE. PC's targetfp: %v, ATD's target fp: %v", c.TargetFingerprint, s.TargetFingerprint)
		return
	}
	if s.Expiry < nowts {
		return
	}
	inElectsBloom := c.ElectsBloom.TestString(s.SourceFingerprint)
	inDisqualifiesBloom := c.DisqualifiesBloom.TestString(s.SourceFingerprint)
	// If both matches
	if inElectsBloom && inDisqualifiesBloom {
		return // Can't do much here.
	}
	if inElectsBloom {
		// Only in elects bloom and ...
		if s.Type == Signal_Elect {
			// Signal is elect, which means ...
			return // ... we've already added this before.
		}
		// Signal is not elect, which means it flipped from elect to disqualify.
		c.ElectsCount--
		c.DisqualifiesCount++
		c.DisqualifiesBloom.AddString(s.SourceFingerprint)
	}
	if inDisqualifiesBloom {
		if s.Type == Signal_Disqualify {
			// We've already added this in
			return
		}
		// In DQ bloom but not a DQ signal, which means flipped from DQ to elect.
		c.DisqualifiesCount--
		c.ElectsCount++
		c.ElectsBloom.AddString(s.SourceFingerprint)
	}
	// In none
	if !inElectsBloom && !inDisqualifiesBloom {
		if s.Type == Signal_Elect {
			c.ElectsCount++
			c.ElectsBloom.AddString(s.SourceFingerprint)
		}
		if s.Type == Signal_Disqualify {
			c.DisqualifiesCount++
			c.DisqualifiesBloom.AddString(s.SourceFingerprint)
		}
	}
	// Add self flags, if needed.
	if s.Self && !c.SelfPubliclyVoted {
		c.SelfPubliclyVoted = true
		c.SelfPublicVoteDirection = s.Type
		c.SelfFingerprint = s.Fingerprint
		c.SelfCreation = s.Creation
		c.SelfLastUpdate = s.LastUpdate
	}
	// If any count has dropped below the threshold of zero, mark them back to zero
	if c.ElectsCount < 0 {
		c.ElectsCount = 0
	}
	if c.DisqualifiesCount < 0 {
		c.DisqualifiesCount = 0
	}
}

/////////////////////////
// Batched Compiled Signals
/////////////////////////

type CPTBatch []CompiledPT

func (cbatch *CPTBatch) Insert(pts []PublicTrustSignal, nowts int64) {
	for k := range pts {
		// cpt := cbatch.FindOrCreate(pts[k].TargetFingerprint)
		var cpt *CompiledPT
		i := cbatch.Find(pts[k].TargetFingerprint)
		if i == -1 {
			cpt = NewCPT(pts[k].TargetFingerprint, nowts)
			cpt.Insert(pts[k], nowts)
			*cbatch = append(*cbatch, *cpt)
		} else {
			cpt = &(*cbatch)[i]
			cpt.Insert(pts[k], nowts)
		}
	}
}

func (cbatch *CPTBatch) Find(targetfp string) int {
	for k := range *cbatch {
		if targetfp == (*cbatch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

func (cbatch *CPTBatch) FindObj(targetfp string) CompiledPT {
	i := cbatch.Find(targetfp)
	if i != -1 {
		return (*cbatch)[i]
	}
	return CompiledPT{}
}

type CCNBatch []CompiledCN

func (cbatch *CCNBatch) Insert(cns []CanonicalNameSignal, nowts int64) {
	for k := range cns {
		// ccn := cbatch.FindOrCreate(cns[k].TargetFingerprint)
		var ccn *CompiledCN
		i := cbatch.Find(cns[k].TargetFingerprint)
		if i == -1 {
			ccn = NewCCN(cns[k].TargetFingerprint, nowts)
			ccn.Insert(cns[k], nowts)
			*cbatch = append(*cbatch, *ccn)
		} else {
			ccn = &(*cbatch)[i]
			ccn.Insert(cns[k], nowts)
		}
	}
}

func (cbatch *CCNBatch) Find(targetfp string) int {
	for k := range *cbatch {
		if targetfp == (*cbatch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

func (cbatch *CCNBatch) FindObj(targetfp string) CompiledCN {
	i := cbatch.Find(targetfp)
	if i != -1 {
		return (*cbatch)[i]
	}
	return CompiledCN{}
}

type CF451Batch []CompiledF451

func (cbatch *CF451Batch) Insert(f451s []F451Signal, nowts int64) {
	for k := range f451s {
		// cf451 := cbatch.FindOrCreate(f451s[k].TargetFingerprint)
		var cf451 *CompiledF451
		i := cbatch.Find(f451s[k].TargetFingerprint)
		if i == -1 {
			cf451 = NewCF451(f451s[k].TargetFingerprint, nowts)
			cf451.Insert(f451s[k], nowts)
			*cbatch = append(*cbatch, *cf451)
		} else {
			cf451 = &(*cbatch)[i]
			cf451.Insert(f451s[k], nowts)
		}
	}
}

func (cbatch *CF451Batch) Find(targetfp string) int {
	for k := range *cbatch {
		if targetfp == (*cbatch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

func (cbatch *CF451Batch) FindObj(targetfp string) CompiledF451 {
	i := cbatch.Find(targetfp)
	if i != -1 {
		return (*cbatch)[i]
	}
	return CompiledF451{}
}

type CPEBatch []CompiledPE

func (cbatch *CPEBatch) Insert(pes []PublicElectSignal, nowts int64) {
	for k := range pes {
		var cpe *CompiledPE
		i := cbatch.Find(pes[k].TargetFingerprint)
		if i == -1 {
			cpe = NewCPE(pes[k].TargetFingerprint, nowts)
			cpe.Insert(pes[k], nowts)
			*cbatch = append(*cbatch, *cpe)
		} else {
			cpe = &(*cbatch)[i]
			cpe.Insert(pes[k], nowts)
		}
	}
}

func (cbatch *CPEBatch) Find(targetfp string) int {
	for k := range *cbatch {
		if targetfp == (*cbatch)[k].TargetFingerprint {
			return k
		}
	}
	return -1
}

func (cbatch *CPEBatch) FindObj(targetfp string) CompiledPE {
	i := cbatch.Find(targetfp)
	if i != -1 {
		return (*cbatch)[i]
	}
	return CompiledPE{}
}
