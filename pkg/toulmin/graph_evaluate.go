//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each warrant node and returns verdicts.
// Default method is Matrix. Use EvalOption to enable Trace, Duration, or Recursive method.
// Returns an error if the defeat graph contains a cycle.
func (g *Graph) Evaluate(claim any, ground any, opts ...EvalOption) ([]EvalResult, error) {
	opt := resolveOption(opts)
	ctx, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range g.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		if opt.Trace {
			ctx.reset()
		}
		verdict := ctx.evalRule(r.Name, claim, ground, opt)
		if ctx.err != nil {
			return nil, ctx.err
		}
		if !ctx.active[r.Name] {
			continue
		}
		result := EvalResult{Name: shortName(r.Name), Verdict: verdict, Evidence: ctx.evidence[r.Name]}
		if opt.Trace {
			result.Trace = collectTrace(ctx.trace)
		}
		results = append(results, result)
	}
	return results, nil
}
