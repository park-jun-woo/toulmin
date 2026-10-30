//ff:func feature=engine type=engine control=sequence
//ff:what Attacks — declares that this rule attacks the target rule (optional force, default 1)
package toulmin

// Attacks declares a defeat edge: this rule attacks target.
// The optional force F in [0, 1] (default 1) scales the attack: a rebuttal
// (Counter or Rule attacker) adds F·(v+1)/2 to the target's denominator, and an
// exception (Except) multiplies the target's applicability by 1 − F·(v+1)/2.
// Panics if more than one force is given, or force is NaN or outside [0, 1].
func (r *Rule) Attacks(target *Rule, force ...float64) {
	f := attackForce(force)
	r.graph.defeats = append(r.graph.defeats, defeatEdge{
		from:  r.id,
		to:    target.id,
		force: f,
	})
}
