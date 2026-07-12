//ff:func feature=tangl type=validator control=selection pattern=recursive
//ff:what walkExpr — preorder-visit every Expr node of a Rules condition tree
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// walkExpr visits every Expr node of a condition tree in preorder, calling fn
// on each. It recurses into Logic/Either terms, a Not term, and an Exists
// condition. Operand children (Compare.Lhs/Rhs, Exists.Source) are not Expr
// nodes; callers inspect them via walkOperand from within fn.
func walkExpr(e ast.Expr, fn func(ast.Expr)) {
	if e == nil {
		return
	}
	fn(e)
	switch v := e.(type) {
	case ast.Logic:
		for _, t := range v.Terms {
			walkExpr(t, fn)
		}
	case ast.Either:
		for _, t := range v.Terms {
			walkExpr(t, fn)
		}
	case ast.Not:
		walkExpr(v.Term, fn)
	case ast.Exists:
		walkExpr(v.Cond, fn)
	}
}
