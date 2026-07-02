//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what renderLogicExpr — joins a Logic node's Terms with && or ||
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderLogicExpr joins a Logic node's Terms with && ("and") or ||
// ("or"), recursing into each term via renderExpr.
func renderLogicExpr(l ast.Logic) (string, error) {
	sep := " && "
	if l.Op == "or" {
		sep = " || "
	}
	parts := make([]string, 0, len(l.Terms))
	for _, t := range l.Terms {
		s, err := renderExpr(t)
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return "(" + strings.Join(parts, sep) + ")", nil
}
