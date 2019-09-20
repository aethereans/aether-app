package cmd

import (
	"aether-core/aether/backend/beapiserver"
	"aether-core/aether/backend/dispatch"
	"aether-core/aether/backend/feapiconsumer"
	"aether-core/aether/backend/responsegenerator"
	"aether-core/aether/backend/server"
	// "aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/ports"
	"aether-core/aether/services/scheduling"
	"aether-core/aether/services/upnp"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	var loggingLevel int
	var backendAPIPort int
	var backendAPIPublic bool
	var adminFeAddr string
	var adminFePk string
	var imprint bool
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Global logging level of the app.")
	cmdRun.Flags().IntVarP(&backendAPIPort, "backendapiport", "", 0, "Sets the port that the backend will attempt to serve the backend API output from. If this port is occupied, it will pick another, therefore it's not safe to assume that this will be the actual backend API port.")
	cmdRun.Flags().BoolVarP(&backendAPIPublic, "backendapipublic", "", false, "If you set this to true, your node will expose the backend api port to the public internet, as well. If not, it will be only served locally. Defaults to false. The reason you might want this is to put the backend on a VPS and make your frontend connect to it, so that it can stay online 24/7.")
	cmdRun.Flags().StringVarP(&adminFeAddr, "adminfeaddr", "", "127.0.0.1:45001", "Spawner FE Address is the address of the frontend that spawns this backend instance. The backend will reach out to this address to tell that it is ready at which port.")
	cmdRun.Flags().StringVarP(&adminFePk, "adminfepk", "", "", "Spawner FE Public Key is the public key of the frontend instance that is spawning the backend process. This is useful to give, because if admin needs to change (ex: when you want to monitor the status of the backend from a different machine than you've installed) you can move your FE config to the new machine, run the FE from the new machine and it will update the admin FE address because it can authenticate with the key.")
	cmdRun.Flags().BoolVarP(&imprint, "imprint", "", false, "Makes the app create the user directory, a default user config, and quit. This is useful when you want to make changes to the default config before the app starts running.")
	cmdRoot.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Start a full-fledged Mim node that tracks the network and responds to requests.",
	Long: `Start a full-fledged Mim node that tracks the network and responds to requests. This is the default entry point if you want to use this Mim backend normally.

This will do three main things:

- Actively start fetching from other nodes and constructing the network head in the local machine
- Expose a local API for the frontend app to peruse, so that the content fetched over the network is available for consumption
- Expose an API to the external world that serves the data this computer has under the rules set by the Mim protocol.
`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := EstablishConfigs(cmd)
		showIntro() // This isn't first because it needs configs to show app version.
		if flags.imprint.value.(bool) == true {
			logging.Logf(1, "The configuration file is created at: %v/backend/backend_config.json. \nYou should make your changes, and run the binary without the flag to start running.", globals.BackendConfig.GetUserDirectory())
			fmt.Printf("The configuration file is created at: %v/backend/backend_config.json. \nYou should make your changes, and run the binary without the flag to start running.\n", globals.BackendConfig.GetUserDirectory())
			os.Exit(0)
		}
		persistence.CreateDatabase()
		persistence.CheckDatabaseReady()
		startSchedules()
		handleTemporaryConfigUpdates()
		gotValidPort := make(chan bool)
		go beapiserver.StartBackendServer(gotValidPort)
		<-gotValidPort // Only proceed after this is true.
		feapiconsumer.SendBackendReady()
		collectAmbientStatusData()
		feapiconsumer.SendBackendAmbientStatus()
		go server.StartMimServer()
		/*----------  Signal handling  ----------*/
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		go func() {
			sig := <-sigs
			logging.Logf(1, "The backend received a signal. Signal: %v", sig)
			shutdown()
		}()
		<-done
	},
}

func handleTemporaryConfigUpdates() {
	// TODO FUTURE Remove this in dev.12 - this is to set this for one time from the prior default, if the user has not changed the default.
	if globals.BackendConfig.GetMaxInboundConns() == 5 {
		globals.BackendConfig.SetMaxInboundConns(10)
	}
	if globals.BackendConfig.GetDispatchExclusionExpiryForLiveAddress() == 3 {
		globals.BackendConfig.SetDispatchExclusionExpiryForLiveAddress(5)
	}
}

