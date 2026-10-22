//ff:func feature=engine type=engine control=sequence
//ff:what Run — declares that this node Runs the sub-graph g when it is Active
package toulmin

// Run declares an execution edge: when this node is Active, g is Run (ctx flows down).
// It is the execution-composition counterpart of Attacks (which is a defeat edge).
// g must not be nil — a nil sub-graph is a registration error and panics.
func (r *Rule) Run(g *Graph) *Rule {
	if g == nil {
		panic("toulmin: Run requires a non-nil sub-graph")
	}
	r.graph.rules[r.idx].RunGraph = g
	return r
}
