//ff:func feature=tangl type=util control=sequence
//ff:what EncodeNode — encode a NodePacket (frame + TID + type + header + flag-ordered slots + SPAN)
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// EncodeNode encodes p as a Node packet (§5.1). The body is TID, NodeType, the
// header word (Role in bits 15-14, Level in bits 10-9, Claim object kind in bits
// 8-7, Flags in bits 6-0), then the conditional slots emitted in flag order (F_FUNC, F_QUAL,
// F_CHECK, F_CLAIM, F_VERDICT, F_WITH; F_STRIKE has no slot), then a trailing
// SPAN. A flag whose field is absent is a hard error so the stream never emits a
// truncated slot.
func EncodeNode(p NodePacket) ([]word.Word, error) {
	var objKind ObjKind
	if p.Claim != nil {
		objKind = p.Claim.Kind
	}
	header := word.Word(uint16(p.Role)<<14 | uint16(p.Level)<<9 | uint16(objKind)<<7 | uint16(p.Flags))
	body := []word.Word{word.Word(p.TID), word.Word(p.Type), header}

	if p.Flags&codebook.NodeFlagFunc != 0 {
		if p.Func == nil {
			return nil, errors.New("tangeul/packet: F_FUNC set but Func is nil")
		}
		body = append(body, encodeFuncRef(*p.Func)...)
	}
	if p.Flags&codebook.NodeFlagQual != 0 {
		if p.Qual == nil {
			return nil, errors.New("tangeul/packet: F_QUAL set but Qual is nil")
		}
		body = append(body, word.Word(*p.Qual))
	}
	if p.Flags&codebook.NodeFlagCheck != 0 {
		if p.Check == nil {
			return nil, errors.New("tangeul/packet: F_CHECK set but Check is nil")
		}
		body = append(body, word.Word(*p.Check))
	}
	if p.Flags&codebook.NodeFlagClaim != 0 {
		if p.Claim == nil {
			return nil, errors.New("tangeul/packet: F_CLAIM set but Claim is nil")
		}
		cw, err := encodeClaim(*p.Claim)
		if err != nil {
			return nil, err
		}
		body = append(body, cw...)
	}
	if p.Flags&codebook.NodeFlagVerdict != 0 {
		if p.Verdict == nil {
			return nil, errors.New("tangeul/packet: F_VERDICT set but Verdict is nil")
		}
		body = append(body, word.Word(*p.Verdict))
	}
	if p.Flags&codebook.NodeFlagWith != 0 {
		body = append(body, word.Word(len(p.With)))
		for _, t := range p.With {
			body = append(body, word.Word(t))
		}
	}

	s := word.EncodeSpan(p.Span)
	body = append(body, s[0], s[1], s[2])
	return append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: uint32(len(body))}), body...), nil
}
