//ff:func feature=tangl type=parser control=sequence
//ff:what parseInlineRuleItem — parse one Rules entry, one-line or nested-tree form
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseInlineRuleItem parses "`name` when <field> <op> <value>" (one line) or
// "`name` when" with a nested either/and/or/not condition tree.
func parseInlineRuleItem(it item, path string) (ast.InlineRule, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.InlineRule{}, errAt(path, it.Line, "expected backtick-quoted rule name, got %q", it.Text)
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.InlineRule{}, errAt(path, it.Line, "expected 'when' after rule name")
	}
	rest = strings.TrimSpace(rest)
	var cond ast.Expr
	var err error
	if rest == "" {
		cond, err = parseCondList(it.Children, path)
	} else {
		cond, err = parseCompare(rest, it.Line, path)
	}
	if err != nil {
		return ast.InlineRule{}, err
	}
	return ast.InlineRule{Name: name, Cond: cond, Line: it.Line}, nil
}
