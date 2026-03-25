//ff:func feature=engine type=engine control=sequence
//ff:what Except — registers an exception rule and returns its reference
package toulmin

// Except registers an exception (defeater) rule in the graph and returns a *Rule reference.
// Use With() to add Specs and Qualifier() to set weight.
func (g *Graph) Except(fn any) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, nil)
	idx := len(g.rules)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: 1.0,
		Strength:  Defeater,
		Fn:        wrapped,
	})
	g.roles[id] = "except"
	return &Rule{id: id, graph: g, idx: idx, fn: fn}
}
