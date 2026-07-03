//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what splitEdges — separate label edges from relation edges of a statement
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// splitEdges partitions a statement's decoded edges into label edges (node ->
// Text spellings) and relation edges (requires/has/means/do/run/…), the two
// carrying different meaning to the builders.
func splitEdges(edges []packet.EdgePacket) (labelEdges, rel []packet.EdgePacket) {
	for _, e := range edges {
		if e.Prop == codebook.PropLabel {
			labelEdges = append(labelEdges, e)
		} else {
			rel = append(rel, e)
		}
	}
	return labelEdges, rel
}
