package fingerprinting_test

import (
	"aether-core/aether/io/api"
	// "fmt"
	// "log"
	"os"
	// "strings"
	"testing"
)

// Infrastructure, setup and teardown

// Initialise the test object.

var newboard api.Board

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func setup() {

	newboard.Fingerprint = ""
	newboard.Creation = 4564654
	newboard.Name = "my board name"
	newboard.Description = "my board description"
	newboard.ProofOfWork = "MIM1:20::::QkaMjkJbvXInQLtW:1166891:"
}

func teardown() {
}

func ValidateTest(expected interface{}, actual interface{}, t *testing.T) {
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

// Tests

// CreateFingerprint tests come first because they are used in the CreatePoW tests as verification.

// // CreateFingerprint tests

func TestCreateFingerprint_Success(t *testing.T) {
	newboard.CreateFingerprint()
	if newboard.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	isValid := newboard.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestCreateFingerprintWithPriorFP_Success(t *testing.T) {
	newboard.Fingerprint = "there is a prior fingerprint"
	newboard.CreateFingerprint()
	if newboard.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	if newboard.Fingerprint == "there is a prior fingerprint" {
		t.Errorf("Newly created fingerprint wasn't inserted")
	}
	isValid := newboard.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestCreateFingerprintWithPriorFPAndChanges_Success(t *testing.T) {
	newboard.Fingerprint = "there is a prior fingerprint"
	newboard.CreateFingerprint()
	intermediateFp := newboard.Fingerprint
	newboard.Name = "my name changed!"
	newboard.CreateFingerprint()
	if newboard.Fingerprint == "" {
		t.Errorf("Fingerprint wasn't created")
	}
	if newboard.Fingerprint == "there is a prior fingerprint" {
		t.Errorf("Newly created fingerprint wasn't inserted")
	}
	if newboard.Fingerprint == intermediateFp {
		t.Errorf("Fingerprint on object is still from the first call of the CreateFingerprint.")
	}
	isValid := newboard.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

// // Verify Fingerprint tests

func TestVerifyFingerprint_Success(t *testing.T) {
	newboard.CreateFingerprint()
	isValid := newboard.VerifyFingerprint()
	if isValid == false {
		t.Errorf("Created fingerprint is invalid.")
	}
}

func TestVerifyFingerprint_Fail(t *testing.T) {
	newboard.Fingerprint = "my fake fingerprint"
	isValid := newboard.VerifyFingerprint()
	if isValid == true {
		t.Errorf("Verify failed to detect invalid fingerprint")
	}
}
