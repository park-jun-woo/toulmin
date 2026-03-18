//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — lazily evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts
// for warrant nodes. Funcs are cached across warrant evaluations.
func (b *GraphBuilder) Evaluate(claim any, ground any) []EvalResult {
	ctx := newEvalContext(b.rules, b.defeats, b.roles)
	var results []EvalResult
	for _, r := range b.rules {
		if !isWarrant(ctx.edges, r.Strength, r.Name) {
			continue
		}
		verdict := ctx.calc(r.Name, claim, ground, 0)
		if !ctx.active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Evidence: ctx.evidence[r.Name]})
	}
	return results
}
