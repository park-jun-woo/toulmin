//ff:func feature=engine type=engine control=sequence
//ff:what Run — rejects run cycles then recursively evaluates and fires each node's handler
package toulmin

// Run pre-evaluates the whole graph (full pass) then fires each Active node's run handler and,
// for Active nodes with a declared sub-graph, Runs that sub-graph (execution composition).
// At the top-level entry it rejects static run cycles once via detectRunCycle, then resolves
// options (preserving option errors such as Recursive), forces Trace and Duration off so the
// full pass runs on shared non-reset state, and delegates the recursive dispatch to runDepth.
// It returns each warrant's result plus the whole graph's flat trace (every node, registration order).
func (g *Graph) Run(ctx Context, opts ...EvalOption) ([]EvalResult, []TraceEntry, error) {
	if err := detectRunCycle(g); err != nil {
		return nil, nil, err
	}
	opt, err := resolveOption(opts)
	if err != nil {
		return nil, nil, err
	}
	opt.Trace = false
	opt.Duration = false
	return g.runDepth(ctx, opt, 0)
}
