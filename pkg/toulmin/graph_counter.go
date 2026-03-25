//ff:func feature=engine type=engine control=sequence
//ff:what Counter — registers a counter rule and returns its reference
package toulmin

// Counter registers a counter (rebuttal) rule in the graph and returns a *Rule reference.
// Use With() to add Specs and Qualifier() to set weight.
func (g *Graph) Counter(fn any) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, nil)
	idx := len(g.rules)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: 1.0,
		Strength:  Defeasible,
		Fn:        wrapped,
	})
	g.roles[id] = "counter"
	return &Rule{id: id, graph: g, idx: idx, fn: fn}
}
