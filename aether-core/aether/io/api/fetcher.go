// API > Fetcher
// This file implements the methods that fetch the data from remotes. Mind that this is only for fetching, the lifecycle and the checks on whether the remote node is available for fetching is handled in dispatch. It deals with getting the data in, it does not deal with decisions on what actions to take (intro, update, search), neither it does with what method to use (get, post).

package api

import (
	"aether-core/aether/services/fingerprinting"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	// "github.com/jmoiron/sqlx"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	// "context"
	// "github.com/davecgh/go-spew/spew"
	// "reflect"
	// "github.com/libp2p/go-reuseport"
)

// Exists checks whether a given item exists in the current DB. This is here because we cannot import persistence due to import cycle being formed, and this is the only place this is being used.
func ExistsInDB(entityType string, fp Fingerprint, lu Timestamp) bool {
	if globals.BackendTransientConfig.ShutdownInitiated {
		return true
	}
	var tableName string
	var result bool
	if entityType == "board" {
		tableName = "Boards"
	} else if entityType == "thread" {
		tableName = "Threads"
	} else if entityType == "post" {
		tableName = "Posts"
	} else if entityType == "vote" {
		tableName = "Votes"
	} else if entityType == "key" {
		tableName = "PublicKeys"
	} else if entityType == "truststate" {
		tableName = "Truststates"
	} else {
		logging.Log(1, fmt.Sprintf("ExistsInDB does not support the entity type you provided. You provided: %s", entityType))
		return false
	}
	qStr := fmt.Sprintf("SELECT count(1) FROM %s WHERE (Fingerprint = ? AND LastUpdate = ?)", tableName)
	rows, err := globals.DbInstance.Queryx(qStr, fp, lu)
	defer rows.Close() // In case of premature exit.
	if err != nil {
		logging.Log(1, fmt.Sprintf("ExistsInDB errored out. Error: %s\n", err))
		return false
	}
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			logging.Log(1, fmt.Sprintf("ExistsInDB errored out. Error: %s\n", err))
			return false
		}
	}
	rows.Close()
	return result
}

func InsertApiResponseToResponse(response Response, apiresp ApiResponse) Response {
	response.Boards = apiresp.ResponseBody.Boards
	response.Threads = apiresp.ResponseBody.Threads
	response.Posts = apiresp.ResponseBody.Posts
	response.Votes = apiresp.ResponseBody.Votes
	response.Keys = apiresp.ResponseBody.Keys
	response.Truststates = apiresp.ResponseBody.Truststates
	response.Addresses = apiresp.ResponseBody.Addresses

	response.BoardIndexes = apiresp.ResponseBody.BoardIndexes
	response.ThreadIndexes = apiresp.ResponseBody.ThreadIndexes
	response.PostIndexes = apiresp.ResponseBody.PostIndexes
	response.VoteIndexes = apiresp.ResponseBody.VoteIndexes
	response.KeyIndexes = apiresp.ResponseBody.KeyIndexes
	response.TruststateIndexes = apiresp.ResponseBody.TruststateIndexes
	response.AddressIndexes = apiresp.ResponseBody.AddressIndexes

	response.BoardManifests = apiresp.ResponseBody.BoardManifests
	response.ThreadManifests = apiresp.ResponseBody.ThreadManifests
	response.PostManifests = apiresp.ResponseBody.PostManifests
	response.VoteManifests = apiresp.ResponseBody.VoteManifests
	response.KeyManifests = apiresp.ResponseBody.KeyManifests
	response.TruststateManifests = apiresp.ResponseBody.TruststateManifests
	response.AddressManifests = apiresp.ResponseBody.AddressManifests

	response.CacheLinks = apiresp.Results

	if response.MostRecentSourceTimestamp < apiresp.Timestamp {
		response.MostRecentSourceTimestamp = apiresp.Timestamp
	}
	return response
}

func concatResponses(response Response, response2 Response) Response {
	/*
		This is how append works: (the first slice to be added, you don't need "...")
		test1 := []string{"a","b"}
			test2 := []string{"c","d"}
			test3 := []string{}

			test3 = append(test1, test2...)
			fmt.Println(test3)
			-> [a,b,c,d]
	*/
	var resp Response
	resp.Boards = append(
		response.Boards, response2.Boards...)
	resp.Threads = append(
		response.Threads, response2.Threads...)
	resp.Posts = append(
		response.Posts, response2.Posts...)
	resp.Votes = append(
		response.Votes, response2.Votes...)
	resp.Keys = append(
		response.Keys, response2.Keys...)
	resp.Truststates = append(
		response.Truststates, response2.Truststates...)
	resp.Addresses = append(
		response.Addresses, response2.Addresses...)

	resp.BoardIndexes = append(
		response.BoardIndexes, response2.BoardIndexes...)
	resp.ThreadIndexes = append(
		response.ThreadIndexes, response2.ThreadIndexes...)
	resp.PostIndexes = append(
		response.PostIndexes, response2.PostIndexes...)
	resp.VoteIndexes = append(
		response.VoteIndexes, response2.VoteIndexes...)
	resp.KeyIndexes = append(
		response.KeyIndexes, response2.KeyIndexes...)
	resp.TruststateIndexes = append(
		response.TruststateIndexes, response2.TruststateIndexes...)
	resp.AddressIndexes = append(
		response.AddressIndexes, response2.AddressIndexes...)

	resp.BoardManifests = append(
		response.BoardManifests, response2.BoardManifests...)
	resp.ThreadManifests = append(
		response.ThreadManifests, response2.ThreadManifests...)
	resp.PostManifests = append(
		response.PostManifests, response2.PostManifests...)
	resp.VoteManifests = append(
		response.VoteManifests, response2.VoteManifests...)
	resp.KeyManifests = append(
		response.KeyManifests, response2.KeyManifests...)
	resp.TruststateManifests = append(
		response.TruststateManifests, response2.TruststateManifests...)
	resp.AddressManifests = append(
		response.AddressManifests, response2.AddressManifests...)

	resp.CacheLinks = append(
		response.CacheLinks, response2.CacheLinks...)

	if response.MostRecentSourceTimestamp < response2.MostRecentSourceTimestamp {
		resp.MostRecentSourceTimestamp = response2.MostRecentSourceTimestamp
	} else {
		resp.MostRecentSourceTimestamp = response.MostRecentSourceTimestamp
	}
	return resp
}

