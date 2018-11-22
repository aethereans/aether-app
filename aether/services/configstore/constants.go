// Services > ConfigStore > Constants
// These are the constants that cannot be user-changed without a recompile. There should be nothing here that you might want to change - if something is here, it means:
// a) there is a better way to change it elsewhere, or
// b) changing it breaks the app in a way that is unrecoverable or causes undefined behaviour

package configstore

import ()

// Hardcoded version numbers specific to this build
const (
	clientVersionMajor   = 2
	clientVersionMinor   = 0
	clientVersionPatch   = 0
	clientName           = "Aether"
	protocolVersionMajor = 1
	protocolVersionMinor = 0
)

// Bootstrapper of last resort, if no other bootstrapper is given or found. If the user or a library higher up in the stack provides a bootstrapper, that will be used instead.
const (
	DefaultBootstrapperLocation    = "bootstrap.getaether.net"
	DefaultBootstrapperSublocation = ""
	DefaultBootstrapperPort        = 443
)
