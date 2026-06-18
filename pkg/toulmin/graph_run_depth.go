//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what runDepth — recursively evaluates, fires RunOn handlers, and Runs Active nodes' sub-graphs
package toulmin

import "fmt"

// runMaxDepth caps recursive execution composition as a runtime backstop for cases
// static cycle detection cannot see (e.g. future dynamic Run wiring).
const runMaxDepth = 64

// runDepth evaluates the graph (full pass), assembles every node's entry in registration
// order from the shared evalContext, wraps them in a Trace, then for each Active node
// (Activated && Verdict>0) in registration order: (a) fires its RunOn handler with the whole
// Run's Trace, then (b) if the node has a RunGraph, Runs that sub-graph with the same ctx at
// depth+1. ctx flows down (shared mutable); each level builds its own Trace. Sub-graph verdicts
// stay isolated — only errors propagate, wrapped as `run ... → ...`.
func (g *Graph) runDepth(ctx Context, opt EvalOption, depth int) ([]EvalResult, Trace, error) {
	if depth > runMaxDepth {
		return nil, Trace{}, fmt.Errorf("toulmin: run depth exceeded %d (possible runaway composition)", runMaxDepth)
	}
	results, ec, err := g.evaluate(ctx, opt, true)
	if err != nil {
		return nil, Trace{}, err
	}
	entries := g.buildTraceEntries(ec, ctx)
	tr := Trace{nodes: entries, ctx: ctx}
	for i := range g.rules {
		self := entries[i]
		if !(self.Activated && self.Verdict > 0) {
			continue
		}
		if rerr := runNode(&g.rules[i], self, tr, ctx, opt, depth); rerr != nil {
			return results, tr, rerr
		}
	}
	return results, tr, nil
}
