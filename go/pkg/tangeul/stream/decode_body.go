//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what decodeBody — decode a statement's NODE/LINK records and return them with the shared span
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// decodeBody decodes the span-bearing body records of one statement into nodes
// and edges, returning the span they share. A statement with no body record has
// no span and is an error.
func decodeBody(d *dec, body []rawPacket) ([]toulmin.NodePacket, []toulmin.EdgePacket, toulmin.Span, error) {
	var nodes []toulmin.NodePacket
	var edges []toulmin.EdgePacket
	var span toulmin.Span
	for i, rp := range body {
		s, err := appendDecoded(d, &nodes, &edges, rp)
		if err != nil {
			return nil, nil, toulmin.Span{}, err
		}
		if i == 0 {
			span = s
		}
	}
	if len(body) == 0 {
		return nil, nil, toulmin.Span{}, errors.New("tangeul/stream: statement has no span-bearing record")
	}
	return nodes, edges, span, nil
}
