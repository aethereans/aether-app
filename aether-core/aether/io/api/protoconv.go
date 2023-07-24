// API > ProtoConv
// This library handles conversions of API Structs from/to protobufs.

package api

import (
	pb "aether-core/aether/protos/mimapi"
)

//////////////////////////////////
// Field set protobuf conversions
//////////////////////////////////

func (fs *ProvableFieldSet) Protobuf() *pb.Provable {
	return &pb.Provable{
		Fingerprint: string(fs.Fingerprint),
		Creation:    int64(fs.Creation),
		ProofOfWork: string(fs.ProofOfWork),
		Signature:   string(fs.Signature),
	}
}

func (ufs *UpdateableFieldSet) Protobuf() *pb.Updateable {
	return &pb.Updateable{
		LastUpdate:        int64(ufs.LastUpdate),
		UpdateProofOfWork: string(ufs.UpdateProofOfWork),
		UpdateSignature:   string(ufs.UpdateSignature),
	}
}

//////////////////////////////////
// Sub-entity protobuf conversions
//////////////////////////////////

func (e *BoardOwner) Protobuf() *pb.BoardOwner {
	return &pb.BoardOwner{
		KeyFingerprint: string(e.KeyFingerprint),
		Expiry:         int64(e.Expiry),
		Level:          int32(e.Level),
	}
}

func (e *Fingerprint) Protobuf() string {
	return string(*e)
}

func (e *Timestamp) Protobuf() int64 {
	return int64(*e)
}

//////////////////////////////////
// Batch protobuf conversions
//////////////////////////////////

func BOSliceToProtobuf(bos []BoardOwner) []*pb.BoardOwner {
	if bos == nil {
		return []*pb.BoardOwner{}
	}
	var pbos []*pb.BoardOwner

	for key, _ := range bos {
		pbos = append(pbos, bos[key].Protobuf())
	}
	return pbos
}

func FPSliceToProtobuf(fps []Fingerprint) []string {
	if fps == nil {
		return []string{}
	}
	var fpstr []string

	for key, _ := range fps {
		fpstr = append(fpstr, string(fps[key]))
	}
	return fpstr
}

//////////////////////////////////
// Core entity to protobuf conversions
//////////////////////////////////

