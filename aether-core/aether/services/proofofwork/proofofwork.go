// Services > ProofOfWork
// This module provides proof of work functions, to both create and validate. Proof of work creation requires a private key to be provided, as proof of work result itself also needs to be signed separately.

package proofofwork

import (
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ed25519"
	// "encoding/json"
	"errors"
	"fmt"
	"math"
	// "math/rand"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/signaturing"
	"math/big"
	"strconv"
	"strings"
	"time"
)

/*
The workings of this can be outlined as such:

0) The way hashcash works is that we check for partial hash collisions. That means, in simpler terms, we are looking to create a hash with a particular number of zeroes at the left side. The number of zeroes is the difficulty of the hashcash token - the more zeroes it has, the harder it was to generate.

1) Create a salt. This makes Hashcash random even if the text is exactly the same.
2) Calculate how many digits in hex (mind that one 0 digit in hex refers to 4 zeros in binary). This is floored, so if we need 4.5 hex digits (i.e. 22 binary digits) it will get floored to 20 binary digits = 4 hex digits. Don't worry, we will also be checking for the last 2 digits later.
3) Create the zeroed strings for both hex and binary.
4) Create the hash of the salt, the input, and the counter, combined.
5) Check whether the generated hash in hex has the number of zeroes that the hex. If not, increment the counter and try again.
6) If hex has enough digits, convert to binary, and check if binary has enough digits. Binary digits >= hex digits * 4. If so, congrats! You got a winner. If not, increment the counter and try again.
7) Format the result as needed, sign if needed, and return.
*/

// Constants

// Used for salt generation used in creation of proof of work.
const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const PADDING7 = "0000000"
const PADDING6 = "000000"
const PADDING5 = "00000"
const PADDING4 = "0000"
const PADDING3 = "000"
const PADDING2 = "00"
const PADDING1 = "0"

// Low level functions

func mimHash(input string) []byte {
	// Mim hash: SHA256 repeated 3 times.
	// Thrice hashing, with the hope that it won't be trivial to convert an ASIC miner to work on this.
	inputByte := []byte(input)
	firstPass := sha256.New()
	secondPass := sha256.New()
	thirdPass := sha256.New()
	firstPass.Write(inputByte)
	secondPass.Write(firstPass.Sum(nil))
	thirdPass.Write(secondPass.Sum(nil))
	// Compute hash with h.sum
	return thirdPass.Sum(nil)
}

func convertToBinaryString(input []byte) string {
	var result string
	for _, byt := range input {
		// result = result + fmt.Sprintf("%b", byt)
		// switch case zero padding here.
		switch len(fmt.Sprintf("%b", byt)) {
		case 1:
			result += PADDING7 + fmt.Sprintf("%b", byt)
		case 2:
			result += PADDING6 + fmt.Sprintf("%b", byt)
		case 3:
			result += PADDING5 + fmt.Sprintf("%b", byt)
		case 4:
			result += PADDING4 + fmt.Sprintf("%b", byt)
		case 5:
			result += PADDING3 + fmt.Sprintf("%b", byt)
		case 6:
			result += PADDING2 + fmt.Sprintf("%b", byt)
		case 7:
			result += PADDING1 + fmt.Sprintf("%b", byt)
		case 8:
			result += fmt.Sprintf("%b", byt)
		}
	}
	return result
}

var bailoutTimeSeconds int

// Mid level functions

