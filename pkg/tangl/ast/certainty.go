//ff:type feature=tangl type=model
//ff:what Certainty — a confidence gate clause (`if <op> <N>% certain`)
package ast

// Certainty is the `if <op> <N>% certain` gate on a do edge.
// Op is one of: "at least", "above", "less than", "at most".
type Certainty struct {
	Op      string `json:"op"`
	Percent int    `json:"percent"`
}
