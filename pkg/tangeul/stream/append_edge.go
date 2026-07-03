//ff:func feature=tangl type=util control=sequence
//ff:what appendEdge — encode an Edge packet and append its words to a body buffer
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// appendEdge encodes ep and appends the resulting words to *body.
func appendEdge(body *[]word.Word, ep packet.EdgePacket) error {
	w, err := packet.EncodeEdge(ep)
	if err != nil {
		return err
	}
	*body = append(*body, w...)
	return nil
}