// makeProxyDialer returns blank if proxy support is not enabled. Heads up, proxy doesn't have support for setting timeouts and keepalive.
func makeProxyDialer() proxy.Dialer {
	var pdialer proxy.Dialer
	var proxyErr error
	if globals.BackendConfig.GetSOCKS5ProxyEnabled() {
		auth := proxy.Auth{
			User:     globals.BackendConfig.GetSOCKS5ProxyUsername(),
			Password: globals.BackendConfig.GetSOCKS5ProxyPassword(),
		}
		if len(auth.User) > 0 || len(auth.Password) > 0 {
			pdialer, proxyErr = proxy.SOCKS5("tcp4", globals.BackendConfig.GetSOCKS5ProxyAddress(), &auth, proxy.Direct)
		} else {
			pdialer, proxyErr = proxy.SOCKS5("tcp4", globals.BackendConfig.GetSOCKS5ProxyAddress(), nil, proxy.Direct)
		}
		if proxyErr != nil {
			logging.LogCrash(fmt.Sprintf(
				"Cannot connect to the proxy. This was the error: %v \nInstead of connecting with no proxy, the application will terminate so as not to create a false impression of being connected via the proxy. If you want to revert back to the default no-proxy setting, go to %s/backend/backend_config.json and set SOCKS5ProxyEnabled key to 'false'.",
				proxyErr, globals.BackendConfig.GetUserDirectory()))
		}
	}
	return pdialer
}

func generateDialFunc(conn *net.Conn) func(network, address string) (net.Conn, error) {
	if conn != nil {
		// This is a reverse connection we decided to accept. We are using the connection opened by the remote to make an outbound request.
		dialFunc := func(network, address string) (net.Conn, error) {
			return *conn, nil
		}
		return dialFunc
	}
	if globals.BackendConfig.GetSOCKS5ProxyEnabled() {
		// We have a proxy. We'll use the proxy to make an outbound request.
		dialer := makeProxyDialer()
		return dialer.Dial
	}
	// We have neither a reverse open nor a proxy. Return the regular dialer we created at init.
	return d.Dial
}

// Basic, reusable instances of transport and client.

var d *net.Dialer
var t *http.Transport
var c *http.Client

// postInit initiates on the first run of Fetch. These cannot be set in init because init runs before the backendconfig is ready.
func postInit() {
	t = &http.Transport{
		TLSHandshakeTimeout: globals.BackendConfig.GetTLSHandshakeTimeout(),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Is that not insecure? Not in this specific case. We're not using TLS as a means of identifying the remote, just for encrypting the pipe. See note at tlscerts library.
		},
	}
	c = &http.Client{
		Transport: t,
		Timeout:   globals.BackendConfig.GetConnectionTimeout(),
	}
	d = &net.Dialer{
		Timeout:   globals.BackendConfig.GetTCPConnectTimeout(),
		KeepAlive: 30 * time.Second,
	}
	fetchRanBefore = true
	protv = globals.BackendConfig.GetProtURLVersion()
}

var fetchRanBefore bool
var protv string

// Fetch is the most basic access method. It returns bytes. This should almost never be called directly outside this package.
func Fetch(host string, subhost string, port uint16, location string, method string, postBody []byte, reverseConn *net.Conn) ([]byte, error) {
	if !fetchRanBefore {
		postInit()
	}
	if reverseConn != nil {
		// If this is a reverse connection, it's operating outside a dialer, a transport, or a client. This means we have to manually handle the timeouts.
		// This timeout is overly generous, because it's not about how much time a remote needs, it's more about if something goes wrong, it will automatically be able to clear the connection and not have it stuck in limbo forever.
		(*reverseConn).SetDeadline(time.Now().Add(10 * time.Minute))
	}
	t.Dial = generateDialFunc(reverseConn)
	var prot string
	if globals.BackendTransientConfig.TLSEnabled {
		prot = "https://"
	} else {
		prot = "http://"
	}
	hostAsString := string(host)
	// If IpV6, we need to bracket the address.
	if toolbox.IsIPv6String(hostAsString) {
		hostAsString = fmt.Sprintf("[%s]", hostAsString)
	}
	var fullLink string
	if len(subhost) > 0 {
		fullLink = fmt.Sprint(
			prot, hostAsString, ":", strconv.Itoa(int(port)), "/", subhost, "/", protv, "/", location)
	} else {
		fullLink = fmt.Sprint(
			prot, hostAsString, ":", strconv.Itoa(int(port)), "/", protv, "/", location)
	}
	// logging.Logf(3, "Fetch is being called for: %v", fullLink)
	// if strings.Contains(fullLink, "127.0.0.1") {
	// 	logging.Log(2, fmt.Sprintf("Fetch is being called for the URL: %s. ReverseConn: %v", fullLink, reverseConn != nil))
	// }
	// if strings.Contains(fullLink, "127.0.0.1") {
	// 	logging.Log(1, fmt.Sprintf("Fetch is being called for the URL: %s", fullLink))
	// }
	var err error
	var resp *http.Response
	if method == "GET" {
		resp, err = c.Get(fullLink)
	} else if method == "POST" {
		resp, err = c.Post(fullLink, "application/json", bytes.NewReader(postBody))
	} else {
		defer resp.Body.Close()
		return []byte{}, errors.New("Unsupported HTTP method. Available methods are: GET, POST")
	}
	if err != nil {
		if resp != nil {
			defer resp.Body.Close()
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "target machine actively refused it") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"The host refused the connection. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Timeout exceeded. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"i/o timeout. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "connection reset by peer") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Connection reset by peer. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "EOF") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"The remote crashed or is shutting down. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "server gave HTTP response to HTTPS client") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Remote gave HTTP response to HTTPS client. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "tls: oversized record received with length") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Remote gave HTTP response to HTTPS client. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "remote error: tls: handshake failure") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Remote is unavailable. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "TLS handshake timeout") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Remote timed out while doing TLS handshake. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "can't assign requested address") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Attempted to dial an undialable address. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location))
		} else if strings.Contains(err.Error(), "use of closed network connection") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"The connection we attempted to use was already closed. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location,
					", ReverseConn: ", fmt.Sprintf("%#v", reverseConn)))
		} else if strings.Contains(err.Error(), "invalid argument") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"Something bad happened in the OS level under us and the OS sent us an 'invalid argument'. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location,
					", ReverseConn: ", fmt.Sprintf("%#v", reverseConn)))
		} else if strings.Contains(err.Error(), "no route to host") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"No route to host. Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location,
					", ReverseConn: ", fmt.Sprintf("%#v", reverseConn)))
		} else if strings.Contains(err.Error(), "network is unreachable") || strings.Contains(err.Error(), "A socket operation was attempted to an unreachable network") {
			return []byte{}, errors.New(
				fmt.Sprint(
					"network is unreachable (Is this an IPv6 address?). Host:", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location,
					", ReverseConn: ", fmt.Sprintf("%#v", reverseConn)))
		} else {
			// FUTURE: This should be very rare, since we cover the major cases above. If this happens, we should log.
			errorStr := fmt.Sprintf("An unknown error occurred in Fetch. This fetch is terminated. Error: %v \nEnvironment:\n Host: %v, Subhost: %v, Port: %v, Location: %v, Method: %v, POSTBody: %v, ReverseConn: %s", err, host, subhost, port, location, method, postBody, fmt.Sprintf("%#v", reverseConn))
			logging.Logf(1, errorStr)
			// fmt.Println("Fatal error in api.Fetch. Quitting.")
			// fmt.Printf("Environment: Host: %v, Subhost: %v, Port: %v, Location: %v, Method: %v, POSTBody: %v, ReverseConn: %s\n", host, subhost, port, location, method, postBody, fmt.Sprintf("%#v", reverseConn))
			// fmt.Println(err)
			// logging.LogCrash(err)
			return []byte{}, errors.New(errorStr)
		}
	}
	// server gave HTTP response to HTTPS client
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		limitedReader := &io.LimitedReader{
			R: resp.Body,
			N: int64(globals.BackendConfig.GetMaxInboundPageSizeKb() * 1000)}
		// *1000 because kb > b
		body, err := ioutil.ReadAll(limitedReader)
		if err != nil {
			// logging.LogCrash(err)
			logging.Logf(1, "Fetch error: %v", err)
		}
		return body, nil
	} else if resp.StatusCode == 501 {
		/*
			HTTP 501: Not implemented
			This usually happens when you try to do a non-static sync on a node that is static, thus isn't supporting POST requests.
		*/
		return []byte{}, errors.New(
			fmt.Sprint(
				"This remote does not implement this HTTP method. Likely this was an attempt to do a live sync with a static node. Host:", host,
				", Subhost: ", subhost,
				", Port: ", port,
				", Location: ", location,
				", Request method: ", method,
			))
	} else if resp.StatusCode == 429 {
		/*
			HTTP 429: Too many requests
			This happens when the remote node is alive and present, but has received more requests than it is willing to serve. W
		*/
		errString := fmt.Sprint(
			"This remote has received too many requests. Wait and try again later. Received status code: 429. Host:", host,
			", Subhost: ", subhost,
			", Port: ", port,
			", Location: ", location,
			", Request method: ", method,
		)
		logging.Logf(2, errString)
		return []byte{}, errors.New(errString)
	} else {
		logging.Logf(2, "FULL LINK IN FETCH FOR THIS FAILED REQUEST: \n%s\nStatus Code: %v, Method: %v\n", fullLink, resp.StatusCode, method)
		return []byte{}, errors.New(
			fmt.Sprint(
				"Non-200 status code returned from Fetch. Received status code: ", resp.StatusCode,
				", Host: ", host,
				", Subhost: ", subhost,
				", Port: ", port,
				", Location: ", location,
				", Method: ", method))
	}
	return []byte{}, errors.New("This should never have happened.")
}

