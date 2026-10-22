//ff:func feature=tangl type=parser control=sequence
//ff:what parseInlineRuleItem — parse one Rules entry, one-line or nested-tree form
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseInlineRuleItem parses "`name` when <field> <op> <value>" (one line) or
// "`name` when" with a nested either/and/or/not condition tree.
func parseInlineRuleItem(it item, path string) (ast.InlineRule, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.InlineRule{}, errAt(path, it.Line, "expected backtick-quoted rule name, got %q", it.Text)
	}
	if using, host := takeKeyword(rest, "uses"); host {
		ref, with, qual, grounds, tail, err := parseUsingClause(using, path, it.Line)
		if err != nil {
			return ast.InlineRule{}, err
		}
		if qual != nil {
			return ast.InlineRule{}, errAt(path, it.Line, "qualified belongs to a case binding, not a rule declaration")
		}
		if strings.TrimSpace(tail) != "" {
			return ast.InlineRule{}, errAt(path, it.Line, "unexpected trailing text: %q", tail)
		}
		return ast.InlineRule{Name: name, Using: &ref, With: with, Grounds: grounds, Line: it.Line}, nil
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.InlineRule{}, errAt(path, it.Line, "expected 'when' after rule name (or 'uses' for a host rule)")
	}
	rest = strings.TrimSpace(rest)
	cond, err := parseInlineRuleHead(rest, it, path)
	if err != nil {
		return ast.InlineRule{}, err
	}
	return ast.InlineRule{Name: name, Cond: cond, Line: it.Line}, nil
}
