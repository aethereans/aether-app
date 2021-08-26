// Backend > Dispatch > Ping
// This file is the subsystem that decides on which remotes to connect to.

package dispatch

import (
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "aether-core/aether/services/safesleep"
	// "errors"
	"fmt"
	// "strings"
	// tb "aether-core/aether/services/toolbox"
	// "net"
	"time"
)

// Pinger goes through the list of extant nodes, and pings the status endpoints to see if they are online. If no response is provided in X seconds, the node is offline. It returns a set of online nodes.
// We need to do this in batches of 100. Otherwise we end up with "socket: too many open files" error.
func Pinger(fullAddressesSlice []api.Address) []api.Address {
	// Paginate addresses first. We batch these into pages of 100, because it's very easy to run into too many open files error if you just dump it through.
	logging.Log(2, fmt.Sprintf("Pinger is called for this number of addresses: %d", len(fullAddressesSlice)))
	var pages [][]api.Address
	dataSet := fullAddressesSlice
	pgSize := globals.BackendConfig.GetPingerPageSize()
	numPages := len(dataSet)/pgSize + 1
	var allUpdatedAddresses []api.Address
	// The division above is floored.
	for i := 0; i < numPages; i++ {
		beg := i * pgSize
		var end int
		// This is to protect from 'slice bounds out of range'
		if (i+1)*pgSize > len(dataSet) {
			end = len(dataSet)
		} else {
			end = (i + 1) * pgSize
		}
		pageData := dataSet[beg:end]
		var page []api.Address
		page = pageData
		pages = append(pages, page)
	}
	// For every page,
	for i, _ := range pages {
		// If there's a shutdown in progress, break and exit.
		if globals.BackendTransientConfig.ShutdownInitiated {
			return []api.Address{}
		}
		logging.Log(2, fmt.Sprintf("This pinger page has this many addresses to ping: %d", len(pages[i])))
		// Run the core logic.
		addrs := pages[i]
		outputChan := make(chan api.Address)
		for j, _ := range addrs {
			// Check if shutdown was initiated.
			if globals.BackendTransientConfig.ShutdownInitiated {
				break // Stop processing and return
			}
			logging.Log(3, fmt.Sprintf("Pinging the address at %v:%d", addrs[j].Location, addrs[j].Port))
			go Ping(addrs[j], outputChan)
		}
		var updatedAddresses []api.Address
		// We will receive as many addresses as answers. Every time something is put into a channel, this will fire, if the channel is empty, it will block.
		for i := 0; i < len(addrs); i++ {
			var a api.Address
			a = <-outputChan
			updatedAddresses = append(updatedAddresses, a)
		}
		allUpdatedAddresses = append(allUpdatedAddresses, updatedAddresses...)
	}
	// Clean blanks.
	logging.Log(2, fmt.Sprintf("All updated addresses count (this should be the same as goroutine count: %d", len(allUpdatedAddresses)))
	var cleanedAllUpdatedAddresses []api.Address
	for i, _ := range allUpdatedAddresses {
		if allUpdatedAddresses[i].Location != "" {
			// The location is not blank. This is an actual updated address.
			cleanedAllUpdatedAddresses = append(cleanedAllUpdatedAddresses, allUpdatedAddresses[i])
		}
	}
	logging.Log(2, fmt.Sprintf("Cleaned addresses count (this should be the same as the online addresses count: %d", len(cleanedAllUpdatedAddresses)))

	return cleanedAllUpdatedAddresses
}

// Ping runs a Check and returns the result. If there is an error, it returns a blank address.
func Ping(addr api.Address, processedAddresses chan<- api.Address) {
	logging.Log(3, fmt.Sprintf("Connection attempt started: %v:%v", addr.Location, addr.Port))
	var blankAddr api.Address
	if addr.LastSuccessfulPing > api.Timestamp(time.Now().Add(time.Duration(-2)*time.Minute).Unix()) {
		// If it's been less than 2 minutes since we last pinged this address. We'll just pass this ping to not create excessive traffic.
		logging.Logf(2, "We pinged this address already in the last 2 minutes. Skipping this ping and using the past result. Address: %s/%s:%d", addr.Location, addr.Sublocation, addr.Port)
		// Mark it timestamped as now, and send it back.
		// addr.LastSuccessfulPing = api.Timestamp(time.Now().Unix())
		processedAddresses <- addr
		return
	}
	updatedAddr, _, _, _, err := Check(addr, nil, "ping")
	if err != nil {
		updatedAddr = blankAddr
		logging.Log(3, err)
	}
	processedAddresses <- updatedAddr
}
