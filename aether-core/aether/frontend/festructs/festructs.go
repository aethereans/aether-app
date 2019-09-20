// Frontend > FeStructs
// These are the struct definitions that are shared. Everything should be defined here so that we can import this package from many places. Be very careful adding application-level imports here (golang native packages and external dependencies are fine.)

package festructs

import (
	"aether-core/aether/frontend/search"
	"aether-core/aether/protos/feobjects"
	"aether-core/aether/services/ca"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "github.com/willf/bloom"
	pbstructs "aether-core/aether/protos/mimapi"
	"math"
	"sort"
	"sync"
	"time"
)

var (
	CBoardIndexCache  = CBoardBatch{}
	CThreadIndexCache = CThreadBatch{}
	CPostIndexCache   = CPostBatch{}
	CUserIndexCache   = CUserBatch{}
)

func CommitSearchIndexes(wg *sync.WaitGroup) {
	defer wg.Done()
	logging.Logf(2, "Commit search indexes enters. We have %v boards, %v threads, %v posts, %v users.", len(CBoardIndexCache), len(CThreadIndexCache), len(CPostIndexCache), len(CUserIndexCache))

	start := time.Now()
	batchSize := 500

	/* Batch boards */
	var bb []CBoardBatch
	for batchSize < len(CBoardIndexCache) {
		CBoardIndexCache, bb = CBoardIndexCache[batchSize:], append(bb, CBoardIndexCache[0:batchSize:batchSize])
	}
	bb = append(bb, CBoardIndexCache)

	/* Batch threads */
	var tb []CThreadBatch
	for batchSize < len(CThreadIndexCache) {
		CThreadIndexCache, tb = CThreadIndexCache[batchSize:], append(tb, CThreadIndexCache[0:batchSize:batchSize])
	}
	tb = append(tb, CThreadIndexCache)

	/* Batch posts */
	var pb []CPostBatch
	for batchSize < len(CPostIndexCache) {
		CPostIndexCache, pb = CPostIndexCache[batchSize:], append(pb, CPostIndexCache[0:batchSize:batchSize])
	}
	pb = append(pb, CPostIndexCache)

	/* Batch users */
	var ub []CUserBatch
	for batchSize < len(CUserIndexCache) {
		CUserIndexCache, ub = CUserIndexCache[batchSize:], append(ub, CUserIndexCache[0:batchSize:batchSize])
	}
	ub = append(ub, CUserIndexCache)

	/* Commit board batches */
	for batchIndex, _ := range bb {
		ib := search.NewBatch()
		for k, _ := range bb[batchIndex] {
			ib.Index(bb[batchIndex][k].SearchId(), bb[batchIndex][k])
		}
		err := search.CommitBatch(ib)
		if err != nil {
			logging.LogCrash(err)
		}
	}

	/* Commit thread batches */
	for batchIndex, _ := range tb {
		ib := search.NewBatch()
		for k, _ := range tb[batchIndex] {
			ib.Index(tb[batchIndex][k].SearchId(), tb[batchIndex][k])
		}
		err := search.CommitBatch(ib)
		if err != nil {
			logging.LogCrash(err)
		}
	}

	/* Commit post batches */
	for batchIndex, _ := range pb {
		ib := search.NewBatch()
		for k, _ := range pb[batchIndex] {
			ib.Index(pb[batchIndex][k].SearchId(), pb[batchIndex][k])
		}
		err := search.CommitBatch(ib)
		if err != nil {
			logging.LogCrash(err)
		}
	}

	/* Commit user batches */
	for batchIndex, _ := range ub {
		ib := search.NewBatch()
		for k, _ := range ub[batchIndex] {
			ib.Index(ub[batchIndex][k].SearchId(), ub[batchIndex][k])
		}
		err := search.CommitBatch(ib)
		if err != nil {
			logging.LogCrash(err)
		}
	}

	elapsed := time.Since(start)
	logging.Logf(2, "Indexing is complete for the whole batch, board, thread, post, user. Count: %v. Took: %v", len(CBoardIndexCache)+len(CThreadIndexCache)+len(CPostIndexCache)+len(CUserIndexCache), elapsed)
	CBoardIndexCache = CBoardBatch{}
	CThreadIndexCache = CThreadBatch{}
	CPostIndexCache = CPostBatch{}
	CUserIndexCache = CUserBatch{}
}

// Compiled types

type CompiledPost struct {
	Fingerprint            string `storm:"id"`
	Board                  string
	Thread                 string
	Parent                 string
	SelfCreated            bool
	Body                   string
	CompiledContentSignals CompiledContentSignals
	Owner                  CompiledUser
	Creation               int64
	LastUpdate             int64
	Meta                   string
}

// BleveType satisfies the bleve Classifier interface so that Bleve knows how to parse this to index for search.
func (c CompiledPost) BleveType() string {
	return "post"
}

// SearchId gives the Id on which we'll save this entity. This has enough data to find what we want in a fast manner, and serve it to the user.
func (c CompiledPost) SearchId() string {
	// A path that defines a post is board, thread, parent, self fp. No user fp.
	bid, err := search.MakeSearchId("Post", c.Board, c.Thread, c.Parent, c.Fingerprint, "")
	if err != nil {
		logging.Logf(1, "Making search ID failed. Error: %v", err)
		return ""
	}
	return bid
}

func (c *CompiledPost) IndexForSearch() {
	// search.Index(c.SearchId(), c)
	CPostIndexCache = append(CPostIndexCache, *c)
}

func (c *CompiledPost) DeleteFromSearchIndex() {
	search.Delete(c.SearchId())
}

func NewCPost(rp *pbstructs.Post) CompiledPost {
	return CompiledPost{
		Fingerprint: rp.GetProvable().GetFingerprint(),
		Board:       rp.GetBoard(),
		Thread:      rp.GetThread(),
		Parent:      rp.GetParent(),
		SelfCreated: rp.GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
		Body:        rp.GetBody(),
		Meta:        rp.GetMeta(),
		// Half-baked ones
		Owner: CompiledUser{
			Fingerprint: rp.GetOwner(),
		},
		Creation:   rp.GetProvable().GetCreation(),
		LastUpdate: rp.GetUpdateable().GetLastUpdate(),
	}
	// Needs: Compiledcontentsignals, owner, bymod, byop, blocked, approved flags
}

// RefreshContentSignals refreshes an existing compiled post's compileduser and signals.
func (c *CompiledPost) RefreshContentSignals(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, nowts int64) {
	cs := CompiledContentSignals{}
	cs.Insert(c.Fingerprint, catds, cfgs, cmas, nowts)
	// Move values that need to be retained
	cs.SelfModIgnored = c.CompiledContentSignals.SelfModIgnored
	c.CompiledContentSignals = cs
}

func (c *CompiledPost) RefreshUserHeader(boardSpecificUserHeaders CUserBatch) {
	for k, _ := range boardSpecificUserHeaders {
		if boardSpecificUserHeaders[k].Fingerprint == c.Owner.Fingerprint {
			c.Owner = boardSpecificUserHeaders[k]
			return
		}
	}
	// Not found in the local headers. Seek from global headers.
	uhc := UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in RefreshContentSignals>One")
	err := globals.KvInstance.One("Fingerprint", c.Owner.Fingerprint, &uhc)
	if err != nil {
		logging.Logf(2, "We could not find the user of this entity in local or global user headers scopes. Entity Fingerprint: %#v, User Fingerprint: %v", c.Fingerprint, c.Owner.Fingerprint)
		return
	}
	i := uhc.Users.Find(c.Owner.Fingerprint)
	if i != -1 {
		c.Owner = uhc.Users[i]
		return
	}
	logging.LogCrashf("This should never happen. We've found a user header carrier that matches the user fingerprint requested by this entity, but within the UHC, there was no CompiledUser matching the fingerprint. Entity: %#v, User Fingerprint: %v, UserHeaderCarrier: %#v", c, c.Owner.Fingerprint, uhc)
}

func (c *CompiledPost) Insert(ce CompiledPost) {
	if c.LastUpdate < ce.LastUpdate {
		*c = ce
		c.IndexForSearch()
	}
}

func (c *CompiledPost) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier, tc *ThreadCarrier) {
	c.RefreshUserHeader(boardSpecificUserHeaders)
	c.RefreshContentSignals(catds, cfgs, cmas, nowts)
	c.RefreshExogenousContentSignals(bc, tc)
}

