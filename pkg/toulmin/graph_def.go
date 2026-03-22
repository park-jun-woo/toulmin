//ff:type feature=engine type=model
//ff:what GraphDef — runtime graph definition for LoadGraph and YAML loading
package toulmin

// GraphDef defines a graph structure for dynamic loading.
type GraphDef struct {
	Graph   string         `yaml:"graph"`
	Rules   []GraphRuleDef `yaml:"rules"`
	Defeats []GraphEdgeDef `yaml:"defeats"`
}
