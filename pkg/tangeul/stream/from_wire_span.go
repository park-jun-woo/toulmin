//ff:func feature=tangl type=util control=sequence
//ff:what fromWireSpan — widen a 16-bit wire word.Span back to a logical doc.Span
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// fromWireSpan widens a wire word.Span (16-bit offset/length, anchored at the
// Source Text payload) back to a logical doc.Span. The anchor TID is dropped:
// the reader replays spans directly against Doc.Source, so only the byte range
// matters.
func fromWireSpan(s word.Span) doc.Span {
	return doc.Span{Off: uint32(s.Off), Len: uint32(s.Len)}
}
