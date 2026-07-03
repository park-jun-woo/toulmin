//ff:func feature=tangl type=parser control=sequence
//ff:what renderCompareKo — render a leaf comparison as its KO predicate (`` `f` 가 <v> <연산자꼴> ``)
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderCompareKo renders a Compare leaf as its KO SOV predicate, the inverse of
// parseCompareKo. The bare boolean Compare{Op:"equals", Value:"true"} — into which
// the KO parser absorbs a bare name `x` (이)면 (Phase 003) — renders back to the bare
// `x` 이면 form (or `x` 이 아니고 when negated) so re-encoding is identity; any other
// operator renders `field` 가 [<value> [와]] <predicate> with 받침-selected 이/가 and 와/과
// particles and the phrase from compareKoPhrase. negated selects the Not-wrapped
// operator variant (같지 않으면 etc.).
func renderCompareKo(c ast.Compare, negated bool) (string, error) {
	field := "`" + c.Field + "`"
	if c.Op == "equals" && c.Value == "true" {
		if negated {
			return field + " " + josaFor(c.Field, "이", "가") + " 아니고", nil
		}
		return field + " " + josaFor(c.Field, "이면", "면"), nil
	}
	phrase, hasValue, objParticle, ok := compareKoPhrase(c.Op, negated)
	if !ok {
		return "", fmt.Errorf("ko: no KO phrase for operator %q (negated=%v)", c.Op, negated)
	}
	s := field + " " + josaFor(c.Field, "이", "가")
	if hasValue {
		s += " " + c.Value
		if objParticle {
			s += " " + josaFor(c.Value, "과", "와")
		}
	}
	return s + " " + phrase, nil
}
