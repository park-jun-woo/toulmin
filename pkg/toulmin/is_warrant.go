//ff:func feature=engine type=engine control=sequence
//ff:what isWarrant — determines if a rule is a warrant based on edges and strength
package toulmin

// isWarrant returns true if the rule is not an attacker and not a defeater.
func isWarrant(edges map[string][]string, strength Strength, name string) bool {
	if strength == Defeater {
		return false
	}
	for _, attackers := range edges {
		for _, aid := range attackers {
			if aid == name {
				return false
			}
		}
	}
	return true
}
