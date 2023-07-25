// Backend > Dispatch > Purgatory

// This library is tasked with filtering data that comes in from the network.

/*
  # Why?
  Imagine this. You're receiving data from a node. It sends you some stuff that's in network head, but also other things, that are older than it. Now, in all likelihood, the node has a pretty good reason to send it: those are very likely sent because they are ancestors of some of the entities within the network head. But we don't actually know that.

  What can happen is this — you hit a malicious node where all data is three years old, but it acts as if it is current. In other words, you ask for the network head, and it sends you very old stuff, stuff outside *your* network head.

  Now, the node might not be explicitly malicious. It might be misconfigured. Some guy might have gotten the wrong idea and set his / her network head to 3 years, thinking it will help the network (it won't). But for all intents and purposes that is the same for the receiver.

  What the receiver needs to do is to accept entities that are legitimate ancestors of data that is being sent, and decline those that are not. In more concrete terms:

  a) Get all data, and put all data that is older than your network head into purgatory.
  b) Get all the rest, and create index forms of those data.
  c) For each item in the purgatory, do a graph search, and try to reach a child that is within your network head. If not possible, delete the item, it is of no use to us. (This search happens only within the data given in this specific sync, it does not scan your own database.)
  d) If found, commit the item into database.

  Shortcut case: If there is nothing that is newer than your network head, bail and commit nothing.

  Q: Wouldn't you be able to do this in place? Why do you need to hold the entities themselves?

  A: No, because entities cross-reference each other. This needs to collect items as the sync goes on, and at the end of the sync, do a cross-reference and then insert.
*/

package dispatch

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"

	// "fmt"
	"sync"
	"time"
)

type Purgatory struct {
	lock              sync.Mutex
	BoardsPurg        []api.Board
	ThreadsPurg       []api.Thread
	PostsPurg         []api.Post
	VotesPurg         []api.Vote
	KeysPurg          []api.Key
	TruststatesPurg   []api.Truststate
	BoardIndexes      []api.BoardIndex
	ThreadIndexes     []api.ThreadIndex
	PostIndexes       []api.PostIndex
	VoteIndexes       []api.VoteIndex
	KeyIndexes        []api.KeyIndex
	TruststateIndexes []api.TruststateIndex
}

func (p *Purgatory) indexOf(item api.Provable) int {
	switch entity := item.(type) {
	case *api.Board:
		for key := range p.BoardsPurg {
			if p.BoardsPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *api.Thread:
		for key := range p.ThreadsPurg {
			if p.ThreadsPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *api.Post:
		for key := range p.PostsPurg {
			if p.PostsPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *api.Vote:
		for key := range p.VotesPurg {
			if p.VotesPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *api.Key:
		for key := range p.KeysPurg {
			if p.KeysPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	case *api.Truststate:
		for key := range p.TruststatesPurg {
			if p.TruststatesPurg[key].Fingerprint == entity.Fingerprint {
				return key
			}
		}
	}
	return -1
}

// func (p *Purgatory) remove(item api.Provable) { // there's something weird here.
// 	switch item.(type) {
// 	case *api.Board:
// 		i := p.indexOf(item)
// 		if i != -1 {
// 			newPurg := append(p.BoardsPurg[0:i], p.BoardsPurg[i+1:len(p.BoardsPurg)]...)
// 			p.BoardsPurg = newPurg
// 		}
// 	case *api.Thread:
// 		i := p.indexOf(item)
// 		if i != -1 {
// 			newPurg := append(p.ThreadsPurg[0:i], p.ThreadsPurg[i+1:len(p.ThreadsPurg)]...)
// 			p.ThreadsPurg = newPurg
// 		}
// 	case *api.Post:
// 		i := p.indexOf(item)
// 		if i != -1 {
// 			newPurg := append(p.PostsPurg[0:i], p.PostsPurg[i+1:len(p.PostsPurg)]...)
// 			p.PostsPurg = newPurg
// 		}
// 	case *api.Vote:
// 		i := p.indexOf(item)
// 		fmt.Println("i in remove")
// 		fmt.Println(i)
// 		if i != -1 {
// 			newPurg := append(p.VotesPurg[0:i], p.VotesPurg[i+1:len(p.VotesPurg)]...)
// 			p.VotesPurg = newPurg
// 		}
// 	case *api.Key:
// 		i := p.indexOf(item)
// 		if i != -1 {
// 			newPurg := append(p.KeysPurg[0:i], p.KeysPurg[i+1:len(p.KeysPurg)]...)
// 			p.KeysPurg = newPurg
// 		}
// 	case *api.Truststate:
// 		i := p.indexOf(item)
// 		if i != -1 {
// 			newPurg := append(p.TruststatesPurg[0:i], p.TruststatesPurg[i+1:len(p.TruststatesPurg)]...)
// 			p.TruststatesPurg = newPurg
// 		}
// 	}
// }

type vertex struct {
	entityType   string
	fingerprint  api.Fingerprint
	owner        api.Fingerprint
	lastModified api.Timestamp // creation or last update, whichever is never
}

func cnvToVertex(item api.Provable) vertex {
	v := vertex{}
	v.fingerprint = item.GetFingerprint()
	v.lastModified = item.GetLastModified()
	v.entityType = item.GetEntityType()
	v.owner = item.GetOwner()
	return v
}

func cnvToVertexIdx(item api.ProvableIndex) vertex {
	v := vertex{}
	v.fingerprint = item.GetFingerprint()
	v.lastModified = item.GetLastModified()
	v.entityType = item.GetEntityType()
	v.owner = item.GetOwner()
	return v
}

func (p *Purgatory) getDirectDescendants(vs []vertex) []vertex {
	var children []vertex

	for key1 := range vs {
		switch vs[key1].entityType {
		case "board":
			for key := range p.ThreadIndexes {
				if p.ThreadIndexes[key].Board == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.ThreadIndexes[key]))
				}
			}
		case "thread":
			for key := range p.PostIndexes {
				if p.PostIndexes[key].Thread == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.PostIndexes[key]))
				}
			}
		case "post":
			for key := range p.PostIndexes {
				if p.PostIndexes[key].Parent == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.PostIndexes[key]))
				}
			}
		case "vote":
			// no descendants
		case "key":
			for key := range p.BoardIndexes {
				if p.BoardIndexes[key].Owner == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.BoardIndexes[key]))
				}
			}
			for key := range p.ThreadIndexes {
				if p.ThreadIndexes[key].Owner == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.ThreadIndexes[key]))
				}
			}
			for key := range p.PostIndexes {
				if p.PostIndexes[key].Owner == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.PostIndexes[key]))
				}
			}
			for key := range p.VoteIndexes {
				if p.VoteIndexes[key].Owner == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.VoteIndexes[key]))
				}
			}
			for key := range p.TruststateIndexes {
				if p.TruststateIndexes[key].Owner == vs[key1].fingerprint {
					children = append(children, cnvToVertexIdx(&p.TruststateIndexes[key]))
				}
			}
		case "truststate":
			// no descendants.
			// This means if a truststate is outside the network head, it is not taken in.
		}
		logging.Logf(2, "For item :%#v, Direct descendants found: %#v", vs[key1], children)
	}
	return children
}

