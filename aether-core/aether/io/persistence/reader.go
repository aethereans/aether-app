// Persistence > Reader
// This file collects all functions that read from a database. The Server uses this API, as well as the UI.

package persistence

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"

	// "aether-core/aether/services/toolbox"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// ReadNode provides the ability to seek a specific node.
func ReadNode(fingerprint api.Fingerprint) (DbNode, error) {
	var n DbNode
	if len(fingerprint) > 0 {
		query, args, err := sqlx.In("SELECT * FROM Nodes WHERE Fingerprint IN (?);", fingerprint)
		if err != nil {
			return n, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		if err != nil {
			return n, err
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			err := rows.StructScan(&n)
			if err != nil {
				return n, err
			}
		}
		rows.Close()
	}
	if len(n.Fingerprint) == 0 {
		return n, fmt.Errorf("The node you have asked for could not be found. You asked for: %s", fingerprint)
	}
	return n, nil
}

// enforceReadValidity enforces that, in a ReadX function (medium level API below), either a time range or a list of fingerprints are asked, and not both.
func enforceReadValidity(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp, privilegedSource bool) error {
	if privilegedSource {
		return nil
	}
	valid := false
	/*
		IF:
			A: BeginTS blank AND EndTS blank AND fingerprint filter extant: GOOD. Fingerprint search.

			B: BeginTS OR EndTS is blank (or both are not), and fingerprint filter blank: GOOD. One-way bounded  or two-way bounded time range search.

			C: All fields are blank: GOOD. Time range search where BeginTS is last cache generation timestamp or network head.

			D: Anything else: BAD.
	*/
	if (beginTimestamp == 0 && endTimestamp == 0 && len(fingerprints) > 0) ||
		((beginTimestamp != 0 || endTimestamp != 0) && len(fingerprints) == 0) ||
		(beginTimestamp == 0 && endTimestamp == 0 && len(fingerprints) == 0) {
		valid = true
	}
	if !valid {
		return fmt.Errorf("You can either search for a time range, or for fingerprint(s). You can't do both or neither at the same time - you have to do one. Asked fingerprints: %#v, BeginTimestamp: %s, EndTimestamp: %s", fingerprints, strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(endTimestamp)))
	}
	return nil
}

// SanitiseTimeRange validates and cleans the time range used in ReadX functions (medium level API below)
func SanitiseTimeRange(
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	now api.Timestamp,
	privilegedSource bool,
) (api.Timestamp, api.Timestamp, error) {
	// Heads up: the difference between privileged and unprivileged paths is that the privileged path doesn't restrict the oldest data given to last cache generation timestamp, or to the network head, whichever is newer. It will return any and all data that is available in the backend.
	if privilegedSource {
		// If begin is newer than now, make it now
		if beginTimestamp > now {
			beginTimestamp = now
		}
		// if end is zero, or newer than now, make it now
		if endTimestamp == 0 || endTimestamp > now {
			endTimestamp = now
		}
		// if begin is larger than end, flip
		if beginTimestamp > endTimestamp {
			beginTimestamp, endTimestamp = endTimestamp, beginTimestamp
		}
		if beginTimestamp == endTimestamp {
			// Heads up, this request *can* return a result, assuming there are entities whose LastRef is exactly that second. So this is an unlikely, but not an invalid configuration.
		}
		// In summary, begin can be <zero-now>, end can be <zero-now>, and we flip begin and end if it's a negative range. This means we don't gate based on last cache generation timestamp, nor do we gate based on the network head, network memory, or local memory. Simply, if if exists in the database, it will be returned.
		return beginTimestamp, endTimestamp, nil
	}
	// If not privileged:
	// If there is no end timestamp, the end is now.
	if endTimestamp == 0 || endTimestamp > now {
		endTimestamp = now
	}
	// If the begin is newer than the end, flip. We haven't started to enforce limits yet, so the change here will be entirely coming from the remote.
	if beginTimestamp > endTimestamp {
		return beginTimestamp, endTimestamp, fmt.Errorf("Your BeginTimestamp is larger than your EndTimestamp. BeginTimestamp: %s, EndTimestamp: %s", strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(endTimestamp)))
	}
	// Internal processing starts.

	// If beginTimestamp is older than our last cache, start from the end of the last cache.
	if beginTimestamp < api.Timestamp(globals.BackendConfig.GetLastCacheGenerationTimestamp()) {
		beginTimestamp = api.Timestamp(globals.BackendConfig.GetLastCacheGenerationTimestamp())
		endTimestamp = now // Because, in thecase of begin 3 and end 5, begin going to 145000000 will make begin much bigger than end. Prevent that by moving the end also.
	}
	if beginTimestamp == 0 {
		// If there are no caches, lastCache will be 0 and this will return everything in the database. To prevent this, we limit the results to the duration of the network head.
		nhd := globals.BackendConfig.GetNetworkHeadDays()
		delta := time.Duration(nhd) * time.Hour * 24
		beginTimestamp = api.Timestamp(time.Now().Add(-delta).Unix())
	}
	//If beginTimestamp is in the future, return error.
	if beginTimestamp > now {
		return beginTimestamp, endTimestamp, fmt.Errorf("Your beginTimestamp is in the future. BeginTimestamp: %s, Now: %s", strconv.Itoa(int(beginTimestamp)), strconv.Itoa(int(now)))
	}
	// End of internal processing
	// After we do these things, if we end up with a begin timestamp that is newer than the end, the end timestamp will be 'now'. This can happen in the case where both the start and end timestamps are within the cached period.
	if beginTimestamp > endTimestamp {
		// This is in the case a start is given but an end is given empty, in which case it's from start to now.
		endTimestamp = now
	}
	return beginTimestamp, endTimestamp, nil
}

