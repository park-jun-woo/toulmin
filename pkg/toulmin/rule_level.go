//ff:func feature=engine type=engine control=sequence
//ff:what Level — sets the crisp reporting level for this rule (overrides role default)
package toulmin

// Level sets the rule's crisp reporting level and returns the rule for chaining.
// Calling it marks LevelSet true so the explicit value (including Unset) always
// wins over the role-based default in defaultLevelForRole.
func (r *Rule) Level(l Level) *Rule {
	r.graph.rules[r.idx].Level = l
	r.graph.rules[r.idx].LevelSet = true
	return r
}
