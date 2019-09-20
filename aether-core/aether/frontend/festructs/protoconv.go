// FEStructs > Protoconv

package festructs

// This is where the conversion methods for our FE entities to protobufs sit.

import (
	pb "aether-core/aether/protos/feobjects"
)

func (e *CompiledBoard) Protobuf() *pb.CompiledBoardEntity {
	return &pb.CompiledBoardEntity{
		Fingerprint:            e.Fingerprint,
		SelfCreated:            e.SelfCreated,
		Name:                   e.Name,
		Description:            e.Description,
		CompiledContentSignals: e.CompiledContentSignals.Protobuf(),
		Owner:                  e.Owner.Protobuf(),
		BoardOwners:            e.BoardOwners,
		Creation:               e.Creation,
		LastUpdate:             e.LastUpdate,
		Meta:                   e.Meta,
		ThreadsCount:           int32(e.ThreadsCount),
		UserCount:              int32(e.UserCount),
	}
}

func (e *CompiledThread) Protobuf() *pb.CompiledThreadEntity {
	return &pb.CompiledThreadEntity{
		Fingerprint:            e.Fingerprint,
		Board:                  e.Board,
		SelfCreated:            e.SelfCreated,
		Name:                   e.Name,
		Body:                   e.Body,
		Link:                   e.Link,
		CompiledContentSignals: e.CompiledContentSignals.Protobuf(),
		Owner:                  e.Owner.Protobuf(),
		Creation:               e.Creation,
		LastUpdate:             e.LastUpdate,
		Meta:                   e.Meta,
		PostsCount:             int32(e.PostsCount),
		Score:                  e.Score,
		ViewMeta_BoardName:     e.ViewMeta_BoardName,
	}
}

func (e *CompiledPost) Protobuf() *pb.CompiledPostEntity {
	return &pb.CompiledPostEntity{
		Fingerprint:            e.Fingerprint,
		Board:                  e.Board,
		Thread:                 e.Thread,
		Parent:                 e.Parent,
		SelfCreated:            e.SelfCreated,
		Body:                   e.Body,
		CompiledContentSignals: e.CompiledContentSignals.Protobuf(),
		Owner:                  e.Owner.Protobuf(),
		Creation:               e.Creation,
		LastUpdate:             e.LastUpdate,
		Meta:                   e.Meta,
	}
}

func (e *CompiledContentSignals) Protobuf() *pb.CompiledContentSignalsEntity {
	return &pb.CompiledContentSignalsEntity{
		TargetFingerprint:  e.TargetFingerprint,
		Upvotes:            int32(e.Upvotes),
		Downvotes:          int32(e.Downvotes),
		SelfUpvoted:        e.SelfUpvoted,
		SelfDownvoted:      e.SelfDownvoted,
		SelfATDFingerprint: e.SelfATDFingerprint,
		SelfATDCreation:    e.SelfATDCreation,
		SelfATDLastUpdate:  e.SelfATDLastUpdate,
		Reports:            ExplainedSignalSliceToProtobuf(e.Reports),
		SelfReported:       e.SelfReported,
		ModBlocks:          ExplainedSignalSliceToProtobuf(e.ModBlocks),
		ModApprovals:       ExplainedSignalSliceToProtobuf(e.ModApprovals),
		SelfModApproved:    e.SelfModApproved,
		SelfModBlocked:     e.SelfModBlocked,
		SelfModIgnored:     e.SelfModIgnored,
		ByMod:              e.ByMod,
		ByFollowedPerson:   e.ByFollowedPerson,
		ByOP:               e.ByOP,
		ModBlocked:         e.ModBlocked,
		ModApproved:        e.ModApproved,
	}
}

func (e *ExplainedSignal) Protobuf() *pb.ExplainedSignalEntity {
	return &pb.ExplainedSignalEntity{
		SourceFp:   e.SourceFp,
		Reason:     e.Reason,
		Creation:   e.Creation,
		LastUpdate: e.LastUpdate,
	}
}

func ExplainedSignalSliceToProtobuf(exps []ExplainedSignal) []*pb.ExplainedSignalEntity {
	if exps == nil {
		return []*pb.ExplainedSignalEntity{}
	}
	pbexps := []*pb.ExplainedSignalEntity{}
	for key, _ := range exps {
		pbexps = append(pbexps, exps[key].Protobuf())
	}
	return pbexps
}

func (e *CompiledUserSignals) Protobuf() *pb.CompiledUserSignalsEntity {
	return &pb.CompiledUserSignalsEntity{
		TargetFingerprint:      e.TargetFingerprint,
		Domain:                 e.Domain,
		FollowedBySelf:         e.FollowedBySelf,
		BlockedBySelf:          e.BlockedBySelf,
		FollowerCount:          int32(e.FollowerCount),
		CanonicalName:          e.CanonicalName,
		CNameSourceFingerprint: e.CNameSourceFingerprint,
		SelfPEFingerprint:      e.SelfPEFingerprint,
		SelfPECreation:         e.SelfPECreation,
		SelfPELastUpdate:       e.SelfPELastUpdate,
		MadeModBySelf:          e.MadeModBySelf,
		MadeNonModBySelf:       e.MadeNonModBySelf,
		MadeModByDefault:       e.MadeModByDefault,
		MadeModByNetwork:       e.MadeModByNetwork,
		MadeNonModByNetwork:    e.MadeNonModByNetwork,
	}
}

