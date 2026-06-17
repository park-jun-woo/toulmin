//ff:func feature=engine type=engine control=sequence
//ff:what OnActive — registers the active event handler for this rule
package toulmin

// OnActive sets the handler fired when this node is Active and returns the rule for chaining.
func (r *Rule) OnActive(h NodeHandler) *Rule {
	r.graph.rules[r.idx].OnActive = h
	return r
}
