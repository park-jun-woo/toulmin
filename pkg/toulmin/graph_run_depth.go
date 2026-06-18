//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what runDepth — recursively evaluates, fires RunOn handlers, and Runs Active nodes' sub-graphs
package toulmin

import "fmt"

// runMaxDepth caps recursive execution composition as a runtime backstop for cases
// static cycle detection cannot see (e.g. future dynamic Run wiring).
const runMaxDepth = 64

// runDepth evaluates the graph (full pass), assembles a flat trace of every node in
// registration order from the shared evalContext, then for each Active node (Activated &&
// Verdict>0) in registration order: (a) fires its RunOn handler with self and the whole trace,
// then (b) if the node has a RunGraph, Runs that sub-graph with the same ctx at depth+1.
// ctx flows down (shared mutable); each level builds its own trace. Sub-graph verdicts stay
// isolated — only errors propagate, wrapped as `run ... → ...`.
func (g *Graph) runDepth(ctx Context, opt EvalOption, depth int) ([]EvalResult, []TraceEntry, error) {
	if depth > runMaxDepth {
		return nil, nil, fmt.Errorf("toulmin: run depth exceeded %d (possible runaway composition)", runMaxDepth)
	}
	results, ec, err := g.evaluate(ctx, opt, true)
	if err != nil {
		return nil, nil, err
	}
	trace := make([]TraceEntry, len(g.rules))
	for i := range g.rules {
		name := g.rules[i].Name
		role := g.roles[name]
		if role == "" {
			role = inferRole(ec.strMap, ec.attackerSet, name)
		}
		trace[i] = TraceEntry{
			Name:      shortName(name),
			Role:      role,
			Activated: ec.active[name],
			Qualifier: ec.qualMap[name],
			Verdict:   ec.verdictCache[name],
			Evidence:  ec.evidence[name],
			Ground:    ctx,
			Specs:     ec.specsMap[name],
		}
	}
	for i := range g.rules {
		meta := &g.rules[i]
		self := trace[i]
		if !(self.Activated && self.Verdict > 0) {
			continue
		}
		if meta.RunOn != nil {
			if herr := safeCallHandler(meta.RunOn, ctx, self, trace); herr != nil {
				return results, trace, fmt.Errorf("runOn %q: %w", self.Name, herr)
			}
		}
		if meta.RunGraph != nil {
			if _, _, err := meta.RunGraph.runDepth(ctx, opt, depth+1); err != nil {
				return results, trace, fmt.Errorf("run %q → %q: %w", self.Name, meta.RunGraph.name, err)
			}
		}
	}
	return results, trace, nil
}
