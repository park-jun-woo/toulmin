//ff:func feature=tangl type=util control=sequence
//ff:what labelEdge — build a label Edge (0x010F) binding a node TID to a Text TID
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// labelEdge encodes a label Edge (PropLabel, §2.3) whose src is a node TID and
// whose dst is a Text TID, anchored at span. Backtick names and function names
// ride as label Text plus this edge (README §6, phase design) so the original
// spelling is the lossless source of truth while VerbCode stays unbound.
func labelEdge(src, textTID uint16, span word.Span) ([]word.Word, error) {
	return packet.EncodeEdge(packet.EdgePacket{
		Prop:   codebook.PropLabel,
		Src:    packet.TID(src),
		DstTID: packet.TID(textTID),
		Span:   span,
	})
}
