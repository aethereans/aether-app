// Backend > Swarmtest
// This package / file contains the testing framework that creates an e2e swarm testing environment for the backend nodes by spawning multiple nodes and letting them interact with each other.

/*
  This package does a few things in order.
  - Requests the database generator to generate a random database with given parameters. This requires node to be set up and dependencies installed on the local machine.
  - Run a static server with the generated database.
  - Spin up the backend node with the appropriate command, so that the backend node will ingest the data from the static node into its own database. Then let it die.
  - Start the metrics receiver
  - Start the backend node again, this time with full-on status. Inject all non-first nodes with the address of the starter node, and see what happens.

*/

package main

import (
	// pb "aether-core/aether/backend/metrics/proto"
	sms "aether-core/aether/backend/swarmtest/simplemetricsserver"
	"aether-core/aether/services/ports"
	"aether-core/aether/services/toolbox"
	"encoding/json"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Utility functions

func saveFileToDisk(fileContents []byte, path string, filename string) {
	ioutil.WriteFile(fmt.Sprint(path, "/", filename), fileContents, 0755)
}

// We could do proper flag parsing but it feels unnecessary here, tbh. If this test grows in size we could probably do that.
type settingsStruct struct {
	swarmsize       int // number of nodes to be created and tested against.
	testdurationsec int // how many seconds the main test will run for. This does not include the time it takes to prime the swarm nodes from the donor nodes.
	staticnodeloc   string
	swarmplanloc    string
	dbsize          string
}

type node struct {
	appname           string // appname determines which folder it will be saved within the OS. Changing this is our way of making sure multiple instances of the app can work separately from each other.
	generatedDataPath string // Data path that the static generator will put the resulting static node into.
	staticServerPort  int
	externalPort      int
}

var settings settingsStruct

func setDefaults() {
	toolbox.CreatePath("Runtime-Generated-Files")
	settings.swarmsize = 5
	settings.testdurationsec = 600 //5/450 is the test
	// xs (20k obj per node), s (0.5m obj/n), m
	settings.dbsize = "xs"
	// settings.staticnodeloc = "Runtime-Generated-Files/temp_generated_data"
	settings.staticnodeloc = "/Users/Helios/Documents/temp_generated_data"
	spl, err := filepath.Abs("Runtime-Generated-Files/swarm-plan.json")
	if err != nil {
		log.Fatal(fmt.Sprintf("Converting the swarm plan to an absolute file path has failed. Error: %s", err))
	}
	settings.swarmplanloc = spl
}

// generateSwarmNames generates the appIdentifier names the swarm nodes will use. This is the main thing that allows swarm nodes to occupy different folders in terms of their db, settings.
func generateSwarmNames() []node {
	var nodes []node
	swarmPorts := ports.GetFreePorts(settings.swarmsize)
	for i := 0; i < settings.swarmsize; i++ {
		n := node{}
		n.appname = fmt.Sprintf("Aether-%d", i)
		n.staticServerPort = 17000 + i
		n.externalPort = swarmPorts[i]
		nodes = append(nodes, n)
	}
	return nodes
}

func getDbSize() string {
	if settings.dbsize == "xs" {
		return "--xsmall"
	} else if settings.dbsize == "s" {
		return "--small"
	} else if settings.dbsize == "m" {
		return "--medium"
	} else {
		log.Fatal(fmt.Sprintf("Unknown database size requested. Size requested: %s", settings.dbsize))
		return ""
	}
}

// generateNodeData generates the random data that every swarm node will be seeded with. Every node will have its own set of data.
func generateNodeData(n *node) {
	log.Printf("Generating the required random database for %s.", n.appname)
	nodeTempFolder := settings.staticnodeloc
	nodeTempPath := fmt.Sprintf("%s/%s", nodeTempFolder, n.appname)
	// Set up a folder with the appropriate name for each node that is being run
	toolbox.CreatePath(nodeTempFolder)
	abspath, err := filepath.Abs(nodeTempPath)
	if err != nil {
		log.Fatal(err)
	}
	n.generatedDataPath = abspath
	// Check if there is anything that exists in the folder. If so, skip this.
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		cmd := exec.Command("node", "main.js", getDbSize(), "--nosign", abspath)
		cmd.Dir = "../../../Documentation/database-generator/"
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// Uncomment above to see the generator output
		cmd.Run()
		log.Printf("Random database generation for %s is complete.", n.appname)
	} else {
		log.Printf("Skipping data generation because the static node generation folder for %s already exists. If you want it regenerated, please delete the appropriate folder.", n.appname)
	}
}

