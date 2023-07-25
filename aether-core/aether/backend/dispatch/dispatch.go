// Backend > Dispatch
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"

	// "aether-core/aether/services/safesleep"
	"errors"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	// "github.com/pkg/errors"
	// "aether-core/aether/services/toolbox"
	"strings"
	// "time"
	// "net"
)

/*
Dispatcher is the big thing here.
One thing to keep thinking about, this behaviour of the dispatch to get one online node that is not excluded, might actually create 'islands' that only connect to each other.
To be able to diagnose this, I might need to build a tool that visualises the connections between the nodes.. Just to make sure that there are no islands.
*/

var (
	scoutAttempts = 50
	// How many times scout will try to do a sync.
)

// NeighbourWatch keeps in sync with all our neighbours.
func NeighbourWatch() {
	logging.Log(2, "NeighbourWatch triggers.")
	/*=======================================
	=            Bootstrap check            =
	=======================================*/
	/*
		Check if bootstrap needs to run. If so, that's what needs to happen before anything else.

		This is important because the first bootstrap can fail. If that happens, we want to be able to continue trying without waiting another 6 hours for the next bootstrap tick.

		This function checks whether enough time has passed since the last successful bootstrap. If it hasn't it just bails without doing anything.
	*/
	Bootstrap()
	/*=====  End of Bootstrap check  ======*/

	loc, subloc, port := globals.BackendTransientConfig.NeighboursList.Pop()
	a := api.Address{
		Location:    api.Location(loc),
		Sublocation: api.Location(subloc),
		Port:        port,
	}
	// The pop was blank
	if isBlank(a) {
		Scout()
		return
	}
	// The pop was not blank...
	err := connect(a)
	// ...but it failed
	if err != nil {
		// Run the scout and it'll try to connect to something up to 10 times.
		scoutErr := Scout()
		if scoutErr != nil {
			logging.Logf(1, "NeighbourWatch: Connect failed. Error: %#v", scoutErr)
			// log and bail.
			return
		}
	}
	// if err != nil {
	// 	// if it failed because the remote was too busy,
	// 	if strings.Contains(err.Error(), "Received status code: 429") {
	// 		// add the failed node as exclusion, and attempt to find another.
	// 		excl := []api.Address{a}
	// 		Scout(&excl, false)
	// 		return
	// 	}
	// 	// if it failed because for any other reason,
	// 	logging.Logf(1, "NeighbourWatch: Connect failed. Error: %#v", err)
	// 	// log and bail.
	// 	return
	// }
	// our connect succeeded, all is well. return.
	// return
}

/*
Scout finds nodes that are online and that we have not connected to, and connects to them - thus adding them to the neighbourhood.

If we receive an error that is not 'too busy', we stop and wait for the next cycle. If we get a too busy error, we try to find a new node to connect to up to maxNeighbourWatchRetriesIfRemoteTooBusy times.

Rationale #1 is that if you get a too busy response, it means that the response returned fairly quick. If you get something else, you might have already spent considerable time in sync, and it might be better to just let go and wait for the next pop.

Rationale #2 is that if a sync failure happened, except in the case of 'remote too busy', it has a chance that it might be because of the local machine being under too much stress, such as DB writes timing out. In most cases, waiting a bit to reduce the pressure might help. In any case, breathlessly retrying again and again will definitely won't.

So we retry when we 100% know it's because of the remote (remote too busy case), but we let goo and let the next tick attempt a new sync. A tick is usually around a minute, so it's a good way to just wait a bit, give a little breathing room to the node, and start fresh.
*/
func Scout() error {
	logging.Log(2, "Scout triggers.")
	addrs, err := findOnlineNodesV2(0, -1, -1, nil, false)
	if err != nil {
		errStr := fmt.Sprintf("Scout: address search failed. Error: %v", err)
		logging.Log(1, errStr)
		return errors.New(errStr)
	}
	// addrs := GetUnconnAddr(100, nil, true)
	if len(addrs) == 0 {
		logging.Log(2, "Scout got no unconnected addresses. Bailing.")
		return errors.New("Scout got no unconnected addresses. Bailing.")
	}
	attempts := 0
	for k := range addrs {
		attempts++
		if attempts > scoutAttempts {
			break
		}
		logging.Logf(1, "Scout: Connection attempt #%v.", attempts)
		err := connect(addrs[k])
		if err != nil {
			logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
			continue
		}
		return nil
		// if err == nil {
		// 	return nil
		// }
		// // Failed
		// logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
	}
	// for i := 0; i < scoutAttempts; i++ {
	// 	logging.Logf(1, "Scout is attempting to connect to a node. Attempt #%v", i+1)
	// 	if i+1 > len(addrs) {
	// 		errText := "Scout doesn't have any more addresses to try. Bailing."
	// 		logging.Log(2, errText)
	// 		return errors.New(errText)
	// 	}
	// 	err := connect(addrs[i])
	// 	// if connect succeeded, all is well.
	// 	if err == nil {
	// 		return nil
	// 	}
	// 	// Failed
	// 	logging.Logf(1, "Scout: Connect failed. Error: %#v", err)
	// }
	// If we've come here without returning, scout tried its best, but could not connect to anything.
	allFailedError := fmt.Errorf("Scout failed because all %v nodes we have tried has failed.", len(addrs))
	logging.Logf(1, "Scout: Connect failed. Error: %#v", allFailedError)
	return fmt.Errorf("Scout: Connect failed. Error: %#v", allFailedError)
}