type OptionalReadInputs struct {
	// Board
	Board_Name string
	// Thread
	Thread_Board string
	// Post
	Post_Board  string
	Post_Thread string
	Post_Parent string
	// Vote
	Vote_Board         string
	Vote_Thread        string
	Vote_Target        string
	Vote_TypeClass     int
	Vote_Type          int
	Vote_NoDescendants bool
	// Key
	Key_Name string
	// Truststate
	Truststate_Target    string
	Truststate_Domain    string
	Truststate_TypeClass int
	Truststate_Type      int
	// All provables (all except Address)
	// Heads up / limit / offsets are only defined over owner key based search. I'll basically have to build a query builder for SQL to make this work with all other options, and I'd rather not have to do that. if you end up needing limit / offset on something else, just write the query here.
	AllProvables_Owner  string
	AllProvables_Limit  int
	AllProvables_Offset int
}

// Read is the high level API for DB reads. It provides filtering support. It can return multiple types if requested by the embeds.
func Read(
	entityType string, // boards, threads, posts, votes, addresses, keys, truststates
	fingerprints []api.Fingerprint,
	embeds []string,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	privilegedSource bool, // coming from a frontend or from an internal request not associated with serving an external one. If enabled, this will unlock expensive requests. It will enable being able to query past last cache generation timestamp, and querying with *both* fingerprints and time ranges, not just only one of those.
	opts *OptionalReadInputs,
) (api.Response, error) {
	if globals.BackendTransientConfig.ShutdownInitiated {
		return api.Response{}, nil
	}
	if opts == nil {
		opts = &OptionalReadInputs{ // -1: not specified (0 is a valid value for type, means it was reverted to default from something else.)
			// Board
			// Thread
			Thread_Board: "",
			// Post
			Post_Board:  "",
			Post_Thread: "",
			Post_Parent: "",
			// Vote
			Vote_Board:         "",
			Vote_Thread:        "",
			Vote_Target:        "",
			Vote_TypeClass:     -1,
			Vote_Type:          -1,
			Vote_NoDescendants: false,
			// Key
			// Truststate
			Truststate_Target:    "",
			Truststate_Domain:    "",
			Truststate_TypeClass: -1,
			Truststate_Type:      -1,
		}
	}
	var result api.Response
	now := api.Timestamp(time.Now().Unix())
	// Fingerprints search and start/end timestamp search are mutually exclusive. Make sure that is enforced.
	err := enforceReadValidity(fingerprints, beginTimestamp, endTimestamp, privilegedSource)
	if err != nil {
		return result, err
	}
	sanitisedBeginTimestamp, sanitisedEndTimestamp, err2 := SanitiseTimeRange(beginTimestamp, endTimestamp, now, privilegedSource)
	if err2 != nil {
		return result, err2
	}
	// This thing below is for embeds. This is the container within which we fill the fingerprints for the item requested (board fps, etc.) as []api.Provable. It's used to do []api.Board to []api.Provable transition essentially.
	var provableArr []api.Provable
	// Now we switch based on the entity type.
	switch entityType {
	case "boards":
		entities, err := ReadBoards(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.AllProvables_Owner, opts.Board_Name, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Boards = entities
		// Convert the result to []api.Provable

		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}

	case "threads":
		entities, err := ReadThreads(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.Thread_Board, opts.AllProvables_Owner, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Threads = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "posts":
		entities, err := ReadPosts(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.Post_Board, opts.Post_Thread, opts.Post_Parent, opts.AllProvables_Owner, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Posts = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}

	case "votes":
		entities, err := ReadVotes(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.Vote_TypeClass, opts.Vote_Type, opts.Vote_Board, opts.Vote_Thread, opts.Vote_Target, opts.Vote_NoDescendants, opts.AllProvables_Owner, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Votes = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "addresses":
		return result, errors.New(fmt.Sprint("You tried to supply an address into the high level Read API. This API only provides reads for entities that fulfil the api.Provable interface. Please use ReadAddress directly."))
	case "keys":
		entities, err := ReadKeys(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.AllProvables_Owner, opts.Key_Name, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Keys = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	case "truststates":
		entities, err := ReadTruststates(fingerprints, sanitisedBeginTimestamp, sanitisedEndTimestamp, opts.Truststate_TypeClass, opts.Truststate_Type, opts.Truststate_Target, opts.Truststate_Domain, opts.AllProvables_Owner, opts.AllProvables_Limit, opts.AllProvables_Offset)
		if err != nil {
			return result, err
		}
		result.Truststates = entities
		// Convert the result to []api.Provable
		for i, _ := range entities {
			provableArr = append(provableArr, &entities[i])
		}
	}
	// We deal with filling the embedded fields. Embed handler has all the code for the different types of embeds.
	embedErr := handleEmbeds(provableArr, &result, embeds)
	if embedErr != nil {
		return result, embedErr
	}
	return result, nil
}

// Embeds APIs. These provide the supporting embeds for the primary calls at the medium level. This is consumed by the high level API above.

// Embed helpers used by the Read high level API above.

func existsInEmbed(asked string, embeds []string) bool {
	if len(embeds) > 0 {
		for i, _ := range embeds {
			if embeds[i] == asked {
				return true
			}
		}
	}
	return false
}

func handleEmbeds(entities []api.Provable, result *api.Response, embeds []string) error {
	// This holds the results of the first embed so we can add it to the main results before sending it into the keys. See the comment below for context.
	var firstEmbedCache []api.Provable
	if existsInEmbed("threads", embeds) {
		thr, err := ReadThreadEmbed(entities)
		if err != nil {
			return err
		}
		result.Threads = thr
		for i, _ := range thr {
			firstEmbedCache = append(firstEmbedCache, &thr[i])
		}

	}

	if existsInEmbed("posts", embeds) {
		posts, err := ReadPostEmbed(entities)
		if err != nil {
			return err
		}
		result.Posts = posts
		for i, _ := range posts {
			firstEmbedCache = append(firstEmbedCache, &posts[i])
		}
	}
	if existsInEmbed("votes", embeds) {
		votes, err := ReadVoteEmbed(entities)
		if err != nil {
			return err
		}
		result.Votes = votes
		for i, _ := range votes {
			firstEmbedCache = append(firstEmbedCache, &votes[i])
		}
	}
	// Keys being at the end is significant. This is always the last processed embed, because, for example, if you do Board with Threads embed, both boards and threads will have keys to refer to, and without the thread keys the data will be incomplete. So the keys need to read both the original data, and the data coming from any other embeds when providing the keys.

	//TOTHINK: Embedding needs to be able to go five levels deep. Boards, (Threads, Posts, Votes, Keys, Truststates). This is useful when a node needs a copy of a board, and the whole board and anything it links can be dissected from another node and sent over.

	// But for the time being, let's treat the key as a special case, as if the key are not fully provided, the embeds with more than one layer don't work at all. The embedded objects will not be able to be validated otherwise. The five layer embed thing could be constructed from a series of queries but the absence of keys for the first embed is a serious problem.
	if existsInEmbed("keys", embeds) {
		keys, err := ReadKeyEmbed(entities, firstEmbedCache) // <- firstEmbedCache
		if err != nil {
			return err
		}
		result.Keys = keys
	}
	return nil
}

// Core read embeds.

// ReadThreadEmbed gets the threads linked from the entities provided.
// Only available for: Boards
func ReadThreadEmbed(entities []api.Provable) ([]api.Thread, error) {
	var arr []api.Thread
	var dbArr []DbThread
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the thread embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Board:
		// Only defined for boards. No other entity has thread embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Threads WHERE Board IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var entity DbThread
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Thread{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Thread))
		}
	}
	return arr, nil
}

