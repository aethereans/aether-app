package cmd

import (
	"aether-core/aether/backend/dispatch"
	// "aether-core/aether/backend/metrics"
	"aether-core/aether/backend/server"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/create"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/scheduling"

	// "github.com/fatih/color"
	// "aether-core/aether/services/ports"
	"encoding/json"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	var orgName string
	var appName string
	var loggingLevel int
	var port int
	var externalIp string
	var bootstrapIp string
	var bootstrapPort int
	var bootstrapType int
	var syncAndQuit bool
	var printToStdout bool
	var metricsDebugMode bool
	var swarmPlan string
	var killTimeout int
	var swarmNodeId int
	var fpCheckEnabled bool
	var sigCheckEnabled bool
	var powCheckEnabled bool
	var pageSigCheckEnabled bool
	var tlsEnabled bool
	var allowLocalhostRemotes bool
	cmdOrchestrate.Flags().StringVarP(&orgName, "orgname", "", "Air Labs", "Global transient org name for the app.")
	cmdOrchestrate.Flags().StringVarP(&appName, "appname", "", "Aether", "Global transient app name for the app.")
	cmdOrchestrate.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Global logging level of the app.")
	cmdOrchestrate.Flags().IntVarP(&port, "port", "", 49999, "The port the external world can use to communicate with this node.")
	cmdOrchestrate.Flags().StringVarP(&externalIp, "externalip", "", "0.0.0.0", "The IP address the external world can use to communicate with this node.")
	cmdOrchestrate.Flags().StringVarP(&bootstrapIp, "bootstrapip", "", "127.0.0.1", "The bootstrap node's ip address that will be inserted into the local database at the start. If you provide this, you also have to provide port and type.")
	cmdOrchestrate.Flags().IntVarP(&bootstrapPort, "bootstrapport", "", 51000, "The bootstrap node's port that will be inserted into the local database at the start. If you provide this, you also have to provide ip address and type.")
	cmdOrchestrate.Flags().IntVarP(&bootstrapType, "bootstraptype", "", 255, "The bootstrap node's type (static or live, 2 or 255) that will be inserted into the local database at the start. If you provide this, you also have to provide ip address and port.")
	cmdOrchestrate.Flags().BoolVarP(&syncAndQuit, "syncandquit", "", false, "The only thing that happens is that we connect to the bootstrap node, ingest everything, and then quit. No ongoing processing.")
	cmdOrchestrate.Flags().BoolVarP(&printToStdout, "printtostdout", "", false, "Route log output to stdout. This will make the logging library write to stdout instead of log. This is useful in the case of orchestrate if you want to see logs from individual backend nodes.")
	cmdOrchestrate.Flags().BoolVarP(&metricsDebugMode, "metricsdebugmode", "", false, "Enable sending of debug-mode metrics. These metrics are designed to provide a debuggable view of what the node is doing.")
	cmdOrchestrate.Flags().StringVarP(&swarmPlan, "swarmplan", "", "", "This flag allows you to load a swarm plan to your swarm nodes. This swarm plan does have a list of TO-FROM node connections with certain delays, so that you can schedule connections to happen in a fashion that is pre-mediated. This allows you to kickstart a few node connections and see how network behaves based on new data, for example.")
	cmdOrchestrate.Flags().IntVarP(&killTimeout, "killtimeout", "", 120, "If given, this sets up a maximum lifetime in seconds for the node. This is useful in swarm tests in which all swarm nodes have to exit so that the test can move on to the data analysis stage. ")
	cmdOrchestrate.Flags().IntVarP(&swarmNodeId, "swarmnodeid", "", 0, "This flag allows you to set the terminal output colour for the orchestrator. This is useful, mostly because you might be running a swarm test in a single machine, outputting into the same terminal, and you want to be able to discern which messages are coming from one node or the other.")
	cmdOrchestrate.Flags().BoolVarP(&fpCheckEnabled, "fpcheckenabled", "", true, "Setting this to false will disable fingerprint checks on entities. True by default.")
	cmdOrchestrate.Flags().BoolVarP(&sigCheckEnabled, "sigcheckenabled", "", true, "Setting this to false will disable signature checks on entities. True by default.")
	cmdOrchestrate.Flags().BoolVarP(&powCheckEnabled, "powcheckenabled", "", true, "Setting this to false will disable signature checks on entities. True by default.")
	cmdOrchestrate.Flags().BoolVarP(&pageSigCheckEnabled, "pagesigcheckenabled", "", true, "Setting this to false will disable page signature checks on pages. True by default.")
	cmdOrchestrate.Flags().BoolVarP(&tlsEnabled, "tlsenabled", "", true, "Setting this to false will disable the TLS encryption layer in peer to peer connections. This is for debug purposes only, mainnet nodes will refuse to connect to or accept connections from any remote with TLS disabled.")
	cmdOrchestrate.Flags().BoolVarP(&allowLocalhostRemotes, "allowlocalhostremotes", "", false, "Setting this to true will allow localhost remotes to be saved into the database. This is useful for swarm testing.")
	cmdRoot.AddCommand(cmdOrchestrate)
}

