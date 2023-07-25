// Backend > Metrics

// This package is the metrics service defined and used by the metrics package. Metrics is used in non-release (alpha, beta, etc.) versions to observe network behaviour. It does not collect any information regarding the user, only about the backend that the user is using and how it is behaving in the network.

package metrics

import (
	pb "aether-core/aether/backend/metrics/proto"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "log"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"

	// "github.com/davecgh/go-spew/spew"
	"strconv"
	"time"
)

const (
	addr            = "127.0.0.1:19999"
	metricsDisabled = true
)

// StartConnection establishes a connection with the metrics server. This is not a singleton and it has to be closed when the metrics delivery is complete.
func StartConnection() (pb.MetricsServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logging.Log(1, fmt.Sprintf("Did not connect: %v", err))
	}
	// defer conn.Close()
	c := pb.NewMetricsServiceClient(conn)
	return c, conn
}

func requestMetricsToken(client pb.MetricsServiceClient, machine *pb.Machine) *pb.Machine_MetricsToken {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// machine := getAnonymousMachineIdentifier()
	metricsToken, err := client.RequestMetricsToken(ctx, machine)
	if err != nil {
		logging.Log(1, fmt.Sprintf("Could not get Metrics Token: %v", err))
	}
	// anonymousMachineIdentifier.MetricsToken = metricsToken
	return metricsToken
}

var anonymousMachineIdentifier *pb.Machine

// generateAnonymousMachineIdentifier generates the metrics identifier for this computer, once every start. Only when this is empty a new one is generated. This is only empty at the start. Configtype is either backend or frontend. This automatically asks for metrics token if there is none.
func getAnonymousMachineIdentifier(client pb.MetricsServiceClient, configType string) *pb.Machine {
	if len(anonymousMachineIdentifier.GetNodeId()) != 0 {
		return anonymousMachineIdentifier
	}
	var proto pb.Machine
	nodeid := globals.BackendConfig.GetNodeId()
	// fmt.Println(nodeid)
	proto.GetNodeId()
	proto.NodeId = nodeid
	subprots := globals.BackendConfig.GetServingSubprotocols() // returns subprotocolshim
	var protobufFormattedSubprots []*pb.Protocol
	for _, val := range subprots {
		var subprot pb.Protocol
		subprot.Name = val.Name
		subprot.VersionMajor = int32(val.VersionMajor)
		subprot.VersionMinor = int32(val.VersionMinor)
		subprot.SupportedEntities = val.SupportedEntities
		protobufFormattedSubprots = append(protobufFormattedSubprots, &subprot)
	}
	proto.Protocols = protobufFormattedSubprots
	proto.Client = &pb.Client{
		Name:         globals.BackendConfig.GetClientName(),
		VersionMajor: int32(globals.BackendConfig.GetClientVersionMajor()),
		VersionMinor: int32(globals.BackendConfig.GetClientVersionMinor()),
		VersionPatch: int32(globals.BackendConfig.GetClientVersionPatch())}
	proto.Address = globals.BackendConfig.GetExternalIp()
	proto.Port = int32(globals.BackendConfig.GetExternalPort())
	if configType == "backend" {
		if len(globals.BackendConfig.GetMetricsToken()) > 0 {
			proto.MetricsToken.Token = globals.BackendConfig.GetMetricsToken()
		} else {
			token := requestMetricsToken(client, &proto)
			proto.MetricsToken = token
			if token != nil {
				globals.BackendConfig.SetMetricsToken(token.Token)
			}
		}
	} else if configType == "frontend" {
		if len(globals.FrontendConfig.GetMetricsToken()) > 0 {
			proto.MetricsToken.Token = globals.FrontendConfig.GetMetricsToken()
		} else {
			token := requestMetricsToken(client, &proto)
			proto.MetricsToken = token
			if token != nil {
				globals.FrontendConfig.SetMetricsToken(token.Token)
			}
		}
	} else {
		logging.LogCrash(fmt.Sprintf("You either didn't provide a configType in getAnonymousMachineIdentifier, or provided one that doesn't match any existing configtype. You provided: %s", configType))
	}
	anonymousMachineIdentifier = &proto
	return anonymousMachineIdentifier
}

