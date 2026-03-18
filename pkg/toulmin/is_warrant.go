//ff:func feature=engine type=engine control=sequence
//ff:what isWarrant — determines if a rule is a warrant based on edges and strength
package toulmin

// isWarrant returns true if the rule is not an attacker and not a defeater.
func isWarrant(attackerSet map[string]bool, strength Strength, name string) bool {
	if strength == Defeater {
		return false
	}
	return !attackerSet[name]
}