var cmdOrchestrate = &cobra.Command{
	Use:   "orchestrate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		flags := EstablishConfigs(cmd)
		// Prep the database
		persistence.CreateDatabase()
		persistence.CheckDatabaseReady()
		showIntro()
		// startup()
		// do things here
		// If all of them are changed, or if none of them are changed, we're good. If SOME of them are changed, we crash.
		if !((flags.bootstrapIp.changed && flags.bootstrapPort.changed && flags.bootstrapType.changed) || (!flags.bootstrapIp.changed && !flags.bootstrapPort.changed && !flags.bootstrapType.changed)) {
			logging.LogCrash(fmt.Sprintf("Please provide either all of the parts of a bootstrap node (IP:Port/Type) or none. You've provided:%s:%d, Type: %d",
				flags.bootstrapIp.value.(string),
				flags.bootstrapPort.value.(int),
				flags.bootstrapType.value.(int)))
		}
		if flags.syncAndQuit.value.(bool) &&
			flags.bootstrapIp.changed &&
			flags.bootstrapPort.changed {
			addr := constructCallAddress(api.Location(flags.bootstrapIp.value.(string)), flags.bootstrapPort.value.(int), flags.bootstrapType.value.(int))
			addrs := []api.Address{addr}
			errs := persistence.InsertOrUpdateAddresses(&addrs)
			if len(errs) > 0 {
				logging.LogCrash(fmt.Sprintf("These errors were encountered on InsertOrUpdateAddress attempt: %s", errs))
			}
			// First, verify external port, so that our metrics will report the right port.
			// We just want to connect, pull and quit.
			err := dispatch.Sync(addr, []string{}, nil)
			if err != nil {
				logging.LogCrash(err)
			}
			// dispatch.Bootstrap()
			logging.Log(2, fmt.Sprintf("We've gotten everything in the node %#v.", addr))
		} else {
			startSchedules()
			if flags.swarmPlan.changed {
				scheduleSwarmPlan(flags.swarmPlan.value.(string))
			}
			if flags.killTimeout.changed {
				logging.Log(1, fmt.Sprintf("This node set to shut down in %d seconds.", flags.killTimeout.value.(int)))
				scheduling.ScheduleOnce(func() {
					shutdown()
				}, time.Duration(flags.killTimeout.value.(int))*time.Second)
			}
			server.StartMimServer()
		}
	},
}

// Orchestrate endpoint will allow us to first generate random data, then pull that into the local database.

func constructCallAddress(ip api.Location, port int, addrtype int) api.Address {
	subprots := []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
	addr, err := create.CreateAddress(ip, "", 4, uint16(port), uint8(addrtype), 1, 1, 1, 0, subprots, 2, 0, 0, "Aether", "")
	if err != nil {
		logging.LogCrash(err)
	}
	addr.SetVerified(true)
	return addr
}

func parsePlansForThisNode(planAsByte []byte) []PlanCommand {
	// This is just JSON parsing without a backing struct. The other alternative was copying over the struct (I don't want to have a swarmtest dependency here, I can't import from there), so this is arguably cleaner.
	// var f interface{}
	var plancmd []PlanCommand

	err2 := json.Unmarshal(planAsByte, &plancmd)
	if err2 != nil {
		logging.LogCrash(fmt.Sprintf("The swarm plan JSON parsing failed. Error: %s", err2))
	}
	// sch := f.([]interface{})
	// var plansForThisNode []map[string]interface{}
	var plansForThisNode []PlanCommand
	for _, val := range plancmd {
		// valMapped := val.(map[string]interface{})
		if val.FromNodeAppName == globals.BackendTransientConfig.AppIdentifier {
			plansForThisNode = append(plansForThisNode, val)
		}
	}
	return plansForThisNode
}

// scheduleSwarmPlan finds and reads through the swarm plan json file, and determines which schedules apply to it. For those which apply, it inserts into the scheduler logic.
func scheduleSwarmPlan(planloc string) {
	// Read the file
	planAsByte, err := ioutil.ReadFile(planloc)
	if err != nil {
		logging.LogCrash(fmt.Sprintf("The swarm plan document could not be read. Error: %s", err))
	}
	// Parse the plans relevant to this specific node based on the AppIdentifier
	plans := parsePlansForThisNode(planAsByte)
	// Insert the plans into the scheduler.
	for _, plan := range plans {
		if plan.CommandName == "connect" {
			logging.Log(1, fmt.Sprintf("This node is going to attempt to connect to the address: %s:%d in %v", plan.ToIp, plan.ToPort, plan.TriggerAfter))
			if plan.Force {
				scheduling.ScheduleOnce(selectCmdFunc("connect_force", plan), plan.TriggerAfter)
			} else {
				scheduling.ScheduleOnce(selectCmdFunc("connect_nonforce", plan), plan.TriggerAfter)
			}
		} else if plan.CommandName == "cachegen" {
			logging.Log(1, fmt.Sprintf("This node is going to generate caches in %v", plan.TriggerAfter))
			scheduling.ScheduleOnce(selectCmdFunc("cachegen", plan), plan.TriggerAfter)
		} else if plan.CommandName == "reverseopen" {
			logging.Logf(1, "This node will attempt to request inbound sync from remote: %s:%v", plan.ToIp, plan.ToPort)
			scheduling.ScheduleOnce(selectCmdFunc("reverseopen", plan), plan.TriggerAfter)
		}

	}
}
