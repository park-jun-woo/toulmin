//ff:func feature=engine type=engine control=sequence
//ff:what RunOn — registers the node's run handler (fires when Active)
package toulmin

// RunOn sets the handler fired when this node is Active and returns the rule for chaining.
func (r *Rule) RunOn(h NodeHandler) *Rule {
	r.graph.rules[r.idx].RunOn = h
	return r
}
