//ff:func feature=engine type=engine control=sequence
//ff:what inferRole — infers rule role from RuleMeta fields
package toulmin

// inferRole determines the role of a rule from its metadata.
func inferRole(r RuleMeta) string {
	if r.Strength == Defeater {
		return "defeater"
	}
	if len(r.Defeats) > 0 {
		return "rebuttal"
	}
	return "warrant"
}
