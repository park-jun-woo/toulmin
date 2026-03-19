//ff:type feature=engine type=model
//ff:what Rule — opaque reference to a registered rule in a Graph
package toulmin

// Rule is an opaque reference to a rule registered in a Graph.
// It is returned by Warrant, Rebuttal, and Defeater, and used
// by Defeat to declare relationships without repeating rule definitions.
type Rule struct {
	id string
}