// collectAmbientStatusData gathers information from different pieces of the system to gather existing data, to send it at boot.
/*
	I know you're gonna think this is a neat function, can we put this on a schedule and make it update the client? No, you can't do that, because the function below actually resets the values to defaults. So if you actually end up running this, it will override existing values for stuff like status. So you can accidentally say the database is available while it is inserting, etc.
*/
func collectAmbientStatusData() {
	/*----------  Network, In, outbounds  ----------*/

	/*Last outbound connection duration and last outbound conn timestamp is set in sync library */
	feapiconsumer.BackendAmbientStatus.LastInboundConnTimestamp = globals.BackendTransientConfig.Bouncer.GetLastInboundSyncTimestamp(false)
	feapiconsumer.BackendAmbientStatus.InboundsCount15 = int32(len(globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(15, true)))

	feapiconsumer.BackendAmbientStatus.OutboundsCount15 = int32(len(globals.BackendTransientConfig.Bouncer.GetOutboundsInLastXMinutes(15, true)))
	/*----------  Network misc  ----------*/
	/* UPNP status is set in the UPNP library */
	// feapiconsumer.BackendAmbientStatus.UPNPStatus = "Idle"
	feapiconsumer.BackendAmbientStatus.LocalNodeExternalIP = globals.BackendConfig.GetExternalIp() // todo
	feapiconsumer.BackendAmbientStatus.LocalNodeExternalPort = int32(globals.BackendConfig.GetExternalPort())
	/*----------  Database  ----------*/
	/*Db Status, LastInsertDurationSeconds and LastDbInsertTimestamp is set in writer at batchInsert */
	feapiconsumer.BackendAmbientStatus.DatabaseStatus = "Available"
	feapiconsumer.BackendAmbientStatus.DbSizeMb = int64(globals.GetDbSize())
	feapiconsumer.BackendAmbientStatus.MaxDbSizeMb = int64(globals.BackendConfig.GetMaxDbSizeMb())
	/*----------  Caching  ----------*/
	feapiconsumer.BackendAmbientStatus.CachingStatus = "Idle"
	feapiconsumer.BackendAmbientStatus.LastCacheGenerationTimestamp = globals.BackendConfig.GetLastCacheGenerationTimestamp()
}

func startSchedules() {
	// logging.Logf(1, "UserDir: %v", globals.BackendConfig.GetUserDirectory())
	logging.Log(1, "Setting up cyclical tasks is starting.")
	defer logging.Log(1, "Setting up cyclical tasks is complete.")
	/*
		Ordered by initial delay:
		Verify external port: T+0 	(immediately)
		Neighbourhood dispatcher 			T+0: 	(immediately)

		UPNP Port mapper: 		T+0 	(immediately)
		Explorer dispatcher 			T+10:
		Address Scanner: 			T+15m
		Cache generator: 			T+30m
	*/
	// Before doing anything, you need to validate the external port. This function takes a second or so, and it needs to block the runtime execution because if two routines call it separately, it causes a race condition. After the first initialisation, however, this function becomes safe for concurrent use.
	ports.VerifyBackendPorts()
	// UPNP tries to port map every 10 minutes.

	globals.BackendTransientConfig.StopUPNPCycle = scheduling.ScheduleRepeat(func() { upnp.MapPort() }, 10*time.Minute, time.Duration(0), nil)
	dispatch.Bootstrap() // This will run only if needed.
	// The dispatcher that seeks live nodes runs every minute.

	globals.BackendTransientConfig.StopNeighbourhoodCycle = scheduling.ScheduleRepeat(func() { dispatch.NeighbourWatch() }, 1*time.Minute, time.Duration(0), nil)

	globals.BackendTransientConfig.StopExplorerCycle = scheduling.ScheduleRepeat(func() { dispatch.Explore() }, 10*time.Minute, time.Duration(10)*time.Minute, nil)

	if !globals.BackendConfig.GetPreventOutboundReverseRequests() {
		// We start the inbound connection watch (which makes outbound reverse connection requests) only if the outbound reverse requests are not disabled.
		globals.BackendTransientConfig.StopInboundConnectionCycle = scheduling.ScheduleRepeat(func() { dispatch.InboundConnectionWatch() }, 1*time.Minute, time.Duration(5)*time.Minute, nil)
	}

	// Address scanner goes through all prior unconnected addresses and attempts to connect to them to establish a relationship. It starts 30 minutes after a node is started, so that the node will actually have a chance to collect some addresses to check.
	// globals.BackendTransientConfig.StopAddressScannerCycle = scheduling.ScheduleRepeat(func() { dispatch.AddressScanner() }, 2*time.Hour, time.Duration(15)*time.Minute, nil)

	globals.BackendTransientConfig.StopNetworkScanCycle = scheduling.ScheduleRepeat(func() { dispatch.DoNetworkScan() }, 10*time.Minute, time.Duration(0)*time.Minute, nil)

	// Attempt cache generation every hour, but it will be pre-empted if the last cache generation is less than 23 hours old, and if the node is not tracking the head. So that this will run effectively every day, only.

	globals.BackendTransientConfig.StopCacheGenerationCycle = scheduling.ScheduleRepeat(func() { responsegenerator.GenerateCaches() }, 5*time.Minute, time.Duration(5)*time.Minute, nil)

	globals.BackendTransientConfig.StopBadlistRefreshCycle = scheduling.ScheduleRepeat(func() { configstore.BadlistInstance.Refresh() }, 6*time.Hour, time.Duration(0)*time.Minute, nil)

}