func deliverBackendMetrics(client pb.MetricsServiceClient, metrics *pb.Metrics) *pb.MetricsDeliveryResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	machine := getAnonymousMachineIdentifier(client, "backend")
	// var metrics pb.Metrics
	metrics.Machine = machine
	r, err := client.UploadMetrics(ctx, metrics)
	if err != nil {
		logging.Log(1, fmt.Sprintf("Could not deliver Metrics: %v", err))
	}
	return r
}

// insertEntitiesAsMetricForm inserts the entities the batchInsert takes into the database into the metrics form to be sent over the wire. This is only used in the debug mode, regular nodes do not send information about what is being sent and received even when regular metrics are enabled.
func insertEntitiesAsMetricForm(input []interface{}, proto *pb.Metrics) {
	if proto.Persistence == nil {
		proto.Persistence = &pb.Persistence{}
	}
	for _, val := range input {
		e := pb.Entity{}
		switch item := val.(type) {
		case api.Board:
			e.EntityType = pb.Entity_BOARD
			e.Fingerprint = string(item.Fingerprint)
			e.LastUpdate = int64(item.LastUpdate)
		case api.Thread:
			e.EntityType = pb.Entity_THREAD
			e.Fingerprint = string(item.Fingerprint)
		case api.Post:
			e.EntityType = pb.Entity_POST
			e.Fingerprint = string(item.Fingerprint)
		case api.Vote:
			e.EntityType = pb.Entity_VOTE
			e.Fingerprint = string(item.Fingerprint)
			e.LastUpdate = int64(item.LastUpdate)
		case api.Key:
			e.EntityType = pb.Entity_KEY
			e.Fingerprint = string(item.Fingerprint)
			e.LastUpdate = int64(item.LastUpdate)
		case api.Truststate:
			e.EntityType = pb.Entity_TRUSTSTATE
			e.Fingerprint = string(item.Fingerprint)
			e.LastUpdate = int64(item.LastUpdate)
		case api.Address:
			e.EntityType = pb.Entity_ADDRESS
			e.AddressLocation = string(item.Location)
			e.AddressSublocation = string(item.Sublocation)
			e.AddressPort = int32(item.Port)
		}
		proto.Persistence.ArrivedEntitiesSinceLastMetricsDbg = append(proto.Persistence.ArrivedEntitiesSinceLastMetricsDbg, &e)
	}
	// spew.Dump(proto)
}

// insertNodeEntityAsMetricForm converts the dbnode entity that comes from the node insert from the db writer. This is nice,because it allows us to keep track of what node each node would connect to.
func insertNodeEntityAsMetricForm(input map[string]string, proto *pb.Metrics) {
	if proto.Persistence == nil {
		proto.Persistence = &pb.Persistence{}
	}
	n := pb.NodeEntity{}
	n.Fingerprint = input["Fingerprint"]
	blc, err := strconv.Atoi(input["BoardsLastCheckin"])
	tlc, err2 := strconv.Atoi(input["ThreadsLastCheckin"])
	plc, err3 := strconv.Atoi(input["PostsLastCheckin"])
	vlc, err4 := strconv.Atoi(input["VotesLastCheckin"])
	klc, err5 := strconv.Atoi(input["KeysLastCheckin"])
	trlc, err6 := strconv.Atoi(input["TruststatesLastCheckin"])
	alc, err7 := strconv.Atoi(input["AddressesLastCheckin"])
	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		logging.LogCrash(fmt.Sprintf("Conversion from Node DB type to NodeEntity Metric type has failed. This means the database for the node in question has been tampered with. Quitting Errors for field conversions: %s, %s, %s, %s, %s, %s, %s NodeId: ", err, err2, err3, err4, err5, err6, err7, input["Fingerprint"]))
	}
	n.BoardsLastCheckin = int64(blc)
	n.ThreadsLastCheckin = int64(tlc)
	n.PostsLastCheckin = int64(plc)
	n.VotesLastCheckin = int64(vlc)
	n.KeysLastCheckin = int64(klc)
	n.TruststatesLastCheckin = int64(trlc)
	n.AddressesLastCheckin = int64(alc)
	var nArr []*pb.NodeEntity

	nArr = append(nArr, &n)
	proto.Persistence.NodeInsertionsSinceLastMetricsDbg = nArr
	// spew.Dump(proto)
}

/* This is the metric collector function. Before doing anything else, it checks whether the metrics collection is allowed by the user. If not, it's a no-op. This does not send the metrics, it just readies the metrics page to be sent.
 */

