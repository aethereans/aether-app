// Frontend > FEStructs > Truststate Signal Reader
// This library is tasked with reading the high level ts signals from the backend. For example, low level read is a truststate read from the backend, a high level read is read of a namemap.

package festructs

import (
	// "aether-core/aether/frontend/festructs"
	// pb "aether-core/aether/protos/beapi"
	"aether-core/aether/frontend/beapiconsumer"
	pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"

	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "golang.org/x/net/context"
	// "google.golang.org/grpc"
	"aether-core/aether/services/metaparse"
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
- [1] publictrust ([2] follow, [2] block) [private trust without releasing signals also works]
- [2] naming ([1] nameassign) [only available for CAs]
- [3] f451 ([1] censorassign) [only available for CAs]
- [4] publicelect ([1] elect [2] disqualify) [private elect without releasing singals also works]

In Types, type=0 is the default state that has no value. so if you want to revert after creation, you revert to 0

In Typeclasses, typeclass=0 is undefined, therefore if you see that, some internal library has sent a malformed request.

If you want to get all types in a typeclass, specify type as -1.
*/

// Three types of query to support: get all since timestamp, get pts for a certain target since timestamp, get all pts for a certain target since timestamp in a certain board.
func GetPTs(targetfp, domainfp string, startts, nowts int64) []PublicTrustSignal {
	rawSignals := getTsBasedSignal(targetfp, domainfp, startts, nowts, 1, -1)
	var sgns []PublicTrustSignal

	for k := range rawSignals {
		sgns = append(sgns, PublicTrustSignal{
			BaseTruststateSignal: BaseTruststateSignal{
				BaseSignal: BaseSignal{
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
				Expiry: rawSignals[k].GetExpiry(),
				Domain: rawSignals[k].GetDomain(),
			},
		})
	}
	return sgns
}

func GetCNs(targetfp, domainfp string, startts, nowts int64) []CanonicalNameSignal {
	rawSignals := getTsBasedSignal(targetfp, domainfp, startts, nowts, 2, -1)
	var sgns []CanonicalNameSignal

	for k := range rawSignals {
		tsmeta, err := metaparse.ReadMeta("Truststate", rawSignals[k].GetMeta())
		if err != nil {
			logging.Logf(2, "We failed to parse this Meta field. Raw Meta field: %v, Entity: %v Error: %v", targetfp, err)
		}
		cname := ""
		if tsmeta != nil {
			cname = tsmeta.(*metaparse.TruststateMeta).CanonicalName
		}
		sgns = append(sgns, CanonicalNameSignal{
			BaseTruststateSignal: BaseTruststateSignal{
				BaseSignal: BaseSignal{
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
				Expiry: rawSignals[k].GetExpiry(),
				Domain: rawSignals[k].GetDomain(),
			},
			CanonicalName: cname,
		})
	}
	return sgns
}

func GetF451s(targetfp, domainfp string, startts, nowts int64) []F451Signal {
	rawSignals := getTsBasedSignal(targetfp, domainfp, startts, nowts, 3, -1)
	var sgns []F451Signal

	for k := range rawSignals {
		sgns = append(sgns, F451Signal{
			BaseTruststateSignal: BaseTruststateSignal{
				BaseSignal: BaseSignal{
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
				Expiry: rawSignals[k].GetExpiry(),
				Domain: rawSignals[k].GetDomain(),
			},
		})
	}
	return sgns
}

func GetPEs(targetfp, domainfp string, startts, nowts int64) []PublicElectSignal {
	rawSignals := getTsBasedSignal(targetfp, domainfp, startts, nowts, 4, -1)
	var sgns []PublicElectSignal

	for k := range rawSignals {
		sgns = append(sgns, PublicElectSignal{
			BaseTruststateSignal: BaseTruststateSignal{
				BaseSignal: BaseSignal{
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
				Expiry: rawSignals[k].GetExpiry(),
				Domain: rawSignals[k].GetDomain(),
			},
		})
	}
	return sgns
}

/*
- [1] publictrust ([2] follow, [2] block) [private trust without releasing signals also works]
- [2] naming ([1] nameassign) [only available for CAs]
- [3] f451 ([1] censorassign) [only available for CAs]
- [4] publicelect ([1] elect [2] disqualify) [private elect without releasing singals also works]
*/

func getTsBasedSignal(targetfp, domainfp string, startts, nowts int64, etypeclass, etype int) []*pbstructs.Truststate {
	return beapiconsumer.GetTruststates(startts, nowts, []string{}, etypeclass, etype, targetfp, domainfp, false, false)
}
