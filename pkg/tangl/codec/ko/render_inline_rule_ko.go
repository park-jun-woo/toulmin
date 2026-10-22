//ff:func feature=tangl type=parser control=sequence
//ff:what renderInlineRuleKo — render a 규칙 entry (one-line single condition or 다음이면 nested tree)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderInlineRuleKo renders a tangl:규칙 entry. A single-term condition — a leaf
// Compare or a negated Compare, which parseCondTermKo accepts inline — uses the
// one-line form "<조건> `name` 이다"; a Logic and/or chain or an Either group uses the
// multi-line form "다음이면 `name` 이다" plus the nested condition list. This is the
// inverse of parseInlineRuleKo, whose inline branch parses exactly one condition
// term while chains and either groups must be nested list items.
func renderInlineRuleKo(r ast.InlineRule) (string, error) {
	if r.Using != nil {
		s := "`" + r.Name + "` 은 "
		if len(r.With) > 0 {
			s += renderWithKo(r.With) + " "
		}
		if len(r.Grounds) > 0 {
			s += renderGroundsKo(r.Grounds) + " "
		}
		s += renderRefKo(*r.Using) + " 쓰는 규칙이다"
		return s, nil
	}
	if pred, ok, err := renderCondPredKo(r.Cond); err != nil {
		return "", err
	} else if ok {
		return pred + " `" + r.Name + "` 이다", nil
	}
	lines, err := renderCondListKo(r.Cond, 1)
	if err != nil {
		return "", err
	}
	return "다음이면 `" + r.Name + "` 이다\n" + strings.Join(lines, "\n"), nil
}
