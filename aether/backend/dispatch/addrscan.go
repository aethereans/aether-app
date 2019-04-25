// Backend > Dispatch > AddrScan
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/aether/io/api"
	pers "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "aether-core/aether/services/safesleep"
	"errors"
	"fmt"
	// "github.com/pkg/errors"
	// "strings"
	"aether-core/aether/services/toolbox"
	// "net"
	// "github.com/davecgh/go-spew/spew"
	"sync"
	"time"
)

func getAllAddresses(isDesc bool) ([]api.Address, error) {
	searchType := ""
	if isDesc {
		searchType = "all_desc"
	} else {
		searchType = "all_asc"
	}
	resp, err := pers.ReadAddresses("", "", 0, 0, 0, 0, 0, 0, searchType)
	if err != nil {
		errors.New(fmt.Sprintf("getAllAddresses in AddressScanner failed.", err))
	}
	return resp, nil
}

func filterByLastSuccessfulPing(addrs []api.Address, scanStart api.Timestamp) []api.Address {
	live := []api.Address{}
	cutoff := api.Timestamp(time.Unix(int64(scanStart), 0).Add(-2 * time.Minute).Unix())
	// Cutoff is 2 minutes before the threshold, because our pinger accepts a node whose last successful ping was within 2 minutes as online.
	for key, _ := range addrs {
		if addrs[key].LastSuccessfulPing >= cutoff {
			live = append(live, addrs[key])
		}
	}
	return live
}

func filterByAddressType(addrType uint8, addrs []api.Address) ([]api.Address, []api.Address) {
	filteredAddrs := []api.Address{}
	remainder := []api.Address{}
	for key, _ := range addrs {
		if addrs[key].Type == addrType {
			filteredAddrs = append(filteredAddrs, addrs[key])
		} else {
			remainder = append(remainder, addrs[key])
		}
	}
	return filteredAddrs, remainder
}

func removeAddr(addr api.Address, addrs []api.Address) []api.Address {
	for i := len(addrs) - 1; i >= 0; i-- {
		if addr.Location == addrs[i].Location &&
			addr.Sublocation == addrs[i].Sublocation &&
			addr.Port == addrs[i].Port {
			addrs = append(addrs[0:i], addrs[i+1:len(addrs)]...)
		}
	}
	return addrs
}

func updateAddrs(addrs []api.Address) ([]api.Address, error) {
	updatedAddrs := Pinger(addrs)
	err := pers.AddrTrustedInsert(&updatedAddrs)
	if err != nil {
		return []api.Address{}, errors.New(fmt.Sprintf("updateAddrs encountered an error in AddrTrustedInsert.", err))
	}
	return updatedAddrs, nil
}

// var (
// 	lastNetworkScan int64
// )

/*
	AddrType:
		-1: we don't care
		0, 1, 2 ...: the usual address type

	Reqtype:
		0, -1: we don't care
		-2: Attempt giving unconnected, but if not available, anything works.
*/
// func findOnlineNodes(count int, reqType, addrType int, excl *[]api.Address, requestNetworkScan bool) ([]api.Address, error) {
// 	start := api.Timestamp(time.Now().Unix())
// 	addrs, err := getAllAddresses(true) // desc - last synced first primary, last pinged first secondary sort
// 	if err != nil {
// 		return []api.Address{}, errors.New(fmt.Sprintf("findOnlineNodes: getAllAddresses within this function failed.", err))
// 	}
// 	// logging.Logf(1, "All addresses: %s", )
// 	// logging.LogObj(2, "All addresses", Dbg_convertAddrSliceToNameSlice(addrs))
// 	if addrType > -1 {
// 		addrs, _ = filterByAddressType(uint8(addrType), addrs)
// 	}
// 	// logging.LogObj(2, "Filtered addresses", Dbg_convertAddrSliceToNameSlice(addrs))
// 	if excl != nil {
// 		for _, addr := range *excl {
// 			addrs = removeAddr(addr, addrs)
// 		}
// 	}
// 	updatedAddrs := []api.Address{}

// 	/*============================================
// 	=            Network scan request            =
// 	============================================*/
// 	/*
// 		Determine if we can run a network scan or not.
// 		- If not requested, we cannot.
// 		- If we ran one via scout in the last 10 mins, we cannot.
// 		- Else, we can.

// 		Why?

// 		Network scans are expensive - you're hitting between 100 and 1000 last known addresses for their status, potentially. Normally, it'll stop once you find one that works, but it's done in batches of  100, so minimum number of attempts before it returns is 100. We want to throttle this, doing it every 10 minutes or so at most.
// 	*/
// 	doNetworkScan := false
// 	if requestNetworkScan {
// 		now := time.Now().Unix()
// 		if lastNetworkScan < toolbox.CnvToCutoffMinutes(10) {
// 			// If last scan is older than 10 minutes ago, we can do it.
// 			doNetworkScan = true
// 			lastNetworkScan = now
// 		}
// 	}
// 	/*=====  End of Network scan request  ======*/

