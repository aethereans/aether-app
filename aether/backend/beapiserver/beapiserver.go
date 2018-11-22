package beapiserver

import (
	pb "aether-core/aether/protos/beapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
	"fmt"
	// "golang.org/x/net/context"
	"net"
	"strconv"
)

func StartBackendServer(gotValidPort chan bool) {
	logging.Logf(1, "Starting Backend API server.")
	defer logging.Logf(1, "Stopped Backend API server.")
	if globals.BackendConfig.GetBackendAPIPublic() {
		logBackendAPIServingPubliclyWarning()
	}
	var backendSourceIp string
	if globals.BackendConfig.GetBackendAPIPublic() {
		backendSourceIp = "0.0.0.0" // listen publicly
	} else {
		backendSourceIp = "127.0.0.1" // listen locally
	}
	sourceAddr := fmt.Sprint(backendSourceIp, ":", globals.BackendConfig.GetBackendAPIPort())
	listener, err := net.Listen("tcp4", sourceAddr)
	if err != nil {
		// Couldn't bind to the given port. Ask the OS to give us a free port instead, and save that port.
		var err2 error
		sourceAddr = fmt.Sprint(backendSourceIp, ":", 0)
		listener, err2 = net.Listen("tcp4", sourceAddr)
		if err2 != nil {
			logging.LogCrashf("The backend API Server could not start listening. Err: %v", err2)
		}
		_, sPort, _ := net.SplitHostPort(listener.Addr().String())
		sPortAsInt, _ := strconv.Atoi(sPort)
		globals.BackendConfig.SetBackendAPIPort(sPortAsInt)
	}
	gotValidPort <- true
	maxmsgsize := grpc.MaxMsgSize(toolbox.MaxInt32) // 12gb. debug TODO
	maxrecv := grpc.MaxRecvMsgSize(toolbox.MaxInt32)
	maxsend := grpc.MaxSendMsgSize(toolbox.MaxInt32)
	s := grpc.NewServer(maxmsgsize, maxrecv, maxsend)
	pb.RegisterBackendAPIServer(s, &server{})
	err2 := s.Serve(listener)
	if err2 != nil {
		logging.LogCrashf("Failed to serve the Backend API. %v", err2)
	}
}

func logBackendAPIServingPubliclyWarning() {
	logging.Logf(1, `

///////////////////////////////////////////////////////////////////////////////
HEADS UP: You are serving this backend API (that talks to frontends, not the one who talks to the Mim network) publicly on the Internet at port %d of this machine's external IP.
If this is not your intent, do not provide the 'backendapipublic' flag at startup. This setting is false by default.
///////////////////////////////////////////////////////////////////////////////
  `, globals.BackendConfig.GetBackendAPIPort())
}
