//ff:func feature=engine type=engine control=sequence
//ff:what Defeater — adds a defeater rule to the graph builder
package toulmin

// Defeater adds a defeater rule to the graph. Qualifier defaults to 1.0.
func (b *GraphBuilder) Defeater(fn func(any, any) (bool, any), qualifier ...float64) *GraphBuilder {
	q := 1.0
	if len(qualifier) > 0 {
		q = qualifier[0]
	}
	id := funcID(fn)
	b.rules = append(b.rules, RuleMeta{
		Name:      id,
		Qualifier: q,
		Strength:  Defeater,
		Fn:        fn,
	})
	b.roles[id] = "defeater"
	return b
}
