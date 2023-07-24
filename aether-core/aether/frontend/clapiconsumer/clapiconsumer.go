// Frontend > ClientAPI Client
// This package is the client side of the Client's GRPC API. This is the API frontend uses to send the client some frontend health related information, updates, etc.

package clapiconsumer

import (
	"aether-core/aether/frontend/festructs"
	"aether-core/aether/frontend/kvstore"
	"aether-core/aether/io/api"
	pb "aether-core/aether/protos/clapi"
	"aether-core/aether/protos/feobjects"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"

	// "aether-core/aether/services/scheduling"
	// "aether-core/aether/services/toolbox"
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "sync"
	// "time"
)

func StartClientAPIConnection() (pb.ClientAPIClient, *grpc.ClientConn) {
	clAddr := fmt.Sprint(globals.FrontendConfig.GetClientAPIAddress(), ":", globals.FrontendConfig.GetClientPort())
	conn, err := grpc.Dial(clAddr, grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(toolbox.MaxInt32),
			grpc.MaxCallSendMsgSize(toolbox.MaxInt32),
		),
	)
	if err != nil {
		logging.Logf(1, "Could not connect to the client API service. Error: %v", err)
	}
	c := pb.NewClientAPIClient(conn)
	return c, conn
}

func SendFrontendReady() {
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := pb.FEReadyRequest{
		Address: "127.0.0.1",
		Port:    int32(globals.FrontendConfig.GetFrontendAPIPort()),
	}
	_, err := c.FrontendReady(ctx, &payload)
	if err != nil {
		logging.Logf(1, "SendFrontendReady encountered an error. Err: %v", err)
	}
}

func DeliverAmbients() {
	logging.Logf(2, "Deliver ambients is called, FE>CL, Cl receiver port is: %v", globals.FrontendConfig.GetClientPort())
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := pb.AmbientsRequest{
		Boards: festructs.GetCurrentAmbients().Protobuf(),
	}
	_, err := c.DeliverAmbients(ctx, &payload)
	if err != nil {
		logging.Logf(1, "DeliverAmbients encountered an error. Err: %v", err)
	}
}

var ClientIsReadyForConnections bool

/*
	Note for the future:
	We tried the rate limiter here, and it didn't work. Why? because we aren't diffing and merging the status updates in a way in the frontend, every piece of the app uses its own subset of the thing, and when we override one subset on top of another, what happens is that some other fields that the other set has updated gets overwritten.

	In the future, we need to make it so that a send is just the communication of the whole state, and not a set of intrinsic, different events. When that's ready, we can add a rate limiter here.

	Thankfully, we don't need a rate limiter, it's only a nice-to-have.
*/

// type ambientStatusRateLimiter struct {
// 	lock                       sync.Mutex
// 	LastSend                   int64
// 	SendScheduled              bool
// 	MinimumSendIntervalSeconds int
// }

// var arl ambientStatusRateLimiter = ambientStatusRateLimiter{
// 	MinimumSendIntervalSeconds: 0,
// }

func SendAmbientStatus(cas *pb.AmbientStatusPayload) {
	// arl.lock.Lock()
	// defer arl.lock.Unlock()
	logging.Logf(2, "SendAmbientStatus is called")
	if cas != nil {
		updateAmbientStatus(cas)
		// If it is nil we just use the extant ambient status in fe transient config
	}
	// /*----------  Decide on rate limiting.  ----------*/
	// cutoff := toolbox.CnvToFutureCutoffSeconds(arl.MinimumSendIntervalSeconds)
	// now := time.Now().Unix()
	// if now < cutoff {
	// 	// We haven't had enough time since the last ambient status send.
	// 	logging.Logf(1, "This ambient status send has happened within the rate limiter period and thus rate limited.")
	// 	if arl.SendScheduled {
	// 		logging.Logf(1, "We already have a send scheduled for the next interval. We have nothing to do here, it will automtically send at schedule.")
	// 		return
	// 	}
	// 	// We are rate limited, but the next event isn't scheduled. Schedule it.
	// 	logging.Logf(1, "We got rate limited, but there is no future schedule. placing a future schedule to fire.")
	// 	scheduling.ScheduleOnce(func() {
	// 		logging.Logf(1, "Scheduled send ambient status is running.")
	// 		SendAmbientStatus(cas)
	// 		arl.SendScheduled = false
	// 	}, time.Duration(cutoff-now+1)*time.Second) // +1 to avoid it being too fast and getting denied to reschedule to instant again.
	// 	arl.SendScheduled = true
	// 	return
	// }
	// /*----------  END rate limiting  ----------*/
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	payload := globals.FrontendTransientConfig.CurrentAmbientStatus
	// logging.Logf(1, "flag 1")
	_, err := c.SendAmbientStatus(ctx, &payload)
	// logging.Logf(1, "flag 2")
	if err != nil {
		logging.Logf(1, "SendAmbientStatus encountered an error. Err: %v", err)
	}
}

