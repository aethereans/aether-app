// Backend > Routines > Check
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	// "aether-core/aether/backend/responsegenerator"
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/ca"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/toolbox"
	// "aether-core/aether/services/logging"
	// tb "aether-core/aether/services/toolbox"
	// "aether-core/aether/services/verify"
	"errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/fatih/color"
	"net"
	// "strconv"
	"strings"
	"time"
)

// Check is the short routine that reaches out to a node to see if it is online, and if so, pull the node data. This returns an updated api.Address object. Sync logic uses check as a starting point.
func Check(a api.Address, reverseConn *net.Conn) (api.Address, bool, api.ApiResponse, bool, error) {
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check is being called for: %s:%d\n", a.Location, a.Port)
	// }
	NODE_STATIC := false
	directlyConnectible := false
	/*
		If the port is 0, the node doesn't want us to connect. Skip directly.
	*/
	if a.Port == 0 {
		return a, false, api.ApiResponse{}, false, nil
	}
	/*
	   - Status GET to check if the node is online.
	*/
	_, err := api.Fetch(string(a.Location), string(a.Sublocation), a.Port, "status", "GET", []byte{}, reverseConn)
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check error: %#v\n", err)
	// }
	if err != nil && (strings.Contains(err.Error(), "Client.Timeout exceeded") || strings.Contains(err.Error(), "i/o timeout")) {
		// CASE: NO RESPONSE
		// The node is offline. It can actually be offline or just too slow, but for our purposes it's the same. The timeout can be set from globals.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, directlyConnectible, err
	} else if err != nil {
		// CASE: CONNECTION REFUSED
		// This is where 'connection refused' would go.
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, directlyConnectible, err
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through first pass for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - The node is online. Ask for node data.
	   (This is a legitimate user of GetPageRaw because the other entities that use Check sometimes need NodeId and other fields within it.)
	*/
	apiResp, err2 := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{}, reverseConn)
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check error: %#v\n", err2)
	// }
	if err2 != nil {
		return api.Address{}, NODE_STATIC, apiResp, directlyConnectible, err2
	}
	permissible, permissionErr := permissible(apiResp, a)
	if !permissible {
		return api.Address{}, NODE_STATIC, api.ApiResponse{}, directlyConnectible, permissionErr
	}
	// if apiResp.NodeId == api.Fingerprint(globals.BackendConfig.GetNodeId()) {
	/*
	   This node is using the same NodeId as we do. This is, in most cases, a node connecting to itself over a loopback interface. Most router will not allow their own address to be pinged from within network, but in testing and in other rare occasions this can happen.
	*/
	// 	return api.Address{}, NODE_STATIC, api.ApiResponse{}, directlyConnectible, errors.New(fmt.Sprintf("This node appears to have found itself through a loopback interface, or via calling its own IP. IP: %s:%d", a.Location, a.Port))
	// }
	if apiResp.Address.Type == 255 || // static node
		apiResp.Address.Type == 254 || // static boostrap
		apiResp.Address.Type == 253 { // static CA
		NODE_STATIC = true
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through GET for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - If the node is not static, present yourself.
	*/
	var postApiResp api.ApiResponse
	if !NODE_STATIC {
		// apiReq := responsegenerator.GeneratePrefilledApiResponse()
		apiReq := api.ApiResponse{}
		apiReq.Prefill()
		signingErr := apiReq.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, directlyConnectible, signingErr
		}
		apiReq.CreatePoW()
		reqAsJson, jsonErr := apiReq.ToJSON()
		// reqAsJson, jsonErr := responsegenerator.ConvertApiResponseToJson(&apiReq)
		if jsonErr != nil {
			return api.Address{}, NODE_STATIC, apiResp, directlyConnectible, jsonErr
		}
		var err3 error
		postApiResp, err3 = api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "POST", reqAsJson, reverseConn) // Raw call instead of regular one because we need access to the inbound remote timestamp.
		// if a.Location == "127.0.0.1" {
		// 	fmt.Printf("Check error: %#v\n", err3)
		// }
		if err3 != nil {
			// Mind that this can fail for verification failure also (if 4 entities in a page fails verification, the page fails verification. This is a page that actually has entities.)
			return api.Address{}, NODE_STATIC, apiResp, directlyConnectible, errors.New(fmt.Sprintf("Getting POST Endpoint in Check() routine for this entity type failed. Endpoint type: %s, Error: %s", "node", err3))
		}
		// Guard against the case where the node says something on GET, and something else on POST.
		if postApiResp.NodePublicKey != apiResp.NodePublicKey {
			return api.Address{}, NODE_STATIC, apiResp, directlyConnectible, errors.New(fmt.Sprintf("This node is declaring itself to be different nodes in GET and POST requests. Endpoint type: %s, GET NodePublicKey: %v, POST NodePublicKey: %v", "node", apiResp.NodePublicKey, postApiResp.NodePublicKey))
		}
	}
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through POST for: %s:%d\n", a.Location, a.Port)
	// }
	/*
	   - Collect the newly built address data.
	*/
	var addr api.Address
	var lastSuccessfulPing api.Timestamp
	if NODE_STATIC {
		addr = apiResp.Address
		lastSuccessfulPing = apiResp.Timestamp
	} else {
		addr = postApiResp.Address // addr is what comes from remote, a is local.
		lastSuccessfulPing = api.Timestamp(time.Now().Unix())
	}
	addr = *insertFirstPartyAddressData(&addr, &a, lastSuccessfulPing, reverseConn)
	// if a.Location == "127.0.0.1" {
	// 	fmt.Printf("Check made it through all. Resulting address: %#v\n", addr)
	// }
	if reverseConn != nil {
		directlyConnectible = checkDirectConnectivity(addr, apiResp.NodePublicKey)
	} else {
		directlyConnectible = true
	}
	return addr, NODE_STATIC, apiResp, directlyConnectible, nil
}

