//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRenderClosureExpr — tests node-kind dispatch for every Expr and the nil error
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderClosureExpr(t *testing.T) {
	cmp := ast.Compare{Field: "x", Op: "equals", Value: "1"}
	cases := []struct {
		node ast.Expr
		want string
	}{
		{cmp, "tanglCompareVal("},
		{ast.Logic{Op: "and", Terms: []ast.Expr{cmp}}, "tanglAnd("},
		{ast.Either{Terms: []ast.Expr{cmp}}, "tanglOr("},
		{ast.Not{Term: cmp}, "tanglNot("},
		{ast.Exists{Source: ast.Path{Steps: []string{"xs"}}, Binder: "e", Cond: cmp}, "tanglAny("},
		{ast.Call{Fn: ast.Ref{Name: "p"}}, "tanglHolds("},
	}
	for _, c := range cases {
		got, err := renderClosureExpr(c.node, nil, 0)
		if err != nil || !strings.Contains(got, c.want) {
			t.Errorf("renderClosureExpr(%T) should emit %s: %q (%v)", c.node, c.want, got, err)
		}
	}
	if _, err := renderClosureExpr(nil, nil, 0); err == nil {
		t.Error("want error for a nil expression node")
	}
}
