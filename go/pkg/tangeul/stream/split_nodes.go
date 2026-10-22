//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what splitNodes — separate auxiliary Func nodes from the primary node of a statement
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// splitNodes partitions a statement's decoded nodes into materialized Func nodes
// (which carry function references, not AST elements) and primary nodes (at most
// one, the element the statement declares).
func splitNodes(nodes []toulmin.NodePacket) (funcNodes, primaries []toulmin.NodePacket) {
	for _, n := range nodes {
		if n.Type == toulmin.NodeFunc {
			funcNodes = append(funcNodes, n)
		} else {
			primaries = append(primaries, n)
		}
	}
	return funcNodes, primaries
}
