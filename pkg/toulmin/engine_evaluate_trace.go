//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — lazily evaluates rules by graph traversal and returns verdicts with per-warrant trace
package toulmin

// EvaluateTrace traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts with
// per-warrant trace containing only relevant rules. State is reset per warrant.
// Returns an error if the defeat graph contains a cycle.
func (e *Engine) EvaluateTrace(claim any, ground any) ([]EvalResult, error) {
	ctx, err := newEvalContext(e.rules, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range e.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		ctx.reset()
		verdict := ctx.calcTrace(r.Name, claim, ground)
		if !ctx.active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Evidence: ctx.evidence[r.Name], Trace: ctx.trace})
	}
	return results, nil
}
