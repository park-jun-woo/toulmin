//ff:func feature=tangl type=codegen control=selection pattern=recursive
//ff:what renderExpr — compiles a Rules condition tree into a Go boolean expression
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderExpr compiles a Rules condition tree into a single Go boolean
// expression: Compare leaves become tanglCompare calls, Logic/Either
// become parenthesized &&/|| groups, and Not becomes a negation. It
// recurses into Logic/Either/Not terms.
func renderExpr(e ast.Expr) (string, error) {
	switch v := e.(type) {
	case ast.Compare:
		return renderCompareExpr(v), nil
	case ast.Not:
		inner, err := renderExpr(v.Term)
		if err != nil {
			return "", err
		}
		return "!(" + inner + ")", nil
	case ast.Logic:
		return renderLogicExpr(v)
	case ast.Either:
		return renderEitherExpr(v)
	default:
		return "", fmt.Errorf("gen: unsupported expression node %T", e)
	}
}
