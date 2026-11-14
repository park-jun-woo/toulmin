//ff:func feature=tangl type=model control=sequence
//ff:what exprNode — marks Not as an Expr implementation
package ast

// exprNode marks Not as an Expr implementation.
func (n Not) exprNode() {}
