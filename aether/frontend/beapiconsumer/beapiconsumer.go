// Frontend > BackendAPIConsumer
// This package consumes the backend's consumer API.

package beapiconsumer

import (
	pb "aether-core/aether/protos/beapi"
	pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func StartBackendAPIConnection() (pb.BackendAPIClient, *grpc.ClientConn) {
	beAddr := fmt.Sprint(globals.FrontendConfig.GetBackendAPIAddress(), ":", globals.FrontendConfig.GetBackendAPIPort())
	conn, err := grpc.Dial(beAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(toolbox.MaxInt32)))
	if err != nil {
		logging.Logf(1, "Could not connect to the backend API service. Error: %v", err)
	}
	c := pb.NewBackendAPIClient(conn)
	return c, conn
}

func createRequesterId() *pb.RequesterId {
	rid := pb.RequesterId{}
	rid.AccessToken = "" // todo
	rid.Nonce = ""       // todo
	rid.PublicKey = globals.FrontendConfig.GetMarshaledFrontendPublicKey()
	rid.Timestamp = time.Now().Unix()
	return &rid
}

/*
	bypassCache:
		Skip cache even if prefilled for this cycle and ready
	readFromBackendIfEmpty:
		Pull from the backend if the result is empty (the result can legitimately be empty - empty doesn't mean not found in cache)
*/

