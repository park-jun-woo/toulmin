//ff:func feature=engine type=engine control=sequence
//ff:what Qualifier — sets the qualifier weight for this rule
package toulmin

// Qualifier sets the rule's confidence weight and returns the rule for chaining.
func (r *Rule) Qualifier(q float64) *Rule {
	r.graph.rules[r.idx].Qualifier = q
	return r
}
