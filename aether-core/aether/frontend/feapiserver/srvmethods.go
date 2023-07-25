// Frontend > FrontendAPI
// This package implements the FrontendAPI GRPC server.

package feapiserver

import (
	"aether-core/aether/frontend/beapiconsumer"
	"aether-core/aether/frontend/clapiconsumer"
	"aether-core/aether/frontend/festructs"
	"aether-core/aether/frontend/refresher"

	// "aether-core/aether/frontend/objpool"
	"aether-core/aether/frontend/inflights"
	// "aether-core/aether/io/api"
	"aether-core/aether/io/api"
	"aether-core/aether/protos/beapi"
	"aether-core/aether/protos/clapi"
	pb "aether-core/aether/protos/feapi"
	"aether-core/aether/protos/feobjects"
	mimapi "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

type server struct{}

func (s *server) BackendReady(ctx context.Context, req *pb.BEReadyRequest) (*pb.BEReadyResponse, error) {
	fmt.Printf("Backend sent us a backendready request. Req: %#v Context: %#v\n", req, ctx)
	peer, ok := peer.FromContext(ctx)
	if !ok {
		logging.Logf(1, "The backend peer address on the GRPC connection in BackendReady could not be read. Something went very wrong, and we are not able to communicate with the backend. Peer: %#v", peer)
		return nil, fmt.Errorf("The backend peer address on the GRPC connection in BackendReady could not be read. Something went very wrong, and we are not able to communicate with the backend. Peer: %#v", peer)
	}
	beAddr := peer.Addr.(*net.TCPAddr)
	backendIp := beAddr.IP.String()
	globals.FrontendConfig.SetBackendAPIAddress(backendIp)
	globals.FrontendConfig.SetBackendAPIPort(int(req.GetPort()))
	globals.FrontendTransientConfig.BackendReady = true
	// Send the ack to the client that we are all ready.
	clapiconsumer.SendFrontendReady()
	resp := pb.BEReadyResponse{}
	return &resp, nil
}

func (s *server) GetThreadAndPosts(ctx context.Context, req *pb.ThreadAndPostsRequest) (*pb.ThreadAndPostsResponse, error) {
	// Get board
	fp := req.GetBoardFingerprint()
	// Get the board carrier
	bc := festructs.BoardCarrier{}
	logging.Logf(3, "Single read happens in BackendReady>One")
	err := globals.KvInstance.One("Fingerprint", fp, &bc)
	if err != nil {
		logging.Logf(1, "Getting BoardCarrier for in GetThreadAndPosts encountered an error. Error: %v", err)
	}
	b := festructs.CompiledBoard{}
	for key, _ := range bc.Boards {
		if bc.Boards[key].Fingerprint == fp {
			b = bc.Boards[key]
		}
	}
	resp := pb.ThreadAndPostsResponse{}
	resp.Board = b.Protobuf()
	// Check if board is subbed / notify enabled.
	subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Board.Fingerprint)
	resp.Board.Subscribed = subbed
	resp.Board.Notify = notify
	resp.Board.LastSeen = lastseen
	// Get thread
	tfp := req.GetThreadFingerprint()
	// Get the thread carrier
	tc := festructs.ThreadCarrier{}
	logging.Logf(3, "Single read happens in GetThreadAndPosts>One")
	err2 := globals.KvInstance.One("Fingerprint", tfp, &tc)
	if err2 != nil {
		logging.Logf(1, "Getting ThreadCarrier for in GetThreadAndPosts encountered an error. Error: %v", err2)
	}
	resp.Thread = tc.MakeTree(false, false) // do not show deleted, do not show orphans
	return &resp, nil
}

