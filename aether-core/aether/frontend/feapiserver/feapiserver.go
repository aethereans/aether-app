// Frontned > FEApiServer
// This package handles the frontend server both client and backend utilises.

package feapiserver

import (
	pb "aether-core/aether/protos/feapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
	// "time"
)

var sourceIP = "127.0.0.1"

func StartFrontendServer(gotValidPort chan bool) {
	logging.Logf(1, "Starting the Frontend API server. ")
	defer logging.Logf(1, "The frontend API server is shut down.")
	sourceAddr := fmt.Sprint(sourceIP, ":", globals.FrontendConfig.GetFrontendAPIPort())
	listener, err := net.Listen("tcp4", sourceAddr)
	if err != nil {
		// Couldn't bind to the given port. Ask the OS to give us a free port instead, and save that port.
		var err2 error
		sourceAddr = fmt.Sprint(sourceIP, ":", 0)
		listener, err2 = net.Listen("tcp4", sourceAddr)
		if err2 != nil {
			logging.LogCrashf("The Frontend Server could not start listening. Err: %v", err2)
		}
		_, sPort, _ := net.SplitHostPort(listener.Addr().String())
		sPortAsInt, _ := strconv.Atoi(sPort)
		globals.FrontendConfig.SetFrontendAPIPort(sPortAsInt)
	}
	gotValidPort <- true
	maxmsgsize := grpc.MaxMsgSize(toolbox.MaxInt32) // 12gb. debug TODO
	maxrecv := grpc.MaxRecvMsgSize(toolbox.MaxInt32)
	maxsend := grpc.MaxSendMsgSize(toolbox.MaxInt32)
	s := grpc.NewServer(maxmsgsize, maxrecv, maxsend)
	// s := grpc.NewServer()
	pb.RegisterFrontendAPIServer(s, &server{})
	err2 := s.Serve(listener)
	if err2 != nil {
		logging.LogCrashf("Failed to serve the Frontend API. %v", err2)
	}
}
