package cmd

import (
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/tlscerts"
	"aether-core/aether/services/toolbox"
	"fmt"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	"os"
	"path/filepath"
	// "strings"
)

// cmdRoot represents the base command when called without any subcommands
var cmdRoot = &cobra.Command{
	Use:   "aether-backend",
	Short: "Aether Runtime Environment communicates with other computers using Mim-based subprotocols, persist objects received over the network, and serves the objects to other requesters as specified by the Mim protocol spec.",
	Long: `Aether Runtime Environment is the standalone executable that handles any Mim-based app's communication with the external world. After spinning up a running instance of this executable, it will act as a database for your app that you can query through appropriate local APIs.

For example: The app 'Aether' uses ARE to communicate in c0 and dweb Mim subprotocols. For more information, please see https://getaether.net. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`You've attempted to run the Aether Runtime Environment without any commands. ARE requires some variables to be passed to it to be able to do what you want.

Please run "aether-backend --help" to see all available options.
`)
	},
}

// This is called by main.main(). It only needs to happen once to the cmdRoot.
func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Below are methods used in the cmd library. These are not specifically related to root cmd.

func showIntro() {
	// color.Cyan(`
	//                  1ttfffLLLLLLLLLffft
	//              11111ttfffLLLLLLLLLffftt111
	//           111ttfLLLCCGGG000000GGGCCLLLfft111
	//        1111ffLLCG00880000GGGG00008880GCCLLft111
	//       11tfLLCG0880GCCLLLLLCCLLLLLCCGG0880CLLft111
	//     111fLLC0880CCLLLLLLLLL08CLLLLLLLLLCG880GLLft11
	//    11tLLCG88GCLCCLLLLLLLLL08CLLLLLLLLCCLLG0@0CLLf11
	//   11tLLC0@0CLLL080CLLLLLLLG8LLLLLLLCG88CLLLG88GLLf11
	//   1tLLC8@GLLLLLCG880CLLLLLG8LLLLLCG880CLLLLLC88GLLf1
	//  11LLL0@GLLLLLLLLLG080CLLLG8LLLCG88GCLLLLLLLLC88CLLt1
	//  1tLLG@8LLLLLLLLLLLLC080CLG8LCG80GCLLLLLLLLLLLG@0LLL1
	//  1LLL0@GLLLLLLLLLLLLLLCG8008080GCLLLLLLLLLLLLLC8@CLLt
	//  1LLL0@CLLG000000000000G0@@@@800000000000000CLL0@CLLt
	//  1LLL0@GLLLCCCCCCCCCCCCG0088080CCCCCCCCCCCCCLLC8@CLLt
	//  1tLLG@8LLLLLLLLLLLLCG08GLG8LC080CLLLLLLLLLLLLG@0LLL1
	//  11LLC8@GLLLLLLLLLCG80GLLLG8LLLC080CLLLLLLLLLC8@GLLt1
	//   1tLLC8@GLLLLLLC080GLLLLLG8LLLLLC080GCLLLLLC8@GLLf1
	//   11tLLC8@0CLLLG80GLLLLLLLG8LLLLLLLC080CLLLG88GLLf11
	//    11tLLCG88GLLCCLLLLLLLLL08CLLLLLLLLCCLLC088CLLf11
	//     111fLLC0880CLLLLLLLLLL08CLLLLLLLLLCG080GLLft111
	//      111tfLLCG8880GCCLLLLLCCLLLLLLCCG0880CLLft1111
	//        1111ffLLCG0088000GGGGGGG0008800CCLLft111
	//           1111tffLLCCGG00000000GGGCLLLftt1111
	//              11111ttfffLLLLLLLLLffftt11111
	//                  1ttfffLLLLLLLLLffft
	//    `)

	colorSet := color.New(color.FgWhite)
	colorSet = colorSet.Add(color.Bold)
	colorSet = colorSet.Add(color.BgCyan)
	untetheredStateString := ""
	if !globals.BackendTransientConfig.TetheredToFrontend {
		untetheredStateString = "This node is running untethered, without a frontend."
	}
	colorSet.Printf(`

   __    __     __     __    __
  /\ "-./  \   /\ \   /\ "-./  \
  \ \ \-./\ \  \ \ \  \ \ \-./\ \
   \ \_\ \ \_\  \ \_\  \ \_\ \ \_\
    \/_/  \/_/   \/_/   \/_/  \/_/

     Aether Runtime Environment
     App version: %s
     Protocol version: %s
     %s

`, fmt.Sprintf(
		"v%d.%d.%d",
		globals.BackendConfig.GetClientVersionMajor(),
		globals.BackendConfig.GetClientVersionMinor(),
		globals.BackendConfig.GetClientVersionPatch()),
		fmt.Sprintf(
			"v%d.%d",
			globals.BackendConfig.GetProtocolVersionMajor(),
			globals.BackendConfig.GetProtocolVersionMinor()),
		untetheredStateString)
}

