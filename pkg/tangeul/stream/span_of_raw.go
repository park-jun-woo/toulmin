//ff:func feature=tangl type=util control=selection
//ff:what spanOfRaw — read the SPAN anchor of a framed Node or Edge packet
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// spanOfRaw decodes just enough of a Node or Edge rawPacket to return its
// trailing SPAN anchor, which statement grouping uses to detect the boundary
// between statements (all packets of one statement share the block span). A
// non-span-bearing class is a hard error — callers only ask for Node/Edge.
func spanOfRaw(rp rawPacket) (word.Span, error) {
	switch rp.Class {
	case word.ClassNode:
		p, _, err := packet.DecodeNode(rp.Words)
		return p.Span, err
	case word.ClassEdge:
		p, _, err := packet.DecodeEdge(rp.Words)
		return p.Span, err
	default:
		return word.Span{}, errors.New("tangeul/stream: spanOfRaw on a non-span packet")
	}
}
