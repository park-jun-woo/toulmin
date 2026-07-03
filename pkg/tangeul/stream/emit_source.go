//ff:func feature=tangl type=util control=sequence
//ff:what emitSource — encode the Source Text payload (authored bytes verbatim, the SPAN anchor)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitSource encodes the Source Text payload (§7.4): the authored bytes captured
// verbatim (no Unicode normalization) as one Text payload that every packet's
// SPAN anchors into. It is the sole Payload emitted before the graph packets, so
// the reader can lift it off the head of the stream and replay spans against it.
func emitSource(src []byte) []word.Word {
	return packet.EncodeText(packet.TextPacket{Bytes: src})
}
