//ff:func feature=engine type=engine control=sequence
//ff:what plan — returns the graph's evalPlan, building it on the first evaluation after the graph was last changed
package toulmin

// plan returns the evaluation plan of the graph as it stands. The plan is
// built once and kept until the graph changes (touch). Two evaluations that
// start at the same time may each build one; both are the same and one stays.
func (g *Graph) plan() *evalPlan {
	if p := g.compiled.Load(); p != nil {
		return p
	}
	p := &evalPlan{}
	p.static, p.err = newEvalContext(g.rules, g.defeats, g.supersedes, g.roles)
	g.compiled.Store(p)
	return p
}
