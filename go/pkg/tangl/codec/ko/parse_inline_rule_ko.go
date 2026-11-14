//ff:func feature=tangl type=parser control=sequence
//ff:what parseInlineRuleKo — parse one 규칙 entry ('<조건> `n` 이다' or '다음이면 `n` 이다')
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseInlineRuleKo parses a tangl:규칙 entry. The sentence-final 이다 and the
// trailing backtick rule name are consumed first; the remaining head is the
// condition. A "다음이면" head opens a nested condition tree (the item's children);
// any other head is a single inline condition ending in ~면.
func parseInlineRuleKo(it scan.Item, path string) (ast.InlineRule, error) {
	if rule, ok, err := parseHostRuleKo(it, path); ok || err != nil {
		return rule, err
	}
	body, ok := takeSuffixKo(it.Text, "이다")
	if !ok {
		return ast.InlineRule{}, errAtKo(path, it.No, "expected '이다' at end of rule, got %q", it.Text)
	}
	name, prefix, ok := takeTrailingBacktick(body)
	if !ok {
		return ast.InlineRule{}, errAtKo(path, it.No, "expected backtick-quoted rule name before '이다'")
	}
	prefix = strings.TrimSpace(prefix)
	var cond ast.Expr
	var err error
	if prefix == "다음이면" {
		cond, err = parseCondListKo(it.Children, path)
	} else {
		cond, err = parseCondTermKo(prefix, nil, it.No, path)
	}
	if err != nil {
		return ast.InlineRule{}, err
	}
	return ast.InlineRule{Name: name, Cond: cond, Line: it.No}, nil
}
