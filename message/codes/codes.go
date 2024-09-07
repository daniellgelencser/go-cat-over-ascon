package codes

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/plgd-dev/go-coap/v3/pkg/math"
)

// A Code is an unsigned 16-bit coap code as defined in the coap spec.
type Code uint16

// Request Codes
const (
	GET       Code = 1
	POST      Code = 2
	PUT       Code = 3
	DELETE    Code = 4
	PROOF     Code = 5  //<- attest response: Empty (ACK)
	PROVE     Code = 6  //<- attest response: ProofNotFound | Unauthorized | Proof
	HANDSHAKE Code = 31 //<- attest response: Empty (Server Hello)
)

// Response Codes
const (
	Empty                   Code = 0
	Created                 Code = 65  // 2 *32 + 01
	Deleted                 Code = 66  // 2 *32 + 02
	Valid                   Code = 67  // 2 *32 + 03
	Changed                 Code = 68  // 2 *32 + 04
	Content                 Code = 69  // 2 *32 + 05
	Proof                   Code = 70  // 2 *32 + 06
	Continue                Code = 95  // 2 *32 + 31
	BadRequest              Code = 128 // 4 *32 + 00
	Unauthorized            Code = 129 // 4 *32 + 01
	BadOption               Code = 130 // 4 *32 + 02
	Forbidden               Code = 131 // 4 *32 + 03
	NotFound                Code = 132 // 4 *32 + 04
	MethodNotAllowed        Code = 133 // 4 *32 + 05
	NotAcceptable           Code = 134 // 4 *32 + 06
	RequestEntityIncomplete Code = 136 // 4 *32 + 08
	PreconditionFailed      Code = 140 // 4 *32 + 12
	RequestEntityTooLarge   Code = 141 // 4 *32 + 13
	UnsupportedMediaType    Code = 143 // 4 *32 + 15
	ProofNotFound           Code = 144 // 4 *32 + 16 <- attest send nonce
	TooManyRequests         Code = 157 // 4 *32 + 29
	InternalServerError     Code = 160 // 5 *32 + 00
	NotImplemented          Code = 161 // 5 *32 + 01
	BadGateway              Code = 162 // 5 *32 + 02
	ServiceUnavailable      Code = 163 // 5 *32 + 03
	GatewayTimeout          Code = 164 // 5 *32 + 04
	ProxyingNotSupported    Code = 165 // 5 *32 + 05
)

// Signaling Codes for TCP
const (
	CSM     Code = 225
	Ping    Code = 226
	Pong    Code = 227
	Release Code = 228
	Abort   Code = 229
)

const _maxCode = 255

var _maxCodeLen int

var strToCode = map[string]Code{
	`"GET"`:                                GET,
	`"POST"`:                               POST,
	`"PUT"`:                                PUT,
	`"DELETE"`:                             DELETE,
	`"Created"`:                            Created,
	`"Deleted"`:                            Deleted,
	`"Valid"`:                              Valid,
	`"Changed"`:                            Changed,
	`"Content"`:                            Content,
	`"BadRequest"`:                         BadRequest,
	`"Unauthorized"`:                       Unauthorized,
	`"BadOption"`:                          BadOption,
	`"Forbidden"`:                          Forbidden,
	`"NotFound"`:                           NotFound,
	`"MethodNotAllowed"`:                   MethodNotAllowed,
	`"NotAcceptable"`:                      NotAcceptable,
	`"PreconditionFailed"`:                 PreconditionFailed,
	`"RequestEntityTooLarge"`:              RequestEntityTooLarge,
	`"UnsupportedMediaType"`:               UnsupportedMediaType,
	`"TooManyRequests"`:                    TooManyRequests,
	`"InternalServerError"`:                InternalServerError,
	`"NotImplemented"`:                     NotImplemented,
	`"BadGateway"`:                         BadGateway,
	`"ServiceUnavailable"`:                 ServiceUnavailable,
	`"GatewayTimeout"`:                     GatewayTimeout,
	`"ProxyingNotSupported"`:               ProxyingNotSupported,
	`"Capabilities and Settings Messages"`: CSM,
	`"Ping"`:                               Ping,
	`"Pong"`:                               Pong,
	`"Release"`:                            Release,
	`"Abort"`:                              Abort,
}

func getMaxCodeLen() int {
	// maxLen uint32 as string binary representation: "0b" + 32 digits
	maxLen := 34
	for k := range strToCode {
		kLen := len(k)
		if kLen > maxLen {
			maxLen = kLen
		}
	}
	return maxLen
}

func init() {
	_maxCodeLen = getMaxCodeLen()
}

// UnmarshalJSON unmarshals b into the Code.
func (c *Code) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return errors.New("nil receiver passed to UnmarshalJSON")
	}

	if len(b) > _maxCodeLen {
		return fmt.Errorf("invalid code: input too large(length=%d)", len(b))
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= _maxCode {
			return fmt.Errorf("invalid code: %q", ci)
		}
		*c = math.CastTo[Code](ci)
		return nil
	}

	if jc, ok := strToCode[string(b)]; ok {
		*c = jc
		return nil
	}
	return fmt.Errorf("invalid code: %v", b)
}