// 	if doNetworkScan {
// 		logging.Logf(1, "Doing a network scan within findOnlineNodes.")
// 		var addrUpdateErr error
// 		updatedAddrs, addrUpdateErr = updateAddrs(addrs)
// 		if addrUpdateErr != nil {
// 			logging.Logf(1, "findOnlineNodes: updateAddress within this function failed.", addrUpdateErr)
// 		}
// 	} else {
// 		logging.Logf(1, "Not doing a network scan within findOnlineNodes.")
// 		updatedAddrs = addrs
// 	}
// 	// logging.Logf(2, "Updated addresses: %s", Dbg_convertAddrSliceToNameSlice(updatedAddrs))
// 	liveNodes := []api.Address{}
// 	if doNetworkScan {
// 		logging.Logf(2, "Live nodes count before filter by last successful ping: %v: live nodes: %v", len(updatedAddrs), updatedAddrs)
// 		liveNodes = filterByLastSuccessfulPing(updatedAddrs, start)
// 		// ^ This is gated behind network scan, because it checks for anything within the last 2 minutes, i.e. that got flagged as live in that scan. If you don't do the scan, then this will always return 0.
// 		logging.Logf(2, "Live nodes count: %v, live nodes: %v", len(liveNodes), liveNodes)
// 	} else {
// 		liveNodes = updatedAddrs
// 	}
// 	// logging.Logf(2, "Live addresses: %s", Dbg_convertAddrSliceToNameSlice(updatedAddrs))
// 	if reqType == -2 {
// 		// logging.Logf(1, "Live nodes are these. Live nodes: %s", Dbg_convertAddrSliceToNameSlice(liveNodes))
// 		logging.Log(1, "ReqType = -2, we are looking for nonconnected addrs.")
// 		nonconnecteds, connecteds := pickUnconnectedAddrs(liveNodes)
// 		// logging.Logf(1, "nonconnecteds: %v", nonconnected)
// 		// if len(nonconnected) > 0 {
// 		// 	// logging.Logf(1, "ReqType = -2, we found some nonconnected onlines. Let's pull from those first. Found: %s", Dbg_convertAddrSliceToNameSlice(nonconnected))
// 		// 	liveNodes = nonconnected
// 		// }
// 		liveNodes = append(nonconnecteds, connecteds...)
// 	}
// 	// logging.Logf(2, "Nonconnected live nodes count: %v, nodes: %v", len(liveNodes), liveNodes)
// 	if len(liveNodes) == 0 { // If zero, bail.
// 		return liveNodes, errors.New("This database has no addresses online.")
// 	}

// 	/*=====================================================================
// 	=            Enforce exclusions for too-recently-connected            =
// 	=====================================================================*/
// 	l := []api.Address{}
// 	for k, _ := range liveNodes {
// 		if !dpe.IsExcluded(liveNodes[k]) {
// 			l = append(l, liveNodes[k])
// 		}
// 	}
// 	liveNodes = l
// 	/*=====  End of Enforce exclusions for too-recently-connected  ======*/
// 	if count > len(liveNodes) || count == 0 {
// 		count = len(liveNodes)
// 	}
// 	/*
// 		Using random vs top x depends on the network scan type, because if you don't do a network scan, what you have at this stage is the whole 1000 item address database that wasn't checked for. Picking off of that database at random is not going to give you anything
// 	*/
// 	if doNetworkScan {
// 		logging.Logf(2, "This database has %v addresses online and not excluded: %v", len(liveNodes), liveNodes)
// 		rands := toolbox.GetInsecureRands(len(liveNodes), count)
// 		logging.Logf(2, "This is the random number generated for address selection: %v", rands)
// 		selected := []api.Address{}
// 		for _, val := range rands {
// 			selected = append(selected, (liveNodes)[val])
// 		}
// 		logging.Logf(2, "This is the selected address to connect to: %v", selected)
// 		return selected, nil
// 	} else {
// 		logging.Logf(2, "This database returned %v addresses as potentially online without a scan: %v", len(liveNodes), liveNodes)
// 		return liveNodes[0:count], nil
// 	}
// }

