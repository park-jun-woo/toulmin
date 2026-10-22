//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRoleFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRoleFromAST(t *testing.T) {
	for _, v := range []ast.Role{ast.GeneralRule, ast.CounterRule, ast.ExceptRule} {
		c, err := RoleFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := RoleToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := RoleFromAST(ast.Role(9)); err == nil {
		t.Fatal("unknown value: want error")
	}
}
