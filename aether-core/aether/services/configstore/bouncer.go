// Services > ConfigStore > Bouncer
// This module controls how many remotes can connect to this computer at the same time, and how many outbound syncs can be happening at the same time.

package configstore

import (
	"aether-core/aether/services/toolbox"
	"fmt"
	"log"
	// "github.com/davecgh/go-spew/spew"
	"aether-core/aether/services/extverify"
	// "log"
	"sync"
	"time"
)

// These two are local variables that only affect this specific library. Since there is no reason to modify them from the outside, this is not brought into the main settings JSON.
const (
	activeInboundLeaseDurationSeconds  = 60  // 1m
	activeOutboundLeaseDurationSeconds = 900 // 15m
	activePingLeaseDurationSeconds     = 10  // 10s

	historyInboundLeaseDurationSeconds  = 86400 // 1d
	historyOutboundLeaseDurationSeconds = 86400 // 1d

	minimumActivesFlushIntervalSeconds = 30   // 30s
	minimumHistoryFlushIntervalSeconds = 3600 // 1h
	// ^ This relates to dropping from history, not adding to history. Adding to history happens on active flush.
)

type Bouncer struct {
	lock             sync.Mutex
	Inbounds         []ConnectionRecord
	Outbounds        []ConnectionRecord
	Pings            []ConnectionRecord
	InboundHistory   []ConnectionRecord
	OutboundHistory  []ConnectionRecord
	ActivesLastFlush Timestamp
	HistoryLastFlush Timestamp
}

type ConnectionRecord struct {
	Location            string
	Sublocation         string
	Port                uint16
	FirstAccess         Timestamp
	LastAccess          Timestamp
	Inbound_ReverseConn bool
	// ^ If inbound and if reverse conn, true
	Inbound_ReverseConn_Successful bool
	// ^ If inbound and if reverse conn and if successful, true
	Outbound_ReverseConn bool
	// ^ If outbound and if done in response to reverse conn, true
	Outbound_Successful bool
	// ^ If outbound and if successful, true
	ConnDurationSeconds float64
}

func (n *ConnectionRecord) equal(c ConnectionRecord) bool {
	// fmt.Printf("Equal was called between %v and %v, and the result was %v\n", n, c, eqresult)
	return n.Location == c.Location &&
		n.Sublocation == n.Sublocation &&
		n.Inbound_ReverseConn == c.Inbound_ReverseConn &&
		n.Outbound_ReverseConn == c.Outbound_ReverseConn
	// n.Port == c.Port && // This is causing a little too much pain. Let's remove that and see what happens.
}
func (n *ConnectionRecord) hasActiveInboundLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(activeInboundLeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}
}
func (n *ConnectionRecord) hasActiveOutboundLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(activeOutboundLeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}
}

func (n *ConnectionRecord) hasActivePingLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(activePingLeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}
}

func (n *ConnectionRecord) hasHistoryInboundLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(historyInboundLeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}
}
func (n *ConnectionRecord) hasHistoryOutboundLease() bool {
	cutoff := Timestamp(time.Now().Add(-(time.Duration(historyOutboundLeaseDurationSeconds) * time.Second)).Unix())
	if n.LastAccess > cutoff {
		return true
	} else {
		return false
	}
}

func (n *Bouncer) indexOf(direction string, loc, subloc string, port uint16, inboundRev, outboundRev bool) int {
	switch direction {
	case "inbound":
		for key, _ := range n.Inbounds {
			if n.Inbounds[key].equal(ConnectionRecord{Location: loc, Sublocation: subloc, Port: port, Inbound_ReverseConn: inboundRev, Outbound_ReverseConn: outboundRev}) {
				return key
			}
		}
		return -1
	case "outbound":
		for key, _ := range n.Outbounds {
			if n.Outbounds[key].equal(ConnectionRecord{Location: loc, Sublocation: subloc, Port: port, Inbound_ReverseConn: inboundRev, Outbound_ReverseConn: outboundRev}) {
				return key
			}
		}
		return -1
	case "ping":
		for key, _ := range n.Pings {
			if n.Pings[key].equal(ConnectionRecord{Location: loc, Sublocation: subloc, Port: port}) {
				return key
			}
		}
		return -1
	default:
		return -1
	}
}

