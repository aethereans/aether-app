// Aether Core Backend
// Application start. Starts all the necessary components. The main body of the regular application start is at cmd/run. See that for actual start procedure.

package main

import (
	"aether-core/aether/backend/cmd"
)

func main() {
	cmd.Execute()
}
