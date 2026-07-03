//ff:func feature=tangl type=util control=sequence
//ff:what DecodeNode — decode a Node packet, verifying the framed length, returning words consumed
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// DecodeNode decodes a Node packet from the head of ws, returning the packet and
// the total words consumed (frame header + body). The header word yields Role,
// the Claim object kind and Flags; slots are read in flag order and the trailing
// SPAN closes the body. A body whose slots do not consume exactly Len words
// (excluding the SPAN) is a hard error (self-contained framing check, §design 2).
func DecodeNode(ws []word.Word) (NodePacket, int, error) {
	f, head, err := word.DecodeFrame(ws)
	if err != nil {
		return NodePacket{}, 0, err
	}
	if f.Class != word.ClassNode {
		return NodePacket{}, 0, fmt.Errorf("tangeul/packet: expected Node class, got %d", uint8(f.Class))
	}
	total := head + int(f.Len)
	if len(ws) < total {
		return NodePacket{}, 0, errors.New("tangeul/packet: Node body shorter than framed length")
	}
	body := ws[head:total]
	if len(body) < 3 {
		return NodePacket{}, 0, errors.New("tangeul/packet: Node body requires TID, NodeType and header")
	}
	header := uint16(body[2])
	p := NodePacket{
		TID:   TID(body[0]),
		Type:  codebook.NodeType(body[1]),
		Role:  codebook.RoleCode(header >> 14),
		Flags: codebook.NodeFlags(header & 0x7F),
	}
	objKind := ObjKind((header >> 7) & 0x3)
	pos := 3

	if p.Flags&codebook.NodeFlagFunc != 0 {
		fr, err := decodeFuncRef(body[pos:])
		if err != nil {
			return NodePacket{}, 0, err
		}
		p.Func = &fr
		pos += 2
	}
	if p.Flags&codebook.NodeFlagQual != 0 {
		if len(body) < pos+1 {
			return NodePacket{}, 0, errors.New("tangeul/packet: F_QUAL slot missing")
		}
		v := word.VFix(body[pos])
		p.Qual = &v
		pos++
	}
	if p.Flags&codebook.NodeFlagCheck != 0 {
		if len(body) < pos+1 {
			return NodePacket{}, 0, errors.New("tangeul/packet: F_CHECK slot missing")
		}
		t := TID(body[pos])
		p.Check = &t
		pos++
	}
	if p.Flags&codebook.NodeFlagClaim != 0 {
		c, n, err := decodeClaim(body[pos:], objKind)
		if err != nil {
			return NodePacket{}, 0, err
		}
		p.Claim = &c
		pos += n
	}
	if p.Flags&codebook.NodeFlagVerdict != 0 {
		if len(body) < pos+1 {
			return NodePacket{}, 0, errors.New("tangeul/packet: F_VERDICT slot missing")
		}
		v := word.VFix(body[pos])
		p.Verdict = &v
		pos++
	}
	if p.Flags&codebook.NodeFlagWith != 0 {
		if len(body) < pos+1 {
			return NodePacket{}, 0, errors.New("tangeul/packet: F_WITH count missing")
		}
		count := int(body[pos])
		pos++
		if len(body) < pos+count {
			return NodePacket{}, 0, errors.New("tangeul/packet: F_WITH TID list truncated")
		}
		p.With = make([]TID, count)
		for i := 0; i < count; i++ {
			p.With[i] = TID(body[pos+i])
		}
		pos += count
	}

	if pos+3 != len(body) {
		return NodePacket{}, 0, errors.New("tangeul/packet: Node slot words do not match framed length")
	}
	s, err := word.DecodeSpan(body[pos:])
	if err != nil {
		return NodePacket{}, 0, err
	}
	p.Span = s
	return p, total, nil
}