func GetBoards(start, end int64, fingerprints []string, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Board {
	if servableFromCache(start, end, bypassCache) {
		// If there is a cache that is already established, use that.
		result := queryBoardsCache(cacheQuery{Fingerprints: fingerprints})
		if !readFromBackendIfEmpty {
			// If we haven't asked to pull if empty, return in both empty and filled cases.
			return result
		}
		// If readFromBackendIfEmpty == true and result is > 0, return
		if len(result) > 0 {
			return result
		}
		// If readFromBackendIfEmpty == true and result == 0, fall through and do the regular pull.
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.BoardsRequest{RequesterId: createRequesterId(), Filters: &pb.Filters{LastRefTimeRange: &pb.TimeRange{Start: start, End: end}, Fingerprints: &pb.Fingerprints{Fingerprints: fingerprints}}}
	resp, err := c.GetBoards(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetBoards encountered an error. Error: %v", err)
	}
	r := resp.GetBoards()
	if r != nil {
		return validateBoards(r)
	}
	return []*pbstructs.Board{}
}

func GetThreads(start, end int64, fingerprints []string, parentfp string, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Thread {
	if servableFromCache(start, end, bypassCache) {
		// If there is a cache that is already established, use that.
		result := queryThreadsCache(cacheQuery{Fingerprints: fingerprints, Thread_Board: parentfp})
		if !readFromBackendIfEmpty {
			// If we haven't asked to pull if empty, return in both empty and filled cases.
			return result
		}
		// If readFromBackendIfEmpty == true and result is > 0, return
		if len(result) > 0 {
			return result
		}
		// If readFromBackendIfEmpty == true and result == 0, fall through and do the regular pull.
	} else {
		logging.Logf(2, "This thread request cannot be served from cache. Start: %v, end: %v, readfromBackendifempty: %v, bypasscache:%v", start, end, readFromBackendIfEmpty, bypassCache)
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.ThreadsRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			LastRefTimeRange: &pb.TimeRange{Start: start, End: end},
			Fingerprints:     &pb.Fingerprints{Fingerprints: fingerprints},
			GraphFilters:     &pb.GraphFilters{Board: parentfp},
		}}
	resp, err := c.GetThreads(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetThreads encountered an error. Error: %v", err)
	}
	r := resp.GetThreads()
	if r != nil {
		return validateThreads(r)
	}
	return []*pbstructs.Thread{}
}

func GetPosts(start, end int64, fingerprints []string, parentfp, parenttype string, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Post {
	if servableFromCache(start, end, bypassCache) {
		cq := cacheQuery{
			Fingerprints: fingerprints,
		}
		if parenttype == "thread" {
			cq.Post_Thread = parentfp
		}
		if parenttype == "post" {
			cq.Post_Parent = parentfp
		}
		if parenttype == "board" {
			cq.Post_Board = parentfp
		}
		result := queryPostsCache(cq)
		if !readFromBackendIfEmpty {
			return result
		}
		if len(result) > 0 {
			return result
		}
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.PostsRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			LastRefTimeRange: &pb.TimeRange{
				Start: start, End: end},
			Fingerprints: &pb.Fingerprints{Fingerprints: fingerprints},
			GraphFilters: &pb.GraphFilters{},
		}}
	if parenttype == "board" {
		req.Filters.GraphFilters.Board = parentfp
	}
	/*
		^ This actually happens, and in that case it scans the post's board field.
		There was a 'this can't possibly happen comment here, and it commented out the above. It can. Commenting that out had caused the board user counts to go out of whack, as that made the backend DB return all posts in the DB, not just the ones in that specific board.
	*/
	if parenttype == "thread" {
		req.Filters.GraphFilters.Thread = parentfp
	}
	if parenttype == "post" {
		req.Filters.GraphFilters.Parent = parentfp
	}
	resp, err := c.GetPosts(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetPosts encountered an error. Error: %v", err)
	}
	r := resp.GetPosts()
	if r != nil {
		return validatePosts(r)
	}
	return []*pbstructs.Post{}
}

func GetVotes(start, end int64, fingerprints []string, parentfp, parenttype, targetfp string, etypeclass, etype int, noDescendants bool, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Vote {
	if servableFromCache(start, end, bypassCache) {
		cq := cacheQuery{
			Fingerprints:       fingerprints,
			Vote_Target:        targetfp,
			Vote_TypeClass:     etypeclass,
			Vote_Type:          etype,
			Vote_NoDescendants: noDescendants,
		}
		if parenttype == "board" {
			cq.Vote_Board = parentfp
		}
		if parenttype == "thread" {
			cq.Vote_Thread = parentfp
		}
		result := queryVotesCache(cq, parenttype)
		if !readFromBackendIfEmpty {
			return result
		}
		if len(result) > 0 {
			return result
		}
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.VotesRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			LastRefTimeRange: &pb.TimeRange{
				Start: start,
				End:   end,
			},
			Fingerprints: &pb.Fingerprints{
				Fingerprints: fingerprints,
			},
			TypeFilters: &pb.TypeFilters{
				TypeClass: int32(etypeclass),
				Type:      int32(etype),
			},
			GraphFilters: &pb.GraphFilters{
				Target:        targetfp,
				NoDescendants: noDescendants,
			},
		},
	}
	if parenttype == "board" {
		// votes for all threads and all posts under those threads
		req.Filters.GraphFilters.Board = parentfp
	}
	if parenttype == "thread" {
		// votes for all posts who are children of the given parent thread
		req.Filters.GraphFilters.Thread = parentfp
		// req.Filters.GraphFilters.Target = parentfp
	}
	resp, err := c.GetVotes(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetVotes encountered an error. Error: %v", err)
	}
	r := resp.GetVotes()
	if r != nil {
		return validateVotes(r)
	}
	return []*pbstructs.Vote{}
}

func GetKeys(start, end int64, fingerprints []string, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Key {
	if servableFromCache(start, end, bypassCache) {
		// If there is a cache that is already established, use that.
		result := queryKeysCache(cacheQuery{Fingerprints: fingerprints})
		if !readFromBackendIfEmpty {
			// If we haven't asked to pull if empty, return in both empty and filled cases.
			return result
		}
		// If readFromBackendIfEmpty == true and result is > 0, return
		if len(result) > 0 {
			return result
		}
		// If readFromBackendIfEmpty == true and result == 0, fall through and do the regular pull.
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.KeysRequest{RequesterId: createRequesterId(), Filters: &pb.Filters{LastRefTimeRange: &pb.TimeRange{Start: start, End: end}, Fingerprints: &pb.Fingerprints{Fingerprints: fingerprints}}}
	resp, err := c.GetKeys(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetKeys encountered an error. Error: %v", err)
	}
	r := resp.GetKeys()
	if r != nil {
		return validateKeys(r)
	}
	return []*pbstructs.Key{}
}

func GetTruststates(start, end int64, fingerprints []string, etypeclass, etype int, targetfp, domainfp string, readFromBackendIfEmpty, bypassCache bool) []*pbstructs.Truststate {
	if servableFromCache(start, end, bypassCache) {
		cq := cacheQuery{
			Fingerprints:         fingerprints,
			Truststate_Target:    targetfp,
			Truststate_Domain:    domainfp,
			Truststate_TypeClass: etypeclass,
			Truststate_Type:      etype,
		}
		result := queryTruststatesCache(cq)
		if !readFromBackendIfEmpty {
			return result
		}
		if len(result) > 0 {
			return result
		}
	} else {
		logging.Logf(2, "This truststate request cannot be served from cache. Start: %v, end: %v, fingerprint: %v, etypeclass: %v, etype: %v, targetfp: %v, domainfp: %v, readfromBackendifempty: %v, bypasscache:%v", start, end, fingerprints, etypeclass, etype, targetfp, domainfp, readFromBackendIfEmpty, bypassCache)
	}
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.TruststatesRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			LastRefTimeRange: &pb.TimeRange{
				Start: start,
				End:   end,
			},
			Fingerprints: &pb.Fingerprints{
				Fingerprints: fingerprints,
			},
			TypeFilters: &pb.TypeFilters{
				TypeClass: int32(etypeclass),
				Type:      int32(etype),
			},
			GraphFilters: &pb.GraphFilters{
				Target: targetfp,
				Domain: domainfp,
			},
		},
	}
	resp, err := c.GetTruststates(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetTruststates encountered an error. Error: %v", err)
		// panic("let's see the stack")
	}
	r := resp.GetTruststates()
	if r != nil {
		return validateTruststates(r)
	}
	return []*pbstructs.Truststate{}
}

func GetBoardThreadsCount(fp string) int {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.BoardThreadsCountRequest{
		RequesterId: createRequesterId(),
		Fingerprint: fp}
	resp, err := c.GetBoardThreadsCount(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetBoardThreadsCount encountered an error. Error: %v", err)
	}
	return int(resp.GetCount())
}

func GetThreadPostsCount(fp string) int {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.ThreadPostsCountRequest{
		RequesterId: createRequesterId(),
		Fingerprint: fp}
	resp, err := c.GetThreadPostsCount(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetThreadPostsCount encountered an error. Error: %v", err)
	}
	return int(resp.GetCount())
}

func GetBoardsByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Board {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.BoardsRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetBoards(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetBoards()
	if r != nil {
		return validateBoards(r)
	}
	return []*pbstructs.Board{}
}

func GetThreadsByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Thread {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.ThreadsRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetThreads(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetThreads()
	if r != nil {
		return validateThreads(r)
	}
	return []*pbstructs.Thread{}
}

func GetPostsByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Post {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.PostsRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetPosts(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetPosts()
	if r != nil {
		return validatePosts(r)
	}
	return []*pbstructs.Post{}
}

func GetVotesByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Vote {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.VotesRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetVotes(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetVotes()
	if r != nil {
		return validateVotes(r)
	}
	return []*pbstructs.Vote{}
}

func GetKeysByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Key {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.KeysRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetKeys(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetKeys()
	if r != nil {
		return validateKeys(r)
	}
	return []*pbstructs.Key{}
}

func GetTruststatesByKeyFingerprint(ownerfp string, limit, offset int) []*pbstructs.Truststate {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req := pb.TruststatesRequest{
		RequesterId: createRequesterId(),
		Filters: &pb.Filters{
			GraphFilters: &pb.GraphFilters{
				Owner:         ownerfp,
				NoDescendants: true,
				Limit:         int32(limit),
				Offset:        int32(offset),
			},
		},
	}
	resp, err := c.GetTruststates(ctx, &req)
	if err != nil {
		logging.Logf(1, "GetEntityByKeyFingerprint encountered an error. Error: %v", err)
	}
	r := resp.GetTruststates()
	if r != nil {
		return validateTruststates(r)
	}
	return []*pbstructs.Truststate{}
}

/*----------  Backend minted content intake  ----------*/

func SendMintedContent(req *pb.MintedContentPayload) (statusCode int) {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req.RequesterId = createRequesterId()
	resp, err := c.SendMintedContent(ctx, req)
	if err != nil {
		logging.Logf(1, "SendMintedContent encountered an error. Error: %v", err)
	}
	r := int(resp.GetStatus().GetStatusCode())
	return r
}

/*----------  Backend send connect request  ----------*/

func SendConnectToRemoteRequest(req *pb.ConnectToRemoteRequest) (statusCode int, errorMessage string) {
	c, conn := StartBackendAPIConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), globals.FrontendConfig.GetGRPCServiceTimeout())
	defer cancel()
	req.RequesterId = createRequesterId()
	resp, err := c.SendConnectToRemoteRequest(ctx, req)
	if err != nil {
		logging.Logf(1, "SendConnectToRemoteRequest encountered an error. Error: %v", err)
	}
	r := int(resp.GetStatus().GetStatusCode())
	errMessage := resp.GetStatus().GetErrorMessage()
	return r, errMessage
}
