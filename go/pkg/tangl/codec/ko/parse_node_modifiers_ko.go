//ff:func feature=tangl type=parser control=sequence
//ff:what parseNodeModifiersKo — consume the fixed-order [기준][근거][가중치][보고] adnominal modifiers
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseNodeModifiersKo consumes a node declaration's optional adnominal modifiers
// in their fixed canonical order — '`a` [그리고 `b`]* 을 기준으로 하는' (With),
// '`g` [그리고 `h`]* 를 근거로 하는' (Grounds, legacy migration input),
// '가중치 <float> 의' (Qualified), '(실패를|검토를) 보고하는' (Level) — and returns the
// unconsumed remainder (the optional '<fn> 쓰는' clause, or empty). Commas between
// modifiers and the particles 을/를/의 are peeled by skipParticlesKo. Out-of-order
// modifiers are left unconsumed and surface as the caller's '쓰는' error.
func parseNodeModifiersKo(n *ast.Node, s string, line int, path string) (string, error) {
	s = skipParticlesKo(s)
	if with, rest, ok := parseWithModifierKo(s); ok {
		n.With = with
		s = skipParticlesKo(rest)
	}
	if grounds, rest, ok := parseGroundModifierKo(s); ok {
		n.Grounds = grounds
		s = skipParticlesKo(rest)
	}
	if rest, ok := takeKeywordKo(s, "가중치"); ok {
		q, rest2, err := parseQualifiedKo(rest, line, path)
		if err != nil {
			return "", err
		}
		n.Qualified = &q
		s = skipParticlesKo(rest2) // 후행 독립 조사 '의' 소거
	}
	if lvl, rest, ok, err := parseReportKo(s, line, path); err != nil {
		return "", err
	} else if ok {
		n.Level = lvl
		s = skipParticlesKo(rest)
	}
	return strings.TrimSpace(s), nil
}