// ReadPostEmbed gets the posts linked from the existing entities provided.
// Only available for: Threads
func ReadPostEmbed(entities []api.Provable) ([]api.Post, error) {
	var arr []api.Post
	var dbArr []DbPost
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the post embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Thread:
		// Only defined for threads. No other entity has post embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Posts WHERE Thread IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var entity DbPost
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Post{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Post))
		}
	}
	return arr, nil
}

// ReadVoteEmbed gets the votes linking to the entities provided.
// Only available for: Posts
func ReadVoteEmbed(entities []api.Provable) ([]api.Vote, error) {
	var arr []api.Vote
	var dbArr []DbVote
	var entityFingerprints []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the vote embed is empty."))
		return arr, nil
	}
	switch entity := entities[0].(type) {
	// entity: typed API object.
	case *api.Post:
		// Only defined for posts. No other entity has vote embeds.
		entity = entity // Stop complaining
		for i, _ := range entities {
			entityFingerprints = append(entityFingerprints, entities[i].GetFingerprint())
		}
		query, args, err := sqlx.In("SELECT DISTINCT * FROM Votes WHERE Target IN (?);", entityFingerprints)
		if err != nil {
			return arr, err
		}
		rows, err := globals.DbInstance.Queryx(query, args...)
		if err != nil {
			return arr, err
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var entity DbVote
			err := rows.StructScan(&entity)
			if err != nil {
				return []api.Vote{}, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
		for _, entity := range dbArr {
			apiEntity, err := DBtoAPI(entity)
			if err != nil {
				// Log the problem and go to the next iteration without saving this one.
				logging.Log(1, err)
				continue
			}
			arr = append(arr, apiEntity.(api.Vote))
		}
	}
	return arr, nil
}

// ReadKeyEmbed gets the keys linked from the existing entities provided.
// Only available for: Boards, Threads, Posts, Truststates
func ReadKeyEmbed(entities []api.Provable, firstEmbedCache []api.Provable) ([]api.Key, error) {
	var arr []api.Key
	var dbArr []DbKey
	var entityOwners []api.Fingerprint
	if len(entities) == 0 {
		logging.Log(1, fmt.Sprintf("The entities list given to the key embed is empty."))
		return arr, nil
	}
	entities = append(entities, firstEmbedCache...)
	for i, _ := range entities {
		switch entity := entities[i].(type) {
		// entity: typed API object.
		case *api.Board:
			entityOwners = append(entityOwners, entity.GetOwner())
			for j, _ := range entity.BoardOwners {
				entityOwners = append(entityOwners, entity.BoardOwners[j].KeyFingerprint)
			}
		case *api.Thread, *api.Post, *api.Truststate:
			entityOwners = append(entityOwners, entity.GetOwner())
		}
	}
	// The thing below is the same as read keys.
	query, args, err := sqlx.In("SELECT DISTINCT * FROM PublicKeys WHERE Fingerprint IN (?);", entityOwners)
	if err != nil {
		return arr, err
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return arr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbKey
		err := rows.StructScan(&entity)
		if err != nil {
			return []api.Key{}, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close() // Close it ASAP, do not call any new DB queries while still in this.
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Key))
	}
	return arr, nil
}

// Medium Level API. You should not use these directly. Use the high level API (above) because that one has the embed support and returns a proper api.Response object.

// ReadBoards reads threads from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadBoards(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	ownerfp, name string, limit, offset int,
) ([]api.Board, error) {
	var arr []api.Board
	dbArr, err := ReadDbBoards(fingerprints, beginTimestamp, endTimestamp, ownerfp, name, limit, offset)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Board))
	}
	return arr, nil
}

