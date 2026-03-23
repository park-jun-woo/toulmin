//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — lazily evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts
// for warrant nodes. Use EvalOption to enable Trace, Duration, or Recursive method.
// Returns an error if the defeat graph contains a cycle.
func (e *Engine) Evaluate(claim any, ground any, opts ...EvalOption) ([]EvalResult, error) {
	opt := resolveOption(opts)
	ctx, err := newEvalContext(e.rules, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range e.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		if opt.Trace {
			ctx.reset()
		}
		verdict := ctx.evalRule(r.Name, claim, ground, opt)
		if !ctx.active[r.Name] {
			continue
		}
		result := EvalResult{Name: r.Name, Verdict: verdict, Evidence: ctx.evidence[r.Name]}
		if opt.Trace {
			result.Trace = ctx.trace
		}
		results = append(results, result)
	}
	return results, nil
}
