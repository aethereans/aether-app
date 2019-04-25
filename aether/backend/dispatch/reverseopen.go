package dispatch

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/tcpmim"
	"aether-core/aether/services/toolbox"
	"errors"
	"fmt"
	"net"
	"time"
)

/*
The way we open a reverse connection is that we open a conn to the local server, and we open a conn to the remote server, and after sending the TCPMim message to request reverse open, we pipe one conn to another.

							C1	(connToLocal)								C2 (connToRemote)
LOCAL SERVER <--> LOCAL END <PIPE> LOCAL END <--> REMOTE SERVER
^ Local Remote    ^ Local Local    ^ Local Local  ^ Remote Remote

Why? Because it allows us to inject data to be sent to the remote server (the reverse open request) and *then* set up the pipe . It also saves us from having to deal with SO_REUSEADDR and SO_REUSEPORT and their subtly different platform-specific implementations.
*/

func RequestInboundSync(host string, subhost string, port uint16) error {
	/*
		As of dev.6, Reverse opens are targeted at the port number - 1. Nodes run one HTTP server at the port, and a TCP server at port - 1 reserved for reverse opens.
	*/
	dev6plusPort := port - 1
	// dev.6+ code path
	logging.Logf(1, "Attempting to request inbound sync from remote: %s:%v", host, dev6plusPort)
	to := fmt.Sprint(host, ":", dev6plusPort)
	connToRemote, err := net.Dial("tcp4", to)
	if err != nil {
		errText := fmt.Sprintf("Request inbound sync failed while attempting to establish a connection to the remote. Error: %v", err)
		logging.Logf(1, errText)
		return errors.New(errText)
	}
	localSrvAddr := fmt.Sprint(":", globals.BackendConfig.GetExternalPort())
	connToLocal, err := net.Dial("tcp4", localSrvAddr)
	if err != nil {
		errText := fmt.Sprintf("Request inbound sync failed while attempting to establish a connection to the local server. Error: %v", err)
		logging.Logf(1, errText)
		return errors.New(errText)
	}
	// Set the values to transient config so that the server will be able to check if an incoming conn is a reverse conn.
	c1LocalLocalAddr, c1LocalLocalPort := toolbox.SplitHostPort(connToLocal.LocalAddr().String())
	globals.BackendTransientConfig.ReverseConnData.C1LocalLocalAddr = c1LocalLocalAddr
	globals.BackendTransientConfig.ReverseConnData.C1LocalLocalPort = c1LocalLocalPort
	mimMsg := tcpmim.MakeMimMessage(tcpmim.ReverseOpenRequest)
	fmt.Fprintf(connToRemote, string(mimMsg))
	logging.Logf(1, "Established pipe: (Local End) R: %v -> L: %v >[Pipe]> R: %v > L: %v (Remote End)",
		connToLocal.RemoteAddr().String(),
		connToLocal.LocalAddr().String(),
		connToRemote.LocalAddr().String(),
		connToRemote.RemoteAddr().String(),
	)
	start := time.Now()
	// Set timeouts to infinite - both are successful.
	pipe(connToRemote, connToLocal)
	// The remote will auto-close the connection, or the local server will, or it will just timeout on its own based on inactivity.
	elapsed := time.Since(start)
	fmt.Printf("Reverse conn took: %v\n", elapsed)
	es := ReverseConnInfo.GetEndStatus()
	switch es {
	case "SUCCESSFUL":
		// It worked out, the remote completed a sync, and told us so.
		return nil
	case "REFUSED":
		// The remote refused to reverse-connect to us.
		return errors.New("The remote connection request was refused.")
	case "FAILED", "":
		// The remote failed while doing a sync, or it never responded in the first place.
		return errors.New("The remote connection failed in sync, or had trouble establishing a connection.")
	default:
		// The server set the reverse connection end status into an unknown state.
		return errors.New(fmt.Sprintf("The local server set an end status to this reverse open connection, but it was not recognised. End status: %v", es))
	}
}