// Start methods for Mim.

// EstablishConfigs establishes the configs (both transient and permanent) based on the flags provided through the CLI. Some flags are available directly (if they're local variables) and some of them are saved into transient config, made available globally until the app quits. If you need to have the data that is provided by the flag used in multiple places, create a new item in the transient config and put it there, it will be made available to the whole app. If you need the value once (ex: inserting a value into a database) then you can just use the value from the flags struct.
func EstablishConfigs(cmd *cobra.Command) flags {
	// Cmd can be nil, in which case it's running under a testing environment.
	flgs := flags{}
	if cmd != nil {
		flgs = renderFlags(cmd)
	}
	// Transient configs are established before permanent (saved to file) configs because appname and orgname in the transient configs determine where permanent configs get saved to. This is useful when running swarm tests, because specifying these effectively makes a swarm spawn save configs and caches into a different location than what it would normally not save.
	globals.BackendTransientConfig = &configstore.Btc
	globals.BackendTransientConfig.SetDefaults()
	// globals.FrontendTransientConfig = &configstore.Ftc
	// globals.FrontendTransientConfig.SetDefaults()
	// Set the transient config data.
	if cmd != nil {
		if flgs.appName.changed {
			globals.BackendTransientConfig.AppIdentifier = flgs.appName.value.(string)
		}
		if flgs.orgName.changed {
			globals.BackendTransientConfig.OrgIdentifier = flgs.orgName.value.(string)
		}
	} else { // cmd == nil, this is a unit test.
		globals.BackendTransientConfig.AppIdentifier = "A-UnitTest"
	}
	// Establish permanent configs.
	becfg, err := configstore.EstablishBackendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	becfg.Cycle()
	globals.BackendConfig = becfg
	// fecfg, err := configstore.EstablishFrontendConfig()
	// if err != nil {
	// 	logging.LogCrash(err)
	// }
	// fecfg.Cycle()
	// globals.FrontendConfig = fecfg
	// Generate TLS keys
	tlscerts.Generate()
	// Determine whether the configs have been manipulated by flags. If so, disable editing of permanent configs for this session.
	if cmd != nil && flagsChanged(cmd) {
		globals.BackendTransientConfig.PermConfigReadOnly = true
		// globals.FrontendTransientConfig.PermConfigReadOnly = true
	}
	if cmd != nil && flgs.loggingLevel.changed {
		// Start setting permanent configs. These are NO-OPs if the permament config is read only.
		globals.BackendConfig.SetLoggingLevel(flgs.loggingLevel.value.(int))
	}

	// If the permanent config is read only, we probably should tell.
	if globals.BackendTransientConfig.PermConfigReadOnly {
		// This is double because we want to both print on the console, and have it in the logs. Also providing even the default value explicitly through the command line triggers a changed=true, so even if you do logginglevel=0 (whose default is also 0), the configs will end up read only.
		fmt.Println("Configuration read only. Configuration for both backend and the frontend will be treated as read only because command line flags were provided for this run.")
		logging.Log(1, fmt.Sprint("Configuration read only. Configuration for both backend and the frontend will be treated as read only because command line flags were provided for this run."))
	}
	if flgs.port.changed {
		// Heads up: this looks like it's not working because we have a global switch that disables writes to the config file if any cl flags are given. We need to make that switch a little less aggressive.
		globals.BackendConfig.SetExternalPort(flgs.port.value.(int))
	}
	if flgs.externalIp.changed {
		globals.BackendConfig.SetExternalIp(flgs.externalIp.value.(string))
	}
	// These are booleans but still gated, because the values can be nil in case of testing, not only true/false.
	if flgs.printToStdout.changed {
		globals.BackendTransientConfig.PrintToStdout = flgs.printToStdout.value.(bool)
	}
	if flgs.metricsDebugMode.changed {
		globals.BackendTransientConfig.MetricsDebugMode = flgs.metricsDebugMode.value.(bool)
	}
	if flgs.swarmNodeId.changed {
		globals.BackendTransientConfig.SwarmNodeId = flgs.swarmNodeId.value.(int)
	} else {
		globals.BackendTransientConfig.SwarmNodeId = -1 // If not given,disable
	}
	if flgs.appName.changed {
		// Also change the client name so that the name change communicates out into the analytics server when under orchestrate test harness. This is different than AppIdentifier which determines the folders that the node saves to the local drive.
		globals.BackendConfig.SetClientName(flgs.appName.value.(string))
	}
	if flgs.fpCheckEnabled.changed {
		globals.BackendTransientConfig.FingerprintCheckEnabled = flgs.fpCheckEnabled.value.(bool)
	}
	if flgs.sigCheckEnabled.changed {
		globals.BackendTransientConfig.SignatureCheckEnabled = flgs.sigCheckEnabled.value.(bool)
	}
	if flgs.powCheckEnabled.changed {
		globals.BackendTransientConfig.ProofOfWorkCheckEnabled = flgs.powCheckEnabled.value.(bool)
	}
	if flgs.pageSigCheckEnabled.changed {
		globals.BackendTransientConfig.PageSignatureCheckEnabled = flgs.pageSigCheckEnabled.value.(bool)
	}
	if flgs.tlsEnabled.changed {
		globals.BackendTransientConfig.TLSEnabled = flgs.tlsEnabled.value.(bool)
	}
	if flgs.backendAPIPort.changed {
		globals.BackendConfig.SetBackendAPIPort(flgs.backendAPIPort.value.(int))
	}
	if flgs.backendAPIPublic.changed {
		globals.BackendConfig.SetBackendAPIPublic(flgs.backendAPIPublic.value.(bool))
	}

	if flgs.adminFeAddr.changed {
		globals.BackendConfig.SetAdminFrontendAddress(flgs.adminFeAddr.value.(string))
		globals.BackendTransientConfig.TetheredToFrontend = true
	}

	if flgs.adminFePk.changed {
		globals.BackendConfig.SetAdminFrontendPublicKey(flgs.adminFePk.value.(string))
		globals.BackendTransientConfig.TetheredToFrontend = true
	}

	if flgs.allowLocalhostRemotes.changed {
		globals.BackendTransientConfig.AllowLocalhostRemotes = flgs.allowLocalhostRemotes.value.(bool)
	}

	// Set up the DB Instance so that we get access to the database.
	if globals.BackendConfig.GetDbEngine() == "sqlite" {
		dbLoc := filepath.Join(globals.BackendConfig.GetSQLiteDBLocation(), "AetherDB.db")
		if !toolbox.FileExists(dbLoc) {
			// Db doesn't exist. Make sure that the bootstrap timer and event horizon is reset. Those values depend on the database, and if the DB is deleted while the user settings are not, they can prevent a bootstrap from happening as it should. In the other case where the database isn't created yet, these calls are idempotent.
			fmt.Println("The database was deleted or is not created yet. Setting event horizon and last successful live, static, bootstrap timestamps to 0.\n")
			logging.Logf(1, "The database was deleted or is not created yet. Setting event horizon and last successful live, static, bootstrap timestamps to 0.\n")
			globals.BackendConfig.ResetEventHorizon()
			globals.BackendConfig.ResetLastLiveAddressConnectionTimestamp()
			globals.BackendConfig.ResetLastStaticAddressConnectionTimestamp()
			globals.BackendConfig.ResetLastBootstrapAddressConnectionTimestamp()
		}
		conn, err := sqlx.Connect(
			"sqlite3", dbLoc)
		if err != nil {
			fmt.Println("Error: Failed to open the SQLite database in dbLoc:", dbLoc)
			logging.LogCrash(err)
		}
		globals.DbInstance = conn
		/*
			SetMaxOpenConns set to 1 is critical.

			Why?

			To keep SQLite happy, and discipline. This makes the app crash if you end up attempting to do more than one read OR write at the same time. Unfortunately, SQLite behaves very unpredictably under any other condition. See here: https://gist.github.com/mrnugget/0eda3b2b53a70fa4a894

			Everything can be in 3 amounts: none, 1, or n. If you go into n, you'll eventually leak DB connections and under high enough load, it'll all come crashing down.

			Strictly speaking, this is only useful in the case of SQLite (which does not handle concurrency) and not in the case of MySQL or any other database backend you might use. But this needs to be there, so that it enforces, a) you have good connection hygiene, b) the code you write is at some point portable to SQLite, and thus by proxy to desktop, regular users.

			This is the magic sauce that makes SQLite not break under load.
		*/
		globals.DbInstance.SetMaxOpenConns(1)

	} else if globals.BackendConfig.GetDbEngine() == "mysql" {
		// If you want to use the MySQL, create a 'AetherDB' in your MySQL instance and insert the username / password here.
		/*
			MySQL connection string syntax:
			root:PASSWORD@tcp(l:3306)/sqlx_test
		*/
		mysqlConnectionString := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/",
			globals.BackendConfig.GetDbUsername(),
			globals.BackendConfig.GetDbPassword(),
			globals.BackendConfig.GetDbIp(),
			globals.BackendConfig.GetDbPort())
		fmt.Println(mysqlConnectionString)
		globals.DbInstance = sqlx.MustConnect("mysql", mysqlConnectionString)
		/*
			Single-connection restriction does not apply to MySQL. Relaxing it helps in server situations.
		*/
		globals.DbInstance.SetMaxOpenConns(100)
	} else {
		logging.LogCrash(fmt.Sprintf("Storage engine you've inputted is not supported. Please change it from the backend user config into something that is supported. You've provided: %s", globals.BackendConfig.GetDbEngine()))
	}
	// Delete all cached post responses from the last run if it wasn't cleaned properly.
	globals.BackendTransientConfig.POSTResponseRepo.DeleteAllFromDisk()
	return flgs
}
