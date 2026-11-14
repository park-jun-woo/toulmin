//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCondPredKoCallExists — predicate calls and single-compare exists render as one line; others defer
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderCondPredKoCallExists(t *testing.T) {
	got, ok, err := renderCondPredKo(ast.Call{Fn: ast.Ref{Alias: "정책", Name: "f"}})
	if err != nil || !ok || got != "`정책`.`f` 가 참이면" {
		t.Fatalf("call = %q ok=%v err=%v", got, ok, err)
	}
	ex := ast.Exists{Source: ast.Path{Steps: []string{"파일들"}}, Binder: "파일", Cond: ast.Compare{Lhs: ast.Path{Steps: []string{"파일", "이름"}}, Op: "is not empty"}}
	got, ok, err = renderCondPredKo(ex)
	if err != nil || !ok || got == "" {
		t.Fatalf("exists = %q ok=%v err=%v", got, ok, err)
	}
	ex.Cond = ast.Logic{Op: "and", Terms: []ast.Expr{ex.Cond, ex.Cond}}
	if got, ok, err := renderCondPredKo(ex); ok || err != nil || got != "" {
		t.Fatalf("multi-term exists must defer: %q ok=%v err=%v", got, ok, err)
	}
}
