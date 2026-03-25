//ff:func feature=engine type=engine control=sequence
//ff:what Backing — sets the backing for this rule and updates ruleID
package toulmin

// Backing sets the rule's judgment criteria and returns the rule for chaining.
// Updates the rule's identity since backing is part of ruleID.
func (r *Rule) Backing(b Backing) *Rule {
	meta := &r.graph.rules[r.idx]
	oldID := meta.Name
	meta.Backing = b
	newID := ruleID(r.fn, b)
	meta.Name = newID
	r.id = newID
	if role, ok := r.graph.roles[oldID]; ok {
		delete(r.graph.roles, oldID)
		r.graph.roles[newID] = role
	}
	return r
}
