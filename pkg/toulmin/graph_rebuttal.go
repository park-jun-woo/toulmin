//ff:func feature=engine type=engine control=sequence
//ff:what Rebuttal — registers a rebuttal rule and returns its reference
package toulmin

// Rebuttal registers a rebuttal rule in the graph and returns a *Rule reference.
func (g *Graph) Rebuttal(fn any, backing any, qualifier float64) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, backing)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: qualifier,
		Strength:  Defeasible,
		Backing:   backing,
		Fn:        wrapped,
	})
	g.roles[id] = "rebuttal"
	return &Rule{id: id}
}