// startServingStaticNodeAsDataDonor starts the generated data as a static node of its own, so that our swarm node that we want to get up to speed will be able to sync with it and get its data.
func startServingStaticNodeAsDataDonor(n node) *http.Server {
	srv := &http.Server{Addr: fmt.Sprint(":", n.staticServerPort)}
	fs := http.FileServer(http.Dir(n.generatedDataPath))
	// http.HandleFunc("/", fs)
	// http.Handle("/", fs)
	svmux := http.NewServeMux()
	svmux.Handle("/", fs)
	srv.Handler = svmux

	go func(svmux *http.ServeMux) {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("HTTP Server at the port %d is closing. %s", n.staticServerPort, err)
		}
	}(svmux)

	// returning reference so caller can call Shutdown()
	return srv
}

// insertDataIntoBackendNodeInstance starts our live node in a special mode called --syncandquit, and it points to the static node we spun up as the data donor (bootstrapper). This makes our swarm node pull in all the data we want from our data donor.
func insertDataIntoBackendNodeInstance(n node) {
	// Here, we call our orchestrate command to pull in the data that we need into the db of the given node.
	/*
		The command we need is:
		wipeterm && go run main.go orchestrate --appname="A2Test" --bootstrapip 127.0.0.1 --bootstrapport 9000 --bootstraptype 255 --logginglevel 1 --syncandquit
	*/
	log.Printf("Database insert of randomly generated node data is starting for the node %s", n.appname)
	cmd := exec.Command(
		"go", "run", "main.go", "orchestrate",
		fmt.Sprintf("--appname=%s", n.appname),
		"--bootstrapip=127.0.0.1",
		fmt.Sprintf("--bootstrapport=%d", n.staticServerPort),
		"--bootstraptype=255",
		"--logginglevel=1",
		"--printtostdout",
		fmt.Sprintf("--port=%d", n.externalPort),
		"--metricsdebugmode",
		fmt.Sprintf("--swarmplan=%s", settings.swarmplanloc),
		"--pagesigcheckenabled=false",
		"--fpcheckenabled=false",
		"--powcheckenabled=false",
		"--sigcheckenabled=false",
		"--tlsenabled=false",
		"--syncandquit",
	)
	cmd.Stdout = os.Stdout
	cmd.Dir = "../../../../aether-core/aether/backend/"
	cmd.Run()
	log.Printf("Database insert of randomly generated node data is complete for the node %s", n.appname)
}

type PlanCommand struct {
	CommandName     string
	FromNodeAppName string
	ToIp            string
	ToPort          int
	ToType          int
	TriggerAfter    time.Duration
	Force           bool
	// Force: false: insert into the addresses DB and let node discover connectivity itself
	// Force: true: directly trigger sync()
}

// findPort finds the port of a given node in the available nodes.
func findPort(name string, nodes *[]node) (int, error) {
	for _, n := range *nodes {
		if n.appname == name {
			return n.externalPort, nil
		}
	}
	return 0, fmt.Errorf("The bootstrap node name (AppIdentifier) you've given does not exist in the Swarm. You've given: %s", name)
}