func (s *server) GetBoardAndThreads(ctx context.Context, req *pb.BoardAndThreadsRequest) (*pb.BoardAndThreadsResponse, error) {
	start := time.Now()
	fp := req.GetBoardFingerprint()
	// Get the board carrier
	bc := festructs.BoardCarrier{}
	logging.Logf(3, "Single read happens in GetBoardAndThreads>One")
	err := globals.KvInstance.One("Fingerprint", fp, &bc)
	if err != nil {
		logging.Logf(1, "Getting BoardCarrier for in GetBoardAndThreads encountered an error. Error: %v", err)
	}
	b := festructs.CompiledBoard{}
	for key, _ := range bc.Boards {
		if bc.Boards[key].Fingerprint == fp {
			b = bc.Boards[key]
		}
	}
	resp := pb.BoardAndThreadsResponse{}
	resp.Board = b.Protobuf()
	// Check if board is subbed / notify enabled.
	subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Board.Fingerprint)
	resp.Board.Subscribed = subbed
	resp.Board.Notify = notify
	resp.Board.LastSeen = lastseen

	threads := festructs.CThreadBatch{}
	for k1, _ := range bc.Threads {
		// Filter out the moddeletes / modapprovals based on the ruleset.
		if bc.Threads[k1].Board == fp {
			if bc.Threads[k1].CompiledContentSignals.ModApproved || bc.Threads[k1].CompiledContentSignals.SelfModApproved {
				threads = append(threads, bc.Threads[k1])
				continue
			}
			if bc.Threads[k1].CompiledContentSignals.ModBlocked || bc.Threads[k1].CompiledContentSignals.SelfModBlocked {
				continue
			}
			threads = append(threads, bc.Threads[k1])
		}
	}

	// If sort by new, we sort it here. Default sort (the one saved to disk) is popular sort.
	if req.GetSortThreadsByNew() {
		threads.SortByCreation()
	}
	// Convert all threads to protos
	var tprotos []*feobjects.CompiledThreadEntity

	for k, _ := range threads {
		tprotos = append(tprotos, threads[k].Protobuf())
	}
	resp.Threads = tprotos
	elapsed := time.Since(start)
	logging.Logf(1, "Sending board and threads took: %v", elapsed)
	return &resp, nil
}

func (s *server) GetAllBoards(ctx context.Context, req *pb.AllBoardsRequest) (*pb.AllBoardsResponse, error) {
	fmt.Println("We received a get all boards request.")
	var boards []festructs.BoardCarrier
	start := time.Now()
	err := globals.KvInstance.All(&boards)
	if err != nil {
		logging.Logcf(1, "Getting all boards from KvInstance encountered an error. Error: %v", err)
	}

	cb := festructs.CBoardBatch{}
	for key, _ := range boards {
		for k2, _ := range boards[key].Boards {
			item := boards[key].Boards[k2]
			cb = append(cb, item)
		}
	}
	cb.SortByThreadsCount()
	var cproto []*feobjects.CompiledBoardEntity

	for k, _ := range cb {
		item := cb[k].Protobuf()
		cproto = append(cproto, item)
		subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(item.Fingerprint)
		whitelisted := globals.FrontendConfig.ContentRelations.SFWList.IsSFWListedBoard(item.Fingerprint)
		item.Subscribed = subbed
		item.Notify = notify
		item.LastSeen = lastseen
		item.SFWListed = whitelisted
	}

	fmt.Printf("Number of items found in get all boards: %v\n", len(boards))
	resp := pb.AllBoardsResponse{AllBoards: cproto}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return &resp, nil
}

