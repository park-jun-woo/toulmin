//ff:func feature=engine type=engine control=sequence
//ff:what Warrant — adds a warrant rule to the graph builder
package toulmin

// Warrant adds a warrant rule to the graph. Qualifier defaults to 1.0.
func (b *GraphBuilder) Warrant(fn func(any, any) (bool, any), qualifier ...float64) *GraphBuilder {
	q := 1.0
	if len(qualifier) > 0 {
		q = qualifier[0]
	}
	id := funcID(fn)
	b.rules = append(b.rules, RuleMeta{
		Name:      id,
		Qualifier: q,
		Strength:  Defeasible,
		Fn:        fn,
	})
	b.roles[id] = "warrant"
	return b
}
