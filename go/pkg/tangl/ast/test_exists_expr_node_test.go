//ff:func feature=tangl type=model control=sequence
//ff:what TestExists_ExprNode — tests Exists.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestExists_ExprNode(t *testing.T) {
	e := Exists{}
	e.exprNode()

	var x Expr = Exists{}
	x.exprNode()
}
