//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — lazily evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts
// for warrant nodes. Funcs are cached across warrant evaluations.
// Returns an error if the defeat graph contains a cycle.
func (e *Engine) Evaluate(claim any, ground any, opts ...EvalOption) ([]EvalResult, error) {
	ctx, err := newEvalContext(e.rules, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range e.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		verdict := ctx.calc(r.Name, claim, ground)
		if !ctx.active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Evidence: ctx.evidence[r.Name]})
	}
	return results, nil
}
