// Services > UPNP
// This module provides UPNP port mapping functionality for routers, so that a node that is behind a router can still be accessed by other nodes.

package upnp

import (
	"aether-core/aether/backend/feapiconsumer"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"
	extUpnp "github.com/NebulousLabs/go-upnp"
	// "time"
)

// var router *extUpnp.IGD
// var err error

func MapPort() {
	feapiconsumer.BackendAmbientStatus.UPNPStatus = "In progress"
	feapiconsumer.SendBackendAmbientStatus()       // send first state
	defer feapiconsumer.SendBackendAmbientStatus() // send end state
	// External port verification has moved to the start of the main routine, so we do no longer need to do that here. We can assume that at this point, the external port is verified.
	router, err := extUpnp.Discover()
	if err != nil {
		// Either could not be found, or connected to the internet directly.
		logging.Log(3, fmt.Sprintf("A router to port map could not be found. This computer could be directly connected to the Internet without a router. Error: %s", err.Error()))
		feapiconsumer.BackendAmbientStatus.UPNPStatus = "Mapping failed, no router, or router uncooperative"
		return
	}
	extIp, err2 := router.ExternalIP()
	if err2 != nil {
		// External IP finding failed.
		logging.Log(1, fmt.Sprintf("External IP of this machine could not be determined. Error: %s", err2.Error()))
	} else {
		globals.BackendConfig.SetExternalIp(extIp)
		logging.Log(1, fmt.Sprintf("This computer's external IP is %s", globals.BackendConfig.GetExternalIp()))
	}
	// Attempt to map the external port
	err3 := router.Forward(globals.BackendConfig.GetExternalPort(), "Aether")
	if err3 != nil {
		// Router is there, but port mapping failed.
		logging.Log(1, fmt.Sprintf("In an attempt to port map, the router was found, but the port mapping failed for the backend port. Error: %s", err3.Error()))
		feapiconsumer.BackendAmbientStatus.UPNPStatus = "Mapping failed, router did not accept"
		return
		// (We can return here, because if the external port mapping [serving other nodes] has failed, there is no point in attempting to map the port that serves the frontend)
	}
	logging.Log(1, fmt.Sprintf("Port mapping was successful. We mapped backend port %d to this computer.", globals.BackendConfig.GetExternalPort()))

	// Attempt to map the reverse open port (which is external port - 1)
	err4 := router.Forward(globals.BackendConfig.GetExternalPort()-1, "Aether Reverse")
	if err4 != nil {
		// Router is there, but port mapping failed.
		logging.Log(1, fmt.Sprintf("In an attempt to port map, the router was found, but the port mapping failed for the backend reverse open port. Error: %s", err4.Error()))
		feapiconsumer.BackendAmbientStatus.UPNPStatus = "Mapping failed, router did not accept"
		return
	}

	logging.Log(1, fmt.Sprintf("Port mapping was successful. We mapped backend reverse open port %d to this computer.", globals.BackendConfig.GetExternalPort()-1))

	// If the backend api is public, attempt to map the backend API port
	if globals.BackendConfig.GetBackendAPIPublic() {
		err4 := router.Forward(globals.BackendConfig.GetBackendAPIPort(), "Aether Backend API Port")
		if err4 != nil {
			// Router is there, but port mapping failed.
			logging.Log(1, fmt.Sprintf("In an attempt to port map, the router was found, but the port mapping failed for the Backend API port. Error: %s", err3.Error()))
			feapiconsumer.BackendAmbientStatus.UPNPStatus = "Mapping failed, router did not accept"
			return
		}
		logging.Log(1, fmt.Sprintf("Port mapping was successful. We mapped backend API port %d to this computer.", globals.BackendConfig.GetBackendAPIPort()))
	}
	feapiconsumer.BackendAmbientStatus.UPNPStatus = "Successful"
}