// RefreshExogenousContentSignals is where we compile and calculate the content signals that depend on external entitites.
func (c *CompiledPost) RefreshExogenousContentSignals(bc *BoardCarrier, tc *ThreadCarrier) {
	/*
		The signals processed here are:
		ByMod            bool
		ByFollowedPerson bool
		ByBlockedPerson  bool
		ByOP             bool
		ModBlocked       bool
		ModApproved      bool

		creator's userheader: (ask to own userheader)
		- following / blocked state,
		- bymod state

		thread owner fingerprint (ask to thread carrier)
		- byop state

		issuer's userheader (ask to board carrier)
		- modblock / modapprove state
	*/

	us := &c.Owner.CompiledUserSignals

	/*----------  ByMod state  ----------*/
	c.CompiledContentSignals.ByMod = isMod(us)

	/*----------  Following / blocked state  ----------*/
	c.CompiledContentSignals.ByFollowedPerson = us.FollowedBySelf
	c.CompiledContentSignals.ByBlockedPerson = us.BlockedBySelf

	/*----------  ByOp state  ----------*/
	if tc != nil {
		// ^ Nil gate, because this can be also used by a thread, not just a post
		if tc.Threads[0].Owner.Fingerprint == c.Owner.Fingerprint {
			c.CompiledContentSignals.ByOP = true
		}
	}
	/*----------  Modblock / modapprove state  ----------*/
	// Behaviour: if at least one modblock, block it, if there is at least one modapprove, unblock it. so if something is both modblocked and modapproved, it will be visible.
	// Approvals
	for k, _ := range c.CompiledContentSignals.ModApprovals {
		sourcefp := c.CompiledContentSignals.ModApprovals[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModApproved = true
		}
	}
	// Blocks
	for k, _ := range c.CompiledContentSignals.ModBlocks {
		sourcefp := c.CompiledContentSignals.ModBlocks[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModBlocked = true
		}
	}
}

func isMod(us *CompiledUserSignals) bool {
	/*
		These signals have different levels of weight. The order is like this, from weakest to strongest:
		- Default made mod
		- Network made mod
		- Network made non-mod
		- Self made mod
		- Self made non-mod
	*/
	isMod := false
	if us.MadeModByDefault {
		isMod = true
	}
	if us.MadeModByNetwork {
		isMod = true
	}
	if us.MadeNonModByNetwork {
		isMod = false
	}
	if us.MadeModBySelf {
		isMod = true
	}
	if us.MadeNonModBySelf {
		isMod = false
	}
	return isMod
}

// Batch compiled sets

type CPostBatch []CompiledPost

// IndexForSearch adds all entities in this batch into the search index.
func (batch *CPostBatch) IndexForSearch() {
	for k, _ := range *batch {
		CPostIndexCache = append(CPostIndexCache, (*batch)[k])
	}
	// if len(*batch) == 0 {
	// 	return
	// }
	// go func() {
	// 	ib := search.NewBatch()
	// 	for k, _ := range *batch {
	// 		ib.Index((*batch)[k].SearchId(), (*batch)[k])
	// 	}
	// 	err := search.CommitBatch(ib)
	// 	if err != nil {
	// 		logging.LogCrash(err)
	// 	}
	// 	logging.Logf(1, "Indexing is complete for this post batch. Count: %v", len(*batch))
	// }()
}

func (batch *CPostBatch) Insert(ces []CompiledPost) {
	toBeIndexed := CPostBatch{}
	for k, _ := range ces {
		i := batch.Find(ces[k].Fingerprint)
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(ces[k])
			continue
		}
		// Doesn't exist in our batch. Add it.
		*batch = append(*batch, ces[k])
		toBeIndexed = append(toBeIndexed, ces[k])
	}
	// Index for new additions. (Index for edits happen inside singular insert.)
	toBeIndexed.IndexForSearch()
}

func (batch *CPostBatch) InsertFromProtobuf(ces []*pbstructs.Post) {
	toBeIndexed := CPostBatch{}
	for k, _ := range ces {
		if ces[k] == nil {
			continue
		}
		i := batch.Find(ces[k].GetProvable().GetFingerprint())
		if i != -1 {
			// Exists in our batch. Check lastUpdate, if new the new one is more recent, update it.
			if (*batch)[i].LastUpdate < ces[k].GetUpdateable().GetLastUpdate() {
				// Trigger insert. It'll only update if the lastupdate is newer.
				(*batch)[i].Insert(NewCPost(ces[k]))
			}
			continue
		}
		// Doesn't exist in our batch. Add it.
		newpost := NewCPost(ces[k])
		*batch = append(*batch, newpost)
		toBeIndexed = append(toBeIndexed, newpost)
	}
	toBeIndexed.IndexForSearch()
}

func (batch *CPostBatch) Find(postfp string) int {
	for k, _ := range *batch {
		if postfp == (*batch)[k].Fingerprint {
			return k
		}
	}
	return -1
}

func (batch *CPostBatch) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier, tc *ThreadCarrier) {
	for k, _ := range *batch {
		(*batch)[k].Refresh(catds, cfgs, cmas, boardSpecificUserHeaders, nowts, bc, tc)
	}
}

func (batch *CPostBatch) Sort() {
	sort.Slice((*batch), func(i, j int) bool {
		iVoteDelta := ((*batch)[i].CompiledContentSignals.Upvotes - (*batch)[i].CompiledContentSignals.Downvotes)
		jVoteDelta := ((*batch)[j].CompiledContentSignals.Upvotes - (*batch)[j].CompiledContentSignals.Downvotes)
		// If they're not the same, use the vote delta sort.
		if iVoteDelta-jVoteDelta != 0 {
			return iVoteDelta > jVoteDelta
		}
		// If they are, make it so that the newer post is sorted higher
		return max((*batch)[i].Creation, (*batch)[i].LastUpdate) > max((*batch)[j].Creation, (*batch)[j].LastUpdate)
	})
}

func (batch *CPostBatch) ToProtobuf() []*feobjects.CompiledPostEntity {
	protos := []*feobjects.CompiledPostEntity{}
	for k, _ := range *batch {
		p := (*batch)[k].Protobuf()
		protos = append(protos, p)
	}
	return protos
}

type CompiledThread struct {
	Fingerprint            string `storm:"id"`
	Board                  string
	SelfCreated            bool
	Name                   string
	Body                   string
	Link                   string
	CompiledContentSignals CompiledContentSignals
	Owner                  CompiledUser
	Creation               int64
	LastUpdate             int64
	Meta                   string
	PostsCount             int
	Score                  float64
	ViewMeta_BoardName     string
}

func (c CompiledThread) BleveType() string {
	return "thread"
}

func (c CompiledThread) SearchId() string {
	// A path that defines a thread is board, threadfp. No user fp.
	bid, err := search.MakeSearchId("Thread", c.Board, "", "", c.Fingerprint, "")
	if err != nil {
		logging.Logf(1, "Making search ID failed. Error: %v", err)
		return ""
	}
	return bid
}

func (c *CompiledThread) IndexForSearch() {
	// search.Index(c.SearchId(), c)
	CThreadIndexCache = append(CThreadIndexCache, *c)
}

func (c *CompiledThread) DeleteFromSearchIndex() {
	search.Delete(c.SearchId())
}

func NewCThread(rp *pbstructs.Thread) CompiledThread {
	return CompiledThread{
		Fingerprint: rp.GetProvable().GetFingerprint(),
		Board:       rp.GetBoard(),
		SelfCreated: rp.GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
		Name:        rp.GetName(),
		Body:        rp.GetBody(),
		Link:        rp.GetLink(),
		Meta:        rp.GetMeta(),
		// Half-baked ones
		Owner: CompiledUser{
			Fingerprint: rp.GetOwner(),
		},
		Creation:   rp.GetProvable().GetCreation(),
		LastUpdate: rp.GetUpdateable().GetLastUpdate(),
	}
	// Needs: Compiledcontentsignals, owner, bymod, byop, blocked, approved flags
}

