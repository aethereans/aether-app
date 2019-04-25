// Backend > FrontendAPIClient
// This package sends requests to the frontend API server to let know of the backend status.

// Heads up: this API talks only to the *admin* frontend, not to other frontends that might be using this.

package feapiconsumer

import (
	pb "aether-core/aether/protos/feapi"
	"aether-core/aether/protos/feobjects"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	// "fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func StartFrontendAPIConnection() (pb.FrontendAPIClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(globals.BackendConfig.GetAdminFrontendAddress(), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(toolbox.MaxInt32)))
	if err != nil {
		logging.Logf(1, "Could not connect to the frontend API service. Error: %v", err)
	}
	c := pb.NewFrontendAPIClient(conn)
	return c, conn
}

func SendBackendReady() {
	if !globals.BackendTransientConfig.TetheredToFrontend {
		return // If untethered, we don't send anything out on our own - there is no one to send to.
	}
	c, conn := StartFrontendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.BackendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := pb.BEReadyRequest{
		Address: globals.BackendConfig.GetExternalIp(),
		Port:    int32(globals.BackendConfig.GetBackendAPIPort()),
	}
	_, err := c.BackendReady(ctx, &payload)
	if err != nil {
		logging.Logf(1, "SendBackendReady encountered an error. Err: %v", err)
	}
}

// The way this works is that you set whatever you need into the backend ambient status here, and calling SendBackendAmbientStatus will send it over. The reason for that is when you send a partial backend ambient status, it will actually delete the other values from the client, because the client cannot know if the value was set to its zero value, or wasn't set at all.
var BackendAmbientStatus feobjects.BackendAmbientStatus

func SendBackendAmbientStatus() {
	if !globals.BackendTransientConfig.TetheredToFrontend {
		return // If untethered, we don't send anything out on our own - there is no one to send to.
	}
	c, conn := StartFrontendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.BackendConfig.GetGRPCServiceTimeout())
	defer cancel()
	if len(BackendAmbientStatus.BackendConfigLocation) == 0 {
		BackendAmbientStatus.BackendConfigLocation = globals.GetBackendConfigLocation()
	}
	/*----------  Update the relevant backend data.  ----------*/
	// Last Inbound
	BackendAmbientStatus.LastInboundConnTimestamp = globals.BackendTransientConfig.Bouncer.GetLastInboundSyncTimestamp(false)
	// Last 15m inbounds
	BackendAmbientStatus.InboundsCount15 = int32(len(globals.BackendTransientConfig.Bouncer.GetInboundsInLastXMinutes(15, true)))
	// Last 15m outbounds
	BackendAmbientStatus.OutboundsCount15 = int32(len(globals.BackendTransientConfig.Bouncer.GetOutboundsInLastXMinutes(15, true)))
	/*----------  END Update the relevant backend data.  ----------*/
	payload := pb.BackendAmbientStatusPayload{
		BackendAmbientStatus: &BackendAmbientStatus,
	}
	_, err := c.SendBackendAmbientStatus(ctx, &payload)
	if err != nil {
		logging.Logf(1, "SendBackendAmbientStatus encountered an error. Err: %v", err)
	}
}
