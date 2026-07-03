//ff:func feature=engine type=engine control=selection
//ff:what defaultLevelForRole — maps a rule role to its default crisp Level
package toulmin

// defaultLevelForRole returns the default Level for a rule role when the rule
// did not set one explicitly. Only violation roles get a non-Unset default:
// counter→Fail, except→Review; every other role (including "rule") →Unset.
func defaultLevelForRole(role string) Level {
	switch role {
	case "counter":
		return Fail
	case "except":
		return Review
	default:
		return Unset
	}
}
