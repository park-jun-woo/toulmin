//ff:func feature=tangl type=util control=selection
//ff:what appendDecoded — decode a body record into the node or edge list, register its Edge Id, return its span
package stream

import (
	"errors"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// appendDecoded decodes one statement-body record (NODE or LINK), registers its
// Edge Id, appends it to the matching list and returns its SPAN.
func appendDecoded(d *dec, nodes *[]toulmin.NodePacket, edges *[]toulmin.EdgePacket, rp rawPacket) (toulmin.Span, error) {
	switch rp.Class {
	case classNode:
		np, err := toulmin.DecodeNode(rp.Rec, d.wire)
		if err == nil {
			err = claimID(d, np.TID)
		}
		*nodes = append(*nodes, np)
		return np.Span, err
	case classLink:
		ep, err := toulmin.DecodeEdge(rp.Rec, d.wire)
		if err == nil {
			err = claimID(d, ep.ID)
		}
		*edges = append(*edges, ep)
		return ep.Span, err
	default:
		return toulmin.Span{}, errors.New("tangeul/stream: unexpected record class in statement body")
	}
}
