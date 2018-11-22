// Services > RandomHashGen
// This module provides a random hash generation function. At this stage, this is not worth its own package, but it makes sense to do so here for the purpose of avoiding import cycles.
package randomhashgen

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

// GenerateInsecureRandomHash generates a *pseudo* random hash that we can use for non-cryptographically-significant places. The benefit of this is that it's much faster than a cryptographically secure one. Essentially, so long as you don't need a random hash to use as part of an encryption scheme, you're good.
// (If you're implementing an encryption scheme, please reconsider, and try to use an existing, vetted implementation.)
func GenerateInsecureRandomHash() (string, error) {
	const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	saltBytes := make([]byte, 16)
	for i := range saltBytes {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		if err != nil {
			return "", errors.New(fmt.Sprint(
				"Random number generator generated an error. err: ", err))
		}
		saltBytes[i] = LETTERS[int(randNum.Int64())]
	}
	calculator := sha256.New()
	calculator.Write(saltBytes)
	resultHex := fmt.Sprintf("%x", calculator.Sum(nil))
	return resultHex, nil
}

// GenerateSecureRandomHash does not exist because we are (thankfully) not building our own encryption.
