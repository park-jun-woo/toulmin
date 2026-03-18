//ff:func feature=engine type=engine control=sequence
//ff:what Defeat — declares a defeat edge in the graph builder
package toulmin

// Defeat declares a defeat edge: from attacks to.
func (b *GraphBuilder) Defeat(from func(any, any) (bool, any), to func(any, any) (bool, any)) *GraphBuilder {
	b.defeats = append(b.defeats, defeatEdge{
		from: funcID(from),
		to:   funcID(to),
	})
	return b
}