func (e *CompiledUser) Protobuf() *pb.CompiledUserEntity {
	return &pb.CompiledUserEntity{
		Fingerprint:         e.Fingerprint,
		NonCanonicalName:    e.NonCanonicalName,
		Info:                e.Info,
		Expiry:              e.Expiry,
		Creation:            e.Creation,
		LastUpdate:          e.LastUpdate,
		LastRefreshed:       e.LastRefreshed,
		Meta:                e.Meta,
		CompiledUserSignals: e.CompiledUserSignals.Protobuf(),
	}
}

func (e *AmbientBoard) Protobuf() *pb.AmbientBoardEntity {
	abe := pb.AmbientBoardEntity{
		Fingerprint:          e.Fingerprint,
		Name:                 e.Name,
		LastUpdate:           e.LastUpdate,
		LastSeen:             e.LastSeen,
		Notify:               e.Notify,
		UserCount:            e.UserCount,
		LastNewThreadArrived: e.LastNewThreadArrived,
	}
	return &abe
}

func (e *AmbientBoardBatch) Protobuf() []*pb.AmbientBoardEntity {
	abes := []*pb.AmbientBoardEntity{}
	for key, _ := range e.Boards {
		abes = append(abes, e.Boards[key].Protobuf())
	}
	return abes
}

func (e *CompiledNotification) Protobuf() *pb.CompiledNotification {
	cnProto := pb.CompiledNotification{
		Type:                    pb.NotificationType(int32(e.Type)),
		Text:                    e.Text,
		ResponsePosts:           e.ResponsePosts,
		ParentThread:            e.ParentThread.Protobuf(),
		ParentPost:              e.ParentPost.Protobuf(),
		CreationTimestamp:       e.CreationTimestamp,
		NewestResponseTimestamp: e.NewestResponseTimestamp,
		Read:                    e.Read,
	}
	if len(e.ResponsePostsUsers) == 1 {
		// We only send this if there is only one user. If there are multiple, we do not send that data.
		u := e.ResponsePostsUsers[e.ResponsePosts[0]]
		unp := []*pb.CUserUsername{u.Protobuf()}
		cnProto.ResponsePostsUsers = unp
	}
	return &cnProto
}

func (e *CUserUsername) Protobuf() *pb.CUserUsername {
	proto := pb.CUserUsername{
		SourceCUser: e.SourceCUser,
		Username:    e.Username,
		Canonical:   e.Canonical,
	}
	return &proto
}

func (e *CNotificationsList) Protobuf() []*pb.CompiledNotification {
	cns := []*pb.CompiledNotification{}
	for key, _ := range *e {
		cns = append(cns, (*e)[key].Protobuf())
	}
	return cns
}

func (e *ReportsTabEntry) Protobuf() *pb.ReportsTabEntry {
	proto := pb.ReportsTabEntry{
		Fingerprint:   e.Fingerprint,
		BoardPayload:  e.BoardPayload.Protobuf(),
		ThreadPayload: e.ThreadPayload.Protobuf(),
		PostPayload:   e.PostPayload.Protobuf(),
		Timestamp:     e.Timestamp,
	}
	return &proto
}

func (e *ModActionsTabEntry) Protobuf() *pb.ModActionsTabEntry {
	proto := pb.ModActionsTabEntry{
		Fingerprint:   e.Fingerprint,
		BoardPayload:  e.BoardPayload.Protobuf(),
		ThreadPayload: e.ThreadPayload.Protobuf(),
		PostPayload:   e.PostPayload.Protobuf(),
		Timestamp:     e.Timestamp,
	}
	return &proto
}

func (e *ReportsTabEntryBatch) Protobuf() []*pb.ReportsTabEntry {
	eProtos := []*pb.ReportsTabEntry{}
	for k, _ := range *e {
		eProtos = append(eProtos, (*e)[k].Protobuf())
	}
	return eProtos
}

/*----------  Batch conversions  ----------*/

func (e *CBoardBatch) Protobuf() []*pb.CompiledBoardEntity {
	eProtos := []*pb.CompiledBoardEntity{}
	for k, _ := range *e {
		eProtos = append(eProtos, (*e)[k].Protobuf())
	}
	return eProtos
}

func (e *CThreadBatch) Protobuf() []*pb.CompiledThreadEntity {
	eProtos := []*pb.CompiledThreadEntity{}
	for k, _ := range *e {
		eProtos = append(eProtos, (*e)[k].Protobuf())
	}
	return eProtos
}

func (e *CPostBatch) Protobuf() []*pb.CompiledPostEntity {
	eProtos := []*pb.CompiledPostEntity{}
	for k, _ := range *e {
		eProtos = append(eProtos, (*e)[k].Protobuf())
	}
	return eProtos
}

func (e *CUserBatch) Protobuf() []*pb.CompiledUserEntity {
	eProtos := []*pb.CompiledUserEntity{}
	for k, _ := range *e {
		eProtos = append(eProtos, (*e)[k].Protobuf())
	}
	return eProtos
}
