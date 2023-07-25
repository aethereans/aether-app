// Frontend > Inflights > Ingestor

// This service reads from the inflight queue and mints the objects based on the client requests.

package inflights

import (
	"aether-core/aether/frontend/beapiconsumer"
	"aether-core/aether/frontend/clapiconsumer"

	// "aether-core/aether/frontend/festructs"
	"aether-core/aether/frontend/refresher"
	"aether-core/aether/io/api"
	"aether-core/aether/protos/beapi"
	pbstructs "aether-core/aether/protos/mimapi"
	"aether-core/aether/services/create"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"encoding/json"
	"time"
)

/*
Ingest reads through all of the inflights list and mints them. When this is called, it will attempt to consume *all* of the inflights, not a particular segment of it.

Ingest runs when the app is opened, to get rid of anything that might be waiting in the queue when the app was closed the last time. It also runs every time when an entity is added to the ingest queue.

(Committing to kvstore on its own does not automatically trigger an ingest because we're doing it in quite a few places.)
*/
func (o *inflights) Ingest() {
	if o.ingestRunning {
		return
	}
	// The reason we're not using a regular mutex is that we want the other call to terminate, not wait until the mutex is clear, so that the other parts of the application can go forward.
	// o.ingestLock.Lock()
	// defer o.ingestLock.Unlock()
	o.ingestRunning = true
	defer func() {
		o.ingestRunning = false
	}()
	o.ingestRanOnce = true
	// ^ This is useful because we want to run ingest on app open as well.
	logging.Logf(1, "Ingest has started in its goroutine.")
	defer logging.Logf(1, "Ingest is done.")
IngestorLoop:
	for {
		switch e := o.getNextItem().(type) {
		case *InflightBoard:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case *InflightThread:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case *InflightPost:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case *InflightVote:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case *InflightKey:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case *InflightTruststate:
			switch e.Status.EventType {
			case "CREATE":
				e.ingestCreate(o)
			case "UPDATE":
				e.ingestUpdate(o)
			default:
				logging.Logf(1, "The event type of this inflight entity could not be determined by the ingestor. Entity: %#v", e)
			}
		case nil:
			logging.Logf(1, "We've reached the end of the ingest queue. Breaking out of the for loop.")
			break IngestorLoop // If we receive something empty, we just break this.
		default:
			logging.Logf(1, "The type of this inflight entity could not be determined by the ingestor Terminating the ingestor loop. Entity: %#v", e)
			break IngestorLoop
		}
		o.commit()
	}
}

func (o *inflights) getNextItem() interface{} {
	// o.lock.Lock()
	// defer o.lock.Unlock()
	// This is the api used by the minter and it will pull the next item in the queue. it is going to be a pointer, and it won't remove the item from the list. the minter needs to set the state of the next item to 'minting'
	now := time.Now().Unix()
	var oldestEntity interface{}
	oldestTs := now
	for k, _ := range o.InflightBoards {
		if e := o.InflightBoards[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightBoards[k])
		}
	}
	for k, _ := range o.InflightThreads {
		if e := o.InflightThreads[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightThreads[k])
		}
	}
	for k, _ := range o.InflightPosts {
		if e := o.InflightPosts[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightPosts[k])
		}
	}
	for k, _ := range o.InflightVotes {
		if e := o.InflightVotes[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightVotes[k])
		}
	}
	for k, _ := range o.InflightKeys {
		if e := o.InflightKeys[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightKeys[k])
		}
	}
	for k, _ := range o.InflightTruststates {
		if e := o.InflightTruststates[k]; e.Status.RequestedTimestamp <= oldestTs && (e.Status.StatusText == STATUS_WAITING || e.Status.StatusText == STATUS_MINTING) {
			oldestEntity = interface{}(&o.InflightTruststates[k])
		}
	}
	logging.Logf(1, "Returned oldest entity is: %#v", oldestEntity)
	return oldestEntity
}

/*----------  Specific ingest functions  ----------*/

/*----------  Board  ----------*/

