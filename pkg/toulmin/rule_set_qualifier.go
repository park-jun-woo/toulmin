//ff:func feature=engine type=engine control=sequence
//ff:what Qualifier — sets the qualifier weight for this rule
package toulmin

// Qualifier sets the rule's confidence weight [0.0, 1.0] and returns the rule for chaining.
func (r *Rule) Qualifier(q float64) *Rule {
	if q < 0.0 || q > 1.0 {
		panic("toulmin: qualifier must be between 0.0 and 1.0")
	}
	r.graph.rules[r.idx].Qualifier = q
	return r
}
