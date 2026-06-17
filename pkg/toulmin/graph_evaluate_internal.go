//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what evaluate — shared graph evaluation for Evaluate (lazy) and Run (full)
package toulmin

// evaluate traverses the defeats graph and returns warrant results with the evalContext.
// When full is true, every node is calc'd so inactive nodes also have their state filled.
// Returns an error if the defeat graph contains a cycle or a rule panics.
func (g *Graph) evaluate(ctx Context, opt EvalOption, full bool) ([]EvalResult, *evalContext, error) {
	ec, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		return nil, nil, err
	}
	var results []EvalResult
	for _, r := range g.rules {
		if !isWarrant(ec.attackerSet, r.Strength, r.Name) {
			continue
		}
		if opt.Trace {
			ec.reset()
		}
		verdict := ec.evalRule(r.Name, ctx, opt)
		if ec.err != nil {
			return nil, nil, ec.err
		}
		if !ec.active[r.Name] {
			continue
		}
		result := EvalResult{Name: shortName(r.Name), Verdict: verdict, Evidence: ec.evidence[r.Name]}
		if opt.Trace {
			result.Trace = collectTrace(ec.trace)
		}
		results = append(results, result)
	}
	if full {
		ec.fillAll(g.rules, ctx)
		if ec.err != nil {
			return nil, ec, ec.err
		}
	}
	return results, ec, nil
}
