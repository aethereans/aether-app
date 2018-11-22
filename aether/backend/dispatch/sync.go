// Backend > Routines > Sync
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	// "aether-core/aether/backend/responsegenerator"
	"aether-core/aether/backend/feapiconsumer"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/scheduling"
	tb "aether-core/aether/services/toolbox"
	// "aether-core/aether/services/verify"
	"aether-core/aether/backend/metrics"
	"errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"net"
	// "strconv"
	"strings"
	"time"
)

// Bouncer gate.
// isAllowed attempts to get an outbound lease, and returns a function that can be used to terminate that lease if the lease is granted. The function return is a little bit fancier than I'd normally like, but in this case, it's useful because it captures the values I've calculated within this function.
func isAllowed(a api.Address, reverseConn *net.Conn) (bool, func(wasSuccessful *bool), func()) {
	var loc, subloc string
	var port uint16
	var isReverseConn bool
	if reverseConn != nil {
		host, _, _ := net.SplitHostPort((*reverseConn).RemoteAddr().String())
		loc = host
		subloc = ""
		port = a.Port // untrusted, but so long as it's consistent, doesn't matter.
		isReverseConn = true
	} else {
		loc = string(a.Location)
		subloc = string(a.Sublocation)
		port = a.Port
		isReverseConn = false
	}
	leaseTerminator := func(wasSuccessful *bool) {
		// That *bool instead of bool is important. It needs to accept a pointer so that we can change the value of that pointer later. Otherwise, the initial value of that bool (false) would be captured when defer statement is executed and would never be flipped.
		// logging.Logf(1, "Release lease was called for %v", port)
		globals.BackendTransientConfig.Bouncer.ReleaseOutboundLease(loc, subloc, port, *wasSuccessful, isReverseConn)
		// logging.Logf(1, "Outbound leases after release lease was executed: %s", spew.Sdump(globals.BackendTransientConfig.Bouncer.Outbounds))
	}
	leaseRenewer := func() {
		// This function, when called, renews an existing lease. If called while the lease is active, it will be successful 100% of the time. If called after, it will request a new lease, and it might fail. So consider the system load and possible delayed execution and call this while there are some minutes left before lease expires.

		// Failing to renew the lease marks the connection as failed, but it does not actually stop an ongoing sync.
		globals.BackendTransientConfig.Bouncer.RequestOutboundLease(loc, subloc, port, isReverseConn)
	}
	// spew.Dump(globals.BackendTransientConfig.Bouncer.Outbounds)
	allowed := globals.BackendTransientConfig.Bouncer.RequestOutboundLease(loc, subloc, port, isReverseConn)
	if !allowed {
		logging.Logf(1, "This request to the remote %v:%v was not allowed. This is our outbound list currently open: %#v", loc, port, globals.BackendTransientConfig.Bouncer.Outbounds)
	}
	return allowed, leaseTerminator, leaseRenewer
}