func shutdown() {
	logging.Log(1, "Shutdown initiated. Entering lameduck mode and stopping all scheduled tasks and routines. This will take a couple seconds.")
	// Initiate lameduck mode. This will start declining and inbound and outbound requests, as well as reverse connections requests. Ongoing database actions can still be processed.
	globals.BackendTransientConfig.LameduckInitiated = true
	// logging.Log(1, "Waiting a second to let all network i/o close gracefully...")
	// time.Sleep(time.Duration(1) * time.Second)
	// Initiate shutdown. At this point, if anything is still being written into the database, they will attempt to exit gracefully.
	// Flip the delay terminators in the case they haven't started actually running yet.

	globals.BackendTransientConfig.ShutdownInitiated = true
	// Stop routines
	globals.BackendTransientConfig.StopNeighbourhoodCycle <- true
	logging.Logf(1, "StopNeighbourhoodCycle is done.")
	globals.BackendTransientConfig.StopExplorerCycle <- true
	logging.Logf(1, "StopExplorerCycle is done.")
	globals.BackendTransientConfig.StopInboundConnectionCycle <- true
	logging.Logf(1, "StopInboundConnectionCycle is done.")
	// globals.BackendTransientConfig.StopAddressScannerCycle <- true
	// logging.Logf(1, "StopAddressScannerCycle is done.")
	globals.BackendTransientConfig.StopNetworkScanCycle <- true
	logging.Logf(1, "StopNetworkScanCycle is done.")
	globals.BackendTransientConfig.StopUPNPCycle <- true
	logging.Logf(1, "StopUPNPCycle is done.")
	globals.BackendTransientConfig.StopCacheGenerationCycle <- true
	logging.Logf(1, "StopCacheGenerationCycle is done.")
	globals.BackendTransientConfig.StopBadlistRefreshCycle <- true
	logging.Logf(1, "StopBadlistRefreshCycle is done.")
	// logging.Logf(1, "Inbounds: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.Inbounds))
	// logging.Logf(1, "Outbounds: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.Outbounds))
	// logging.Logf(1, "InboundHistory: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.InboundHistory))
	// logging.Logf(1, "OutboundHistory: %s\n", spew.Sdump(globals.BackendTransientConfig.Bouncer.OutboundHistory))
	// logging.Logf(1, "Last Inbound Sync Timestamp: %v", globals.BackendTransientConfig.Bouncer.GetLastInboundSyncTimestamp(false))
	// logging.Logf(1, "Last Successful Outbound Sync Timestamp: %v", globals.BackendTransientConfig.Bouncer.GetLastOutboundSyncTimestamp(true))
	// logging.Logf(1, "Inbounds in the last 5 minutes: %v", len(globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(5)))
	// logging.Logf(1, "Successful outbounds in the last 5 minutes: %v", len(globals.BackendTransientConfig.Bouncer.GetOutboundsInLastXMinutes(5, true)))
	// logging.Logf(2, "Nonces: %v", globals.BackendTransientConfig.Nonces)

	logging.Log(1, "Waiting a second to let DB close gracefully...")
	time.Sleep(time.Duration(1) * time.Second) // Wait 5 seconds to let DB tasks complete.
	// And after that, we shut down the database.
	globals.DbInstance.Close()
	defer func() {
		// The functions that access DB can panic after the DB is closed. But after DB is closed, we don't care - the DB is out of harm's way and the only state that remains at this phase is the transient state, and that's going to be wiped out a few nanoseconds later. Recover from any panics.
		recResult := recover()
		if recResult != nil {
			logging.Logf(1, "Recovered from a panic at the end of the shutdown after DB close. A panic here can be caused by a process being interrupted. In most cases, it's normal behaviour and nothing to worry about. Panic'd error: %#v", recResult)
		}
	}()
	// We delete at shutdown and at boot, just in case deletion at shutdown didn't work.
	globals.BackendTransientConfig.POSTResponseRepo.DeleteAllFromDisk()
	logging.Log(1, "Shutdown is complete. Bye.")
	os.Exit(0)
}
