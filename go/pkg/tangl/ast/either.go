//ff:type feature=tangl type=model
//ff:what Either — an `either` group of sub-expressions
package ast

// Either groups Terms under an `either` block.
type Either struct {
	Terms []Expr `json:"terms"`
}