func (s *server) SetClientAPIServerPort(ctx context.Context, req *pb.SetClientAPIServerPortRequest) (*pb.SetClientAPIServerPortResponse, error) {
	logging.Logf(1, "Client said hello.")
	logging.Logf(1, "We received a set client api server port request. Old port was: %v and the new one is %v", globals.FrontendConfig.GetClientPort(), req.Port)
	globals.FrontendConfig.SetClientPort(int(req.Port))
	clapiconsumer.DeliverAmbients()
	inflights := inflights.GetInflights()
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.ClientIsReadyForConnections = true
	clapiconsumer.SendAmbientStatus(&as)
	clapiconsumer.PushLocalUserAmbient() // todo let's disable this for a minute.
	// SendAmbients(false)
	clapiconsumer.SendHomeView()
	clapiconsumer.SendPopularView()
	clapiconsumer.SendNewView()
	clapiconsumer.SendNotifications()
	clapiconsumer.SendOnboardCompleteStatus()
	clapiconsumer.SendModModeEnabledStatus()
	clapiconsumer.SendAlwaysShowNSFWListStatus()
	clapiconsumer.SendExternalContentAutoloadDisabledStatus()
	clapiconsumer.SendSFWListDisabledStatus()
	resp := pb.SetClientAPIServerPortResponse{}
	return &resp, nil
}

func (s *server) SetBoardSignal(ctx context.Context, req *pb.BoardSignalRequest) (*pb.BoardSignalResponse, error) {
	// logging.Logf(1, "We've received a set board signal request.")
	cr := globals.FrontendConfig.GetContentRelations()
	committed := cr.SetBoardSignal(req.Fingerprint, req.Subscribed, req.Notify, req.LastSeen, req.LastSeenOnly)
	globals.FrontendConfig.SetContentRelations(cr)
	resp := pb.BoardSignalResponse{Committed: committed}
	clapiconsumer.DeliverAmbients()
	return &resp, nil
}

func (s *server) GetUserAndGraph(ctx context.Context, req *pb.UserAndGraphRequest) (*pb.UserAndGraphResponse, error) {
	fp := req.GetFingerprint()
	resp := pb.UserAndGraphResponse{}
	resp.UserEntityRequested = req.GetUserEntityRequested()
	resp.UserBoardsRequested = req.GetUserBoardsRequested()
	resp.UserThreadsRequested = req.GetUserThreadsRequested()
	resp.UserPostsRequested = req.GetUserPostsRequested()
	if req.GetUserEntityRequested() {
		uh := festructs.UserHeaderCarrier{}
		logging.Logf(3, "Single read happens in GetUserAndGraph>One")
		err := globals.KvInstance.One("Fingerprint", fp, &uh)
		if err != nil {
			logging.Logf(1, "Getting User Header Carrier for GetUserAndGraph encountered an error. Error: %v", err)
		}
		u := festructs.CompiledUser{}
		for key, _ := range uh.Users {
			if uh.Users[key].Fingerprint == fp {
				u = uh.Users[key]
			}
		}
		resp.User = u.Protobuf()
	}
	// logging.Logf(1, "resp: %v", resp)
	return &resp, nil
}

func (s *server) SendContentEvent(ctx context.Context, req *pb.ContentEventPayload) (*pb.ContentEventResponse, error) {
	logging.Logf(1, "We've received a content event. Event: %v", *req)
	inflights := inflights.GetInflights()
	inflights.Insert(*req)
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
	resp := pb.ContentEventResponse{}
	return &resp, nil
}

func (s *server) SendSignalEvent(ctx context.Context, req *pb.SignalEventPayload) (*pb.SignalEventResponse, error) {
	logging.Logf(1, "We've received a signal event. Event: %v", *req)
	inflights := inflights.GetInflights()
	inflights.Insert(*req)
	// logging.Logf(1, "Pool: %#v", pool)
	as := clapi.AmbientStatusPayload{Inflights: inflights.Protobuf()}
	clapiconsumer.SendAmbientStatus(&as)
	// SendAmbients(true)
	resp := pb.SignalEventResponse{}
	return &resp, nil
}