// ReadDbBoards returns the search result in DB form. This layer has a few fields like LocalArrival and LastReferenced not exposed to the API layer that allows for internal decision making.
func ReadDbBoards(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	ownerfp, name string, limit, offset int) ([]DbBoard, error) {
	var dbArr []DbBoard
	var query string
	var args []interface{}
	var err error
	opts := reqtype(reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         -1,
		typ:            -1,
		ownerFp:        ownerfp,
		name:           name,
		limit:          limit,
		offset:         offset,
	})
	switch opts {
	case "(name)":
		query, args, err = sqlx.In("SELECT * FROM Boards WHERE Name IN (?);", name)
		if err != nil {
			return dbArr, err
		}
	case "(fp)(ts)":
		query, args, err = sqlx.In("SELECT * FROM Boards WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)":
		query, args, err = sqlx.In("SELECT * FROM Boards WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)":
		query, args, err = sqlx.In("SELECT DISTINCT * from Boards WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from Boards WHERE (Owner = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from Boards WHERE (Owner = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbBoards was invalid. Fps: %v, Start: %v, End: %v", fingerprints, beginTimestamp, endTimestamp)
		return dbArr, fmt.Errorf("The request you've made to ReadDbBoards was invalid. Fps: %v, Start: %v, End: %v All opts: %#v", fingerprints, beginTimestamp, endTimestamp, opts)
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbBoard
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

// ReadThreads reads threads from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadThreads(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	boardfp string,
	ownerfp string, limit, offset int,
) ([]api.Thread, error) {
	var arr []api.Thread
	dbArr, err := ReadDbThreads(fingerprints, beginTimestamp, endTimestamp, boardfp, ownerfp, limit, offset)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Thread))
	}
	return arr, nil
}

