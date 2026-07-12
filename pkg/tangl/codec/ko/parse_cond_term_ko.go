//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseCondTermKo — parse a single condition item (다음 중/이름 단독/이 아니고/비교식)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseCondTermKo parses one condition-list item, in order: a "다음 중" group (ast.Either
// over the item's children), an exists quantifier lookahead (ast.Exists), a bare
// boolean negation "<path> 이 아니고" (ast.Not), a "<FUNC_REF> 가 참이면" host predicate
// (ast.Call), a leaf comparison, or a bare boolean name/path "<path>" ("<path> 이면").
// A bare boolean maps to the EN-canonical Compare{Op:"equals", Value:"true"} and a
// 의-chain path rides in Lhs, so KO and EN parse to the same AST. The exists, 참이면,
// path Lhs, and operand Rhs branches all gate on tokens absent from pre-existing
// documents, so those keep a byte-identical AST.
func parseCondTermKo(text string, children []scan.Item, line int, path string) (ast.Expr, error) {
	text = strings.TrimSpace(text)
	if rest, ok := takeKeywordKo(text, "다음 중"); ok && strings.TrimSpace(rest) == "" {
		inner, err := parseCondListKo(children, path)
		if err != nil {
			return nil, err
		}
		return ast.Either{Terms: []ast.Expr{inner}}, nil
	}
	if expr, ok, err := parseExistsKo(text, children, line, path); ok || err != nil {
		return expr, err
	}
	if before, ok := takeSuffixKo(text, "아니고"); ok {
		if steps, r, okp := parsePathKo(before); okp && strings.TrimSpace(skipParticlesKo(r)) == "" {
			return ast.Not{Term: boolCompareKo(steps)}, nil
		}
	}
	if expr, ok, err := parseCallPredKo(text, line, path); ok || err != nil {
		return expr, err
	}
	if expr, err := parseCompareKo(text, line, path); err == nil {
		return expr, nil
	}
	if steps, r, ok := parsePathKo(text); ok && strings.TrimSpace(skipParticlesKo(r)) == "" {
		return boolCompareKo(steps), nil
	}
	return nil, errAtKo(path, line, "unrecognized condition %q", text)
}
