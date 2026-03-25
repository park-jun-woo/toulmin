//ff:func feature=engine type=engine control=sequence
//ff:what Rule — registers a rule and returns its reference
package toulmin

// Rule registers a rule in the graph and returns a *Rule reference.
// Default backing is nil and qualifier is 1.0. Use Backing() and Qualifier() to override.
func (g *Graph) Rule(fn any) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, nil)
	idx := len(g.rules)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: 1.0,
		Strength:  Defeasible,
		Backing:   nil,
		Fn:        wrapped,
	})
	g.roles[id] = "rule"
	return &Rule{id: id, graph: g, idx: idx, fn: fn}
}