func findOnlineNodesV2(count int, reqType, addrType int, excl *[]api.Address, reverse bool) ([]api.Address, error) {
	logging.Logf(1, "Attempting network scan within findOnlineNodes")
	DoNetworkScan()
	// ^ This will block until the mutex is released, and then when it's released, you can continue with fresh data.
	logging.Logf(1, "Network scan complete within findOnlineNodes.")
	addrs, err := getAllAddresses(true) // desc - last synced first primary, last pinged first secondary sort
	if err != nil {
		return []api.Address{}, errors.New(fmt.Sprintf("findOnlineNodes: getAllAddresses within this function failed.", err))
	}
	if addrType > -1 {
		addrs, _ = filterByAddressType(uint8(addrType), addrs)
	}
	if excl != nil {
		for _, addr := range *excl {
			addrs = removeAddr(addr, addrs)
		}
	}
	liveNodes := addrs
	if reqType == -2 {
		logging.Log(1, "ReqType = -2, we are looking for nonconnected addrs.")
		nonconnecteds, connecteds := pickUnconnectedAddrs(liveNodes)

		liveNodes = append(nonconnecteds, connecteds...)
	}
	// logging.Logf(2, "Nonconnected live nodes count: %v, nodes: %v", len(liveNodes), liveNodes)
	if len(liveNodes) == 0 { // If zero, bail.
		return liveNodes, errors.New("This database has no addresses online.")
	}
	/*=====================================================================
	=            Enforce exclusions for too-recently-connected            =
	=====================================================================*/
	l := []api.Address{}
	if !reverse {
		// If this is a non-reverse (normal) find online request, we use the reverse dispatcher exclusion queue.
		for k, _ := range liveNodes {
			if !dpe.IsExcluded(liveNodes[k]) {
				l = append(l, liveNodes[k])
			}
		}
	} else {
		// If this is a reverse find online request, we use the reverse dispatcher exclusion queue.
		for k, _ := range liveNodes {
			if !reverseDpe.IsExcluded(liveNodes[k]) {
				l = append(l, liveNodes[k])
			}
		}
	}
	liveNodes = l
	/*=====  End of Enforce exclusions for too-recently-connected  ======*/
	if count > len(liveNodes) || count == 0 {
		count = len(liveNodes)
	}
	logging.Logf(3, "This database returned %v addresses as potentially online without a scan: %v", len(liveNodes), liveNodes)
	return liveNodes[0:count], nil
}

func pickUnconnectedAddrs(addrs []api.Address) ([]api.Address, []api.Address) {
	nonconnecteds := []api.Address{}
	connecteds := []api.Address{}
	for key, _ := range addrs {
		if addrs[key].LastSuccessfulSync == 0 {
			nonconnecteds = append(nonconnecteds, addrs[key])
		} else {
			connecteds = append(connecteds, addrs[key])
		}
	}
	return nonconnecteds, connecteds
}

func RefreshAddresses() error {
	addrs, err := getAllAddresses(false) // asc - the oldest unconnected first
	if err != nil {
		return errors.New(fmt.Sprintf("RefreshAddresses: getAllAddresses within this function failed.", err))
	}
	updateAddrs(addrs)
	return nil
}

// func GetOnlineAddresses(
// 	noOfOnlineAddressesRequested int,
// 	exclude []api.Address,
// 	addressType uint8,
// 	forceUnconnected bool,
// ) (
// 	[]api.Address, error,
// ) {
// 	ln, err := findOnlineNodes(noOfOnlineAddressesRequested, 0, int(addressType), &exclude, true)
// 	// spew.Dump(ln)
// 	return ln, err
// }

func AddressScanner() {
	globals.BackendTransientConfig.AddressesScannerActive.Lock()
	defer globals.BackendTransientConfig.AddressesScannerActive.Unlock()
	err := RefreshAddresses()
	if err != nil {
		logging.Log(1, fmt.Sprintf("AddressScanner failed. Error: %v", err))
		// return errors.Wrap(err, "AddressScanner failed.")
	}
}

var (
	lastNetworkScan     int64
	networkScannerMutex sync.Mutex
)

func DoNetworkScan() {
	networkScannerMutex.Lock()
	defer networkScannerMutex.Unlock()
	if lastNetworkScan > toolbox.CnvToCutoffMinutes(10) {
		logging.Logf(1, "We've done a network scan less than 10 minutes ago. Skipping.")
		return
	}
	defer func() { lastNetworkScan = time.Now().Unix() }()
	addrs, err := getAllAddresses(true) // desc - last synced first primary, last pinged first secondary sort
	if err != nil {
		logging.Logf(1, "DoNetworkScan: getAllAddresses within this function failed.", err)
		return
	}
	var addrUpdateErr error
	_, addrUpdateErr = updateAddrs(addrs)
	if addrUpdateErr != nil {
		logging.Logf(1, "findOnlineNodes: updateAddress within this function failed.", addrUpdateErr)
		return
	}
}

func GetUnconnAddr(count int, excl *[]api.Address, doNetworkScan bool) []api.Address {
	addrs, err := findOnlineNodesV2(count, -2, -1, excl, false)
	// fmt.Println(len(addrs))
	if err != nil {
		logging.Log(1, fmt.Sprintf("Unconnected address search failed. Error: %v", err))
		return []api.Address{}
	}
	return addrs
}

// func Dbg_convertAddrSliceToNameSlice(nodes []api.Address) []string {
// 	names := []string{}
// 	for _, val := range nodes {
// 		if val.Client.ClientName != "" { // If this is not a completely nonconnected node with no data
// 			names = append(names, val.Client.ClientName)
// 		}
// 	}
// 	return names
// }
