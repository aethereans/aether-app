// Services > SafeSleep
// Safesleep provides a sleep function that does not prevent application shutdown.

package safesleep

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"errors"
	"fmt"
	"time"
)

func Sleep(dur time.Duration, terminator *bool) error {
	var shutdownIndicator *bool
	if globals.BackendTransientConfig != nil {
		shutdownIndicator = &globals.BackendTransientConfig.ShutdownInitiated
	}
	if globals.FrontendTransientConfig != nil {
		shutdownIndicator = &globals.FrontendTransientConfig.ShutdownInitiated
	}
	if terminator == nil {
		f := false
		terminator = &f
	}
	// Split time.duration into 10 second intervals
	sec := int(dur.Seconds())
	var blocks int
	if sec <= 5 {
		// Sleep for the exact time if it's less than or exactly 10 seconds.
		time.Sleep(dur)
		if *terminator || *shutdownIndicator {
			logging.Log(2, fmt.Sprintf("Sleep terminator was flipped true, so SafeSleep is exiting. Duration was: %s", dur))
			return errors.New("Sleep terminator was flipped true. Please exit gracefully.")
		} else {
			return nil
		}
	} else {
		blocks = sec / 1
	}
	for i := 0; i < blocks; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		if *terminator || *shutdownIndicator {
			logging.Log(2, fmt.Sprintf("Sleep terminator was flipped true, SafeSleep is exiting. Duration was: %s", dur))
			return errors.New("Sleep terminator was flipped true. Please exit gracefully.")
		}
	}
	return nil
}
