//ff:func feature=engine type=engine control=sequence
//ff:what touch — drops the graph's evalPlan after a change to its rules or edges
package toulmin

// touch must be called by everything that changes what newEvalContext reads:
// the rules (name, function, qualifier, strength, specs, level), the roles,
// the defeat edges and the supersede edges.
func (g *Graph) touch() {
	g.compiled.Store(nil)
}
