//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseCondTerm — parse a single condition item (either/not/compare)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCondTerm parses one condition-list item: an "any of ..." quantifier, an
// "either" group, a "not" negation, a "FUNC_REF holds" host predicate, or a
// leaf comparison.
func parseCondTerm(it item, path string) (ast.Expr, error) {
	text := strings.TrimSpace(it.Text)
	if _, _, ok := takeQuantifier(text); ok {
		return parseExistsTerm(text, it, path)
	}
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
	if expr, ok, err := tryParseCallTerm(text, it.Line, path); ok || err != nil {
		return expr, err
	}
	return parseCompare(text, it.Line, path)
}