func getMostRecentLastModified(vs []vertex) api.Timestamp {
	var mostRecentLm api.Timestamp
	for key := range vs {
		if vs[key].lastModified > mostRecentLm {
			mostRecentLm = vs[key].lastModified
		}
	}
	return mostRecentLm
}

// ok, what ancestors are available to us?
/*
  updated by list

  #board is updated by:
  threads
  posts

  #thread is updated by
  posts

  #post is updated by
  posts

  #vote is updated by
  nothing

  key is updated by
  boards
  threads
  thread through board
  posts
  post through thread
  post through thread through board
  post through post .. through post
  votes
  truststates owner
  truststates target

  truststate is updated by
  nothing
*/

// we have two conditions, get children, and so long as the NH isn't achieved, or we run out of items, in the children, we repeat. We also should have a max depth, something super high, but finite. So that if somebody comes up with a cute little idea of coming up with circular references, it won't break this.
func (p *Purgatory) verify(item api.Provable) bool {
	cutoff := api.Timestamp(globals.BackendConfig.GetEventHorizonTimestamp())
	// nhD := globals.BackendConfig.GetNetworkHeadDays()
	// nhCutoff := api.Timestamp(toolbox.CnvToCutoffDays(nhD))
	v := cnvToVertex(item)
	toBeSearched := []vertex{v}
	itercount := 0
	var mostRecentLastModified api.Timestamp
	for cutoff > mostRecentLastModified {
		if itercount > 1000 {
			return false
		}
		dd := p.getDirectDescendants(toBeSearched)
		if len(dd) == 0 {
			// logging.Logf(2, "This item ran out of ancestors before achieving cutoff. Item: %#v", v)
			return false
		} else {
			toBeSearched = dd
		}
		mrlm := getMostRecentLastModified(dd)
		if mrlm > mostRecentLastModified {
			mostRecentLastModified = mrlm
		}
		itercount++
	}
	logging.Logf(2, "This item was successfully mapped to a child within the network head. Item: %#v", v)
	return true
}

