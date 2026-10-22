//ff:func feature=tangl type=parser control=sequence
//ff:what renderCompareKo — render a leaf comparison as its KO predicate (`` `f` 가 <v> <연산자꼴> ``)
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderCompareKo renders a Compare leaf as its KO SOV predicate, the inverse of
// parseCompareKo. The left side is a bare `field` when Lhs is nil (the byte-identical
// path for pre-existing documents) or a rendered operand otherwise; the right side is
// the verbatim Value when Rhs is nil or a rendered operand otherwise. The bare boolean
// Compare{Op:"equals", Value:"true"} (Rhs nil) renders back to the bare `x` 이면 form
// (or `x` 이 아니고 when negated) so re-encoding is identity; any other operator renders
// <lhs> 가 [<value> [와]] <predicate> with 받침-selected 이/가 and 와/과 particles and the
// phrase from compareKoPhrase. negated selects the Not-wrapped operator variant.
func renderCompareKo(c ast.Compare, negated bool) (string, error) {
	lhs, lhsBase, err := compareLhsTextKo(c)
	if err != nil {
		return "", err
	}
	if c.Op == "equals" && c.Value == "true" && c.Rhs == nil {
		if negated {
			return lhs + " " + josaFor(lhsBase, "이", "가") + " 아니고", nil
		}
		return lhs + " " + josaFor(lhsBase, "이면", "면"), nil
	}
	phrase, hasValue, objParticle, ok := compareKoPhrase(c.Op, negated)
	if !ok {
		return "", fmt.Errorf("ko: no KO phrase for operator %q (negated=%v)", c.Op, negated)
	}
	s := lhs + " " + josaFor(lhsBase, "이", "가")
	if !hasValue {
		return s + " " + phrase, nil
	}
	val, valBase, err := compareRhsTextKo(c)
	if err != nil {
		return "", err
	}
	s += " " + val
	if objParticle {
		s += " " + josaFor(valBase, "과", "와")
	}
	return s + " " + phrase, nil
}
