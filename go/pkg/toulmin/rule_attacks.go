//ff:func feature=engine type=engine control=sequence
//ff:what Attacks — declares that this rule attacks the target rule
package toulmin

// Attacks declares a defeat edge: this rule attacks target.
func (r *Rule) Attacks(target *Rule) {
	r.graph.defeats = append(r.graph.defeats, defeatEdge{
		from: r.id,
		to:   target.id,
	})
}
