//ff:func feature=tangl type=util control=sequence
//ff:what DecodeText — decode a Text payload, checking format and byte-length padding, returning words consumed
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// DecodeText decodes a Text payload from the head of ws, returning the packet
// and the total words consumed. It enforces the UTF-8 format code and that the
// data words are exactly ceil(byteLen/2) — i.e. the padding is a single trailing
// 0x00 at most. TID is left zero (stream-assigned).
func DecodeText(ws []word.Word) (TextPacket, int, error) {
	f, head, err := word.DecodeFrame(ws)
	if err != nil {
		return TextPacket{}, 0, err
	}
	if f.Class != word.ClassPayload {
		return TextPacket{}, 0, fmt.Errorf("tangeul/packet: expected Payload class, got %d", uint8(f.Class))
	}
	total := head + int(f.Len)
	if len(ws) < total {
		return TextPacket{}, 0, errors.New("tangeul/packet: Text body shorter than framed length")
	}
	body := ws[head:total]
	if len(body) < 3 {
		return TextPacket{}, 0, errors.New("tangeul/packet: Text body requires format and byte length")
	}
	if TextFormat(body[0]) != TextFormatUTF8 {
		return TextPacket{}, 0, fmt.Errorf("tangeul/packet: unsupported Text format 0x%04X", uint16(body[0]))
	}
	byteLen, err := word.DecodeNum32(body[1:3])
	if err != nil {
		return TextPacket{}, 0, err
	}
	data := body[3:]
	if int(byteLen) != len(data)*2 && int(byteLen) != len(data)*2-1 {
		return TextPacket{}, 0, errors.New("tangeul/packet: Text byte length inconsistent with padded word count")
	}
	return TextPacket{Bytes: word.PutWords(data)[:byteLen]}, total, nil
}
