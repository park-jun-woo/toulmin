//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseCondTermKo — parse a single condition item (다음 중/이름 단독/이 아니고/비교식)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseCondTermKo parses one condition-list item: a "다음 중" group (ast.Either over
// the item's children), a bare boolean-name negation "`x` 이 아니고" (ast.Not), a
// leaf comparison, or a bare boolean name "`x`" ("`x` 이면", 이면 dropped as a
// particle). A bare boolean name maps to the EN-canonical Compare{Op:"equals",
// Value:"true"} so KO and EN parse to the same AST.
func parseCondTermKo(text string, children []scan.Item, line int, path string) (ast.Expr, error) {
	text = strings.TrimSpace(text)
	if rest, ok := takeKeywordKo(text, "다음 중"); ok && strings.TrimSpace(rest) == "" {
		inner, err := parseCondListKo(children, path)
		if err != nil {
			return nil, err
		}
		return ast.Either{Terms: []ast.Expr{inner}}, nil
	}
	if before, ok := takeSuffixKo(text, "아니고"); ok {
		if name, r, okb := takeBacktick(before); okb && strings.TrimSpace(skipParticlesKo(r)) == "" {
			return ast.Not{Term: ast.Compare{Field: name, Op: "equals", Value: "true"}}, nil
		}
	}
	if expr, err := parseCompareKo(text, line, path); err == nil {
		return expr, nil
	}
	if name, r, ok := takeBacktick(text); ok && strings.TrimSpace(skipParticlesKo(r)) == "" {
		return ast.Compare{Field: name, Op: "equals", Value: "true"}, nil
	}
	return nil, errAtKo(path, line, "unrecognized condition %q", text)
}