// GetPageRaw returns a raw page from the cache. This returns the entire page, not just the data. This is useful for functions that need to be aware of the page's metadata.
func GetPageRaw(host string, subhost string, port uint16, location string, method string, postBody []byte, reverseConn *net.Conn) (ApiResponse, error) {
	var apiresp ApiResponse
	result, err := Fetch(host, subhost, port, location, method, postBody, reverseConn)
	if err != nil {
		return apiresp, err
	}
	err2 := json.Unmarshal(result, &apiresp)
	if err2 != nil {
		return apiresp, errors.New(
			fmt.Sprint(
				"The JSON that arrived over the network is malformed. JSON: ", string(result),
				", Host: ", host,
				", Subhost: ", subhost,
				", Port: ", port,
				", Location: ", location))
	}
	// Map over everything you have.
	// if method == "POST" {
	// 	logging.Log(2, fmt.Sprintf("We've made a POST request to the endpoint %s and this was its body: %#v", location, string(postBody)))
	// }
	// if method == "POST" {
	// 	apiresp.Dump() // let's see
	// }
	pageVerified, err := apiresp.VerifySignature() // If signature check is disabled, this will always return true.
	if err != nil {
		return ApiResponse{}, fmt.Errorf("Page signature verification failed with an error. Error: %s", err)
	}
	if !pageVerified {
		return ApiResponse{}, errors.New("Page signature verification failed. The signature does not match.")
	}
	if len(apiresp.NodePublicKey) > 0 {
		apiresp.NodeId = Fingerprint(fingerprinting.Create(apiresp.NodePublicKey))
	} else {
		/*
			This makes it more obvious that the given field in the database is catchall for all nodes without a node public key. This should never happen in production because by default this check is enabled, and can only be disabled via a command line flag, which forces the app into read-only configs mode.
		*/
		apiresp.NodeId = "NODEID FOR NODE(S) WITH EMPTY NODEPUBLICKEY"
	}
	errs := apiresp.Verify()
	if len(errs) == 1 && strings.Contains(errs[0].Error(), "This ApiResponse failed the boundary check") {
		return ApiResponse{}, errs[0]
	}
	if len(errs) >= 3 {
		var errStrs []string

		for _, err := range errs {
			errStrs = append(errStrs, err.Error())
		}
		logging.Log(1, fmt.Sprintf("This page has 3 or more entities who has failed verification. Errors: %#v", errStrs))
		return ApiResponse{}, fmt.Errorf("This page has 3 or more entities who has failed verification")
	}
	return apiresp, nil
}

// GetPage gets a page from a cache. This returns the data on the provided page.
func GetPage(host string, subhost string, port uint16, location string, method string, postBody []byte, reverseConn *net.Conn) (Response, time.Duration, error) {
	var apiresp ApiResponse
	var response Response
	var start time.Time
	var elapsed time.Duration
	if method == "POST" {
		start = time.Now()
	}
	apiresp, err := GetPageRaw(host, subhost, port, location, method, postBody, reverseConn)
	if err != nil {
		return response, elapsed, err // elapsed is unset, set only on POST below.
	}
	if method == "POST" {
		elapsed = time.Since(start)
	}
	response = InsertApiResponseToResponse(response, apiresp)
	return response, elapsed, nil // elapsed potentially unset.
}

func countManifests(resp Response) {
	b, t, p, v, k, ts, a := 0, 0, 0, 0, 0, 0, 0
	for _, val := range resp.BoardManifests {
		b = b + len(val.Entities)
	}
	for _, val := range resp.ThreadManifests {
		t = t + len(val.Entities)
	}
	for _, val := range resp.PostManifests {
		p = p + len(val.Entities)
	}
	for _, val := range resp.VoteManifests {
		v = v + len(val.Entities)
	}
	for _, val := range resp.KeyManifests {
		k = k + len(val.Entities)
	}
	for _, val := range resp.TruststateManifests {
		ts = ts + len(val.Entities)
	}
	for _, val := range resp.AddressManifests {
		a = a + len(val.Entities)
	}
	logging.Logf(2, "generateHitlist manifestResponse result returned these: \nB: %v, T: %v, P: %v, V: %v, K: %v, TS: %v, A: %v", b, t, p, v, k, ts, a)
}

