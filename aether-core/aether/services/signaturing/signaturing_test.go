package signaturing_test

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/signaturing"
	// "crypto/elliptic"
	"aether-core/aether/backend/cmd"
	"golang.org/x/crypto/ed25519"
	// "fmt"
	// "encoding/hex"
	"os"
	"testing"
)

// Infrastructure, setup and teardown

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func setup() {
	cmd.EstablishConfigs(nil)
}

func teardown() {
}

// Tests

func TestCreateKeyPair_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	if len(signaturing.MarshalPrivateKey(*privKey)) == 0 {
		t.Errorf("Created key pair is invalid.")
	}
}

func TestSign_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	signature, err2 := signaturing.Sign("this is my input", privKey)
	if err2 != nil || len(signature) == 0 {
		t.Errorf("Signing failed. Err: '%s'", err2)
	}
}

func TestVerify_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	signature, err2 := signaturing.Sign("this is my input", privKey)
	if err2 != nil {
		t.Errorf("Signing failed. Err: '%s'", err2)
	}
	marshaledPubKey := signaturing.MarshalPublicKey(privKey.Public().(ed25519.PublicKey))
	signatureIsValid := signaturing.Verify("this is my input", signature, string(marshaledPubKey))
	if signatureIsValid != true {
		t.Errorf("Signature verification failed.")
	}
}

func TestVerify_Fail(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	signature, err2 := signaturing.Sign("this is my input", privKey)
	if err2 != nil {
		t.Errorf("Signing failed. Err: '%s'", err2)
	}
	signatureIsValid := signaturing.Verify("this is my input", signature, "fake pub key")
	if signatureIsValid != false {
		t.Errorf("This signature was supposed to fail, but it did not.")
	}
}

func TestBoardCreateVerifySignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.ProofOfWork = "my fake pow"
	newboard2.EntityVersion = 1
	err2 := newboard2.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := signaturing.MarshalPublicKey(privKey.Public().(ed25519.PublicKey))
		result, err3 := newboard2.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestBoardCreateVerifySignatureWithPriorSignature_Success(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.Signature = "my fake signature"
	newboard2.EntityVersion = 1
	err2 := newboard2.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := signaturing.MarshalPublicKey(privKey.Public().(ed25519.PublicKey))
		result, err3 := newboard2.VerifySignature(marshaledPubKey)
		if err3 != nil {
			t.Errorf("Test failed, err: '%s'", err3)
		} else if result != true {
			t.Errorf("Test failed, this Signature should be valid but it is not.")
		}
	}
}

func TestBoardCreateVerifySignature_InvalidPubKey_Fail(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.ProofOfWork = "my fake pow"
	newboard2.EntityVersion = 1
	err2 := newboard2.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		signature := string(newboard2.Signature)
		signatureIsValid := signaturing.Verify("this is my input", signature, "fake pub key")
		if signatureIsValid != false {
			t.Errorf("This signature was supposed to fail, but it did not.")
		}
	}
}

func TestBoardCreateVerifySignature_InvalidSignature_Fail(t *testing.T) {
	privKey, err := signaturing.CreateKeyPair()
	if err != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err)
	}
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	newboard2.ProofOfWork = "my fake pow"
	newboard2.EntityVersion = 1
	err2 := newboard2.CreateSignature(privKey)
	if err2 != nil {
		t.Errorf("Test failed, err: '%s'", err2)
	} else {
		marshaledPubKey := signaturing.MarshalPublicKey(privKey.Public().(ed25519.PublicKey))
		signatureIsValid := signaturing.Verify("this is my input", "fake signature", marshaledPubKey)
		if signatureIsValid != false {
			t.Errorf("This signature was supposed to fail, but it did not.")
		}
	}
}

func TestBoardCreateUpdateSignature_Success(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	newboard.EntityVersion = 1
	err := newboard.CreateSignature(privKey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdateSignature(privKey)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey := signaturing.MarshalPublicKey(privKey.Public().(ed25519.PublicKey))
			result, err3 := newboard.VerifySignature(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}