// Sync is the core logic of a single connection. It pulls updates from a remote node and patches it to the current node.
func Sync(a api.Address, lineup []string, reverseConn *net.Conn) error {
	//////////
	// PREP //
	//////////

	// MUTEX
	globals.BackendTransientConfig.SyncActive.Lock()
	defer globals.BackendTransientConfig.SyncActive.Unlock()

	// REVERSE CONN STATUS

	// If reverse connection is present, terminate it at the end of the sync.
	if reverseConn != nil {
		// This is not a very nice thing to do, since this abruptly cuts off the connection. TLS in this connection can error out and say 'use of closed network connection'. Oh well. We just closed the connection, we have no plans to use it anyway. FUTURE: Find a nicer way to do this.
		// Set the deadline to now..
		defer (*reverseConn).SetDeadline(time.Now())
		// Terminate the reverse connection when done.
		defer (*reverseConn).Close()
	}

	// LEASES

	// Request lease
	var syncSuccessful bool
	var syncExited bool

	allowed, releaseLease, renewLease := isAllowed(a, reverseConn)
	if !allowed {
		logging.Logf(1, "Sync() failed to secure an outbound lease. Addr: %#v, isReverseConn: %v", a, reverseConn != nil)
		return errors.New(fmt.Sprintf("Sync() failed to secure an outbound lease. Addr: %#v, isReverseConn: %v", a, reverseConn != nil))
	}

	// Set the defer to release the lease when the sync is done, either via failure or success. (We set syncSuccessful to true when it's successfully completed.)
	defer releaseLease(&syncSuccessful)
	// Set up the outbound lease renewal. Outbound leases expire in 15 minutes, so we'll renew the lease every 10. If an outbound lease expires, the outbound connection is marked and saved as a failure, even if it actually succeeded. So renewing outbound lease as long as we need is good, because at the end we will explicitly terminate the lease (with release lease above)

	// It runs every 10 minutes with an initial 10 minute delay.
	cancelRepeatOutboundLeaseRenewal := scheduling.ScheduleRepeat(
		func() { renewLease() }, 10*time.Minute, time.Duration(10)*time.Minute, &syncExited)

	defer func() {
		cancelRepeatOutboundLeaseRenewal <- true
		logging.Logf(2, "Outbound lease renewal cancelled at the end of the sync.")
	}()
	defer func() { syncExited = true }() // This should always be true at the end. This tells safesleep that it needs to stop sleeping and let schedule be cancelled. This needs to happen after the defer for cancelling repeat outbound renewal, because defers are LIFO.
	///////////////
	// PREP DONE //
	///////////////

	// --------------------
	// Steps
	// - Fetch /status GET to see if the node is online.
	// - Fetch /status POST to see if the node data is valid. Save the Address with the updated last online.
	// - Check if there is a record of the node in the nodes table. If not, create it.
	// - For every entity endpoint, hit the caches that have a later end date than the timestamp for that entity endpoint.
	// - If the node is not static, do posts requets with the update timestamp of that endpoint. The remote will automatically filter the response down to entities that came after the end of the last cache.
	// - At the completion of every endpoint (get + post), save the timestamp.

	logging.Log(2, fmt.Sprintf("SYNC STARTED with node: %s:%d", a.Location, a.Port))
	start := time.Now()
	if reverseConn != nil {
		(*reverseConn).SetDeadline(time.Now().Add(30 * time.Second))
	}
	addr, NODE_STATIC, apiResp, directlyConnectible, err := Check(a, reverseConn)
	if err != nil {
		logging.Logf(1, "Sync errored out. Error: %v", err)
		return err
	}
	if reverseConn != nil {
		(*reverseConn).SetDeadline(time.Now().Add(10 * time.Minute))
	}
	// Establish purgatory. This is where we keep received items that are older than our network head. At the end of the sync, we will take a look at those items and determine if they're ancestor of something that arrived in the sync. If so, we'll insert them as the last step of the sync. If not so, we'll discard them.
	p := Purgatory{}

	// FULLY TRUSTED ADDRESS ENTRY
	// Anything here will be committed in and will write over existing data, since all of this data is either coming from a first-party remote, or from the client.
	addrs := []api.Address{addr}
	if directlyConnectible {
		/*
			DirectlyConnectible means the address is directly accessible. So this is either a sync that we initiated, or is reverse sync request we've accepted and will be doing over a reverse connection open, but still is directly accessible (i.e. the remote is not behind a restrictive NAT that did not accept our port mapping requests.)

			If the node is not directly connectible, there is no point saving the address for this node, because it's guaranteed to be not connectible. They should fix their NAT and make it accept our UPNP port mapping requests.
		*/
		errs := persistence.InsertOrUpdateAddresses(&addrs)
		if len(errs) > 0 {
			err := errors.New(fmt.Sprintf("Some errors were encountered when the Sync attempted InsertOrUpdateAddresses. Sync aborted. Errors: %s", errs))
			logging.Log(1, err)
			abortClr := color.New(color.FgWhite, color.BgRed)
			logging.Log(1, abortClr.Sprintf("SYNC ABORTED. Err: %s", err))
			return err
		}
	}
	var n persistence.DbNode
	var err4 error
	n, err4 = persistence.ReadNode(api.Fingerprint(apiResp.NodeId))
	if err4 != nil && strings.Contains(err4.Error(), "The node you have asked for could not be found") {
		// Node does not exist in the DB. Create it and commit it to DB.
		n.Fingerprint = api.Fingerprint(apiResp.NodeId)
		err5 := persistence.InsertNode(n)
		if err5 != nil {
			// DB commit error, or node was using the same id as ours.
			return err5
		}
	} else if err4 != nil {
		// We have an error in node query and it's not 'node not found'
		return err4
	}
	// Send the connection state to the metrics server.
	firstSync := false
	if n.AddressesLastCheckin == 0 {
		firstSync = true
	} else {
		logging.Logf(2, "This is not a first sync. Addresses Last Checkin: %#v, Node data: %#v", n.AddressesLastCheckin, n)
	}
	metrics.SendConnState(addr, true, firstSync, nil)
	// metrics server end for conn state
	c := startMetricsContainer(apiResp, addr, n, reverseConn != nil)
	openClr := color.New(color.FgWhite, color.BgYellow)
	logging.Log(2, generateStartMessage(c, openClr))
	// For every endpoint, hit the caches. If the node is not static, hit the POSTs too.
	endpoints := map[string]api.Timestamp{
		"boards":      n.BoardsLastCheckin,
		"threads":     n.ThreadsLastCheckin,
		"posts":       n.PostsLastCheckin,
		"votes":       n.VotesLastCheckin,
		"addresses":   n.AddressesLastCheckin,
		"keys":        n.KeysLastCheckin,
		"truststates": n.TruststatesLastCheckin}
	logging.Log(2, fmt.Sprintf("SYNC:PULL STARTED with data from node: %s:%d", a.Location, a.Port))
	logging.Log(2, fmt.Sprintf("Endpoints: %#v", endpoints))
	ims := []persistence.InsertMetrics{}
	// callOrder := []string{"addresses", "votes", "truststates", "posts", "threads", "boards", "keys"}
	callOrder := constructCallOrder(addr, lineup)
	for _, endpointName := range callOrder {
		addrSatiated := false
		logging.Logf(1, "Getting: %s", endpointName)
		if endpointName == "addresses" && !NODE_STATIC {

			//////////////
			//// POST //// (only addresses)
			//////////////

			// We have a special provision for addresses. Unlike others, addresses endpoint scan needs to start from the most recent, and move backwards, because most recent addresses are much more valuable than the older ones. We also have a limit of 100 addresses downloaded at every sync, which means we will stop downloading when we reach the number. If we download all and then pick 100, then we can potentially end up downloading a lot of unwanted data, that would be not that useful.
			// fmt.Println("Addresses endpoint special provision enters.")
			start := time.Now()
			var elapsed time.Duration
			postResp, timeToFirstResponse, err := api.GetPOSTEndpoint(string(a.Location), string(a.Sublocation), a.Port, endpointName, endpoints[endpointName], reverseConn)
			if err != nil {
				logging.Logf(1, "GetPOSTEndpoint inside Sync has errored out. Error: %v", err)
			}
			elapsed = time.Since(start)
			if len(postResp.Addresses) >= 100 {
				// fmt.Println("post address response satiated, won't hit get")
				addrSatiated = true
				postResp.Addresses = postResp.Addresses[0:100]
			}
			c.AddressesPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
			c.AddressesPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
			c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.AddressesPOSTTimeToFirstResponse
			c.AddressesSinglePage = len(postResp.CacheLinks) == 0
			postIface := prepareForBatchInsert(&postResp)
			im, err := persistence.BatchInsert(*postIface)
			if err != nil {
				logging.Logf(1, "Addresses POST BatchInsert inside Sync has errored out. Error: %v", err)
			}
			ims = append(ims, im)
			endpoints[endpointName] = postResp.MostRecentSourceTimestamp
		}
		if addrSatiated && endpointName == "addresses" {
			// We've checked addresses and we're all full for addresses. We won't check the address GET endpoint. Skip addresses GET portion of this sync.
			continue
		}
		start := time.Now()
		if globals.BackendConfig.GetScaledMode() {
			/*
				First check if we're in the scaled mode. If so, skip this part - we'll only sync addresses until we're out of the scaled mode.
				Why?
				Scaled mode means that the node is under so much disk pressure that the event horizon (the threshold of history deletion that can move forwards or backwards in time) has touched the network head, which renders this node one that is not able to provide a full network head to its peers. In the future, in this mode the node will switch to a mode where it only tracks the boards and people followed by its users, but for now, it temporarily stops accepting new content until the network head moves far enough ahead that event horizon can reduce the DB size to under maximum allowable.
				To prepare for that moment, though, we keep updating the addresses tables. Since that table is limited to 1000 addresses, it takes up a constant space.
				(This also appropriately skips setting up the timestamps, so that it won't set timestamps for things that it did not sync.)
			*/
			logging.Logf(1, "This node is in scaled mode, so it's skipping sync with this remote. Remote: %s:%d", a.Location, a.Port)
			continue
		}

		/////////////
		//// GET //// (all)
		/////////////

		// Do an endpoint GET with the timestamp. (Mind that the timestamp is being provided into the GetGETEndpoint, it will only fetch stuff after that timestamp.)
		logging.Log(2, fmt.Sprintf("Asking for entity type: %s", endpointName))
		resp, err6 := api.GetGETEndpoint(string(a.Location), string(a.Sublocation), a.Port, endpointName, endpoints[endpointName], reverseConn)
		if err6 != nil {
			logging.Log(2, fmt.Sprintf("Getting GET Endpoint for the entity type '%s' failed. Error: %s, Address: %#v", endpointName, err6, a))
		}
		logging.Log(3, fmt.Sprintf("Response to be moved to the interface pack: %#v", resp))
		elapsed := time.Since(start) // We end this counter before DB insert starts, because this is the network-time counter.
		// Move the objects into an interface to prepare them to be committed.
		// Address specific
		if len(resp.Addresses) > 100 {
			resp.Addresses = resp.Addresses[0:100]
		}
		p.Filter(&resp) // Filter through purgatory. Older items will be held in purgatory and removed from the resp. At the end of the sync, we'll deal with the items in the purgatory.
		iface := prepareForBatchInsert(&resp)
		// Save the response to the database.
		im, err := persistence.BatchInsert(*iface)
		if err != nil {
			logging.Logf(1, "GET BatchInsert inside Sync has errored out. Entity Type: %v, Error: %v", endpointName, err)
		}
		ims = append(ims, im)
		// Set the last checkin timestamp for each entity type to the beginning of this process. (We will update this later before committing the node checkin set based on the POST response receipts, if any)
		// Check if the apiResp.Timestamp is newer or older than the timestamp we have. It might actually be older,because we might have received a POST response from this node, and that might have been a later Timestamp than that of the last cache's.

		if endpoints[endpointName] < resp.MostRecentSourceTimestamp {
			endpoints[endpointName] = resp.MostRecentSourceTimestamp
		}
		// Insert the metrics into the container.
		if endpointName == "boards" {
			c.BoardsGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "threads" {
			c.ThreadsGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "posts" {
			c.PostsGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "votes" {
			c.VotesGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "keys" {
			c.KeysGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "truststates" {
			c.TruststatesGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		} else if endpointName == "addresses" {
			c.AddressesGETNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
		}
		// GET portion of this sync is done. Now on to POST requests.

		//////////////
		//// POST //// (excl. addresses)
		//////////////

		// POST requests can have two types of responses. If the results of that POST request is few enough, the data might just be provided as a response to the post request directly. Or, if there are many pages of results, the remote saves these into a folder that is available for the next half hour or so, and sends back the link to that folder. The two cases below deal with this.
		if !NODE_STATIC {
			// Because we do the address post request at the top.
			if endpointName == "addresses" {
				continue
			}
			// Generate the POST request.
			// POST request is essentially an ApiResponse converted to JSON. This can have fields like:
			// "filters": [
			//  {"type":"timestamp", "values": ["0", "1483641920"]}
			//  ]
			// which allows us to filter. But if you create an empty request for POST to an entity endpoint, it will give you all the entities for that endpoint since the last cache generation, automatically. There are no filters required for that kind of query.
			start := time.Now()
			var elapsed time.Duration
			postResp, timeToFirstResponse, err := api.GetPOSTEndpoint(string(a.Location), string(a.Sublocation), a.Port, endpointName, endpoints[endpointName], reverseConn)
			elapsed = time.Since(start)
			p.Filter(&postResp)
			postIface := prepareForBatchInsert(&postResp)
			im, err := persistence.BatchInsert(*postIface)
			if err != nil {
				logging.Logf(1, "POST BatchInsert inside Sync has errored out. Entity Type: %v, Error: %v", endpointName, err)
			}
			ims = append(ims, im)
			var singlePage bool
			if len(postResp.CacheLinks) == 0 {
				singlePage = true
			}
			endpoints[endpointName] = postResp.MostRecentSourceTimestamp
			// Insert the metrics into the container.
			if endpointName == "boards" {
				c.BoardsPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.BoardsPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.BoardsPOSTTimeToFirstResponse
				c.BoardsSinglePage = singlePage
			} else if endpointName == "threads" {
				c.ThreadsPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.ThreadsPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.ThreadsPOSTTimeToFirstResponse
				c.ThreadsSinglePage = singlePage
			} else if endpointName == "posts" {
				c.PostsPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.PostsPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.PostsPOSTTimeToFirstResponse
				c.PostsSinglePage = singlePage
			} else if endpointName == "votes" {
				c.VotesPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.VotesPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.VotesPOSTTimeToFirstResponse
				c.VotesSinglePage = singlePage
			} else if endpointName == "keys" {
				c.KeysPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.KeysPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.KeysPOSTTimeToFirstResponse
				c.KeysSinglePage = singlePage
			} else if endpointName == "truststates" {
				c.TruststatesPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.TruststatesPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.TruststatesPOSTTimeToFirstResponse
				c.TruststatesSinglePage = singlePage
			} else if endpointName == "addresses" {
				c.AddressesPOSTNetworkTime = tb.Round(elapsed.Seconds(), 0.1)
				c.AddressesPOSTTimeToFirstResponse = tb.Round(timeToFirstResponse.Seconds(), 0.1)
				c.TotalNetworkRemoteWait = c.TotalNetworkRemoteWait + c.AddressesPOSTTimeToFirstResponse
				c.AddressesSinglePage = singlePage
			}
		}
	}
	// Here, after all the endpoint pulls are complete, we process the purgatory and commit it separately.
	iface := p.Process()
	// Save the response to the database.
	im, err := persistence.BatchInsert(iface)
	if err != nil {
		logging.Logf(1, "Purgatory BatchInsert inside Sync has errored out. Error: %v", err)
	}
	ims = append(ims, im)
	// Purgatory end.
	logging.Log(2, fmt.Sprintf("SYNC:PULL COMPLETE with data from node: %s:%d", a.Location, a.Port))
	// Both POST and GETs are committed into the database. We now need to save the Node LastCheckin timestamps into the database.
	n.BoardsLastCheckin = endpoints["boards"]
	n.ThreadsLastCheckin = endpoints["threads"]
	n.PostsLastCheckin = endpoints["posts"]
	n.VotesLastCheckin = endpoints["votes"]
	n.AddressesLastCheckin = endpoints["addresses"]
	n.KeysLastCheckin = endpoints["keys"]
	n.TruststatesLastCheckin = endpoints["truststates"]
	err9 := persistence.InsertNode(n)
	if err9 != nil {
		return err9
	}
	if directlyConnectible {
		addrs[0].LastSuccessfulPing = api.Timestamp(time.Now().Unix())
		addrs[0].LastSuccessfulSync = api.Timestamp(time.Now().Unix())
		errs2 := persistence.InsertOrUpdateAddresses(&addrs)
		if len(errs2) > 0 {
			err := errors.New(fmt.Sprintf("Some errors were encountered when the Sync attempted InsertOrUpdateAddresses. Sync aborted. Errors: %s", errs2))
			logging.Log(1, err)
			abortClr := color.New(color.FgWhite, color.BgRed)
			logging.Log(1, abortClr.Sprintf("SYNC ABORTED. Err: %s", err))
			return err
		}
	}
	logging.Log(2, "Inserted the last successful sync stamp at the end of the sync.")
	logging.Log(2, fmt.Sprintf("SYNC COMPLETE with node: %s:%d. It took %d seconds", a.Location, a.Port, int(time.Since(start).Seconds())))
	closeClr := color.New(color.FgBlack, color.BgWhite)
	logging.Log(1, generateCloseMessage(c, closeClr, &ims, int(time.Since(start).Seconds()), true))
	// Send the connection state to the metrics server.
	metrics.SendConnState(addr, false, firstSync, &ims)
	// Insert the appropriate markers to the config
	switch addr.Type {
	case 2:
		globals.BackendConfig.SetLastLiveAddressConnectionTimestamp(time.Now().Unix())
	case 3, 254:
		globals.BackendConfig.SetLastBootstrapAddressConnectionTimestamp(time.Now().Unix())
	case 255:
		globals.BackendConfig.SetLastStaticAddressConnectionTimestamp(time.Now().Unix())
	}
	syncSuccessful = true

	/*----------  Send sync metrics to frontend  ----------*/
	feapiconsumer.BackendAmbientStatus.LastOutboundDurationSeconds = int32(time.Since(start).Seconds())
	feapiconsumer.BackendAmbientStatus.LastOutboundConnTimestamp = time.Now().Unix()
	feapiconsumer.SendBackendAmbientStatus()
	/*----------  And all done!  ----------*/

	return nil
}

/*
//////////
Internal functions
//////////
*/

func constructCallOrder(remote api.Address, lineup []string) []string {
	// All mim nodes support addresses to enable proper protocol function.
	supported := []string{"addresses"}
	// supported := []string{"addresses", "addresses"}
	// supported := []string{"boards"}
	// return supported
	if len(lineup) == 0 {
		// If not specified, all entities are allowed.
		//FUTURE: This needs to read from a central somewhere else — otherwise when we add a new entity, we're going to forget adding it here and it's gonna be a lot of unnecessary pain to find it out.
		lineup = []string{"vote", "truststate", "post", "thread", "board", "key"}
	}
	availableSubprots := remote.Protocol.Subprotocols
	for _, val := range availableSubprots {
		if val.Name == "c0" {
			if tb.IndexOf(tb.Singular("votes"), val.SupportedEntities) != -1 &&
				tb.IndexOf("vote", lineup) != -1 {
				supported = append(supported, "votes")
			}
			if tb.IndexOf(tb.Singular("truststates"), val.SupportedEntities) != -1 &&
				tb.IndexOf("truststate", lineup) != -1 {
				supported = append(supported, "truststates")
			}
			if tb.IndexOf(tb.Singular("posts"), val.SupportedEntities) != -1 &&
				tb.IndexOf("post", lineup) != -1 {
				supported = append(supported, "posts")
			}
			if tb.IndexOf(tb.Singular("threads"), val.SupportedEntities) != -1 &&
				tb.IndexOf("thread", lineup) != -1 {
				supported = append(supported, "threads")
			}
			if tb.IndexOf(tb.Singular("boards"), val.SupportedEntities) != -1 &&
				tb.IndexOf("board", lineup) != -1 {
				supported = append(supported, "boards")
			}
			if tb.IndexOf(tb.Singular("keys"), val.SupportedEntities) != -1 &&
				tb.IndexOf("key", lineup) != -1 {
				supported = append(supported, "keys")
			}
		}
	}
	return supported
}

// prepareForBatchInsert verifies the items in this response container, and converts it to the correct form BatchInsert accepts.
func prepareForBatchInsert(r *api.Response) *[]interface{} {
	resp := *r
	var carrier []interface{}
	for i, _ := range resp.Boards {
		carrier = append(carrier, resp.Boards[i])
	}
	for i, _ := range resp.Threads {
		carrier = append(carrier, resp.Threads[i])
	}
	for i, _ := range resp.Posts {
		carrier = append(carrier, resp.Posts[i])
	}
	for i, _ := range resp.Votes {
		carrier = append(carrier, resp.Votes[i])
	}
	for i, _ := range resp.Keys {
		carrier = append(carrier, resp.Keys[i])
	}
	for i, _ := range resp.Truststates {
		carrier = append(carrier, resp.Truststates[i])
	}
	for i, _ := range resp.Addresses {
		carrier = append(carrier, resp.Addresses[i])
	}
	return &carrier
}