func generateHitlist(host string, subhost string, port uint16, location string, reverseConn *net.Conn) (map[int]bool, error) {
	start := time.Now()
	manifestResponse, err := getManifestOfCache(host, subhost, port, location, reverseConn)
	// logging.Logf(1, "Manifest Response: %#v", manifestResponse)
	if err != nil {
		return make(map[int]bool), fmt.Errorf("Error raised from GetManifestOfCache inside generateHitlist. Error: %s", err)
	}
	countManifests(manifestResponse)

	// Look at everything in the index and find the things that we want to pull. Page Number : bool pairs help us find which pages to hit.
	allPgs := make(map[int]bool)
BoardLoop:
	for key := range manifestResponse.BoardManifests {
		for _, val := range manifestResponse.BoardManifests[key].Entities {
			if !ExistsInDB("board", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.BoardManifests[key].Page)] = true
				continue BoardLoop
			}
		}
	}
ThreadLoop:
	for key := range manifestResponse.ThreadManifests {
		for _, val := range manifestResponse.ThreadManifests[key].Entities {
			if !ExistsInDB("thread", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.ThreadManifests[key].Page)] = true
				continue ThreadLoop
			}
		}
	}
PostLoop:
	for key := range manifestResponse.PostManifests {
		for _, val := range manifestResponse.PostManifests[key].Entities {
			if !ExistsInDB("post", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.PostManifests[key].Page)] = true
				continue PostLoop
			}
		}
	}
VoteLoop:
	for key := range manifestResponse.VoteManifests {
		for _, val := range manifestResponse.VoteManifests[key].Entities {
			if !ExistsInDB("vote", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.VoteManifests[key].Page)] = true
				continue VoteLoop
			}
		}
	}
KeyLoop:
	for key := range manifestResponse.KeyManifests {
		for _, val := range manifestResponse.KeyManifests[key].Entities {
			if !ExistsInDB("key", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.KeyManifests[key].Page)] = true
				continue KeyLoop
			}
		}
	}
TruststateLoop:
	for key := range manifestResponse.TruststateManifests {
		for _, val := range manifestResponse.TruststateManifests[key].Entities {
			if !ExistsInDB("truststate", val.Fingerprint, val.LastUpdate) {
				// Grab the whole page and insert into to-be-fetched queue, DB will remove useless stuff.
				allPgs[int(manifestResponse.TruststateManifests[key].Page)] = true
				continue TruststateLoop
			}
		}
	}
	elapsed := time.Since(start)
	logging.Logf(2, "GenerateHitlist V1 time spent: %#v\n", elapsed.String())
	return allPgs, nil
}

// GetCache returns an entire cache. This is useful to pull a cache from the remote. This is a single thread process, it does go through the pages in order.  We could bombard the remote with goroutines, but on a larger scale, that would be called a DDoS of the remote node, so we shouldn't do that.
func GetCache(host string, subhost string, port uint16, location string, isAddr bool, reverseConn *net.Conn) (Response, error) {
	var response Response
	// Get the first raw page (because we need to access pagination),
	pageResp, err := GetPageRaw(host, subhost, port, fmt.Sprint(location, "/0.json"), "GET", []byte{}, reverseConn)
	if err != nil && strings.Contains(err.Error(), "Received status code: 404") {
		return response, errors.New(
			fmt.Sprint(
				"The first page of the cache does not exist. This cache likely does not exist.",
				", Host: ", host,
				", Subhost: ", subhost,
				", Port: ", port,
				", Location: ", location))
	} else if err != nil {
		// If the first page is faulty, bail.
		return response, err
	}
	// And look at the page count, so we know how many times to iterate.
	pageCount := pageResp.Pagination.Pages
	// Convert this raw page response to page response data for merge.
	response = InsertApiResponseToResponse(response, pageResp)
	// Create a counter for missing pages. If 3 of them come one after another, bail.
	// Address specific
	addrCount := 0
	brokenPageCounter := 0
	// Iterate over all of the pages, starting from 1 (we already cleared the 0)
	for i := uint64(1); i <= pageCount; i++ { // Pagination starts from 0
		pageResp2, _, err := GetPage(host, subhost, port,
			fmt.Sprint(location, "/", i, ".json"), "GET", []byte{}, reverseConn)
		if err != nil {
			logging.Logf(2, "GetPage returned this error: Err: %v", err)
			brokenPageCounter++
			if brokenPageCounter >= 3 {
				logging.Log(1, fmt.Sprint(
					"3 or more broken pages, either missing or verification failures. Stopping the download of this cache.",
					", Host: ", host,
					", Subhost: ", subhost,
					", Port: ", port,
					", Location: ", location,
					", Last page number: ", i))
				return response, errors.New(
					fmt.Sprint(
						"3 or more broken pages, either missing or verification failures. Stopping the download of this cache.",
						", Host: ", host,
						", Subhost: ", subhost,
						", Port: ", port,
						", Location: ", location,
						", Last page number: ", i))
			}
		}
		// And save into the response.
		response = concatResponses(response, pageResp2)
		// Address specific
		if isAddr {
			addrCount = addrCount + len(pageResp2.Addresses)
			if addrCount >= 100 {
				// fmt.Println("This address cache download bailed beacuse we have enough addresses.")
				break
			}
		}
	}
	return response, nil
}

// mapEndpointToEndpointAddress generates the address that needs to be called for the endpoint that is being requested.
func mapEndpointToEndpointAddress(endpoint string) string {
	endpointsMap := map[string]string{
		"boards":      "c0/boards",
		"threads":     "c0/threads",
		"posts":       "c0/posts",
		"votes":       "c0/votes",
		"addresses":   "addresses", // Addresses is a mim entity, not a c0 entity.
		"keys":        "c0/keys",
		"truststates": "c0/truststates"}
	epAddress := endpointsMap[endpoint]
	// If we don't know which endpoint this is, attempt to call it directly.
	if epAddress == "" {
		epAddress = endpoint
	}
	return epAddress
}

// GetManifestGatedCache hits the manifests of the cache to determine which pages of the cache this computer needs to hit. This is useful in the case where you expect less than 50% of the cache will be downloaded. Mind that this adds a database check dependency (to know which one of these things we have at hand) and it will have to download the manifests for that cache, so it's a tradeoff.
func GetManifestGatedCache(host string, subhost string, port uint16, location string, endpoint string, reverseConn *net.Conn) (Response, error) {
	start := time.Now()
	allPgs, err := generateHitlist(host, subhost, port, location, reverseConn)
	if err != nil && strings.Contains(err.Error(), "Non-200 status code returned from Fetch") {
		// Manifest doesn't exist for this cache.
		logging.Log(1, fmt.Sprintf("This cache does not have a manifest. We'll be downloading the full cache. Host %s, Subhost: %s, Port: %d, Location: %s", host, subhost, port, location))
		resp, err2 := GetCache(host, subhost, port, location, endpoint == "addresses", reverseConn)
		return resp, err2
	} else if err != nil {
		logging.Log(1, fmt.Errorf("Error raised from generateHitlist inside GetManifestGatedCache. Error: %s", err))
		return Response{}, fmt.Errorf("Error raised from generateHitlist inside GetManifestGatedCache. Error: %s", err)
	}
	logging.Log(2, fmt.Sprintf("The pages we have to make a call to are: %#v\n", allPgs))

	// For each page we have for this post response, hit the main cache and gather the data.
	mainResp := Response{}
	for key := range allPgs {
		loc := fmt.Sprint(location, "/", key, ".json")
		logging.Log(2, fmt.Sprintf("Making a request to %s\n", loc))
		resp, _, err := GetPage(host, subhost, port, loc, "GET", []byte{}, reverseConn)
		if err != nil {
			return Response{}, err
		}
		mainResp = concatResponses(mainResp, resp)
	}
	elapsed := time.Since(start)
	logging.Logf(2, "GetManifestGatedCache V1 took this long: %v", elapsed.String())
	return mainResp, nil
}

