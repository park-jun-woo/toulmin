//ff:type feature=tangl type=model
//ff:what Not — negation of a sub-expression
package ast

// Not negates Term.
type Not struct {
	Term Expr `json:"term"`
}
