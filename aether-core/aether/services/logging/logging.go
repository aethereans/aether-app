// Services > Logging
// Logging is the universal logger. This library is responsible for checking whether logging to a file (or to stderr) is enabled, and if so, will process logs as such.

package logging

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/toolbox"
	"fmt"
	"github.com/fatih/color"
	"log"
	// "runtime"
)

// Prevents hitting the config every time we need this. We use this a lot.
type loggingCache struct {
	loggingInitialised  bool
	logginglevel        int
	conftypeInitialised bool
	conftype            string
	colorinitialised    bool
	color               *color.Color
}

var llcache loggingCache

// Log prints to the standard logger.
func Log(level int, input interface{}) {
	if getLoggingLevel() >= level {
		// If print to stdout is enabled, instead of logging, route to stdout. This means it's running in a swarm setup that wants the results that way for collation.
		if getPrintToStdout() {
			if getSwarmNodeId() != -1 {
				fmt.Printf("%d: %s\n", getSwarmNodeId(), input)
			} else {
				fmt.Println(input)
			}
		} else {
			// If not routed to stdout, log normally.
			log.Println(input)
		}
	}
}

func Logf(level int, input string, v ...interface{}) {
	if getLoggingLevel() >= level {
		// If print to stdout is enabled, instead of logging, route to stdout. This means it's running in a swarm setup that wants the results that way for collation.
		if getPrintToStdout() {
			if getSwarmNodeId() != -1 {
				fmt.Printf("%d: %s\n", getSwarmNodeId(), fmt.Sprintf(input, v...))
			} else {
				fmt.Printf("%s\n", fmt.Sprintf(input, v...))
			}
		} else {
			// If not routed to stdout, log normally.
			log.Printf("%s\n", fmt.Sprintf(input, v...))
		}
	}
}

func Logcf(level int, input string, v ...interface{}) {
	if !llcache.colorinitialised {
		llcache.color = color.New(color.FgHiWhite, color.BgHiBlack)
		llcache.colorinitialised = true
	}
	Logf(level, llcache.color.Sprintf(input, v...))
}

func LogCrash(input interface{}) {
	// If we are already shutting down, do not crash.
	if getShutdownInitiated() {
		return
	}
	log.Println(toolbox.DumpStack())
	log.Fatal(input)
}

func LogCrashf(input string, v ...interface{}) {
	if getShutdownInitiated() {
		return
	}
	Logf(0, fmt.Sprintf(input, v...))
	LogCrash(fmt.Sprintf(input, v...))
}

func LogObj(level int, objName string, input interface{}) {
	Logf(level, "%s: %#v", objName, input)
}

// These methods below allow the routines above to not care about whether it's a BE or a FE.

func getLoggingLevel() int {
	if llcache.loggingInitialised {
		return llcache.logginglevel
	}
	var ll int
	if getConfType() == "backend" {
		// backend
		ll = globals.BackendConfig.GetLoggingLevel()
	} else {
		// frontend
		ll = globals.FrontendConfig.GetLoggingLevel()
	}
	llcache.logginglevel = ll
	llcache.loggingInitialised = true
	return ll
}

func getShutdownInitiated() bool {
	if getConfType() == "backend" {
		return globals.BackendTransientConfig.ShutdownInitiated
	} else {
		// frontend
		return globals.FrontendTransientConfig.ShutdownInitiated
	}
}

func getSwarmNodeId() int {
	if getConfType() == "backend" {
		return globals.BackendTransientConfig.SwarmNodeId
	} else {
		// frontend
		return -1
	}
}

func getPrintToStdout() bool {
	if getConfType() == "backend" {
		return globals.BackendTransientConfig.PrintToStdout
	} else {
		// frontend
		return globals.FrontendTransientConfig.PrintToStdout
	}
}

func getConfType() string {
	if llcache.conftypeInitialised {
		return llcache.conftype
	}
	beConfInitd := globals.BackendConfig != nil
	feConfInitd := globals.FrontendConfig != nil
	if beConfInitd && feConfInitd {
		panic("Both FE and BE configs are initialised. This is a programming error. Please only init one of those.")
	}
	llcache.conftypeInitialised = true
	if beConfInitd {
		llcache.conftype = "backend"
		return llcache.conftype
	}
	llcache.conftype = "frontend"
	return llcache.conftype
}
