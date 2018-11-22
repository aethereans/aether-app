package api_test

import (
	"aether-core/aether/backend/cmd"
	"aether-core/aether/io/api"
	// "aether-core/aether/services/configstore"
	"aether-core/aether/services/create"
	"aether-core/aether/services/globals"
	// "aether-core/aether/services/logging"
	"aether-core/aether/services/signaturing"
	"crypto/elliptic"
	"encoding/hex"
	// "fmt"
	"os"
	"strings"
	"testing"
)

// Infrastructure, setup and teardown

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

var MarshaledPubKey string

func setup() {
	cmd.EstablishConfigs(nil)
	globals.BackendConfig.SetMinimumPoWStrengths(16)
	MarshaledPubKey = hex.EncodeToString(elliptic.Marshal(elliptic.P521(), globals.FrontendConfig.GetUserKeyPair().PublicKey.X, globals.FrontendConfig.GetUserKeyPair().PublicKey.Y))
}

func teardown() {
}

// Tests

func TestVerify_Success(t *testing.T) {
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"keyfp", MarshaledPubKey, "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	err2 := api.Verify(&thr)
	if err2 != nil {
		t.Errorf("Object verification process failed. Error: '%#v\n'", err2)
	}
	if thr.GetVerified() != true {
		t.Errorf("This object should be valid, but it is invalid. Object: '%#v\n'", thr)
	}
}

func TestVerify_BrokenFingerprint_Fail(t *testing.T) {
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint", MarshaledPubKey, "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	thr.Board = "I'm changing the board of this thread to fail the test."
	err2 := api.Verify(&thr)
	errMessage := "Fingerprint of this entity is invalid"
	if err2 == nil || thr.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_BrokenPoW1_Fail(t *testing.T) {
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint", MarshaledPubKey, "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	thr.ProofOfWork = "I'm changing the board pow to fail the test."
	// Re-shrink-wrap
	thr.CreateFingerprint()
	errMessage := "PoW had more or less fields than expected"
	err2 := api.Verify(&thr)
	if err2 == nil || thr.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
	// fmt.Printf("%#v\n", thr)
	// fmt.Printf("%#v\n", thr.GetVerified())
}

func TestVerify_BrokenPoW2_Fail(t *testing.T) {
	// Changing a mutable element, but not actually running update.
	entity, err :=
		create.CreateBoard(
			"board name",
			"my board owner fingerprint", MarshaledPubKey,
			*new([]api.BoardOwner),
			"my board description", "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	entity.Description = "Changing an element that is not protected by fingerprint"
	err2 := api.Verify(&entity)
	errMessage := "This proof of work is invalid or malformed"
	if err2 == nil || entity.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_BrokenSignature_Fail(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	thr, err :=
		create.CreateThread(
			"my board fingerprint",
			"my thread name",
			"my thread body",
			"my thread link",
			"my owner fingerprint", MarshaledPubKey, "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%s'", err)
	}
	// Re-shrink-wrap
	thr.CreateSignature(privKey) // Signing it with a new key
	thr.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), 20)
	thr.CreateFingerprint()
	errMessage := "This signature is invalid, but no reason given as to why"
	err2 := api.Verify(&thr)
	if err2 == nil || thr.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
	// fmt.Printf("%#v\n", thr)
	// fmt.Printf("%#v\n", thr.GetVerified())
}

// Test update verify success / fail