// GetGETEndpoint returns an entire endpoint from the remote node.
func GetGETEndpoint(host string, subhost string, port uint16, endpoint string, lastCheckin Timestamp, reverseConn *net.Conn) (Response, error) {
	// This is where the mapping for an endpoint to its respective subprotocol folder is mapped. Below this level, you have to supply your own subprotocol string.
	logging.Log(2, fmt.Sprintf("GetGETEndpoint was called for the endpoint: %s", endpoint))
	epAddress := mapEndpointToEndpointAddress(endpoint)
	var response Response
	// Get raw page, because we need to access index links.
	result, err := getIndexOfEndpoint(host, subhost, port, epAddress, reverseConn)
	// Map the timestamp of the index onto the response we're generating, in case we might not have any caches (this can happen if our internal timestamp for this cache is newer than the last cache's timestamp.)
	response.MostRecentSourceTimestamp = result.MostRecentSourceTimestamp
	indexes := result.CacheLinks
	if err != nil {
		return response, errors.New(
			fmt.Sprint(
				"Get Endpoint failed because it couldn't get the index of the endpoint.",
				", Error: ", err,
				", Host: ", host,
				", Subhost: ", subhost,
				", Port: ", port,
				", Endpoint: ", endpoint))
	}
	// A broken cache can happen because the cache has underlying missing pages, or pages that has failed verification. At the level of the endpoint, it does not matter why the cache has failed, only that it failed. if there are enough failures, we bail.
	brokenCacheCounter := 0
	// Address specific
	addrCount := 0
	for _, val := range indexes {
		// If the cache does end after our last checkin timestamp, we want to read that cache.
		// ----------------- Why? -------------------------
		// Example:
		// Assume lastcheckin is 5
		// Assume caches are: 1-2, 2-3, 3-4, 4-5, 5-6, 6-7.
		// We want 4-5, 5-6, 6-7.
		// 5 6 7 (ends)
		// 5,6,7 > lastcheckin = true.
		// ------------------------------------------------
		if val.EndsAt >= lastCheckin {
			// Get the first page of the cache.
			cache, err := GetManifestGatedCache(host, subhost, port, fmt.Sprint(epAddress, "/", val.ResponseUrl), endpoint, reverseConn)
			// cache, err := GetCache(host, subhost, port,
			// 	fmt.Sprint(epAddress, "/", val.ResponseUrl), endpoint == "addresses", reverseConn)
			// cache, err := GetCache(host, subhost, port,
			// 	fmt.Sprint(epAddress, "/", val.ResponseUrl))
			response = concatResponses(response, cache)
			if err != nil {
				brokenCacheCounter++ // We never reset this within this endpoint call.
				if brokenCacheCounter >= 3 {
					return response, errors.New(
						fmt.Sprint(
							"3 or more cache failures in the same endpoint. Stopping the download of this endpoint.",
							", Error: ", err,
							", Host: ", host,
							", Subhost: ", subhost,
							", Port: ", port,
							", Endpoint: ", endpoint,
							", Cache link: ", fmt.Sprint(endpoint, "/", val.ResponseUrl)))
				}
			}
			// Address specific
			if endpoint == "addresses" {
				addrCount = addrCount + len(response.Addresses)
				if addrCount >= 100 {
					// fmt.Println("This address endpoint download bailed beacuse we have enough addresses.")
					break
				}
			}
		}

	}
	boardCount := len(response.Boards)
	threadCount := len(response.Threads)
	postCount := len(response.Posts)
	voteCount := len(response.Votes)
	addressCount := len(response.Addresses)
	keysCount := len(response.Keys)
	truststatesCount := len(response.Truststates)
	// logging.Log(1, fmt.Sprintf("Response for the endpoint %s was %#v\n", endpoint, response))
	logging.Log(2, fmt.Sprintf("GetGETEndpoint returned for the endpoint: %s. Number of items: Boards: %d, Threads: %d, Posts: %d, Votes: %d, Addresses: %d, Keys: %d, Truststates: %d", endpoint, boardCount, threadCount, postCount, voteCount, addressCount, keysCount, truststatesCount))

	return response, nil
}

/*
	HEADS UP:
	Post "Pingpong"
	We flip this to true if the POST response was 1 page. Why does it matter? Because our logs provide only candidate inputs, and if the result is only one page, it is not paginated, thus doesn't have indexes, thus it is not subject to 'exists in db' checks. By the virtue of getting the result, you've already downloaded it (i.e. you can't save bandwidth by not hitting the pages, it is already one page and you already have that page), so there is no point in doing this checks as DB will check it already anyway. But it creates a confusing side effect where it looks like the node has sent unnecessary data exclusively for the endpoints that were single-page because unneeded data wasn't filtered out by the local.

	In actuality, the unnecessary data is being sent by all endpoints in POST, but when there is a multipage post response, the local client can go through the index and not hit the pages it does not need. So they look like zeroes. But that check doesn't happen for single-page responses.

	Why is unnecessary data being sent?

	Let's look at an example.

	t0 1 > 0 : FirstS : N1 gets 31 new keys (total 62) ts:t0

	t1 1 > 0 : Resync : nothing  ts:t1

	t2 0 > 1 : FirstS : N0 gets 62 keys (the keys unique to 1, plus all 0's keys) 31 NEW KEYS @ N0 at ts:t2

	t3 1 > 0 : Resync : N1 gets 31 keys (because the ts N1 has for N0 is ts:t1, and all the keys N0 just got from N1 are at ts:t2)

	t4 0 > 1 : nothing ts:t4

	t5 1 > 0 : nothing ts:t5

	The crux is at t3, the data in N1 got transmitted to N0, and N0 transmitted it back to N1. This is because Mim connections are one-way. They are not two-way syncs. The nodes cannot retain information of which node they got the results from, because that node might have been reset, and might actually need that information. The local nodes cannot communicate what they have in their database, because that would be a privacy violation.

	Why is this not an issue?
	- It does not cost bandwidth except in the case of a single page. Single page means very little bandwidth is used. Any response of decent size will be multipart, and in multipart responses, the local node will check the index, discover that the things in the index are those that it already has, and it won't download the pages.
	- This only happens to POST responses, and not GET. POST responses are the 'tip of the spear', the usual way of process is that you exhaust the caches with GET first, and then get the delta from the end of the last cache to now via the POST response. Caches are generated frequently, the time range that the POST response will cover will be in the order of hours, not days / weeks. Most of the network traffic is GET.

*/