// Refresh refreshes an existing compiled thread's userheadercarrier and signals.
func (c *CompiledThread) RefreshContentSignals(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, nowts int64) {
	cs := CompiledContentSignals{}
	cs.Insert(c.Fingerprint, catds, cfgs, cmas, nowts)
	// Move values that need to be retained
	cs.SelfModIgnored = c.CompiledContentSignals.SelfModIgnored
	c.CompiledContentSignals = cs
}

// RefreshUserHeader needs board fingerprint because user signals within user headers are scoped, global scope is available to all and without a board fp, those don't depend on any board scope. But within the boards, there is a scope that is based on elections, assignments, people choosing to trust certain people as mods only within that board, and that scope needs to be applied over.
func (c *CompiledThread) RefreshUserHeader(boardSpecificUserHeaders CUserBatch) {
	for k, _ := range boardSpecificUserHeaders {
		if boardSpecificUserHeaders[k].Fingerprint == c.Owner.Fingerprint {
			c.Owner = boardSpecificUserHeaders[k]
			return
		}
	}
	// Not found in the local headers. Seek from global headers.
	uhc := UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in RefreshUserHeader>One")
	err := globals.KvInstance.One("Fingerprint", c.Owner.Fingerprint, &uhc)
	if err != nil {
		logging.Logf(2, "We could not find the user of this entity in local or global user headers scopes. Entity Fingerprint: %#v, User Fingerprint: %v", c.Fingerprint, c.Owner.Fingerprint)
		return
	}
	i := uhc.Users.Find(c.Owner.Fingerprint)
	if i != -1 {
		c.Owner = uhc.Users[i]
		return
	}
	logging.LogCrashf("This should never happen. We've found a user header carrier that matches the user fingerprint requested by this entity, but within the UHC, there was no CompiledUser matching the fingerprint. Entity: %#v, User Fingerprint: %v, UserHeaderCarrier: %#v", c, c.Owner.Fingerprint, uhc)
}

func (c *CompiledThread) Insert(ce CompiledThread) {
	if c.LastUpdate < ce.LastUpdate {
		*c = ce
		c.IndexForSearch()
	}
}

// CalcScore calculates the rank score for this thread.
func (c *CompiledThread) CalcScore() {
	// We need: upvotes, downvotes, current timestamp, creation
	voteScore := c.CompiledContentSignals.Upvotes - c.CompiledContentSignals.Downvotes
	orderOfMagnitude := math.Log10(math.Max(1, math.Abs(float64(voteScore))))
	sign := 0
	if voteScore > 0 {
		sign = 1
	}
	if voteScore < 0 {
		sign = -1
	}
	sec := c.Creation - 1533081600 // > Here we go again, Gordon Freeman
	score := (float64(sign) * orderOfMagnitude) + (float64(sec) / 42300)
	// > Approximate half life of Sodium-24
	c.Score = score
}

func (c *CompiledThread) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier) {
	c.RefreshUserHeader(boardSpecificUserHeaders)
	c.RefreshContentSignals(catds, cfgs, cmas, nowts)
	c.RefreshExogenousContentSignals(bc)
	c.CalcScore()
}

// RefreshExogenousContentSignals is where we compile and calculate the content signals that depend on external entitites.
func (c *CompiledThread) RefreshExogenousContentSignals(bc *BoardCarrier) {
	// ^ c: compiledthread, tc: threadcarrier. different things.
	/*
		The signals processed here are:
		ByMod            bool
		ByFollowedPerson bool
		ByBlockedPerson  bool
		ByOP             bool
		ModBlocked       bool
		ModApproved      bool

		creator's userheader: (ask to own userheader)
		- following / blocked state,
		- bymod state

		issuer's userheader (ask to board carrier)
		- modblock / modapprove state
	*/

	us := &c.Owner.CompiledUserSignals

	/*----------  ByMod state  ----------*/
	c.CompiledContentSignals.ByMod = isMod(us)

	/*----------  Following / blocked state  ----------*/
	c.CompiledContentSignals.ByFollowedPerson = us.FollowedBySelf
	c.CompiledContentSignals.ByBlockedPerson = us.BlockedBySelf

	/*----------  ByOp state  ----------*/
	/*----------  We don't do ByOp state in this one  ----------*/

	/*----------  Modblock / modapprove state  ----------*/
	// Behaviour: if at least one modblock, block it, if there is at least one modapprove, unblock it. so if something is both modblocked and modapproved, it will be visible.
	// Approvals
	for k, _ := range c.CompiledContentSignals.ModApprovals {
		sourcefp := c.CompiledContentSignals.ModApprovals[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModApproved = true
		}
	}
	// Blocks
	for k, _ := range c.CompiledContentSignals.ModBlocks {
		sourcefp := c.CompiledContentSignals.ModBlocks[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModBlocked = true
		}
	}
}

// Batch thread

type CThreadBatch []CompiledThread

// IndexForSearch adds all entities in this batch into the search index.
func (batch *CThreadBatch) IndexForSearch() {
	for k, _ := range *batch {
		CThreadIndexCache = append(CThreadIndexCache, (*batch)[k])
	}

	// if len(*batch) == 0 {
	// 	return
	// }
	// go func() {
	// 	ib := search.NewBatch()
	// 	for k, _ := range *batch {
	// 		ib.Index((*batch)[k].SearchId(), (*batch)[k])
	// 	}
	// 	err := search.CommitBatch(ib)
	// 	if err != nil {
	// 		logging.LogCrash(err)
	// 	}
	// 	logging.Logf(1, "Indexing is complete for this thread batch. Count: %v", len(*batch))
	// }()
}

func (batch *CThreadBatch) Insert(cthreads []CompiledThread) {
	toBeIndexed := CThreadBatch{}
	for k, _ := range cthreads {
		i := batch.Find(cthreads[k].Fingerprint)
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(cthreads[k])
			continue
		}
		// Doesn't exist in our batch. Add it.
		*batch = append(*batch, cthreads[k])
		toBeIndexed = append(toBeIndexed, cthreads[k])
	}
	// Index for new additions. (Index for edits happen inside singular insert.)
	toBeIndexed.IndexForSearch()
}

func (batch *CThreadBatch) InsertFromProtobuf(cthreads []*pbstructs.Thread) bool {
	var hasNewThreads bool
	toBeIndexed := CThreadBatch{}
	for k, _ := range cthreads {
		if cthreads[k] == nil {
			continue
		}
		i := batch.Find(cthreads[k].GetProvable().GetFingerprint())
		if i != -1 {
			// logging.Logf(1, "found in the compiled thread batch, updating")
			// logging.Logf(1, "This is what we were searching for: %v, and this is what we found: %v", cthreads[k], (*batch)[i])
			// TOFIX: why is the first thing even present? this is the first ever load. deal with that.
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(NewCThread(cthreads[k]))
			continue
		}
		// Doesn't exist in our batch. Add it.
		newthread := NewCThread(cthreads[k])
		if !newthread.SelfCreated {
			hasNewThreads = true
			/*
				^ We use this to determine whether we want to highlight the above scope board in the UI if it has new threads, but only if we did not create that thread ourselves.
			*/
		}
		*batch = append(*batch, newthread)
		toBeIndexed = append(toBeIndexed, newthread)
	}
	toBeIndexed.IndexForSearch()
	return hasNewThreads
}

func (batch *CThreadBatch) Find(threadfp string) int {
	for k, _ := range *batch {
		if threadfp == (*batch)[k].Fingerprint {
			return k
		}
	}
	return -1
}

func (batch *CThreadBatch) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier) {
	for k, _ := range *batch {
		(*batch)[k].Refresh(catds, cfgs, cmas, boardSpecificUserHeaders, nowts, bc)
	}
}

// Sort sorts the threads in the batch according to their score.
func (batch *CThreadBatch) SortByScore() {
	sort.Slice((*batch), func(i, j int) bool {
		return (*batch)[i].Score > (*batch)[j].Score
	})
}

func (batch *CThreadBatch) SortByCreation() {
	sort.Slice((*batch), func(i, j int) bool {
		return (*batch)[i].Creation > (*batch)[j].Creation
	})
}

func (batch *CThreadBatch) ToProtobuf() []*feobjects.CompiledThreadEntity {
	protos := []*feobjects.CompiledThreadEntity{}
	for k, _ := range *batch {
		p := (*batch)[k].Protobuf()
		protos = append(protos, p)
	}
	return protos
}