func TestVerify_UpdatedItemSuccess(t *testing.T) {
	board, err :=
		create.CreateBoard(
			"my board name",
			"random key fp", MarshaledPubKey,
			[]api.BoardOwner{},
			"my board description", "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	updatereq := create.BoardUpdateRequest{}
	updatereq.Entity = &board
	updatereq.DescriptionUpdated = true
	updatereq.NewDescription = "I changed the board description!"
	create.UpdateBoard(updatereq)
	err2 := api.Verify(&board)
	if err2 != nil {
		t.Errorf("Object verification process failed. Error: '%#v\n'", err2)
	}
	if board.GetVerified() != true {
		t.Errorf("This object should be valid, but it is invalid. Object: '%#v\n'", board)
	}
}

func TestVerify_UpdatedItemFailure_Pow(t *testing.T) {
	// Failed to call the update request.
	board, err :=
		create.CreateBoard(
			"my board name",
			"random key ", MarshaledPubKey,
			[]api.BoardOwner{},
			"my board description", "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is a mutable field, the fingerprint does not protect against it. The field to fail first will be pow.
	board.Description = "new description"
	err2 := api.Verify(&board)
	errMessage := "This proof of work is invalid or malformed"
	if err2 == nil || board.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

func TestVerify_UpdatedItemFailure_Fingerprint(t *testing.T) {
	// Failed to call the update request.
	board, err :=
		create.CreateBoard(
			"my board name",
			"randomkeyfp", MarshaledPubKey,
			[]api.BoardOwner{},
			"my board description", "", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is an immutable field, it will trip up the fingerprint first.
	board.Name = "New name"
	err2 := api.Verify(&board)
	errMessage := "Fingerprint of this entity is invalid"
	if err2 == nil || board.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

// The test below no longer applies, because we are no longer checking the key object at verification. It is not being checked against the key embedded in the object itself.

// func TestVerify_UpdatedItemFailure_Signature(t *testing.T) {
// 	// Failed to call the update request.
// 	keyEntity, err3 := create.CreateKey(
// 		"", MarshaledPubKey, "", "", "", "")
// 	if err3 != nil {
// 		t.Errorf("Object creation failed. Err: '%s'", err3)
// 	}
// 	board, err :=
// 		create.CreateBoard(
// 			"my board name",
// 			keyEntity.Fingerprint, keyEntity.Key,
// 			[]api.BoardOwner{},
// 			"my board description", "", "")
// 	if err != nil {
// 		t.Errorf("Object creation failed. Err: '%#v\n'", err)
// 	}
// 	// Since this is an immutable field, it will trip up the fingerprint first.
// 	keyEntity.Fingerprint = "changed key entity fingerprint"
// 	err2 := api.Verify(&board)
// 	errMessage := "A wrong key is provided for this signature"
// 	if err2 == nil || board.GetVerified() == true {
// 		t.Errorf("Expected an error to be raised from this test.")
// 	} else if !strings.Contains(err2.Error(), errMessage) {
// 		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
// 	}
// }

// This checks whether mutable item is allowed change. (but POW will fail due to not reminting the pow and signature)
func TestVerify_UpdatedItemFailure_Mutable_Fingerprint(t *testing.T) {
	// Failed to call the update request.
	board, err :=
		create.CreateBoard(
			"my board name",
			"my key fingerprint", MarshaledPubKey,
			[]api.BoardOwner{},
			"my board description", "my old meta field", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is an immutable field, it will trip up the fingerprint first.
	board.Meta = "my new meta field"
	err2 := api.Verify(&board)
	errMessage := "This proof of work is invalid or malformed"
	if err2 == nil || board.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

// This checks whether immutable item is guarded appropriately against change.
func TestVerify_UpdatedItemFailure_Immutable_Fingerprint(t *testing.T) {
	// Failed to call the update request.
	board, err :=
		create.CreateBoard(
			"my board name",
			"my board owner fingerpint", MarshaledPubKey,
			[]api.BoardOwner{},
			"my board description", "my old meta field", "")
	if err != nil {
		t.Errorf("Object creation failed. Err: '%#v\n'", err)
	}
	// Since this is an immutable field, it will trip up the fingerprint first.
	board.Owner = "changed my board owner fingerpint"
	err2 := api.Verify(&board)
	errMessage := "Fingerprint of this entity is invalid"
	if err2 == nil || board.GetVerified() == true {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}

// We can't really test the verify improper version (below) because it's behind the fp / sig / pow enabled gate. if you disable that gate, all of them just return true, regardless of version. if you disable it, the only way I can mock  a different version object would be to allow the versions allowed to be dynamic, which is not necessary.

// func TestVerify_ImproperVersion_Fail(t *testing.T) {
// 	var p api.Post
// 	p.SetVerified(true)
// 	p.Fingerprint = "yo"
// 	p.Body = "body"
// 	p.EntityVersion = 99999999999999 // we don't support this version
// 	p.Owner = "owner"
// 	p.OwnerPublicKey = "ownerpk"
// 	p.Creation = 1
// 	p.Signature = "sig"
// 	p.ProofOfWork = "pow"
// 	p.LastUpdate = 2
// 	p.UpdateProofOfWork = "updatepow"
// 	p.UpdateSignature = "updatesig"
// 	p.Board = "boardpk"
// 	p.Thread = "threadpk"
// 	p.Parent = "yo2"
// 	// We're just testing for version gating right now.
// 	globals.BackendTransientConfig.FingerprintCheckEnabled = false
// 	globals.BackendTransientConfig.SignatureCheckEnabled = false
// 	globals.BackendTransientConfig.ProofOfWorkCheckEnabled = false
// 	err2 := api.Verify(&p)
// 	// Set it back up for the rest of the tests.
// 	globals.BackendTransientConfig.FingerprintCheckEnabled = true
// 	globals.BackendTransientConfig.SignatureCheckEnabled = true
// 	globals.BackendTransientConfig.ProofOfWorkCheckEnabled = true
// 	errMessage := "Fingerprint of this entity is invalid"
// 	if err2 == nil || p.GetVerified() == true {
// 		t.Errorf("Expected an error to be raised from this test.")
// 	} else if !strings.Contains(err2.Error(), errMessage) {
// 		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
// 	}
// }

func TestCreate_ImproperVersion_Fail(t *testing.T) {
	// generate post manually so we can attempt to set the version to what we want.
	var p api.Post
	p.SetVerified(true)
	p.Fingerprint = "yo"
	p.Body = "body"
	p.EntityVersion = 99999999999999 // we don't support this version
	p.Owner = "owner"
	p.OwnerPublicKey = "ownerpk"
	p.Creation = 1
	p.Signature = "sig"
	p.ProofOfWork = "pow"
	p.LastUpdate = 2
	p.UpdateProofOfWork = "updatepow"
	p.UpdateSignature = "updatesig"
	p.Board = "boardpk"
	p.Thread = "threadpk"
	p.Parent = "yo2"
	err2 := create.Bake(&p)
	errMessage := "Signature creation of this version of this entity is not supported in this version of the app"
	if err2 == nil {
		t.Errorf("Expected an error to be raised from this test.")
	} else if !strings.Contains(err2.Error(), errMessage) {
		t.Errorf("Test returned an error that did not include the expected one. Error: '%s', Expected error: '%s'", err2, errMessage)
	}
}
