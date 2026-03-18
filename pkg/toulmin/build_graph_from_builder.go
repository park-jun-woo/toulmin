//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what buildGraphFromBuilder — creates RuleGraph from activated rules and builder defeat edges
package toulmin

// buildGraphFromBuilder creates a RuleGraph using builder-defined defeat edges
// instead of RuleMeta.Defeats.
func buildGraphFromBuilder(activated []RuleMeta, defeats []defeatEdge) *RuleGraph {
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
	for _, d := range defeats {
		targetNode, ok := graph.Nodes[d.to]
		if !ok {
			continue
		}
		if targetNode.Strength == Strict {
			continue
		}
		if _, ok := graph.Nodes[d.from]; !ok {
			continue
		}
		graph.Edges[d.to] = append(graph.Edges[d.to], d.from)
	}
	return graph
}
