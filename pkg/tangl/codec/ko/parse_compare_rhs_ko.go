//ff:func feature=tangl type=parser control=sequence
//ff:what parseCompareRhsKo — attach a comparison right side as an operand Rhs or a literal Value
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseCompareRhsKo parses the middle text of a KO comparison (between the left
// operand and the peeled operator suffix) as the right side. After dropping its
// leading 조사, a right side that begins with an operand introducer (a backtick or
// 기준) becomes cmp.Rhs — with any residual bytes past the operand rejected as an
// error (a backtick literal-list abuse guard, mirroring the EN 우변 규범); anything
// else stays a literal Value string, trailing 와/과 object particle removed and a pair
// of surrounding double quotes stripped. It is the
// KO mirror of pkg/tangl/parser.withCompareRhs, shared by the 서술형 and 관형형 parsers.
func parseCompareRhsKo(cmp ast.Compare, rhs string, line int, path string) (ast.Compare, error) {
	mid := strings.TrimSpace(skipParticlesKo(rhs))
	op, rest, ok, err := parseOperandKo(mid, line, path)
	if err != nil {
		return ast.Compare{}, err
	}
	if !ok {
		cmp.Value = normalizeLiteralValueKo(cmp.Op, trimTrailingParticleKo(mid))
		return cmp, nil
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Compare{}, errAtKo(path, line, "unexpected trailing text after operand: %q", rest)
	}
	cmp.Rhs = op
	return cmp, nil
}

func normalizeLiteralValueKo(op, value string) string {
	if op != "is in" {
		return stripLiteralQuotesKo(value)
	}
	items := strings.Split(value, ",")
	for i := range items {
		items[i] = stripLiteralQuotesKo(strings.TrimSpace(items[i]))
	}
	return strings.Join(items, ", ")
}

func stripLiteralQuotesKo(value string) string {
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}
	return value
}
