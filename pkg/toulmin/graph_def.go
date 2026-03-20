//ff:type feature=engine type=model
//ff:what GraphDef — runtime graph definition for LoadGraph
package toulmin

// GraphDef defines a graph structure for dynamic loading.
type GraphDef struct {
	Graph   string
	Rules   []GraphRuleDef
	Defeats []GraphEdgeDef
}

// GraphRuleDef defines a single rule entry.
type GraphRuleDef struct {
	Name      string
	Role      string  // "warrant", "rebuttal", "defeater"
	Qualifier float64 // 0 means default 1.0
}

// GraphEdgeDef defines a single defeat edge.
type GraphEdgeDef struct {
	From string
	To   string
}