// The things missing from Key: expiry, info, meta
type CompiledUser struct {
	Fingerprint         string `storm:"id"`
	NonCanonicalName    string
	Info                string
	Expiry              int64
	Creation            int64
	LastUpdate          int64
	LastRefreshed       int64
	Meta                string
	CompiledUserSignals CompiledUserSignals
}

func (c CompiledUser) BleveType() string {
	return "user"
}

func (c CompiledUser) SearchId() string {
	// A path that defines a board is fingerprint. no board, thread, post, parent, userfp.
	bid, err := search.MakeSearchId("User", "", "", "", c.Fingerprint, "")
	if err != nil {
		logging.Logf(1, "Making Search ID failed. Error: %v", err)
		return ""
	}
	return bid
}

func (c *CompiledUser) IndexForSearch() {
	// search.Index(c.SearchId(), c)
	CUserIndexCache = append(CUserIndexCache, *c)
}

func (c *CompiledUser) DeleteFromSearchIndex() {
	search.Delete(c.SearchId())
}

func NewCUser(u *pbstructs.Key, nowts int64) CompiledUser {
	return CompiledUser{
		Fingerprint:      u.GetProvable().GetFingerprint(),
		NonCanonicalName: u.GetName(),
		Info:             u.GetInfo(),
		Expiry:           u.GetExpiry(),
		Creation:         u.GetProvable().GetCreation(),
		LastUpdate:       u.GetUpdateable().GetLastUpdate(),
		Meta:             u.GetMeta(),
		LastRefreshed:    nowts,
	}
	// needs: compiledusersignals
}

// Refresh refreshes an existing compiled thread's userheadercarrier and signals.
func (c *CompiledUser) RefreshUserSignals(
	cpts *CPTBatch, ccns *CCNBatch, cf451s *CF451Batch, cpes *CPEBatch, localDefaultMods []string, domainfp string, totalPop int) {
	cs := CompiledUserSignals{}
	cs.Insert(c.Fingerprint, domainfp, localDefaultMods, totalPop, cpts, ccns, cf451s, cpes)
	c.CompiledUserSignals = cs
}

// Insert is a full-on override - anything from the prior compiled user will be wiped out, including signals. If you want a soft merge where signals are merged, not replaced with the new signals, see InsertWithSignalMerge.
func (c *CompiledUser) Insert(ce CompiledUser) {
	if c.LastUpdate < ce.LastUpdate {
		*c = ce
		c.IndexForSearch()
	}
}

// InsertWithSignalMerge is useful when you want to merge a global user header with a community specific user header. It does a SUM type merge where signals are summed. (The normal merge just overwrites the older signals with the newer, it does not merge.)
func (c *CompiledUser) InsertWithSignalMerge(ce CompiledUser) {
	extantSignals := c.CompiledUserSignals
	oncomingSignals := ce.CompiledUserSignals
	c.Insert(ce) // this part will only tigger if c.lastupdate < ce.lastupdate, so regardless of which you merge into which, you'll end up with the most up to date result.
	c.CompiledUserSignals = extantSignals
	c.CompiledUserSignals.Merge(oncomingSignals)
}

func (c *CompiledUser) Refresh(
	cpts *CPTBatch, ccns *CCNBatch, cf451s *CF451Batch, cpes *CPEBatch, localDefaultMods []string, domainfp string, totalPop int) {
	c.RefreshUserSignals(cpts, ccns, cf451s, cpes, localDefaultMods, domainfp, totalPop)
}

/*----------  Utility methods for CompiledUser  ----------*/

type CUserUsername struct {
	SourceCUser string
	Username    string
	Canonical   bool
}

func (c *CompiledUser) GetUsername() CUserUsername {
	un := CUserUsername{}
	un.SourceCUser = c.Fingerprint
	un.Username = c.NonCanonicalName
	if len(c.CompiledUserSignals.CanonicalName) > 0 {
		un.Username = c.CompiledUserSignals.CanonicalName
		un.Canonical = true
	}
	return un
}

/*----------  END Utility methods for CompiledUser  ----------*/

// Batch user

type CUserBatch []CompiledUser

// IndexForSearch adds all entities in this batch into the search index.
func (batch *CUserBatch) IndexForSearch() {
	for k, _ := range *batch {
		CUserIndexCache = append(CUserIndexCache, (*batch)[k])
	}
	// if len(*batch) == 0 {
	// 	return
	// }
	// go func() {
	// 	ib := search.NewBatch()
	// 	for k, _ := range *batch {
	// 		ib.Index((*batch)[k].SearchId(), (*batch)[k])
	// 	}
	// 	err := search.CommitBatch(ib)
	// 	if err != nil {
	// 		logging.LogCrash(err)
	// 	}
	// 	logging.Logf(1, "Indexing is complete for this user batch. Count: %v", len(*batch))
	// }()
}

func (batch *CUserBatch) Insert(cusers []CompiledUser) {
	toBeIndexed := CUserBatch{}
	for k, _ := range cusers {
		i := batch.Find(cusers[k].Fingerprint)
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(cusers[k])
			continue
		}
		// Doesn't exist in our batch. Add it.
		*batch = append(*batch, cusers[k])
		toBeIndexed = append(toBeIndexed, cusers[k])
	}
	// Index for new additions. (Index for edits happen inside singular insert.)
	toBeIndexed.IndexForSearch()
}

func (batch *CUserBatch) InsertWithSignalMerge(cusers []CompiledUser) {
	for k, _ := range cusers {
		i := batch.Find(cusers[k].Fingerprint)
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].InsertWithSignalMerge(cusers[k])
			continue
		}
		// Doesn't exist in our batch. Add it.
		*batch = append(*batch, cusers[k])
	}
	// Heads up: this does not index anything for search, since signals are not indexed, so this action cannot create a change in the index.
}

func (batch *CUserBatch) InsertFromProtobuf(cusers []*pbstructs.Key, nowts int64) {
	toBeIndexed := CUserBatch{}
	for k, _ := range cusers {
		if cusers[k] == nil {
			continue
		}
		i := batch.Find(cusers[k].GetProvable().GetFingerprint())
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(NewCUser(cusers[k], nowts))
			continue
		}
		// Doesn't exist in our batch. Add it.
		newuser := NewCUser(cusers[k], nowts)
		*batch = append(*batch, newuser)
		toBeIndexed = append(toBeIndexed, newuser)
	}
	toBeIndexed.IndexForSearch()
}

func (batch *CUserBatch) Find(userfp string) int {
	for k, _ := range *batch {
		if userfp == (*batch)[k].Fingerprint {
			return k
		}
	}
	return -1
}

func (batch *CUserBatch) Refresh(cpts *CPTBatch, ccns *CCNBatch, cf451s *CF451Batch, cpes *CPEBatch, localDefaultMods []string, domainfp string, totalPop int) {
	for k, _ := range *batch {
		(*batch)[k].Refresh(cpts, ccns, cf451s, cpes, localDefaultMods, domainfp, totalPop)
	}
}

type CompiledBoard struct {
	Fingerprint            string `storm:"id"`
	SelfCreated            bool
	Name                   string
	Description            string
	CompiledContentSignals CompiledContentSignals
	Owner                  CompiledUser
	BoardOwners            []string
	Creation               int64
	LastUpdate             int64
	Meta                   string
	LocalScopeUserHeaders  CUserBatch
	// ^ This carries entitlements specific to this specific board for users. If a person is a mod of this board, this is where his or her modship flag gets stored in.
	ThreadsCount int
	UserCount    int
}

func (c CompiledBoard) BleveType() string {
	return "board"
}

func (c CompiledBoard) SearchId() string {
	// A path that defines a board is fingerprint. no board, thread, post, parent, userfp.
	bid, err := search.MakeSearchId("Board", "", "", "", c.Fingerprint, "")
	if err != nil {
		logging.Logf(1, "Making search ID failed. Error: %v", err)
		return ""
	}
	return bid
}

func (c *CompiledBoard) IndexForSearch() {
	// search.Index(c.SearchId(), c)
	CBoardIndexCache = append(CBoardIndexCache, *c)
}

func (c *CompiledBoard) DeleteFromSearchIndex() {
	search.Delete(c.SearchId())
}

