//ff:func feature=engine type=engine control=sequence
//ff:what Defeater — registers a defeater rule and returns its reference
package toulmin

// Defeater registers a defeater rule in the graph and returns a *Rule reference.
func (g *Graph) Defeater(fn any, backing any, qualifier float64) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, backing)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: qualifier,
		Strength:  Defeater,
		Backing:   backing,
		Fn:        wrapped,
	})
	g.roles[id] = "defeater"
	return &Rule{id: id}
}
