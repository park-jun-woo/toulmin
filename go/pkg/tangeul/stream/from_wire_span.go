//ff:func feature=tangl type=util control=sequence
//ff:what fromWireSpan — narrow a wire toulmin.Span to a logical doc.Span, rejecting values past 32 bits
package stream

import (
	"errors"
	"math"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// fromWireSpan converts a wire toulmin.Span back to a logical doc.Span. doc.Span is
// 32-bit, so an offset or length beyond that is reported instead of truncated.
func fromWireSpan(s toulmin.Span) (doc.Span, error) {
	if s.Off > math.MaxUint32 || s.Len > math.MaxUint32 {
		return doc.Span{}, errors.New("tangeul/stream: span exceeds 32-bit source offsets")
	}
	return doc.Span{Off: uint32(s.Off), Len: uint32(s.Len)}, nil
}
