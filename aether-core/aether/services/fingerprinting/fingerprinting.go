// Services > Fingerprinting
// This module handles the creation and verification of fingerprints

package fingerprinting

import (
	"crypto/sha256"
	"fmt"
)

/*
This should be relatively simple. The key insight here is that we would need to grab the whole thing (similar to proof of work) and feed it to the fingerprinting function. This is similar to how the PoW takes in the whole thing converted to a single string.

*/

func Create(input string) string {
	// Create a fingerprint from the string given.
	inputByte := []byte(input)
	calculator := sha256.New()
	calculator.Write(inputByte)
	resultHex := fmt.Sprintf("%x", calculator.Sum(nil))
	return resultHex
}

func Verify(input string, fingerprint string) bool {
	calculatedFingerprint := Create(input)
	if calculatedFingerprint == fingerprint {
		return true
	} else {
		return false
	}
}
