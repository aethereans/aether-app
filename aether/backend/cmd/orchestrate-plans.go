// Cmd > Orchestrate > OrchestratePlans
// This is a library of plans that the swarm nodes can receive from the swarm coordinator for testing.

package cmd

import (
	"aether-core/aether/backend/dispatch"
	"aether-core/aether/backend/responsegenerator"
	// "aether-core/aether/backend/metrics"
	// "aether-core/aether/backend/server"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	// "aether-core/aether/services/create"
	// "aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "aether-core/aether/services/scheduling"
	"github.com/fatih/color"
	// "aether-core/aether/services/ports"
	// "encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/spf13/cobra"
	// "io/ioutil"
	"time"
)

type PlanCommand struct {
	CommandName     string
	FromNodeAppName string
	ToIp            string
	ToPort          int
	ToType          int
	TriggerAfter    time.Duration
	Force           bool
	// Force: false: insert into the addresses DB and let node discover connectivity itself
	// Force: true: directly trigger sync()
}

func selectCmdFunc(str string, plan PlanCommand) func() {
	// Inserts the given address to the database, so that the node can discover the address itself. The sync will typically begin when the next cycle of the cycle hits.
	nonforcedInsert := func() {
		clr := color.New(color.FgYellow)
		logging.Log(1, clr.Sprintf("Injecting the address: %s:%d into the database now.", plan.ToIp, plan.ToPort))
		// color.Cyan(fmt.Sprintf("Injecting the address: %s:%d into the database now.", ip, port))
		addr := constructCallAddress(api.Location(plan.ToIp), plan.ToPort, plan.ToType)
		addrs := []api.Address{addr}
		errs := persistence.InsertOrUpdateAddresses(&addrs)
		if len(errs) > 0 {
			logging.LogCrash(fmt.Sprintf("These errors were encountered on InsertOrUpdateAddress attempt: %s", errs))
		}
	}
	// Triggers a sync with the given address directly, so that the sync begins immediately.
	forcedInsert := func() {
		clr := color.New(color.FgYellow)
		logging.Log(1, clr.Sprintf("Starting the sync with the address: %s:%d now. (Forced)", plan.ToIp, plan.ToPort))
		addr := constructCallAddress(api.Location(plan.ToIp), plan.ToPort, plan.ToType)
		dispatch.Sync(addr, []string{}, nil)
	}
	generateCaches := func() {
		clr := color.New(color.FgYellow)
		logging.Log(1, clr.Sprintf("Starting to generate caches now."))
		// trigger cachegen.
		responsegenerator.GenerateCaches()
	}
	requestReverseOpen := func() {
		clr := color.New(color.FgYellow)
		logging.Log(1, clr.Sprintf("Starting the reverse open sync request with the address: %s:%d now.", plan.ToIp, plan.ToPort))
		api.RequestInboundSync(plan.ToIp, "", uint16(plan.ToPort))
	}
	if str == "connect_force" {
		return forcedInsert
	} else if str == "connect_nonforce" {
		return nonforcedInsert
	} else if str == "cachegen" {
		return generateCaches
	} else if str == "reverseopen" {
		return requestReverseOpen
	} else {
		return func() {}
	}
}
