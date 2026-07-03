//ff:func feature=tangl type=util control=sequence
//ff:what toWireSpan — convert a logical doc.Span to a wire word.Span, rejecting 16-bit overflow
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// srcAnchorTID is the SPAN anchor TID every packet points at — the Source Text
// payload (§7.4). The reader replays spans directly against Doc.Source, so the
// value is nominal; it is kept distinct from the primary node range for clarity.
const srcAnchorTID uint16 = 0

// toWireSpan converts a logical doc.Span (32-bit offset/length into Source) to a
// wire word.Span anchored at the Source Text payload. word.Span offset/length
// are provisionally 16 bits (§11-1), so a Source larger than 65535 bytes is a
// hard error rather than a silent truncation.
func toWireSpan(s doc.Span) (word.Span, error) {
	if s.Off > 0xFFFF || s.Len > 0xFFFF || s.Off+s.Len > 0xFFFF {
		return word.Span{}, fmt.Errorf("tangeul/stream: span [%d,%d) exceeds 16-bit wire width", s.Off, s.Off+s.Len)
	}
	return word.Span{TID: srcAnchorTID, Off: uint16(s.Off), Len: uint16(s.Len)}, nil
}