// updateAmbientStatus partially updates the parts of the live ambient status. So effectively if you make an update to the inflights, this one makes it so that the update doesn't delete the existing but older ambient statuses from backend and frontend.
func updateAmbientStatus(currentAmbientStatus *pb.AmbientStatusPayload) {
	as := globals.FrontendTransientConfig.CurrentAmbientStatus
	if bas := currentAmbientStatus.GetBackendAmbientStatus(); bas != nil {
		as.BackendAmbientStatus = bas
	}
	if fas := currentAmbientStatus.GetFrontendAmbientStatus(); fas != nil {
		as.FrontendAmbientStatus = fas
	}
	if ifl := currentAmbientStatus.GetInflights(); ifl != nil {
		as.Inflights = ifl
	}
	globals.FrontendTransientConfig.CurrentAmbientStatus = as
	// logging.Logf(1, "Current ambient status: %v", as)
}

/*----------  Ambient Local User Data  ----------*/

func SendAmbientLocalUserEntity(localUserExists bool, localUser *feobjects.CompiledUserEntity) {
	logging.Logf(1, "AmbientLocalUserEntity is called")

	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	alu := pb.AmbientLocalUserEntityPayload{}
	alu.LocalUserExists = localUserExists
	alu.LocalUserEntity = localUser
	logging.Logf(1, "alu: %#v", alu)
	_, err := c.SendAmbientLocalUserEntity(ctx, &alu)
	if err != nil {
		logging.Logf(1, "AmbientLocalUserEntity encountered an error. Err: %v", err)
	}
}

/*----------  Higher level methods  ----------*/
/*
These methods aren't 1-1 matches to the gRPC API.
*/

// pushLocalUserAmbient reads from the configstore, and if local user doesn't exist there, bails. If it does, it attempts to read the compiled user header with the same fingerprint. If that fails, the entity exists but not found, and no data is sent until the next attempt.
func PushLocalUserAmbient() {
	alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
	localUserExists := false
	var fp string
	if len(alu) == 0 {
		SendAmbientLocalUserEntity(false, nil)
		return
	}
	localUserExists = true
	var key api.Key
	json.Unmarshal([]byte(alu), &key)
	fp = string(key.Fingerprint)
	uh := festructs.UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in PushLocalUserAmbient>One")
	err := globals.KvInstance.One("Fingerprint", fp, &uh)
	if err != nil {
		logging.Logf(1, "Getting the compiled user entity in PushLocalUserAmbient failed. Error: %v", err)
		// If it exists but not found in the compiled store, that means it hasn't been compiled yet. In this case, we wait and not push anything so that the client can keep its 'loading' state.
		return
	}
	u := festructs.CompiledUser{}
	for key, _ := range uh.Users {
		if uh.Users[key].Fingerprint == fp {
			u = uh.Users[key]
		}
	}
	uproto := u.Protobuf()
	SendAmbientLocalUserEntity(localUserExists, uproto)
	return
}

var FrontendAmbientStatus feobjects.FrontendAmbientStatus

