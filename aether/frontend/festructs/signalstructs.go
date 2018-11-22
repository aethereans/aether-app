package festructs

import ()

/////////////////////////
// Base for Vote and Truststate based signals
/////////////////////////

type BaseSignal struct {
	Fingerprint       string
	Creation          int64
	LastUpdate        int64
	TargetFingerprint string
	SourceFingerprint string
	LastRefreshed     int64
	TypeClass         int
	Type              int
	Self              bool
}

/////////////////////////
// Raw Content Signal forms
/////////////////////////

type BaseVoteSignal BaseSignal

type AddsToDiscussionSignal struct {
	BaseVoteSignal
}

type FollowsGuidelinesSignal struct {
	BaseVoteSignal
	Reason string
}

type ModActionsSignal struct {
	BaseVoteSignal
	Reason string
}

func (s *FollowsGuidelinesSignal) CnvToExplainedSignal() ExplainedSignal {
	e := ExplainedSignal{}
	e.Reason = s.Reason
	e.SourceFp = s.SourceFingerprint
	e.Creation = s.BaseVoteSignal.Creation
	e.LastUpdate = s.BaseVoteSignal.LastUpdate
	return e
}

func (s *ModActionsSignal) CnvToExplainedSignal() ExplainedSignal {
	e := ExplainedSignal{}
	e.Reason = s.Reason
	e.SourceFp = s.SourceFingerprint
	e.Creation = s.BaseVoteSignal.Creation
	e.LastUpdate = s.BaseVoteSignal.LastUpdate
	return e
}

/////////////////////////
// Raw User Signal forms
/////////////////////////

type BaseTruststateSignal struct {
	BaseSignal
	Domain string
	Expiry int64
}

type PublicTrustSignal struct {
	BaseTruststateSignal
}

type CanonicalNameSignal struct {
	BaseTruststateSignal
	CanonicalName string
}

type F451Signal struct {
	BaseTruststateSignal
}

type PublicElectSignal struct {
	BaseTruststateSignal
}
