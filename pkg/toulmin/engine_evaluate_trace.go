//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — lazily evaluates rules by graph traversal and returns verdicts with per-warrant trace
package toulmin

// EvaluateTrace traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts with
// per-warrant trace containing only relevant rules. State is reset per warrant.
func (e *Engine) EvaluateTrace(claim any, ground any) []EvalResult {
	ctx := newEvalContext(e.rules, nil, nil)
	var results []EvalResult
	for _, r := range e.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		ctx.reset()
		verdict := ctx.calcTrace(r.Name, claim, ground, 0)
		if !ctx.active[r.Name] {
			continue
		}
		results = append(results, EvalResult{Name: r.Name, Verdict: verdict, Evidence: ctx.evidence[r.Name], Trace: ctx.trace})
	}
	return results
}
