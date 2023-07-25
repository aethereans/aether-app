// Services > Nonces

// This service handles the creation, retention and deletion of nonces.

package nonces

import (
	"aether-core/aether/services/toolbox"
	"log"
	"sync"
	"time"
)

/*
Why?

Nonces prevent replay attacks.

Imagine this scenario. A remote sends you a request. This is a POST request specifying an intricate query. This request is signed by the remote's public key, so we know that it's coming from a valid node.

Somehow, that request gets captured. Normally, the connections between nodes are encrypted by TLS, so there is no practical way this can easily happen. But let's say it happened, and the connection got MiTM'ed. The attacker got a copy of the remote's request to us.

What can this attacker, let's call him Mallory, can do to us? Well, it can send that copy of that request thousands of times to us, creating a DDoS that we would need to deal with. It cannot tamper with the request, though, since it is signed. Obviously, the local node will notice that it's being asked to do something off, and it will ban this PK for a period of time. But the PK that is banned will be the PK of the original remote, not Mallory's, since Mallory is just resenting the original remote's request to us many times.

So what Mallory achieved is he managed to get us to ban an innocent third party because he managed to MiTM the connection and sent a copy of the request. That's not the end of the world, but it's also not that great.

So, to prevent that, we are doing two things:

a) We are adding nonces to every request.

b) We are putting a limit on the difference between the local and the remote (called clock skew). This will come handy later.

What do these things do?

Well, nonces prevent resending of the same request. When you look at a nonce, and see that nonce in the list of nonces you've received in the past, you know it was a duplicate request, and deny it. Mallory just got denied.

This creates a problem though. That means you have to keep track of every nonce you've ever received to be able to check against the list, and that list is going to grow very large.

So we have a maximum allowed clock skew. This means, if the difference between our own current time and the remote's current time is larger than maximum allowed clock skew, it won't be accepted.

Why does this help us? Let's say Mallory got a hold of the request with a timestamp of now. He can only use this request to make replay attacks until the maximum allowed clock skew (MACS) is exceeded. If our MACS is 20 minutes, Mallory's copy of the request will be rendered unacceptable 20 minutes later. That means any request whose timestamp is outside +20 -20 boundaries is declined by default. This is great! Because what it means is that we do not have to keep nonces for longer than that range, anything outside that range is actually auto-declined. Anything inside the range is accepted ONLY if the nonce is something we've not seen before.

One gotcha is that it's not just in the past, but also in the future. That is, if your MACS is 20 minutes, that gives you a 40 minutes range, 20 minutes into the past and 20 minutes into the future. So your nonce deletion point should definitely be higher than 40. It's a good call to have it be 2.5x the MACS.

So limiting clock skew makes it so that we only have to keep the nonces that are within the clock skew range.

*/

const (
	maximumAllowedClockSkewMinutes = 20 // MACS
	// ^ We only accept requests that are generated up to this minutes old, or up to this minutes into the future. If the remote's UTC clock and local machine's UTC clock is skewed more than this minutes, the local machine will not accept the request. If this value is 20, for example, it will accept requests timestamped up to 20 minutes into the past, and 20 minutes into the future, meaning it captures 40 minutes.
	maxAllowedRequestsWithinClockSkew = 358
	// ^ This is a rate limiting feature. Effectively, we count the number of nonces that we've received, and if it's above this, we stop responding to that node. This is a naive implementation, but if it turns out to be a problem, we can write some logic that spreads the requests evenly.

	//300 per remote is a best guess, assuming the worst case that a remote does not know any other nodes and it hits you every minute, every sync is 7 post requests (B T P V K TS A), so 50 (MACS*2.5) x 7 = 350. +7 for the worst case where last flush has happened exactly 60 seconds ago, 357. Anything beyond that is very likely malicious.
	minFlushIntervalMinutes = 1
	// ^ We want to flush at most this often, not more often than that.
	nonceExpirationMinutes = maximumAllowedClockSkewMinutes * 2.5
)

var (
	lastflush int64
)

// Initialiser

func NewRemotesNonces() RemotesNonces {
	rn := RemotesNonces{}
	rn.NoncesMap = make(map[publicKey][]nonce)
	return rn
}

// Public types

type RemotesNonces struct {
	lock      sync.Mutex
	NoncesMap map[publicKey][]nonce // needs init
}

// Private types
type publicKey string
type nonce struct {
	nStr     string
	creation int64
}

// checks whether this nonce can be used for entry.
func (n *nonce) afterCutoff(cutoff int64) bool {
	return n.creation > cutoff
}

func (rn *RemotesNonces) flush(cutoff int64) {
	// We don't want to run it every time it's called, only when some time has passed.
	if lastflush > toolbox.CnvToCutoffMinutes(minFlushIntervalMinutes) {
		return
	}
	newRn := make(map[publicKey][]nonce)
	// Go through every pk,
	for key := range rn.NoncesMap {
		// And every nonce,
		for i := range rn.NoncesMap[key] {
			// If the nonce creation is after cutoff (i.e. still within MACS)
			if rn.NoncesMap[key][i].afterCutoff(cutoff) {
				// add it to the new list.
				newRn[key] = append(newRn[key], rn.NoncesMap[key][i])
			}
		}
	}
	// Make new list the main list.
	rn.NoncesMap = newRn
	lastflush = time.Now().Unix()
}

// IsValid checks for nonce validity. If nonextant, it is valid. If extant and within the expiration, it is valid. If pk is extant and nonce has changed, check minimum replacement age, and if it's older than that, replace and it is valid. One PK can only have one nonce
func (rn *RemotesNonces) IsValid(pk, nonceStr string, apiRespTimestamp int64) bool {
	rn.lock.Lock()
	defer rn.lock.Unlock()
	// Guard against empty.
	if len(pk) == 0 || len(nonceStr) == 0 || apiRespTimestamp == 0 {
		return false
	}
	// Guard against too long nonces.
	if len(nonceStr) > 64 {
		return false
	}
	// Guard against clock skew.
	futureClockSkewCutoff := toolbox.CnvToFutureCutoffMinutes(maximumAllowedClockSkewMinutes)
	pastClockSkewCutoff := toolbox.CnvToCutoffMinutes(maximumAllowedClockSkewMinutes)
	if !(apiRespTimestamp < futureClockSkewCutoff &&
		apiRespTimestamp > pastClockSkewCutoff) {
		log.Printf("This API response failed the nonce check because the provided timestamp is too old or too new. Current time: %v, Timestamp given: %v", time.Now().Unix(), apiRespTimestamp)
		return false
	}
	// Set cutoffs.
	expCutoff := toolbox.CnvToCutoffMinutes(nonceExpirationMinutes)
	pkey := publicKey(pk)
	// Flush to remove nonces older than expiration.
	rn.flush(expCutoff)
	// Guard against too many requests (rate limiter)
	if len(rn.NoncesMap[pkey]) > maxAllowedRequestsWithinClockSkew {
		return false
	}
	// Check if the nonce exists.
	for i := range rn.NoncesMap[pkey] {
		if rn.NoncesMap[pkey][i].nStr == nonceStr {
			return false // We already have this nonce, this is a replay.
		}
	}
	// It doesn't exist. Add it to our library, so it can't be reused.
	rn.NoncesMap[pkey] = append(rn.NoncesMap[pkey],
		nonce{
			nStr:     nonceStr,
			creation: time.Now().Unix(),
		})
	return true
}
