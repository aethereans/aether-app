package extverify

import (
	"net/http"
)

// Verifier is the external verifier interface. Any verifier provided in should be able to satisfy this.
type VerifierInterface interface {
	GetRemoteIP(h http.Header) string
	IsAllowedRemoteIP(ip string) bool
	Invalidate(url string)
}

var Verifier VerifierInterface
