//ff:type feature=engine type=model
//ff:what GraphDef — runtime graph definition for LoadGraph and YAML loading
package toulmin

// GraphDef defines a graph structure for dynamic loading.
type GraphDef struct {
	Graph   string         `yaml:"graph"`
	Rules   []GraphRuleDef `yaml:"rules"`
	Defeats []GraphEdgeDef `yaml:"defeats"`
}

// GraphRuleDef defines a single rule entry.
type GraphRuleDef struct {
	Name      string  `yaml:"name"`
	Role      string  `yaml:"role"`      // "warrant", "rebuttal", "defeater"
	Qualifier float64 `yaml:"qualifier"` // 0 means default 1.0
}

// GraphEdgeDef defines a single defeat edge.
type GraphEdgeDef struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}
