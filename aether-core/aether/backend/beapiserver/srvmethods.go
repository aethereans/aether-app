// Backend > BackendAPI
// This service is the interface of the frontend(s) to the backend. Through this, the BE acts as a traditional database and abstracts away the distribution mechanism.

package beapiserver

import (
	"aether-core/aether/backend/dispatch"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	pb "aether-core/aether/protos/beapi"
	"aether-core/aether/services/create"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
	"fmt"
	"golang.org/x/net/context"
)

type server struct{}

func (s *server) RequestBackendAccess(
	ctx context.Context, req *pb.AccessRequest) (
	*pb.AccessResponse, error) {
	fmt.Println("We received an access request.")
	resp := pb.AccessResponse{}
	resp.Status = &pb.Status{StatusCode: 200}
	return &resp, nil
}

func (s *server) GetBoards(
	ctx context.Context, req *pb.BoardsRequest) (*pb.BoardsResponse, error) {
	resp := pb.BoardsResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once.
	start := req.GetFilters().GetLastRefTimeRange().GetStart()
	end := req.GetFilters().GetLastRefTimeRange().GetEnd()
	fps := req.GetFilters().GetFingerprints()
	if fps == nil {
		fps = &pb.Fingerprints{}
	}
	apiFps := []api.Fingerprint{}
	for key, _ := range fps.Fingerprints {
		apiFps = append(apiFps, api.Fingerprint(fps.Fingerprints[key]))
	}
	apiStart := api.Timestamp(start)
	apiEnd := api.Timestamp(end)
	opts := persistence.OptionalReadInputs{
		AllProvables_Owner:  req.GetFilters().GetGraphFilters().GetOwner(),
		AllProvables_Limit:  int(req.GetFilters().GetGraphFilters().GetLimit()),
		AllProvables_Offset: int(req.GetFilters().GetGraphFilters().GetOffset()),
		Board_Name:          req.GetFilters().GetGraphFilters().GetName(),
	}
	result, _ := persistence.Read("boards", apiFps, []string{}, apiStart, apiEnd, true, &opts)
	for key, _ := range result.Boards {
		r := result.Boards[key].Protobuf()
		resp.Boards = append(resp.Boards, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func (s *server) GetThreads(
	ctx context.Context, req *pb.ThreadsRequest) (*pb.ThreadsResponse, error) {
	resp := pb.ThreadsResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once. If we end up with a primitive type at the end of the chain, we don't have to nil check.
	start := req.GetFilters().GetLastRefTimeRange().GetStart()
	end := req.GetFilters().GetLastRefTimeRange().GetEnd()
	fps := req.GetFilters().GetFingerprints().GetFingerprints()
	apiFps := []api.Fingerprint{}
	for key, _ := range fps {
		apiFps = append(apiFps, api.Fingerprint(fps[key]))
	}
	apiStart := api.Timestamp(start)
	apiEnd := api.Timestamp(end)
	boardfp := req.GetFilters().GetGraphFilters().GetBoard()
	opts := persistence.OptionalReadInputs{
		Thread_Board:        boardfp,
		AllProvables_Owner:  req.GetFilters().GetGraphFilters().GetOwner(),
		AllProvables_Limit:  int(req.GetFilters().GetGraphFilters().GetLimit()),
		AllProvables_Offset: int(req.GetFilters().GetGraphFilters().GetOffset()),
	}
	result, _ := persistence.Read("threads", apiFps, []string{}, apiStart, apiEnd, true, &opts)
	for key, _ := range result.Threads {
		r := result.Threads[key].Protobuf()
		resp.Threads = append(resp.Threads, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func (s *server) GetPosts(
	ctx context.Context, req *pb.PostsRequest) (*pb.PostsResponse, error) {
	resp := pb.PostsResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once.
	start := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetStart())
	end := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetEnd())
	fps := req.GetFilters().GetFingerprints().GetFingerprints()
	apiFps := []api.Fingerprint{}
	for key, _ := range fps {
		apiFps = append(apiFps, api.Fingerprint(fps[key]))
	}
	result, _ := persistence.Read("posts", apiFps, []string{}, start, end, true, &persistence.OptionalReadInputs{
		Post_Board:          req.GetFilters().GetGraphFilters().GetBoard(),
		Post_Thread:         req.GetFilters().GetGraphFilters().GetThread(),
		Post_Parent:         req.GetFilters().GetGraphFilters().GetParent(),
		AllProvables_Owner:  req.GetFilters().GetGraphFilters().GetOwner(),
		AllProvables_Limit:  int(req.GetFilters().GetGraphFilters().GetLimit()),
		AllProvables_Offset: int(req.GetFilters().GetGraphFilters().GetOffset()),
	})
	for key, _ := range result.Posts {
		r := result.Posts[key].Protobuf()
		resp.Posts = append(resp.Posts, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func (s *server) GetVotes(
	ctx context.Context, req *pb.VotesRequest) (*pb.VotesResponse, error) {
	resp := pb.VotesResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once.
	start := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetStart())
	end := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetEnd())
	fps := req.GetFilters().GetFingerprints().GetFingerprints()
	apiFps := []api.Fingerprint{}
	for key, _ := range fps {
		apiFps = append(apiFps, api.Fingerprint(fps[key]))
	}
	result, _ := persistence.Read("votes", apiFps, []string{}, start, end, true,
		&persistence.OptionalReadInputs{
			Vote_Board:          req.GetFilters().GetGraphFilters().GetBoard(),
			Vote_Thread:         req.GetFilters().GetGraphFilters().GetThread(),
			Vote_Target:         req.GetFilters().GetGraphFilters().GetTarget(),
			Vote_TypeClass:      int(req.GetFilters().GetTypeFilters().GetTypeClass()),
			Vote_Type:           int(req.GetFilters().GetTypeFilters().GetType()),
			Vote_NoDescendants:  req.GetFilters().GetGraphFilters().GetNoDescendants(),
			AllProvables_Owner:  req.GetFilters().GetGraphFilters().GetOwner(),
			AllProvables_Limit:  int(req.GetFilters().GetGraphFilters().GetLimit()),
			AllProvables_Offset: int(req.GetFilters().GetGraphFilters().GetOffset()),
		})
	for key, _ := range result.Votes {
		r := result.Votes[key].Protobuf()
		resp.Votes = append(resp.Votes, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func (s *server) GetKeys(
	ctx context.Context, req *pb.KeysRequest) (*pb.KeysResponse, error) {
	resp := pb.KeysResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once.
	start := req.GetFilters().GetLastRefTimeRange().GetStart()
	end := req.GetFilters().GetLastRefTimeRange().GetEnd()
	fps := req.GetFilters().GetFingerprints()
	if fps == nil {
		fps = &pb.Fingerprints{}
	}
	apiFps := []api.Fingerprint{}
	for key, _ := range fps.Fingerprints {
		apiFps = append(apiFps, api.Fingerprint(fps.Fingerprints[key]))
	}
	apiStart := api.Timestamp(start)
	apiEnd := api.Timestamp(end)
	opts := persistence.OptionalReadInputs{
		AllProvables_Owner:  req.GetFilters().GetGraphFilters().GetOwner(),
		AllProvables_Limit:  int(req.GetFilters().GetGraphFilters().GetLimit()),
		AllProvables_Offset: int(req.GetFilters().GetGraphFilters().GetOffset()),
		Key_Name:            req.GetFilters().GetGraphFilters().GetName(),
	}
	result, _ := persistence.Read("keys", apiFps, []string{}, apiStart, apiEnd, true, &opts)
	for key, _ := range result.Keys {
		r := result.Keys[key].Protobuf()
		resp.Keys = append(resp.Keys, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func (s *server) GetTruststates(
	ctx context.Context, req *pb.TruststatesRequest) (*pb.TruststatesResponse, error) {
	resp := pb.TruststatesResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	// The reason we use Get.. is that it can push our null check to the end, so we only have to nil check it once.
	start := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetStart())
	end := api.Timestamp(req.GetFilters().GetLastRefTimeRange().GetEnd())
	fps := req.GetFilters().GetFingerprints().GetFingerprints()
	apiFps := []api.Fingerprint{}
	for key, _ := range fps {
		apiFps = append(apiFps, api.Fingerprint(fps[key]))
	}
	result, _ := persistence.Read("truststates", apiFps, []string{}, start, end, true,
		&persistence.OptionalReadInputs{
			Truststate_Target:    req.GetFilters().GetGraphFilters().GetTarget(),
			Truststate_Domain:    req.GetFilters().GetGraphFilters().GetDomain(),
			Truststate_TypeClass: int(req.GetFilters().GetTypeFilters().GetTypeClass()),
			Truststate_Type:      int(req.GetFilters().GetTypeFilters().GetType()),
			AllProvables_Owner:   req.GetFilters().GetGraphFilters().GetOwner(),
			AllProvables_Limit:   int(req.GetFilters().GetGraphFilters().GetLimit()),
			AllProvables_Offset:  int(req.GetFilters().GetGraphFilters().GetOffset()),
		})
	for key, _ := range result.Truststates {
		r := result.Truststates[key].Protobuf()
		resp.Truststates = append(resp.Truststates, &r)
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

// checkRequestAllowed checks whether the given type is allowed. This is our future point of check gate where we take a look at nonces, PKs, timestamps, and whether the backend is too busy to respond. This is unnecessary when the backend is local, but useful if one wants to provide it as a public service. TODO FUTURE. (This is needed to have publicly accesible backends serving multiple frontends)
func requestAllowed(req interface{}) bool {
	switch req.(type) {
	case *pb.BoardsRequest:
		return true
	case *pb.ThreadsRequest:
		return true
	case *pb.PostsRequest:
		return true
	case *pb.VotesRequest:
		return true
	case *pb.KeysRequest:
		return true
	case *pb.TruststatesRequest:
		return true
	case *pb.MintedContentPayload:
		return true
	case *pb.ConnectToRemoteRequest:
		return true
	default:
		return false
	}
}

// GetBoardThreadsCount counts all threads in a board without a time limit. This will give you all stuff that is available in the local memory. The results of this is not cached, so it will directly hit the backend. If you do this in too many parallel threads, the backend will start to send you 'connection refused's as you exceed the maximum number of simultaneous connections. Be careful with that.
func (s *server) GetBoardThreadsCount(
	ctx context.Context, req *pb.BoardThreadsCountRequest) (*pb.BoardThreadsCountResponse, error) {
	resp := pb.BoardThreadsCountResponse{Count: 0}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	ct := persistence.GetBoardThreadsCount(req.Fingerprint)
	resp.Count = int32(ct)
	resp.Status.StatusCode = 200
	return &resp, nil
}

// GetThreadPostsCount counts all posts in a thread without a time limit. This will give you all stuff that is available in the local memory. The results of this is not cached, so it will directly hit the backend. If you do this in too many parallel threads, the backend will start to send you 'connection refused's as you exceed the maximum number of simultaneous connections. Be careful with that.
func (s *server) GetThreadPostsCount(
	ctx context.Context, req *pb.ThreadPostsCountRequest) (*pb.ThreadPostsCountResponse, error) {
	resp := pb.ThreadPostsCountResponse{Status: &pb.Status{}, Count: 0}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	ct := persistence.GetThreadPostsCount(req.Fingerprint)
	resp.Count = int32(ct)
	resp.Status.StatusCode = 200
	return &resp, nil
}

// SendMintedContent receives the content a frontend has minted and it adds it into the database if the content is valid.
func (s *server) SendMintedContent(
	ctx context.Context, req *pb.MintedContentPayload) (*pb.MintedContentResponse, error) {
	resp := pb.MintedContentResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	allItems := []interface{}{}
	boardsProto := req.GetBoards()
	for k, _ := range boardsProto {
		e := api.Board{}
		e.FillFromProtobuf(*boardsProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	threadsProto := req.GetThreads()
	for k, _ := range threadsProto {
		e := api.Thread{}
		e.FillFromProtobuf(*threadsProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	postsProto := req.GetPosts()
	for k, _ := range postsProto {
		e := api.Post{}
		e.FillFromProtobuf(*postsProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	votesProto := req.GetVotes()
	for k, _ := range votesProto {
		e := api.Vote{}
		e.FillFromProtobuf(*votesProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	keysProto := req.GetKeys()
	for k, _ := range keysProto {
		e := api.Key{}
		e.FillFromProtobuf(*keysProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	truststatesProto := req.GetTruststates()
	for k, _ := range truststatesProto {
		e := api.Truststate{}
		e.FillFromProtobuf(*truststatesProto[k])
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
			resp.Status.StatusCode = 400 // HTTP 400 Bad Request
			return &resp, nil
		}
		allItems = append(allItems, interface{}(e))
	}
	addressesProto := req.GetAddresses()
	if len(addressesProto) > 0 {
		logging.LogCrashf("Addresses insert is not yet implemented.")
		// todo: this should only work if the data is coming from the admin frontend.
	}
	// for k, _ := range addressesProto {
	// 	e := api.Address{}
	// 	e.FillFromProtobuf(*addressesProto[k])
	// 	err2 := api.Verify(api.Provable(&e))
	// 	if err2 != nil {
	// 		logging.Logf(1, "Verification of an entity received from the frontend failed. Error: %v", err2)
	// 		resp.Status.StatusCode = 400 // HTTP 400 Bad Request
	//    return &resp, nil
	// 	}
	// 	allItems = append(allItems, interface{}(e))
	// }
	// ^ TODO FUTURE when needed.
	persistence.BatchInsert(allItems)
	resp.Status.StatusCode = 200
	globals.BackendTransientConfig.NewContentCommitted = true
	// This flag will make the backend attempt to push out the new content fast as possible by triggering a reverse open.
	return &resp, nil
}

// SendConnectToRemoteRequest is a request from an admin frontend asking the backend to connect to a specific remote node. This triggers a sync with that remote immediately, pending completion of any existing sync.
func (s *server) SendConnectToRemoteRequest(
	ctx context.Context, req *pb.ConnectToRemoteRequest) (*pb.ConnectToRemoteResponse, error) {
	resp := pb.ConnectToRemoteResponse{Status: &pb.Status{}}
	if !requestAllowed(req) {
		resp.Status.StatusCode = 401 // HTTP 401 Unauthorised
		return &resp, nil
	}
	logging.Logf(1, "Backend received a connect request to a remote node. Addr: %#v", req.GetAddress())
	addr := constructDirectConnectAddress(
		req.GetAddress().GetLocation(),
		req.GetAddress().GetSublocation(),
		int(req.GetAddress().GetPort()))
	err := dispatch.Sync(addr, []string{}, nil)
	if err != nil {
		resp.Status.StatusCode = 503 // Service unavailable
		resp.Status.ErrorMessage = err.Error()
		return &resp, nil
	}
	resp.Status.StatusCode = 200
	return &resp, nil
}

func constructDirectConnectAddress(loc, subloc string, port int) api.Address {
	subprots := []api.Subprotocol{api.Subprotocol{"c0", 1, 0, []string{"board", "thread", "post", "vote", "key", "truststate"}}}
	addr, err := create.CreateAddress(api.Location(loc), api.Location(subloc), 4, uint16(port), 2, 1, 1, 1, 0, subprots, 2, 0, 0, "Aether", "")
	if err != nil {
		logging.Logf(1, "Constructing direct connect address in response to a connect request from the frontend failed. Error: %v", err)
		return api.Address{}
	}
	addr.SetVerified(true)
	return addr
}
