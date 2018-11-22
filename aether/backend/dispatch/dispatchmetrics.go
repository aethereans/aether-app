// Backend > Routines > DispatchMetrics
// This file contains the dispatch routines that dispatch uses to deal with certain cases such as dealing with an update, encountering a new node, etc.

package dispatch

import (
	// "aether-core/aether/backend/responsegenerator"
	"aether-core/aether/io/api"
	"aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	// "aether-core/aether/services/logging"
	// tb "aether-core/aether/services/toolbox"
	// "aether-core/aether/services/verify"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	// "net"
	// "strconv"
	// "strings"
	"time"
)

// Metrics container for every sync.
type CurrentOutboundSyncMetrics struct {
	BoardsReceived                 int
	BoardsSinglePage               bool
	BoardsGETNetworkTime           float64
	BoardsPOSTNetworkTime          float64
	BoardsPOSTTimeToFirstResponse  float64
	BoardsDBCommitTime             float64
	BoardOwnerDBCommitTime         float64
	BoardOwnerDeletionDBCommitTime float64

	ThreadsReceived                int
	ThreadsSinglePage              bool
	ThreadsGETNetworkTime          float64
	ThreadsPOSTNetworkTime         float64
	ThreadsPOSTTimeToFirstResponse float64
	ThreadsDBCommitTime            float64

	PostsReceived                int
	PostsSinglePage              bool
	PostsGETNetworkTime          float64
	PostsPOSTNetworkTime         float64
	PostsPOSTTimeToFirstResponse float64
	PostsDBCommitTime            float64

	VotesReceived                int
	VotesSinglePage              bool
	VotesGETNetworkTime          float64
	VotesPOSTNetworkTime         float64
	VotesPOSTTimeToFirstResponse float64
	VotesDBCommitTime            float64

	KeysReceived                int
	KeysSinglePage              bool
	KeysGETNetworkTime          float64
	KeysPOSTNetworkTime         float64
	KeysPOSTTimeToFirstResponse float64
	KeysDBCommitTime            float64

	TruststatesReceived                int
	TruststatesSinglePage              bool
	TruststatesGETNetworkTime          float64
	TruststatesPOSTNetworkTime         float64
	TruststatesPOSTTimeToFirstResponse float64
	TruststatesDBCommitTime            float64

	AddressesReceived                int
	AddressesSinglePage              bool
	AddressesGETNetworkTime          float64
	AddressesPOSTNetworkTime         float64
	AddressesPOSTTimeToFirstResponse float64
	AddressesDBCommitTime            float64

	MultipleInsertDBCommitTime float64

	LocalIp                 string
	LocalPort               int
	RemoteIp                string
	RemotePort              int
	TotalDurationSeconds    int
	TotalNetworkRemoteWait  float64
	DbInsertDurationSeconds int
	LocalClientName         string
	RemoteClientName        string
	SyncHistory             string
	IsReverseConn           bool
}

func startMetricsContainer(
	resp api.ApiResponse,
	remoteAddr api.Address,
	n persistence.DbNode, isReverseConn bool) *CurrentOutboundSyncMetrics {
	var c CurrentOutboundSyncMetrics
	c.RemoteIp = string(remoteAddr.Location)
	c.RemotePort = int(remoteAddr.Port)
	c.IsReverseConn = isReverseConn
	c.LocalIp = string(globals.BackendConfig.GetExternalIp())
	c.LocalPort = int(globals.BackendConfig.GetExternalPort())
	c.LocalClientName = globals.BackendTransientConfig.AppIdentifier
	c.RemoteClientName = resp.Address.Client.ClientName
	c.SyncHistory = "First sync"
	if n.BoardsLastCheckin > 0 ||
		n.ThreadsLastCheckin > 0 ||
		n.PostsLastCheckin > 0 ||
		n.VotesLastCheckin > 0 ||
		n.KeysLastCheckin > 0 ||
		n.TruststatesLastCheckin > 0 ||
		n.AddressesLastCheckin > 0 {
		c.SyncHistory = "Resync"
	}
	return &c
}

