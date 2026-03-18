//ff:func feature=engine type=engine control=sequence
//ff:what Rebuttal — adds a rebuttal rule to the graph builder
package toulmin

// Rebuttal adds a rebuttal rule to the graph. Qualifier defaults to 1.0.
func (b *GraphBuilder) Rebuttal(fn func(any, any) bool, qualifier ...float64) *GraphBuilder {
	q := 1.0
	if len(qualifier) > 0 {
		q = qualifier[0]
	}
	name := FuncName(fn)
	b.rules = append(b.rules, RuleMeta{
		Name:      name,
		Qualifier: q,
		Strength:  Defeasible,
		Fn:        fn,
	})
	b.roles[name] = "rebuttal"
	return b
}
