//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what BuildSubgraph — constructs defeats subgraph from activated rules
package toulmin

// BuildSubgraph creates a RuleGraph from activated rules.
// Nodes are created for each activated rule. Defeat edges are added
// unless the target node has Strict strength.
func BuildSubgraph(activated []RuleMeta) *RuleGraph {
	graph := &RuleGraph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}
	for _, r := range activated {
		graph.Nodes[r.Name] = &Node{
			Name:      r.Name,
			Qualifier: r.Qualifier,
			Strength:  r.Strength,
		}
	}
	for _, r := range activated {
		addDefeatEdges(graph, r)
	}
	return graph
}