// Create creates the Hashcash proof of with the given difficulty. This function has an inner loop which adds a random element to the input and tries to find enough zeros at the beginning of the SHA1 hash of the result.
func Create(input string, difficulty int, privKey *ed25519.PrivateKey) (string, error) {
	if difficulty <= 0 {
		return "", nil
		// Any value of 0 or less means there is no PoW implied.
	}
	if bailoutTimeSeconds == 0 {
		if globals.BackendConfig != nil {
			bailoutTimeSeconds = globals.BackendConfig.GetPoWBailoutTimeSeconds()
		}

		if globals.FrontendConfig != nil {
			bailoutTimeSeconds = globals.FrontendConfig.GetPoWBailoutTimeSeconds()
		}
	}
	// First of all, check if BailoutSeconds exists. If this does not exist we have to exit as the allotted maximum time until a PoW is created will be zero.
	difficulty64 := int64(difficulty)
	// if bailoutTimeSeconds == 0 {
	// 	return "", errors.New(fmt.Sprint(
	// 		"Please initialise BailoutSeconds first."))
	// }
	if bailoutTimeSeconds == 0 {
		bailoutTimeSeconds = 600
		// This is TODO to accommodate nameminter, whose config is separate. We should remove backend / frontend config dependency from this library in the future.
	}
	// Calculate number of zeros we need at the beginning of the hash.
	// Hex is floored. So you still need to do additional checks.
	lowerHexDigitsNeeded := int(math.Floor(float64(difficulty64) / float64(4)))
	binDigitsNeeded := difficulty64
	var zeroHexDigits string
	var zeroBinDigits string
	for i := 0; i < lowerHexDigitsNeeded; i++ {
		zeroHexDigits = zeroHexDigits + "0"
	}
	for i := int64(0); i < binDigitsNeeded; i++ {
		zeroBinDigits = zeroBinDigits + "0"
	}
	// Create the counter
	var counter int64
	// Before creating the salt, we need to seed the random number generator first. We check if it is already seeded, we do nothing.
	// fmt.Printf("%#v\n", rand.Seed)
	// Create the salt.
	saltBytes := make([]byte, 16)
	for i := range saltBytes {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		if err != nil {
			return "", errors.New(fmt.Sprint(
				"Random number generator generated an error. err: ", err))
		}
		saltBytes[i] = LETTERS[int(randNum.Int64())]
	}
	// Add salt to the end of the input string.
	inputToBePoWd := strconv.FormatInt(difficulty64, 10) +
		input + string(saltBytes)
	// Take time here
	timeCounter := int(time.Now().Unix())
	bailoutSeconds := bailoutTimeSeconds
	for {
		// This is the tight loop.
		// Check if the bailout time has passed.
		now := int(time.Now().Unix())
		if now-(timeCounter+bailoutSeconds) > 0 {
			return "", errors.New(fmt.Sprint(
				"The timestamp took too long to create."))
		}
		// Compute hash.
		result := mimHash(inputToBePoWd + strconv.FormatInt(counter, 10))
		resultHex := fmt.Sprintf("%x", result)
		// Check hex first.
		if resultHex[:lowerHexDigitsNeeded] == zeroHexDigits {
			// It has passed the hex portion. So we need to check for the binary.
			resultBinary := convertToBinaryString(result)
			if resultBinary[:binDigitsNeeded] == zeroBinDigits {
				// It passed both inspections, success!
				break
			} else {
				counter++
			}
		} else {
			counter++
		}
	}
	// Mind the terminating ":" in case of no signature.
	proofOfWork := "MIM1" + ":" + strconv.FormatInt(difficulty64, 10) + "::::" +
		string(saltBytes) + ":" + strconv.FormatInt(counter, 10) + ":"
	if len(signaturing.MarshalPrivateKey(*privKey)) > 0 {
		// We have a private key. Sign the hash with this key.
		// The result will be in the format of [Rest of PoW]:[Signature]
		result, err := signaturing.Sign(proofOfWork, privKey)
		if err != nil {
			return "", errors.New(fmt.Sprint(
				"This PoW could not be signed. Error:", err))
		}
		proofOfWork = proofOfWork + result
		return proofOfWork, nil
	} else {
		// We don't have a private key. Do not sign the hash.
		// The result will be in the format of [Rest of PoW]: <- mind the last ":"
		return proofOfWork, nil
	}
}

