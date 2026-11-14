//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseRoleExprComma — each role directly followed by a comma keeps the comma in the remainder
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseRoleExprComma(t *testing.T) {
	for _, c := range []struct {
		in   string
		role ast.Role
	}{
		{"  general rule, level fail", ast.GeneralRule},
		{"counter rule, level warn", ast.CounterRule},
		{"except rule,", ast.ExceptRule},
	} {
		role, rest, ok := parseRoleExpr(c.in)
		if !ok || role != c.role || rest[0] != ',' {
			t.Errorf("%q: role=%v rest=%q ok=%v", c.in, role, rest, ok)
		}
	}
}
