//ff:func feature=tangl type=util control=selection
//ff:what DecodeEdge — decode an Edge packet, verifying the framed length, returning words consumed
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// DecodeEdge decodes an Edge packet from the head of ws, returning the packet
// and the total words consumed. The PropCode selects the dst slot layout via
// codebook.PropDstKind; a GATE16 is read iff F_GATE is set. The slot words plus
// the trailing SPAN must consume exactly Len words or it is a hard error.
func DecodeEdge(ws []word.Word) (EdgePacket, int, error) {
	f, head, err := word.DecodeFrame(ws)
	if err != nil {
		return EdgePacket{}, 0, err
	}
	if f.Class != word.ClassEdge {
		return EdgePacket{}, 0, fmt.Errorf("tangeul/packet: expected Edge class, got %d", uint8(f.Class))
	}
	total := head + int(f.Len)
	if len(ws) < total {
		return EdgePacket{}, 0, errors.New("tangeul/packet: Edge body shorter than framed length")
	}
	body := ws[head:total]
	if len(body) < 2 {
		return EdgePacket{}, 0, errors.New("tangeul/packet: Edge body requires PropCode and src")
	}
	p := EdgePacket{Prop: codebook.PropCode(body[0]), Src: TID(body[1])}
	kind, err := codebook.PropDstKind(p.Prop)
	if err != nil {
		return EdgePacket{}, 0, err
	}
	pos := 2
	switch kind {
	case codebook.DstFuncRef:
		fr, err := decodeFuncRef(body[pos:])
		if err != nil {
			return EdgePacket{}, 0, err
		}
		p.DstFunc = &fr
		pos += 2
	case codebook.DstTID, codebook.DstText:
		if len(body) < pos+1 {
			return EdgePacket{}, 0, errors.New("tangeul/packet: Edge TID dst missing")
		}
		p.DstTID = TID(body[pos])
		pos++
	default:
		return EdgePacket{}, 0, errors.New("tangeul/packet: provenance/Block dst is not decoded in v0.1")
	}

	if len(body) < pos+1 {
		return EdgePacket{}, 0, errors.New("tangeul/packet: Edge flag word missing")
	}
	p.Flags = codebook.EdgeFlags(body[pos])
	pos++
	if p.Flags&codebook.EdgeFlagGate != 0 {
		if len(body) < pos+1 {
			return EdgePacket{}, 0, errors.New("tangeul/packet: F_GATE slot missing")
		}
		g, err := word.DecodeGate(body[pos])
		if err != nil {
			return EdgePacket{}, 0, err
		}
		p.Gate = &g
		pos++
	}

	if pos+3 != len(body) {
		return EdgePacket{}, 0, errors.New("tangeul/packet: Edge slot words do not match framed length")
	}
	s, err := word.DecodeSpan(body[pos:])
	if err != nil {
		return EdgePacket{}, 0, err
	}
	p.Span = s
	return p, total, nil
}
