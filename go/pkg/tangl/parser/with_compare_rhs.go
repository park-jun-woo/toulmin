//ff:func feature=tangl type=parser control=sequence
//ff:what withCompareRhs — attach a comparison right side as an operand Rhs or a literal Value
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// withCompareRhs parses the right side of a comparison. When it begins with an
// operand introducer (backtick, "spec", or "result") it becomes cmp.Rhs;
// otherwise it stays a literal Value string. A literal surrounded by matching
// double quotes has only those outer quotes removed; unquoted literals retain
// the pre-existing pass-through behavior.
func withCompareRhs(cmp ast.Compare, rhs string, line int, path string) (ast.Compare, error) {
	op, rest, ok, err := parseOperand(rhs, line, path)
	if err != nil {
		return ast.Compare{}, err
	}
	if !ok {
		cmp.Value = normalizeLiteralValue(cmp.Op, strings.TrimSpace(rhs))
		return cmp, nil
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Compare{}, errAt(path, line, "unexpected trailing text after operand: %q", rest)
	}
	cmp.Rhs = op
	return cmp, nil
}
