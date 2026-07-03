//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what splitNodes — separate auxiliary Func nodes from the primary node of a statement
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// splitNodes partitions a statement's decoded nodes into materialized Func nodes
// (which carry function references, not AST elements) and primary nodes (at most
// one, the element the statement declares).
func splitNodes(nodes []packet.NodePacket) (funcNodes, primaries []packet.NodePacket) {
	for _, n := range nodes {
		if n.Type == codebook.NodeFunc {
			funcNodes = append(funcNodes, n)
		} else {
			primaries = append(primaries, n)
		}
	}
	return funcNodes, primaries
}
