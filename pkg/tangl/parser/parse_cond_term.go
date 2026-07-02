//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseCondTerm — parse a single condition item (either/not/compare)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseCondTerm parses one condition-list item: an "either" group, a "not"
// negation, or a leaf field/operator/value comparison.
func parseCondTerm(it item, path string) (ast.Expr, error) {
	text := strings.TrimSpace(it.Text)
	if text == "either" {
		inner, err := parseCondList(it.Children, path)
		if err != nil {
			return nil, err
		}
		return ast.Either{Terms: []ast.Expr{inner}}, nil
	}
	if rest, ok := takeKeyword(text, "not"); ok {
		inner, err := parseNotTerm(rest, it, path)
		if err != nil {
			return nil, err
		}
		return ast.Not{Term: inner}, nil
	}
	return parseCompare(text, it.Line, path)
}