/*
//////////
Internal functions
//////////
*/

func permissible(apiResp api.ApiResponse, a api.Address) (bool, error) {
	if apiResp.NodeId == api.Fingerprint(globals.BackendConfig.GetNodeId()) {
		/*
		   This node is using the same NodeId as we do. This is, in most cases, a node connecting to itself over a loopback interface. Most router will not allow their own address to be pinged from within network, but in testing and in other rare occasions this can happen.
		*/
		return false, errors.New(fmt.Sprintf("This node appears to have found itself through a loopback interface, or via calling its own IP. IP: %s:%d", a.Location, a.Port))
	}
	if apiResp.Address.Type == 253 || apiResp.Address.Type == 4 {
		// This is a CA node. Make sure that it is a CA we trust, otherwise, terminate the connection.
		if !ca.IsTrustedCAKeyByPK(apiResp.NodePublicKey) {
			return false, errors.New(fmt.Sprintf("This is a CA node, and it is a CA that we do not trust. IP: %s:%d", a.Location, a.Port))
		}
	}
	return true, nil
}

// checkDirectConnectivity checks whether the node given is publicly connectable, and is the node it says it is.
func checkDirectConnectivity(a api.Address, npk string) bool {
	apiResp, err := api.GetPageRaw(string(a.Location), string(a.Sublocation), a.Port, "node", "GET", []byte{}, nil)
	if err != nil {
		// logging.Log(2, "checkDirectConnectivity errored out. Error: %v", err)
		return false
	}
	if apiResp.NodePublicKey == npk {
		return true
	}
	return false
}

// insertFirstPartyAddressData inserts the first-party data we know about this address.
func insertFirstPartyAddressData(inboundAddrPtr *api.Address, localAddrPtr *api.Address, lastSuccessfulPing api.Timestamp, reverseConn *net.Conn) *api.Address {
	addr := *inboundAddrPtr
	if reverseConn != nil {
		// Reverse open. //
		/*
			Heads up, in the reverse open case, these data points are not guaranteed to be reachable at any point in time, the node might be behind NAT be permanently unreachable through this. Do not save these without checking connectivity first.

			In addition, the Port data is **not** trusted. (Though since the host is, this can't be used to DDoS some other node.) Because the port information coming from the reverse connection is the port of that reverse connection, and it is guaranteed that it is wrong since the remote has to do it from a port that its server is not using. Therefore, if we take in the port number from the connection, it won't work. Since we don't actually have any direct data on what the port of the server is, we have to accept data from the remote that declares its server port.

			This all to say, make sure in the case of a reverse connection, this data should ONLY be admitted if it is actually connectable. Try connecting first, and save it only if it works outside the reverse connection. There is no point in saving data that pertains to that specific reverse connection only and won't work outside after that connection is terminated.
		*/
		host, _, _ := net.SplitHostPort((*reverseConn).RemoteAddr().String())
		addr.Location = api.Location(host)
		addr.Port = inboundAddrPtr.Port // Heads up, untrusted data entry, check connectivity before using / saving. See commentary above for details.
	} else {
		// Normal connection initiated by us. //
		addr.Location = localAddrPtr.Location // We know this to be true, because we just connected to it through a.Location. If the remote says it's a different IP, it's lying.
		addr.Sublocation = localAddrPtr.Sublocation
		// Determine IP type from the local address we just used to connect to this remote.
		addr.Port = localAddrPtr.Port // Because we just connected to this port and it worked. If the remote says it's a different port, it's lying.
	}
	if IsIPv4(addr) {
		addr.LocationType = 4
	} else if IsIPv6(addr) {
		addr.LocationType = 6
	} else {
		addr.LocationType = 3 // = URL
	}
	addr.LastSuccessfulPing = lastSuccessfulPing
	addr.EntityVersion = globals.BackendTransientConfig.EntityVersions.Address
	return &addr
}

func IsIPv4(addr api.Address) bool {
	return toolbox.IsIPv4String(string(addr.Location))
}

func IsIPv6(addr api.Address) bool {
	return toolbox.IsIPv6String(string(addr.Location))
}