func (n *Bouncer) insert(direction string, loc, subloc string, port uint16, isReverseConn bool) {
	now := Timestamp(time.Now().Unix())
	entry := ConnectionRecord{Location: loc, Sublocation: subloc, Port: port, FirstAccess: now, LastAccess: now}
	switch direction {
	case "inbound":
		entry.Inbound_ReverseConn = isReverseConn
		n.Inbounds = append(n.Inbounds, entry)
	case "outbound":
		entry.Outbound_ReverseConn = isReverseConn
		n.Outbounds = append(n.Outbounds, entry)
	case "ping":
		n.Pings = append(n.Pings, entry)
	default:
		panic(fmt.Sprintf("You gave an invalid direction to insert in Bouncer. Direction: %v", direction))
	}
}

func (n *Bouncer) removeItem(direction string, i int) {
	finalList := []ConnectionRecord{}
	switch direction {
	case "inbound":
		finalList = append(n.Inbounds[0:i], n.Inbounds[i+1:len(n.Inbounds)]...)
		n.Inbounds[i].ConnDurationSeconds = calcDuration(n.Inbounds[i])
		// move to inbound history
		n.InboundHistory = append(n.InboundHistory, n.Inbounds[i])
		inboundExpiredHook(n.Inbounds[i])
		n.Inbounds = finalList
	case "outbound":
		finalList = append(n.Outbounds[0:i], n.Outbounds[i+1:len(n.Outbounds)]...)
		n.Outbounds[i].ConnDurationSeconds = calcDuration(n.Outbounds[i])
		// move to outbound history
		n.OutboundHistory = append(n.OutboundHistory, n.Outbounds[i])
		outboundExpiredHook(n.Outbounds[i])
		n.Outbounds = finalList
	case "ping":
		finalList = append(n.Pings[0:i], n.Pings[i+1:len(n.Pings)]...)
		n.Pings[i].ConnDurationSeconds = calcDuration(n.Pings[i])
		pingExpiredHook(n.Pings[i])
		n.Pings = finalList
	case "inboundHistory":
		finalList = append(n.InboundHistory[0:i], n.InboundHistory[i+1:len(n.InboundHistory)]...)
		n.InboundHistory[i].ConnDurationSeconds = calcDuration(n.InboundHistory[i])
		n.InboundHistory = finalList
	case "outboundHistory":
		finalList = append(n.OutboundHistory[0:i], n.OutboundHistory[i+1:len(n.OutboundHistory)]...)
		n.OutboundHistory[i].ConnDurationSeconds = calcDuration(n.OutboundHistory[i])
		n.OutboundHistory = finalList
	}
}

func (n *Bouncer) flush() {
	n.flushActives()
	n.flushHistory()
}

func (n *Bouncer) flushActives() {
	// If there's been a flush in the past 1 min, ignore flush. This is because flush is in a hot path, we want to avoid unnecessary repeats.
	if n.ActivesLastFlush > Timestamp(time.Now().Add(-(time.Duration(minimumActivesFlushIntervalSeconds) * time.Second)).Unix()) {
		return
	}
	// Set ActivesLastFlush to now if the gate above passes.
	n.ActivesLastFlush = Timestamp(time.Now().Add(-(time.Duration(minimumActivesFlushIntervalSeconds) * time.Second)).Unix())
	for i := len(n.Inbounds) - 1; i >= 0; i-- {
		if !n.Inbounds[i].hasActiveInboundLease() {
			n.removeItem("inbound", i)
		}
	}
	for i := len(n.Outbounds) - 1; i >= 0; i-- {
		if !n.Outbounds[i].hasActiveOutboundLease() {
			n.removeItem("outbound", i)
		}
	}
	for i := len(n.Pings) - 1; i >= 0; i-- {
		if !n.Pings[i].hasActivePingLease() {
			n.removeItem("ping", i)
		}
	}
}