func (e *Board) Protobuf() pb.Board {
	if e.GetVersion() == 1 {
		return pb.Board{
			Provable:       e.ProvableFieldSet.Protobuf(),
			Name:           e.Name,
			BoardOwners:    BOSliceToProtobuf(e.BoardOwners),
			Description:    e.Description,
			Owner:          e.Owner.Protobuf(),
			OwnerPublicKey: e.OwnerPublicKey,
			EntityVersion:  int32(e.EntityVersion),
			Language:       e.Language,
			Meta:           e.Meta,
			RealmId:        e.RealmId.Protobuf(),
			EncrContent:    e.EncrContent,
			Updateable:     e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Board{}
}

func (e *Thread) Protobuf() pb.Thread {
	if e.GetVersion() == 1 {
		return pb.Thread{
			Provable:       e.ProvableFieldSet.Protobuf(),
			Board:          e.Board.Protobuf(),
			Name:           e.Name,
			Body:           e.Body,
			Link:           e.Link,
			Owner:          e.Owner.Protobuf(),
			OwnerPublicKey: e.OwnerPublicKey,
			EntityVersion:  int32(e.EntityVersion),
			Meta:           e.Meta,
			RealmId:        e.RealmId.Protobuf(),
			EncrContent:    e.EncrContent,
			Updateable:     e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Thread{}
}
func (e *Post) Protobuf() pb.Post {
	if e.GetVersion() == 1 {
		return pb.Post{
			Provable:       e.ProvableFieldSet.Protobuf(),
			Board:          e.Board.Protobuf(),
			Thread:         e.Thread.Protobuf(),
			Parent:         e.Parent.Protobuf(),
			Body:           e.Body,
			Owner:          e.Owner.Protobuf(),
			OwnerPublicKey: e.OwnerPublicKey,
			EntityVersion:  int32(e.EntityVersion),
			Meta:           e.Meta,
			RealmId:        e.RealmId.Protobuf(),
			EncrContent:    e.EncrContent,
			Updateable:     e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Post{}
}
func (e *Vote) Protobuf() pb.Vote {
	if e.GetVersion() == 1 {
		return pb.Vote{
			Provable:       e.ProvableFieldSet.Protobuf(),
			Board:          e.Board.Protobuf(),
			Thread:         e.Thread.Protobuf(),
			Target:         e.Target.Protobuf(),
			Owner:          e.Owner.Protobuf(),
			OwnerPublicKey: e.OwnerPublicKey,
			TypeClass:      int32(e.TypeClass),
			Type:           int32(e.Type),
			EntityVersion:  int32(e.EntityVersion),
			Meta:           e.Meta,
			RealmId:        e.RealmId.Protobuf(),
			EncrContent:    e.EncrContent,
			Updateable:     e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Vote{}
}
func (e *Key) Protobuf() pb.Key {
	if e.GetVersion() == 1 {
		return pb.Key{
			Provable:      e.ProvableFieldSet.Protobuf(),
			Type:          e.Type,
			Key:           e.Key,
			Expiry:        e.Expiry.Protobuf(),
			Name:          e.Name,
			Info:          e.Info,
			EntityVersion: int32(e.EntityVersion),
			Meta:          e.Meta,
			RealmId:       e.RealmId.Protobuf(),
			EncrContent:   e.EncrContent,
			Updateable:    e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Key{}
}
func (e *Truststate) Protobuf() pb.Truststate {
	if e.GetVersion() == 1 {
		return pb.Truststate{
			Provable:       e.ProvableFieldSet.Protobuf(),
			Target:         e.Target.Protobuf(),
			Owner:          e.Owner.Protobuf(),
			OwnerPublicKey: e.OwnerPublicKey,
			TypeClass:      int32(e.TypeClass),
			Type:           int32(e.Type),
			Domain:         e.Domain.Protobuf(),
			Expiry:         e.Expiry.Protobuf(),
			EntityVersion:  int32(e.EntityVersion),
			Meta:           e.Meta,
			RealmId:        e.RealmId.Protobuf(),
			EncrContent:    e.EncrContent,
			Updateable:     e.UpdateableFieldSet.Protobuf(),
		}
	}
	return pb.Truststate{}
}

// Protobuf > API object conversions

//////////////////////////////////
// Protobuf > Field set
//////////////////////////////////

func (e *ProvableFieldSet) FillFromProtobuf(v pb.Provable) {
	e.Fingerprint = Fingerprint(v.Fingerprint)
	e.Creation = Timestamp(v.Creation)
	e.ProofOfWork = ProofOfWork(v.ProofOfWork)
	e.Signature = Signature(v.Signature)
}

func (e *UpdateableFieldSet) FillFromProtobuf(v pb.Updateable) {
	e.LastUpdate = Timestamp(v.LastUpdate)
	e.UpdateProofOfWork = ProofOfWork(v.UpdateProofOfWork)
	e.UpdateSignature = Signature(v.UpdateSignature)
}

//////////////////////////////////
// Protobuf > Sub-entity
//////////////////////////////////

func (e *BoardOwner) FillFromProtobuf(v pb.BoardOwner) {
	e.KeyFingerprint = Fingerprint(v.KeyFingerprint)
	e.Expiry = Timestamp(v.Expiry)
	e.Level = uint8(v.Level)
}

func (e *Fingerprint) FillFromProtobuf(v string) {
	*e = Fingerprint(v)
}

func (e *Timestamp) FillFromProtobuf(v int64) {
	*e = Timestamp(v)
}

//////////////////////////////////
// Protobuf > Sub-entity (batch)
//////////////////////////////////

func BoardOwnerSliceProtoToAPI(pbos []*pb.BoardOwner) []BoardOwner {
	if pbos == nil {
		return []BoardOwner{}
	}
	var abos []BoardOwner

	for key, _ := range pbos {
		bo := BoardOwner{}
		bo.FillFromProtobuf(*pbos[key])
		abos = append(abos, bo)
	}
	return abos
}

func FPSliceProtoToAPI(pfps []string) []Fingerprint {
	if pfps == nil {
		return []Fingerprint{}
	}
	var afps []Fingerprint

	for key, _ := range pfps {
		var fp Fingerprint
		fp.FillFromProtobuf(pfps[key])
		afps = append(afps, fp)
	}
	return afps
}

//////////////////////////////////
// Protobuf > Entity
//////////////////////////////////

func (e *Board) FillFromProtobuf(v pb.Board) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Name = v.GetName()
		e.BoardOwners = BoardOwnerSliceProtoToAPI(v.GetBoardOwners())
		e.Description = v.GetDescription()
		e.Owner = Fingerprint(v.GetOwner())
		e.OwnerPublicKey = v.GetOwnerPublicKey()
		e.EntityVersion = int(v.GetEntityVersion())
		e.Language = v.GetLanguage()
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}

func (e *Thread) FillFromProtobuf(v pb.Thread) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Board = Fingerprint(v.GetBoard())
		e.Name = v.GetName()
		e.Body = v.GetBody()
		e.Link = v.GetLink()
		e.Owner = Fingerprint(v.GetOwner())
		e.OwnerPublicKey = v.GetOwnerPublicKey()
		e.EntityVersion = int(v.GetEntityVersion())
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}

func (e *Post) FillFromProtobuf(v pb.Post) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Board = Fingerprint(v.GetBoard())
		e.Thread = Fingerprint(v.GetThread())
		e.Parent = Fingerprint(v.GetParent())
		e.Body = v.GetBody()
		e.Owner = Fingerprint(v.GetOwner())
		e.OwnerPublicKey = v.GetOwnerPublicKey()
		e.EntityVersion = int(v.GetEntityVersion())
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}

func (e *Vote) FillFromProtobuf(v pb.Vote) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Board = Fingerprint(v.GetBoard())
		e.Thread = Fingerprint(v.GetThread())
		e.Target = Fingerprint(v.GetTarget())
		e.Owner = Fingerprint(v.GetOwner())
		e.OwnerPublicKey = v.GetOwnerPublicKey()
		e.Type = int(v.GetType())
		e.TypeClass = int(v.GetTypeClass())
		e.EntityVersion = int(v.GetEntityVersion())
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}

func (e *Key) FillFromProtobuf(v pb.Key) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Type = v.GetType()
		e.Key = v.GetKey()
		e.Expiry = Timestamp(v.GetExpiry())
		e.Name = v.GetName()
		e.Info = v.GetInfo()
		e.EntityVersion = int(v.GetEntityVersion())
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}

func (e *Truststate) FillFromProtobuf(v pb.Truststate) {
	if v.GetEntityVersion() == 1 {
		pv := ProvableFieldSet{}
		pv.FillFromProtobuf(*v.GetProvable())
		e.ProvableFieldSet = pv
		e.Target = Fingerprint(v.GetTarget())
		e.Owner = Fingerprint(v.GetOwner())
		e.OwnerPublicKey = v.GetOwnerPublicKey()
		e.Type = int(v.GetType())
		e.TypeClass = int(v.GetTypeClass())
		e.Domain = Fingerprint(v.GetDomain())
		e.Expiry = Timestamp(v.GetExpiry())
		e.EntityVersion = int(v.GetEntityVersion())
		e.Meta = v.GetMeta()
		e.RealmId = Fingerprint(v.GetRealmId())
		e.EncrContent = v.GetEncrContent()
		u := UpdateableFieldSet{}
		u.FillFromProtobuf(*v.GetUpdateable())
		e.UpdateableFieldSet = u
	}
}
