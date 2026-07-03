//ff:func feature=engine type=engine control=sequence
//ff:what buildLevelGraph — test helper: builds the Level integration graph fixture
package toulmin

func buildLevelGraph() *Graph {
	g := NewGraph("level")
	g.Rule(lvlWarrant)
	g.Counter(lvlCounter)
	g.Except(lvlExcept)
	g.Counter(lvlOverride).Level(Review)
	g.Counter(lvlOptOut).Level(Unset)
	g.Rule(lvlFact)
	return g
}
