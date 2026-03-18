//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — runs all rules against a claim and returns verdicts
package toulmin

// Evaluate runs all registered rules against the claim/ground pair,
// builds a defeats subgraph from activated rules, and returns verdicts
// for warrant nodes (rules with empty Defeats and non-Defeater strength).
func (e *Engine) Evaluate(claim any, ground any) []EvalResult {
	var activated []RuleMeta
	for _, r := range e.rules {
		if r.Fn(claim, ground) {
			activated = append(activated, r)
		}
	}
	graph := BuildSubgraph(activated)
	var results []EvalResult
	for _, r := range activated {
		if len(r.Defeats) > 0 {
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
