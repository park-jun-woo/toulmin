//ff:func feature=engine type=engine control=sequence
//ff:what With — adds a Spec to this rule (additive, supports chaining)
package toulmin

// With adds a Spec to the rule and returns the rule for chaining.
// Multiple With() calls compose specs. Updates ruleID.
func (r *Rule) With(spec Spec) *Rule {
	meta := &r.graph.rules[r.idx]
	meta.Specs = append(meta.Specs, spec)
	oldID := meta.Name
	newID := ruleID(r.fn, meta.Specs)
	meta.Name = newID
	r.id = newID
	if role, ok := r.graph.roles[oldID]; ok {
		delete(r.graph.roles, oldID)
		r.graph.roles[newID] = role
	}
	return r
}
