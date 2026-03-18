//ff:func feature=engine type=engine control=sequence
//ff:what Defeat — declares a defeat edge in the graph builder
package toulmin

// Defeat declares a defeat edge: from attacks to.
func (b *GraphBuilder) Defeat(from func(any, any) bool, to func(any, any) bool) *GraphBuilder {
	b.defeats = append(b.defeats, defeatEdge{
		from: FuncName(from),
		to:   FuncName(to),
	})
	return b
}
