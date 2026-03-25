//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — lazily evaluates rules by graph traversal and returns verdicts
package toulmin

// Evaluate traverses the defeats graph from each rule node,
// lazily executing rule funcs only when reached. Returns verdicts
// for rule nodes. Use EvalOption to enable Trace, Duration, or Recursive method.
// Returns an error if the defeat graph contains a cycle or a rule panics.
func (e *Engine) Evaluate(ctx Context, opts ...EvalOption) ([]EvalResult, error) {
	opt := resolveOption(opts)
	ec, err := newEvalContext(e.rules, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range e.rules {
		if !isWarrant(ec.attackerSet, r.Strength, r.Name) {
			continue
		}
		if opt.Trace {
			ec.reset()
		}
		verdict := ec.evalRule(r.Name, ctx, opt)
		if ec.err != nil {
			return nil, ec.err
		}
		if !ec.active[r.Name] {
			continue
		}
		result := EvalResult{Name: r.Name, Verdict: verdict, Evidence: ec.evidence[r.Name]}
		if opt.Trace {
			result.Trace = ec.trace
		}
		results = append(results, result)
	}
	return results, nil
}
