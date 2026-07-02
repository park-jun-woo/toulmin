//ff:func feature=tangl type=parser control=sequence
//ff:what parseNotTerm — parse the negated inner expression of a "not" term
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseNotTerm parses the inner expression following "not": either an
// inline comparison ("not `f` is empty") or a nested condition-list ("not"
// alone, with children).
func parseNotTerm(rest string, it item, path string) (ast.Expr, error) {
	if strings.TrimSpace(rest) == "" {
		return parseCondList(it.Children, path)
	}
	return parseCompare(rest, it.Line, path)
}
