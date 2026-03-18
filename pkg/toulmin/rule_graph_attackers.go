//ff:func feature=engine type=engine control=sequence
//ff:what Attackers — returns attacker node IDs for a given node
package toulmin

// Attackers returns the IDs of nodes that attack the given node.
func (g *RuleGraph) Attackers(nodeID string) []string {
	return g.Edges[nodeID]
}
