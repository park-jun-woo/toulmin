//ff:func feature=tangl type=util control=selection
//ff:what EncodeEdge — encode an EdgePacket (frame + PropCode + src + dst + flags + GATE16 + SPAN)
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// EncodeEdge encodes p as an Edge packet (§6.1). The dst slot layout is decided
// by the PropCode via codebook.PropDstKind: a FuncRef dst (do/undo) emits
// DstFunc as two words, every TID or Text dst emits DstTID as one word. A GATE16
// follows the flag word iff F_GATE is set, then the trailing SPAN. Provenance
// (DstBlock) is not encoded in v0.1.
func EncodeEdge(p EdgePacket) ([]word.Word, error) {
	kind, err := codebook.PropDstKind(p.Prop)
	if err != nil {
		return nil, err
	}
	body := []word.Word{word.Word(p.Prop), word.Word(p.Src)}
	switch kind {
	case codebook.DstFuncRef:
		if p.DstFunc == nil {
			return nil, errors.New("tangeul/packet: FuncRef dst PropCode but DstFunc is nil")
		}
		body = append(body, encodeFuncRef(*p.DstFunc)...)
	case codebook.DstTID, codebook.DstText:
		body = append(body, word.Word(p.DstTID))
	default:
		return nil, errors.New("tangeul/packet: provenance/Block dst is not encoded in v0.1")
	}

	body = append(body, word.Word(p.Flags))
	if p.Flags&codebook.EdgeFlagGate != 0 {
		if p.Gate == nil {
			return nil, errors.New("tangeul/packet: F_GATE set but Gate is nil")
		}
		g, err := word.EncodeGate(*p.Gate)
		if err != nil {
			return nil, err
		}
		body = append(body, g)
	}

	s := word.EncodeSpan(p.Span)
	body = append(body, s[0], s[1], s[2])
	return append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: uint32(len(body))}), body...), nil
}
