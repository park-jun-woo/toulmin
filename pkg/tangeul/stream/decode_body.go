//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what decodeBody — decode a statement body's Node/Edge packets and its shared span
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// decodeBody decodes each span-bearing packet of a statement body into typed Node
// and Edge slices and returns their shared block span (all packets of a statement
// carry the same span). An empty body has no span and is a hard error.
func decodeBody(body []rawPacket) ([]packet.NodePacket, []packet.EdgePacket, word.Span, error) {
	var nodes []packet.NodePacket
	var edges []packet.EdgePacket
	var span word.Span
	set := false
	for _, rp := range body {
		s, err := appendDecoded(&nodes, &edges, rp)
		if err != nil {
			return nil, nil, word.Span{}, err
		}
		if !set {
			span, set = s, true
		}
	}
	if !set {
		return nil, nil, word.Span{}, errors.New("tangeul/stream: statement has no span-bearing packet")
	}
	return nodes, edges, span, nil
}
