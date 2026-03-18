//ff:func feature=engine type=engine control=sequence
//ff:what Defeater — adds a defeater rule to the graph builder
package toulmin

// Defeater adds a defeater rule to the graph. Qualifier defaults to 1.0.
func (b *GraphBuilder) Defeater(fn func(any, any) bool, qualifier ...float64) *GraphBuilder {
	q := 1.0
	if len(qualifier) > 0 {
		q = qualifier[0]
	}
	b.rules = append(b.rules, RuleMeta{
		Name:      FuncName(fn),
		Qualifier: q,
		Strength:  Defeater,
		Fn:        fn,
	})
	return b
}