func GetPOSTEndpoint(host string, subhost string, port uint16, endpoint string, lastCheckin Timestamp, reverseConn *net.Conn) (Response, time.Duration, error) {
	// But before anything, we need to create the mapping for the endpoint URLs.
	endpointsMap := map[string]string{
		"boards":      "c0/boards",
		"threads":     "c0/threads",
		"posts":       "c0/posts",
		"votes":       "c0/votes",
		"keys":        "c0/keys",
		"truststates": "c0/truststates",
		"addresses":   "addresses",
	}
	apiReq := ApiResponse{}
	apiReq.Prefill()
	// Here, we need to insert the last sync timestamp into the post request, so that it will be gated appropriately.
	f := Filter{}
	f.Type = "timestamp"
	f.Values = []string{strconv.Itoa(int(lastCheckin)), strconv.Itoa(0)}
	apiReq.Filters = []Filter{f}
	signingErr := apiReq.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if signingErr != nil {
		return Response{}, 0, signingErr
	}
	apiReq.CreatePoW()
	reqAsJson, err := apiReq.ToJSON()
	if err != nil {
		return Response{}, 0, err
	}
	postResp, respDuration, err7 := GetPage(host, subhost, port, endpointsMap[endpoint], "POST", reqAsJson, reverseConn)
	if err7 != nil {
		return Response{}, respDuration, fmt.Errorf("Getting POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", endpoint, err7)
	}
	// presp := dbg_CheckIfExistsInDb(postResp)
	allResults := Response{}
	// Add entities embedded directly into the response to our response container, if any. A response can have both.
	// allResults = concatResponses(allResults, postResp)
	allResults.Insert(&postResp)
	// allResults.Insert(&presp)
	// If there are any cache links, one or multiple, read all of them, and insert.
	// Address-specific. We'll build the structure out if we need to do this for anything other than addresses.
	if endpoint == "addresses" && len(allResults.Addresses) >= 100 {
		return allResults, respDuration, nil
	}
	logging.Logf(1, "We received these response links from the remote to download. %v", postResp.CacheLinks)
	for _, clink := range postResp.CacheLinks {
		if clink.EndsAt > lastCheckin {
			// if true {
			// This cache ends after we have our sync timestamp with this remote. We can benefit from downloading this cachelink.
			logging.Logf(1, "Downloading %s from %s:%d", clink.ResponseUrl, host, port)
			postCacheResp, err8 := GetManifestGatedCache(host, subhost, port, fmt.Sprintf("responses/%s", clink.ResponseUrl), endpoint, reverseConn)
			// We're adding /responses/ because that's where the singular responses will be.
			if err8 != nil {
				return allResults, respDuration, fmt.Errorf("Getting Multi page POST Endpoint for this entity type failed. Endpoint type: %s, Error: %s", endpoint, err8)
			}
			logging.Logf(2, "postCacheResp counts at Get POST endpoint: \nB: %v, T: %v, P: %v, V: %v, K: %v, TS: %v, A: %v",
				len(postCacheResp.Boards),
				len(postCacheResp.Threads),
				len(postCacheResp.Posts),
				len(postCacheResp.Votes),
				len(postCacheResp.Keys),
				len(postCacheResp.Truststates),
				len(postCacheResp.Addresses))
			// Ends here, since we don't want to capture DB time.
			// allResults = concatResponses(allResults, postCacheResp)
			allResults.Insert(&postCacheResp)
			// Address-specific.
			if endpoint == "addresses" && len(allResults.Addresses) >= 100 {
				return allResults, respDuration, nil
			}
		} else {
			logging.Logf(1, "%s was skipped because this container's end is older than our last sync with this node.", clink.ResponseUrl)
		}
	}
	logging.Logf(2, "AllResults counts at the end of Get POST endpoint: \nB: %v, T: %v, P: %v, V: %v, K: %v, TS: %v, A: %v",
		len(allResults.Boards),
		len(allResults.Threads),
		len(allResults.Posts),
		len(allResults.Votes),
		len(allResults.Keys),
		len(allResults.Truststates),
		len(allResults.Addresses))
	return allResults, respDuration, nil
}

// GetRemoteNode downloads the entire remote node data by hitting all endpoints and all caches and all pages within them. This is the bootstrap function. This should be used when the local database is empty and the remote node is new. Never call this when the local database is not empty as that is fairly wasteful.
func GetRemoteNode(host string, subhost string, port uint16, reverseConn *net.Conn) (Response, error) {
	endpoints := []string{
		"boards", "threads", "posts", "votes", "addresses", "keys", "truststates"}
	var response Response
	for _, endpoint := range endpoints {
		resp, err := GetGETEndpoint(host, subhost, port, endpoint, 0, reverseConn)
		response = concatResponses(response, resp)
		if err != nil {
			// GetRemoteNode continues to work under all conditions. It won't stop the sequence for any errors.
			// NOOP for now.
		}
	}
	return response, nil // It won't communicate out any errors.
}

// This function normally does not run on single page responses, but we've added this for debugging purposes. This does have a performance penalty, but makes it so that the numbers visible at the log output are accurate numbers which go into the database.
// func dbg_CheckIfExistsInDb(postResp Response) Response {
// 	t := time.Now()
// 	r := Response{}
// 	for _, val := range postResp.Boards {
// 		if !ExistsInDB("board", val.Fingerprint, val.LastUpdate) {
// 			r.Boards = append(r.Boards, val)
// 		}
// 	}
// 	for _, val := range postResp.Threads {
// 		if !ExistsInDB("thread", val.Fingerprint, val.LastUpdate) {
// 			r.Threads = append(r.Threads, val)
// 		}
// 	}
// 	for _, val := range postResp.Posts {
// 		if !ExistsInDB("post", val.Fingerprint, val.LastUpdate) {
// 			r.Posts = append(r.Posts, val)
// 		}
// 	}
// 	for _, val := range postResp.Votes {
// 		if !ExistsInDB("vote", val.Fingerprint, val.LastUpdate) {
// 			r.Votes = append(r.Votes, val)
// 		}
// 	}
// 	for _, val := range postResp.Keys {
// 		if !ExistsInDB("key", val.Fingerprint, val.LastUpdate) {
// 			r.Keys = append(r.Keys, val)
// 		}
// 	}
// 	for _, val := range postResp.Truststates {
// 		if !ExistsInDB("truststate", val.Fingerprint, val.LastUpdate) {
// 			r.Truststates = append(r.Truststates, val)
// 		}
// 	}
// 	postResp.Boards = r.Boards
// 	postResp.Threads = r.Threads
// 	postResp.Posts = r.Posts
// 	postResp.Votes = r.Votes
// 	postResp.Keys = r.Keys
// 	postResp.Truststates = r.Truststates
// 	elapsed := time.Since(t)
// 	// logging.Logf(1, "dbg_CheckIfExistsInDb time spent: %#v.", elapsed.String())
// 	return postResp
// }