func NewCBoard(rp *pbstructs.Board) CompiledBoard {
	cb := CompiledBoard{
		Fingerprint: rp.GetProvable().GetFingerprint(),
		SelfCreated: rp.GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
		Name:        rp.GetName(),
		Description: rp.GetDescription(),
		// Half-baked ones
		Owner: CompiledUser{
			Fingerprint: rp.GetOwner(),
		},
		Creation:   rp.GetProvable().GetCreation(),
		LastUpdate: rp.GetUpdateable().GetLastUpdate(),
		Meta:       rp.GetMeta(),
	}
	if bo := rp.GetBoardOwners(); len(bo) > 0 {
		for k, _ := range bo {
			cb.BoardOwners = append(cb.BoardOwners, bo[k].GetKeyFingerprint())
		}
	}
	return cb
	// Needs: Compiledcontentsignals, owner, bymod, byop, blocked, approved flags
}

// GetUserHeader attempts to get the local user header if available within that board scope, if not, it attempts to get the global user header, if present.
func (cb *CompiledBoard) GetUserHeader(fp string) CompiledUser {
	for k, _ := range cb.LocalScopeUserHeaders {
		if cb.LocalScopeUserHeaders[k].Fingerprint == fp {
			return cb.LocalScopeUserHeaders[k]
		}
	}
	uhc := UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in GetUserHeader>One")
	err := globals.KvInstance.One("Fingerprint", fp, &uhc)
	if err != nil {
		logging.Logf(1, "We could not get the requested user from the global user headers. Error: %v, We asked for: %v", err, fp)
		return CompiledUser{}
	}
	if len(uhc.Users) > 0 {
		return uhc.Users[0]
	}
	return CompiledUser{}
}

func (cb *CompiledBoard) GetDefaultMods() []string {
	var dm []string
	dm = append(dm, cb.Owner.Fingerprint)
	dm = append(dm, cb.BoardOwners...)
	// To map and back again to remove dedupes.
	m := make(map[string]bool)
	for k, _ := range dm {
		m[dm[k]] = true
	}
	result := []string{}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// RefreshContentSignals refreshes an existing compiled board's userheadercarrier and signals.
func (c *CompiledBoard) RefreshContentSignals(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, nowts int64) {
	cs := CompiledContentSignals{}
	cs.Insert(c.Fingerprint, catds, cfgs, cmas, nowts)
	// Move values that need to be retained
	cs.SelfModIgnored = c.CompiledContentSignals.SelfModIgnored
	c.CompiledContentSignals = cs
}

func (c *CompiledBoard) RefreshUserHeader(boardSpecificUserHeaders CUserBatch) {
	for k, _ := range boardSpecificUserHeaders {
		if boardSpecificUserHeaders[k].Fingerprint == c.Owner.Fingerprint {
			c.Owner = boardSpecificUserHeaders[k]
			return
		}
	}
	// Not found in the local headers. Seek from global headers.
	uhc := UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in RefreshUserHeader>One")
	err := globals.KvInstance.One("Fingerprint", c.Owner.Fingerprint, &uhc)
	if err != nil {
		logging.Logf(2, "We could not find the user of this entity in local or global user headers scopes. Entity Fingerprint: %#v, User Fingerprint: %v", c.Fingerprint, c.Owner.Fingerprint)
		return
	}
	i := uhc.Users.Find(c.Owner.Fingerprint)
	if i != -1 {
		c.Owner = uhc.Users[i]
		return
	}
	logging.LogCrashf("This should never happen. We've found a user header carrier that matches the user fingerprint requested by this entity, but within the UHC, there was no CompiledUser matching the fingerprint. Entity: %#v, User Fingerprint: %v, UserHeaderCarrier: %#v", c, c.Owner.Fingerprint, uhc)
}

func (c *CompiledBoard) Insert(ce CompiledBoard) {
	if c.LastUpdate < ce.LastUpdate {
		*c = ce
		c.IndexForSearch()
	}
}

func (c *CompiledBoard) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier) {
	c.RefreshUserHeader(boardSpecificUserHeaders)
	c.RefreshContentSignals(catds, cfgs, cmas, nowts)
	c.RefreshExogenousContentSignals(bc)
}

// RefreshExogenousContentSignals is where we compile and calculate the content signals that depend on external entitites.
func (c *CompiledBoard) RefreshExogenousContentSignals(bc *BoardCarrier) {
	// ^ c: compiledthread, tc: threadcarrier. different things.
	/*
		The signals processed here are:
		ByMod            bool
		ByFollowedPerson bool
		ByBlockedPerson  bool
		ByOP             bool
		ModBlocked       bool
		ModApproved      bool

		creator's userheader: (ask to own userheader)
		- following / blocked state,
		- bymod state

		issuer's userheader (ask to board carrier)
		- modblock / modapprove state
	*/

	us := &c.Owner.CompiledUserSignals

	/*----------  ByMod state  ----------*/
	c.CompiledContentSignals.ByMod = isMod(us)
	// ^ In a board, this will always be default true

	/*----------  Following / blocked state  ----------*/
	c.CompiledContentSignals.ByFollowedPerson = us.FollowedBySelf
	c.CompiledContentSignals.ByBlockedPerson = us.BlockedBySelf

	/*----------  ByOp state  ----------*/
	/*----------  We don't do ByOp state in this one  ----------*/

	/*----------  Modblock / modapprove state  ----------*/
	// Behaviour: if at least one modblock, block it, if there is at least one modapprove, unblock it. so if something is both modblocked and modapproved, it will be visible.
	// Approvals
	for k, _ := range c.CompiledContentSignals.ModApprovals {
		sourcefp := c.CompiledContentSignals.ModApprovals[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModApproved = true
		}
	}
	// Blocks
	for k, _ := range c.CompiledContentSignals.ModBlocks {
		sourcefp := c.CompiledContentSignals.ModBlocks[k].SourceFp
		b := &CompiledBoard{}
		for k, _ := range bc.Boards {
			if bc.Boards[k].Fingerprint == bc.Fingerprint {
				b = &bc.Boards[k]
			}
		}
		uh := b.GetUserHeader(sourcefp)
		if isMod(&uh.CompiledUserSignals) {
			c.CompiledContentSignals.ModBlocked = true
		}
	}
}

type CBoardBatch []CompiledBoard

// IndexForSearch adds all entities in this batch into the search index.
func (batch *CBoardBatch) IndexForSearch() {
	for k, _ := range *batch {
		CBoardIndexCache = append(CBoardIndexCache, (*batch)[k])
	}
	// if len(*batch) == 0 {
	// 	return
	// }
	// go func() {
	// 	ib := search.NewBatch()
	// 	for k, _ := range *batch {
	// 		ib.Index((*batch)[k].SearchId(), (*batch)[k])
	// 	}
	// 	err := search.CommitBatch(ib)
	// 	if err != nil {
	// 		logging.LogCrash(err)
	// 	}
	// 	logging.Logf(1, "Indexing is complete for this board batch. Count: %v", len(*batch))
	// }()
}

func (batch *CBoardBatch) Insert(cboards []CompiledBoard) {
	toBeIndexed := CBoardBatch{}
	for k, _ := range cboards {
		i := batch.Find(cboards[k].Fingerprint)
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(cboards[k])
			continue
		}
		// Doesn't exist in our batch. Add it.
		*batch = append(*batch, cboards[k])
		toBeIndexed = append(toBeIndexed, cboards[k])
	}
	// Index for new additions. (Index for edits happen inside singular insert.)
	toBeIndexed.IndexForSearch()
}

func (batch *CBoardBatch) InsertFromProtobuf(cboards []*pbstructs.Board) {
	toBeIndexed := CBoardBatch{}
	for k, _ := range cboards {
		if cboards[k] == nil {
			continue
		}
		i := batch.Find(cboards[k].GetProvable().GetFingerprint())
		if i != -1 {
			// Trigger insert. It'll only update if the lastupdate is newer.
			(*batch)[i].Insert(NewCBoard(cboards[k]))
			continue
		}
		// Doesn't exist in our batch. Add it.
		newboard := NewCBoard(cboards[k])
		*batch = append(*batch, newboard)
		toBeIndexed = append(toBeIndexed, newboard)
	}
	toBeIndexed.IndexForSearch()
}

func (batch *CBoardBatch) Find(boardfp string) int {
	for k, _ := range *batch {
		if boardfp == (*batch)[k].Fingerprint {
			return k
		}
	}
	return -1
}

