//ff:func feature=tangl type=codegen control=selection pattern=recursive
//ff:what renderClosureExpr — compiles a Rules condition node into a witness closure literal
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureExpr compiles one Rules condition node into a Go closure
// literal of type func() (bool, any) — the (verdict, witness) contract shared
// with pkg/tangl/eval. It dispatches per node kind and recurses through the
// per-node renderers into Logic/Either/Not terms and an Exists condition.
func renderClosureExpr(e ast.Expr, env map[string]binderInfo, depth int) (string, error) {
	switch v := e.(type) {
	case ast.Compare:
		return renderClosureCompare(v, env)
	case ast.Logic:
		return renderClosureLogic(v, env, depth)
	case ast.Either:
		return renderClosureEither(v, env, depth)
	case ast.Not:
		return renderClosureNot(v, env, depth)
	case ast.Exists:
		return renderClosureExists(v, env, depth)
	case ast.Call:
		return renderClosureCall(v), nil
	default:
		return "", fmt.Errorf("gen: unsupported expression node %T", e)
	}
}
