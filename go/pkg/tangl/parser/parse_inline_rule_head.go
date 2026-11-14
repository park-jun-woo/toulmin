//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseInlineRuleHead — dispatch a rule head's one-line/nested condition body
package parser

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// parseInlineRuleHead dispatches the condition body of a "`name` when <rest>"
// rule head: an empty rest parses the nested children as a condition list; a
// quantifier head parses an existential term; a "FUNC_REF holds" head parses a
// host predicate; anything else stays the pre-existing one-line comparison so
// non-quantifier, non-call heads keep a byte-identical AST.
func parseInlineRuleHead(rest string, it item, path string) (ast.Expr, error) {
	if rest == "" {
		return parseCondList(it.Children, path)
	}
	if _, _, ok := takeQuantifier(rest); ok {
		return parseExistsTerm(rest, it, path)
	}
	if expr, consumed, err := tryParseCallTerm(rest, it.Line, path); consumed || err != nil {
		return expr, err
	}
	return parseCompare(rest, it.Line, path)
}