// getManifestOfCache gets the manifest of a cache. Location is the url up to cache name.
func getManifestOfCache(
	host string, subhost string, port uint16, location string, reverseConn *net.Conn) (Response, error) {
	logging.Log(2, fmt.Sprintf("Making a request to %s\n", fmt.Sprint(location, "/manifest/0.json")))
	firstManifestPage, err := GetPageRaw(
		host, subhost, port, fmt.Sprint(location, "/manifest/0.json"), "GET", []byte{}, reverseConn)
	if err != nil {
		return Response{}, err
	}
	var resp Response
	resp = InsertApiResponseToResponse(resp, firstManifestPage)
	if firstManifestPage.Pagination.Pages > 0 {
		for i := uint64(1); i <= firstManifestPage.Pagination.Pages; i++ {
			logging.Log(2, fmt.Sprintf("Making a request to %s\n", fmt.Sprint(location, "/manifest/", i, ".json")))
			page, err := GetPageRaw(host, subhost, port,
				fmt.Sprint(location, "/manifest/", i, ".json"), "GET", []byte{}, reverseConn)
			if err != nil {
				return Response{}, err
			}
			var pgResp Response
			pgResp = InsertApiResponseToResponse(pgResp, page)
			resp = concatResponses(resp, pgResp)
		}
	}
	return resp, nil
}

// getIndexOfCache gets the index of a cache. Location is the url up to cache name.
func getIndexOfCache(
	host string, subhost string, port uint16, location string, reverseConn *net.Conn) (Response, error) {
	firstIndexPage, err := GetPageRaw(
		host, subhost, port, fmt.Sprint(location, "/index/0.json"), "GET", []byte{}, reverseConn)
	if err != nil {
		return Response{}, err
	}
	var resp Response
	resp = InsertApiResponseToResponse(resp, firstIndexPage)
	if firstIndexPage.Pagination.Pages > 0 {
		for i := uint64(1); i <= firstIndexPage.Pagination.Pages; i++ {
			page, err := GetPageRaw(host, subhost, port,
				fmt.Sprint(location, "/index/", i, ".json"), "GET", []byte{}, reverseConn)
			if err != nil {
				return Response{}, err
			}
			var pgResp Response
			pgResp = InsertApiResponseToResponse(pgResp, page)
			resp = concatResponses(resp, pgResp)
		}
	}
	return resp, nil
}

