//ff:type feature=engine type=engine
//ff:what Graph — defeats graph with rule registration and evaluation
package toulmin

// Graph accumulates rules and defeat edges for graph construction and evaluation.
type Graph struct {
	name    string
	rules   []RuleMeta
	roles   map[string]string
	defeats []defeatEdge
}
