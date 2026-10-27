//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRoleToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRoleToAST(t *testing.T) {
	for _, v := range []ast.Role{ast.GeneralRule, ast.CounterRule, ast.ExceptRule} {
		c, _ := RoleFromAST(v)
		if got, err := RoleToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := RoleToAST(toulmin.RoleCode(9)); err == nil {
		t.Fatal("unknown code: want error")
	}
}