func (s *server) GetUncompiledEntityByKey(ctx context.Context, req *pb.UncompiledEntityByKeyRequest) (*pb.UncompiledEntityByKeyResponse, error) {
	logging.Logf(1, "We've received an uncompiled entity by key request. Event: %v", *req)
	switch req.GetEntityType() {
	case pb.UncompiledEntityType_BOARD:
		var entities []*mimapi.Board

		if boardName := req.GetBoardName(); len(boardName) > 0 {
			entities = beapiconsumer.GetBoardsByName(boardName)
		} else {
			entities = beapiconsumer.GetBoardsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		}
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_BOARD,
			Boards:     entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_THREAD:
		entities := beapiconsumer.GetThreadsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_THREAD,
			Threads:    entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_POST:
		entities := beapiconsumer.GetPostsByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_POST,
			Posts:      entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_VOTE:
		entities := beapiconsumer.GetVotesByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_VOTE,
			Votes:      entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_KEY:
		var entities []*mimapi.Key

		if keyName := req.GetKeyName(); len(keyName) > 0 {
			entities = beapiconsumer.GetKeysByName(keyName)
		} else {
			entities = beapiconsumer.GetKeysByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))

		}
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType: pb.UncompiledEntityType_KEY,
			Keys:       entities,
		}
		return &resp, nil
	case pb.UncompiledEntityType_TRUSTSTATE:
		entities := beapiconsumer.GetTruststatesByKeyFingerprint(req.GetOwnerFingerprint(), int(req.GetLimit()), int(req.GetOffset()))
		resp := pb.UncompiledEntityByKeyResponse{
			EntityType:  pb.UncompiledEntityType_TRUSTSTATE,
			Truststates: entities,
		}
		return &resp, nil
	default:
		return nil, errors.New("Entity type could not be determined")
	}
}

func (s *server) SendInflightsPruneRequest(ctx context.Context, req *pb.InflightsPruneRequest) (*pb.InflightsPruneResponse, error) {
	inflights := inflights.GetInflights()
	inflights.Prune()
	inflights.PushChangesToClient()
	resp := pb.InflightsPruneResponse{}
	return &resp, nil
}

// If we receive ambient status data from the backend, we just forward it directly to the server.
func (s *server) SendBackendAmbientStatus(ctx context.Context, req *pb.BackendAmbientStatusPayload) (*pb.BackendAmbientStatusResponse, error) {
	// Regardless of whether the client can receive it or not, we want to trigger a refresh if the backend has completed a bootstrap.
	if req.GetBackendAmbientStatus().GetTriggerBootstrapRefresh() {
		go func() {
			req.BackendAmbientStatus.TriggerBootstrapRefresh = false // We don't let it go through, because if this is kept true, every send backend ambient status will cause a refresh cycle.
			globals.FrontendTransientConfig.CurrentAmbientStatus.FrontendAmbientStatus.BootstrapRefreshInProgress = true
			globals.FrontendTransientConfig.CurrentAmbientStatus.FrontendAmbientStatus.BootstrapRefreshComplete = false
			if clapiconsumer.ClientIsReadyForConnections {
				clapiconsumer.SendAmbientStatus(&globals.FrontendTransientConfig.CurrentAmbientStatus)
			}
			refresher.Refresh()
			globals.FrontendTransientConfig.CurrentAmbientStatus.FrontendAmbientStatus.BootstrapRefreshInProgress = false
			globals.FrontendTransientConfig.CurrentAmbientStatus.FrontendAmbientStatus.BootstrapRefreshComplete = true
			if clapiconsumer.ClientIsReadyForConnections {
				clapiconsumer.SendAmbientStatus(&globals.FrontendTransientConfig.CurrentAmbientStatus)
			}
		}()
	}
	globals.FrontendTransientConfig.CurrentAmbientStatus.BackendAmbientStatus = req.BackendAmbientStatus
	// for !clapiconsumer.ClientIsReadyForConnections {
	// 	time.Sleep(time.Millisecond * 100)
	// }
	// If the backend has finished a bootstrap, it will ask to trigger a refresh.
	if clapiconsumer.ClientIsReadyForConnections {
		/*
			If the client told us that it is ready, we send this in. Otherwise, the data is already saved into frontend transient config, and it will be sent in with the next ambient status send.
		*/
		clapiconsumer.SendAmbientStatus(&globals.FrontendTransientConfig.CurrentAmbientStatus)
	}
	resp := pb.BackendAmbientStatusResponse{}
	return &resp, nil
}