func (batch *CBoardBatch) Refresh(catds *CATDBatch, cfgs *CFGBatch, cmas *CMABatch, boardSpecificUserHeaders CUserBatch, nowts int64, bc *BoardCarrier) {
	for k, _ := range *batch {
		(*batch)[k].Refresh(catds, cfgs, cmas, boardSpecificUserHeaders, nowts, bc)
	}
}

func (batch *CBoardBatch) GetDefaultMods() []string {
	var dms []string
	for k, _ := range *batch {
		dms = append(dms, (*batch)[k].GetDefaultMods()...)
	}
	return dms
}

func (batch *CBoardBatch) GetBoardSpecificUserHeaders() CUserBatch {
	b := CUserBatch{}
	for k, _ := range *batch {
		b = append(b, (*batch)[k].LocalScopeUserHeaders...)
	}
	return b
}

func (batch *CBoardBatch) SortByThreadsCount() {
	sort.Slice((*batch), func(i, j int) bool {
		return (*batch)[i].ThreadsCount > (*batch)[j].ThreadsCount
	})
}

////////////////////
// Compiled Signals
///////////////////

// Compiled User Signals

// Self refers to the local user. SelfFollowed = followed by the local user.

type CompiledUserSignals struct {
	TargetFingerprint string
	Domain            string
	// These are the compiled 'final' decisions based on the collections of signals we have.
	FollowedBySelf         bool // check private follows list too
	BlockedBySelf          bool // check private blocks list too
	FollowerCount          int  // This is on an ongoing basis, within network memory. In practice, this means new follower count in the last 6 months.
	CanonicalName          string
	CNameSourceFingerprint string
	// Self, PE data for domain given above
	SelfPEFingerprint string
	SelfPECreation    int64
	SelfPELastUpdate  int64
	// Mod signals
	MadeModBySelf bool
	/*
		^ Any mod that is accepted as one because of the local user input. This trumps everything else, if a local user does this, it doesn't matter what the network thinks.

		- Global scope

			- A user that was publicelect:modelect'ed by the local user without a domain (without a specific board given in domain), which means it's valid for all boards
			- A user that was privately modelected by the local user without a domain (without releasing a publicelect:modelect signal).
			- A censor that was assigned by a CA that the local user trusts.

		- Local scope

			 - A user that was publicelect:modelect'ed by the local user with a domain, which means the modship is valid for that domain only.
			 - A user that was privately modelected by the local user, with a domain.

		Heads up:

		- When the local user releases a publictrust:modelect or publictrust:moddisqualify, those count as public votes on whether the local user thinks that person should be a mod or that they should be impeached, and it will count based on the scope that you put in the domain. So when you vote for a person either way, your decision will immediately start to apply for you. For example, if you vote for a person to elect him/her as a mod for a specific board, that person immediately becomes a mod for you for that board. In other words, if you want to get somebody elected a mod, you have to live with him/her being a mod for you also.

	*/
	MadeNonModBySelf bool
	/*
		^ Any mod that was disqualified from being a mod because of local user input.

		- Global scope

			- The local user did a publicelect:moddisqualify on a user without a domain.
			- The local user did a private moddisqualify on the user without a domain.
			- The local user did a private moddisqualify on the censor.

		- Local scope

			- The local user did a publicelect:moddisqualify on a user with a domain, which disqualified that user from being a mod within that board for the local user.
			- The local user did a private disqualify on that user with a domain.

		Heads up:

		- By default, the elect/disqualify actions are both public and private for all users except censors. For elect/disqualify on censors, the action is by default private (but can be made public also.) The reason why is that the decision to ignore censorship should be a personal one, and in some contexts, the local user might not want to broadcast that s/he is ignoring censorship, especially if the censor does have some real-world authority over them. The network is anonymous, so even if these actions were broadcast it would be nigh impossible to link a specific user key to a specific real world person for anybody. That said, out of an abundance of caution and with the appreciation of the fact that nothing is impossible given enough money and determination, publictrust actions toward censors are made private by default.

		- Small detail: These actions will only be private if you have a truststate from a CA that you trust that say that specific user is a censor. So if that 'this key is a censor' message hasn't arrived from a CA you trust to you yet, or that it is an appointed censor from a CA that you do not trust (in which case you won't have that 'this key is a censor' message because your backend will immediately decline all CA-specific messages originating from that CA), your publictrust actions for that key will be public by default. Ultimately, these actions are private only if that key is an otherwise valid censor for you.

	*/
	MadeModByDefault bool
	/*
		^ Any mod that is accepted as by-default mod. This is set under these conditions:

		- Global scope (all boards)

			- None. There are no mods in the global scope by default.

		- Local scope (in a board)

			- The user is an admin (founder) of a board, or a person that is assigned a mod by the admin of a board. An admin can change board entity (name, description, etc) and moderate content. A mod can moderate content. Both mods and admins are subject to impeachment of their mod privileges, but admins retain the ability to edit the mods list, description and meta fields even when impeached.

			(If an admin actually gets impeached, though, it might be time for that community to just move on and found a new board.)

	*/
	// MadeNonModByDefault bool // Doesn't exist. Everybody has a chance of becoming a mod for any scope.
	MadeModByNetwork bool
	/*
		^ Any mod that is accepted as mod because of the network input. This can be because:

		- Global scope

			- The network has elected a mod based on publictrust:modelects.

		- Local scope

			- The participants of that board has elected a mod based on publictrust:modelects.

	*/
	MadeNonModByNetwork bool
	/*
		^ Any mod that is made not one because of the network decision. This can be because:

		- Global scope

			- The network has disqualified a user from being a mod because of the publictrust:moddisqualify votes.

		- Local scope

			- Participants of that board disqualified a user from being a mod in that board by casting publictrust:moddisqualify votes. This is impeachment.
	*/

	/*
		These modsignals are applied in this order, highest priority being at the top.

		- Self signals. As a guiding principle, nothing coming over from the network can override local user decisions. Local user is fully sovereign on what his/her device does with the data received from the network, and neither network nor defaults can possibly override this.

		- Network signals. These signals represent the decisions of the scope in question, whether be it global (the whole network has voted a certain way) or local (the participants of the board has voted a certain way).

		- Defaults. These assumptions represent some implicit trust that is there by default unless explicitly overridden by the network or by the user. For example, the founder of a board (admin) is a default mod in that board, as well as mods that s/he assigns. That said, if the local user decides that they're not up to the task, or the network decides as such, they can be impeached.
	*/

	/*
		As a final note, all of these network votes are based on public trustsate entities, and they expire around 6 months. The voting process is an ongoing one, and it never stops. Effectively, this means when you vote for someone to be elected a mod for a specific board, you will cast a public vote (unless you've disabled public voting) and save a private one. The public vote will live on for six months, visible to everyone. The private one will be permanent, and it'll be saved in your local profile, so even when your public vote expires, that person will remain a mod for you. For him to remain a network-elected mod, though, he needs to stay above 50% of votes being positive on an ongoing basis, since the voting is always an ongoing process, and when your vote expires, it'll be dropped from the count.

		Your vote can also can cease to be counted when you stop contributing to that board within a 2-week timeframe. So for a vote to be counted, you need to be, and remain, an active participant of the network, and that your vote needs to have not expired.

		Network makes decisions at 50% vote threshold and at 5% population threshold. That means the result of a vote is only valid if at least 5% of the population of that scope has cast a public vote one way or the other.

		For example, if you are an elected mod, you need to stay above 50% approval, and 5% of the active participants of that board should have voted. If for any reason either of these conditions cease to be true, you are reverted back to your non-mod state.

		If you are a default-mod, you are impeached if 50+% votes for your impeachment. This vote is valid only if at least 5% of the population of that board has voted. If either of those conditions cease to be true, you are reverted back to your default mod state.

		In essence, it takes ongoing conviction by people who have a stake in that scope to keep a network decision in effect. If that conviction evaporates one way or the other, the decision ceases to be in effect. The intent is to not make the children suffer the consequences of poor decisions of their forefathers.
	*/
}

