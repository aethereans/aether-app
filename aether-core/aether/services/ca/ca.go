// Services > CA
// This service handles the validation of certificate authorities.

package ca

var (
	trustedCAsPKs = []string{"142ae631bcd46ab6e548c9b2d9494c7e30c65b4389e541e31bbe31e0ae47515e"}
	trustedCAFps  = []string{"5460a18d7dd4c6b078199f5ee8ee70037f877166b07e9596cf26deb73223e15c"}
)

// IsTrustedCAKeyByPK checks whether a key is one of our trusted CA keys based on the PK. Mind that this function does not actually check whether the message that this CA key came in is valid, so you should run this after you've otherwise validated the message and you know the message is signed properly by the key that you're checking.
func IsTrustedCAKeyByPK(publicKey string) bool {
	// return true // TODO (debug)
	for key := range trustedCAsPKs {
		if publicKey == trustedCAsPKs[key] {
			return true
		}
	}
	return false
}

func IsTrustedCAKeyByPKWithPriority(publicKey string) (bool, int) {
	// return true // TODO (debug)
	for key := range trustedCAsPKs {
		if publicKey == trustedCAsPKs[key] {
			return true, key
		}
	}
	return false, -1
}

// IsTrustedCAKeyByFp checks whether a key is one of our trusted CA keys based on the Fingerprint. Mind that this function does not actually check whether the message that this CA key came in is valid, so you should run this after you've otherwise validated the message and you know the message is signed properly by the key that you're checking.
func IsTrustedCAKeyByFp(publicKey string) bool {
	// return true // TODO (debug)
	for key := range trustedCAFps {
		if publicKey == trustedCAFps[key] {
			return true
		}
	}
	return false
}

func IsTrustedCAKeyByFpWithPriority(publicKey string) (bool, int) {
	// return true // TODO (debug)
	for key := range trustedCAFps {
		if publicKey == trustedCAFps[key] {
			return true, key
		}
	}
	return false, -1
}
