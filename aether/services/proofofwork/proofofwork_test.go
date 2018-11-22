package proofofwork_test

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/signaturing"
	// "fmt"
	// "log"
	"crypto/elliptic"
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"os"
	"strings"
	"testing"
)

// Infrastructure, setup and teardown

// Initialise the test object.

var newboard api.Board
var signedNewboard api.Board
var signedNewboardUpdated api.Board
var signedNewBoardUpdatedPubkey string

var invalidPoWBoard api.Board
var weakPoWBoard api.Board

var userKey api.Key

var fakeSignedBoard api.Board

var minPowStrength int64

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func setup() {

	becfg, err := configstore.EstablishBackendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	becfg.Cycle()
	globals.BackendConfig = becfg

	fecfg, err := configstore.EstablishFrontendConfig()
	if err != nil {
		logging.LogCrash(err)
	}
	fecfg.Cycle()
	globals.FrontendConfig = fecfg

	// Set up the min PoW strengths from services. This is normally in main()
	globals.BackendConfig.SetMinimumPoWStrengths(16)

	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	newboard.ProofOfWork = ""
	newboard.ProofOfWork = "MIM1:20::::hvMazmkOQUvYriEB:630538:"
	// To regenerate:
	// newboard.CreatePoW(new(ed25519.PrivateKey), 20)
	// fmt.Println(newboard.ProofOfWork)

	signedNewboard.Fingerprint = "my random fingerprint3"
	signedNewboard.Creation = 4564654
	signedNewboard.Name = "my board name"
	signedNewboard.Description = "my board description"
	signedNewboard.ProofOfWork = "MIM1:20::::cJSTqJNnTcTcYgzH:2947609:3338744f899d411e399e35b7f23a48a43790142ff8533e13f5383991700788108c32564e773bc40e484ac915dbf059edc97e605b66587c4c4e70143922a96f3de0-0133ea996835ece4dd192165bd0a60f78872bf7ce84d8b0481cdb772e711dde0a69426cab70ce50d4385b5637068322167e1fcce77159f14d37794d31f91d1a4e59f"

	// Marshaled pub key for this is:
	//0400fa3aa273d67a3161069414879677f0ecf554a13edaa7d85387fe8fafaacea4beb1daf0284efa5f6ac0bc5181f3d20f748d53f79f3c1979cc58b7210dd7a161b0a4007bfb1315a6e78cd43855a4a40401fdc977a01c9f53616b52c92b16dc42fab04b84cc1788fbebc4f1449634cffe35e457680c432f28232d49aa8fc73fde8f374525

	// To regenerate:
	// privKey, _ := signaturing.CreateKeyPair()
	// signedNewboard.CreatePoW(privKey, 20)
	// marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
	// fmt.Println(signedNewboard.ProofOfWork)
	// fmt.Println(marshaledPubKey)

	// privKey, _ := signaturing.CreateKeyPair()
	// fmt.Println(hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y)))
	signedNewBoardUpdatedPubkey = "04011c1b73221ac0afbd404ab0aa86ed7ea99c2042b2f05581bccd6b0e321edeab2eb5a56cbf2b0a952aed53cc18c47b2511552d1613eb710ea05b0b2646590e9b96ea01317c365886e4ec7302ad0ca12a667676c6fe2aff3afeeaa5439823976bfb2e9b60e787bceba8a2b529adb7a1057d60b054b62e333f8f84c37a9272ae3372e63952"
	signedNewboardUpdated.Fingerprint = "my random fingerprint"
	signedNewboardUpdated.Creation = 4564654
	signedNewboardUpdated.Name = "my board name"
	signedNewboardUpdated.Description = "description"
	// signedNewboardUpdated.CreatePoW(privKey, 20)
	// fmt.Println(signedNewboardUpdated.ProofOfWork)
	signedNewboardUpdated.Description = "I updated this board's description"
	// signedNewboardUpdated.CreateUpdatePoW(privKey, 20)
	// fmt.Println(signedNewboardUpdated.UpdateProofOfWork)
	signedNewboardUpdated.ProofOfWork = "MIM1:20::::nEGYzEPiiPlurnCt:630323:37767a77da2d0d6cb33042ad1e5841dd1acbac99ae13375555d5ffded151734d8493e39027259c2c0d5edc01f628f31b2a2e2ebf86952638992856b07c51fafa76-01ea31d15880cfc7c5d2519878f12c94fb4c8d7deebec6a06549a93d124ac754fc7087e4d1fbcb8a9d1bf2407fbf87b28d0c3aa0c7f4aea4319328866c090b9453d5"
	signedNewboardUpdated.UpdateProofOfWork = "MIM1:20::::DQXchTAxFPlMMAXl:745645:a59a5df7c787b5742242de80b18850ec0dd5130ee9bd4d2beb3ae8d4e3bd8ced1342cba34767abc41df791a37d5b18a41ffaddc7d0929c398745703b49436ea581-c7f29291c813157c339560f1d349bf5d2d62d968f98950a0ff61db4ff6927120c03c0b2f80539b8b028b09d84b33b5a9df0c9b73d3397ad0108f871cfb2bd35dbd"

	invalidPoWBoard.Fingerprint = "my random fingerprint"
	invalidPoWBoard.Creation = 4564654
	invalidPoWBoard.Name = "my board name"
	invalidPoWBoard.Description = "my board description"
	invalidPoWBoard.ProofOfWork = "MIM1:21::::QkaMjkJbvXInQLtW:1166891:"

	weakPoWBoard.Fingerprint = "my random fingerprint"
	weakPoWBoard.Creation = 4564654
	weakPoWBoard.Name = "my board name"
	weakPoWBoard.Description = "my board description"
	weakPoWBoard.ProofOfWork = "MIM1:18::::ZVqtpNqdZXkNcSpc:290736:"

	fakeSignedBoard.Fingerprint = "my random fingerprint"
	fakeSignedBoard.Creation = 4564654
	fakeSignedBoard.Name = "my board name"
	fakeSignedBoard.Description = "my board description"
	fakeSignedBoard.ProofOfWork = "MIM1:20::::xDQPQMOBXYIMCDvE:1912024:fake key"

}

