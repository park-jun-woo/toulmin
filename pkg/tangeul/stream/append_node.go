//ff:func feature=tangl type=util control=sequence
//ff:what appendNode — encode a Node packet and append its words to a body buffer
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// appendNode encodes p and appends the resulting words to *body, the shared
// span-bearing region a statement's Node/Edge packets accumulate into.
func appendNode(body *[]word.Word, p packet.NodePacket) error {
	w, err := packet.EncodeNode(p)
	if err != nil {
		return err
	}
	*body = append(*body, w...)
	return nil
}