func (s *server) RequestAmbientStatus(ctx context.Context, req *pb.AmbientStatusRequest) (*pb.AmbientStatusResponse, error) {
	clapiconsumer.SendAmbientStatus(nil)
	resp := pb.AmbientStatusResponse{}
	return &resp, nil
}

func (s *server) RequestHomeView(ctx context.Context, req *pb.HomeViewRequest) (*pb.HomeViewResponse, error) {
	clapiconsumer.SendHomeView()
	resp := pb.HomeViewResponse{}
	return &resp, nil
}

func (s *server) RequestPopularView(ctx context.Context, req *pb.PopularViewRequest) (*pb.PopularViewResponse, error) {
	clapiconsumer.SendPopularView()
	resp := pb.PopularViewResponse{}
	return &resp, nil
}

func (s *server) RequestNewView(ctx context.Context, req *pb.NewViewRequest) (*pb.NewViewResponse, error) {
	clapiconsumer.SendNewView()
	resp := pb.NewViewResponse{}
	return &resp, nil
}

func (s *server) RequestNotifications(ctx context.Context, req *pb.NotificationsRequest) (*pb.NotificationsResponse, error) {
	clapiconsumer.SendNotifications()
	resp := pb.NotificationsResponse{}
	return &resp, nil
}

func (s *server) SetNotificationsSignal(ctx context.Context, req *pb.NotificationsSignalPayload) (*pb.NotificationsSignalResponse, error) {
	if req.GetMarkAllAsRead() {
		festructs.NotificationsSingleton.MarkAllAsRead()
	}
	if req.GetSeen() {
		festructs.NotificationsSingleton.MarkSeen()
	}
	if fp := req.GetReadItemFingerprint(); len(fp) > 0 {
		festructs.NotificationsSingleton.MarkRead(fp)
	}
	clapiconsumer.SendNotifications()
	resp := pb.NotificationsSignalResponse{}
	return &resp, nil
}

func (s *server) SetOnboardComplete(ctx context.Context, req *pb.OnboardCompleteRequest) (*pb.OnboardCompleteResponse, error) {
	globals.FrontendConfig.SetOnboardComplete(req.GetOnboardComplete())
	clapiconsumer.SendOnboardCompleteStatus()
	resp := pb.OnboardCompleteResponse{}
	return &resp, nil
}

func (s *server) SendAddress(ctx context.Context, req *pb.SendAddressPayload) (*pb.SendAddressResponse, error) {
	logging.Logf(1, "We've received a send address request. Event: %v", *req)
	beReq := beapi.ConnectToRemoteRequest{}
	beReq.Address = req.GetAddress()
	sc, errMessage := beapiconsumer.SendConnectToRemoteRequest(&beReq)
	resp := pb.SendAddressResponse{StatusCode: int32(sc), ErrorMessage: errMessage}
	return &resp, nil
}

func (s *server) SendFEConfigChanges(ctx context.Context, req *pb.FEConfigChangesPayload) (*pb.FEConfigChangesResponse, error) {
	logging.Logf(1, "We've received a FE config change request. Event: %v", *req)
	ApplyFEConfigChanges(req)
	clapiconsumer.SendModModeEnabledStatus()
	clapiconsumer.SendAlwaysShowNSFWListStatus()
	clapiconsumer.SendExternalContentAutoloadDisabledStatus()
	clapiconsumer.SendSFWListDisabledStatus()
	resp := pb.FEConfigChangesResponse{}
	return &resp, nil
}

