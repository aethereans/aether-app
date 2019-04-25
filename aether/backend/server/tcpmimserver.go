package server

import (
	"aether-core/aether/backend/dispatch"
	"aether-core/aether/io/api"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/tcpmim"
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

/*
	StartTCPMimServer is the quickest way to boot up a TCPMim server that can respond to reverse open requests.
*/
func StartTCPMimServer() {
	logging.Logf(1, "StartTCPMimServer enters.")
	// For now, the only use case we have for TCPMim is reverse open requests. If this node does not accept inbound reverse open requests, TCPMim server is not started.
	if globals.BackendConfig.GetDeclineInboundReverseRequests() {
		logging.Logf(0, "TCPMimServer: Not starting because user has chosen to not allow inbound reverse open requests.")
		return
	}
	tms := TCPMimServer{}
	go tms.Serve()
}

type TCPMimServerConfig struct {
	Network string
	Host    string
	Port    uint16
	Timeout time.Duration
}

type TCPMimServer struct {
	Config *TCPMimServerConfig
}

func (t *TCPMimServer) Serve() {
	if t.Config == nil {
		logging.Logf(0, "TCPMimServer: No config given. Using defaults.")
		cfg := TCPMimServerConfig{
			Network: "tcp4",
			Host:    "",
			Port:    globals.BackendConfig.GetExternalPort() - 1,
			Timeout: 10 * time.Second,
		}
		t.Config = &cfg
	}
	addr := t.Config.Host + ":" + strconv.Itoa(int(t.Config.Port))
	l, err := net.Listen(t.Config.Network, addr)
	if err != nil {
		logging.Logf(0, "TCPMimServer: Listener had an error and is exiting. Err: %v", err)
		return
	}
	defer l.Close()
	logging.Logf(0, "TCPMimServer: Started listening for %v on %v", t.Config.Network, addr)
	for {
		// If we're lameduck, exit the server.
		if globals.BackendTransientConfig.LameduckInitiated || globals.BackendTransientConfig.ShutdownInitiated {
			return
		}
		conn, err := l.Accept()
		if err != nil {
			logging.Logf(0, "TCPMIMServer: Accept conn failed. Err: %v", err)
			if conn != nil {
				conn.Close()
			}
			continue
		}
		conn.SetDeadline(time.Now().Add(t.Config.Timeout))
		go t.HandleConn(conn)
	}
}

/*
MaybeStartSync checks whether we have a slot allowed in our outbound gate. If so, this will claim a slot (lease), and it will start the sync. This requesting outbound lease logic used here also happens in the sync itself. This is fine, because requesting a lease, if one is present, is idempotent. Likewise, returning a lease is idempotent if the lease has already been returned.

This does not actually make use of lease renewal and return functions provided here. The reason why is that we just get the lease here, and when the appropriate sync enters, it's going to claim the lease we started here and take over the maintenance functions like those. It will also terminate that lease as needed.
*/
func MaybeStartSync(reverseConn *net.Conn) error {
	/*=================================================
	=            Requesting outbound lease            =
	=================================================*/
	allowed, _, _ := dispatch.OutboundAllowed(api.Address{}, reverseConn)
	if !allowed {
		api.SendReverseOpenStatusRefused(reverseConn)
		// ^ The connection is closed within this.

		// time.Sleep(10 * time.Second)
		// (*reverseConn).Close()
		errMessage := fmt.Sprintf("We don't have an open outbound lease to respond to this reverse open request, so we declined it. Connection: %#v", reverseConn)
		logging.Logf(1, errMessage)
		return errors.New(errMessage)
	}
	/*=====  End of Requesting outbound lease  ======*/
	err := dispatch.Sync(api.Address{}, []string{}, reverseConn)
	// ^ The connection is closed within this.

	// time.Sleep(10 * time.Second)
	// (*reverseConn).Close()
	return err
}

func (t *TCPMimServer) HandleConn(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := bufio.NewReader(conn)
	// We only have one message for now. We can skip the protocol parser.
	msg := make([]byte, 9)
	_, err := buf.Read(msg)
	if err != nil {
		logging.Logf(0, "TCPMIMServer: HandleConn: Error reading: %v", err.Error())
		conn.Close()
		return
	}
	if tcpmim.ParseMimMessage(msg) != tcpmim.ReverseOpenRequest {
		logging.Logf(0, "TCPMIMServer: HandleConn: Not a known TCPMim message. Message: %v", string(msg))
		conn.Close()
	}
	// logging.Logf(0, "DEBUG TCPMIM: %v, as bytes: %v, source: %v", string(msg), msg, conn.RemoteAddr().String())
	// This is a reverse open request.
	logging.Logf(0, "TCPMIMServer: HandleConn: This is a TCPMim reverse open request from %v. Going into MaybeSync.", conn.RemoteAddr().String())
	err2 := MaybeStartSync(&conn)
	if err2 != nil {
		logging.Logf(0, "TCPMIMServer: HandleConn: MaybeStartSync errored out. Error: %v", err2)
		return
	}
	logging.Logf(0, "TCPMIMServer: HandleConn: Reverse open sync completed successfully.")
}
