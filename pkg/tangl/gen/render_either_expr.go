//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what renderEitherExpr — renders an Either group's Terms as one parenthesized group
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderEitherExpr renders an Either group's Terms (the parser always
// produces exactly one, the merged and/or tree inside the "either" block)
// as a single parenthesized group; multiple terms are joined with || as a
// defensive fallback.
func renderEitherExpr(e ast.Either) (string, error) {
	parts := make([]string, 0, len(e.Terms))
	for _, t := range e.Terms {
		s, err := renderExpr(t)
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return "(" + strings.Join(parts, " || ") + ")", nil
}
