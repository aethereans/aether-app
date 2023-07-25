// This is a naive implementation of the metrics server for swarm testing purposes.

package simplemetricsserver

import (
	pb "aether-core/aether/backend/metrics/proto"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Buf map[int64][]pb.Metrics // There can be multiple metrics pages arriving in the same UNIX timestamp, hence the []slice.

var ConnStateBuf []pb.ConnState
var DbStateBuf []pb.DbState

type server struct{}

func (s *server) RequestMetricsToken(ctx context.Context, remoteMachine *pb.Machine) (*pb.Machine_MetricsToken, error) {
	// fmt.Printf("A message was received from the node w/ port: %d. It's requesting a metrics token.\n", remoteMachine.GetPort())
	// saveNode(remoteMachine.Client.GetName(), int(remoteMachine.GetPort()))
	metricsToken := pb.Machine_MetricsToken{Token: "testtoken"}
	return &metricsToken, nil
}

func (s *server) UploadMetrics(ctx context.Context, metrics *pb.Metrics) (*pb.MetricsDeliveryResponse, error) {
	// fmt.Printf("A message was received from the node w/ port: %d. It's sending metrics.\n", metrics.GetMachine().GetPort())
	// This saves inbound metrics into a file, so that we will have a record of what swarm nodes are doing in the network.
	now := time.Now().Unix()
	Buf[now] = append(Buf[now], *metrics)
	// spew.Dump(metrics)
	return &pb.MetricsDeliveryResponse{}, nil
}

func (s *server) SendConnectionState(ctx context.Context, connState *pb.ConnState) (*pb.MetricsDeliveryResponse, error) {
	// fmt.Println("We received a connection state.")
	// fmt.Println(connState)
	ConnStateBuf = append(ConnStateBuf, *connState)
	return &pb.MetricsDeliveryResponse{}, nil
}

func (s *server) SendDbState(ctx context.Context, dbState *pb.DbState) (*pb.MetricsDeliveryResponse, error) {
	// fmt.Println("We received a connection state.")
	// fmt.Println(connState)
	DbStateBuf = append(DbStateBuf, *dbState)
	return &pb.MetricsDeliveryResponse{}, nil
}

func StartListening() {
	Buf = make(map[int64][]pb.Metrics)
	fmt.Println("Metrics server started listening.")
	lis, err := net.Listen("tcp", "127.0.0.1:19999")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	fmt.Println("Metrics server stopped listening.")

}

// METRICS MANIPULATION

// StructuredBuffer is the struct that collects the node input and db input events that are being reported by the metrics framework from the swarm nodes.
type StructuredBuffer struct {
	StartTimestamp int64 // This is our zero point.
	Nodes          []BufNode
}

type BufNode struct {
	NodeId          string
	DbInputEvents   []DbInputEvent
	NodeInputEvents []NodeInputEvent
}

type DbInputEvent struct {
	Timestamp     int64
	InputtedItems []*pb.Entity
}

type NodeInputEvent struct {
	Timestamp    int64
	InputtedNode *pb.NodeEntity
}

// AddNodeIdIfNotExtant makes a list of nodes that has reported data so far.
func (sbuf *StructuredBuffer) AddNodeIdIfNotExtant(nodeid string) {
	var exists bool
	bufNode := BufNode{}
	for _, val := range sbuf.Nodes {
		if val.NodeId == nodeid {
			exists = true
		}
	}
	if !exists {
		bufNode = BufNode{NodeId: nodeid}
		sbuf.Nodes = append(sbuf.Nodes, bufNode)
	}
}

// findNodeKeyInSBuf finds a specific node in the metrics data.
func (sbuf *StructuredBuffer) findNodeKeyInSBuf(nodeid string) int {
	for key, val := range sbuf.Nodes {
		if val.NodeId == nodeid {
			return key
		}
	}
	log.Fatal("This should never happen. Func: findNodeKeyInSBuf")
	return -1
}

var startTime int64

// structureBufData converts raw metrics data into something easier to read. This data is still humongous, however, and needs further processing.
func structureBufData(rawBuf map[int64][]pb.Metrics) *StructuredBuffer {
	sBuf := StructuredBuffer{}
	sBuf.StartTimestamp = startTime
	// For each timestamp we have:
	for ts, val := range rawBuf {
		// Check every metrics page inside. For every metrics page:
		for _, metricsPage := range val {
			// Below is here because we need to check machine node ids to access the appropriate bufnode.
			sBuf.AddNodeIdIfNotExtant(metricsPage.Machine.NodeId)
			// Look whether the nodeid exists in this sbuf and add as needed.
			// Determine if this page is a node insert or db insert.
			nkey := sBuf.findNodeKeyInSBuf(metricsPage.Machine.NodeId)
			if metricsPage.Persistence.NodeInsertionsSinceLastMetricsDbg != nil {
				// Node insert.
				nodeInputEvent := NodeInputEvent{
					Timestamp:    ts,
					InputtedNode: metricsPage.Persistence.NodeInsertionsSinceLastMetricsDbg[0]}
				sBuf.Nodes[nkey].NodeInputEvents = append(sBuf.Nodes[nkey].NodeInputEvents, nodeInputEvent)
			} else if metricsPage.Persistence.ArrivedEntitiesSinceLastMetricsDbg != nil {
				// Db Insert
				dbInputEvent := DbInputEvent{Timestamp: ts, InputtedItems: metricsPage.Persistence.ArrivedEntitiesSinceLastMetricsDbg}
				sBuf.Nodes[nkey].DbInputEvents = append(sBuf.Nodes[nkey].DbInputEvents, dbInputEvent)
			}
		}
	}
	return &sBuf
}

// CONN STATE MANIPULATION

type CSFinalData struct {
	TotalTestTimeSeconds            int64
	NetworkReachedEquilibrium       bool
	NetworkReachedDBSizeEquilibrium bool
	EndDbSizeMb                     int64
	TimeToEquilibriumSeconds        int64
	NumberOfNodes                   int
	EquilibriumTime                 string
	Nodes                           []CSNode
	DbStates                        []pb.DbState
}

func (f *CSFinalData) AddNode(n CSNode) {
	idx := f.indexOfNode(n)
	if idx == -1 {
		f.Nodes = append(f.Nodes, n)
	} else {
		f.Nodes[idx].Connections = append(f.Nodes[idx].Connections, n.Connections...)
	}
}

func (f *CSFinalData) indexOfNode(n CSNode) int {
	for key, val := range f.Nodes {
		if val.Name == n.Name {
			return key
		}
	}
	return -1
}

type CSNode struct {
	Name                 string `json:",omitempty"`
	ReachedEquilibrium   bool
	EndDbSizeMb          int64
	EquilibriumTime      string   `json:",omitempty"`
	EquilibriumTimestamp int64    `json:",omitempty"`
	FirstSyncCount       int      `json:",omitempty"`
	Connections          []CSConn `json:",omitempty"`
}

type CSConn struct {
	RemoteName string `json:",omitempty"`
	Duration   string `json:",omitempty"`
	FirstSync  bool   `json:",omitempty"`
	State      bool   `json:",omitempty"`
	Timestamp  int64  `json:",omitempty"`
	StartTime  string `json:",omitempty"`
	// Objects    CSObjects `json:",omitempty"`
	Receipt     string `json:",omitempty"`
	DbSizeAtEnd int    `json:",omitempty"`
	Diff        int
}

// type CSObjects struct {
// 	Boards      int64 `json:",omitempty"`
// 	Threads     int64 `json:",omitempty"`
// 	Posts       int64 `json:",omitempty"`
// 	Votes       int64 `json:",omitempty"`
// 	Keys        int64 `json:",omitempty"`
// 	Truststates int64 `json:",omitempty"`
// 	Addresses   int64 `json:",omitempty"`
// }

func generateReceipt(obj pb.Objects) string {
	return fmt.Sprintf(
		"B: %d, T: %d, P: %d, V: %d, K: %d, TS: %d, A: %d",
		obj.Boards, obj.Threads, obj.Posts, obj.Votes, obj.Keys, obj.Truststates, obj.Addresses)
}

func connConvert(rawConn pb.OrchestrateConn) CSConn {
	c := CSConn{}
	c.RemoteName = rawConn.ToNodeName
	c.State = rawConn.State
	c.Timestamp = rawConn.Timestamp
	c.FirstSync = rawConn.FirstSync
	c.Receipt = generateReceipt(*rawConn.Objects)
	c.DbSizeAtEnd = int(rawConn.CurrentDbSizeMb)
	// obj := CSObjects{}
	// obj.Boards = rawConn.Objects.Boards
	// obj.Threads = rawConn.Objects.Threads
	// obj.Posts = rawConn.Objects.Posts
	// obj.Votes = rawConn.Objects.Votes
	// obj.Keys = rawConn.Objects.Keys
	// obj.Truststates = rawConn.Objects.Truststates
	// obj.Addresses = rawConn.Objects.Addresses
	// c.Objects = obj
	return c
}

func ProcessConnectionStates(rawData []pb.ConnState, rawDbStateData []pb.DbState, startTs int64) CSFinalData {
	endTs := int64(time.Now().Unix())
	finalData := CSFinalData{}
	finalData.NetworkReachedEquilibrium = true
	finalData.NetworkReachedDBSizeEquilibrium = true
	var dbSize int64
	var lastEqTimestamp int64
	for _, val := range rawData {
		n := CSNode{Name: val.Machine.Client.Name}
		n.Connections = append(n.Connections, connConvert(*val.Connection))
		finalData.AddNode(n)
	}
	for key := range finalData.Nodes {
		// Get Db Size.
		fi, _ := os.Stat(filepath.Join("/Users/Helios/Library/Application Support/Air Labs", finalData.Nodes[key].Name, "backend/AetherDB.db"))
		// get the size
		size := fi.Size() / 1000000
		finalData.Nodes[key].EndDbSizeMb = size
		finalconns := processCSConns(finalData.Nodes[key].Connections)
		finalData.Nodes[key].Connections = finalconns
		eq, tStr, fsc := calcCSNodeEquilibriumState(len(finalData.Nodes), finalconns)
		finalData.Nodes[key].ReachedEquilibrium = eq
		finalData.Nodes[key].FirstSyncCount = fsc
		if !eq {
			finalData.NetworkReachedEquilibrium = false
		} else { // has reache equilibrium
			// finalData.Nodes[key].EquilibriumTimestamp = tStr
			finalData.Nodes[key].EquilibriumTime = time.Unix(tStr, 0).String()
			if tStr > lastEqTimestamp {
				lastEqTimestamp = tStr
			}
		}
		// Db size equilibrium calculation.
		if dbSize == 0 {
			dbSize = finalData.Nodes[key].EndDbSizeMb
		}
		if !(dbSize+1 >= finalData.Nodes[key].EndDbSizeMb &&
			dbSize-1 <= finalData.Nodes[key].EndDbSizeMb) {
			// If not in the +1 to -1 range, fail.
			finalData.NetworkReachedDBSizeEquilibrium = false
		}
	}
	if finalData.NetworkReachedEquilibrium {
		finalData.EquilibriumTime = time.Unix(lastEqTimestamp, 0).String()
	}
	if finalData.NetworkReachedDBSizeEquilibrium {
		finalData.EndDbSizeMb = dbSize
	}
	finalData.NumberOfNodes = len(finalData.Nodes)
	finalData.TotalTestTimeSeconds = endTs - startTs
	finalData.TimeToEquilibriumSeconds = lastEqTimestamp - startTs
	// finalData.DbStates = processDbStateData(rawDbStateData)
	finalData.DbStates = rawDbStateData
	return finalData
}

func calcCSNodeEquilibriumState(numberOfNodes int, conns []CSConn) (bool, int64, int) {
	hasFullData := false
	firstSyncsCount := 0
	var lastFirstSyncTimestamp int64
	for key, val := range conns {
		if val.FirstSync {
			firstSyncsCount++
			if val.Timestamp > lastFirstSyncTimestamp {
				lastFirstSyncTimestamp = val.Timestamp
			}
		}
		conns[key].Timestamp = 0
	}
	if firstSyncsCount == numberOfNodes || firstSyncsCount == numberOfNodes+1 {
		// or +1 comes from there being a bootstrap node available possibly.
		// this is not -1 as you would expect removing the node itself from conns list. +1 comes from the donor node we use before the beginning.
		hasFullData = true
	}
	return hasFullData, lastFirstSyncTimestamp, firstSyncsCount
}

// CSConns above has open and close saved as separate conns. We want to make it so that we merge them and fill in the timestamps.
func processCSConns(conns []CSConn) []CSConn {
	var finalConns []CSConn

	var opens []CSConn

	var closes []CSConn

	for _, val := range conns {
		if val.State == true {
			opens = append(opens, val)
		} else {
			closes = append(closes, val)
		}
	}
	for _, val := range opens {
		cls := findClose(closes, val)
		finalConn := merge(val, cls)
		finalConns = append(finalConns, finalConn)
	}
	return finalConns
}

func findClose(conns []CSConn, conn CSConn) CSConn {
	smallestCloseAfterConnStart := int64(time.Now().Unix())
	resultKey := -1
	for key, val := range conns {
		if val.Timestamp >= conn.Timestamp &&
			val.Timestamp < smallestCloseAfterConnStart {
			smallestCloseAfterConnStart = val.Timestamp
			resultKey = key
		}
	}
	if resultKey != -1 {
		return conns[resultKey]
	}
	return CSConn{}
}

func merge(o, cl CSConn) CSConn {
	o.Duration = time.Duration(time.Duration(cl.Timestamp-o.Timestamp) * time.Second).String()
	o.State = false
	o.StartTime = time.Unix(o.Timestamp, 0).String()
	o.Receipt = cl.Receipt
	o.Diff = cl.DbSizeAtEnd - o.DbSizeAtEnd
	o.DbSizeAtEnd = cl.DbSizeAtEnd
	return o
}