func SendFrontendAmbientStatus() {
	if len(FrontendAmbientStatus.FrontendConfigLocation) == 0 {
		FrontendAmbientStatus.FrontendConfigLocation = globals.GetFrontendConfigLocation()
	}
	FrontendAmbientStatus.SFWListDisabled = globals.FrontendConfig.GetSFWListDisabled()
	as := pb.AmbientStatusPayload{
		FrontendAmbientStatus: &FrontendAmbientStatus,
	}
	SendAmbientStatus(&as)
}

/*----------  Views senders  ----------*/

func SendHomeView() {
	logging.Logf(1, "SendHomeView is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	hvc := festructs.HomeViewCarrier{}
	logging.Logf(3, "Single read happens in SendHomeView>One")
	err := globals.KvInstance.One("Id", 1, &hvc)
	if err != nil {
		logging.Logf(1, "Home view fetch in SendHomeView encountered an error. Error: %v", err)
		return
	}
	var thr []*feobjects.CompiledThreadEntity

	for k, _ := range hvc.Threads {
		thr = append(thr, hvc.Threads[k].Protobuf())
	}
	hvp := pb.HomeViewPayload{Threads: thr}
	_, err2 := c.SendHomeView(ctx, &hvp)
	if err2 != nil {
		logging.Logf(1, "SendHomeView encountered an error. Err: %v", err2)
	}
}

func SendPopularView() {
	logging.Logf(1, "SendPopularView is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	hvc := festructs.PopularViewCarrier{}
	logging.Logf(3, "Single read happens in SendPopularView>One")
	err := globals.KvInstance.One("Id", 1, &hvc)
	if err != nil {
		logging.Logf(1, "Popular view fetch in SendPopularView encountered an error. Error: %v", err)
		return
	}
	var thr []*feobjects.CompiledThreadEntity

	for k, _ := range hvc.Threads {
		thr = append(thr, hvc.Threads[k].Protobuf())
	}
	hvp := pb.PopularViewPayload{Threads: thr}
	_, err2 := c.SendPopularView(ctx, &hvp)
	if err2 != nil {
		logging.Logf(1, "SendPopularView encountered an error. Err: %v", err2)
	}
}

func SendNewView() {
	logging.Logf(1, "SendNewView is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	nvc := festructs.NewViewCarrier{}
	logging.Logf(3, "Single read happens in SendNewView>One")
	err := globals.KvInstance.One("Id", 1, &nvc)
	if err != nil {
		logging.Logf(1, "New view fetch in SendNewView encountered an error. Error: %v", err)
		return
	}
	var thrs []*feobjects.CompiledThreadEntity

	for k, _ := range nvc.Threads {
		protoEntity := nvc.Threads[k].Protobuf()
		// Get board name for thread
		ab := festructs.AmbientBoard{}
		logging.Logf(3, "Single read happens in SendNewView>One>Post>Board")
		err := globals.KvInstance.One("Fingerprint", protoEntity.GetBoard(), &ab)
		if err != nil {
			logging.Logf(1, "Trying to get the parent board of this post in the search results errored out #1. Error: %v ProtoEntity: %#v", err, *protoEntity)
			continue
		}
		protoEntity.ViewMeta_BoardName = ab.Name
		thrs = append(thrs, protoEntity)
	}
	var psts []*feobjects.CompiledPostEntity

	for k, _ := range nvc.Posts {
		protoEntity := nvc.Posts[k].Protobuf()
		// Get board name for post
		ab := festructs.AmbientBoard{}
		logging.Logf(3, "Single read happens in SendNewView>One>Post>Board")
		err := globals.KvInstance.One("Fingerprint", protoEntity.GetBoard(), &ab)
		if err != nil {
			logging.Logf(1, "Trying to get the parent board of this post in the search results errored out #2. Error: %v ProtoEntity: %#v", err, *protoEntity)
			continue
		}
		protoEntity.ViewMeta_BoardName = ab.Name
		// Get thread name for post
		tc := festructs.ThreadCarrier{}
		logging.Logf(3, "Single read happens in SendNewView>One>Post>Thread")
		err2 := globals.KvInstance.One("Fingerprint", protoEntity.GetThread(), &tc)
		if err2 != nil {
			logging.Logf(1, "Trying to get the parent thread of this post in the search results errored out. Error: %v", err2)
			continue
		}
		protoEntity.ViewMeta_ThreadName = tc.Threads[0].Name
		psts = append(psts, protoEntity)
	}
	nvp := pb.NewViewPayload{
		Threads: thrs,
		Posts:   psts,
	}
	_, err2 := c.SendNewView(ctx, &nvp)
	if err2 != nil {
		logging.Logf(1, "SendNewView encountered an error. Err: %v", err2)
	}
}

func SendNotifications() {
	logging.Logf(1, "SendNotifications is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	nList, lastSeen := festructs.NotificationsSingleton.Listify()
	nListProto := nList.Protobuf()
	notifications := pb.NotificationsPayload{Notifications: nListProto, LastSeen: lastSeen}
	_, err2 := c.SendNotifications(ctx, &notifications)
	if err2 != nil {
		logging.Logf(1, "SendNotifications encountered an error. Err: %v", err2)
	}
}

/*----------  Onboarding  ----------*/

func SendOnboardCompleteStatus() {
	logging.Logf(1, "SendOnboardCompleteStatus is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.OnboardCompleteStatusPayload{OnboardComplete: globals.FrontendConfig.GetOnboardComplete()}
	_, err2 := c.SendOnboardCompleteStatus(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendOnboardCompleteStatus encountered an error. Err: %v", err2)
	}
}

/*----------  Mod mode enabled status  ----------*/
func SendModModeEnabledStatus() {
	logging.Logf(1, "SendModModeEnabledStatus is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.ModModeEnabledStatusPayload{ModModeEnabled: globals.FrontendConfig.GetModModeEnabled()}
	_, err2 := c.SendModModeEnabledStatus(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendModModeEnabledStatus encountered an error. Err: %v", err2)
	}
}

/*----------- Always Show NSFW Lists ------------*/
func SendAlwaysShowNSFWListStatus() {
	logging.Logf(1, "SendAlwaysShowNSFWListStatus is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.AlwaysShowNSFWListStatusPayload{AlwaysShowNSFWList: globals.FrontendConfig.GetAlwaysShowNSFWList()}
	_, err2 := c.SendAlwaysShowNSFWListStatus(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendAlwaysShowNSFWListStatus encountered an error. Err: %v", err2)
	}
}

/*----------  External content autoload enabled status  ----------*/
func SendExternalContentAutoloadDisabledStatus() {
	logging.Logf(1, "SendExternalContentAutoloadDisabledStatus is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.ExternalContentAutoloadDisabledStatusPayload{ExternalContentAutoloadDisabled: globals.FrontendConfig.GetExternalContentAutoloadDisabled()}
	_, err2 := c.SendExternalContentAutoloadDisabledStatus(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendExternalContentAutoloadDisabledStatus encountered an error. Err: %v", err2)
	}
}

func SendSFWListDisabledStatus() {
	logging.Logf(1, "SendSFWListDisabledStatus is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.SFWListDisabledStatusPayload{SFWListDisabled: globals.FrontendConfig.GetSFWListDisabled()}
	_, err2 := c.SendSFWListDisabledStatus(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendSFWListDisabledStatus encountered an error. Err: %v", err2)
	}
}

/*----------  Send search results  ----------*/

func SendSearchResult(searchType, searchQuery string) {
	logging.Logf(1, "SendSearchResult is called")
	c, conn := StartClientAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	resp := pb.SearchResultPayload{
		SearchType: searchType,
	}
	/*==============================================
	=            Search for the content            =
	==============================================*/
	switch searchType {
	case "Board":
		r, scoreMap, err := kvstore.SearchBoards(searchQuery)
		if err != nil {
			logging.Logf(1, "This search errored out. Type: %v, Query: %v, Error: %v", searchType, searchQuery, err)
		}
		resp.Boards = r.Protobuf()
		for k, _ := range resp.Boards {
			subbed, notify, lastseen := globals.FrontendConfig.ContentRelations.IsSubbedBoard(resp.Boards[k].Fingerprint)
			whitelisted := globals.FrontendConfig.ContentRelations.SFWList.IsSFWListedBoard(resp.Boards[k].Fingerprint)
			resp.Boards[k].Subscribed = subbed
			resp.Boards[k].Notify = notify
			resp.Boards[k].LastSeen = lastseen
			resp.Boards[k].SFWListed = whitelisted
			resp.Boards[k].ViewMeta_SearchScore = scoreMap[resp.Boards[k].Fingerprint]
		}
	case "Content": // Content = Thread + Post
		posts, threads, scoreMap, err := kvstore.SearchContent(searchQuery)
		if err != nil {
			logging.Logf(1, "This search errored out. Type: %v, Query: %v, Error: %v", searchType, searchQuery, err)
		}
		resp.Threads = threads.Protobuf()
		resp.Posts = posts.Protobuf()
		// Add whitelist data and board name, search score to the threads
		for k, _ := range resp.Threads {
			resp.Threads[k].ViewMeta_SFWListed = globals.FrontendConfig.ContentRelations.SFWList.IsSFWListedBoard(resp.Threads[k].Board)
			ab := festructs.AmbientBoard{}
			logging.Logf(3, "Single read happens in SendSearchResult>One>Thread>Board")
			err := globals.KvInstance.One("Fingerprint", resp.Threads[k].Board, &ab)
			if err != nil {
				logging.Logf(1, "Trying to get the parent board of this thread in the search results errored out. Error: %v", err)
				continue
			}
			resp.Threads[k].ViewMeta_BoardName = ab.Name
			resp.Threads[k].ViewMeta_SearchScore = scoreMap[resp.Threads[k].Fingerprint]
		}

		// Add whitelist data and scores to the posts
		for k, _ := range resp.Posts {
			resp.Posts[k].ViewMeta_SFWListed = globals.FrontendConfig.ContentRelations.SFWList.IsSFWListedBoard(resp.Posts[k].Board)
			// Get board name
			ab := festructs.AmbientBoard{}
			logging.Logf(3, "Single read happens in SendSearchResult>One>Post>Board")
			err := globals.KvInstance.One("Fingerprint", resp.Posts[k].Board, &ab)
			if err != nil {
				logging.Logf(1, "Trying to get the parent board of this post in the search results errored out #3. Error: %v", err)
				continue
			}
			resp.Posts[k].ViewMeta_BoardName = ab.Name
			// Get thread name
			tc := festructs.ThreadCarrier{}
			logging.Logf(3, "Single read happens in SendSearchResult>One>Post>Thread")
			err2 := globals.KvInstance.One("Fingerprint", resp.Posts[k].Thread, &tc)
			if err2 != nil {
				logging.Logf(1, "Trying to get the parent thread of this post in the search results errored out. Error: %v", err2)
				continue
			}
			resp.Posts[k].ViewMeta_ThreadName = tc.Threads[0].Name
			resp.Posts[k].ViewMeta_SearchScore = scoreMap[resp.Posts[k].Fingerprint]
		}
	case "User":
		r, scoreMap, err := kvstore.SearchUsers(searchQuery)
		if err != nil {
			logging.Logf(1, "This search errored out. Type: %v, Query: %v, Error: %v", searchType, searchQuery, err)
		}
		resp.Users = r.Protobuf()
		// Add whitelist data and scores to the posts
		for k, _ := range resp.Users {
			resp.Users[k].ViewMeta_SearchScore = scoreMap[resp.Users[k].Fingerprint]
		}
	default:
		logging.Logf(1, "The search type given by the client is not understood. Given: %v", searchType)
	}
	/*=====  End of Search for the content  ======*/
	_, err2 := c.SendSearchResult(ctx, &resp)
	if err2 != nil {
		logging.Logf(1, "SendSearchResult encountered an error. Err: %v", err2)
	}
}
