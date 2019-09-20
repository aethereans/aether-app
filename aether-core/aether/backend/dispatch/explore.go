// Backend > Routines > Explore
// This file contains the explore routine in dispatch. Explore is the routine that helps us discover new nodes, and it also makes sure that we occasionally check our static nodes and bootstrappers as well.

package dispatch

import (
	"aether-core/aether/io/persistence"
	// "aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
)

var (
	ticker = 0
)

// Explore is a function that reaches further into the network than our neighbourhood watch. At every call, it finds a live regular remote that is new to us and does a sync, and at the end it inserts it into our neighbourhood. At every 6 ticks, it syncs with our statics, and at every 36 ticks, it syncs with bootstrappers. An explore tick is not always a minute, it depends on the schedule, but it is something around 10 minutes. So that means refreshing statics every hour, and hitting bootstrappers every 6 hours.
func Explore() {
	if ticker%36 == 0 && ticker != 0 {
		ticker = 0

		////////////////////////////////////////////
		// per 36 ticks
		// 3 live bootstrap, 1 static bootstrap node
		// all live CAs, all static CAs
		////////////////////////////////////////////

		// call 3 live and 1 static bootstrap nodes and sync with them.
		liveBs, err := persistence.ReadAddresses("", "", 0, 0, 0, 3, 0, 3, "limit")
		if err != nil {
			logging.Logf(1, "There was an error when we tried to read live bootstrapper addresses for Explore schedule. Error: %#v", err)
		}
		staticBs, err2 := persistence.ReadAddresses("", "", 0, 0, 0, 1, 0, 254, "limit")
		if err2 != nil {
			logging.Logf(1, "There was an error when we tried to read static bootstrapper addresses for Explore schedule. Error: %#v", err2)
		}
		bsAddrs := append(liveBs, staticBs...)
		for key, _ := range bsAddrs {
			Sync(bsAddrs[key], []string{}, nil)
		}
		// call all CA nodes and sync with them. These should be fairly short. We are not limiting them to x number of CAs because each CA will likely have their own data only. (We terminate the connection without sync if it's a CA that we do not trust.)
		liveCA, err3 := persistence.ReadAddresses("", "", 0, 0, 0, 0, 0, 4, "limit")
		if err3 != nil {
			logging.Logf(1, "There was an error when we tried to read live CA addresses for Explore schedule. Error: %#v", err3)
		}
		staticCA, err4 := persistence.ReadAddresses("", "", 0, 0, 0, 0, 0, 253, "limit")
		if err4 != nil {
			logging.Logf(1, "There was an error when we tried to read static CA addresses for Explore schedule. Error: %#v", err4)
		}
		caAddrs := append(liveCA, staticCA...)
		for key, _ := range caAddrs {
			Sync(caAddrs[key], []string{}, nil)
		}
	} else if ticker%6 == 0 && ticker != 0 {

		////////////////////////////////////////////
		// per 6 ticks
		// all static nodes
		////////////////////////////////////////////

		// go through all statics to see if there are any updates.
		statics, err := persistence.ReadAddresses("", "", 0, 0, 0, 4, 0, 255, "limit") // get all static nodes we know of.
		if err != nil {
			logging.Logf(1, "There was an error when we tried to read static addresses for Explore schedule. Error: %#v", err)
		}
		for key, _ := range statics {
			Sync(statics[key], []string{}, nil)
		}
	} else {
		////////////////////////////////////////////
		// per every tick
		// all live nodes
		////////////////////////////////////////////
		// find a new node that we haven't synced before, and sync with it.
		Scout()
	}
	ticker++
	// globals.BackendTransientConfig.ExplorerTick++
}
