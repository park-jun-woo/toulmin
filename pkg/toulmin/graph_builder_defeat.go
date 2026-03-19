//ff:func feature=engine type=engine control=sequence
//ff:what Defeat — declares a defeat edge in the graph builder (backing nil)
package toulmin

// Defeat declares a defeat edge: from attacks to.
// Uses funcID only (backing nil). For same-function different-backing, use DefeatWith.
func (b *GraphBuilder) Defeat(from any, to any) *GraphBuilder {
	b.defeats = append(b.defeats, defeatEdge{
		from: ruleID(from, nil),
		to:   ruleID(to, nil),
	})
	return b
}
