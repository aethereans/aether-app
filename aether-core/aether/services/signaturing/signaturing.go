// Services > Fingerprinting
// This module handles the creation of signatures, signing entities, and checking for signatures.

package signaturing

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/ed25519"
)

/*
This is where we create keys, sign entities, and validate them. This does have a dependency on the config system, in that the user's private key is needs to be saved to the configuration, which is a subsystem of its own.
*/

func MarshalPublicKey(key ed25519.PublicKey) string {
	return hex.EncodeToString(key)
}

func UnmarshalPublicKey(hexKey string) (ed25519.PublicKey, error) {
	keyAsByte, err := hex.DecodeString(hexKey)
	if err != nil {
		return ed25519.PublicKey{}, err
	}
	return ed25519.PublicKey(keyAsByte), nil
}

func MarshalPrivateKey(key ed25519.PrivateKey) string {
	return hex.EncodeToString(key)
}

func UnmarshalPrivateKey(hexKey string) (ed25519.PrivateKey, error) {
	keyAsByte, err := hex.DecodeString(hexKey)
	if err != nil {
		return ed25519.PrivateKey{}, err
	}
	return ed25519.PrivateKey(keyAsByte), nil
}

func CreateKeyPair() (*ed25519.PrivateKey, error) {
	_, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return new(ed25519.PrivateKey), errors.New(fmt.Sprint(
			"Key pair generation failed. err: ", err))
	}
	return &privKey, nil
	// Note: accessing the public key is shown below.
	// var pubKey ed25519.PublicKey
	// pubKey = privKey.PublicKey
}

func Sign(input string, privKey *ed25519.PrivateKey) (string, error) {
	// This creates an ed25519 signature for the input provided.
	// Mind that signatures are generated on hashes of items, not the item itself. So the first step in either signing or verifying is to hash the input provided.
	inputByte := []byte(input)
	hasher := sha256.New()
	hasher.Write(inputByte)
	hash := hasher.Sum(nil)
	sigAsByte, err := privKey.Sign(nil, hash, crypto.Hash(0))
	if err != nil {
		return "", errors.New(fmt.Sprint(
			"Signing failed. err: ", err))
	}
	signature := hex.EncodeToString(sigAsByte)
	return signature, nil
}

func Verify(input string, signature string, pubKey string) bool {
	// This verifies the input provided by a given signature and public key.
	// Mind that signatures are generated on hashes of items, not the item itself. So the first step in either signing or verifying is to hash the input provided.
	if !(len(signature) > 0 && len(pubKey) > 0) {
		return false
	}
	inputByte := []byte(input)
	hasher := sha256.New()
	hasher.Write(inputByte)
	hash := hasher.Sum(nil)
	sigAsByte, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	pk, err := UnmarshalPublicKey(pubKey)
	if err != nil {
		return false
	}
	valid := ed25519.Verify(pk, hash, sigAsByte)
	return valid
}
