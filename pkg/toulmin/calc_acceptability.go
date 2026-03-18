//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what CalcAcceptability — h-Categoriser recursive verdict computation [-1, 1]
package toulmin

const maxDepth = 100

// CalcAcceptability computes the verdict for a node using Amgoud's
// weighted h-Categoriser. Returns [-1.0, +1.0].
// +1.0 = violation confirmed, 0.0 = undecided, -1.0 = fully rebutted.
func CalcAcceptability(nodeID string, graph *RuleGraph, depth int) float64 {
	if depth >= maxDepth {
		return 0.0
	}
	node := graph.Nodes[nodeID]
	attackerSum := 0.0
	for _, attackerID := range graph.Attackers(nodeID) {
		raw := (CalcAcceptability(attackerID, graph, depth+1) + 1.0) / 2.0
		attackerSum += raw
	}
	raw := node.Qualifier / (1.0 + attackerSum)
	return 2*raw - 1
}