func (p *Purgatory) process() {
	var newB []api.Board

	var newT []api.Thread

	var newP []api.Post

	var newV []api.Vote

	var newK []api.Key

	var newTs []api.Truststate

	for key := range p.BoardsPurg {
		if p.verify(&p.BoardsPurg[key]) {
			newB = append(newB, p.BoardsPurg[key])
		}
	}
	for key := range p.ThreadsPurg {
		if p.verify(&p.ThreadsPurg[key]) {
			newT = append(newT, p.ThreadsPurg[key])
		}
	}
	for key := range p.PostsPurg {
		if p.verify(&p.PostsPurg[key]) {
			newP = append(newP, p.PostsPurg[key])
		}
	}
	for key := range p.VotesPurg {
		if p.verify(&p.VotesPurg[key]) {
			newV = append(newV, p.VotesPurg[key])
		}
	}
	for key := range p.KeysPurg {
		if p.verify(&p.KeysPurg[key]) {
			newK = append(newK, p.KeysPurg[key])
		}
	}
	for key := range p.TruststatesPurg {
		if p.verify(&p.TruststatesPurg[key]) {
			newTs = append(newTs, p.TruststatesPurg[key])
		}
	}
	p.BoardsPurg = newB
	p.ThreadsPurg = newT
	p.PostsPurg = newP
	p.VotesPurg = newV
	p.KeysPurg = newK
	p.TruststatesPurg = newTs
}

func (p *Purgatory) convertAllToIface() []interface{} {
	var carrier []interface{}
	for i := range p.BoardsPurg {
		carrier = append(carrier, p.BoardsPurg[i])
	}
	for i := range p.ThreadsPurg {
		carrier = append(carrier, p.ThreadsPurg[i])
	}
	for i := range p.PostsPurg {
		carrier = append(carrier, p.PostsPurg[i])
	}
	for i := range p.VotesPurg {
		carrier = append(carrier, p.VotesPurg[i])
	}
	for i := range p.KeysPurg {
		carrier = append(carrier, p.KeysPurg[i])
	}
	for i := range p.TruststatesPurg {
		carrier = append(carrier, p.TruststatesPurg[i])
	}
	return carrier
}

// Process goes through all items in the purgatory, and sends back items it finds valid to be committed into the DB.
func (p *Purgatory) Process() []interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()
	start := time.Now()
	p.process()
	elapsed := time.Since(start)
	logging.Logf(2, "This purgatory process run took %vs.", toolbox.Round(elapsed.Seconds(), 0.1))
	logging.Logf(2, "At the end of this purgatory process run, this is our purgatory: B: %v, T: %v, P: %v, V: %v, K: %v, TS: %v.\n", len(p.BoardsPurg), len(p.ThreadsPurg), len(p.PostsPurg), len(p.VotesPurg), len(p.KeysPurg), len(p.TruststatesPurg))
	resultAsIface := p.convertAllToIface()
	logging.Logf(2, "This is the length of the result that is going to be inserted after the purgatory process. Length: %v", len(resultAsIface))
	return resultAsIface
}

// func calcGate(item api.Provable) api.Timestamp {
// 	lu := item.GetLastUpdate()
// 	cr := item.GetCreation()
// 	var gate api.Timestamp
// 	if lu > cr {
// 		gate = lu
// 	} else {
// 		gate = cr
// 	}
// 	return gate
// }

