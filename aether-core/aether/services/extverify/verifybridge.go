// +build !extvenabled

package extverify

import (
	"log"
	"net/http"
)

type StubVerifier struct{}

func (v *StubVerifier) GetRemoteIP(h http.Header) string {
	log.Fatal("You have attempted to enable external verifier gate on a build that was not compiled with the appropriate flags to include the bridge. If you'd like to use it, please recompile with the appropriate flags.")
	return ""
}

func (v *StubVerifier) IsAllowedRemoteIP(ip string) bool {
	log.Fatal("You have attempted to enable external verifier gate on a build that was not compiled with the appropriate flags to include the bridge. If you'd like to use it, please recompile with the appropriate flags.")
	return false
}

func (v *StubVerifier) Invalidate(url string) {
}

func init() {
	stub := StubVerifier{}
	Verifier = VerifierInterface(&stub)

}
