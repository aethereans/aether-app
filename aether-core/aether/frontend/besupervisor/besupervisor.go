// Frontend > BackendSupervisor

// This package handles the supervisory tasks related to the backend this frontend is the admin (admin) of.

package besupervisor

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	// "time"
)

var BackendDaemon *exec.Cmd

var localBackendRestartAttempts int

func StartLocalBackend() {
	if globals.FrontendConfig.GetLocalDevBackendEnabled() {
		// Development
		BackendDaemon = exec.Command("go", constructArgs(true)...)
		BackendDaemon.Stdout = os.Stdout
		BackendDaemon.Stderr = os.Stderr
		BackendDaemon.Dir = "../../../aether-core/aether/backend"
		logging.Log(1, "Local backend being started")
		err := BackendDaemon.Run()
		if err != nil {
			logging.Logf(1, "Local backend had an error. Err: %v", err)
		}
	} else {
		// Production
		BackendDaemon = exec.Command(generateBackendPath(), constructArgs(false)...)
		BackendDaemon.Stdout = os.Stdout
		BackendDaemon.Stderr = os.Stderr
		logging.Log(1, "Local backend being started")
		err2 := BackendDaemon.Run()
		if err2 != nil {
			logging.Logf(1, "Local backend had an error. Err: %v", err2)
		}
	}
	logging.Log(1, "Local backend exited.")
	if globals.FrontendTransientConfig.ShutdownInitiated {
		// If this is flipped true, the backend closed as a result of that - it did not crash. We exit the app here.
		os.Exit(0)
	}
	if localBackendRestartAttempts < 3 {
		localBackendRestartAttempts++
		logging.Log(1, "Attempting to restart the local backend.")
		StartLocalBackend()
	} else {
		logging.Log(1, "Local backend crashed more than 3 times in this run, something went very wrong. Killing the frontend.")
		os.Exit(1)
	}
}

// generateBackendPath generates where the backend executable should be. This is OS-dependent, so we unfortunately have to go with good old regular OS-level if/else.
func generateBackendPath() string {
	// dir, err := os.Getwd()
	// if err != nil {
	// 	logging.Logf(1, "Getting working directory had an error. Error: %v", err)
	// }
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logging.Logf(1, "Getting working directory had an error. Error: %v", err)
	}
	bePath := ""
	// X64 (64 bit)
	if runtime.GOARCH == "amd64" {
		// X64 OS X
		if runtime.GOOS == "darwin" {
			bePath = filepath.Join(dir, "aether-backend-mac-x64")
		}
		// X64 Windows
		if runtime.GOOS == "windows" {
			bePath = filepath.Join(dir, "aether-backend-win-x64")
		}
		// X64 Linux
		if runtime.GOOS == "linux" {
			bePath = filepath.Join(dir, "aether-backend-linux-x64")
		}
	}
	// X86 (32 bit)
	if runtime.GOARCH == "386" {
		// OS X
		if runtime.GOOS == "darwin" {
			bePath = filepath.Join(dir, "aether-backend-mac-ia32")
		}
		// Windows
		if runtime.GOOS == "windows" {
			bePath = filepath.Join(dir, "aether-backend-windows-ia32")
		}
		// Linux
		if runtime.GOOS == "linux" {
			bePath = filepath.Join(dir, "aether-backend-linux-ia32")
		}
	}
	// 64 bit ARM
	if runtime.GOARCH == "arm64" {
		// OS X
		if runtime.GOOS == "darwin" {
			bePath = filepath.Join(dir, "aether-backend-mac-arm64")
		}
		// Windows
		if runtime.GOOS == "windows" {
			bePath = filepath.Join(dir, "aether-backend-windows-arm64")
		}
		// Linux
		if runtime.GOOS == "linux" {
			bePath = filepath.Join(dir, "aether-backend-linux-arm64")
		}
	}
	// 32 bit ARM
	if runtime.GOARCH == "arm" {
		// OS X
		if runtime.GOOS == "darwin" {
			bePath = filepath.Join(dir, "aether-backend-mac-arm32")
		}
		// Windows
		if runtime.GOOS == "windows" {
			bePath = filepath.Join(dir, "aether-backend-windows-arm32")
		}
		// Linux
		if runtime.GOOS == "linux" {
			bePath = filepath.Join(dir, "aether-backend-linux-arm32")
		}
	}
	return bePath
}

func constructArgs(development bool) []string {
	var baseCmd []string

	// backendLogginglevel := 1
	if development {
		compilerTags := ""
		/*
		  {{ COMPILE INSTRUCTIONS }}
		  To run in extvenabled in development, you need to comment out the line below
		*/
		// compilerTags = "extvenabled"
		// ^^^^^ This line

		baseCmd = []string{"run", "-tags", compilerTags, "main.go", "run"}
	} else {
		baseCmd = []string{"run"}
		// backendLogginglevel = 0
	}
	fesrvaddr := "127.0.0.1"
	// fesrvport := globals.FrontendTransientConfig.FrontendServerPort
	fesrvport := globals.FrontendConfig.GetFrontendAPIPort()
	fePublicKey := globals.FrontendConfig.GetMarshaledFrontendPublicKey()
	// baseCmd = append(baseCmd, fmt.Sprintf("--logginglevel=%d", backendLogginglevel))
	baseCmd = append(baseCmd, fmt.Sprintf("--adminfeaddr=%s:%d", fesrvaddr, fesrvport))
	baseCmd = append(baseCmd, fmt.Sprintf("--adminfepk=%s", fePublicKey))
	return baseCmd
}
