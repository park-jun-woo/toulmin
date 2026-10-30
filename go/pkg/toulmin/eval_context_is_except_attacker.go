//ff:func feature=engine type=engine control=sequence
//ff:what isExceptAttacker — reports whether an attacker acts as an exception (defeater) gate
package toulmin

// isExceptAttacker reports whether attacker id is an exception (Except / Defeater),
// judged by its Strength or its registered role. Exceptions gate the target's
// applicability multiplicatively; every other attacker is a rebuttal in the denominator.
func (ec *evalContext) isExceptAttacker(id string) bool {
	return ec.strMap[id] == Defeater || ec.roleMap[id] == "except"
}
