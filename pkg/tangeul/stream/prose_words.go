//ff:func feature=tangl type=util control=sequence
//ff:what proseWords — encode a PROSE control carrying a source span (§7.4)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// proseWords encodes a PROSE Control packet carrying the wire span of s. It backs
// both standalone prose blocks and the heading span that follows a SECTION.
func proseWords(s doc.Span) ([]word.Word, error) {
	ws, err := toWireSpan(s)
	if err != nil {
		return nil, err
	}
	return packet.EncodeControl(packet.ControlPacket{Code: codebook.ControlProse, Span: &ws})
}
