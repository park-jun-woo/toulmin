//ff:func feature=tangl type=util control=sequence
//ff:what addText — assign the next Text TID, append the Text payload to the prefix, return the TID
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// addText issues the next segment-local Text TID, appends a UTF-8 Text payload
// carrying s to the statement prefix, and returns the TID. Because a Text
// payload's TID is never encoded, every label/means/with reference points at
// this stream-assigned value and the reader recovers it by walking Text payloads
// in the same order — so addText must be the sole issuer of Text TIDs and must
// append in issue order.
func addText(e *enc, prefix *[]word.Word, s string) uint16 {
	tid := e.textNext
	e.textNext++
	*prefix = append(*prefix, packet.EncodeText(packet.TextPacket{Bytes: []byte(s)})...)
	return tid
}
