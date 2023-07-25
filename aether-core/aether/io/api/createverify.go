// API > Create / Verify
// This file provides the creation and verification commands that we use for each entity.

package api

import (
	// "fmt"
	// "aether-core/aether/services/fingerprinting"
	"aether-core/aether/services/ca"
	"aether-core/aether/services/configstore"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/randomhashgen"
	"aether-core/aether/services/signaturing"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/ed25519"
	// "github.com/davecgh/go-spew/spew"
)

func isFrontend() bool {
	if globals.BackendTransientConfig == nil {
		return true
	}
	return false
}

// // Create ProofOfWork
func (b *Board) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if b.GetVersion() == 1 {
		return createBoardPoW_V1(b, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", b)
	}
}

func (t *Thread) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if t.GetVersion() == 1 {
		return createThreadPoW_V1(t, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", t)
	}
}

func (p *Post) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if p.GetVersion() == 1 {
		return createPostPoW_V1(p, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", p)
	}
}

func (v *Vote) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if v.GetVersion() == 1 {
		return createVotePoW_V1(v, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", v)
	}
}

func (k *Key) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if k.GetVersion() == 1 {
		return createKeyPoW_V1(k, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", k)
	}
}

func (ts *Truststate) CreatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if ts.GetVersion() == 1 {
		return createTruststatePoW_V1(ts, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts)
	}
}

// Create UpdateProofOfWork

func (b *Board) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if b.GetVersion() == 1 {
		return createBoardUpdatePoW_V1(b, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", b)
	}
}

func (t *Thread) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if t.GetVersion() == 1 {
		return createThreadUpdatePoW_V1(t, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", t)
	}
}

func (p *Post) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if p.GetVersion() == 1 {
		return createPostUpdatePoW_V1(p, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", p)
	}
}

func (v *Vote) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if v.GetVersion() == 1 {
		return createVoteUpdatePoW_V1(v, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", v)
	}
}

func (k *Key) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if k.GetVersion() == 1 {
		return createKeyUpdatePoW_V1(k, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", k)
	}
}

func (ts *Truststate) CreateUpdatePoW(keyPair *ed25519.PrivateKey, difficulty int) error {
	if ts.GetVersion() == 1 {
		return createTruststateUpdatePoW_V1(ts, keyPair, difficulty)
	} else {
		return fmt.Errorf("PoW update creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts)
	}
}

// Verify ProofOfWork

func (b *Board) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if b.GetVersion() == 1 {
		return verifyBoardPoW_V1(b, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false, nil
	}
}

func (t *Thread) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if t.GetVersion() == 1 {
		return verifyThreadPoW_V1(t, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false, nil
	}
}

func (p *Post) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if p.GetVersion() == 1 {
		return verifyPostPoW_V1(p, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false, nil
	}
}

func (v *Vote) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if v.GetVersion() == 1 {
		return verifyVotePoW_V1(v, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false, nil
	}
}

func (k *Key) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if k.GetVersion() == 1 {
		return verifyKeyPoW_V1(k, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false, nil
	}
}

func (ts *Truststate) VerifyPoW(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
		return true, nil
	}
	if ts.GetVersion() == 1 {
		return verifyTruststatePoW_V1(ts, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false, nil
	}
}

// Create Fingerprint

func (b *Board) CreateFingerprint() error {
	if b.GetVersion() == 1 {
		createBoardFp_V1(b)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", b)
	}
}

func (t *Thread) CreateFingerprint() error {
	if t.GetVersion() == 1 {
		createThreadFp_V1(t)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", t)
	}
}

func (p *Post) CreateFingerprint() error {
	if p.GetVersion() == 1 {
		createPostFp_V1(p)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", p)
	}
}

func (v *Vote) CreateFingerprint() error {
	if v.GetVersion() == 1 {
		createVoteFp_V1(v)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", v)
	}
}

func (k *Key) CreateFingerprint() error {
	if k.GetVersion() == 1 {
		createKeyFp_V1(k)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", k)
	}
}

func (ts *Truststate) CreateFingerprint() error {
	if ts.GetVersion() == 1 {
		createTruststateFp_V1(ts)
		return nil
	} else {
		return fmt.Errorf("Fingerprint creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts)
	}
}

// Verify Fingerprint
func (b *Board) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if b.GetVersion() == 1 {
		return verifyBoardFingerprint_V1(b)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false
	}
}

func (t *Thread) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if t.GetVersion() == 1 {
		return verifyThreadFingerprint_V1(t)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false
	}
}

func (p *Post) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if p.GetVersion() == 1 {
		return verifyPostFingerprint_V1(p)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false
	}
}

