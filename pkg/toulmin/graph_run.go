//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what Run — evaluates the graph then fires each node's event handler with a shared RunView
package toulmin

import "fmt"

// Run pre-evaluates the whole graph (full pass) then fires each node's event handler.
// Firing order is g.rules registration order (deterministic). Before any handler fires,
// it builds one immutable RunView snapshot of every node's final event and shares it with
// all handlers. If a handler returns an error or panics, Run stops immediately and returns
// the RunView built before dispatch. Run forces Trace and Duration off so the full pass
// runs on shared, non-reset state.
func (g *Graph) Run(ctx Context, opts ...EvalOption) ([]EvalResult, RunView, error) {
	opt, err := resolveOption(opts)
	if err != nil {
		return nil, nil, err
	}
	opt.Trace = false
	opt.Duration = false
	results, ec, err := g.evaluate(ctx, opt, true)
	if err != nil {
		return nil, nil, err
	}
	view := newRunView(g, ec)
	events := view.order
	for i := range events {
		ne := events[i]
		meta := &g.rules[i]
		h := selectHandler(meta, ne.Type)
		if h == nil {
			continue
		}
		if err := safeCallHandler(h, ctx, ne, view); err != nil {
			return results, view, fmt.Errorf("handler %q (%v): %w", ne.Name, ne.Type, err)
		}
	}
	return results, view, nil
}
