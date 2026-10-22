//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseExistsTerm — parse an "any of <source> as <binder> satisfies <cond>" quantifier
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseExistsTerm parses an existential quantifier term. It consumes the
// leading quantifier keyword itself: "any" is implemented, while "all",
// "count", and "sum" are reserved and rejected with a clear error. When the
// text after "satisfies" is empty the body is the item's nested children
// (parsed as a condition list); otherwise it is a one-line comparison.
func parseExistsTerm(text string, it item, path string) (ast.Expr, error) {
	kw, rest, ok := takeQuantifier(text)
	if !ok {
		return nil, errAt(path, it.Line, "expected a quantifier, got %q", text)
	}
	if kw != "any" {
		return nil, errAt(path, it.Line, "quantifier %q is reserved and not yet supported", kw)
	}
	rest, ok = takeKeyword(rest, "of")
	if !ok {
		return nil, errAt(path, it.Line, "expected 'of' after 'any'")
	}
	src, rest, ok, err := parseOperand(rest, it.Line, path)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errAt(path, it.Line, "expected a collection operand after 'any of'")
	}
	rest, ok = takeKeyword(rest, "as")
	if !ok {
		return nil, errAt(path, it.Line, "expected 'as' after the collection operand")
	}
	binder, rest, ok := takeBacktick(rest)
	if !ok {
		return nil, errAt(path, it.Line, "expected a backtick-quoted binder after 'as'")
	}
	rest, ok = takeKeyword(rest, "satisfies")
	if !ok {
		return nil, errAt(path, it.Line, "expected 'satisfies' after the binder")
	}
	var cond ast.Expr
	if strings.TrimSpace(rest) == "" {
		cond, err = parseCondList(it.Children, path)
	} else {
		cond, err = parseCompare(rest, it.Line, path)
	}
	if err != nil {
		return nil, err
	}
	return ast.Exists{Source: src, Binder: binder, Cond: cond}, nil
}
