//ff:func feature=engine type=engine control=sequence
//ff:what Warrant — registers a warrant rule and returns its reference
package toulmin

// Warrant registers a warrant rule in the graph and returns a *Rule reference.
// fn accepts both func(any,any,any)(bool,any) and legacy func(any,any)(bool,any).
// backing is the rule's judgment criteria (Toulmin backing). Use nil if not needed.
// qualifier is the rule's confidence weight.
func (g *Graph) Warrant(fn any, backing Backing, qualifier float64) *Rule {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, backing)
	g.rules = append(g.rules, RuleMeta{
		Name:      id,
		Qualifier: qualifier,
		Strength:  Defeasible,
		Backing:   backing,
		Fn:        wrapped,
	})
	g.roles[id] = "warrant"
	return &Rule{id: id}
}
