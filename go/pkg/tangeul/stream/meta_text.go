//ff:func feature=tangl type=util control=sequence
//ff:what metaText — the string of a Meta record whose payload is exactly one byte-string Value
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/core"
)

// metaText returns the text of a Meta record whose payload is exactly one
// embedded Value of the even/odd codec pair evenCodec.
func metaText(r core.Record, evenCodec uint8) (string, error) {
	v, n, err := core.DecodeRecord(r.Payload)
	if err != nil {
		return "", err
	}
	if n != len(r.Payload) {
		return "", errors.New("tangeul/stream: Meta payload has words after its Value")
	}
	b, err := core.DecodeBytesValue(v, evenCodec)
	return string(b), err
}
