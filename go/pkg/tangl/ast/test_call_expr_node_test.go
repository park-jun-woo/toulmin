//ff:func feature=tangl type=model control=sequence
//ff:what TestCall_ExprNode — tests Call.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestCall_ExprNode(t *testing.T) {
	c := Call{}
	c.exprNode()

	var e Expr = Call{}
	e.exprNode()
}
