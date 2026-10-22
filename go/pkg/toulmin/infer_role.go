//ff:func feature=engine type=engine control=sequence
//ff:what inferRole — infers rule role from strength and attacker set
package toulmin

// inferRole determines the role of a rule from its strength and
// whether it appears as an attacker in any defeat edge.
func inferRole(strMap map[string]Strength, attackerSet map[string]bool, id string) string {
	if strMap[id] == Defeater {
		return "except"
	}
	if attackerSet[id] {
		return "counter"
	}
	return "rule"
}