func ReadDbThreads(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	boardfp string,
	ownerfp string, limit, offset int,
) ([]DbThread, error) {
	// logging.Logf(1, "read db threads hits with owner: %v", ownerfp)
	var dbArr []DbThread
	var query string
	var args []interface{}
	var err error
	reqtyp := reqtype(reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         -1,
		typ:            -1,
		parentBoardFp:  boardfp,
		ownerFp:        ownerfp,
		limit:          limit,
		offset:         offset,
	})
	switch reqtyp {
	case "(fp)(ts)":
		query, args, err = sqlx.In("SELECT * FROM Threads WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)":
		query, args, err = sqlx.In("SELECT * FROM Threads WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)":
		query, args, err = sqlx.In("SELECT DISTINCT * from Threads WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(pbfp)": // lastref + parent board fingerprint
		// Hot path, removed some niceties (select distinct, order by) that we don't necessarily use for faster pull
		query, args, err = sqlx.In("SELECT * from Threads WHERE (LastReferenced >= ? AND LastReferenced <= ? ) AND Board = ?", beginTimestamp, endTimestamp, boardfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from Threads WHERE (Owner = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from Threads WHERE (Owner = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbThreads was invalid. Reqtype: %v, Fps: %v, Start: %v, End: %v", reqtyp, fingerprints, beginTimestamp, endTimestamp)
		return dbArr, fmt.Errorf("The request you've made to ReadDbThreads was invalid. Fps: %v, Start: %v, End: %v", fingerprints, beginTimestamp, endTimestamp)
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbThread
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

// ReadPosts reads posts from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadPosts(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	boardfp, threadfp, parentfp string,
	ownerfp string, limit, offset int,
) ([]api.Post, error) {
	var arr []api.Post
	dbArr, err := ReadDbPosts(fingerprints, beginTimestamp, endTimestamp, boardfp, threadfp, parentfp, ownerfp, limit, offset)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Post))
	}
	return arr, nil
}

func ReadDbPosts(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	boardfp, threadfp, parentfp string,
	ownerfp string, limit, offset int,
) ([]DbPost, error) {
	var dbArr []DbPost
	var query string
	var args []interface{}
	var err error
	opts := reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         -1,
		typ:            -1,
		parentBoardFp:  boardfp,
		parentThreadFp: threadfp,
		parentPostFp:   parentfp,
		ownerFp:        ownerfp,
		limit:          limit,
		offset:         offset,
	}
	rtype := reqtype(opts)
	switch rtype {
	case "(fp)(ts)":
		query, args, err = sqlx.In("SELECT * FROM Posts WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)":
		query, args, err = sqlx.In("SELECT * FROM Posts WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)":
		query, args, err = sqlx.In("SELECT DISTINCT * from Posts WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(pbfp)":
		query, args, err = sqlx.In("SELECT * from Posts WHERE (LastReferenced >= ? AND LastReferenced <= ? ) AND (Board = ?)", beginTimestamp, endTimestamp, boardfp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(ptfp)":
		query, args, err = sqlx.In("SELECT * from Posts WHERE (LastReferenced >= ? AND LastReferenced <= ? ) AND (Thread = ?)", beginTimestamp, endTimestamp, threadfp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(ppfp)":
		query, args, err = sqlx.In("SELECT * from Posts WHERE (LastReferenced >= ? AND LastReferenced <= ? ) AND (Parent = ?)", beginTimestamp, endTimestamp, parentfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from Posts WHERE (Owner = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from Posts WHERE (Owner = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbPosts was invalid. Fps: %v, Start: %v, End: %v", fingerprints, beginTimestamp, endTimestamp)
		return dbArr, fmt.Errorf("The request you've made to ReadDbPosts was invalid. Fps: %v, Start: %v, End: %v", fingerprints, beginTimestamp, endTimestamp)
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbPost
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

// ReadVotes reads votes from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadVotes(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	vtypeclass, vtype int,
	boardfp, threadfp, targetfp string,
	noDescendants bool,
	ownerfp string, limit, offset int,
) ([]api.Vote, error) {
	var arr []api.Vote
	dbArr, err := ReadDbVotes(
		fingerprints, beginTimestamp, endTimestamp,
		vtypeclass, vtype,
		boardfp, threadfp, targetfp,
		noDescendants, ownerfp, limit, offset,
	)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Vote))
	}
	return arr, nil
}

func ReadDbVotes(
	fingerprints []api.Fingerprint,
	beginTimestamp, endTimestamp api.Timestamp,
	vtypeclass, vtype int,
	boardfp, threadfp, targetfp string,
	noDescendants bool,
	ownerfp string, limit, offset int,
) ([]DbVote, error) {
	var dbArr []DbVote
	var query string
	var args []interface{}
	var err error
	opts := reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         vtypeclass,
		typ:            vtype,
		parentBoardFp:  boardfp,
		parentThreadFp: threadfp,
		targetFp:       targetfp,
		noDescendants:  noDescendants,
		ownerFp:        ownerfp,
		limit:          limit,
		offset:         offset,
	}
	reqtyp := reqtype(opts)
	// fmt.Println("reqtype for votes:")
	// fmt.Println(reqtyp)
	// fmt.Println("reqtype opts:")
	// fmt.Println(opts)
	switch reqtyp {
	case "(fp)(ts)": // list of fingerprints, timestamps
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)": // list of fingerprints
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)": // timestamps
		query, args, err = sqlx.In("SELECT DISTINCT * from Votes WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)": // timestamps, typeclass
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", vtypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(ty)": //timestamps, typeclass, type
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (TypeClass = ? AND Type = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", vtypeclass, vtype, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(tc)(ty)": // typeclass, type
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (TypeClass = ? AND Type = ?);", vtypeclass, vtype)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(pbfp)": // timestamps, typeclass, parent board fp
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (Board = ?) AND (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", boardfp, vtypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(ptfp)": // timestamps, typeclass, parent thread fp
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (Thread = ?) AND (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", threadfp, vtypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(tafp)": // timestamps, typeclass, target fp
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (Target = ?) AND (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", targetfp, vtypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(pbfp)(nodesc)": // timestamps, typeclass, parent board fp (useful when you want to get only votes to thread entity, but not to its posts)
		query, args, err = sqlx.In("SELECT * FROM Votes WHERE (Board = ?) AND (Target = Thread) AND (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", boardfp, vtypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from Votes WHERE (Owner = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from Votes WHERE (Owner = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbVotes was invalid. Fps: %v, Start: %v, End: %v, ReqType: %v", fingerprints, beginTimestamp, endTimestamp, reqtyp)
		return dbArr, fmt.Errorf("The request you've made to ReadDbVotes was invalid. Fps: %v, Start: %v, End: %v ReqType: %v", fingerprints, beginTimestamp, endTimestamp, reqtyp)
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbVote
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

func readDbAddressesBasicSearch(Location api.Location, Sublocation api.Location, Port uint16) (*[]DbAddress, error) {
	var dbArr []DbAddress
	if len(Location) > 0 && Port > 0 { // Regular address search.
		rows, err := globals.DbInstance.Queryx("SELECT * from Addresses WHERE Location = ? AND Sublocation = ? AND Port = ?", Location, Sublocation, Port)
		if err != nil {
			return &dbArr, err
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var entity DbAddress
			err := rows.StructScan(&entity)
			if err != nil {
				return &dbArr, err
			}
			dbArr = append(dbArr, entity)
		}
		rows.Close()
	}
	return &dbArr, nil
}

func readDbAddressesFirstXResultsSearch(maxResults int, offset int, addrType uint8) (*[]DbAddress, error) {
	var dbArr []DbAddress
	query := ""
	var rows *sqlx.Rows
	var err error
	if maxResults == 0 {
		query = "SELECT * from Addresses WHERE AddressType = ? ORDER BY LocalArrival DESC OFFSET ?"
		rows, err = globals.DbInstance.Queryx(query, addrType, offset)
	} else if maxResults > 0 {
		query = "SELECT * from Addresses WHERE AddressType = ? ORDER BY LocalArrival DESC LIMIT ? OFFSET ?"
		rows, err = globals.DbInstance.Queryx(query, addrType, maxResults, offset)
	} else {
		// if negative value
		return &dbArr, errors.New("You've provided a negative maxResults value to address search.")
	}
	// First X results search.
	// You have to provide a addrtype, if you search for 0, that will find the nodes you haven't connected yet.
	if err != nil {
		return &dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbAddress
		err := rows.StructScan(&entity)
		if err != nil {
			return &dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return &dbArr, nil
}

func readDbAddressesTimeRangeSearch(
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	offset int,
	searchType string) (*[]DbAddress, error) {
	var dbArr []DbAddress
	// Time range search
	// This should result in:
	// - Entities that has landed to local after the beginning and before the end
	// If the end timestamp is 0, it's assumed that endTs is right now.
	var endTs api.Timestamp
	if endTimestamp == 0 {
		endTs = api.Timestamp(time.Now().Unix())
	} else {
		endTs = endTimestamp
	}
	var rangeSearchColumn string
	// Options: all, connected.
	if searchType == "timerange_all" {
		// Default case. Default is LocalArrival.
		rangeSearchColumn = "LocalArrival"
	} else if searchType == "timerange_lastsuccessfulping" {
		rangeSearchColumn = "LastSuccessfulPing"
	} else if searchType == "timerange_lastsuccessfulsync" {
		rangeSearchColumn = "LastSuccessfulSync"
	} else {
		return &dbArr, fmt.Errorf("You have provided an invalid time range search type. You provided: %s", searchType)
	}
	query := fmt.Sprintf("SELECT DISTINCT * from Addresses WHERE (%s > ? AND %s < ?) ORDER BY %s DESC", rangeSearchColumn, rangeSearchColumn, rangeSearchColumn)
	rows, err := globals.DbInstance.Queryx(query, beginTimestamp, endTs)
	if err != nil {
		return &dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbAddress
		err := rows.StructScan(&entity)
		if err != nil {
			return &dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return &dbArr, nil
}

// readAddressContainerResponse is a direct prepared response to a container generation request, whether be it for a cache generation or POST response generation.
func readAddressContainerResponse(beg, end api.Timestamp, addrType uint8, limit int) (*[]DbAddress, error) {
	// Filter by time range given, sort by fixed elements, and limit the results to limit
	var results []DbAddress

	q := "SELECT * FROM Addresses WHERE (LastSuccessfulPing > ? AND LastSuccessfulPing < ? AND AddressType = ?) ORDER BY LastSuccessfulSync DESC, LastSuccessfulPing DESC LIMIT ?"
	r, err := globals.DbInstance.Queryx(q, beg, end, addrType, limit)
	defer r.Close() // In case of premature exit.
	if err != nil {
		return &results, err
	}
	for r.Next() {
		var a DbAddress
		err := r.StructScan(&a)
		if err != nil {
			return &results, err
		}
		results = append(results, a)
	}
	r.Close()
	// logging.Logf(2, "Addr Response generated to serve in POST: %#v", Dbg_convertAddrSliceToNameSlice(results))
	return &results, nil
}

func readDBAddressesAll(isDesc bool) (*[]DbAddress, error) {
	var results []DbAddress

	q := ""
	if isDesc {
		q = "SELECT * FROM Addresses ORDER BY LastSuccessfulSync DESC, LastSuccessfulPing DESC, LocalArrival DESC LIMIT ?"
	} else {
		q = "SELECT * FROM Addresses ORDER BY LastSuccessfulSync ASC, LastSuccessfulPing ASC, LocalArrival ASC LIMIT ?"
	}
	r, err := globals.DbInstance.Queryx(q, globals.BackendConfig.MaxAddressTableSize)
	defer r.Close()
	if err != nil {
		return &results, err
	}
	for r.Next() {
		var a DbAddress
		err := r.StructScan(&a)
		if err != nil {
			return &results, err
		}
		results = append(results, a)
	}
	r.Close()
	return &results, nil
}

func ReadDbAddresses(
	loc, subloc api.Location, port uint16,
	beg, end api.Timestamp, limit, offset int,
	addrType uint8, searchType string) (*[]DbAddress, error) {
	if searchType == "container_generate" {
		live, err1 := readAddressContainerResponse(beg, end, 2, (limit/10)*8)
		bs, err2 := readAddressContainerResponse(beg, end, 3, (limit / 10))
		static, err3 := readAddressContainerResponse(beg, end, 255, (limit / 10))
		var all []DbAddress

		all = append(all, (*live)...)
		all = append(all, (*bs)...)
		all = append(all, (*static)...)
		if err1 != nil || err2 != nil || err3 != nil {
			errs := []error{}
			return &all, fmt.Errorf("Some errors appeared while trying to prepare this address response to a remote request. Errors: %#v", errs)
		}
		return &all, nil
	} else if searchType == "basic" {
		return readDbAddressesBasicSearch(loc, subloc, port)
	} else if searchType == "limit" {
		return readDbAddressesFirstXResultsSearch(limit, offset, addrType)
	} else if searchType == "timerange_all" ||
		searchType == "timerange_lastsuccessfulping" ||
		searchType == "timerange_lastsuccessfulsync" {
		return readDbAddressesTimeRangeSearch(beg, end, offset, searchType)
	} else if searchType == "all_desc" {
		return readDBAddressesAll(true)
	} else if searchType == "all_asc" {
		return readDBAddressesAll(false)
	} else {
		return &[]DbAddress{}, errors.New("You have requested data from ReadAddresses in an invalid configuration.")
	}
}

// ReadAddresses reads addresses from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
// FUTURE: these should eventually do a join on subprotocols to be able to filter by subprotocol items. But for now, since we're not using that, there's no join, which is faster.
func ReadAddresses(
	loc, subloc api.Location, port uint16,
	beg, end api.Timestamp, limit, offset int,
	addrType uint8, searchType string) ([]api.Address, error) {
	var arr []api.Address
	res, err := ReadDbAddresses(loc, subloc, port, beg, end, limit, offset, addrType, searchType)
	if err != nil {
		return arr, err
	}
	for _, entity := range *res {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Address))
	}
	return arr, nil
}

// ReadKeys reads keys from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadKeys(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	ownerfp, name string, limit, offset int,
) ([]api.Key, error) {
	var arr []api.Key
	dbArr, err := ReadDbKeys(fingerprints, beginTimestamp, endTimestamp, ownerfp, name, limit, offset)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Key))
	}
	return arr, nil
}

func ReadDbKeys(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	ownerfp, name string, limit, offset int,
) ([]DbKey, error) {
	var dbArr []DbKey
	var query string
	var args []interface{}
	var err error
	opts := reqtype(reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         -1,
		typ:            -1,
		ownerFp:        ownerfp,
		name:           name,
		limit:          limit,
		offset:         offset,
	})
	switch opts {
	case "(name)":
		query, args, err = sqlx.In("SELECT * FROM PublicKeys WHERE Name IN (?);", name)
		if err != nil {
			return dbArr, err
		}
	case "(fp)(ts)":
		query, args, err = sqlx.In("SELECT * FROM PublicKeys WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)":
		query, args, err = sqlx.In("SELECT * FROM PublicKeys WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)":
		query, args, err = sqlx.In("SELECT DISTINCT * from PublicKeys WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from PublicKeys WHERE (Fingerprint = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from PublicKeys WHERE (Fingerprint = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbKeys was invalid. Fps: %v, Start: %v, End: %v", fingerprints, beginTimestamp, endTimestamp)
		return dbArr, fmt.Errorf("The request you've made to ReadDbKeys was invalid. Fps: %v, Start: %v, End: %v All opts: %#v", fingerprints, beginTimestamp, endTimestamp, opts)
	}
	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbKey
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

// ReadTrustStates reads trust states from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.
func ReadTruststates(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	tstypeclass, tstype int,
	targetfp, domainfp string,
	ownerfp string, limit, offset int,
) ([]api.Truststate, error) {
	var arr []api.Truststate
	dbArr, err := ReadDbTruststates(fingerprints, beginTimestamp, endTimestamp, tstypeclass, tstype, targetfp, domainfp, ownerfp, limit, offset)
	if err != nil {
		return arr, err
	}
	for _, entity := range dbArr {
		apiEntity, err := DBtoAPI(entity)
		if err != nil {
			// Log the problem and go to the next iteration without saving this one.
			logging.Log(1, err)
			continue
		}
		arr = append(arr, apiEntity.(api.Truststate))
	}
	return arr, nil
}

func ReadDbTruststates(
	fingerprints []api.Fingerprint,
	beginTimestamp api.Timestamp,
	endTimestamp api.Timestamp,
	tstypeclass, tstype int,
	targetfp, domainfp string,
	ownerfp string, limit, offset int,
) ([]DbTruststate, error) {
	var dbArr []DbTruststate
	var query string
	var args []interface{}
	var err error
	opts := reqtypeOpts{
		fingerprints:   fingerprints,
		beginTimestamp: beginTimestamp,
		endTimestamp:   endTimestamp,
		tclass:         tstypeclass,
		typ:            tstype,
		targetFp:       targetfp,
		domainFp:       domainfp,
		ownerFp:        ownerfp,
		limit:          limit,
		offset:         offset,
	}
	switch reqtype(opts) {
	case "(fp)(ts)": // fingerprint + timespan
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE Fingerprint IN (?) AND (LastReferenced >= ? AND LastReferenced <= ?);", fingerprints, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(fp)": // fingerprint
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE Fingerprint IN (?);", fingerprints)
		if err != nil {
			return dbArr, err
		}
	case "(ts)": // timespan
		query, args, err = sqlx.In("SELECT DISTINCT * from Truststates WHERE (LastReferenced >= ? AND LastReferenced <= ? ) ORDER BY LastReferenced DESC", beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)": // timespan + typeclass
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE (TypeClass = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", tstypeclass, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(ty)": // timespan + typeclass + type
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE (TypeClass = ? AND Type = ?) AND (LastReferenced >= ? AND LastReferenced <= ?);", tstypeclass, tstype, beginTimestamp, endTimestamp)
		if err != nil {
			return dbArr, err
		}
	case "(tc)(ty)": // typeclass + type
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE (TypeClass = ? AND Type = ?);", tstypeclass, tstype)
		if err != nil {
			return dbArr, err
		}
	case "(tc)(tafp)": // typeclass + target
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE TypeClass = ? AND Target = ?;", tstypeclass, targetfp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(tafp)": // timespan + typeclass + target
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE TypeClass = ? AND (LastReferenced >= ? AND LastReferenced <= ?) AND Target = ?;", tstypeclass, beginTimestamp, endTimestamp, targetfp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(dofp)": // timespan + typeclass + target
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE TypeClass = ? AND (LastReferenced >= ? AND LastReferenced <= ?) AND Domain = ?;", tstypeclass, beginTimestamp, endTimestamp, domainfp)
		if err != nil {
			return dbArr, err
		}
	case "(tc)(tafp)(dofp)": // typeclass + target + domain
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE TypeClass = ? AND Target = ? AND Domain = ?;", tstypeclass, targetfp, domainfp)
		if err != nil {
			return dbArr, err
		}
	case "(ts)(tc)(tafp)(dofp)": // timespan + typeclass + target + domain
		query, args, err = sqlx.In("SELECT * FROM Truststates WHERE TypeClass = ? AND (LastReferenced >= ? AND LastReferenced <= ?) AND Target = ? AND Domain = ?;", tstypeclass, beginTimestamp, endTimestamp, targetfp, domainfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)":
		query, args, err = sqlx.In("SELECT * from Truststates WHERE (Owner = ?) ORDER BY LastReferenced DESC", ownerfp)
		if err != nil {
			return dbArr, err
		}
	case "(ownr)(lim-offs)":
		query, args, err = sqlx.In("SELECT * from Truststates WHERE (Owner = ?) ORDER BY LastReferenced DESC LIMIT ? OFFSET ?", ownerfp, limit, offset)
		if err != nil {
			return dbArr, err
		}
	default:
		logging.Logf(1, "The request you've made to ReadDbTruststates was invalid. Fps: %v, Start: %v, End: %v, Opts: %#v", fingerprints, beginTimestamp, endTimestamp, opts)
		return dbArr, fmt.Errorf("The request you've made to ReadDbTruststates was invalid. Fps: %v, Start: %v, End: %v, Opts: %#v", fingerprints, beginTimestamp, endTimestamp, opts)
	}

	rows, err := globals.DbInstance.Queryx(query, args...)
	if err != nil {
		return dbArr, err
	}
	defer rows.Close() // In case of premature exit.
	for rows.Next() {
		var entity DbTruststate
		err := rows.StructScan(&entity)
		if err != nil {
			return dbArr, err
		}
		dbArr = append(dbArr, entity)
	}
	rows.Close()
	return dbArr, nil
}

// The Reader functions that return DB instances, rather than API ones.

// ReadDBBoardOwners reads board owners from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

// This is left as a single-select, not multiple, because it already supports returning multiple entities, and there is no demand for these to be fetched in bulk.
func ReadDBBoardOwners(BoardFingerprint api.Fingerprint,
	KeyFingerprint api.Fingerprint) ([]DbBoardOwner, error) {
	var arr []DbBoardOwner
	// If this query is without a key fingerprint (we want all addresses with that board fingerprint), change the query as such.
	if KeyFingerprint == "" {
		rows, err := globals.DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ?", BoardFingerprint)
		if err != nil {
			logging.Log(1, err)
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var boardOwner DbBoardOwner
			err := rows.StructScan(&boardOwner)
			if err != nil {
				logging.Log(1, err)
			}
			arr = append(arr, boardOwner)
		}
		rows.Close()
	} else {
		rows, err := globals.DbInstance.Queryx("SELECT * from BoardOwners WHERE BoardFingerprint = ? AND KeyFingerprint = ?", BoardFingerprint, KeyFingerprint)
		if err != nil {
			logging.Log(1, err)
		}
		defer rows.Close() // In case of premature exit.
		for rows.Next() {
			var boardOwner DbBoardOwner
			err := rows.StructScan(&boardOwner)
			if err != nil {
				logging.Log(1, err)
			}
			arr = append(arr, boardOwner)
		}
		rows.Close()
	}
	return arr, nil
}

// ReadDBSubprotocols reads the subprotocols of a given address from the database. Even when there is a single result, it will still be arriving in an array to provide a consistent API.

func ReadDBSubprotocols(Location api.Location, Sublocation api.Location, Port uint16) ([]DbSubprotocol, error) {
	var fpArr []api.Fingerprint
	rows, err := globals.DbInstance.Queryx("SELECT * from AddressesSubprotocols WHERE AddressLocation = ? AND AddressSublocation = ? AND AddressPort = ?", Location, Sublocation, Port)
	if err != nil {
		logging.Log(1, err)
	}
	defer rows.Close() // In case of premature exit.
	// Get the Subprotocol fingerprints from the junction table.
	for rows.Next() {
		var dbAddressSubprot DbAddressSubprotocol
		err := rows.StructScan(&dbAddressSubprot)
		if err != nil {
			logging.Log(1, err)
		}
		fpArr = append(fpArr, dbAddressSubprot.SubprotocolFingerprint)
	}
	rows.Close()
	// For each fingerprint, get the matching subprotocol.
	var subprotArr []DbSubprotocol
	for _, val := range fpArr {
		rows, err := globals.DbInstance.Queryx("SELECT * from Subprotocols WHERE Fingerprint = ?", val)
		if err != nil {
			logging.Log(1, err)
		}
		defer rows.Close() // In case of premature exit.

		for rows.Next() {
			var subprot DbSubprotocol
			err := rows.StructScan(&subprot)
			if err != nil {
				logging.Log(1, err)
			}
			subprotArr = append(subprotArr, subprot)
		}
		rows.Close()
	}
	return subprotArr, nil
}

type reqtypeOpts struct {
	fingerprints   []api.Fingerprint
	beginTimestamp api.Timestamp
	endTimestamp   api.Timestamp
	tclass         int
	typ            int
	parentBoardFp  string
	parentThreadFp string
	parentPostFp   string
	targetFp       string
	domainFp       string
	noDescendants  bool
	ownerFp        string
	name           string
	limit          int
	offset         int
}

func reqtype(opts reqtypeOpts) string {
	// logging.Logf(1, "Start:%v, End: %v, Fingerprints: %v", beginTimestamp, endTimestamp, fingerprints)
	// Determine what kind of request we have.
	var rtype string
	if len(opts.name) > 0 {
		rtype = rtype + "(name)"
		return rtype // Name search does not accept any other arguments.
	}
	if len(opts.fingerprints) > 0 {
		rtype = rtype + "(fp)"
	}
	if opts.beginTimestamp > 0 || opts.endTimestamp > 0 {
		rtype = rtype + "(ts)"
	}
	if opts.tclass != -1 {
		rtype = rtype + "(tc)"
	}
	if opts.typ != -1 {
		rtype = rtype + "(ty)"
	}
	if len(opts.parentBoardFp) > 0 {
		rtype = rtype + "(pbfp)"
	}
	if len(opts.parentThreadFp) > 0 {
		rtype = rtype + "(ptfp)"
	}
	if len(opts.parentPostFp) > 0 {
		rtype = rtype + "(ppfp)"
	}
	if len(opts.targetFp) > 0 {
		rtype = rtype + "(tafp)"
	}
	if len(opts.domainFp) > 0 {
		rtype = rtype + "(dofp)"
	}
	if opts.noDescendants {
		rtype = rtype + "(nodesc)"
	}
	if len(opts.ownerFp) > 0 {
		// If owner is defined, it overrides - so that there can be no combination of anything prior with owner.
		// rtype = rtype + "(ownr)"
		rtype = "(ownr)"
	}
	if opts.limit > 0 || opts.offset > 0 {
		rtype = rtype + "(lim-offs)" // These are pairs.
	}
	return rtype
}

/*

var id int
err = db.Get(&id, "SELECT count(*) FROM place")

*/

func GetBoardThreadsCount(fp string) int {
	var count int
	err := globals.DbInstance.Get(&count, "SELECT COUNT(1) from Threads where Board= ?", fp)
	if err != nil {
		logging.Log(1, err)
	}
	return count
}

func GetThreadPostsCount(fp string) int {
	var count int
	err := globals.DbInstance.Get(&count, "SELECT COUNT(1) from Posts where Thread= ?", fp)
	if err != nil {
		logging.Log(1, err)
	}
	return count
}

// func Dbg_convertAddrSliceToNameSlice(nodes []DbAddress) []string {
// 	names := []string{}
// 	for _, val := range nodes {
// 		names = append(names, val.ClientName)
// 	}
// 	return names
// }

type dbCounts struct {
	Boards      int32
	Threads     int32
	Posts       int32
	Votes       int32
	Keys        int32
	Truststates int32
	Addresses   int32
}

func Dbg_ReadDatabaseCounts() dbCounts {
	b, t, p, v, k, ts, a := 0, 0, 0, 0, 0, 0, 0
	globals.DbInstance.Get(&b, "SELECT count(Fingerprint) FROM Boards")
	globals.DbInstance.Get(&t, "SELECT count(Fingerprint) FROM Threads")
	globals.DbInstance.Get(&p, "SELECT count(Fingerprint) FROM Posts")
	globals.DbInstance.Get(&v, "SELECT count(Fingerprint) FROM Votes")
	globals.DbInstance.Get(&k, "SELECT count(Fingerprint) FROM PublicKeys")
	globals.DbInstance.Get(&ts, "SELECT count(Fingerprint) FROM Truststates")
	globals.DbInstance.Get(&a, "SELECT count(Location) FROM Addresses")
	return dbCounts{int32(b), int32(t), int32(p), int32(v), int32(k), int32(ts), int32(a)}
}
