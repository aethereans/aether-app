package api_test

import (
	"aether-core/aether/backend/cmd"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/signaturing"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519"
)

// Tests will be run with the current protocol version.

// USAGE
// go test -nodeloc=/Users/Helios/Desktop/generated\ nodes/node-newest_4/static_mim_node

// Infrastructure, setup and teardown

var testNodeAddress string
var testNodePort uint16
var nodeLocation string
var protv string

func TestMain(m *testing.M) {
	// Create the database and configs.
	cmd.EstablishConfigs(nil)
	persistence.CreateDatabase()
	persistence.CheckDatabaseReady()
	protv = globals.BackendConfig.GetProtURLVersion()
	globals.BackendTransientConfig.FingerprintCheckEnabled = false
	globals.BackendTransientConfig.SignatureCheckEnabled = false
	globals.BackendTransientConfig.ProofOfWorkCheckEnabled = false
	globals.BackendTransientConfig.PageSignatureCheckEnabled = false
	globals.BackendTransientConfig.PermConfigReadOnly = true
	globals.FrontendTransientConfig.PermConfigReadOnly = true
	globals.BackendConfig.SetMinimumPoWStrengths(5)
	globals.BackendConfig.SetLoggingLevel(0)
	testNodeAddress = "127.0.0.1"
	testNodePort = 8089
	setup(testNodeAddress, testNodePort)
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func copyFile(source string, dest string) (err error) {
	sourceF, err := os.Open(source)
	defer sourceF.Close()
	destF, err := os.Create(dest)
	defer destF.Close()
	_, err = io.Copy(destF, sourceF)
	return err
}

func copyDirectory(source string, dest string) (err error) {
	info, _ := os.Stat(source)
	os.MkdirAll(dest, info.Mode())
	dir, _ := os.Open(source)
	objs, err := dir.Readdir(-1)
	for _, obj := range objs {
		sourceFilePointer := source + "/" + obj.Name()
		destinationFilePointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			// If directory, create sub directories
			err = copyDirectory(sourceFilePointer, destinationFilePointer)
		} else {
			// If file, just copy over
			err = copyFile(sourceFilePointer, destinationFilePointer)
		}
	}
	return err
}

func editExistingJson(source string, regex string, newStringFragment string) {
	r, _ := regexp.Compile(regex)
	data, _ := ioutil.ReadFile(source)
	dataStr := string(data)
	dataStr = r.ReplaceAllString(dataStr, newStringFragment)
	ioutil.WriteFile(source, []byte(dataStr), 0755)
}

