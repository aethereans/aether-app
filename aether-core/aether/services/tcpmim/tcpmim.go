// Services > TCPMim

// This library parses TCP-only Mim commands. Most of Mim is not raw TCP, it's built on HTTP. TCPMim is a separate protocol that runs on raw TCP.

// This is used for things that are impossible to do in HTTP, like doing reverse connection opens. If you want to start reading Mim code, this is not what you're looking for. This is TCPMim, not Mim.

// Do not use this to add any functionality that can be added via regular HTTP methods. This place is for stuff that requires raw TCP reads.

// Tread lightly and only use this when all other methods are exhausted. Many, many things can go wrong here if you're not paranoid. Sterling Archer says: this is officially Danger Zone™.

/*
Sample Mim message:
MIM 9 ROR

Maximum message size:
255 bytes (as of now)

No delimiter - the messages come length prefixed, and if a message is not long enough as its prefix we timeout. If it's longer, we bail with invalid message.
*/

package tcpmim

// MimMsg type

type TCPMimMessage uint

func (r TCPMimMessage) String() string {
	codePoints := [...]string{
		"InvalidMessage",
		"UnknownMessage",
		"ReverseOpenRequest",
		// This set has to match the set in const() and its order.
	}
	if r < InvalidMessage || r > ReverseOpenRequest {
		return "Invalid Code point for Mim message."
	}
	return codePoints[r]
}

// TCPMim message types
const (
	InvalidMessage     TCPMimMessage = 0
	UnknownMessage     TCPMimMessage = 1
	ReverseOpenRequest TCPMimMessage = 2
	// This set has to match the set in codePoints and its order.
)

// TCPMim max values
const (
	maxMimMsgSize = ^uint8(0)
)

// ParseMimMessage handles parsing of TCP mim messages.
func ParseMimMessage(rawmsg []byte) TCPMimMessage {
	// Check if given slice is longer than max size of a Mim message. If so, somebody (me) effed up upstream and sent us more than maximum possible bytes.
	if len(rawmsg) > int(maxMimMsgSize) {
		return InvalidMessage
	}
	// Check if Mim message
	if string(rawmsg[0:3]) != "MIM" {
		return InvalidMessage
	}
	// Check the declared size of the message
	ln := uint8(rawmsg[4:5][0])
	if ln > maxMimMsgSize {
		return InvalidMessage
	}
	if int(ln) != len(rawmsg) {
		return InvalidMessage
	}
	msgBody := rawmsg[6:ln] // Header: "MIM X " X being uint8. = 6
	if string(msgBody) == "ROR" {
		return ReverseOpenRequest
	}
	return UnknownMessage
}

func MakeMimMessage(codePoint TCPMimMessage) []byte {
	msgCode := []byte{}
	if codePoint == 2 {
		msgCode = append(msgCode, []byte("ROR")...)
	}
	msgHeader := []byte("MIM ")
	msgBody := append([]byte(" "), msgCode...)
	msgLen := len(msgHeader) + len(msgBody) + 1 // +1 for itself, +1 for the \n at the end as the delimiter.
	msgHeader = append(msgHeader, uint8(msgLen))
	msg := append(msgHeader, msgBody...)
	return msg
}
