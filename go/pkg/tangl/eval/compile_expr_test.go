//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestCompileExpr — tests node-type dispatch for every Expr kind and the nil error
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompileExpr(t *testing.T) {
	res := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
		"\x00p": func(toulmin.Context, toulmin.Specs) (bool, any) { return true, nil },
	}}
	truthy := ast.Compare{Field: "x", Op: "equals", Value: "1"}
	nodes := []ast.Expr{
		truthy,
		ast.Logic{Op: "and", Terms: []ast.Expr{truthy}},
		ast.Either{Terms: []ast.Expr{truthy}},
		ast.Not{Term: truthy},
		ast.Exists{Source: ast.Path{Steps: []string{"xs"}}, Binder: "e", Cond: truthy},
		ast.Call{Fn: ast.Ref{Name: "p"}},
	}
	for _, n := range nodes {
		if _, err := compileExpr(n, res); err != nil {
			t.Fatalf("compileExpr(%T): %v", n, err)
		}
	}
	if _, err := compileExpr(nil, res); err == nil {
		t.Fatal("want error for nil expression node")
	}
}
