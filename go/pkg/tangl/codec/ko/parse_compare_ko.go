//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCompareKo — parse a leaf '`field` 가 <값> <연산자꼴>' comparison
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCompareKo parses a KO leaf comparison. The SOV operator is the trailing
// phrase, so compareOpsKo (longest first) drives the match: the 서술형 phrase (o.Ko)
// is peeled off the end, the left operand is read from the front by parseCompareLhsKo
// (a single-step backtick path is demoted to Field for AST back-compat), and the
// middle is attached by parseCompareRhsKo as an operand Rhs or a verbatim Value. A
// negated phrase yields ast.Not{Compare}; 비어 있지 않으면 yields the direct "is not
// empty" operator.
func parseCompareKo(cond string, line int, path string) (ast.Expr, error) {
	cond = strings.TrimSpace(cond)
	for _, o := range compareOpsKo {
		before, ok := takeSuffixKo(cond, o.Ko)
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
	return nil, errAtKo(path, line, "unrecognized comparison operator in %q", cond)
}
