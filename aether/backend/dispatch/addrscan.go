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

func filterByType(addrType int, addrs []api.Address) ([]api.Address, []api.Address) {
	if addrType <= -1 {
		return addrs, []api.Address{}
	}
	filteredAddrs := []api.Address{}
	remainder := []api.Address{}
	for key, _ := range addrs {
		if addrs[key].Type == uint8(addrType) {
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

// if count == 0, we do a full-range search and return all live nodes.
// addrType == -1 : Give me ANY nodes
// addrType == -2 : Attempt to give me the nodes I've not synced before. If not, anything works.
// Why the negative numbers? because anything 0<= is used as a real value in addrType.
func findOnlineNodes(count int, addrType int, excl *[]api.Address) ([]api.Address, error) {
	start := api.Timestamp(time.Now().Unix())
	addrs, err := getAllAddresses(true) // desc - last synced first
	if err != nil {
		return []api.Address{}, errors.New(fmt.Sprintf("findOnlineNodes: getAllAddresses within this function failed.", err))
	}
	// logging.Logf(1, "All addresses: %s", )
	// logging.LogObj(2, "All addresses", Dbg_convertAddrSliceToNameSlice(addrs))
	addrs, _ = filterByType(addrType, addrs)
	// logging.LogObj(2, "Filtered addresses", Dbg_convertAddrSliceToNameSlice(addrs))
	if excl != nil {
		for _, addr := range *excl {
			addrs = removeAddr(addr, addrs)
		}
	}
	updatedAddrs, err := updateAddrs(addrs)
	// logging.Logf(2, "Updated addresses: %s", Dbg_convertAddrSliceToNameSlice(updatedAddrs))
	if err != nil {
		errors.New(fmt.Sprintf("findOnlineNodes: updateAddress within this function failed.", err))
	}
	logging.Logf(2, "Live nodes count before filter by last successful ping: %v: live nodes: %v", len(updatedAddrs), updatedAddrs)
	liveNodes := filterByLastSuccessfulPing(updatedAddrs, start)
	// logging.Logf(2, "Live addresses: %s", Dbg_convertAddrSliceToNameSlice(updatedAddrs))
	if count == 0 { // count == 0: return everything found.
		return liveNodes, nil
	}
	logging.Logf(2, "Live nodes count: %v, live nodes: %v", len(liveNodes), liveNodes)
	if addrType == -2 {
		// logging.Logf(1, "Live nodes are these. Live nodes: %s", Dbg_convertAddrSliceToNameSlice(liveNodes))
		logging.Log(1, "AddrType = -2, we are looking for nonconnected addrs.")
		nonconnected := pickUnconnectedAddrs(liveNodes)
		// logging.Logf(1, "nonconnecteds: %v", nonconnected)
		if len(nonconnected) != 0 {
			// logging.Logf(1, "AddrType = -2, we found some nonconnected onlines. Let's pull from those first. Found: %s", Dbg_convertAddrSliceToNameSlice(nonconnected))
			liveNodes = nonconnected
		}
	}
	logging.Logf(2, "Nonconnected live nodes count: %v, nodes: %v", len(liveNodes), liveNodes)
	if len(liveNodes) == 0 { // If zero, bail.
		return liveNodes, errors.New("This database has no addresses online.")
	}
	logging.Logf(2, "This database has %v addresses online: %v", len(liveNodes), liveNodes)
	rands := toolbox.GetInsecureRands(len(liveNodes), count)
	logging.Logf(2, "This is the random number generated for address selection: %v", rands)
	selected := []api.Address{}
	for _, val := range rands {
		selected = append(selected, (liveNodes)[val])
	}
	logging.Logf(2, "This is the selected address to connect to: %v", selected)
	return selected, nil
}

func pickUnconnectedAddrs(addrs []api.Address) []api.Address {
	nonconnecteds := []api.Address{}
	for key, _ := range addrs {
		if addrs[key].LastSuccessfulSync == 0 {
			nonconnecteds = append(nonconnecteds, addrs[key])
		}
	}
	return nonconnecteds
}

func RefreshAddresses() error {
	addrs, err := getAllAddresses(false) // asc - the oldest unconnected first
	if err != nil {
		return errors.New(fmt.Sprintf("RefreshAddresses: getAllAddresses within this function failed.", err))
	}
	updateAddrs(addrs)
	return nil
}

func GetOnlineAddresses(
	noOfOnlineAddressesRequested int,
	exclude []api.Address,
	addressType uint8,
	forceUnconnected bool,
) (
	[]api.Address, error,
) {
	ln, err := findOnlineNodes(noOfOnlineAddressesRequested, int(addressType), &exclude)
	// spew.Dump(ln)
	return ln, err
}

func AddressScanner() {
	globals.BackendTransientConfig.AddressesScannerActive.Lock()
	defer globals.BackendTransientConfig.AddressesScannerActive.Unlock()
	err := RefreshAddresses()
	if err != nil {
		logging.Log(1, fmt.Sprintf("AddressScanner failed. Error: %v", err))
		// return errors.Wrap(err, "AddressScanner failed.")
	}
}

func GetUnconnAddr(count int, excl *[]api.Address) []api.Address {
	addrs, err := findOnlineNodes(count, -2, excl)
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