func ApplyFEConfigChanges(req *pb.FEConfigChangesPayload) {
	if req.GetModModeEnabledIsSet() {
		globals.FrontendConfig.SetModModeEnabled(req.GetModModeEnabled())
	}
	if req.GetAlwaysShowNSFWListIsSet() {
		globals.FrontendConfig.SetAlwaysShowNSFWList(req.GetAlwaysShowNSFWList())
	}
	if req.GetExternalContentAutoloadDisabledIsSet() {
		globals.FrontendConfig.SetExternalContentAutoloadDisabled(req.GetExternalContentAutoloadDisabled())
	}
	// Added functionality to be able to turn on and off SFW Lists
	if req.GetSFWListDisabledIsSet() {
		globals.FrontendConfig.SetSFWListDisabled(req.GetSFWListDisabled())
	}
}

func (s *server) RequestBoardReports(ctx context.Context, req *pb.BoardReportsRequest) (*pb.BoardReportsResponse, error) {
	var threadCarriers []festructs.ThreadCarrier

	err := globals.KvInstance.Find("ParentFingerprint", req.GetBoardFingerprint(), &threadCarriers)
	if err != nil {
		logging.Logf(1, "Fetching threads of this board to get the reports failed. Error: %v Board FP: %v", err, req.GetBoardFingerprint())
	}
	var rtes []*feobjects.ReportsTabEntry

	for k, _ := range threadCarriers {
		// Get all reportes threads and posts in this thread carrier
		thrs := getReportedThreads(threadCarriers[k].Threads)
		psts := getReportedPosts(threadCarriers[k].Posts)
		// And convert them to ReportsTabEntries, then protobufs
		for k2, _ := range thrs {
			entry := festructs.NewReportsTabEntryFromThread(&thrs[k2])
			rtes = append(rtes, entry.Protobuf())
		}
		for k3, _ := range psts {
			entry := festructs.NewReportsTabEntryFromPost(&psts[k3])
			rtes = append(rtes, entry.Protobuf())
		}
	}

	resp := pb.BoardReportsResponse{
		ReportsTabEntries: rtes,
	}
	return &resp, nil
}

func getReportedThreads(sl []festructs.CompiledThread) []festructs.CompiledThread {
	var reported []festructs.CompiledThread

	for k, _ := range sl {
		if len(sl[k].CompiledContentSignals.Reports) > 0 && !sl[k].CompiledContentSignals.SelfModIgnored {
			reported = append(reported, sl[k])
		}
	}
	return reported
}

func getReportedPosts(sl []festructs.CompiledPost) []festructs.CompiledPost {
	var reported []festructs.CompiledPost

	for k, _ := range sl {
		if len(sl[k].CompiledContentSignals.Reports) > 0 && !sl[k].CompiledContentSignals.SelfModIgnored {
			reported = append(reported, sl[k])
		}
	}
	return reported
}

func (s *server) RequestBoardModActions(ctx context.Context, req *pb.BoardModActionsRequest) (*pb.BoardModActionsResponse, error) {
	logging.Logf(1, "request board mod actions called.")
	var threadCarriers []festructs.ThreadCarrier

	err := globals.KvInstance.Find("ParentFingerprint", req.GetBoardFingerprint(), &threadCarriers)
	if err != nil {
		logging.Logf(1, "Fetching threads of this board to get the ModActions failed. Error: %v Board FP: %v", err, req.GetBoardFingerprint())
	}
	var rtes []*feobjects.ModActionsTabEntry

	for k, _ := range threadCarriers {
		// Get all ModActioned threads and posts in this thread carrier
		thrs := getModActionedThreads(threadCarriers[k].Threads)
		psts := getModActionedPosts(threadCarriers[k].Posts)
		// And convert them to ModActionsTabEntries, then protobufs
		for k2, _ := range thrs {
			entry := festructs.NewModActionsTabEntryFromThread(&thrs[k2])
			rtes = append(rtes, entry.Protobuf())
		}
		for k3, _ := range psts {
			entry := festructs.NewModActionsTabEntryFromPost(&psts[k3])
			rtes = append(rtes, entry.Protobuf())
		}
	}

	resp := pb.BoardModActionsResponse{
		ModActionsTabEntries: rtes,
	}
	return &resp, nil
}

