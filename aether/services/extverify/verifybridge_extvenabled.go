// +build extvenabled

package extverify

import (
	"aether-core/support/extbridge"
)

func init() {
	verifier := extbridge.NewVerifier()
	Verifier = VerifierInterface(&verifier)
}
