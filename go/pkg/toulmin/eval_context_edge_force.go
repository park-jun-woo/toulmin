//ff:func feature=engine type=engine control=sequence
//ff:what edgeForce — returns the force of the i-th attack edge on id (default 1)
package toulmin

// edgeForce returns the force F of edges[id][i]. Edges derived from
// RuleMeta.Defeats carry no stored force and default to 1.
func (ec *evalContext) edgeForce(id string, i int) float64 {
	fs := ec.forces[id]
	if i < len(fs) {
		return fs[i]
	}
	return 1.0
}
