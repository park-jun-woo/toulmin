//ff:func feature=tangl type=parser control=sequence
//ff:what mergeLogic — combine two condition expressions with and/or
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// mergeLogic combines left and right into a Logic expression using op.
func mergeLogic(left ast.Expr, op string, right ast.Expr) ast.Expr {
	return ast.Logic{Op: op, Terms: []ast.Expr{left, right}}
}
