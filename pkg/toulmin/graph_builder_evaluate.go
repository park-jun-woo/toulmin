//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — runs graph builder rules against claim/ground and returns verdicts
package toulmin

// Evaluate runs all rules in the graph against the claim/ground pair,
// builds a defeats subgraph from activated rules with builder-defined edges,
// and returns verdicts for non-attacking nodes.
func (b *GraphBuilder) Evaluate(claim any, ground any) []EvalResult {
	var activated []RuleMeta
	for _, r := range b.rules {
		if r.Fn(claim, ground) {
			activated = append(activated, r)
		}
	}
	graph := buildGraphFromBuilder(activated, b.defeats)
	var results []EvalResult
	attackers := collectAttackers(b.defeats)
	for _, r := range activated {
		if attackers[r.Name] {
			continue
		}
		if r.Strength == Defeater {
			continue
		}
		verdict := CalcAcceptability(r.Name, graph, 0)
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict})
	}
	return results
}
