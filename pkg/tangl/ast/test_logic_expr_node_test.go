//ff:func feature=tangl type=model control=sequence
//ff:what TestLogic_ExprNode — tests Logic.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestLogic_ExprNode(t *testing.T) {
	l := Logic{}
	l.exprNode()

	var x Expr = Logic{}
	x.exprNode()
}