func teardown() {
}

func ValidateTest(expected interface{}, actual interface{}, t *testing.T) {
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

// Tests

// VerifyPoW tests come first because they are used in the CreatePoW tests as verification.

// // VerifyPoW tests

func TestVerifyPoW_Success_WithoutKey(t *testing.T) {
	result, err := newboard.VerifyPoW("")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if result != true {
		t.Errorf("Test failed, this PoW should be valid but it is not.")
	}
}

func TestVerifyPoW_Success_WithKey(t *testing.T) {
	marshaledPubKey := "0400fa3aa273d67a3161069414879677f0ecf554a13edaa7d85387fe8fafaacea4beb1daf0284efa5f6ac0bc5181f3d20f748d53f79f3c1979cc58b7210dd7a161b0a4007bfb1315a6e78cd43855a4a40401fdc977a01c9f53616b52c92b16dc42fab04b84cc1788fbebc4f1449634cffe35e457680c432f28232d49aa8fc73fde8f374525"
	// fmt.Printf("%#v\n", signedNewboard)
	result, err := signedNewboard.VerifyPoW(marshaledPubKey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if result != true {
		t.Errorf("Test failed, this PoW should be valid but it is not.")
	}
}

func TestVerifyPoW_Fail_SignatureKeyMismatch_WithKey(t *testing.T) {
	_, err := signedNewboard.VerifyPoW("fake key")
	errMessage := "The signature of this PoW is invalid. The PoW signature and the public key provided does not match."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_Malformed(t *testing.T) {
	_, err := invalidPoWBoard.VerifyPoW("")
	errMessage := "This proof of work is invalid or malformed."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_NotStrongEnough(t *testing.T) {
	// For the test, ncrease min pow strength to something above the board.
	globals.BackendConfig.SetMinimumPoWStrengths(30)
	_, err := weakPoWBoard.VerifyPoW("")
	globals.BackendConfig.SetMinimumPoWStrengths(16)
	errMessage := "This proof of work is not strong enough."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_KeyProvidedButPoWIsUnsigned(t *testing.T) {
	_, err := newboard.VerifyPoW("fake key")
	errMessage := "A key is provided, but the PoW is unsigned."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_PoWIsSignedButKeyIsNotProvided(t *testing.T) {
	_, err := signedNewboard.VerifyPoW("")
	errMessage := "The PoW is signed, but a key is not provided."
	if err == nil {
		t.Errorf("Did not return error on missing key.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_UnrecognisedPoWVersion(t *testing.T) {
	// Create a misversioned PoW and attach it. If the version is wrong, we bail without any sort of processing other than parsing.
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM2:21::::QkaMjkJbvXInQLtW:1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This proof of work is in a format Mim does not support."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_ParsedDifficultyIsNegative(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:-20::::QkaMjkJbvXInQLtW:1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This proof of work is invalid or malformed. (Negative parsed difficulty)"
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_ParsedCounterIsNegative(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtW:-1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This proof of work is invalid or malformed. (Negative counter)"
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_ParsingFailure(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtW:a1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "PoW parsing failed, this PoW is invalid."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_DataInFieldsThatShouldBeEmpty(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20:AA:AA:AA:QkaMjkJbvXInQLtW:a1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This proof of work either has fields that should be empty and is not, or it does have empty fields which it should not."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_NoDataInFieldsThatShouldBeNotEmpty(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1::::::a1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This proof of work either has fields that should be empty and is not, or it does have empty fields which it should not."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_PoWTooLong(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtWAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA:1166891:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "This PoW is longer than maximum allowed 1024 characters."
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_TooManyFields(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtW:1166891::A:A:A"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "PoW had more or less fields than expected"
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

func TestVerifyPoW_Fail_TooFeWFields(t *testing.T) {
	var brokenVersionPoWBoard api.Board
	brokenVersionPoWBoard.Fingerprint = "my random fingerprint"
	brokenVersionPoWBoard.Creation = 4564654
	brokenVersionPoWBoard.Name = "my board name"
	brokenVersionPoWBoard.Description = "my board description"
	brokenVersionPoWBoard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtW:"
	_, err := brokenVersionPoWBoard.VerifyPoW("")
	errMessage := "PoW had more or less fields than expected"
	if err == nil {
		t.Errorf("Did not return error on malformed PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

// // VerifyUpdatePow tests

// Not testing the failure cases here because they all take root from proofofwork.Create, and all the failure cases in that is tested by non-update create function tests.

func TestVerifyUpdatePoW_Success_WithoutKey(t *testing.T) {
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my description"
	// newboard.CreatePoW(new(ed25519.PrivateKey), 20)
	// fmt.Println(newboard.ProofOfWork)
	newboard.ProofOfWork = "MIM1:20::::VcoyLilzglhVKYdG:687101:"
	newboard.Description = "my updated description"
	// newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 20)
	// fmt.Println(newboard.UpdateProofOfWork)
	newboard.UpdateProofOfWork = "MIM1:20::::fhUxvaQuQpePjkwr:3315783:"
	result, err := newboard.VerifyPoW("")
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if result != true {
		t.Errorf("Test failed, this PoW should be valid but it is not.")
	}
}

func TestVerifyUpdatePoW_Success_WithKey(t *testing.T) {

	// marshaledPubKey := "0401c6166de14f9a2e698951591c9c68f2b93b3f3c9f60c93785d0464d88845bfa4b849954a7b93f7e9c4a5f40acb8c8f59ea2183069ca9cdb57b4e528aaa86ebd58ec01e164ae48f0be3d63ed69e309760337047aea8b3d37f2c08f53e50937a94c2d44b2182055f34c4cef0aa5eb653c0af32493c65d8c95c79170742b352acd7af1647c"
	result, err := signedNewboardUpdated.VerifyPoW(signedNewBoardUpdatedPubkey)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else if result != true {
		t.Errorf("Test failed, this PoW should be valid but it is not.")
	}
}

func TestVerifyUpdatePoW_Fail_UpdatePoWInvalid(t *testing.T) {
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my updated description"
	newboard.ProofOfWork = "MIM1:20::::pLBjxwHwpcHNVGBk:928329:"
	newboard.UpdateProofOfWork = "MIM1:20::::uSLzxFcWbbNXaOzXS:2032317:"
	_, err := newboard.VerifyPoW("")
	errMessage := "This proof of work is invalid or malformed."
	if err == nil {
		t.Errorf("Did not return error on malformed update PoW.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

// // CreatePoW tests

func TestCreatePoW_Success_WithoutKey(t *testing.T) {
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	err := newboard2.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newboard2.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestCreatePoW_RunTwice_Success_WithoutKey(t *testing.T) {
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	err := newboard2.CreatePoW(new(ed25519.PrivateKey), 20)
	// We run the CreatePoW twice to make sure the second run is idempotent (that it properly removes the first PoW before running so as to not include the old PoW in the input to the new PoW)
	err2 := newboard2.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil || err2 != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		result, err := newboard2.VerifyPoW("")
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}
}

func TestCreatePoW_SaltIsDifferent_Success(t *testing.T) {
	var newboard2 api.Board
	newboard2.Fingerprint = "my random fingerprint2"
	newboard2.Creation = 4564654
	newboard2.Name = "my board name"
	newboard2.Description = "my board description2"
	err := newboard2.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	}
	pow := newboard2.ProofOfWork
	if pow == "MIM1:20::::xPLDnJObCsNVlgTe:165571:" {
		t.Errorf("Test failed, this PoW's salt is the same.")
	}
}

func TestCreatePoW_Success_WithKey(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var signedNewboard api.Board
	signedNewboard.Fingerprint = "my random fingerprint3"
	signedNewboard.Creation = 4564654
	signedNewboard.Name = "my board name"
	signedNewboard.Description = "my board description"
	err := signedNewboard.CreatePoW(privKey, 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
		result, err := signedNewboard.VerifyPoW(marshaledPubKey)
		if err != nil {
			t.Errorf("Test failed, err: '%s'", err)
		} else if result != true {
			t.Errorf("Test failed, this PoW should be valid but it is not.")
		}
	}

}

func TestCreatePoW_Fail_TookTooLong(t *testing.T) {
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint2"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description2"
	// In the unlikely case that your test machine can create a 32 bit hash collision in less than 30 seconds, increase it to 36 or 40. If so, on a completely unrelated note: can I borrow your computer?
	err := newboard.CreatePoW(new(ed25519.PrivateKey), 32)
	errMessage := "The timestamp took too long to create."
	if err == nil {
		t.Errorf("Did not bail after too long a time has passed.")
	} else if !strings.Contains(err.Error(), errMessage) {
		t.Errorf("Test returned an error that was different than the expected one. '%s'", err)
	}
}

// // CreateUpdatePoW tests

// Not testing the failure cases here because they all take root from proofofwork.Create, and all the failure cases in that is tested by non-update create function tests.

func TestCreateUpdatePoW_Success_WithoutKey(t *testing.T) {
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			result, err3 := newboard.VerifyPoW("")
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestCreateUpdatePoW_Success_WithKey(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(privKey, 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(privKey, 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			// fmt.Printf("%#v\n", newboard)
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			// fmt.Printf("%#v\n", marshaledPubKey)
			result, err3 := newboard.VerifyPoW(marshaledPubKey)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestCreateUpdatePoW_SignedByTwoSeparateKeys_WithKey(t *testing.T) {
	// FUTURE: This is an example of an issue that this low-level library cannot catch. In this, the user has two private keys, and it signs the creation with one and update with another. Since we are skipping the creation check after update, this library cannot know. At the higher level, what you should do is that you should look at the ownerfingerprint and check with that.
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	privKey2, err5 := signaturing.CreateKeyPair()
	if err5 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err5)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(privKey, 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(privKey2, 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			marshaledPubKey2 := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey2.PublicKey.X, privKey2.PublicKey.Y))
			result, err3 := newboard.VerifyPoW(marshaledPubKey2)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else if result != true {
				t.Errorf("Test failed, this PoW should be valid but it is not.")
			}
		}
	}
}

func TestCreateUpdatePoW_RunTwice_Success_WithoutKey(t *testing.T) {
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			newboard.Description = "I updated this board's description twice"
			err3 := newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 20)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else {
				result, err4 := newboard.VerifyPoW("")
				if err4 != nil {
					t.Errorf("Test failed, err: '%s'", err4)
				} else if result != true {
					t.Errorf("Test failed, this PoW should be valid but it is not.")
				}
			}
		}
	}
}

func TestCreateUpdatePoW_RunTwice_Success_WithKey(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(privKey, 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(privKey, 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			newboard.Description = "I updated this board's description twice"
			err3 := newboard.CreateUpdatePoW(privKey, 20)
			if err3 != nil {
				t.Errorf("Test failed, err: '%s'", err3)
			} else {
				marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
				result, err4 := newboard.VerifyPoW(marshaledPubKey)
				if err4 != nil {
					t.Errorf("Test failed, err: '%s'", err4)
				} else if result != true {
					t.Errorf("Test failed, this PoW should be valid but it is not.")
				}
			}
		}
	}
}

func TestCreateUpdatePoW_RunTwice_ForgotReUpdatePoW_Fail_WithoutKey(t *testing.T) {
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(new(ed25519.PrivateKey), 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(new(ed25519.PrivateKey), 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			newboard.Description = "I updated this board's description twice"
			// (but I forgot to generated a new UpdatePoW)
			_, err4 := newboard.VerifyPoW("")
			errMessage := "This proof of work is invalid or malformed."
			if err4 == nil {
				t.Errorf("Did not fail after invalid UpdatePoW.")
			} else if !strings.Contains(err4.Error(), errMessage) {
				t.Errorf("Test returned an error that was different than the expected one. '%s'", err4)
			}
		}
	}
}

func TestCreateUpdatePoW_RunTwice_ForgotReUpdatePoW_Fail_WithKey(t *testing.T) {
	privKey, err4 := signaturing.CreateKeyPair()
	if err4 != nil {
		t.Errorf("Key pair creation failed. Err: '%s'", err4)
	}
	var newboard api.Board
	newboard.Fingerprint = "my random fingerprint"
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	err := newboard.CreatePoW(privKey, 20)
	if err != nil {
		t.Errorf("Test failed, err: '%s'", err)
	} else {
		newboard.Description = "I updated this board's description"
		err2 := newboard.CreateUpdatePoW(privKey, 20)
		if err2 != nil {
			t.Errorf("Test failed, err: '%s'", err2)
		} else {
			newboard.Description = "I updated this board's description twice"
			// (but I forgot to generated a new UpdatePoW)
			marshaledPubKey := hex.EncodeToString(elliptic.Marshal(elliptic.P521(), privKey.PublicKey.X, privKey.PublicKey.Y))
			_, err4 := newboard.VerifyPoW(marshaledPubKey)
			errMessage := "This proof of work is invalid or malformed."
			if err4 == nil {
				t.Errorf("Did not fail after invalid UpdatePoW.")
			} else if !strings.Contains(err4.Error(), errMessage) {
				t.Errorf("Test returned an error that was different than the expected one. '%s'", err4)
			}
		}
	}
}