func CollateMetrics(metricName string, payload interface{}) {
	if globals.BackendConfig.GetMetricsLevel() > 0 || globals.BackendTransientConfig.MetricsDebugMode {
		// We are metrics-enabled. Create a metrics message.
		m := globals.BackendTransientConfig.CurrentMetricsPage
		// The metric name has to be the same as how this is represented in the protobuf as a matter of convention.
		if metricName == "ArrivedEntitiesSinceLastMetricsDbg" {
			// This is coming from inbound DB entries.
			// We need both the type check and metricname check because multiple places can return the same type of item.
			input := payload.([]interface{})
			insertEntitiesAsMetricForm(input, &m)
		} else if metricName == "NodeInsertionsSinceLastMetricsDbg" {
			//When a new node is inserted or an extant node is otherwise updated, this is what triggers if debug is enabled.
			input := payload.(map[string]string)
			insertNodeEntityAsMetricForm(input, &m)
		}
		globals.BackendTransientConfig.CurrentMetricsPage = m
	}
}

// This is the place where we send metrics. This place should also involve labelling the metrics with the right tag. (pb.Metrics / request metrics token)
func SendMetrics(client pb.MetricsServiceClient) *pb.MetricsDeliveryResponse {
	result := deliverBackendMetrics(client, &globals.BackendTransientConfig.CurrentMetricsPage)
	globals.BackendTransientConfig.CurrentMetricsPage = pb.Metrics{} // After the send, blank out the metrics page.
	return result
}

func SendConnState(remote api.Address, isOpen bool, firstSync bool, ims *[]persistence.InsertMetrics) {
	if metricsDisabled {
		return
	}
	client, conn := StartConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	machine := getAnonymousMachineIdentifier(client, "backend")
	// var metrics pb.Metrics
	var connState pb.ConnState
	connState.Machine = machine
	oConn := pb.OrchestrateConn{}
	oConn.State = isOpen
	oConn.ToIp = string(remote.Location)
	oConn.ToPort = int32(remote.Port)
	oConn.ToNodeName = remote.Client.ClientName
	oConn.Timestamp = int64(time.Now().Unix())
	oConn.FirstSync = firstSync
	// Handle InsertMetrics to generate a receipt.
	receipt := generateReceipts(ims)
	oConn.Objects = &receipt
	oConn.CurrentDbSizeMb = int32(globals.GetDbSize())
	connState.Connection = &oConn
	_, err := client.SendConnectionState(ctx, &connState)
	if err != nil {
		logging.Log(1, fmt.Sprintf("Could not deliver ConnState: %v", err))
	}
}

func generateReceipts(ims *[]persistence.InsertMetrics) pb.Objects {
	receipt := pb.Objects{}
	im := persistence.InsertMetrics{}
	if ims != nil {
		for _, val := range *ims {
			im.Add(val)
		}
	}
	receipt.Boards = int64(im.BoardsReceived)
	receipt.Threads = int64(im.ThreadsReceived)
	receipt.Posts = int64(im.PostsReceived)
	receipt.Votes = int64(im.VotesReceived)
	receipt.Keys = int64(im.KeysReceived)
	receipt.Truststates = int64(im.TruststatesReceived)
	receipt.Addresses = int64(im.AddressesReceived)
	// logging.Logf(1, "%#v", receipt)
	return receipt
}

// SendDbState sends the counts of all DB entities to the metrics server for debugging purposes. This should be only used for testing, not on prod instances.
func SendDbState() {
	if metricsDisabled {
		return
	}
	return // do not send any for now.
	client, conn := StartConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbc := persistence.Dbg_ReadDatabaseCounts()
	var dbState pb.DbState
	dbState.SwarmNodeName = globals.BackendConfig.GetClientName()
	t := time.Now()
	dbState.Timestamp = t.Unix()
	dbState.Time = time.Unix(dbState.Timestamp, 0).String() // remove nsec etc.
	dbState.BoardsCount = dbc.Boards
	dbState.ThreadsCount = dbc.Threads
	dbState.PostsCount = dbc.Posts
	dbState.VotesCount = dbc.Votes
	dbState.KeysCount = dbc.Keys
	dbState.TruststatesCount = dbc.Truststates
	dbState.AddressesCount = dbc.Addresses
	client.SendDbState(ctx, &dbState)
}