func (s *CompiledUserSignals) Merge(sn CompiledUserSignals) {
	if s.TargetFingerprint != sn.TargetFingerprint {
		return // Not the signals for the same entity! Can't merge.
	}
	s.Domain = sn.Domain
	// ^ Mind that we're not checking the domain to be the same. You *can* merge a local and a global scope signal, but the result will always be local.
	s.FollowedBySelf = s.FollowedBySelf || sn.FollowedBySelf
	s.BlockedBySelf = s.BlockedBySelf || sn.BlockedBySelf
	s.FollowerCount = s.FollowerCount + sn.FollowerCount
	if len(sn.CanonicalName) > 0 && ca.IsTrustedCAKeyByFp(sn.CNameSourceFingerprint) {
		s.CanonicalName = sn.CanonicalName
		// ^ If there's a domain-specific canonical name, we apply that since it takes priority. This is for the future where we might actually have that. As of this code being written (June 2018), we don't have that feature.
		s.CNameSourceFingerprint = sn.CNameSourceFingerprint
	}
	s.MadeModBySelf = s.MadeModBySelf || sn.MadeModBySelf
	s.MadeNonModBySelf = s.MadeNonModBySelf || sn.MadeNonModBySelf
	s.MadeModByDefault = s.MadeModByDefault || sn.MadeModByDefault
	s.MadeModByNetwork = s.MadeModByNetwork || sn.MadeModByNetwork
	s.MadeNonModByNetwork = s.MadeNonModByNetwork || sn.MadeNonModByNetwork

	// Merge data (the input to be merged overrides the existing input)
	s.SelfPEFingerprint = sn.SelfPEFingerprint
	s.SelfPECreation = sn.SelfPECreation
	s.SelfPELastUpdate = sn.SelfPELastUpdate
}

func (s *CompiledUserSignals) Insert(
	targetfp, domainfp string,
	localDefaultMods []string,
	totalPop int,
	cpts *CPTBatch,
	ccns *CCNBatch,
	cf451s *CF451Batch,
	cpes *CPEBatch) {
	if len(targetfp) == 0 || cpts == nil || ccns == nil || cf451s == nil || cpes == nil {
		return
	}
	s.TargetFingerprint = targetfp
	// Look at the feconfig to determine whether the user is followed by self in this current scope.
	cpt := cpts.FindObj(targetfp)
	ccn := ccns.FindObj(targetfp)
	cf451 := cf451s.FindObj(targetfp)
	cpe := cpes.FindObj(targetfp)
	s.FollowedBySelf = globals.FrontendConfig.UserRelations.Following.Find(targetfp, domainfp) != -1
	s.BlockedBySelf = globals.FrontendConfig.UserRelations.Blocked.Find(targetfp, domainfp) != -1
	s.FollowerCount = parseFollowerCount(targetfp, cpt)
	s.MadeModBySelf = globals.FrontendConfig.UserRelations.ModElected.Find(targetfp, domainfp) != -1 && isF451Mod(targetfp, cf451)
	s.MadeNonModBySelf = globals.FrontendConfig.UserRelations.ModDisqualified.Find(targetfp, domainfp) != -1
	s.CanonicalName, s.CNameSourceFingerprint = parseCanonicalName(ccn)
	s.MadeModByDefault = isModByDefault(targetfp, localDefaultMods)
	s.MadeModByNetwork, s.MadeNonModByNetwork = parsePublicElectByNetwork(targetfp, totalPop, cpe)
	s.SelfPEFingerprint = cpe.SelfFingerprint
	s.SelfPECreation = cpe.SelfCreation
	s.SelfPELastUpdate = cpe.SelfLastUpdate
}

func parsePublicElectByNetwork(targetfp string, totalPop int, cpe CompiledPE) (elected, disqualified bool) {
	totalVoteCount := cpe.ElectsCount + cpe.DisqualifiesCount
	totalVoteRequired := int(float64(totalPop) * (float64(globals.FrontendConfig.GetThresholdForElectionValidityPercent()) / 100))
	// ^ This is a frustrating way to do totalpop * 0.05.
	if totalVoteCount < totalVoteRequired {
		return false, false
		// The vote is invalid because not enough people voted.
	}
	totalElectsRequired := int(float64(totalVoteCount) * (float64(globals.FrontendConfig.GetThresholdForElectionWinPercent()) / 100))
	totalDisqualifiesRequired := int(float64(totalVoteCount) * (float64(globals.FrontendConfig.GetThresholdForElectionWinPercent()) / 100))
	if cpe.ElectsCount > cpe.DisqualifiesCount {
		// The vote is valid and more elects than disqualifies
		if cpe.ElectsCount < globals.FrontendConfig.GetMinimumVoteThresholdForElectionValidity() {
			return false, false
			// The are just way too few votes for this election to be valid. Below this threshold, elections start to become erratic and easy to manipulate.
		}
		if cpe.ElectsCount < totalElectsRequired {
			return false, false
			// the vote is valid, has more elects than disqualifies, but it hasn't crossed the win threshold.
			// (This sounds unnecessary on a two way vote, but it is possible when you have a 3 way vote and a 50% threshold. They can all remain at 33% and none of them would win.)
		}
		return true, false
		// the vote is valid, has more elects than disqualifies, and it has crossed the win threshold
	}
	// The vote is valid and more disqualifies than elects
	if cpe.DisqualifiesCount < globals.FrontendConfig.GetMinimumVoteThresholdForElectionValidity() {
		return false, false
		// The are just way too few votes for this election to be valid. Below this threshold, elections start to become erratic and easy to manipulate.
	}
	if cpe.DisqualifiesCount < totalDisqualifiesRequired {
		return false, false
		// the vote is valid, has more disqualifies than elects, but it hasn't crossed the win threshold.
	}
	return false, true
	// the vote is valid, has more disqualifies than elects, and it has crossed the win threshold.
}

func parseFollowerCount(targetfp string, cpt CompiledPT) int {
	var count int
	for k, _ := range cpt.PTs {
		if cpt.PTs[k].Fingerprint == targetfp &&
			cpt.PTs[k].Type == Signal_Follow {
			count++
		}
	}
	return count
}

func isF451Mod(targetfp string, cf451 CompiledF451) bool {
	for k, _ := range cf451.F451s {
		if cf451.F451s[k].TargetFingerprint == targetfp && ca.IsTrustedCAKeyByFp(cf451.F451s[k].SourceFingerprint) {
			return true
		}
	}
	return false
}

func isModByDefault(targetfp string, localDefaultMods []string) bool {
	for k, _ := range localDefaultMods {
		if localDefaultMods[k] == targetfp {
			return true
		}
	}
	return false
}

func parseCanonicalName(ccn CompiledCN) (cname, cnamesource string) {
	highestPrioritySourceKey := -1
	highestPrioritySoFar := 0
	highestPrioritySet := false
	for k, _ := range ccn.CNs {
		isCaKey, priority := ca.IsTrustedCAKeyByFpWithPriority(
			ccn.CNs[k].SourceFingerprint)
		if isCaKey {
			if !highestPrioritySet {
				highestPrioritySet = true
				highestPrioritySoFar = priority
				highestPrioritySourceKey = k
				continue
			}
			if highestPrioritySoFar > priority { // higher number means lower prio.
				highestPrioritySoFar = priority
				highestPrioritySourceKey = k
				continue
			}
		}
	}
	if !highestPrioritySet {
		// None of these keys were from a CA we trusted.
		return "", ""
	}
	return ccn.CNs[highestPrioritySourceKey].CanonicalName, ccn.CNs[highestPrioritySourceKey].SourceFingerprint
}

// Compiled content signals

// CompiledContentSignals is the final compiled form of all signals that relate to a specific entity, baked directly into that content entity.
type CompiledContentSignals struct {
	TargetFingerprint string
	/*----------  Endogenous signals  ----------*/
	// (Signals that are directly generated from Vote entities pointing to the target)
	// ATD
	Upvotes            int
	Downvotes          int
	SelfUpvoted        bool
	SelfDownvoted      bool
	SelfATDFingerprint string
	SelfATDCreation    int64
	SelfATDLastUpdate  int64
	// ^ In aggregate types such as ATDs, we have to carry over the creation, lastupdate and fingerprint to the client, because the client needs those information to be able to edit the signal. In other types we carry those information in the signal entity itself, since they are not aggregated, the client can figure out what to do based on determining which one is self.
	// FG
	Reports      []ExplainedSignal
	SelfReported bool
	// MA
	ModBlocks       []ExplainedSignal
	ModApprovals    []ExplainedSignal
	SelfModBlocked  bool
	SelfModApproved bool
	SelfModIgnored  bool // Only available as self, is not communicated out.

	/*----------  Exogenous Signals  ----------*/
	// (Signals generated after combining the above with external entities - second generation signals)
	ByMod bool
	// ^ Because receiving a mod signal means nothing without knowing we trust that mod or not, which comes from the user entity's signals.
	ByFollowedPerson bool
	ByBlockedPerson  bool
	ByOP             bool
	ModBlocked       bool
	ModApproved      bool

	LastRefreshed int64
}

