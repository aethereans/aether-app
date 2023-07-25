// Services > Ports
// This package provides function related to ports in the local machine, such as finding a free port, or checking whether a port is open for use.

package ports

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"
	"net"
	"strings"
)

// GetFreePort returns a free port that is currently unused in the local system.
func GetFreePort() int {
	a, err := net.ResolveTCPAddr("tcp4", ":0")
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not parse the TCP address in an attempt to get a free port. The error raised was: %s", err))
	}
	l, err := net.ListenTCP("tcp4", a)
	defer l.Close()
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not listen to TCP in an attempt to get a free port. The error raised was: %s", err))
	}
	return l.Addr().(*net.TCPAddr).Port
}

// GetFreePorts returns a number of free ports that are currently unused in the local system.
func GetFreePorts(number int) []int {
	var ports []int

	clashcount := 0
	checkport := func(ports []int, port int) bool {
		for _, val := range ports {
			if val == port {
				clashcount++
				if clashcount > 65535 {
					logging.LogCrash(fmt.Sprintf("This computer does not have enough ports that are free. You've requested %d free ports. ", number))
				}
				return true
			}
		}
		return false
	}
	for i := 0; i < number; i++ {
		port := GetFreePort()
		for checkport(ports, port) {
			port = GetFreePort()
		}
		ports = append(ports, port)
	}
	return ports
}

// CheckPortAvailability checks for whether a port that it is given is currently free to use.
func CheckPortAvailability(port uint16) bool {
	if port == 49999 {
		/*
			We want all nodes to pick a random port that's not the default. Why? Because you can have multiple nodes in different computers behind the same IP address. In that case, both of these nodes will pick 49999 as their default port, and their UPNP port maps override each other - the router will have to map port 49999 to one or the other. That means only one of an arbitrary number of nodes behind a single IP address can be accessible publicly, which is not ideal.

			Allowing the ports to be chosen randomly at first boot allows us to avoid that. Rendering the default port permanently unavailable is a good way to get the app to choose a new port if the port found turns out to be the default one.
		*/
		return false
	}
	a, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		logging.LogCrash(fmt.Sprintf("We could not parse the TCP address in an attempt to check the availability of the given port. The error raised was: %s, The port attempted to be checked was: %d", err, port))
	}
	l, err := net.ListenTCP("tcp4", a)
	defer l.Close()
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") || strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "Only one usage of each socket address (protocol/network address/port) is normally permitted.") { // Last one is Windows.
			logging.Log(1, fmt.Sprintf("Port number %d is already in use. Error: %s", port, err))
			return false
		} else {
			logging.Logf(1, "We attempted to check the availability of the port %d on the current computer and it failed with this error: %v", port, err)
			return false
		}
	}
	return true
}

// VerifyBackendPorts verifies the local port available in the config, and if it is not available, replaces it with one that is. Then it flips the bit to mark the local port as verified. Backend ports are the external port that talks to other Mim nodes, and the backend api port that talks to frontends.
func VerifyBackendPorts() {
	logging.Log(2, "VerifyBackendPorts check is running.")
	defer logging.Log(2, "VerifyBackendPorts check is done.")
	// This check only runs once per start.
	// Check the external port that talks to other Mim nodes
	if !globals.BackendTransientConfig.ExternalPortVerified {
		// Prevent race condition in which any number of calls can enter this before CheckPortAvailability returns.
		globals.BackendTransientConfig.ExternalPortVerified = true
		if CheckPortAvailability(globals.BackendConfig.GetExternalPort()) {
			logging.Log(1, fmt.Sprintf("The backend external port %d is verified to be open and available for use.", globals.BackendConfig.GetExternalPort()))
		} else {
			freeport := GetFreePort()
			logging.Log(1, fmt.Sprintf("The backend external port number for this node has changed. New external port for this node is: %d", freeport))
			globals.BackendConfig.SetExternalPort(freeport)
		}
	}
	// Check the backend API port that talks to frontends
	if !globals.BackendTransientConfig.BackendAPIPortVerified {
		globals.BackendTransientConfig.BackendAPIPortVerified = true
		if CheckPortAvailability(globals.BackendConfig.GetBackendAPIPort()) {
			logging.Log(1, fmt.Sprintf("The backend backend port %d is verified to be open and available for use.", globals.BackendConfig.GetBackendAPIPort()))
		} else {
			freeport := GetFreePort()
			logging.Log(1, fmt.Sprintf("The backend external port number for this node has changed. New external port for this node is: %d", freeport))
			globals.BackendConfig.SetBackendAPIPort(freeport)
		}
	}
}

// VerifyFrontendPorts verifies the frontend ports for use, and if they are not available, replaces with ones that are. Frontend ports are: frontend API port that talks to clients and backends.
func VerifyFrontendPorts() {
	logging.Log(1, "VerifyFrontendPorts check is running.")
	defer logging.Log(1, "VerifyFrontendPorts check is done.")
	// Check the frontend API server port that talks to clients and backends.
	if !globals.FrontendTransientConfig.FrontendAPIPortVerified {
		globals.FrontendTransientConfig.FrontendAPIPortVerified = true
		if CheckPortAvailability(globals.FrontendConfig.GetFrontendAPIPort()) {
			logging.Log(1, fmt.Sprintf("The frontend API server port %d is verified to be open and available for use.", globals.FrontendConfig.GetFrontendAPIPort()))
		} else {
			freeport := GetFreePort()
			logging.Log(1, fmt.Sprintf("The frontend API server port number for this node has changed. New external port for this node is: %d", freeport))
			globals.FrontendConfig.SetFrontendAPIPort(freeport)
		}
	}
}