/*
These are all fairly similar, but most documentation is on the Key entity since that has a special case.
*/
func (o *InflightBoard) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		e, err := create.CreateBoard(
			o.Entity.GetName(),
			api.Fingerprint(o.Entity.GetOwner()),
			GetLocalUserOwnerPk(o.Entity.GetOwner()),
			[]api.BoardOwner{},
			o.Entity.GetDescription(),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in board creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := e.Protobuf()
		var eps []*pbstructs.Board

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(e.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Board created: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightBoard) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		// Get the original entity from the backend
		fp := string(o.Entity.GetProvable().GetFingerprint())
		fps := []string{fp}
		e := beapiconsumer.GetBoards(0, 0, fps, true, true)
		if len(e) == 0 {
			logging.Logf(1, "We have an entity update request, but the origin entity does not exist in the backend.")
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		entity := api.Board{}
		entity.FillFromProtobuf(*e[0])
		// Create the update request
		ur := create.BoardUpdateRequest{}
		ur.Entity = &entity
		ur.BoardOwnersUpdated = false
		ur.NewBoardOwners = []api.BoardOwner{}
		// ^ FUTURE, this is where the board owner insertion can go in.
		ur.DescriptionUpdated = true
		/*
			Heads up, when we eventually end up with multiple fields that can be updated, we need to make it so that these 'updated' fields are set correctly. Otherwise, updating one field and not touching the rest can accidentally wipe out the rest of the fields.
		*/
		ur.NewDescription = o.Entity.GetDescription()
		err := create.UpdateBoard(ur)
		if err != nil {
			logging.Logf(1, "Minting in board creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&entity))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := entity.Protobuf()
		var eps []*pbstructs.Board

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(entity.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Board updated: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Thread  ----------*/

func (o *InflightThread) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		e, err := create.CreateThread(
			api.Fingerprint(o.Entity.GetBoard()),
			o.Entity.GetName(),
			o.Entity.GetBody(),
			o.Entity.GetLink(),
			api.Fingerprint(o.Entity.GetOwner()),
			GetLocalUserOwnerPk(o.Entity.GetOwner()),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in board update encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := e.Protobuf()
		var eps []*pbstructs.Thread

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(e.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Thread created: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		logging.Logf(1, "Refresher is being called by the recompiling frontend state")
		refresher.Refresh()
		logging.Logf(1, "Refresher called by the recompiling frontend state is done.")
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightThread) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		// Get the original entity from the backend
		fp := string(o.Entity.GetProvable().GetFingerprint())
		fps := []string{fp}
		e := beapiconsumer.GetThreads(0, 0, fps, "", true, true)
		if len(e) == 0 {
			logging.Logf(1, "We have an entity update request, but the origin entity does not exist in the backend.")
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		entity := api.Thread{}
		entity.FillFromProtobuf(*e[0])
		// Create the update request
		ur := create.ThreadUpdateRequest{}
		ur.Entity = &entity
		ur.BodyUpdated = true
		ur.NewBody = o.Entity.GetBody()
		err := create.UpdateThread(ur)
		if err != nil {
			logging.Logf(1, "Minting in thread update encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&entity))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := entity.Protobuf()
		var eps []*pbstructs.Thread

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(entity.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Thread updated: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Post  ----------*/

func (o *InflightPost) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		e, err := create.CreatePost(
			api.Fingerprint(o.Entity.GetBoard()),
			api.Fingerprint(o.Entity.GetThread()),
			api.Fingerprint(o.Entity.GetParent()),
			o.Entity.GetBody(),
			api.Fingerprint(o.Entity.GetOwner()),
			GetLocalUserOwnerPk(o.Entity.GetOwner()),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in board creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := e.Protobuf()
		var eps []*pbstructs.Post

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(e.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Post created: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightPost) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		// Get the original entity from the backend
		fp := string(o.Entity.GetProvable().GetFingerprint())
		fps := []string{fp}
		e := beapiconsumer.GetPosts(0, 0, fps, "", "", true, true)
		if len(e) == 0 {
			logging.Logf(1, "We have an entity update request, but the origin entity does not exist in the backend.")
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		entity := api.Post{}
		entity.FillFromProtobuf(*e[0])
		// Create the update request
		ur := create.PostUpdateRequest{}
		ur.Entity = &entity
		ur.BodyUpdated = true
		ur.NewBody = o.Entity.GetBody()
		err := create.UpdatePost(ur)
		if err != nil {
			logging.Logf(1, "Minting in Post update encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&entity))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := entity.Protobuf()
		var eps []*pbstructs.Post

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(entity.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Post updated: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Vote  ----------*/

func (o *InflightVote) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		e, err := create.CreateVote(
			api.Fingerprint(o.Entity.GetBoard()),
			api.Fingerprint(o.Entity.GetThread()),
			api.Fingerprint(o.Entity.GetTarget()),
			api.Fingerprint(o.Entity.GetOwner()),
			GetLocalUserOwnerPk(o.Entity.GetOwner()),
			int(o.Entity.GetTypeClass()),
			int(o.Entity.GetType()),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in board creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := e.Protobuf()
		var eps []*pbstructs.Vote

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(e.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Vote created: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightVote) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		// Get the original entity from the backend
		fp := string(o.Entity.GetProvable().GetFingerprint())
		fps := []string{fp}
		e := beapiconsumer.GetVotes(0, 0, fps, "", "", "", -1, -1, false, true, true)
		if len(e) == 0 {
			logging.Logf(1, "We have an entity update request, but the origin entity does not exist in the backend.")
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		entity := api.Vote{}
		entity.FillFromProtobuf(*e[0])
		// Create the update request
		ur := create.VoteUpdateRequest{}
		ur.Entity = &entity
		ur.TypeUpdated = true
		ur.NewType = int(o.Entity.GetType())
		err := create.UpdateVote(ur)
		if err != nil {
			logging.Logf(1, "Minting in Vote update encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&entity))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := entity.Protobuf()
		var eps []*pbstructs.Vote

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(entity.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
			o.Status.Update(STATUS_COMPLETE)
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Vote updated: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Key  ----------*/

func (o *InflightKey) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	// ^ This is the prior state - if you're seeing the text here, it has already been completed.
	case STATUS_WAITING, STATUS_MINTING:
		// Nothing has happened yet. Start minting.
		/*
		   This is a special case. The only reason a client would ever create a key entity is that the user does not have a key entity already and is just starting out.

		   That means, we need to check whether that is the case (PK needs to be empty), and that we don't already have a key object for this entity instead (in which case, we just return that and mark this complete.)

		   If the FE doesn't have a local key entity, and this is a local key creation request, we need to mint this key, dehydrate it to JSON and stick it to the frontend config. then we send the result to the client, then to the backend.
		*/
		if !isLocalKeyRequest(o) {
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		if localKeyExists(o) {
			logging.Logf(1, "A local user key already exists.")
			o.Status.Update(STATUS_COMPLETE)
			ifl.PushChangesToClient()
			return
		}
		// We're good. Start minting.
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		key, err := create.CreateKey(
			globals.FrontendConfig.GetMarshaledUserPublicKey(),
			o.Entity.GetName(),
			o.Entity.GetInfo(),
			api.Timestamp(o.Entity.GetExpiry()),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in key creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Minted key!: %#v", key)
		err2 := api.Verify(api.Provable(&key))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		/*----------  Special logic for key (insert to feconfig)  ----------*/
		kJson, err := json.Marshal(key)
		if err != nil {
			logging.Logf(1, "The created local user key could not be converted to JSON. Error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		globals.FrontendConfig.SetDehydratedLocalUserKeyEntity(string(kJson))
		/*----------  Special logic for key (send to Client)  ----------*/
		kp := key.Protobuf()
		var eps []*pbstructs.Key

		eps = append(eps, &kp)
		// Stick it to the key refresher and have it compile this (so that we can get the canonical name and other compiled properties, if any)
		observableUniverse := make(map[string]bool)
		observableUniverse[string(key.Fingerprint)] = true
		refresher.RefreshGlobalUserHeaders(eps, time.Now().Unix(), observableUniverse)
		clapiconsumer.PushLocalUserAmbient()
		// Map the object's fingerprint to the request object, so when the BE sends us status updates for this object, we can find it.
		o.Entity.Provable.Fingerprint = string(key.Fingerprint)
		// The object was minted. Send to backend.
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightKey) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		if !isLocalKeyRequest(o) {
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()

		// Pull the original entity from the frontend config
		alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
		if len(alu) == 0 {
			logging.Logf(1, "We have key update request, but the origin key does not exist in the frontend.")
			return
		}
		var key api.Key
		json.Unmarshal([]byte(alu), &key)

		ur := create.KeyUpdateRequest{}
		ur.Entity = &key
		ur.InfoUpdated = true
		ur.NewInfo = o.Entity.GetInfo()
		ur.ExpiryUpdated = false
		ur.NewExpiry = api.Timestamp(0)
		err := create.UpdateKey(ur)

		if err != nil {
			logging.Logf(1, "Minting in key creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&key))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		/*----------  Special logic for key (insert to feconfig)  ----------*/
		kJson, err := json.Marshal(key)
		if err != nil {
			logging.Logf(1, "The created local user key could not be converted to JSON. Error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		globals.FrontendConfig.SetDehydratedLocalUserKeyEntity(string(kJson))
		/*----------  Special logic for key (send to Client)  ----------*/
		kp := key.Protobuf()
		var eps []*pbstructs.Key

		eps = append(eps, &kp)
		observableUniverse := make(map[string]bool)
		observableUniverse[string(key.Fingerprint)] = true
		refresher.RefreshGlobalUserHeaders(eps, time.Now().Unix(), observableUniverse)
		clapiconsumer.PushLocalUserAmbient()
		o.Entity.Provable.Fingerprint = string(key.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Truststate  ----------*/

func (o *InflightTruststate) ingestCreate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		e, err := create.CreateTruststate(
			api.Fingerprint(o.Entity.GetTarget()),
			api.Fingerprint(o.Entity.GetOwner()),
			GetLocalUserOwnerPk(o.Entity.GetOwner()),
			int(o.Entity.GetTypeClass()),
			int(o.Entity.GetType()),
			api.Fingerprint(o.Entity.GetDomain()),
			api.Timestamp(o.Entity.GetExpiry()),
			o.Entity.GetMeta(),
			"")
		if err != nil {
			logging.Logf(1, "Minting in board creation encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&e))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := e.Protobuf()
		var eps []*pbstructs.Truststate

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(e.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Truststate created: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

func (o *InflightTruststate) ingestUpdate(ifl *inflights) {
	switch o.Status.StatusText {
	case STATUS_WAITING, STATUS_MINTING:
		o.Status.Update(STATUS_MINTING)
		ifl.PushChangesToClient()
		// Get the original entity from the backend
		fp := string(o.Entity.GetProvable().GetFingerprint())
		fps := []string{fp}
		e := beapiconsumer.GetTruststates(0, 0, fps, -1, -1, "", "", true, true)
		if len(e) == 0 {
			logging.Logf(1, "We have an entity update request, but the origin entity does not exist in the backend.")
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		entity := api.Truststate{}
		entity.FillFromProtobuf(*e[0])
		// Create the update request
		ur := create.TruststateUpdateRequest{}
		ur.Entity = &entity
		ur.TypeUpdated = true
		ur.NewType = int(o.Entity.GetType())
		ur.ExpiryUpdated = false
		ur.NewExpiry = api.Timestamp(0)
		err := create.UpdateTruststate(ur)
		if err != nil {
			logging.Logf(1, "Minting in Truststate update encountered an error: %v", err)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		err2 := api.Verify(api.Provable(&entity))
		if err2 != nil {
			logging.Logf(1, "Verification after minting failed. Error: %v", err2)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
			return
		}
		ep := entity.Protobuf()
		var eps []*pbstructs.Truststate

		eps = append(eps, &ep)
		o.Entity.Provable.Fingerprint = string(entity.Fingerprint)
		/*----------  Send to backend  ----------*/
		statusCode := SendToBackend(eps)
		if statusCode == 200 {
			logging.Logf(1, "Successfully inserted into the backend database!")
		} else {
			logging.Logf(1, "Insert to database request failed. Status code provided by the database: %v", statusCode)
			o.Status.Update(STATUS_FAILED)
			ifl.PushChangesToClient()
		}
		logging.Logf(1, "Truststate updated: %#v", e)
		o.Status.Update(STATUS_RECOMPILING_FRONTEND)
		ifl.PushChangesToClient()
		fallthrough
	case STATUS_RECOMPILING_FRONTEND:
		refresher.Refresh()
		o.Status.Update(STATUS_COMPLETE)
		ifl.PushChangesToClient()
	default:
		return
	}
}

/*----------  Specific methods special to the key creation  ----------*/

/*
  These are present because the one reason a client would attempt to create a key is for itself, in that it is just joining the network. All other key creation requests are not useful, so we are gating based on those conditions:

  1 - It is a key creation request for the local PK,
  2 - No key has been created before in this feconfig

*/

func isLocalKeyRequest(o *InflightKey) bool {
	// Check whether this relates to the local key.
	switch o.Status.EventType {
	case "CREATE":
		if len(o.Entity.GetProvable().GetFingerprint()) == 0 {
			return true
		}
	case "UPDATE":
		alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
		if len(alu) == 0 {
			logging.Logf(1, "We have key update request, but the origin key does not exist in the frontend.")
			return false
		}
		var key api.Key
		json.Unmarshal([]byte(alu), &key)
		if o.Entity.GetProvable().GetFingerprint() == string(key.Fingerprint) {
			return true
		}
	default:
		logging.Logf(1, "isLocalKeyRequest failed to determine whether this entity request is a CREATE or UPDATE. Entity: %#v", o)
	}
	return false
}

func localKeyExists(o *InflightKey) bool {
	if len(globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()) > 0 {
		return true
		/*
			^ There's no other condition that would make sense.

			Cases:

			Exists in FE, nowhere else: all is well, we return the existing one

			Exists in BE, nowhere else: FE will assume it doesn't exist and create a new one.

			This case can only happen if a user acts maliciously by copying over an existing PK/secret key pair from an existing feconfig.json, delete and restart the app, and stick those to the new feconfig, and then generate. Since there is no benefit to be gained from this, (since entities are keyed to both key fingerprint *and* the PK, PK being same doesn't mean anything since key fp will still be different), this would be mostly an academic exercise. Nevertheless, in the future, we'll add an enforcement logic path into the backend that will scan all keys and ban all that share a PK, just in case. FUTURE

			Exists in CL, nowhere else: Nonsensical, since CL doesn't have state, the only data source it has is the FE.

			Therefore, checking for whether it exists in FE is enough.
		*/
	}
	return false
}

func GetLocalUserOwnerPk(localUserFp string) string {
	alu := globals.FrontendConfig.GetDehydratedLocalUserKeyEntity()
	if len(alu) == 0 {
		return ""
	}
	var key api.Key
	json.Unmarshal([]byte(alu), &key)
	if string(key.Fingerprint) != localUserFp {
		return ""
		// ^ If the client provides a FP that is not the feconfig local user fp, something went wrong, and we return nothing.
	}
	return key.Key
}

/*----------  Send to DB API layer.  ----------*/

// SendToBackend gets the minted entity, converts it to JSON, and sticks it into the appropriate backend endpoint so that the BE can verify and insert into the database.
func SendToBackend(e interface{}) int {
	// receives []*pbstructs.Key{} style objects
	payload := beapi.MintedContentPayload{}
	switch et := e.(type) {
	case []*pbstructs.Board:
		payload.Boards = et
	case []*pbstructs.Thread:
		payload.Threads = et
	case []*pbstructs.Post:
		payload.Posts = et
	case []*pbstructs.Vote:
		payload.Votes = et
	case []*pbstructs.Key:
		payload.Keys = et
	case []*pbstructs.Truststate:
		payload.Truststates = et
	case []*pbstructs.Address:
		payload.Addresses = et
	}
	return beapiconsumer.SendMintedContent(&payload)
}
