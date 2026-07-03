//ff:func feature=tangl type=util control=selection
//ff:what DecodeControl — decode a Control packet, verifying the framed length, returning words consumed
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// DecodeControl decodes a Control packet from the head of ws, returning the
// packet and the total words consumed. PROSE yields a SPAN body; every other
// supported control yields a single Arg word. The body must consume exactly Len
// words or it is a hard error.
func DecodeControl(ws []word.Word) (ControlPacket, int, error) {
	f, head, err := word.DecodeFrame(ws)
	if err != nil {
		return ControlPacket{}, 0, err
	}
	if f.Class != word.ClassControl {
		return ControlPacket{}, 0, fmt.Errorf("tangeul/packet: expected Control class, got %d", uint8(f.Class))
	}
	total := head + int(f.Len)
	if len(ws) < total {
		return ControlPacket{}, 0, errors.New("tangeul/packet: Control body shorter than framed length")
	}
	body := ws[head:total]
	if len(body) < 1 {
		return ControlPacket{}, 0, errors.New("tangeul/packet: Control body requires a subcode")
	}
	p := ControlPacket{Code: codebook.ControlCode(body[0])}
	switch p.Code {
	case codebook.ControlProse:
		if len(body) != 4 {
			return ControlPacket{}, 0, errors.New("tangeul/packet: PROSE control body must be subcode + SPAN")
		}
		s, err := word.DecodeSpan(body[1:])
		if err != nil {
			return ControlPacket{}, 0, err
		}
		p.Span = &s
	default:
		if len(body) != 2 {
			return ControlPacket{}, 0, errors.New("tangeul/packet: Control body must be subcode + Arg")
		}
		p.Arg = uint16(body[1])
	}
	return p, total, nil
}
