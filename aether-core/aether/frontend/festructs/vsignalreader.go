// Frontend > FEStructs > Vote Signal Reader
// This library is tasked with reading high level vote-based signals from the backend data store.

package festructs

import (
	// "aether-core/aether/frontend/festructs"
	// pb "aether-core/aether/protos/beapi"
	"aether-core/aether/frontend/beapiconsumer"
	pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/metaparse"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "golang.org/x/net/context"
	// "google.golang.org/grpc"
	// "github.com/willf/bloom"
)

/*
Structure:

- [1] addstodiscussion ([1] upvote, [2] downvote)
   ^ Typeclass           ^ Type      ^ Type

Votes
- [1] addstodiscussion ([1] upvote, [2] downvote)
- [2] followsguidelines ([1] reporttomod)
- [3] modactions ([1] modblock, [2] modapprove)

Truststates
- [1] publictrust ([1] follow, [2] block) [private trust without releasing signals also works]
- [2] naming ([1] nameassign) [only available for CAs]
- [3] f451 ([1] censorassign) [only available for CAs]
- [4] publicelect ([1] elect [2] disqualify) [private elect without releasing singals also works]

In Types, type=0 is the default state that has no value. so if you want to revert after creation, you revert to 0

In Typeclasses, typeclass=0 is undefined, therefore if you see that, some internal library has sent a malformed request.

If you want to get all types in a typeclass, specify type as -1.
*/

// ATD (VC:1), FG (VC:2), MA (VC:3)

// GetATDs gets all AddsToDiscussion type of votes targeting a given entity, or all ATDs whose target is child of a given parent entity.
func GetATDs(parentfp, parenttype, targetfp string, startts, nowts int64, noDescendants bool) []AddsToDiscussionSignal {
	rawSignals := getVoteBasedSignal(parentfp, parenttype, targetfp, startts, nowts, 1, -1, noDescendants)
	var sgns []AddsToDiscussionSignal

	for k, _ := range rawSignals {
		sgns = append(sgns, AddsToDiscussionSignal{
			BaseVoteSignal: BaseVoteSignal{
				Fingerprint:       rawSignals[k].GetProvable().GetFingerprint(),
				Creation:          rawSignals[k].GetProvable().GetCreation(),
				LastUpdate:        rawSignals[k].GetUpdateable().GetLastUpdate(),
				TargetFingerprint: rawSignals[k].GetTarget(),
				SourceFingerprint: rawSignals[k].GetOwner(),
				TypeClass:         int(rawSignals[k].GetTypeClass()),
				Type:              int(rawSignals[k].GetType()),
				Self:              rawSignals[k].GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
				LastRefreshed:     nowts,
			},
		})
	}
	return sgns
}

func GetFGs(parentfp, parenttype, targetfp string, startts, nowts int64, noDescendants bool) []FollowsGuidelinesSignal {
	rawSignals := getVoteBasedSignal(parentfp, parenttype, targetfp, startts, nowts, 2, -1, noDescendants)
	var sgns []FollowsGuidelinesSignal

	for k, _ := range rawSignals {
		vmeta, err := metaparse.ReadMeta("Vote", rawSignals[k].GetMeta())
		if err != nil {
			logging.Logf(2, "We failed to parse this Meta field. Raw Meta field: %v, Entity: %v Error: %v", targetfp, err)
		}
		fgreason := ""
		if vmeta != nil {
			fgreason = vmeta.(*metaparse.VoteMeta).FGReason
		}
		sgns = append(sgns, FollowsGuidelinesSignal{
			BaseVoteSignal: BaseVoteSignal{
				Fingerprint:       rawSignals[k].GetProvable().GetFingerprint(),
				Creation:          rawSignals[k].GetProvable().GetCreation(),
				LastUpdate:        rawSignals[k].GetUpdateable().GetLastUpdate(),
				TargetFingerprint: rawSignals[k].GetTarget(),
				SourceFingerprint: rawSignals[k].GetOwner(),
				TypeClass:         int(rawSignals[k].GetTypeClass()),
				Type:              int(rawSignals[k].GetType()),
				Self:              rawSignals[k].GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
				LastRefreshed:     nowts,
			},
			Reason: fgreason,
		})
	}
	return sgns
}

func GetMAs(parentfp, parenttype, targetfp string, startts, nowts int64, noDescendants bool) []ModActionsSignal {
	rawSignals := getVoteBasedSignal(parentfp, parenttype, targetfp, startts, nowts, 3, -1, noDescendants)
	var sgns []ModActionsSignal

	for k, _ := range rawSignals {
		vmeta, err := metaparse.ReadMeta("Vote", rawSignals[k].GetMeta())
		if err != nil {
			logging.Logf(2, "We failed to parse this Meta field. Raw Meta field: %v, Entity: %v Error: %v", targetfp, err)
		}
		mareason := ""
		if vmeta != nil {
			mareason = vmeta.(*metaparse.VoteMeta).MAReason
		}
		sgns = append(sgns, ModActionsSignal{
			BaseVoteSignal: BaseVoteSignal{
				Fingerprint:       rawSignals[k].GetProvable().GetFingerprint(),
				Creation:          rawSignals[k].GetProvable().GetCreation(),
				LastUpdate:        rawSignals[k].GetUpdateable().GetLastUpdate(),
				TargetFingerprint: rawSignals[k].GetTarget(),
				SourceFingerprint: rawSignals[k].GetOwner(),
				TypeClass:         int(rawSignals[k].GetTypeClass()),
				Type:              int(rawSignals[k].GetType()),
				Self:              rawSignals[k].GetOwnerPublicKey() == globals.FrontendConfig.GetMarshaledUserPublicKey(),
				LastRefreshed:     nowts,
			},
			Reason: mareason,
		})
	}
	return sgns
}

func getVoteBasedSignal(
	parentfp, parenttype, targetfp string,
	startts, nowts int64, etypeclass, etype int, noDescendants bool) []*pbstructs.Vote {
	return beapiconsumer.GetVotes(startts, nowts, []string{}, parentfp, parenttype, targetfp, etypeclass, etype, noDescendants, false, false)
}
