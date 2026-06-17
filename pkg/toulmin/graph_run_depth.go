//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what runDepth — recursively evaluates, fires handlers, and Runs Active nodes' sub-graphs
package toulmin

import "fmt"

// runMaxDepth caps recursive execution composition as a runtime backstop for cases
// static cycle detection cannot see (e.g. future dynamic Run wiring).
const runMaxDepth = 64

// runDepth evaluates the graph (full pass), builds an immutable RunView, then for each
// node in registration order: (a) fires its event handler (leaf side effect), and
// (b) if the node is Active and has a RunGraph, Runs that sub-graph with the same ctx
// at depth+1. ctx flows down (shared mutable); each level builds its own view. Sub-graph
// verdicts stay isolated — only errors propagate, wrapped as `run ... → ...`.
func (g *Graph) runDepth(ctx Context, opt EvalOption, depth int) ([]EvalResult, RunView, error) {
	if depth > runMaxDepth {
		return nil, nil, fmt.Errorf("toulmin: run depth exceeded %d (possible runaway composition)", runMaxDepth)
	}
	results, ec, err := g.evaluate(ctx, opt, true)
	if err != nil {
		return nil, nil, err
	}
	view := newRunView(g, ec)
	events := view.All()
	for i := range events {
		ne := events[i]
		meta := &g.rules[i]
		var herr error
		if h := selectHandler(meta, ne.Type); h != nil {
			herr = safeCallHandler(h, ctx, ne, view)
		}
		if herr != nil {
			return results, view, fmt.Errorf("handler %q (%v): %w", ne.Name, ne.Type, herr)
		}
		if ne.Type != Active || meta.RunGraph == nil {
			continue
		}
		if _, _, err := meta.RunGraph.runDepth(ctx, opt, depth+1); err != nil {
			return results, view, fmt.Errorf("run %q → %q: %w", ne.Name, meta.RunGraph.name, err)
		}
	}
	return results, view, nil
}
