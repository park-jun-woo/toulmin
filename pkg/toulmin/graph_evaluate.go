//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Evaluate — evaluates rules by graph traversal and returns verdicts
package toulmin

import "fmt"

// Evaluate traverses the defeats graph from each rule node and returns verdicts.
// Default method is Matrix. Use EvalOption to enable Trace, Duration, or Recursive method.
// Returns an error if the defeat graph contains a cycle or a rule panics.
func (g *Graph) Evaluate(ctx Context, opts ...EvalOption) ([]EvalResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("toulmin: ctx must not be nil")
	}
	opt, err := resolveOption(opts)
	if err != nil {
		return nil, err
	}
	ec, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		return nil, err
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
			return nil, ec.err
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
	return results, nil
}