// InboundConnectionWatch takes a look at how many inbound connections we have received in the past 3 minutes. If the number is zero, it triggers a reverse connection open request to a node.
func InboundConnectionWatch() {
	nt := globals.BackendConfig.GetNodeType()
	if nt != 2 {
		// If not a live node, we don't request reverse opens.
		return
	}
	logging.Log(2, "Inbound connection watch triggers.")
	pastInboundConns := globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(3, true)
	if len(pastInboundConns) < 2 || globals.BackendTransientConfig.NewContentCommitted {
		// // Request reverse connect to a node we think is online.
		// online, err := online, err := findOnlineNodesV2(1, -1, -1, nil, true)
		// if err != nil {
		// 	logging.Logf(1, "Find online nodes for InboundConnectionWatch failed. Error: %v", err)
		// }
		//  We do not reverse open to addresses that are URLs. Why? Because they might be behind a HTTP proxy, and the reverse proxy request is not a HTTP request. A node with an URL is probably pretty established as a service - and receiving a lot of requests, it likely doesn't appreciate being asked to connect into some random node on the Internet.

		// 	But beyond that, in the essence, a URL is 'established', very likely to have some sort of reverse proxy like NGINX in the front. Raw Mim is TCP, not HTTP - so under normal circumstances, a HTTP reverse proxy will not pass it through. Even if it passed it through, the response from the responding server is going to be an outbound TLS handshake, which is going to a) seriously confuse a reverse proxy, b) probably break a lot of things in the meanwhile.

		// 	It's good that RawMim exists, because vast, vast majority of the nodes (99.99..%) won't ever be behind a URL. Then this works normally. For those who are, we disable this, which is beneficial to them, since they won't have to respond to reverse open requests.
		// onlineWithoutURLAddresses := []api.Address{}
		// for k, _ := range online {
		// 	if online[k].LocationType != 3 {
		// 		// If not URL
		// 		onlineWithoutURLAddresses = append(onlineWithoutURLAddresses, online[k])
		// 	}
		// }
		// if len(onlineWithoutURLAddresses) > 0 {
		// 	api.RequestInboundSync(string(onlineWithoutURLAddresses[0].Location),
		// 		string(onlineWithoutURLAddresses[0].Sublocation),
		// 		onlineWithoutURLAddresses[0].Port)
		// }
		err := ReverseScout()
		if err != nil {
			logging.Logf(1, "InboundConnectionWatch: ReverseScout failed with the error: %v", err)
			return
		}
		// If no error, we're successful and we can set the content as delivered, at least to this remote.
		globals.BackendTransientConfig.NewContentCommitted = false
	}
}

/*
//////////
Internal functions
//////////
*/

func isBlank(a api.Address) bool {
	return len(a.Location) == 0 &&
		len(a.Sublocation) == 0 &&
		a.Port == 0
}

func connect(a api.Address) error {
	// sync
	if dpe.IsExcluded(a) {
		errText := fmt.Sprintf("Connect failed. We've connected to this address way too recently. We're skipping this connect. Address: %#v", a)
		logging.Logf(2, errText)
		return errors.New(errText)
	}
	maxAttempts := 2
	attempts := 0
	var err error
	for {
		/*
			We try X times to get an outbound lease. If it still doesn't work, we bail. In practice, the second attempt should always work, because second attempt will immediately queue up a sync and it'll get stuck on the mutex. Go releases mutex to the earliest call, so it should always be this one - since the other competing thing for the outbounds is the reversescout, and that does not ever retry getting a lease.
		*/
		attempts++
		if attempts > maxAttempts {
			errText := fmt.Sprintf("Sync failed because we couldn't get an outbound lease after %v tries. Address: %#v, Error: %#v", maxAttempts, a, err)
			logging.Log(2, errText)
			return errors.New(errText)
		}
		err = Sync(a, []string{}, nil)
		if err != nil {
			if strings.Contains(err.Error(), "Failed to secure an outbound lease") {
				// If sync failed because it couldn't get a lease, we try again until it does.
				logging.Logf(1, "Failed to secure an outbound lease. This was attempt #%v", attempts)
				continue
			}
			// If some other error
			dpe.Add(a)
			// ^ We now exclude when it fails too.
			errText := fmt.Sprintf("Sync failed. Address: %#v, Error: %#v", a, err)
			logging.Log(2, errText)
			return errors.New(errText)
		}
		// err is nil
		break
	}
	// Add to exclusions for a while
	dpe.Add(a)
	globals.BackendTransientConfig.NeighboursList.Push(string(a.Location), string(a.Sublocation), a.Port)
	return nil
}

// sameAddress checks if the addresses given are the same
func sameAddress(a1 *api.Address, a2 *api.Address) bool {
	if a1.Location == a2.Location && a1.Sublocation == a2.Sublocation && a1.Port == a2.Port {
		return true
	}
	return false
}

// addrsInGivenSlice checks if the address is extant in a given slice.
func addrsInGivenSlice(addr *api.Address, slc *[]api.Address) bool {
	address := *addr
	slice := *slc
	for i := range slice {
		if sameAddress(&address, &slice[i]) {
			return true
		}
	}
	return false
}
