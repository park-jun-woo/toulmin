//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCompareKo — parse a leaf '`field` 가 <값> <연산자꼴>' comparison
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseCompareKo parses a KO leaf comparison. The SOV operator is the trailing
// phrase, so compareOpsKo (longest first) drives the match: the phrase is peeled
// off the end, the field backtick is read from the front, and the middle (minus
// its trailing 와/과 object particle) is the verbatim value. A negated phrase
// yields ast.Not{Compare}; 비어 있지 않으면 yields the direct "is not empty" operator.
func parseCompareKo(cond string, line int, path string) (ast.Expr, error) {
	cond = strings.TrimSpace(cond)
	for _, o := range compareOpsKo {
		before, ok := takeSuffixKo(cond, o.Ko)
		if !ok {
			continue
		}
		field, rest, ok := takeBacktick(before)
		if !ok {
			return nil, errAtKo(path, line, "expected backtick-quoted field in %q", cond)
		}
		value := trimTrailingParticleKo(strings.TrimSpace(skipParticlesKo(rest)))
		cmp := ast.Compare{Field: field, Op: o.Op, Value: value}
		if o.Negated {
			return ast.Not{Term: cmp}, nil
		}
		return cmp, nil
	}
	return nil, errAtKo(path, line, "unrecognized comparison operator in %q", cond)
}