func (v *Vote) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if v.GetVersion() == 1 {
		return verifyVoteFingerprint_V1(v)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false
	}
}

func (k *Key) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if k.GetVersion() == 1 {
		return verifyKeyFingerprint_V1(k)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false
	}
}

func (ts *Truststate) VerifyFingerprint() bool {
	if !isFrontend() && !globals.BackendTransientConfig.FingerprintCheckEnabled {
		return true
	}
	if ts.GetVersion() == 1 {
		return verifyTruststateFingerprint_V1(ts)
	} else {
		logging.Log(1, fmt.Sprintf("Fingerprint verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false
	}
}

// Signature

func (b *Board) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if b.GetVersion() == 1 {
		return createBoardSignature_V1(b, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", b)
	}
}

func (t *Thread) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if t.GetVersion() == 1 {
		return createThreadSignature_V1(t, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", t)
	}
}

func (p *Post) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if p.GetVersion() == 1 {
		return createPostSignature_V1(p, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", p)
	}
}

func (v *Vote) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if v.GetVersion() == 1 {
		return createVoteSignature_V1(v, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", v)
	}
}

func (k *Key) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if k.GetVersion() == 1 {
		return createKeySignature_V1(k, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", k)
	}
}

func (ts *Truststate) CreateSignature(keyPair *ed25519.PrivateKey) error {
	if ts.GetVersion() == 1 {
		return createTruststateSignature_V1(ts, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts)
	}
}

// Create UpdateSignature

func (b *Board) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if b.GetVersion() == 1 {
		return createBoardUpdateSignature_V1(b, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", b)
	}
}

func (t *Thread) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if t.GetVersion() == 1 {
		return createThreadUpdateSignature_V1(t, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", t)
	}
}

func (p *Post) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if p.GetVersion() == 1 {
		return createPostUpdateSignature_V1(p, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", p)
	}
}

func (v *Vote) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if v.GetVersion() == 1 {
		return createVoteUpdateSignature_V1(v, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", v)
	}
}

func (k *Key) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if k.GetVersion() == 1 {
		return createKeyUpdateSignature_V1(k, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", k)
	}
}

func (ts *Truststate) CreateUpdateSignature(keyPair *ed25519.PrivateKey) error {
	if ts.GetVersion() == 1 {
		return createTruststateUpdateSignature_V1(ts, keyPair)
	} else {
		return fmt.Errorf("Signature creation of this version of this entity is not supported in this version of the app. Entity: %#v", ts)
	}
}

// Verify Signature

func (b *Board) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(b.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if b.GetVersion() == 1 {
		return verifyBoardSignature_V1(b, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", b))
		return false, nil
	}
}

func (t *Thread) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(t.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if t.GetVersion() == 1 {
		return verifyThreadSignature_V1(t, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", t))
		return false, nil
	}
}

func (p *Post) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(p.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if p.GetVersion() == 1 {
		return verifyPostSignature_V1(p, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", p))
		return false, nil
	}
}

func (v *Vote) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(v.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if v.GetVersion() == 1 {
		return verifyVoteSignature_V1(v, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", v))
		return false, nil
	}
}

func (k *Key) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(k.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if k.GetVersion() == 1 {
		return verifyKeySignature_V1(k, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", k))
		return false, nil
	}
}

func (ts *Truststate) VerifySignature(pubKey string) (bool, error) {
	if !isFrontend() && !globals.BackendTransientConfig.SignatureCheckEnabled {
		// If signature check is disabled with a debug flag, then we unconditionally return true.
		return true, nil
	}
	if !isFrontend() && globals.BackendConfig.GetAllowUnsignedEntities() && len(ts.Signature) == 0 {
		// If Allow Unsigned Entities is true, we allow for anonymous posts without signature, but if there is a signature present, we still want to do the signature check. Allow Unsigned Entities does not mean that we will allow invalid signatures.
		return true, nil
	}
	if ts.GetVersion() == 1 {
		return verifyTruststateSignature_V1(ts, pubKey)
	} else {
		logging.Log(1, fmt.Sprintf("Signature verification of this version of this entity is not supported in this version of the app. Entity: %#v", ts))
		return false, nil
	}
}

// Api Response Signature Create / Verify

func (ar *ApiResponse) CreateSignature(keyPair *ed25519.PrivateKey) error {
	// Unlike other signatures, ApiResponse signature includes the key that it is signed by itself, because it does not have a separate fingerprint field. By including the key within the signature, we protect the key under the seal of the signature, as well.
	cpI := *ar
	// Remove signature just in case, if it's been accidentally set.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Create signature
	signature, err := signaturing.Sign(string(res), keyPair)
	if err != nil {
		return err
	}
	ar.Signature = Signature(signature)
	return nil
}

// VerifySignature verifies the signature of the page. Since the public key the page is verified by is within the page itself, it does not need the public key to be given from the outside.
func (ar *ApiResponse) VerifySignature() (bool, error) {
	// 1) Check if signature check is enabled.
	if !isFrontend() && !globals.BackendTransientConfig.PageSignatureCheckEnabled {
		return true, nil
	}
	// 2) Check if required fields are empty.
	if !(len(ar.NodePublicKey) > 0 && len(ar.Signature) > 0) {
		return false, fmt.Errorf("Page signature check is enabled, but the page has some fields (Public Key or Signature) empty. Public Key: %s, Signature: %s", ar.NodePublicKey, ar.Signature)
	}
	// 3) Verify signature.
	cpI := *ar
	var signature string
	// Determine if we are checking for original or update signature
	// Save signature to be verified
	signature = string(cpI.Signature)
	// This happens *after* Signature, so should be empty here.
	cpI.Signature = ""
	// Convert to JSON
	res, _ := json.Marshal(cpI)
	// Verify Signature
	verifyResult := signaturing.Verify(string(res), signature, ar.NodePublicKey)
	// If the Signature is valid
	if verifyResult {
		return true, nil
	} else {
		return false, fmt.Errorf("This signature is invalid, but no reason given as to why. Signature: %s", signature)
	}
}

/*
	Methods below refer to the ApiResponse, which is only used by the backend. It's OK for them to refer exclusively to the backendconfig.
*/

// ApiResponse PoW Create / Verify

func (ar *ApiResponse) CreatePoW() error {
	if ar.GetVersion() == 1 {
		return createApiResponsePoW_V1(ar, globals.BackendConfig.GetBackendKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().ApiResponse)
	} else {
		return fmt.Errorf("PoW creation of this version of this entity is not supported in this version of the app. Entity: %#v", ar)
	}
}

func (ar *ApiResponse) VerifyPoW() (bool, error) {
	// if !globals.BackendTransientConfig.ProofOfWorkCheckEnabled {
	// 	return true, nil
	// }
	// ^ Heads up
	if ar.GetVersion() == 1 {
		return verifyApiResponsePoW_V1(ar, ar.NodePublicKey)
	} else {
		logging.Log(1, fmt.Sprintf("PoW verification of this version of this entity is not supported in this version of the app. Entity: %#v", ar))
		return false, nil
	}
}

// ApiResponse Nonce Create / Verify

func (ar *ApiResponse) CreateNonce() {
	rh, err := randomhashgen.GenerateInsecureRandomHash()
	if err != nil {
		logging.Logf(1, "There was an error in creating a nonce for this apiresponse. Error: %v, ApiResponse: %#v", err, ar)
	}
	ar.Nonce = Nonce(rh)
	// logging.Logf(1, "Nonce creation request, we created this: %v", ar.Nonce)
}

func (ar *ApiResponse) VerifyNonce() bool {
	// logging.Logf(1, "Nonce verification request, we got this: %v", ar.Nonce)
	return globals.BackendTransientConfig.Nonces.IsValid(ar.NodePublicKey, string(ar.Nonce), int64(ar.Timestamp))
}

// Verification for the provable and for the response.
func Verify(e interface{}) error {
	switch entity := e.(type) {
	case Provable:
		encrypted := len(entity.GetEncrContent()) > 0
		if encrypted {
			return fmt.Errorf("This item appears to be encrypted. Please decrypt before requesting verification. EncrContent: %s, Entity: %#v", entity.GetEncrContent(), entity)
		}
		realmed := len(entity.GetRealmId()) > 0
		if realmed {
			return fmt.Errorf("This item appears to belong to a realm that is different than the mainnet. Non-mainnet realms are currently not supported, but might be in the future. RealmId: %s, Entity: %#v", entity.GetRealmId(), entity)
		}
		boundsOk, err := entity.CheckBounds()
		if err != nil {
			return err
		}
		if !boundsOk {
			return fmt.Errorf("Field boundaries of this entity is invalid. Entity: %#v", entity)
		}
		fpOk := entity.VerifyFingerprint()
		if !fpOk {
			return fmt.Errorf("Fingerprint of this entity is invalid. Fingerprint: %s, Entity: %#v\n", entity.GetFingerprint(), entity)
		}
		// Bounds ok, Fp ok
		powOk, err2 := entity.VerifyPoW(entity.GetOwnerPublicKey())
		if err2 != nil {
			return err2
		}
		if !powOk {
			return fmt.Errorf("ProofOfWork of this entity is invalid. ProofOfWork: %s, Entity: %#v\n", entity.GetProofOfWork(), entity)
		}
		// Bounds ok, Fp ok, PoW ok
		sigOk, err3 := entity.VerifySignature(entity.GetOwnerPublicKey())
		if err3 != nil {
			return err3
		}
		if !sigOk {
			return fmt.Errorf("Signature of this entity is invalid. Signature: %s, Entity: %#v\n", entity.GetSignature(), entity)
		}
		// Bounds ok, Fp ok, PoW ok, Sig ok
		entOk := entity.VerifyEntitlements()
		if !entOk {
			return fmt.Errorf("Entitlements of this entity is invalid. This entity is attempting to do something that it is not authorised to do. (Ex: A CA-specific TypeClass from a CA that we do not trust.) Entity: %#v\n", entity)
		}
		badlistOk := entity.NotInBadlist()
		if !badlistOk {
			return fmt.Errorf("This entity is in a badlist, either directly or indirectly (via its parent being in a badlist) Entity: %#v\n", entity)
		}
		entity.SetVerified(true)
		return nil

	case *Address:
		boundsOk, err := entity.CheckBounds()
		if err != nil {
			return err
		}
		if !boundsOk {
			return fmt.Errorf("Field boundaries of this entity is invalid. Entity: %#v", entity)
		}
		badlistOk := entity.NotInBadlist()
		if !badlistOk {
			return fmt.Errorf("This entity is in a badlist, either directly or indirectly (via its parent being in a badlist) Entity: %#v\n", entity)
		}
		// Bounds ok
		entity.SetVerified(true)
		return nil

	default:
		return fmt.Errorf("Verify could not recognise this entity type. Entity: %#v", entity)
	}

}

func (e *Board) VerifyEntitlements() bool {
	return true
}

func (e *Thread) VerifyEntitlements() bool {
	return true
}

func (e *Post) VerifyEntitlements() bool {
	return true
}

func (e *Vote) VerifyEntitlements() bool {
	return true
}

func (e *Key) VerifyEntitlements() bool {
	return true
}

/*
- TypeClass [2] naming ([1] nameassign) [only available for CAs]
- TypeClass [3] f451 ([1] censorassign) [only available for CAs]
*/
func (e *Truststate) VerifyEntitlements() bool {
	if e.TypeClass == 2 || e.TypeClass == 3 {
		if !ca.IsTrustedCAKeyByPK(e.OwnerPublicKey) {
			return false
		}
	}
	return true
}

/*======================================
=            Badlist checks            =
======================================*/

func (e *Board) NotInBadlist() bool {
	fp, ownerfp := string(e.Fingerprint), string(e.Owner)
	return !configstore.BadlistInstance.IsBadBoard(fp, ownerfp)
}

func (e *Thread) NotInBadlist() bool {
	fp, boardfp, ownerfp := string(e.Fingerprint), string(e.Board), string(e.Owner)
	return !configstore.BadlistInstance.IsBadThread(fp, boardfp, ownerfp)
}

func (e *Post) NotInBadlist() bool {
	fp, boardfp, threadfp, parentfp, ownerfp := string(e.Fingerprint), string(e.Board), string(e.Thread), string(e.Parent), string(e.Owner)
	return !configstore.BadlistInstance.IsBadPost(fp, boardfp, threadfp, parentfp, ownerfp)
}

func (e *Vote) NotInBadlist() bool {
	/*
		This costs 4 map accesses: fp, board, thread, target and it checks whether any of those are bad. And since votes are signal entities, they don't even show up on the frontend if their main entity isn't present. This means, effectively, this is a lot of checks for no real world purpose.

		I'm still leaving this code commented our here in the case there is an unforeseen reason this might need to be enabled.
	*/
	// fp, boardfp, threadfp, targetfp := string(e.Fingerprint), string(e.Board), string(e.Thread), string(e.Target)
	// return !configstore.BadlistInstance.IsBadVote(fp, boardfp, threadfp, targetfp)
	return true
}

func (e *Key) NotInBadlist() bool {
	fp := string(e.Fingerprint)
	return !configstore.BadlistInstance.IsBadKey(fp)
}

func (e *Truststate) NotInBadlist() bool {
	fp, targetfp, ownerfp := string(e.Fingerprint), string(e.Target), string(e.Owner)
	return !configstore.BadlistInstance.IsBadTruststate(fp, targetfp, ownerfp)
}

func (e *Address) NotInBadlist() bool {
	loc, subloc, port := string(e.Location), string(e.Sublocation), uint16(e.Port)
	return !configstore.BadlistInstance.IsBadAddress(loc, subloc, port)
}

/*=====  End of Badlist checks  ======*/
