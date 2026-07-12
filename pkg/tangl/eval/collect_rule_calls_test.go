//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestCollectRuleCalls — tests bare-rule-call collection across every Expr kind
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectRuleCalls(t *testing.T) {
	ruleSet := map[string]bool{"a": true, "b": true, "c": true}
	rule := func(n string) ast.Call { return ast.Call{Fn: ast.Ref{Name: n}} }
	cond := ast.Logic{Op: "and", Terms: []ast.Expr{
		rule("a"),
		ast.Either{Terms: []ast.Expr{rule("b")}},
		ast.Not{Term: rule("c")},
		ast.Exists{Source: ast.Path{Steps: []string{"xs"}}, Binder: "e", Cond: rule("a")},
		ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "a"}},
		rule("other"),
		ast.Compare{Field: "x", Op: "equals", Value: "1"},
	}}
	var out []string
	collectRuleCalls(cond, ruleSet, &out)
	want := []string{"a", "b", "c", "a"}
	if len(out) != len(want) {
		t.Fatalf("want %v, got %v", want, out)
	}
	for i := range want {
		if out[i] != want[i] {
			t.Fatalf("want %v, got %v", want, out)
		}
	}
}
