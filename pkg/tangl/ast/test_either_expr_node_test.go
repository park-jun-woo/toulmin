//ff:func feature=tangl type=model control=sequence
//ff:what TestEither_ExprNode — tests Either.exprNode marker method and Expr interface satisfaction
package ast

import "testing"

func TestEither_ExprNode(t *testing.T) {
	e := Either{}
	e.exprNode()

	var x Expr = Either{}
	x.exprNode()
}
