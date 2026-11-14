//ff:func feature=tangl type=model control=sequence
//ff:what TestCompare_ExprNode — tests Compare.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestCompare_ExprNode(t *testing.T) {
	c := Compare{}
	c.exprNode()

	var e Expr = Compare{}
	e.exprNode()
}
