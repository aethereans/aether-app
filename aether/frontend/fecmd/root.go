package fecmd

import (
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	"os"
)

// cmdRoot represents the base command when called without any subcommands
var cmdRoot = &cobra.Command{
	Use:   "aetherfe",
	Short: "Aether FE communicates with the backend using gRPC, caches and compiles data received frmo the backend, and serves this precompiled data to the client.",
	Long: `Aether FE is the caching and compilation layer between the backend and what the user sees in the client. Most client queries will be responded by the frontend, and in the cases where the backend lives in a separate machine, the FE+CL combo acts as the user's client.

For more information, please see https://getaether.net. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`You've attempted to run the Aether Frontend without any commands. Aether FE requires some variables to be passed to it to work.

Please run "aetherfe --help" to see all available options.
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

func EstablishConfigs(cmd *cobra.Command) flags {
	// Cmd can be nil, in which case it's running under a testing environment.
	flgs := flags{}
	if cmd != nil {
		flgs = renderFlags(cmd)
	}
	globals.FrontendTransientConfig = &configstore.Ftc
	globals.FrontendTransientConfig.SetDefaults()

	fecfg, err := configstore.EstablishFrontendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	fecfg.Cycle()
	globals.FrontendConfig = fecfg
	if cmd == nil {
		return flgs
	}
	if flgs.loggingLevel.changed {
		fmt.Println("loggingLevel changed. new value:")
		fmt.Println(flgs.loggingLevel.value)
		globals.FrontendConfig.SetLoggingLevel(flgs.loggingLevel.value.(int))
	}
	if flgs.clientIp.changed {
		globals.FrontendConfig.SetClientAPIAddress(flgs.clientIp.value.(string))
	}
	if flgs.clientPort.changed {
		globals.FrontendConfig.SetClientPort(flgs.clientPort.value.(int))
	}
	if flgs.isDev.changed {
		globals.FrontendConfig.SetLocalDevBackendEnabled(flgs.isDev.value.(bool))
	}
	return flgs
}
