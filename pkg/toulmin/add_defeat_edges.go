//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what addDefeatEdges — adds defeat edges from an attacker to its targets
package toulmin

// addDefeatEdges adds edges from the attacker to each of its defeat targets.
// Edges targeting Strict nodes are rejected.
func addDefeatEdges(graph *RuleGraph, attacker RuleMeta) {
	for _, target := range attacker.Defeats {
		targetNode, ok := graph.Nodes[target]
		if !ok {
			continue
		}
		if targetNode.Strength == Strict {
			continue
		}
		graph.Edges[target] = append(graph.Edges[target], attacker.Name)
	}
}
