//ff:func feature=tangl type=engine control=selection pattern=recursive
//ff:what compileExpr — compiles one Expr node into an evaluating closure
package eval

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// compileExpr compiles one Expr node into an evalFunc, binding every nested
// Call to res. It dispatches on the concrete node type; a nil or unknown node
// is a compile error.
func compileExpr(e ast.Expr, res Resolver) (evalFunc, error) {
	switch v := e.(type) {
	case ast.Compare:
		return compileCompare(v, res)
	case ast.Logic:
		return compileLogic(v, res)
	case ast.Either:
		return compileEither(v, res)
	case ast.Not:
		return compileNot(v, res)
	case ast.Exists:
		return compileExists(v, res)
	case ast.Call:
		return compileCall(v, res)
	default:
		return nil, fmt.Errorf("eval: unsupported or nil expression node %T", e)
	}
}
