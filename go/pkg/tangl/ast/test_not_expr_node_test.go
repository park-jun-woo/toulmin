//ff:func feature=tangl type=model control=sequence
//ff:what TestNot_ExprNode — tests Not.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestNot_ExprNode(t *testing.T) {
	n := Not{}
	n.exprNode()

	var x Expr = Not{}
	x.exprNode()
}
