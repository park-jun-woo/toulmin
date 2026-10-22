//ff:func feature=tangl type=util control=sequence
//ff:what metaRecord — a HEAD Meta record whose payload is one embedded Value
package stream

import "github.com/geul-org/geul/go/pkg/core"

// metaRecord encodes a Meta record of codec whose payload is the single
// embedded Value v (GEUL t002: a Meta payload may embed Value records).
func metaRecord(codec uint8, v core.Record) ([]core.Word, error) {
	vw, err := core.EncodeRecord(v)
	if err != nil {
		return nil, err
	}
	return core.EncodeRecord(core.Record{Kind: core.KindMeta, Codec: codec, Payload: vw})
}
