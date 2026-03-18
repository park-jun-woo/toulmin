//ff:type feature=engine type=model
//ff:what Node — activated rule node in defeats graph
package toulmin

// Node represents an activated rule in the defeats graph.
type Node struct {
	Name      string
	Qualifier float64
	Strength  Strength
}
