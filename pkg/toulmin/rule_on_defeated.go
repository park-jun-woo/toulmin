//ff:func feature=engine type=engine control=sequence
//ff:what OnDefeated — registers the defeated event handler for this rule
package toulmin

// OnDefeated sets the handler fired when this node is Defeated and returns the rule for chaining.
func (r *Rule) OnDefeated(h NodeHandler) *Rule {
	r.graph.rules[r.idx].OnDefeated = h
	return r
}
