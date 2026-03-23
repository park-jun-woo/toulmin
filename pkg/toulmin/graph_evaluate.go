//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node and returns verdicts.
// Default method is Matrix. Pass Recursive for proven h-Categoriser traversal.
// Returns an error if the defeat graph contains a cycle.
func (g *Graph) Evaluate(claim any, ground any, opts ...EvalOption) ([]EvalResult, error) {
	ctx, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range g.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		verdict := ctx.calc(r.Name, claim, ground)
		if !ctx.active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: shortName(r.Name), Verdict: verdict, Evidence: ctx.evidence[r.Name]})
	}
	return results, nil
}
