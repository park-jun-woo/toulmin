//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what EvaluateTrace — lazily evaluates rules by graph traversal and returns verdicts with per-warrant trace
package toulmin

// EvaluateTrace traverses the defeats graph from each warrant node,
// lazily executing rule funcs only when reached. Returns verdicts with
// per-warrant trace containing only relevant rules. State is reset per warrant.
// Returns an error if the defeat graph contains a cycle.
func (g *Graph) EvaluateTrace(claim any, ground any) ([]EvalResult, error) {
	ctx, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		return nil, err
	}
	var results []EvalResult
	for _, r := range g.rules {
		if !isWarrant(ctx.attackerSet, r.Strength, r.Name) {
			continue
		}
		ctx.reset()
		verdict := ctx.calcTrace(r.Name, claim, ground)
		if !ctx.active[r.Name] {
			continue
		}
		trace := make([]TraceEntry, len(ctx.trace))
		for i, te := range ctx.trace {
			te.Name = shortName(te.Name)
			trace[i] = te
		}
		results = append(results, EvalResult{Name: shortName(r.Name), Verdict: verdict, Evidence: ctx.evidence[r.Name], Trace: trace})
	}
	return results, nil
}