func getValidEntity(entityType string) (string, api.Timestamp, api.Timestamp) {
	// getValidEntity provides real entities from the given data set so that the tests can be run on real objects.
	// cache, or any other entity type.
	// if cache, the format is:
	// threads/cache_c4012f95c6b1547656846133acb7c54345ba95e669e13141f38f1fa825fcbce4/
	var returnVal string
	var creation api.Timestamp
	var lastUpdate api.Timestamp
	if entityType == "cache" { // Heads up you're in api_test. This is just for testing - Don't get freaked out by this only being for posts.
		postsDir := fmt.Sprint(nodeLocation, "/", protv, "/c0/posts")
		flist := []string{}
		// Get all paths in that directory.
		filepath.Walk(postsDir,
			func(path string, f os.FileInfo, err error) error {
				_, fname := filepath.Split(path)
				flist = append(flist, fmt.Sprint("c0/posts/", fname))
				return nil
			})
		// fmt.Printf("flist: %#v\n", flist)
		// The first item of the directory.
		firstPostsDirCacheDir := flist[1]
		// The first item of the first item of the directory (the page 0).
		if strings.Contains(firstPostsDirCacheDir, ".DS_Store") {
			firstPostsDirCacheDir = flist[2]
		}
		returnVal = firstPostsDirCacheDir
	} else {
		var dir string
		if entityType == "boards" {
			dir = fmt.Sprint(nodeLocation, "/", protv, "/c0/boards")
		} else if entityType == "posts" {
			dir = fmt.Sprint(nodeLocation, "/", protv, "/c0/posts")
		} else if entityType == "truststates" {
			dir = fmt.Sprint(nodeLocation, "/", protv, "/c0/truststates")
		} else if entityType == "threads_index" { // this is so that we can pull a nonexistent entity from index, delete it from the actual source. This is a test case to have index but not data.
			dir = fmt.Sprint(nodeLocation, "/", protv, "/c0/threads")
		} else {
			log.Fatal("The type of entity that is requested to bring a valid instance of could not be determined.")
		}
		fileList := []string{}
		// Get all paths in that directory.
		filepath.Walk(dir,
			func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
		if len(fileList) == 0 {
			log.Fatal("The directory you have provided does not have a Mim node.")
		} else if len(fileList) == 1 && strings.Contains(fileList[1], ".DS_Store") {
			log.Fatal("The directory you have provided does not have a Mim node.")
		}
		// The first item of the directory.
		firstSubdir := fileList[1]
		// The first item of the first item of the directory (the page 0).
		if strings.Contains(firstSubdir, ".DS_Store") {
			firstSubdir = fileList[2]
		}
		var firstItem string
		if entityType == "threads_index" {
			firstItem = fmt.Sprint(firstSubdir, "/index/0.json")
		} else {
			firstItem = fmt.Sprint(firstSubdir, "/0.json")
		}
		fI, _ := ioutil.ReadFile(firstItem)
		var apiResp api.ApiResponse
		json.Unmarshal(fI, &apiResp)
		if entityType == "boards" {
			returnVal = string(apiResp.ResponseBody.Boards[0].Fingerprint)
			creation = apiResp.ResponseBody.Boards[0].Creation
			lastUpdate = apiResp.ResponseBody.Boards[0].LastUpdate
		} else if entityType == "posts" {
			returnVal = string(apiResp.ResponseBody.Posts[0].Fingerprint)
			creation = apiResp.ResponseBody.Posts[0].Creation
			// Posts do not have last update (they're immutable)
		} else if entityType == "truststates" {
			returnVal = string(apiResp.ResponseBody.Truststates[0].Fingerprint)
			creation = apiResp.ResponseBody.Truststates[0].Creation
			lastUpdate = apiResp.ResponseBody.Truststates[0].LastUpdate
		} else if entityType == "threads_index" {
			returnVal = string(apiResp.ResponseBody.ThreadIndexes[0].Fingerprint)
		}
	}
	return returnVal, creation, lastUpdate
}

func setup(testNodeAddress string, testNodePort uint16) {
	// 1) Set up a web server.
	// 2) Put a static node behind that web server.
	// Parse the nodeloc flag
	flag.StringVar(&nodeLocation, "nodeloc", "", "Node location needed for tests.")
	flag.Parse()

	if len(nodeLocation) == 0 {
		// If no node location is given, assume default. This will break when you move that folder off desktop...
		nodeLocation = "/Users/Helios/Desktop/Hazel Desktop/2015-Q4 /generated nodes/node-newest_16/static_mim_node"
	}

	// // Vote endpoint borkage test setup start.

	// This breaks the index.json of the vote endpoint. This is to test how the protocol behaves under broken endpoints.

	votesIndex := fmt.Sprint(nodeLocation, "/", protv, "/c0/votes/index.json")
	vI, _ := ioutil.ReadFile(votesIndex)
	var votesIndexResp api.ApiResponse
	json.Unmarshal(vI, &votesIndexResp)
	// fmt.Printf("votespage: %#v\n", votesIndexResp)
	for i, _ := range votesIndexResp.Results {
		votesIndexResp.Results[i].ResponseUrl = "invalid url invalid url invalid url"
	}
	result, _ := json.Marshal(votesIndexResp)
	ioutil.WriteFile(votesIndex, result, 0755)
	// fmt.Printf("votespage: %#v\n", result)

	// // Thread remove item and leave its index

	// This removes the item that is pointed at by the index.

	threadIndexFp, _, _ := getValidEntity("threads_index")
	// fmt.Printf("threadIndexFp: %#v\n", threadIndexFp)
	dir := fmt.Sprint(nodeLocation, "/", protv, "/c0/threads")
	fileList := []string{}
	// Get all paths in that directory.
	filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return nil
		})
	if len(fileList) == 0 {
		log.Fatal("The directory you have provided does not have a Mim node.")
	} else if len(fileList) == 1 && strings.Contains(fileList[1], ".DS_Store") {
		log.Fatal("The directory you have provided does not have a Mim node.")
	}
	// The first item of the directory.
	firstSubdir := fileList[1]
	// The first item of the first item of the directory (the page 0).
	if strings.Contains(firstSubdir, ".DS_Store") {
		firstSubdir = fileList[2]
	}
	firstItem := fmt.Sprint(firstSubdir, "/0.json")
	// fmt.Printf("firstItem: %#v\n", firstItem)

	editExistingJson(
		firstItem,
		threadIndexFp, `yo yo this broken`)

	// // Cache tests setup start.

	// For the malformed cache cases, we need to generate the malformed data.
	// 1) Cache: Negative value in pages count.
	// 2) Cache: Missing cache
	// 3) Cache: Huge pages count as a DDoS attack on oneself?
	// Build the directory. We're using posts because it usually has a lot of items.
	postsDir := fmt.Sprint(nodeLocation, "/", protv, "/c0/posts")
	fileList2 := []string{}
	// Get all paths in that directory.
	filepath.Walk(postsDir,
		func(path string, f os.FileInfo, err error) error {
			fileList2 = append(fileList2, path)
			return nil
		})
	if len(fileList2) == 0 {
		log.Fatal("The directory you have provided does not have a Mim node.")
	} else if len(fileList2) == 1 && strings.Contains(fileList2[1], ".DS_Store") {
		log.Fatal("The directory you have provided does not have a Mim node.")
	}
	// The first item of the directory.
	firstPostsDir := fileList2[1]
	// The first item of the first item of the directory (the page 0).
	if strings.Contains(firstPostsDir, ".DS_Store") {
		firstPostsDir = fileList2[2]
	}
	// Create the test directories from that baseline.
	copyDirectory(firstPostsDir, fmt.Sprint(postsDir, "/cache_negative_page_count"))
	copyDirectory(firstPostsDir, fmt.Sprint(postsDir, "/cache_missing_pages"))
	copyDirectory(firstPostsDir, fmt.Sprint(postsDir, "/cache_huge_page_number"))
	// First case: Edit the page count number to be negative.
	editExistingJson(
		fmt.Sprint(postsDir, "/cache_negative_page_count/0.json"),
		`"pages": \d`, `"pages": -16`)
	// Second case: missing pages from the cache.
	os.Remove(fmt.Sprint(postsDir, "/cache_missing_pages/7.json"))
	os.Remove(fmt.Sprint(postsDir, "/cache_missing_pages/34.json"))
	os.Remove(fmt.Sprint(postsDir, "/cache_missing_pages/35.json"))
	os.Remove(fmt.Sprint(postsDir, "/cache_missing_pages/40.json"))
	os.Remove(fmt.Sprint(postsDir, "/cache_missing_pages/55.json"))
	// Third case: Huge page count as a way to DDoS.
	editExistingJson(
		fmt.Sprint(postsDir, "/cache_huge_page_number/0.json"),
		`"pages": \d`, `"pages": 549815`)

	// // Cache tests setup end.

	// // Endpoint tests setup start.

	// votesEndpointIndex := fmt.Sprint(nodeLocation, "/",protv,"/c0/votes/index.json")
	// editExistingJson(
	// 	votesEndpointIndex,
	// 	`cache_20fd7dc96e6ec13b3afcb2d417709e7e407fb47a6d2f6e0c1914af0df17b41a1`, `broken_cache_20fd7dc96e6ec13b3afcb2d417709e7e407fb47a6d2f6e0c1914af0df17b41a1`)
	// editExistingJson(
	// 	votesEndpointIndex,
	// 	`cache_62df5cb262755febac0098a69d0a5c16938b932d939ebd1ed88887bda18052ce`, `broken_cache_62df5cb262755febac0098a69d0a5c16938b932d939ebd1ed88887bda18052ce`)
	// editExistingJson(
	// 	votesEndpointIndex,
	// 	`cache_e6506b867c067fef752db3f36337e5cb0caae4502342c181a72c74d6c9083db0`, `broken_cache_e6506b867c067fef752db3f36337e5cb0caae4502342c181a72c74d6c9083db0`)
	// editExistingJson(
	// 	votesEndpointIndex,
	// 	`cache_c71472ac4d44c1a90e7656f2b0783182a4e07dbe7e584e5900083425ac9a3a50`, `broken_cache_c71472ac4d44c1a90e7656f2b0783182a4e07dbe7e584e5900083425ac9a3a50`)

	fakeEndpointDir := fmt.Sprint(nodeLocation, "/", protv, "/c0/invalidendpoint")
	copyDirectory(postsDir, fakeEndpointDir)

	// // Endpoint tests setup end.

	// // Query tests setup start.
	threadsDir := fmt.Sprint(nodeLocation, "/", protv, "/c0/threads")
	flist := []string{}
	// Get all paths in that directory.
	filepath.Walk(threadsDir,
		func(path string, f os.FileInfo, err error) error {
			flist = append(flist, path)
			return nil
		})
	// fmt.Printf("flist: %#v\n", flist)
	// The first item of the directory.
	firstThreadsCacheDir := flist[1]
	// The first item of the first item of the directory (the page 0).
	if strings.Contains(firstThreadsCacheDir, ".DS_Store") {
		firstThreadsCacheDir = flist[2]
	}
	os.Remove(fmt.Sprint(firstThreadsCacheDir, "/3.json"))

	// // Query tests setup end.

	// Create a HTTP server serving the nodeloc.
	fs := http.FileServer(http.Dir(nodeLocation))
	http.Handle("/", fs)
	http.HandleFunc("/", protv, "/timeouter", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30000 * time.Second)
	})
	http.HandleFunc("/", protv, "/c0/invalid_data.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is some invalid JSON."))
	})
	go http.ListenAndServe(fmt.Sprint(testNodeAddress, ":", testNodePort), nil)
}