func generateStartMessage(c *CurrentOutboundSyncMetrics, clr *color.Color) string {
	openMessage := clr.Sprintf("\nOPEN: %s:%d (%s) >>> %s:%d (%s) - %s ",
		c.LocalIp, c.LocalPort, c.LocalClientName, c.RemoteIp, c.RemotePort,
		c.RemoteClientName, c.SyncHistory,
	)
	return openMessage
}

func generateCloseMessage(c *CurrentOutboundSyncMetrics, clr *color.Color, ims *[]persistence.InsertMetrics, dur int, extended bool) string {
	im := persistence.InsertMetrics{}
	for _, val := range *ims {
		im.Add(val)
	}
	// spew.Dump(im)
	c.BoardsReceived = im.BoardsReceived
	c.ThreadsReceived = im.ThreadsReceived
	c.PostsReceived = im.PostsReceived
	c.PostsReceived = im.PostsReceived
	c.VotesReceived = im.VotesReceived
	c.KeysReceived = im.KeysReceived
	c.TruststatesReceived = im.TruststatesReceived
	c.AddressesReceived = im.AddressesReceived
	c.TotalDurationSeconds = dur
	c.DbInsertDurationSeconds = im.TimeElapsedSeconds
	c.BoardsDBCommitTime = im.BoardsDBCommitTime
	c.ThreadsDBCommitTime = im.ThreadsDBCommitTime
	c.PostsDBCommitTime = im.PostsDBCommitTime
	c.VotesDBCommitTime = im.VotesDBCommitTime
	c.KeysDBCommitTime = im.KeysDBCommitTime
	c.TruststatesDBCommitTime = im.TruststatesDBCommitTime
	c.AddressesDBCommitTime = im.AddressesDBCommitTime
	c.MultipleInsertDBCommitTime = im.MultipleInsertDBCommitTime

	totalEntitiesReceived := c.BoardsReceived + c.ThreadsReceived + c.PostsReceived + c.VotesReceived + c.KeysReceived + c.TruststatesReceived + c.AddressesReceived
	insertDbDetailString := ""
	if totalEntitiesReceived > 0 {
		insertDbDetailString = fmt.Sprintf("  %d Boards, %d Threads, %d Posts, %d Votes, %d Keys, %d Truststates, %d Addresses (All before dedupe)\n  %s %s %s %s %s %s %s",
			c.BoardsReceived,
			c.ThreadsReceived,
			c.PostsReceived,
			c.VotesReceived,
			c.KeysReceived,
			c.TruststatesReceived,
			c.AddressesReceived,
			singlePageSprinter(c, "Boards"),
			singlePageSprinter(c, "Threads"),
			singlePageSprinter(c, "Posts"),
			singlePageSprinter(c, "Votes"),
			singlePageSprinter(c, "Keys"),
			singlePageSprinter(c, "Truststates"),
			singlePageSprinter(c, "Addresses"))
	}
	var longTimeDetailString string
	var shortTimeDetailString string
	var timeDetailString string
	var reverseConnString string
	if c.IsReverseConn {
		reverseConnString = "(Reverse Connection)"
	}
	if c.TotalDurationSeconds > 0 {
		// This is where we collect db time metrics.
		dbTimeDetailString := fmt.Sprintf("\n    Boards:      %.1fs, \n    Threads:     %.1fs, \n    Posts:       %.1fs, \n    Votes:       %.1fs, \n    Keys:        %.1fs, \n    Truststates: %.1fs, \n    Addresses:   %.1fs, \n    Multitype:   %.1fs (likely through purgatory)", c.BoardsDBCommitTime, c.ThreadsDBCommitTime, c.PostsDBCommitTime, c.VotesDBCommitTime, c.KeysDBCommitTime, c.TruststatesDBCommitTime, c.AddressesDBCommitTime, c.MultipleInsertDBCommitTime)
		// network time metrics for GET and POST.
		networkTimeDetailString := fmt.Sprintf("\n    Boards:      G: %.1fs P: %.1fs (PWait: %.1f), \n    Threads:     G: %.1fs P: %.1fs (PWait: %.1f), \n    Posts:       G: %.1fs P: %.1fs (PWait: %.1f), \n    Votes:       G: %.1fs P: %.1fs (PWait: %.1f), \n    Keys:        G: %.1fs P: %.1fs (PWait: %.1f), \n    Truststates: G: %.1fs P: %.1fs (PWait: %.1f), \n    Addresses:   G: %.1fs P: %.1fs (PWait: %.1f).",
			c.BoardsGETNetworkTime, c.BoardsPOSTNetworkTime, c.BoardsPOSTTimeToFirstResponse,
			c.ThreadsGETNetworkTime, c.ThreadsPOSTNetworkTime, c.ThreadsPOSTTimeToFirstResponse,
			c.PostsGETNetworkTime, c.PostsPOSTNetworkTime, c.PostsPOSTTimeToFirstResponse,
			c.VotesGETNetworkTime, c.VotesPOSTNetworkTime, c.VotesPOSTTimeToFirstResponse,
			c.KeysGETNetworkTime, c.KeysPOSTNetworkTime, c.KeysPOSTTimeToFirstResponse,
			c.TruststatesGETNetworkTime, c.TruststatesPOSTNetworkTime, c.TruststatesPOSTTimeToFirstResponse,
			c.AddressesGETNetworkTime, c.AddressesPOSTNetworkTime, c.AddressesPOSTTimeToFirstResponse,
		)
		longTimeDetailString = fmt.Sprintf("\n  DB: %ds (%s) %s \n  Network: %ds %s", c.DbInsertDurationSeconds, globals.BackendConfig.GetDbEngine(), dbTimeDetailString, c.TotalDurationSeconds-c.DbInsertDurationSeconds, networkTimeDetailString)
		shortTimeDetailString = fmt.Sprintf("\n    DB: %ds (%s)  Network: %ds (Wait for remote: %.1fs)", c.DbInsertDurationSeconds, globals.BackendConfig.GetDbEngine(), c.TotalDurationSeconds-c.DbInsertDurationSeconds, c.TotalNetworkRemoteWait)
		if extended {
			timeDetailString = longTimeDetailString
		} else {
			timeDetailString = shortTimeDetailString
		}
	}
	closeMessage := clr.Sprintf("\nCLOSE: %s >>> %s (%s) %s @ %s\n(%s:%d >>> %s:%d) \nReceived: Total: %d. \n%s \nTime: Total: %ds. %s",
		c.LocalClientName, c.RemoteClientName, c.SyncHistory, reverseConnString,
		time.Now().Format(time.RFC1123),
		c.LocalIp, c.LocalPort, c.RemoteIp, c.RemotePort,
		totalEntitiesReceived, insertDbDetailString,
		c.TotalDurationSeconds, timeDetailString)
	return closeMessage
}

func singlePageSprinter(c *CurrentOutboundSyncMetrics, entityType string) string {
	resp := fmt.Sprintf("single page POST")
	if (entityType == "Boards" && c.BoardsSinglePage) ||
		(entityType == "Threads" && c.ThreadsSinglePage) ||
		(entityType == "Posts" && c.PostsSinglePage) ||
		(entityType == "Votes" && c.VotesSinglePage) ||
		(entityType == "Keys" && c.KeysSinglePage) ||
		(entityType == "Truststates" && c.TruststatesSinglePage) ||
		(entityType == "Addresses" && c.AddressesSinglePage) {
		return fmt.Sprintf("(%s %s)", entityType, resp)
	} else {
		return ""
	}
}
