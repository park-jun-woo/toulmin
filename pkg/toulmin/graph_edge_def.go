//ff:type feature=engine type=model
//ff:what GraphEdgeDef — single defeat edge in graph definition
package toulmin

// GraphEdgeDef defines a single defeat edge.
type GraphEdgeDef struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}