// Verify validates whether the given Hashcash token is strong enough to satisfy the given difficulty.
func Verify(input string, pow string, pubKey string) (bool, int, error) {
	// MimHashcash syntax:
	// [version]:[difficulty]:[date]:[input]:[extension]:[salt]:[counter]:[signature]
	// Check PoW length. If longer than 1024 chars, bail.
	if len(pow) > 1024 {
		return false, 0, errors.New(fmt.Sprint(
			"This PoW is longer than maximum allowed 1024 characters. PoW: ", pow))
	}
	// Field count is 8+1, because if there are more fields, we want the excess to create the field 9. If there is a field 9, we bail.
	parsedStrings := strings.SplitN(pow, ":", 9)
	// First, check for field counts. If above 8 or below, bail.
	if len(parsedStrings) != 8 {
		return false, 0, errors.New(fmt.Sprint(
			"PoW had more or less fields than expected. PoW: ", pow))
	}
	// Second, check for whether they are empty or not. Date, input, extension should always be empty. Version, difficulty, salt and counter should always be non-empty.
	if parsedStrings[2] != "" || parsedStrings[3] != "" || parsedStrings[4] != "" || parsedStrings[0] == "" || parsedStrings[1] == "" || parsedStrings[5] == "" || parsedStrings[6] == "" {
		return false, 0, errors.New(fmt.Sprint(
			"This proof of work either has fields that should be empty and is not, or it does have empty fields which it should not. PoW: ", pow))
	}
	// Create proper parse fields.
	var parsedVersion string
	var parsedDifficulty64 int64
	var parsedSalt string
	var parsedCounter int64
	var parsedSignature string
	// Attempt to parse into the properly typed fields.
	parsedVersion = parsedStrings[0]
	parsedDifficulty64, err := strconv.ParseInt(parsedStrings[1], 10, 64)
	parsedDifficulty := int(parsedDifficulty64)
	if err != nil {
		return false, 0, errors.New(fmt.Sprint(
			"PoW parsing failed, this PoW is invalid. Error: ", err))
	}
	if parsedDifficulty < 0 {
		return false, 0, errors.New(fmt.Sprint(
			"This proof of work is invalid or malformed. (Negative parsed difficulty) PoW: ", pow))
	}
	parsedSalt = parsedStrings[5]
	parsedCounter, err2 := strconv.ParseInt(parsedStrings[6], 10, 64)
	if err2 != nil {
		return false, 0, errors.New(fmt.Sprint(
			"PoW parsing failed, this PoW is invalid. Error: ", err2))
	}
	if parsedCounter < 0 {
		return false, 0, errors.New(fmt.Sprint(
			"This proof of work is invalid or malformed. (Negative counter) PoW: ", pow))
	}
	parsedSignature = parsedStrings[7]
	// Parsing complete. Check for PoW validity.
	switch parsedVersion {
	case "MIM1":
		// Add the difficulty, salt and counter to the input string.
		stringToBeVerified := strconv.FormatInt(parsedDifficulty64, 10) +
			input + parsedSalt + strconv.FormatInt(parsedCounter, 10)
		// Hash the outputted string to see the result per SHA256x3.
		result := mimHash(stringToBeVerified)
		resultBinary := convertToBinaryString(result)
		// Check for zeroes at the beginning.
		var zeroBinDigits string
		for i := int(0); i < parsedDifficulty; i++ {
			zeroBinDigits = zeroBinDigits + "0"
		}
		if resultBinary[:parsedDifficulty] == zeroBinDigits {
			// SUCCESS, the PoW has *at least* the given number of zeroes at the beginning. We will still accept the difficulty at the declared level, so that the user gets no free zeroes.

			// Signature check starts here.
			// Check if there is a signature, if a key is provided, and if signature is valid.
			if len(pubKey) > 0 && len(parsedSignature) > 0 {
				// We have both the key and the signature.
				stringToBeSignatureChecked := parsedStrings[0] + ":" + parsedStrings[1] + ":" + parsedStrings[2] + ":" + parsedStrings[3] + ":" + parsedStrings[4] + ":" + parsedStrings[5] + ":" + parsedStrings[6] + ":"
				verifyResult := signaturing.Verify(stringToBeSignatureChecked, parsedSignature, pubKey)
				if verifyResult != true {
					return verifyResult, parsedDifficulty, errors.New(fmt.Sprint(
						"The signature of this PoW is invalid. The PoW signature and the public key provided does not match. PoW: ", pow))
				}
				return verifyResult, parsedDifficulty, nil
			} else if len(pubKey) > 0 {
				// We have the key but no signature in PoW. Bail.
				return false, 0, errors.New(fmt.Sprint(
					"A key is provided, but the PoW is unsigned. PoW: ", pow))
			} else if len(parsedSignature) > 0 {
				// We have the signature but no key is provided. Bail.
				return false, 0, errors.New(fmt.Sprint(
					"The PoW is signed, but a key is not provided."))
			} else {
				// We have neither key nor signature. This means that this is a PoW for an anonymous object.
				return true, parsedDifficulty, nil
			}
		} else {
			return false, 0, errors.New(fmt.Sprint(
				"This proof of work is invalid or malformed. PoW: ", pow))
		}
	default:
		// If this has a different version, bail. This is where we would create the next version's code in, if there is any.
		return false, 0, errors.New(fmt.Sprint(
			"This proof of work is in a format Mim does not support. PoW: ", pow))
	}
}
