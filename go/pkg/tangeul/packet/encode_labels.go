//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what encodeLabels — serialize an EXPR label table: a count word then one UTF8-TEXT Value per label
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// encodeLabels serializes an EXPR label table: a count word, then each label as
// an embedded UTF8-TEXT Value record (verbatim bytes). The table keeps an
// ExprPacket self-contained, so it round-trips without TEXT records.
func encodeLabels(labels []string) ([]core.Word, error) {
	out := []core.Word{core.Word(len(labels))}
	var err error
	for _, l := range labels {
		if out, err = toulmin.PutValue(out, core.EncodeBytesValue(core.CodecUTF8TextEven, []byte(l))); err != nil {
			return nil, err
		}
	}
	return out, nil
}
