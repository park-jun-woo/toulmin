//ff:func feature=engine type=engine control=sequence
//ff:what context — returns an evalContext for one evaluation: the plan's shared maps with empty per-evaluation maps
package toulmin

// context returns the plan's error when the graph has a cycle. A state from
// the pool comes back with its per-evaluation maps already cleared.
func (p *evalPlan) context() (*evalContext, error) {
	if p.err != nil {
		return nil, p.err
	}
	if v := p.pool.Get(); v != nil {
		return v.(*evalContext), nil
	}
	ec := *p.static
	ec.home = p
	ec.reset()
	return &ec, nil
}
