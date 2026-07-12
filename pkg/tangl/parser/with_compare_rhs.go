//ff:func feature=tangl type=parser control=sequence
//ff:what withCompareRhs — attach a comparison right side as an operand Rhs or a literal Value
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// withCompareRhs parses the right side of a comparison. When it begins with an
// operand introducer (backtick, "spec", or "result") it becomes cmp.Rhs;
// otherwise it stays a literal Value string, preserving the pre-existing
// pass-through behavior for literals like 650, "VIP", or a, b, c.
func withCompareRhs(cmp ast.Compare, rhs string, line int, path string) (ast.Compare, error) {
	op, rest, ok, err := parseOperand(rhs, line, path)
	if err != nil {
		return ast.Compare{}, err
	}
	if !ok {
		cmp.Value = strings.TrimSpace(rhs)
		return cmp, nil
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Compare{}, errAt(path, line, "unexpected trailing text after operand: %q", rest)
	}
	cmp.Rhs = op
	return cmp, nil
}
