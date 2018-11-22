package fecmd

import (
	// "aether-core/aether/frontend/beapiconsumer"
	"aether-core/aether/frontend/besupervisor"
	// "aether-core/aether/frontend/clapiconsumer"
	"aether-core/aether/frontend/feapiserver"
	// "aether-core/aether/protos/clapi"
	"aether-core/aether/frontend/festructs"
	// "aether-core/aether/frontend/inflights"
	"aether-core/aether/frontend/kvstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/ports"
	"aether-core/aether/services/scheduling"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"aether-core/aether/frontend/refresher"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	var loggingLevel int
	var clientIp string
	var clientPort int
	var isDev bool
	cmdRun.Flags().IntVarP(&loggingLevel, "logginglevel", "", 0, "Sets the frontend logging level.")
	cmdRun.Flags().StringVarP(&clientIp, "clientip", "", "127.0.0.1", "This is the IP of the client that is starting the frontend instance. THis is almost always 127.0.0.1 since clients and frontends almost always live in the same computer.")
	cmdRun.Flags().IntVarP(&clientPort, "clientport", "", 0, "The port of the client instance starting the frontend. Frontend will call back at this endpoint via GRPC and confirm it's ready.")
	cmdRun.Flags().BoolVarP(&isDev, "isdev", "", false, "If you set this to true, the frontend will be compiled from scratch and it will compile the backend it uses from scratch. This is good for development use, but it will only work within the checked out repo.")
	cmdRoot.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Start an Aether frontend instance that maintains a compiled data source for the frontend to use..",
	Long: `This starts a Aether frontend process. This is the main process that communicates and responds to client requests. This is where all the different pieces of data from the network actually comes together and ends up as boards, threads, posts, users and so on.
`,
	Run: func(cmd *cobra.Command, args []string) {
		EstablishConfigs(cmd)
		// Start frontend kvstore
		kvstore.OpenKVStore()
		defer kvstore.CloseKVStore()
		kvstore.CheckKVStoreReady()
		// Start notifications subsystem
		festructs.InstantiateNotificationsSingleton()
		// start frontend server
		gotValidPort := make(chan bool)
		go feapiserver.StartFrontendServer(gotValidPort)
		<-gotValidPort // Only proceed after this is true.
		go besupervisor.StartLocalBackend()
		for globals.FrontendTransientConfig.BackendReady != true {
			// Block until the backend tells the frontend via gRPC that it is ready.
			time.Sleep(time.Millisecond * 100)
			// ^ If you don't have sleep here, this will block and stop execution forever.
		}
		startSchedules()
		// feapiserver.SendAmbients(false)
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		go func() {
			sig := <-sigs
			logging.Logf(1, "Frontend received a signal: %v", sig)
			handleShutdown()
		}()
		<-done
	},
}

/*
sigs := make(chan os.Signal, 1)
done := make(chan bool, 1)
signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
go func() {
	sig := <-sigs
	logging.Logf(1, "Frontend received a signal: %v", sig)
	logging.Logf(1, "Frontend shutdown routine initiated.")
	globals.FrontendTransientConfig.ShutdownInitia
	kvstore.CloseKVStore()
	logging.Logf(1, "Sending a SIGTERM to the backend...")
	err := besupervisor.BackendDaemon.Process.Signal(syscall.SIGINT)
	if err != nil {
		logging.Logf(1, "Backend process kill returned an error. Err: %v", err)
	}
	logging.Logf(1, "Done sending a SIGTERM to the backend.")
	logging.Logf(1, "Frontend is shut down.")
	time.Sleep(5 * time.Second)
}()
<-done


*/

func handleShutdown() {
	/*
		Heads up, this works when the app is packaged, but not when it is not.

		Why? Because when the app is running in dev, it's using go run ... to compile what you have on the fly. Which means the parent process is "go" which is the compiler, not your process, and *that* process will also receive signals. Even if you don't handle the signals, that process will, and kill the app.

		When compiled, the app acts normally, since it doesn't have the compiler exec() wrapper there.
	*/

	logging.Logf(1, "Frontend shutdown routine initiated.")
	globals.FrontendTransientConfig.ShutdownInitiated = true
	kvstore.CloseKVStore()
	logging.Logf(1, "Sending a SIGTERM to the backend...")
	err := besupervisor.BackendDaemon.Process.Signal(syscall.SIGINT)
	if err != nil {
		logging.Logf(1, "Backend process kill returned an error. Err: %v", err)
	}
	logging.Logf(1, "Done sending a SIGTERM to the backend.")
	logging.Logf(1, "Frontend is shut down.")
}

func startSchedules() {
	logging.Log(1, "Setting up cyclical frontend tasks is starting.")
	defer logging.Log(1, "Setting up cyclical frontend tasks is complete.")
	ports.VerifyFrontendPorts()
	globals.FrontendTransientConfig.StopRefresherCycle = scheduling.ScheduleRepeat(func() {
		start := time.Now()
		refresher.Refresh()
		elapsed := time.Since(start)
		logging.Logcf(1, "We've refreshed the frontend. It took: %s", elapsed)

	}, 2*time.Minute, time.Duration(0), nil)

	// Refresh the SFW list every hour.
	globals.FrontendTransientConfig.StopSFWListUpdateCycle = scheduling.ScheduleRepeat(func() {
		start := time.Now()
		globals.FrontendConfig.ContentRelations.SFWList.Refresh()
		elapsed := time.Since(start)
		logging.Logcf(1, "We've refreshed the SFW list. It took: %s", elapsed)

		// Also - prune the notifications carrier while you're at it
		festructs.NotificationsSingleton.Prune()

	}, 1*time.Hour, time.Duration(0), nil)
}
