//ff:func feature=tangl type=util control=selection
//ff:what spanOfRaw — the SPAN of a NODE or LINK record (decoded at the document widths)
package stream

import (
	"errors"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// spanOfRaw decodes rp just far enough to return its SPAN; only NODE and LINK
// records carry one.
func spanOfRaw(d *dec, rp rawPacket) (toulmin.Span, error) {
	switch rp.Class {
	case classNode:
		p, err := toulmin.DecodeNode(rp.Rec, d.wire)
		return p.Span, err
	case classLink:
		p, err := toulmin.DecodeEdge(rp.Rec, d.wire)
		return p.Span, err
	default:
		return toulmin.Span{}, errors.New("tangeul/stream: spanOfRaw on a span-less record")
	}
}
