//ff:func feature=engine type=engine control=sequence
//ff:what Defeat — declares a defeat edge between two Rule references
package toulmin

// Defeat declares a defeat edge: from attacks to.
// Both arguments are *Rule references returned by Warrant/Rebuttal/Defeater.
func (g *Graph) Defeat(from *Rule, to *Rule) {
	g.defeats = append(g.defeats, defeatEdge{
		from: from.id,
		to:   to.id,
	})
}
