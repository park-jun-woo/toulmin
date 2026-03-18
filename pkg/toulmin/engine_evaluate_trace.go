//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — runs all rules against a claim and returns verdicts with trace
package toulmin

// EvaluateTrace runs all registered rules against the claim/ground pair,
// builds a defeats subgraph from activated rules, and returns verdicts
// with trace for warrant nodes.
func (e *Engine) EvaluateTrace(claim any, ground any) []EvalResult {
	var trace []TraceEntry
	var activated []RuleMeta
	for _, r := range e.rules {
		result := r.Fn(claim, ground)
		trace = append(trace, TraceEntry{
			Name:      r.Name,
			Role:      inferRole(r),
			Activated: result,
			Qualifier: r.Qualifier,
		})
		if result {
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
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Trace: trace})
	}
	return results
}
