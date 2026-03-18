//ff:type feature=engine type=engine
//ff:what GraphBuilder — chainable builder for defeats graph with function identifiers
package toulmin

// GraphBuilder accumulates rules and defeat edges for graph construction.
type GraphBuilder struct {
	name    string
	rules   []RuleMeta
	defeats []defeatEdge
}