// getIndexOfEndpoint gets the cache index of an endpoint.
func getIndexOfEndpoint(
	host string, subhost string, port uint16, endpoint string, reverseConn *net.Conn) (Response, error) {
	EndpointIndexResponse, err := GetPageRaw(
		host, subhost, port, fmt.Sprint(endpoint, "/index.json"), "GET", []byte{}, reverseConn)
	var resp Response
	resp = InsertApiResponseToResponse(resp, EndpointIndexResponse)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// checkForEntityInAnswer is a private function which returns whether an entity exists in a cache result. If so, it returns the entity. If not, it returns nil.
func checkForEntityInAnswer(a Answer, fp Fingerprint, t string) interface{} {
	switch t {
	case "boards":
		var entities []Board
		entities = append(entities, a.Boards...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	case "threads":
		var entities []Thread
		entities = append(entities, a.Threads...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	case "posts":
		var entities []Post
		entities = append(entities, a.Posts...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	case "votes":
		var entities []Vote
		entities = append(entities, a.Votes...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	case "addresses":
		// Nothing happens, as addresses aren't queryable
	case "keys":
		var entities []Key
		entities = append(entities, a.Keys...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	case "truststates":
		var entities []Truststate
		entities = append(entities, a.Truststates...)
		for _, entity := range entities {
			if entity.Fingerprint == fp {
				return entity
			}
		}
	}
	return nil
}

// inTimeRange returns true or false based on whether the value given are within the bounds of the given timestamps.
func inTimeRange(oldest Timestamp, newest Timestamp, val Timestamp) bool {
	if val > oldest && val < newest {
		return true
	} else {
		return false
	}
}

// pullFullEntityFromCache returns the item you have requested by fingerprint from the cache you are pointing at. If no result is found, it will return an empty interface. This could be implemented as a normal GetCache and then search, but that requires the entire cache to be downloaded, whereas this method stops and returns as soon as it can.
func pullFullEntityFromCache(cacheUrl string, cachePage int, fingerprint Fingerprint, t string, host string, subhost string, port uint16, reverseConn *net.Conn) (interface{}, error) {
	if cachePage == 0 {
		// If the cache page is zero, the item we need is either in the first page, or we don't know the cache page, so we need to search.

		// Get the first raw page (because we need to access pagination),
		pageResp, err := GetPageRaw(host, subhost, port, fmt.Sprint(cacheUrl, "/0.json"), "GET", []byte{}, reverseConn)
		if err != nil {
			return nil, errors.New(
				fmt.Sprint(
					"The item is not found at the location the index points to.",
					", Error: ", err,
					", CacheUrl: ", cacheUrl,
					", CachePage: ", cachePage,
					", Fingerprint: ", fingerprint))
		}
		// And look at the page count, so we know how many times to iterate.
		pageCount := pageResp.Pagination.Pages
		// Check the Answer type object to see whether we have it or not.
		entity := checkForEntityInAnswer(pageResp.ResponseBody, fingerprint, t)
		if entity == nil {
			// We haven't found what we wanted on the first page, so we go forward on searching other pages.
			for i := uint64(1); i <= pageCount; i++ { // Pagination starts from 0
				pageResp2, err := GetPageRaw(host, subhost, port,
					fmt.Sprint(cacheUrl, "/", i, ".json"), "GET", []byte{}, reverseConn)
				if err != nil {
					return nil, errors.New(
						fmt.Sprint(
							"The item is not found at the location the index points to.",
							", Error: ", err,
							", CacheUrl: ", cacheUrl,
							", CachePage: ", cachePage,
							", Fingerprint: ", fingerprint))
				}
				// Again, check for whether entity exists on this page.
				entity := checkForEntityInAnswer(pageResp2.ResponseBody, fingerprint, t)
				if entity != nil {
					// If we have an entity that fits the bill, return it and exit.
					return entity, nil
				}
			}
		} else {
			// If we have found what we want on the first page, return it and exit.
			return entity, nil
		}
	} else {
		// If we know the cache page, we can just directly fetch the item.
		pageResp, err := GetPageRaw(host, subhost, port, fmt.Sprint(cacheUrl, "/", cachePage, ".json"), "GET", []byte{}, reverseConn)
		if err != nil {
			return nil, errors.New(
				fmt.Sprint(
					"The item is not found at the location the index points to.",
					", Error: ", err,
					", CacheUrl: ", cacheUrl,
					", CachePage: ", cachePage,
					", Fingerprint: ", fingerprint))
		}
		entity := checkForEntityInAnswer(pageResp.ResponseBody, fingerprint, t)
		return entity, nil
	}

	return nil, errors.New( // If nothing is found, return empty item error.
		fmt.Sprint(
			"The item is not found at the location the index points to.",
			", CacheUrl: ", cacheUrl,
			", CachePage: ", cachePage,
			", Fingerprint: ", fingerprint))
}

// Query struct. Used to provide input to the Query function below.

type QueryData struct {
	EntityType  string
	Fingerprint Fingerprint
	Creation    Timestamp // Can be empty.
	LastUpdate  Timestamp // Last *known* update, can be empty
}

// Query requests an entity from the remote provided. It takes an index form struct. It only returns the requested item or an empty answer.
func Query(host string, subhost string, port uint16, q QueryData, reverseConn *net.Conn) (Response, error) {
	var r Response
	// Before doing anything else, if the type is thread or post, disable LastUpdate. Those items are not updateable.
	epAddress := mapEndpointToEndpointAddress(q.EntityType)
	result, err := getIndexOfEndpoint(host, subhost, port, epAddress, reverseConn)
	endpointIndex := result.CacheLinks
	if err != nil {
		return r, nil
	}
	// Do a range search within all caches that include the last update and creation timestamps. This is where we figure out which caches we need to search.
	var cachesSlice []ResultCache
	for _, cache := range endpointIndex {
		if inTimeRange(cache.StartsFrom, cache.EndsAt, q.Creation) || inTimeRange(cache.StartsFrom, cache.EndsAt, q.LastUpdate) {
			// This adds the endpoints which declare themselves to be in the time range of either Creation, LastUpdate or both. Mind that the creation will be only on one cache, but there may be more than one update.
			cachesSlice = append(cachesSlice, cache)
			// FUTURE: If there is a LastUpdate, checking Creation is inefficient as the result found in the cache pointed at by Creation will not be used. But for purposes of simplicity and to avoid checking for the corner conditions created by having LastUpdate stopping Creation check, I'm leaving it there to be made more efficient in a future refactoring.
		}
	}
	if q.Creation == 0 && q.LastUpdate == 0 {
		// If no data is provided as to when the entity could be, we have to go through all of the data to find it.
		cachesSlice = endpointIndex
	}
	// cachesSlice has all of the caches we have to search now.

CacheIterator: // Naming the for loop CacheIterator.
	for _, cache := range cachesSlice {
		cacheLocation := fmt.Sprint(epAddress, "/", cache.ResponseUrl)
		cIndex, err := getIndexOfCache(host, subhost, port, cacheLocation, reverseConn)
		if err != nil {
			logging.Log(1, fmt.Sprintf("Error in CacheIterator coming from GetIndexOfCache. Error: %s", err))
		}
		// Save the EntityIndexes into proper locations on Response.
		switch q.EntityType {
		case "boards":
			entities := cIndex.BoardIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Boards = append(r.Boards, obj.(Board))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		case "threads":
			entities := cIndex.ThreadIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Threads = append(r.Threads, obj.(Thread))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		case "posts":
			entities := cIndex.PostIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Posts = append(r.Posts, obj.(Post))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		case "votes":
			entities := cIndex.VoteIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Votes = append(r.Votes, obj.(Vote))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		case "addresses":
			// Nothing happens, as addresses aren't queryable
			return r, nil
		case "keys":
			entities := cIndex.KeyIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Keys = append(r.Keys, obj.(Key))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		case "truststates":
			entities := cIndex.TruststateIndexes
			// For each of those entities,
			for _, entityIndex := range entities {
				// Check if this is what we want.
				if entityIndex.Fingerprint == q.Fingerprint {
					// If so, pull the result from cache.
					obj, err := pullFullEntityFromCache(cacheLocation, entityIndex.PageNumber, q.Fingerprint, q.EntityType, host, subhost, port, reverseConn)
					if err != nil {
						return r, errors.New(
							fmt.Sprint(
								"Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available.",
								", Error: ", err,
								", Host: ", host,
								", Subhost: ", subhost,
								", Port: ", port,
								", QueryData: ", q))
					}
					// And put into the proper part of the response.
					r.Truststates = append(r.Truststates, obj.(Truststate))
					// And finally, break the for loop, so it won't look at other caches when it's done.
					break CacheIterator
				}
			}
		}
	}
	return r, nil
}

/*=======================================================
=            Reverse open end status senders            =
=======================================================*/
/*
	These functions communicate the end status at our end (when we're obliging with a reverse open request) so that the remote can know whether our sync succeeded or not. This is a data that we provide to the remote as a courtesy.

	We do not care about responses to these, whether they failed, succeeded, or something entirely else.

	This is where the reverse connection is terminated, also, since these are the last ever calls on a reverse connection. We tried to have it elsewhere, but it causes some issues where it can actually close before we actually send this message through.
*/

func SendReverseOpenStatusRefused(reverseConn *net.Conn) {
	methodUrlEndpoint := "refused"
	sendReverseOpenStatusMessage(methodUrlEndpoint, reverseConn)
}

func SendReverseOpenStatusFailed(reverseConn *net.Conn) {
	methodUrlEndpoint := "failed"
	sendReverseOpenStatusMessage(methodUrlEndpoint, reverseConn)
}

func SendReverseOpenStatusSuccessful(reverseConn *net.Conn) {
	methodUrlEndpoint := "successful"
	sendReverseOpenStatusMessage(methodUrlEndpoint, reverseConn)
}

func sendReverseOpenStatusMessage(methodUrlEndpoint string, reverseConn *net.Conn) {
	logging.Logf(1, "SendReverseOpenStatus: Sending message... Message: '%v', ", methodUrlEndpoint)
	loc := "revconn/" + methodUrlEndpoint
	_, err := Fetch("", "", 0, loc, "GET", []byte(""), reverseConn)
	if err != nil && !strings.Contains(err.Error(), "Non-200 status code returned from Fetch") {
		/*
			^ Non-200 is fine – not only we don't care about what remote does with them, it might also be just old versions where these endpoints don't exist.
		*/
		logging.Logf(1, "SendReverseOpenStatus: Message delivery failed. Message: '%v', Error: %v", methodUrlEndpoint, err)
		return
	}
	logging.Logf(1, "SendReverseOpenStatus: Message delivery succeeded. Message: %v", methodUrlEndpoint)
	(*reverseConn).SetDeadline(time.Now())
	(*reverseConn).Close()
}

/*=====  End of Reverse open end status senders  ======*/
