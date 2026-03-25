//ff:type feature=engine type=model
//ff:what Rule — opaque reference to a registered rule in a Graph
package toulmin

// Rule is an opaque reference to a rule registered in a Graph.
// It is returned by Rule, Counter, and Except, and supports
// Attacks, Backing, and Qualifier method chaining.
type Rule struct {
	id    string
	graph *Graph
	idx   int
	fn    any
}