/*
generateSwarmSchedules generates the connection requests that we need to deal with. This is effectively the "track" we are running, in that it generates a pattern that will request nodes to connect to each other, and see how the information spreads from that point on. At the end of completion, this will save it to disk, and provide a save location so that swarm nodes can read from it.

## Test types:

Simple:
You have one bootstrap node, and all other nodes connect to that node at the same time. What this does is that it gets all the data from the bootstrap node and puts them in each. Now, in the first load, the other nodes might not be able to find each other, because they're connecting at the same time and the bs node might not be able to add them to its database fast enough so that their link communicates to others. However, The second sync with that node a minute later will be able to distribute to each of them each other's addresses, and they should then start syncing.
*/
func generateSwarmSchedules(nodes []node, testType string) {
	var planCommands []PlanCommand
	bsNodeName := "Aether-0" // This is given externally
	bsPort, err := findPort(bsNodeName, &nodes)
	if err != nil {
		log.Fatal(err)
	}
	if testType == "simple" {
		for _, n := range nodes {
			if n.appname != bsNodeName {
				// If not bootstrap node, make a call to bootstrap node in 30 secs
				cr := generateConnectionRequest(n.appname, bsPort, 10, false)
				planCommands = append(planCommands, cr)
			}
		}
	} else if testType == "manual_postresponse_reuse_3_nodes" {
		// Instructions: disable the regular sync. XS size
		// A1 > A0 @ 3sec (no reuse)
		// A0 > A1 @ 30 sec (no reuse)
		// A2 > A0 @ 60sec (votes and posts should have 2 cacheresults)
		// Expected result: A0 should be able to reuse response entirely, and return two caches, one generated on A2's request, one remainder from A1's request, for both posts and votes on XS size.
		if len(nodes) != 3 {
			log.Fatal("This test requires 3 nodes.")
		}
		for _, n := range nodes {
			if n.appname == "Aether-0" {
				cr := generateConnectionRequest(n.appname, nodes[1].externalPort, 30, true)
				planCommands = append(planCommands, cr)
			} else if n.appname == "Aether-1" {
				cr := generateConnectionRequest(n.appname, nodes[0].externalPort, 3, true)
				planCommands = append(planCommands, cr)
			} else if n.appname == "Aether-2" {
				cr := generateConnectionRequest(n.appname, nodes[0].externalPort, 60, true)
				planCommands = append(planCommands, cr)
			}
		}
	} else if testType == "manual_postresponse_plus_directresp_4_nodes" {
		/*
			Instructions: Set Vote page threshold to 15000 and disable regular sync. XS size.
			A0 > A1 - a0 gets 12800 more. it's now eligible for post response generation (t3)
			A1 > A0 - a0 generates first post response (t20)
			A0 > A2 - a0 gets 12800 more (but that's below threshold) (t45)
			A2 > A0 - should result in one post container response and one direct response. (t70)
		*/
		if len(nodes) != 4 {
			log.Fatal("This test requires 4 nodes.")
		}
		for _, n := range nodes {
			if n.appname == "Aether-0" {
				cr := generateConnectionRequest(n.appname, nodes[1].externalPort, 3, true)
				cr2 := generateConnectionRequest(n.appname, nodes[2].externalPort, 45, true)
				planCommands = append(planCommands, cr)
				planCommands = append(planCommands, cr2)
			} else if n.appname == "Aether-1" {
				cr := generateConnectionRequest(n.appname, nodes[0].externalPort, 20, true)
				planCommands = append(planCommands, cr)
			} else if n.appname == "Aether-2" {
				cr := generateConnectionRequest(n.appname, nodes[0].externalPort, 70, true)
				planCommands = append(planCommands, cr)
			}
		}
	} else if testType == "cachegen_test" {
		if len(nodes) != 1 {
			log.Fatal("This test requires 1 node.")
		}
		planCmd := generateCacheGenRequest(nodes[0].appname, 2)
		planCommands = append(planCommands, planCmd)
	} else if testType == "reverseopen" {
		// This is where other tests go in.
		for _, n := range nodes {
			if n.appname != bsNodeName {
				// If not bootstrap node, make a call to bootstrap node in 30 secs
				cr := generateReverseOpenRequest(n.appname, bsPort, 10, false)
				planCommands = append(planCommands, cr)
			}
		}
	} else if testType == "" {
		// This is where other tests go in.
	}

	crAsByte, err := json.MarshalIndent(planCommands, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err2 := ioutil.WriteFile(fmt.Sprint(settings.swarmplanloc), crAsByte, 0755)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func generateConnectionRequest(originName string, toPort int, triggerAfter int, force bool) PlanCommand {
	r := PlanCommand{}
	r.FromNodeAppName = originName
	r.ToIp = "127.0.0.1"
	r.ToPort = toPort
	r.ToType = 2
	r.TriggerAfter = time.Duration(triggerAfter) * time.Second
	r.Force = force
	r.CommandName = "connect"
	return r
}

func generateCacheGenRequest(appname string, triggerAfter int) PlanCommand {
	return PlanCommand{FromNodeAppName: appname, CommandName: "cachegen", TriggerAfter: time.Duration(triggerAfter) * time.Second}
}

func generateReverseOpenRequest(originName string, toPort int, triggerAfter int, force bool) PlanCommand {
	r := PlanCommand{}
	r.FromNodeAppName = originName
	r.ToIp = "127.0.0.1"
	r.ToPort = toPort
	r.ToType = 2
	r.TriggerAfter = time.Duration(triggerAfter) * time.Second
	r.Force = force
	r.CommandName = "reverseopen"
	return r
}

// collectAndSaveResults saves collates and saves the metrics at the end of th e test. This is where we get the insights we want out.
func collectAndSaveResults(startTs int64) {
	// At the end of the test, save the data into a JSON file so it can be replayed.
	// structuredBuf := sms.structureBufData(sms.Buf)
	strBuf := sms.ProcessConnectionStates(sms.ConnStateBuf, sms.DbStateBuf, startTs)
	bufAsByte, err := json.MarshalIndent(strBuf, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	formattedTime := fmt.Sprint(time.Now().Format(time.RFC1123))
	toolbox.CreatePath("Runtime-Generated-Files/Test Results")
	err3 := ioutil.WriteFile(fmt.Sprintf("Runtime-Generated-Files/Test Results/Swarm Results %s.json", formattedTime), bufAsByte, 0755)
	if err3 != nil {
		log.Fatal(err3)
	}
}

func startSwarmNode(appname string, externalPort int, killTimeout int, wg *sync.WaitGroup, swarmNodeId int) {
	log.Printf("We're starting the swarm node with the app name %s at the port %d", appname, externalPort)
	defer wg.Done()
	cmd := exec.Command(
		"go", "run", "main.go", "orchestrate",
		fmt.Sprintf("--appname=%s", appname),
		"--logginglevel=1",
		"--printtostdout",
		fmt.Sprintf("--port=%d", externalPort),
		"--metricsdebugmode",
		"--pagesigcheckenabled=false",
		"--fpcheckenabled=false",
		"--powcheckenabled=false",
		"--sigcheckenabled=false",
		"--allowlocalhostremotes=true",
		fmt.Sprintf("--killtimeout=%d", killTimeout),
		fmt.Sprintf("--swarmplan=%s", settings.swarmplanloc),
		fmt.Sprintf("--swarmnodeid=%d", swarmNodeId))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "../../../../aether-core/aether/backend/"
	err2 := cmd.Run()
	if err2 != nil {
		log.Fatal(fmt.Sprintf("The swarm node %s has crashed with an error. Error: %v", appname, err2))
	}
}

func startSwarmNodes(nodes []node, killTimeout int) {
	// sms.presentnodes does not depend on saveresults - it's immediately available as the metrics come in. That's where we can get the port from.
	// The wait group blocks until all goroutines are complete and exited.
	var wg sync.WaitGroup
	for key, n := range nodes {
		wg.Add(1)
		go startSwarmNode(n.appname, n.externalPort, killTimeout, &wg, key)
	}
	wg.Wait()
	fmt.Println("All swarm nodes have exited per their kill timeouts.")
}

var startTime int64

func main() {
	start := time.Now()
	startTime = start.Unix()
	setDefaults()
	durAstDur := time.Duration(settings.testdurationsec) * time.Second
	fmt.Printf("Started at: %s with Nodes: %v, DbSize (each): %v, Duration: %s. End~: %s \n", start.Format(time.RFC1123), settings.swarmsize, settings.dbsize, durAstDur, time.Now().Add(durAstDur*12/10).Format(time.RFC1123))
	go sms.StartListening()
	// For each node that we have requested
	nodes := generateSwarmNames()
	// spew.Dump(ports.GetFreePorts(100))
	for i, _ := range nodes {
		generateNodeData(&nodes[i])
		serverInstance := startServingStaticNodeAsDataDonor(nodes[i])
		insertDataIntoBackendNodeInstance(nodes[i])
		// After the donor gives to the swarm node the whole load, the donor is killed so as to make sure that swarm nodes are the only online nodes in the petri dish.
		serverInstance.Shutdown(nil)
	}
	// Here, generate the list of connection requests we want to inject to the swarm nodes. This is where we create the connection mapping we want to test live.
	generateSwarmSchedules(nodes, "simple")
	// generateSwarmSchedules(nodes, "reverseopen")
	// Start the nodes, this time without a kill-switch at the end of the load.
	startSwarmNodes(nodes, settings.testdurationsec)
	// spew.Dump(nodes)
	collectAndSaveResults(startTime)
	fmt.Printf("It took %d to run this swarm test. It was set to run for %d seconds. Time is now %s\n", int(time.Since(start).Seconds()), settings.testdurationsec, time.Now().Format(time.RFC1123))
}
