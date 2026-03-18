//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — runs graph builder rules against claim/ground and returns verdicts with trace
package toulmin

// EvaluateTrace runs all rules in the graph against the claim/ground pair,
// builds a defeats subgraph from activated rules with builder-defined edges,
// and returns verdicts with trace for non-attacking nodes.
func (b *GraphBuilder) EvaluateTrace(claim any, ground any) []EvalResult {
	var trace []TraceEntry
	var activated []RuleMeta
	for _, r := range b.rules {
		result := r.Fn(claim, ground)
		trace = append(trace, TraceEntry{
			Name:      r.Name,
			Role:      b.roles[r.Name],
			Activated: result,
			Qualifier: r.Qualifier,
		})
		if result {
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
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Trace: trace})
	}
	return results
}