// Accept goes through the given list of entities, and looks at those that might actually go into the purgatory. It takes them in, and sends back the fingerprints of those that are taken in, so that the receiver can remove it from the list.
func (p *Purgatory) accept(items []api.Provable) {
	// takenIn := []api.Fingerprint{}
	// nhD := globals.BackendConfig.GetNetworkHeadDays()
	// nhCutoff := api.Timestamp(toolbox.CnvToCutoffDays(nhD))
	cutoff := api.Timestamp(globals.BackendConfig.GetEventHorizonTimestamp()) // todo: should purgatory be gated on network head or event horizon?
	for key := range items {
		// gate := calcGate(items[key])
		if cutoff > items[key].GetLastModified() {
			// Entity older than our network head. Enters purgatory.
			switch entity := items[key].(type) {
			case *api.Board:
				p.BoardsPurg = append(p.BoardsPurg, *entity)
			case *api.Thread:
				p.ThreadsPurg = append(p.ThreadsPurg, *entity)
			case *api.Post:
				p.PostsPurg = append(p.PostsPurg, *entity)
			case *api.Vote:
				p.VotesPurg = append(p.VotesPurg, *entity)
			case *api.Key:
				p.KeysPurg = append(p.KeysPurg, *entity)
			case *api.Truststate:
				p.TruststatesPurg = append(p.TruststatesPurg, *entity)
			}
		} else {
			// Entity is within the network head. We don't take it into the purgatory, but we get its index, so we can use it to search for items in the purgatory.
			switch entity := items[key].(type) {
			case *api.Board:
				index := api.BoardIndex{
					Fingerprint: entity.Fingerprint,
					Owner:       entity.Owner,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.BoardIndexes = append(p.BoardIndexes, index)
			case *api.Thread:
				index := api.ThreadIndex{
					Fingerprint: entity.Fingerprint,
					Owner:       entity.Owner,
					Board:       entity.Board,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.ThreadIndexes = append(p.ThreadIndexes, index)
			case *api.Post:
				index := api.PostIndex{
					Fingerprint: entity.Fingerprint,
					Owner:       entity.Owner,
					Board:       entity.Board,
					Thread:      entity.Thread,
					Parent:      entity.Parent,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.PostIndexes = append(p.PostIndexes, index)
			case *api.Vote:
				index := api.VoteIndex{
					Fingerprint: entity.Fingerprint,
					Owner:       entity.Owner,
					Board:       entity.Board,
					Thread:      entity.Thread,
					Target:      entity.Target,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.VoteIndexes = append(p.VoteIndexes, index)
			case *api.Key:
				index := api.KeyIndex{
					Fingerprint: entity.Fingerprint,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.KeyIndexes = append(p.KeyIndexes, index)
			case *api.Truststate:
				index := api.TruststateIndex{
					Fingerprint: entity.Fingerprint,
					Owner:       entity.Owner,
					Target:      entity.Target,
					Creation:    entity.Creation,
					LastUpdate:  entity.LastUpdate,
				}
				p.TruststateIndexes = append(p.TruststateIndexes, index)
			}
		}
	}
}

func (p *Purgatory) removeFromResp(r *api.Response) {
	if len(r.Boards) > 0 {
		var removalIdxs []int

		for key := range p.BoardsPurg {
			idx := r.IndexOf(&p.BoardsPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "board")
	}
	if len(r.Threads) > 0 {
		var removalIdxs []int

		for key := range p.ThreadsPurg {
			idx := r.IndexOf(&p.ThreadsPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "thread")
	}
	if len(r.Posts) > 0 {
		var removalIdxs []int

		for key := range p.PostsPurg {
			idx := r.IndexOf(&p.PostsPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "post")
	}
	if len(r.Votes) > 0 {
		var removalIdxs []int

		for key := range p.VotesPurg {
			idx := r.IndexOf(&p.VotesPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "vote")
	}
	if len(r.Keys) > 0 {
		var removalIdxs []int

		for key := range p.KeysPurg {
			idx := r.IndexOf(&p.KeysPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "key")
	}
	if len(r.Truststates) > 0 {
		var removalIdxs []int

		for key := range p.TruststatesPurg {
			idx := r.IndexOf(&p.TruststatesPurg[key])
			if idx != -1 {
				removalIdxs = append(removalIdxs, idx)
			}
		}
		r.MassRemoveByIndex(removalIdxs, "truststate")
	}
}

func (p *Purgatory) Filter(r *api.Response) {
	p.lock.Lock()
	defer p.lock.Unlock()
	start := time.Now()
	var bProv []api.Provable

	for key := range r.Boards {
		bProv = append(bProv, api.Provable(&r.Boards[key]))
	}
	var tProv []api.Provable

	for key := range r.Threads {
		tProv = append(tProv, api.Provable(&r.Threads[key]))
	}
	var pProv []api.Provable

	for key := range r.Posts {
		pProv = append(pProv, api.Provable(&r.Posts[key]))
	}
	var vProv []api.Provable

	for key := range r.Votes {
		vProv = append(vProv, api.Provable(&r.Votes[key]))
	}
	var kProv []api.Provable

	for key := range r.Keys {
		kProv = append(kProv, api.Provable(&r.Keys[key]))
	}
	var tsProv []api.Provable

	for key := range r.Truststates {
		tsProv = append(tsProv, api.Provable(&r.Truststates[key]))
	}
	p.accept(bProv)
	p.accept(tProv)
	p.accept(pProv)
	p.accept(vProv)
	p.accept(kProv)
	p.accept(tsProv)
	// Here, we have the whole of this response filtered through.
	// Now, we need to take a look at the purgatory, and remove anything in the purgatory from the response.
	p.removeFromResp(r)
	logging.Logf(2, "At the end of this purgatory run, this is our response: B: %v, T: %v, P: %v, V: %v, K: %v, TS: %v.\n", len(r.Boards), len(r.Threads), len(r.Posts), len(r.Votes), len(r.Keys), len(r.Truststates))

	elapsed := time.Since(start)
	logging.Logf(2, "At the end of this purgatory filter run, this is our purgatory: B: %v, T: %v, P: %v, V: %v, K: %v, TS: %v.\n", len(p.BoardsPurg), len(p.ThreadsPurg), len(p.PostsPurg), len(p.VotesPurg), len(p.KeysPurg), len(p.TruststatesPurg))
	logging.Logf(2, "This purgatory filter run took %vs.", toolbox.Round(elapsed.Seconds(), 0.1))
}
