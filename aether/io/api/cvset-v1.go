// API > Create / Verify / EntitySet V1
// This file provides the version specific create and verify methods based on entity versions. This file is for the v1 versions of the objects.

package api

import (
	// "fmt"
	"aether-core/aether/services/ca"
	"aether-core/aether/services/fingerprinting"
	"aether-core/aether/services/globals"
	// "aether-core/aether/services/logging"
	"aether-core/aether/services/proofofwork"
	"aether-core/aether/services/signaturing"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ed25519"
	// "github.com/davecgh/go-spew/spew"
)

/*
This file contains entity-version-specific creation and verification flows for:
		Fingerprinting
		Proof of work
		Signaturing

that pertain to entities:
		Board
		Thread
		Post
		Vote
		Key
		Truststate

for versions:
		v1
*/

// PoW

// // CreatePoW

func createBoardPoW_V1(b *Board, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *b
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	b.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createThreadPoW_V1(t *Thread, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *t
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	t.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createPostPoW_V1(p *Post, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *p
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	p.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createVotePoW_V1(v *Vote, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *v
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	v.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createKeyPoW_V1(k *Key, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *k
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	k.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createTruststatePoW_V1(ts *Truststate, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *ts
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	ts.ProofOfWork = ProofOfWork(pow)
	return nil
}

func createApiResponsePoW_V1(ar *ApiResponse, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *ar
	// Remove the existing proof of work if any exists so as to not accidentally take it as an input to the new proof of work about to be calculated.
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	ar.ProofOfWork = ProofOfWork(pow)
	// logging.Logf(1, "We created an ApiResponse PoW. PoW: %v", ar.ProofOfWork)
	return nil
}

// // Create UpdatePoW

func createBoardUpdatePoW_V1(b *Board, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *b
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	b.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func createThreadUpdatePoW_V1(t *Thread, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *t
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	t.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func createPostUpdatePoW_V1(p *Post, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *p
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	p.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func createVoteUpdatePoW_V1(v *Vote, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *v
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	v.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func createKeyUpdatePoW_V1(k *Key, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *k
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	k.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

func createTruststateUpdatePoW_V1(ts *Truststate, keyPair *ed25519.PrivateKey, difficulty int) error {
	cpI := *ts
	// Updateable
	cpI.UpdateProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create PoW
	pow, err := proofofwork.Create(string(res), difficulty, keyPair)
	if err != nil {
		return err
	}
	ts.UpdateProofOfWork = ProofOfWork(pow)
	return nil
}

// // Verify PoW

func verifyBoardPoW_V1(b *Board, pubKey string) (bool, error) {
	cpI := *b
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().BoardUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().BoardUpdate
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Board
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Board
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func verifyThreadPoW_V1(t *Thread, pubKey string) (bool, error) {
	cpI := *t
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().ThreadUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().ThreadUpdate
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Thread
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Thread
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func verifyPostPoW_V1(p *Post, pubKey string) (bool, error) {
	cpI := *p
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().PostUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().PostUpdate
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Post
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Post
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func verifyVotePoW_V1(v *Vote, pubKey string) (bool, error) {
	cpI := *v
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().VoteUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().VoteUpdate
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Vote
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Vote
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func verifyKeyPoW_V1(k *Key, pubKey string) (bool, error) {
	cpI := *k
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().KeyUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().KeyUpdate
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Key
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Key
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

// Special case below: we drop the PoW requirement to a minimum if it's a CA-specific TypeClassed Truststate, and we trust that CA.
func verifyTruststatePoW_V1(ts *Truststate, pubKey string) (bool, error) {
	cpI := *ts
	var pow string
	var neededStrength int
	// Determine if we are checking for original or update PoW
	if len(cpI.UpdateProofOfWork) > 0 {
		// This is a VerifyUpdatePoW. (The object was subject to some updates.)
		// Updateable
		// Save PoW to be verified
		pow = string(cpI.UpdateProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().TruststateUpdate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().TruststateUpdate
		}
		// The process below allows for trusted CAs to be able to issue entities with lower PoW. Since this acceptance pass is done at the backend, the frontend verification does not need to care about this.
		if !isFrontend() && (ts.TypeClass == 2 || ts.TypeClass == 3) &&
			ca.IsTrustedCAKeyByPK(pubKey) {
			// This truststate is using a CA-specific TypeClass.
			// If it is a CA that we trust, we drop the truststate PoW requirement to minimum.
			neededStrength = globals.BackendTransientConfig.MinimumTrustedPoWStrength
			// If not, we retain the original TS minimum. We'll fail this for not having the entitlement for the CA-specific TypeClass in the future, but not failing it here just to not create confusing error messages.
		}
		// Delete PoW so that the PoW will match
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifyPoW (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save PoW to be verified
		pow = string(cpI.ProofOfWork)
		if isFrontend() {
			neededStrength = globals.FrontendConfig.GetMinimumPoWStrengths().Truststate
		} else {
			neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().Truststate
		}
		if !isFrontend() && (ts.TypeClass == 2 || ts.TypeClass == 3) &&
			ca.IsTrustedCAKeyByPK(pubKey) {
			neededStrength = globals.BackendTransientConfig.MinimumTrustedPoWStrength
		}
		// Delete PoW so that the PoW will match
		cpI.ProofOfWork = ""
	}

	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

func verifyApiResponsePoW_V1(ar *ApiResponse, pubKey string) (bool, error) {
	cpI := *ar
	var pow string
	var neededStrength int
	// Save PoW to be verified
	pow = string(cpI.ProofOfWork)
	neededStrength = globals.BackendConfig.GetMinimumPoWStrengths().ApiResponse
	// Delete PoW so that the PoW will match
	cpI.ProofOfWork = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify PoW
	verifyResult, strength, err := proofofwork.Verify(string(res), pow, pubKey)
	if err != nil {
		return false, err
	}
	// logging.Logf(1, "We verified an ApiResponse PoW. PoW: %v Result: %v", pow, verifyResult)
	// If the PoW is valid
	if verifyResult {
		// Check if satisfies required minimum
		if strength >= neededStrength {
			return true, nil
		} else {
			return false, errors.New(fmt.Sprint(
				"This proof of work is not strong enough. PoW: ", pow))
		}
	} else {
		return false, errors.New(fmt.Sprint(
			"This proof of work is invalid, but no reason given as to why. PoW: ", pow))
	}
}

// Fingerprint

// Create Fp

func createBoardFp_V1(b *Board) {
	cpI := *b
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.BoardOwners = []BoardOwner{}
	cpI.Description = ""
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	b.Fingerprint = Fingerprint(fp)
}

func createThreadFp_V1(t *Thread) {
	cpI := *t
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Body = ""
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	t.Fingerprint = Fingerprint(fp)
}

func createPostFp_V1(p *Post) {
	cpI := *p
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Body = ""
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	p.Fingerprint = Fingerprint(fp)
}

func createVoteFp_V1(v *Vote) {
	cpI := *v
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Type = 0
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	v.Fingerprint = Fingerprint(fp)
}

func createKeyFp_V1(k *Key) {
	cpI := *k
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Info = ""
	cpI.Expiry = 0
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	k.Fingerprint = Fingerprint(fp)
}

func createTruststateFp_V1(ts *Truststate) {
	cpI := *ts
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Type = 0
	cpI.Expiry = 0
	cpI.Meta = ""
	// Remove the existing fingerprint if any exists so as to not accidentally take it as an input to the new fingerprint about to be calculated.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create Fingerprint
	fp := fingerprinting.Create(string(res))
	ts.Fingerprint = Fingerprint(fp)
}

// // Verify Fp

func verifyBoardFingerprint_V1(b *Board) bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *b
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.BoardOwners = []BoardOwner{}
	cpI.Description = ""
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func verifyThreadFingerprint_V1(t *Thread) bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *t
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Body = ""
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func verifyPostFingerprint_V1(p *Post) bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *p
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Body = ""
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func verifyVoteFingerprint_V1(v *Vote) bool {
	if !isFrontend() && globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *v
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Type = 0
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func verifyKeyFingerprint_V1(k *Key) bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *k
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Info = ""
	cpI.Expiry = 0
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

func verifyTruststateFingerprint_V1(ts *Truststate) bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	cpI := *ts
	var fp string
	fp = string(cpI.Fingerprint)
	// Updateable set
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Remove ALL mutable fields
	cpI.Type = 0
	cpI.Expiry = 0
	cpI.Meta = ""
	// Remove the existing fingerprint so that it won't be included as part of the input to be verified.
	cpI.Fingerprint = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Fingerprint
	verifyResult := fingerprinting.Verify(string(res), fp)
	return verifyResult
}

// Signaturing

// // Create Signature

func createBoardSignature_V1(b *Board, keyPair *ed25519.PrivateKey) error {
	cpI := *b
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	b.Signature = Signature(signature)
	return nil
}

func createThreadSignature_V1(t *Thread, keyPair *ed25519.PrivateKey) error {
	cpI := *t
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	t.Signature = Signature(signature)
	return nil
}

func createPostSignature_V1(p *Post, keyPair *ed25519.PrivateKey) error {
	cpI := *p
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	p.Signature = Signature(signature)
	return nil
}

func createVoteSignature_V1(v *Vote, keyPair *ed25519.PrivateKey) error {
	cpI := *v
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	v.Signature = Signature(signature)
	return nil
}

func createKeySignature_V1(k *Key, keyPair *ed25519.PrivateKey) error {
	cpI := *k
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	k.Signature = Signature(signature)
	return nil
}

func createTruststateSignature_V1(ts *Truststate, keyPair *ed25519.PrivateKey) error {
	cpI := *ts
	// Updateable
	cpI.Fingerprint = ""
	cpI.LastUpdate = 0
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	cpI.ProofOfWork = ""
	// Remove existing signature if any so it won't end up in the mix accidentally.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ts.Signature = Signature(signature)
	return nil
}

// // Create UpdateSignature

func createBoardUpdateSignature_V1(b *Board, keyPair *ed25519.PrivateKey) error {
	cpI := *b
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	b.UpdateSignature = Signature(signature)
	return nil
}

func createThreadUpdateSignature_V1(t *Thread, keyPair *ed25519.PrivateKey) error {
	cpI := *t
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	t.UpdateSignature = Signature(signature)
	return nil
}

func createPostUpdateSignature_V1(p *Post, keyPair *ed25519.PrivateKey) error {
	cpI := *p
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	p.UpdateSignature = Signature(signature)
	return nil
}

func createVoteUpdateSignature_V1(v *Vote, keyPair *ed25519.PrivateKey) error {
	cpI := *v
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	v.UpdateSignature = Signature(signature)
	return nil
}

func createKeyUpdateSignature_V1(k *Key, keyPair *ed25519.PrivateKey) error {
	cpI := *k
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	k.UpdateSignature = Signature(signature)
	return nil
}

func createTruststateUpdateSignature_V1(ts *Truststate, keyPair *ed25519.PrivateKey) error {
	cpI := *ts
	// Updateable
	cpI.UpdateProofOfWork = ""
	cpI.UpdateSignature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ts.UpdateSignature = Signature(signature)
	return nil
}

// // Verify Signature

func verifyBoardSignature_V1(b *Board, pubKey string) (bool, error) {
	cpI := *b
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func verifyThreadSignature_V1(t *Thread, pubKey string) (bool, error) {
	cpI := *t
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func verifyPostSignature_V1(p *Post, pubKey string) (bool, error) {
	cpI := *p
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func verifyVoteSignature_V1(v *Vote, pubKey string) (bool, error) {
	cpI := *v
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func verifyKeySignature_V1(k *Key, pubKey string) (bool, error) {
	cpI := *k
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}

func verifyTruststateSignature_V1(ts *Truststate, pubKey string) (bool, error) {
	cpI := *ts
	var signature string
	// Determine if we are checking for original or update signature
	if len(cpI.UpdateSignature) > 0 {
		// This is a VerifyUpdateSignature. (The object was subject to some updates.)
		// Updateable
		// Save Signature to be verified
		signature = string(cpI.UpdateSignature)
		// Delete Signature and PoW so that the Signature check will match
		cpI.UpdateSignature = ""
		cpI.UpdateProofOfWork = ""
	} else {
		// This is a VerifySignature (there is no update on this object.)
		// Updateable
		cpI.Fingerprint = ""
		cpI.LastUpdate = 0
		cpI.UpdateProofOfWork = ""
		cpI.UpdateSignature = ""
		// Save signature to be verified
		signature = string(cpI.Signature)
		// This happens *after* Signature, so should be empty here.
		cpI.ProofOfWork = ""
		cpI.Signature = ""
	}
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, pubKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(
			"This signature is invalid, but no reason given as to why. Signature: ", signature))
	}
}
