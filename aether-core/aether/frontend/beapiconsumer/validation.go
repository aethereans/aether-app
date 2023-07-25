package beapiconsumer

import (
	pbstructs "aether-core/aether/protos/mimapi"
)

func validateBoards(eSet []*pbstructs.Board) []*pbstructs.Board {
	var valids []*pbstructs.Board

	for k := range eSet {
		if boardValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

func validateThreads(eSet []*pbstructs.Thread) []*pbstructs.Thread {
	var valids []*pbstructs.Thread

	for k := range eSet {
		if threadValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

func validatePosts(eSet []*pbstructs.Post) []*pbstructs.Post {
	var valids []*pbstructs.Post

	for k := range eSet {
		if postValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

func validateVotes(eSet []*pbstructs.Vote) []*pbstructs.Vote {
	var valids []*pbstructs.Vote

	for k := range eSet {
		if voteValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

func validateKeys(eSet []*pbstructs.Key) []*pbstructs.Key {
	var valids []*pbstructs.Key

	for k := range eSet {
		if keyValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

func validateTruststates(eSet []*pbstructs.Truststate) []*pbstructs.Truststate {
	var valids []*pbstructs.Truststate

	for k := range eSet {
		if truststateValid(eSet[k]) {
			valids = append(valids, eSet[k])
		}
	}
	return valids
}

/*
TODO

We need to implement validity checks. What we need is pretty much that we need to change this to an apistruct, run verify on that apistruct, and send the result back.

This is useful when you don't trust the backend, though, in the case you're using a public backend. There is no need for this if the backend is running on your local machine, where the process is owned by you. Doing this kind of verification in the frontend *and* the backend means doing 2x the work, so these should only be enabled when the frontend knows it's connecting to a remote, untrusted backend.

So I'm leaving the frontend verification unimplemented for now because we're not yet at a place where there are people hosting backends for common use. Whenever that time comes, this verification system should be activated.

(Heads up - the verifiers in apistructs assume a backend config being present, in the frontend that is not the case. You'll need to untangle that and make them take a backend or a frontend config.)
*/

func boardValid(e *pbstructs.Board) bool {
	return true
}

func threadValid(e *pbstructs.Thread) bool {
	return true
}

func postValid(e *pbstructs.Post) bool {
	return true
}

func voteValid(e *pbstructs.Vote) bool {
	return true
}

func keyValid(e *pbstructs.Key) bool {
	return true
}

func truststateValid(e *pbstructs.Truststate) bool {
	return true
}