func teardown() {
}

func ValidateTest(expected interface{}, actual interface{}, t *testing.T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

// Tests

// Fetch tests

func TestFetch_Success(t *testing.T) {
	httpResp, err :=
		api.Fetch(testNodeAddress, "", testNodePort, "status", "GET", []byte{}, nil)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	expected := ""
	actual := string(httpResp)
	ValidateTest(expected, actual, t)
}

func TestFetch_404(t *testing.T) {
	_, err := api.Fetch(testNodeAddress, "", testNodePort, "this is a nonexistent location", "GET", []byte{}, nil)
	expected := "Non-200 status code returned from Fetch. Received status code: 404, Host: 127.0.0.1, Subhost: , Port: 8089, Location: this is a nonexistent location, Method: GET"
	actual := err.Error()
	ValidateTest(expected, actual, t)
}

func TestFetch_Refused(t *testing.T) {
	_, err := api.Fetch(testNodeAddress, "", 48915, "this is a nonexistent location", "GET", []byte{}, nil)
	expected := "The host refused the connection. Host:127.0.0.1, Subhost: , Port: 48915, Location: this is a nonexistent location"
	actual := err.Error()
	ValidateTest(expected, actual, t)
}

func TestFetch_Timeout(t *testing.T) {
	_, err := api.Fetch(testNodeAddress, "", testNodePort, "timeouter", "GET", []byte{}, nil)
	expected := "Timeout exceeded. Host:127.0.0.1, Subhost: , Port: 8089, Location: timeouter"
	actual := err.Error()
	ValidateTest(expected, actual, t)
}

// Get Page tests
func TestGetPageRaw_Success(t *testing.T) {
	resp, err := api.GetPageRaw(testNodeAddress, "", testNodePort, "c0/boards/index.json", "GET", []byte{}, nil)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Results) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestGetPageRaw_Unparsable(t *testing.T) {
	_, err := api.GetPageRaw(testNodeAddress, "", testNodePort, "c0/invalid_data.json", "GET", []byte{}, nil)
	expected := "The JSON that arrived over the network is malformed. JSON: This is some invalid JSON., Host: 127.0.0.1, Subhost: , Port: 8089, Location: c0/invalid_data.json"
	actual := err.Error()
	ValidateTest(expected, actual, t)
}

func TestGetPage_Success(t *testing.T) {
	resp, _, err := api.GetPage(testNodeAddress, "", testNodePort, "c0/boards/index.json", "GET", []byte{}, nil)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.CacheLinks) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

// Get Cache tests

func TestGetCache_Success(t *testing.T) {
	cacheName, _, _ := getValidEntity("cache")
	// fmt.Printf("cachename: %#v\n", cacheName)
	resp, err := api.GetCache(testNodeAddress, "", testNodePort, cacheName)
	// Pointing out the name directly here is brittle. We have no others, so fix this.
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Posts) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestGetCache_InvalidPageCount_CountNegative(t *testing.T) {
	_, err := api.GetCache(testNodeAddress, "", testNodePort, "c0/posts/cache_negative_page_count/")
	errMessage := "The JSON that arrived over the network is malformed"
	if err == nil {
		t.Errorf("JSON parser failed to catch the error. No error from parser.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("The parser returned an error that was different than the expected one. '%s'", err)
	}
}

func TestGetCache_InvalidPageCount_HugePageCount(t *testing.T) {
	// This also tests for the 3 consequent missing pages safeguard, as the huge fake page count is stopped by the 3 pages after the last real page failing.
	_, err := api.GetCache(testNodeAddress, "", testNodePort, "c0/posts/cache_huge_page_number/")
	errMessage := "3 or more broken pages"
	if err == nil {
		t.Errorf("GetCache failed to stop when 3 missing pages followed each other.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("GetCache returned an error that was different than the expected one. '%s'", err)
	}
}

func TestGetCache_MissingPage(t *testing.T) {
	resp, err := api.GetCache(testNodeAddress, "", testNodePort, "c0/posts/cache_missing_pages/")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Posts) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

// Get Endpoint tests

func TestGetGETEndpoint_Success(t *testing.T) {
	resp, err := api.GetGETEndpoint(testNodeAddress, "", testNodePort, "threads", 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Threads) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestGetGETEndpoint_3ConsequentCachesMissingFailure(t *testing.T) {
	_, err := api.GetGETEndpoint(testNodeAddress, "", testNodePort, "votes", 0)
	errMessage := "3 or more cache failures"
	if err == nil {
		t.Errorf("Did not notice the cache being missing.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestGetGETEndpoint_NonexistentEndpoint(t *testing.T) {
	_, err := api.GetGETEndpoint(testNodeAddress, "", testNodePort, "fakeendpoint", 0)
	errMessage := "Get Endpoint failed because it couldn't get the index of the endpoint."
	if err == nil {
		t.Errorf("Did not notice the endpoint being missing.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestGetGETEndpoint_EndpointNameAndContentsMismatch(t *testing.T) {
	// This test is present to make sure that endpoints have no dependence on their names. The parsing logic should be global.
	resp, err := api.GetGETEndpoint(testNodeAddress, "", testNodePort, "invalidendpoint", 0)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Posts) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

// Get Remote Node tests

func TestGetRemoteNode_Success(t *testing.T) {
	resp, err := api.GetRemoteNode(testNodeAddress, "", testNodePort)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) == 0 ||
		len(resp.Threads) == 0 ||
		len(resp.Posts) == 0 ||
		len(resp.Votes) != 0 || // Votes are deliberately borked.
		len(resp.Addresses) == 0 ||
		len(resp.Keys) == 0 ||
		len(resp.Truststates) == 0 {
		fmt.Printf("Boards: %#v\n", len(resp.Boards))
		fmt.Printf("Threads: %#v\n", len(resp.Threads))
		fmt.Printf("Posts: %#v\n", len(resp.Posts))
		fmt.Printf("Votes: %#v\n", len(resp.Votes))
		fmt.Printf("Addresses: %#v\n", len(resp.Addresses))
		fmt.Printf("Keys: %#v\n", len(resp.Keys))
		fmt.Printf("Truststates: %#v\n", len(resp.Truststates))
		t.Errorf("Test failed, a part of the response is empty.")
	}
}

// Query tests

func TestQuery_Fingerprint_Success(t *testing.T) {
	entityFp, _, _ := getValidEntity("boards")
	data := api.QueryData{"boards", api.Fingerprint(entityFp), 0, 0}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Boards) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestQuery_FingerprintAndCreation_Success(t *testing.T) {
	// fmt.Printf("valid entity: %#v\n", api.Fingerprint(getValidEntity("posts")))
	// Mind that it's asking for something created AFTER 0451102626
	entityFp, creation, _ := getValidEntity("posts")
	data := api.QueryData{"posts", api.Fingerprint(entityFp), creation, 0}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Posts) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestQuery_FingerprintAndCreationAndLastUpdate_Success(t *testing.T) {
	entityFp, creation, lastUpdate := getValidEntity("truststates")
	data := api.QueryData{"truststates", api.Fingerprint(entityFp), creation, lastUpdate}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Truststates) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestQuery_FingerprintAndLastUpdate_Success(t *testing.T) {
	entityFp, _, lastUpdate := getValidEntity("truststates")
	data := api.QueryData{"truststates", api.Fingerprint(entityFp), 0, lastUpdate}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Truststates) == 0 {
		t.Errorf("Test failed, the response is empty.")
	}
}

func TestQuery_NotFound(t *testing.T) {
	data := api.QueryData{"truststates", "0af3473c5a3ae6376f0d3824b16d2ef90510973c75f889557d39f6616ea55535", 0, 0}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Truststates) > 0 {
		t.Errorf("Test failed, the response is not empty. It should be.")
	}
}

func TestQuery_InvalidTimeRange(t *testing.T) {
	data := api.QueryData{"truststates", "7bb882b1e9b679948478266c6ccdd153cb71fbcb2e58bf1237f30d43245eed5d", 1449122236523, 1451543248432}
	resp, err := api.Query(testNodeAddress, "", testNodePort, data)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if len(resp.Truststates) > 0 {
		t.Errorf("Test failed, the response is not empty. It should be.")
	}
}

func TestQuery_TheItemDoesNotExistAtLocationGivenByIndex(t *testing.T) {
	entityFp, _, _ := getValidEntity("threads_index")
	data := api.QueryData{"threads", api.Fingerprint(entityFp), 0, 0}
	_, err := api.Query(testNodeAddress, "", testNodePort, data)
	errMessage := "Could not pull entity from cache. The item is indexed as available in the remote node, but the actual body of the item is not available."
	if err == nil {
		t.Errorf("This should have caused an error.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

// CreatePoW tests

func TestBoardCreatePoW_Success_WithoutKey(t *testing.T) {
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.EntityVersion = 1
	err := newboard2.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newboard2.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestThreadCreatePoW_Success_WithoutKey(t *testing.T) {
	var newthread api.Thread
	newthread.Fingerprint = "my random fingerprint"
	newthread.Creation = 4564654
	newthread.Name = "my thread name"
	newthread.Body = "my thread description"
	newthread.EntityVersion = 1
	err := newthread.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newthread.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestPostCreatePoW_Success_WithoutKey(t *testing.T) {
	var newpost api.Post
	newpost.Fingerprint = "my random fingerprint"
	newpost.Creation = 4564654
	newpost.Parent = "my post parent fingerprint"
	newpost.Body = "my board description"
	newpost.EntityVersion = 1
	err := newpost.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newpost.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestVoteCreatePoW_Success_WithoutKey(t *testing.T) {
	var newvote api.Vote
	newvote.Board = "my random board fingerprint"
	newvote.Creation = 4564654
	newvote.Target = "my vote target fingerprint"
	newvote.EntityVersion = 1
	err := newvote.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newvote.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestKeyCreatePoW_Success_WithoutKey(t *testing.T) {
	var newkey api.Key
	newkey.Type = "my key type"
	newkey.Key = "my key"
	newkey.Name = "my key name"
	newkey.EntityVersion = 1
	err := newkey.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newkey.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestTruststateCreatePoW_Success_WithoutKey(t *testing.T) {
	var newtruststate api.Truststate
	newtruststate.Target = "my truststate target"
	newtruststate.Owner = "my truststate owner"
	newtruststate.Expiry = 4134235
	newtruststate.EntityVersion = 1
	err := newtruststate.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newtruststate.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

// CreateUpdatePoW tests

func TestBoardCreateUpdatePoW_Success_WithoutKey(t *testing.T) {
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	newboard.EntityVersion = 1
	err := newboard.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 16)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			result, err3 := newboard.VerifyPoW("")
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestVoteCreateUpdatePoW_Success_WithoutKey(t *testing.T) {
	var newvote api.Vote
	newvote.Board = "my random board fingerprint"
	newvote.Creation = 4564654
	newvote.Target = "my vote target fingerprint"
	newvote.EntityVersion = 1
	err := newvote.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newvote.Type = 1
		err2 := newvote.CreateUpdatePoW(new(ed25519.PrivateKey), 16)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			result, err3 := newvote.VerifyPoW("")
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestKeyCreateUpdatePoW_Success_WithoutKey(t *testing.T) {
	var newkey api.Key
	newkey.Type = "my key type"
	newkey.Key = "my key"
	newkey.Name = "my key name"
	newkey.EntityVersion = 1
	err := newkey.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newkey.Name = "my name changed"
		err2 := newkey.CreateUpdatePoW(new(ed25519.PrivateKey), 16)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			result, err3 := newkey.VerifyPoW("")
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestTruststateCreateUpdatePoW_Success_WithoutKey(t *testing.T) {
	var newtruststate api.Truststate
	newtruststate.Target = "my truststate target"
	newtruststate.Owner = "my truststate owner"
	newtruststate.Expiry = 4134235
	newtruststate.EntityVersion = 1
	err := newtruststate.CreatePoW(new(ed25519.PrivateKey), 16)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newtruststate.Expiry = 2341512
		err2 := newtruststate.CreateUpdatePoW(new(ed25519.PrivateKey), 16)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			result, err3 := newtruststate.VerifyPoW("")
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

// Create & verify fingerprint tests

func TestBoardCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Board
	item.Fingerprint = "my random fingerprint"
	item.Creation = 4564654
	item.Name = "my board name"
	item.Description = "my board description"
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestThreadCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Thread
	item.Fingerprint = "my random fingerprint"
	item.Creation = 4564654
	item.Name = "my thread name"
	item.Body = "my thread description"
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestPostCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Post
	item.Fingerprint = "my random fingerprint"
	item.Creation = 4564654
	item.Thread = "parent thread fingerprint"
	item.Body = "my post body"
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestVoteCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Vote
	item.Fingerprint = "my random fingerprint"
	item.Target = "my target fingerprint"
	item.Thread = "parent thread fingerprint"
	item.Owner = "owner fingerprint"
	item.Type = 1
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestKeyCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Key
	item.Fingerprint = "my random fingerprint"
	item.Type = "my key type"
	item.Name = "my key name"
	item.Info = "my key info"
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestTruststateCreateVerifyFingerprint_Success(t *testing.T) {
	var item api.Truststate
	item.Fingerprint = "my random fingerprint"
	item.Target = "truststate target"
	item.Owner = "truststate owner"
	item.Expiry = 523523
	item.EntityVersion = 1
	item.CreateFingerprint()
	if item.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := item.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

// Create Signature tests. These include verification as part of the test validation.

func TestBoardCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.ProofOfWork = "my fake pow"
	newboard2.EntityVersion = 1
	err2 := newboard2.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newboard2.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestThreadCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newthread api.Thread
	newthread.Fingerprint = "my random fingerprint2"
	newthread.Creation = 4564654
	newthread.Name = "my thread name"
	newthread.Body = "my thread body"
	newthread.ProofOfWork = "my fake pow"
	newthread.EntityVersion = 1
	err2 := newthread.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newthread.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestPostCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newpost api.Post
	newpost.Fingerprint = "my random fingerprint"
	newpost.Creation = 4564654
	newpost.Parent = "my post parent fingerprint"
	newpost.Body = "my board description"
	newpost.EntityVersion = 1
	err2 := newpost.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newpost.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestVoteCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newvote api.Vote
	newvote.Board = "my random board fingerprint"
	newvote.Creation = 4564654
	newvote.Target = "my vote target fingerprint"
	newvote.EntityVersion = 1
	err2 := newvote.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newvote.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestKeyCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newkey api.Key
	newkey.Type = "my key type"
	newkey.Key = "my key"
	newkey.Name = "my key name"
	newkey.EntityVersion = 1
	err2 := newkey.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newkey.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestTruststateCreateSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newtruststate api.Truststate
	newtruststate.Target = "my truststate target"
	newtruststate.Owner = "my truststate owner"
	newtruststate.Expiry = 4134235
	newtruststate.EntityVersion = 1
	err2 := newtruststate.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err3 := newtruststate.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

// Create Update Signature tests. These have verification as part of the test validation process.

func TestBoardCreateUpdateSignature_Success(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	newboard.EntityVersion = 1
	err := newboard.CreateSignature(privKey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdateSignature(privKey)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			result, err3 := newboard.VerifySignature(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestVoteCreateUpdateSignature_Success(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newvote api.Vote
	newvote.Board = "my random board fingerprint"
	newvote.Creation = 4564654
	newvote.Target = "my vote target fingerprint"
	newvote.EntityVersion = 1
	err := newvote.CreateSignature(privKey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newvote.Type = 1
		err2 := newvote.CreateUpdateSignature(privKey)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			result, err3 := newvote.VerifySignature(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestKeyCreateUpdateSignature_Success(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newkey api.Key
	newkey.Fingerprint = "my test fingerprint"
	newkey.Type = "my key type"
	newkey.Key = "my key"
	newkey.Name = "my key name"
	newkey.EntityVersion = 1
	err5 := newkey.CreateSignature(privKey)
	if err5 != nil {
		t.Errorf("Test failed, err: '%s'", err5)
	} else {
		newkey.Name = "my name changed"
		err2 := newkey.CreateUpdateSignature(privKey)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			result, err3 := newkey.VerifySignature(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestTruststateCreateUpdateSignature_Success(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newtruststate api.Truststate
	newtruststate.Target = "my truststate target"
	newtruststate.Owner = "my truststate owner"
	newtruststate.Expiry = 4134235
	newtruststate.EntityVersion = 1
	err := newtruststate.CreateSignature(privKey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newtruststate.Expiry = 2341512
		err2 := newtruststate.CreateUpdateSignature(privKey)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			result, err3 := newtruststate.VerifySignature(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestApiResponseCreateSignature_Success(t *testing.T) {
	globals.BackendTransientConfig.PageSignatureCheckEnabled = true
	apiResp := api.ApiResponse{}
	apiResp.Prefill()
	// apiResp := responsegenerator.GeneratePrefilledApiResponse()
	err := apiResp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err2 := apiResp.VerifySignature()
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestApiResponseCreateSignature_Fail(t *testing.T) {
	globals.BackendTransientConfig.PageSignatureCheckEnabled = true
	apiResp := api.ApiResponse()
	apiResp.Prefill()
	// apiResp := responsegenerator.GeneratePrefilledApiResponse()
	err := apiResp.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
	apiResp.Entity = "changing the page after it was signed"
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, _ := apiResp.VerifySignature()
		if result != false {
			t.Errorf("Test failed, this Signature should not be valid but it was successfully verified.")
		}
	}
}