func (s *CompiledContentSignals) Insert(
	targetfp string,
	catds *CATDBatch,
	cfgs *CFGBatch,
	cmas *CMABatch,
	nowts int64,
) {
	if len(targetfp) == 0 || catds == nil || cfgs == nil || cmas == nil {
		return
	}
	s.TargetFingerprint = targetfp
	// compile catd related to this entity
	i := catds.Find(targetfp)
	if i != -1 {
		catd := (*catds)[i]
		s.Upvotes = catd.UpvotesCount
		s.Downvotes = catd.DownvotesCount
		if catd.SelfVoted {
			if catd.SelfVoteDirection == Signal_Upvote {
				s.SelfUpvoted = true
			}
			if catd.SelfVoteDirection == Signal_Downvote {
				s.SelfDownvoted = true
			}
			s.SelfATDCreation = catd.SelfCreation
			s.SelfATDFingerprint = catd.SelfFingerprint
			s.SelfATDLastUpdate = catd.SelfLastUpdate
		}
		if s.LastRefreshed < catd.LastRefreshed {
			s.LastRefreshed = catd.LastRefreshed
		}
	}
	// compile fgs: reports
	i2 := cfgs.Find(targetfp)
	if i2 != -1 {
		cfg := (*cfgs)[i2]
		expss := []ExplainedSignal{}
		for k, _ := range cfg.FGs {
			expss = append(expss, cfg.FGs[k].CnvToExplainedSignal())
			// if cfg.FGs[k].Self {
			// 	s.SelfReported = true
			// }
		}
		s.Reports = expss
		s.SelfReported = cfg.SelfReported
	}
	// compile mas
	i3 := cmas.Find(targetfp)
	if i3 != -1 {
		cma := (*cmas)[i3]
		for k, _ := range cma.MAs {
			if cma.MAs[k].Type == Signal_ModBlock {
				s.ModBlocks = append(
					s.ModBlocks, cma.MAs[k].CnvToExplainedSignal())
			}
			if cma.MAs[k].Type == Signal_ModApprove {
				s.ModApprovals = append(
					s.ModApprovals, cma.MAs[k].CnvToExplainedSignal())
			}
		}
		s.SelfModBlocked = cma.SelfModBlocked
		s.SelfModApproved = cma.SelfModApproved
	}
	s.LastRefreshed = nowts
}

type ExplainedSignal struct {
	SourceFp   string
	Reason     string
	Creation   int64
	LastUpdate int64
}

/////////////////////////
// Constant types
/////////////////////////

const (
	// // Vote signal types
	// AddsToDiscussion types
	Signal_Upvote   = 1
	Signal_Downvote = 2
	// FollowsGuidelines types
	Signal_ReportToMod = 1
	// ModActions types
	Signal_ModBlock   = 1
	Signal_ModApprove = 2
)

const (
	// // Truststates signal types
	// PublicTrust types
	Signal_Follow = 1
	Signal_Block  = 2
	// CanonicalName types
	Signal_NameAssign = 1
	// F451 types
	Signal_CensorAssign = 1
	// PublicElect types
	Signal_Elect      = 1
	Signal_Disqualify = 2
)

// InitialiseKvStore initialises the buckets that we have in the store. This runs only once every run.
func InitialiseKvStore() {
	// logging.Logf(1, "KvStore buckets are being created.")
	err1 := globals.KvInstance.Init(&BoardCarrier{})
	if err1 != nil {
		logging.Logf(1, "BoardCarrier init encountered a problem. Error: %v", err1)
	}
	err2 := globals.KvInstance.Init(&ThreadCarrier{})
	if err2 != nil {
		logging.Logf(1, "ThreadCarrier init encountered a problem. Error: %v", err2)
	}
	err3 := globals.KvInstance.Init(&AmbientBoard{})
	if err3 != nil {
		logging.Logf(1, "AmbientBoard init encountered a problem. Error: %v", err3)
	}
	err4 := globals.KvInstance.Init(&GlobalStatisticsCarrier{})
	if err4 != nil {
		logging.Logf(1, "GlobalStatisticsCarrier init encountered a problem. Error: %v", err4)
	}
	err5 := globals.KvInstance.Init(&UserHeaderCarrier{})
	if err5 != nil {
		logging.Logf(1, "UserHeaderCarrier init encountered a problem. Error: %v", err5)
	}
}

/*----------  Reports tab entry  ----------*/

/*
	This is generated after the fact, after the core entities are compiled. It collects all the items with reports and puts them into a sortable payload form.

*/
type ReportsTabEntry struct {
	Fingerprint   string
	BoardPayload  CompiledBoard
	ThreadPayload CompiledThread
	PostPayload   CompiledPost
	Timestamp     int64
}

func NewReportsTabEntryFromBoard(e *CompiledBoard) *ReportsTabEntry {
	return &ReportsTabEntry{
		Fingerprint:  e.Fingerprint,
		BoardPayload: *e,
		Timestamp:    getNewestReportTimestamp(&e.CompiledContentSignals),
	}
}

func NewReportsTabEntryFromThread(e *CompiledThread) *ReportsTabEntry {
	return &ReportsTabEntry{
		Fingerprint:   e.Fingerprint,
		ThreadPayload: *e,
		Timestamp:     getNewestReportTimestamp(&e.CompiledContentSignals),
	}
}

func NewReportsTabEntryFromPost(e *CompiledPost) *ReportsTabEntry {
	return &ReportsTabEntry{
		Fingerprint: e.Fingerprint,
		PostPayload: *e,
		Timestamp:   getNewestReportTimestamp(&e.CompiledContentSignals),
	}
}

func getNewestReportTimestamp(c *CompiledContentSignals) int64 {
	var newest int64
	for k, _ := range c.Reports {
		if stamp := max(c.Reports[k].Creation, c.Reports[k].LastUpdate); stamp > newest {
			newest = stamp
		}
	}
	return newest
}

type ReportsTabEntryBatch []ReportsTabEntry

/*=============================================
=            Mod actions tab entry            =
=============================================*/

/*
	This is generated after the fact, after the core entities are compiled. It collects all the items with reports and puts them into a sortable payload form.

*/
type ModActionsTabEntry struct {
	Fingerprint   string
	BoardPayload  CompiledBoard
	ThreadPayload CompiledThread
	PostPayload   CompiledPost
	Timestamp     int64
}

func NewModActionsTabEntryFromBoard(e *CompiledBoard) *ModActionsTabEntry {
	return &ModActionsTabEntry{
		Fingerprint:  e.Fingerprint,
		BoardPayload: *e,
		Timestamp:    getNewestModActionTimestamp(&e.CompiledContentSignals),
	}
}

func NewModActionsTabEntryFromThread(e *CompiledThread) *ModActionsTabEntry {
	return &ModActionsTabEntry{
		Fingerprint:   e.Fingerprint,
		ThreadPayload: *e,
		Timestamp:     getNewestModActionTimestamp(&e.CompiledContentSignals),
	}
}

func NewModActionsTabEntryFromPost(e *CompiledPost) *ModActionsTabEntry {
	return &ModActionsTabEntry{
		Fingerprint: e.Fingerprint,
		PostPayload: *e,
		Timestamp:   getNewestModActionTimestamp(&e.CompiledContentSignals),
	}
}

// TODO handle modapprovals in the future as well, not just modblocks
func getNewestModActionTimestamp(c *CompiledContentSignals) int64 {
	var newest int64
	for k, _ := range c.ModBlocks {
		if stamp := max(c.ModBlocks[k].Creation, c.ModBlocks[k].LastUpdate); stamp > newest {
			newest = stamp
		}
	}
	return newest
}

type ModActionsTabEntryBatch []ModActionsTabEntry

/*=====  End of Mod actions tab entry  ======*/

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}
