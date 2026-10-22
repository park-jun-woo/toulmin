//ff:type feature=engine type=model
//ff:what Strength — rule strength classification (strict/defeasible/defeater)
package toulmin

// Strength controls whether a node accepts incoming attack edges.
type Strength int

const (
	Defeasible Strength = iota
	Strict
	Defeater
)
