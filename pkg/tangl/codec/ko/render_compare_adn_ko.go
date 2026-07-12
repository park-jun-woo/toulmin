//ff:func feature=tangl type=parser control=sequence
//ff:what renderCompareAdnKo — render a leaf comparison in its 관형형 form for a one-line exists body
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderCompareAdnKo renders a Compare leaf in its 관형형 (attributive) form, the
// inverse of parseCompareAdnKo and the 관형형 mirror of renderCompareKo. It shares the
// same Lhs/Rhs operand handling but ends in the attributive phrase from
// compareKoAdnPhrase; the bare boolean renders `x` 인 (or `x` 이 아닌 when negated). It
// is used for the inner condition of a one-line exists.
func renderCompareAdnKo(c ast.Compare, negated bool) (string, error) {
	lhs, lhsBase, err := compareLhsTextKo(c)
	if err != nil {
		return "", err
	}
	if c.Op == "equals" && c.Value == "true" && c.Rhs == nil {
		if negated {
			return lhs + " " + josaFor(lhsBase, "이", "가") + " 아닌", nil
		}
		return lhs + " 인", nil
	}
	phrase, hasValue, objParticle, ok := compareKoAdnPhrase(c.Op, negated)
	if !ok {
		return "", fmt.Errorf("ko: no KO adnominal phrase for operator %q (negated=%v)", c.Op, negated)
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