/*=====================================
=            Reverse scout            =
=====================================*/
/*
	This is where we retry the reverse connections with multiple remotes, should the first one fail.
*/

var (
	reverseScoutAttempts = 10
	// How many times scout will try to do a reverse open.
)

func reverseConnect(a api.Address) error {
	return RequestInboundSync(string(a.Location), string(a.Sublocation), a.Port)
}

func ReverseScout() error {
	logging.Log(2, "ReverseScout triggers.")
	addrs, err := findOnlineNodesV2(0, -1, -1, nil, true)
	if err != nil {
		errStr := fmt.Sprintf("ReverseScout: address search failed. Error: %v", err)
		logging.Log(1, errStr)
		return errors.New(errStr)
	}
	if len(addrs) == 0 {
		logging.Log(2, "ReverseScout got no unconnected addresses. Bailing.")
		return errors.New("ReverseScout got no unconnected addresses. Bailing.")
	}
	attempts := 0
	for k, _ := range addrs {
		if addrs[k].LocationType == 3 {
			continue
			// If it's an URL (type=3), we don't attempt to do a reverse open.
		}
		if attempts > reverseScoutAttempts {
			break
		}
		attempts++
		logging.Logf(1, "ReverseScout: Reverse connection attempt #%v.", attempts)
		err := reverseConnect(addrs[k])
		if err != nil {
			logging.Logf(1, "ReverseScout: ReverseConnect failed. Error: %v", err)
			continue
		}
		return nil
	}

	allFailedError := errors.New(fmt.Sprintf("ReverseScout failed because all %v nodes we have tried has failed.", attempts))
	logging.Logf(1, "ReverseScout: Connect failed. Error: %#v", allFailedError)
	return errors.New(fmt.Sprintf("ReverseScout: Connect failed. Error: %#v", allFailedError))
}

/*=====  End of Reverse scout  ======*/

/*=================================================
=            Reverse connection status            =
=================================================*/
/*
	Server will update this status based on the status the remote tells us.
	We read this status when the TCP socket is terminated. When we read this, we also clear it for the next cycle.
*/

var (
	ReverseConnInfo rci
)

type rci struct {
	EndStatus string // REFUSED, FAILED, SUCCESSFUL
}

func (i *rci) SetEndStatus(endStatus string) {
	i.EndStatus = endStatus
}

func (i *rci) GetEndStatus() string {
	es := i.EndStatus
	i.EndStatus = ""
	return es
}

/*=====  End of Reverse connection status  ======*/

/*============================================
=            Plumbing - literally            =
============================================*/

func connToChan(conn net.Conn) chan []byte {
	c := make(chan []byte)
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy just so so it doesn't change while read by the recipient
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()
	return c
}

func pipe(c1, c2 net.Conn) {
	chan1 := connToChan(c1)
	chan2 := connToChan(c2)
	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			} else {
				// In the case of a r/w, update expiries.
				updateDeadlines(&c1, &c2)
				// c1.SetDeadline(time.Now().Add(1 * time.Minute))
				// c2.SetDeadline(time.Now().Add(1 * time.Minute))
				c2.Write(b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				return
			} else {
				// In the case of a r/w, update expiries.
				updateDeadlines(&c1, &c2)
				// c1.SetDeadline(time.Now().Add(1 * time.Minute))
				// c2.SetDeadline(time.Now().Add(1 * time.Minute))
				c1.Write(b2)
			}
		}
	}
}

var lastDeadlineUpdate int64

func updateDeadlines(c1, c2 *net.Conn) {
	if lastDeadlineUpdate < time.Now().Add(-30*time.Second).Unix() {
		(*c1).SetDeadline(time.Now().Add(1 * time.Minute))
		(*c2).SetDeadline(time.Now().Add(1 * time.Minute))
	}
}

/*=====  End of Plumbing - literally  ======*/