// For these two below, we also want to do modapprovals in the future, not just modblocks. TODO
func getModActionedThreads(sl []festructs.CompiledThread) []festructs.CompiledThread {
	var modBlocked []festructs.CompiledThread

	for k, _ := range sl {
		if len(sl[k].CompiledContentSignals.ModBlocks) > 0 && !sl[k].CompiledContentSignals.SelfModIgnored {
			modBlocked = append(modBlocked, sl[k])
		}
	}
	return modBlocked
}

func getModActionedPosts(sl []festructs.CompiledPost) []festructs.CompiledPost {
	var modBlocked []festructs.CompiledPost

	for k, _ := range sl {
		if len(sl[k].CompiledContentSignals.ModBlocks) > 0 && !sl[k].CompiledContentSignals.SelfModIgnored {
			modBlocked = append(modBlocked, sl[k])
		}
	}
	return modBlocked
}

// func (s *server) SendMintedUsername_Old(ctx context.Context, req *pb.SendMintedUsernamePayload) (*pb.SendMintedUsernameResponse, error) {
// 	logging.Logf(1, "We've received a send minted username request. Event: %v", *req)
// 	resp := pb.SendMintedUsernameResponse{}
// 	jsonTs := req.GetMintedUsernameRawJSON()
// 	ts := api.Truststate{}
// 	err := json.Unmarshal([]byte(jsonTs), &ts)
// 	tsProto := ts.Protobuf()
// 	tsProtosSlice := []*mimapi.Truststate{&tsProto}
// 	inflights.SendToBackend(interface{}(tsProtosSlice))
// 	if err != nil {
// 		logging.Logf(1, "Unmarshalling from JSON Unique name Truststate failed. Error: %v", err)
// 		return &resp, err
// 	}
// 	return &resp, nil
// }

func (s *server) SendMintedUsernames(ctx context.Context, req *pb.SendMintedUsernamesPayload) (*pb.SendMintedUsernamesResponse, error) {
	logging.Logf(1, "We've received a send minted usernames request. Event: %v", *req)
	resp := pb.SendMintedUsernamesResponse{}
	jsonTs := req.GetMintedUsernamesRawJSON()
	var tses []api.Truststate

	err := json.Unmarshal([]byte(jsonTs), &tses)
	if err != nil {
		logging.Logf(1, "Unmarshalling from JSON Unique name Truststate failed. Error: %v", err)
		return &resp, err
	}
	var tsProtosSlice []*mimapi.Truststate

	for _, ts := range tses {
		tsProto := ts.Protobuf()
		tsProtosSlice = append(tsProtosSlice, &tsProto)
	}
	inflights.SendToBackend(interface{}(tsProtosSlice))
	return &resp, nil
}

func (s *server) SendClientVersion(ctx context.Context, req *pb.ClientVersionPayload) (*pb.ClientVersionResponse, error) {
	logging.Logf(1, "The client is sending us the current client version. Event: %v", *req)
	// Respond with the version we have in the config file
	resp := pb.ClientVersionResponse{}
	resp.LastKnownClientVersion = globals.FrontendConfig.GetLastKnownClientVersion()
	// Put the new version into the config file
	globals.FrontendConfig.SetLastKnownClientVersion(req.GetCurrentClientVersion())
	return &resp, nil
}

func (s *server) SendSearchRequest(ctx context.Context, req *pb.SearchRequestPayload) (*pb.SearchRequestResponse, error) {
	logging.Logf(1, "The client is sending us a search request. Event: %v", *req)
	clapiconsumer.SendSearchResult(req.SearchType, req.SearchQuery)
	// ^ The actual act of computing the search result happens inside this above.
	resp := pb.SearchRequestResponse{}
	return &resp, nil
}
