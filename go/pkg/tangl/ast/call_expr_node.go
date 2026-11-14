//ff:func feature=tangl type=model control=sequence
//ff:what exprNode — marks Call as an Expr implementation (predicate position)
package ast

// exprNode marks Call as an Expr implementation (predicate position).
func (c Call) exprNode() {}
