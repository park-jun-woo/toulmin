//ff:func feature=engine type=engine control=sequence
//ff:what Rebuttal — adds a rebuttal rule to the graph builder
package toulmin

// Rebuttal adds a rebuttal rule to the graph.
// fn accepts both func(any,any,any)(bool,any) and legacy func(any,any)(bool,any).
// backing is the rule's judgment criteria (Toulmin backing). Use nil if not needed.
// qualifier is the rule's confidence weight.
func (b *GraphBuilder) Rebuttal(fn any, backing any, qualifier float64) *GraphBuilder {
	wrapped := toRuleFunc(fn)
	id := ruleID(fn, backing)
	b.rules = append(b.rules, RuleMeta{
		Name:      id,
		Qualifier: qualifier,
		Strength:  Defeasible,
		Backing:   backing,
		Fn:        wrapped,
	})
	b.roles[id] = "rebuttal"
	return b
}
