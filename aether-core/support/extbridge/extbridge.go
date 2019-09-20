package extbridge

import (
	"net"
	"net/http"
	"sync"
)

type CFExtVerify struct {
	lock        sync.Mutex
	IPV4CIDRs   []string
	IPV6CIDRs   []string
	Ranges      []*net.IPNet
	IPV4Source  string
	IPV6Source  string
	LastRefresh int64
}

func NewVerifier() CFExtVerify {
	return CFExtVerify{
		IPV4CIDRs:  []string{},
		IPV6CIDRs:  []string{},
		IPV4Source: "",
		IPV6Source: "",
	}
}

func (v *CFExtVerify) GetRemoteIP(h http.Header) string {
	return ""
}

func (v *CFExtVerify) IsAllowedRemoteIP(ip string) bool {
	return true
}

func (v *CFExtVerify) Invalidate(url string) {
}
