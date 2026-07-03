//ff:func feature=tangl type=util control=selection
//ff:what appendDecoded — decode one body packet (Node or Edge) and append it, returning its span
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// appendDecoded decodes a single body rawPacket as a Node or Edge, appends it to
// the matching slice, and returns its SPAN anchor. A Payload or other class in
// the body is a hard error — the body holds only span-bearing packets.
func appendDecoded(nodes *[]packet.NodePacket, edges *[]packet.EdgePacket, rp rawPacket) (word.Span, error) {
	switch rp.Class {
	case word.ClassNode:
		np, _, err := packet.DecodeNode(rp.Words)
		if err != nil {
			return word.Span{}, err
		}
		*nodes = append(*nodes, np)
		return np.Span, nil
	case word.ClassEdge:
		ep, _, err := packet.DecodeEdge(rp.Words)
		if err != nil {
			return word.Span{}, err
		}
		*edges = append(*edges, ep)
		return ep.Span, nil
	default:
		return word.Span{}, errors.New("tangeul/stream: unexpected packet class in statement body")
	}
}
