//ff:func feature=engine type=engine control=sequence
//ff:what OnInactive — registers the inactive event handler for this rule
package toulmin

// OnInactive sets the handler fired when this node is Inactive and returns the rule for chaining.
func (r *Rule) OnInactive(h NodeHandler) *Rule {
	r.graph.rules[r.idx].OnInactive = h
	return r
}
