//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCompareAdnKo — parse a 관형형 (attributive) leaf comparison for a one-line exists body
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseCompareAdnKo parses a 관형형 (attributive) leaf comparison — the inner
// condition of a one-line exists, which does not end a sentence and so uses the
// ~ㄴ/은/는/인 attributive form instead of the 서술형 ~면. It is the parseCompareKo
// mirror over the o.KoAdn column. The bare-boolean attributive specials (인 → equals
// true, 이 아닌 → its Not) are tried only after the whole table so a bare 아닌 does not
// mis-strip the tail of 이하가 아닌 (the same longest-first rule as the 서술형 parser).
func parseCompareAdnKo(cond string, line int, path string) (ast.Expr, error) {
	cond = strings.TrimSpace(cond)
	for _, o := range compareOpsKo {
		before, ok := takeSuffixKo(cond, o.KoAdn)
		if !ok {
			continue
		}
		cmp, rest, ok, err := parseCompareLhsKo(before, line, path)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errAtKo(path, line, "expected backtick-quoted field in %q", cond)
		}
		cmp.Op = o.Op
		cmp, err = parseCompareRhsKo(cmp, rest, line, path)
		if err != nil {
			return nil, err
		}
		if o.Negated {
			return ast.Not{Term: cmp}, nil
		}
		return cmp, nil
	}
	if before, ok := takeSuffixKo(cond, "아닌"); ok {
		if steps, r, okp := parsePathKo(before); okp && strings.TrimSpace(skipParticlesKo(r)) == "" {
			return ast.Not{Term: boolCompareKo(steps)}, nil
		}
	}
	if before, ok := takeSuffixKo(cond, "인"); ok {
		if steps, r, okp := parsePathKo(before); okp && strings.TrimSpace(skipParticlesKo(r)) == "" {
			return boolCompareKo(steps), nil
		}
	}
	return nil, errAtKo(path, line, "unrecognized adnominal comparison in %q", cond)
}