func (n *Bouncer) flushHistory() {
	if n.HistoryLastFlush > Timestamp(time.Now().Add(-(time.Duration(minimumHistoryFlushIntervalSeconds) * time.Second)).Unix()) {
		return
	}
	n.HistoryLastFlush = Timestamp(time.Now().Add(-(time.Duration(minimumHistoryFlushIntervalSeconds) * time.Second)).Unix())
	for i := len(n.InboundHistory) - 1; i >= 0; i-- {
		if !n.InboundHistory[i].hasHistoryInboundLease() {
			n.removeItem("inboundHistory", i)
		}
	}
	for i := len(n.OutboundHistory) - 1; i >= 0; i-- {
		if !n.OutboundHistory[i].hasHistoryOutboundLease() {
			n.removeItem("outboundHistory", i)
		}
	}
}

func (n *Bouncer) RequestInboundLease(loc, subloc, proxy string, port uint16, isReverseConn bool) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	// If we're lameduck, decline.
	if Btc.LameduckInitiated || Btc.ShutdownInitiated {
		return false
	}
	if bc.GetExternalVerifyEnabled() {
		// log.Printf("External verify enabled. Remote IP: %v", loc)
		if !extverify.Verifier.IsAllowedRemoteIP(proxy) {
			// log.Printf("Remote wasn't allowed in because it wasn't a CF ip. Remote IP: %v", loc)
			return false
		}
	}
	n.flush()
	// fmt.Println("An inbound lease was requested.")
	/*
		See "reverse" bypasses on both the true cases below. If it's a reverse connection, we always accept. If you want to limit reverse connections, do it from where we trigger them, all reverses are locally triggered.
	*/
	leaseIndex := n.indexOf("inbound", loc, subloc, port, isReverseConn, false)
	if leaseIndex != -1 && n.Inbounds[leaseIndex].hasActiveInboundLease() {
		// fmt.Println("Lease was renewed.")
		// fmt.Printf("lease index: %v, inbounds: %#v", leaseIndex, n.Inbounds)
		n.Inbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		return true
	} else {
		if len(n.Inbounds) < bc.GetMaxInboundConns() ||
			isReverseConn {
			n.insert("inbound", loc, subloc, port, isReverseConn)
			// fmt.Println("Lease was granted.")
			allocated, max := n.GetInboundSaturation()
			logf(0, "GIVEN INBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
			return true
		} else {
			// fmt.Println("A lease was denied.")
			allocated, max := n.GetInboundSaturation()
			logf(0, "DENIED INBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
			return false
		}
	}
}

func (n *Bouncer) RequestOutboundLease(loc, subloc string, port uint16, isReverseConn bool) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	// If we're lameduck, decline.
	if Btc.LameduckInitiated || Btc.ShutdownInitiated {
		return false
	}
	if isReverseConn {
		/*
			The name of this function is outbounds, this method below says inbound, that is not a mistake or typo.

			In the case of a reverse open request received from the remote, when this is called, the sync is asking for permission to reach out. So it is response to an inbound reverse request. If our inbound reverse is disabled, then we do not want to reach out in response.
		*/
		if bc.GetDeclineInboundReverseRequests() {
			return false
		}
	}
	// fmt.Println("outbound lease requested")
	// fmt.Println(loc, subloc, port, isReverseConn)
	n.flush()
	leaseIndex := n.indexOf("outbound", loc, subloc, port, false, isReverseConn)
	if leaseIndex != -1 && n.Outbounds[leaseIndex].hasActiveOutboundLease() {
		n.Outbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		return true
	} else {
		if len(n.Outbounds) < bc.GetMaxOutboundConns() { // || isReverseConn
			/*
				Removed: automatically giving reverse connections an outbound lease.
			*/
			n.insert("outbound", loc, subloc, port, isReverseConn)
			allocated, max := n.GetOutboundSaturation()
			logf(0, "GIVEN OUTBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
			return true
		} else {
			allocated, max := n.GetOutboundSaturation()
			logf(0, "DENIED OUTBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
			return false
		}
	}
}

// RequestPingLease allows a ping from a remote to be allowed in. Since this does not make sense in the case of reverse (i.e. if a remote has made an inbound request, it is a sync request, not a ping request.)
func (n *Bouncer) RequestPingLease(loc, subloc, proxy string, port uint16) bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	// If we're lameduck, decline.
	if Btc.LameduckInitiated || Btc.ShutdownInitiated {
		return false
	}
	if bc.GetExternalVerifyEnabled() {
		// log.Printf("External verify enabled. Remote IP: %v", loc)
		if !extverify.Verifier.IsAllowedRemoteIP(proxy) {
			// log.Printf("Remote wasn't allowed in because it wasn't a CF ip. Remote IP: %v", loc)
			return false
		}
	}
	n.flush()
	/*
		See "reverse" bypasses on both the true cases below. If it's a reverse connection, we always accept. If you want to limit reverse connections, do it from where we trigger them, all reverses are locally triggered.
	*/
	leaseIndex := n.indexOf("ping", loc, subloc, port, false, false)
	if leaseIndex != -1 && n.Pings[leaseIndex].hasActivePingLease() {
		// fmt.Println("Lease was renewed.")
		// fmt.Printf("lease index: %v, pings: %#v", leaseIndex, n.Pings)
		n.Pings[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		return true
	} else {
		if len(n.Pings) < bc.GetMaxPingConns() {
			n.insert("ping", loc, subloc, port, false)
			// fmt.Println("Lease was granted.")
			allocated, max := n.GetPingSaturation()
			logf(0, "GIVEN PING LEASE: (%v/%v) %v:%v", allocated, max, loc, port)
			return true
		} else {
			// fmt.Println("A lease was denied.")
			allocated, max := n.GetPingSaturation()
			logf(0, "DENIED PING LEASE: (%v/%v) %v:%v", allocated, max, loc, port)
			return false
		}
	}
}

// ReleaseOutboundLease is idempotent if there is no such lease.
func (n *Bouncer) ReleaseOutboundLease(loc, subloc string, port uint16, wasSuccessful, isReverseConn bool) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.flush()
	leaseIndex := n.indexOf("outbound", loc, subloc, port, false, isReverseConn)
	// fmt.Printf("Lease index was %v\n", leaseIndex)
	if leaseIndex != -1 {
		n.Outbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		n.Outbounds[leaseIndex].Outbound_Successful = wasSuccessful
		n.removeItem("outbound", leaseIndex)
		allocated, max := n.GetOutboundSaturation()
		logf(0, "RELEASED OUTBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
	}
}

// ReleaseInboundLease is idempotent if there is no such lease.
func (n *Bouncer) ReleaseInboundLease(loc, subloc string, port uint16, wasSuccessful, isReverseConn bool) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.flush()
	leaseIndex := n.indexOf("inbound", loc, subloc, port, isReverseConn, false)
	// fmt.Printf("Lease index was %v\n", leaseIndex)
	if leaseIndex != -1 {
		n.Inbounds[leaseIndex].LastAccess = Timestamp(time.Now().Unix())
		if isReverseConn && wasSuccessful {
			n.Inbounds[leaseIndex].Inbound_ReverseConn_Successful = wasSuccessful
		}
		// ^ Inbounds are considered successful when they expire, the only exception being reverse conns, which we assume to be successful only if we are explicitly told so.
		n.removeItem("inbound", leaseIndex)
		allocated, max := n.GetInboundSaturation()
		logf(0, "RELEASED INBOUND LEASE: (%v/%v) %v:%v. Type: %v", allocated, max, loc, port, getConnType(isReverseConn))
	}
}

/*
	ReleaseInboundLease and ReleasePingLease don't exist because those are not in our control. When the remote is done, it just stops asking, and we just expire them.
*/

// Below are the high level methods that comprise of the public API.

func (b *Bouncer) GetLastInboundSyncTimestamp(onlyReverseConn bool) int64 {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.flush()
	ts := Timestamp(0)
	if onlyReverseConn {
		for key, _ := range b.Inbounds {
			if b.Inbounds[key].Inbound_ReverseConn &&
				b.Inbounds[key].LastAccess > ts {
				ts = b.Inbounds[key].LastAccess
			}
		}
		for key, _ := range b.InboundHistory {
			if b.InboundHistory[key].Inbound_ReverseConn &&
				b.InboundHistory[key].LastAccess > ts {
				ts = b.InboundHistory[key].LastAccess
			}
		}
		return int64(ts)
	}
	for key, _ := range b.Inbounds {
		if b.Inbounds[key].LastAccess > ts {
			ts = b.Inbounds[key].LastAccess
		}
	}
	for key, _ := range b.InboundHistory {
		if b.InboundHistory[key].LastAccess > ts {
			ts = b.InboundHistory[key].LastAccess
		}
	}
	return int64(ts)
}

func (b *Bouncer) GetLastOutboundSyncTimestamp(onlySuccessful bool) int64 {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.flush()
	ts := Timestamp(0)
	if onlySuccessful {
		for key, _ := range b.Outbounds {
			if b.Outbounds[key].Outbound_Successful &&
				b.Outbounds[key].LastAccess > ts {
				ts = b.Outbounds[key].LastAccess
			}
		}
		for key, _ := range b.OutboundHistory {
			if b.OutboundHistory[key].Outbound_Successful &&
				b.OutboundHistory[key].LastAccess > ts {
				ts = b.OutboundHistory[key].LastAccess
			}
		}
		return int64(ts)
	}
	for key, _ := range b.Outbounds {
		if b.Outbounds[key].LastAccess > ts {
			ts = b.Outbounds[key].LastAccess
		}
	}
	for key, _ := range b.OutboundHistory {
		if b.OutboundHistory[key].LastAccess > ts {
			ts = b.OutboundHistory[key].LastAccess
		}
	}
	return int64(ts)
}

func (b *Bouncer) GetLastPingTimestamp(onlySuccessful bool) int64 {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.flush()
	ts := Timestamp(0)
	for key, _ := range b.Pings {
		if b.Pings[key].LastAccess > ts {
			ts = b.Pings[key].LastAccess
		}
	}
	return int64(ts)
}

func (b *Bouncer) GetInboundsInLastXMinutes(min uint, onlySuccessful bool) []ConnectionRecord {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.flush()
	cutoff := Timestamp(toolbox.CnvToCutoffMinutes(int(min)))
	results := []ConnectionRecord{}
	/*
		Logic below:
		- A normal (non-reverse) inbound is interpreted successful by default
		- A reverse inbound is successful only if explicitly marked as such
		- If only looking for successful, all normals, and explicitly-marked-as-successful reverses, gated by minute time limit
		- If looking for all, just all of them as gated by time limit
	*/
	if onlySuccessful {
		for key, _ := range b.Inbounds {
			if b.Inbounds[key].LastAccess <= cutoff {
				// If not within the X minutes we want, pass without further processing
				continue
			}
			if !b.Inbounds[key].Inbound_ReverseConn {
				// If normal (not a reverse) conn, we count it as successful by default.
				results = append(results, b.Inbounds[key])
				continue
			}
			// This was a reverse, it is successful *only* if it is explicitly marked as such.
			if b.Inbounds[key].Inbound_ReverseConn_Successful {
				results = append(results, b.Inbounds[key])
			}
		}
		// Same as above, repeated below over inbound history
		for key, _ := range b.InboundHistory {
			if b.InboundHistory[key].LastAccess <= cutoff {
				// If not within the X minutes we want, pass without further processing
				continue
			}
			if !b.InboundHistory[key].Inbound_ReverseConn {
				// If normal (not a reverse) conn, we count it as successful by default.
				results = append(results, b.InboundHistory[key])
				continue
			}
			// This was a reverse, it is successful *only* if it is explicitly marked as such.
			if b.InboundHistory[key].Inbound_ReverseConn_Successful {
				results = append(results, b.InboundHistory[key])
			}
		}
		return results
	}
	// Not only successful, but all
	for key, _ := range b.Inbounds {
		if b.Inbounds[key].LastAccess > cutoff {
			results = append(results, b.Inbounds[key])
		}
	}
	for key, _ := range b.InboundHistory {
		if b.InboundHistory[key].LastAccess > cutoff {
			results = append(results, b.InboundHistory[key])
		}
	}
	return results
}

func (b *Bouncer) GetOutboundsInLastXMinutes(min uint, onlySuccessful bool) []ConnectionRecord {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.flush()
	cutoff := Timestamp(toolbox.CnvToCutoffMinutes(int(min)))
	results := []ConnectionRecord{}
	if onlySuccessful {
		for key, _ := range b.Outbounds {
			if b.Outbounds[key].Outbound_Successful &&
				b.Outbounds[key].LastAccess > cutoff {
				results = append(results, b.Outbounds[key])
			}
		}
		for key, _ := range b.OutboundHistory {
			if b.OutboundHistory[key].Outbound_Successful &&
				b.OutboundHistory[key].LastAccess > cutoff {
				results = append(results, b.OutboundHistory[key])
			}
		}
		return results
	}
	for key, _ := range b.Outbounds {
		if b.Outbounds[key].LastAccess > cutoff {
			results = append(results, b.Outbounds[key])
		}
	}
	for key, _ := range b.OutboundHistory {
		if b.OutboundHistory[key].LastAccess > cutoff {
			results = append(results, b.OutboundHistory[key])
		}
	}
	return results
}

func (b *Bouncer) GetInboundSaturation() (used, total int) {
	return len(b.Inbounds), bc.GetMaxInboundConns()
}
func (b *Bouncer) GetOutboundSaturation() (used, total int) {
	return len(b.Outbounds), bc.GetMaxOutboundConns()
}
func (b *Bouncer) GetPingSaturation() (used, total int) {
	return len(b.Pings), bc.GetMaxPingConns()
}

func getConnType(isReverseConn bool) string {
	if isReverseConn {
		return "Reverse"
	}
	return "Normal"
}

func calcDuration(c ConnectionRecord) float64 {
	fa := int64(c.FirstAccess)
	la := int64(c.LastAccess)
	diff := time.Duration(la-fa) * time.Second
	return toolbox.Round(diff.Seconds(), 0.1)
}

// These are empty in the case you need to take some action when these events happen.

// inboundExpiredHook runs after an inbound connection lease expires, just before it is fully removed from the inbounds list.
func inboundExpiredHook(c ConnectionRecord) {
}

// outboundExpiredHook runs after an outbound connection lease expires, just before it is fully removed from the outbounds list.
func outboundExpiredHook(c ConnectionRecord) {
}

// pingExpiredHook runs after a ping connection lease expires, just before it is fully removed from the pings list.
func pingExpiredHook(c ConnectionRecord) {
}

/*========================================
=            internal helpers            =
========================================*/

func logf(level int, input string, v ...interface{}) {
	if bc.GetLoggingLevel() >= level {
		log.Printf("%s\n", fmt.Sprintf(input, v...))
	}
}

/*=====  End of internal helpers  ======*/
