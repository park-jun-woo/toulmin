//ff:type feature=tangl type=model
//ff:what Exists — an "any of <source> as <binder> satisfies <cond>" quantifier
package ast

// Exists is an existential quantifier over Source: it holds when at least one
// element (bound to Binder) satisfies Cond.
type Exists struct {
	Source Operand `json:"source"`
	Binder string  `json:"binder"`
	Cond   Expr    `json:"cond"`
}
