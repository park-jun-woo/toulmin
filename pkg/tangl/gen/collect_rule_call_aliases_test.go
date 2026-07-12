//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestCollectRuleCallAliases — tests alias collection across every Expr and operand kind
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectRuleCallAliases(t *testing.T) {
	rules := []ast.InlineRule{
		{Name: "r1", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{
				Lhs: ast.Call{Fn: ast.Ref{Alias: "a1", Name: "f"}},
				Op:  "equals",
				Rhs: ast.Call{Fn: ast.Ref{Alias: "a2", Name: "g"}},
			},
			ast.Compare{Field: "x", Op: "is empty"},
			ast.Compare{Lhs: ast.Path{Steps: []string{"x"}}, Op: "is empty"},
			ast.Either{Terms: []ast.Expr{ast.Call{Fn: ast.Ref{Alias: "a3", Name: "h"}}}},
			ast.Not{Term: ast.Call{Fn: ast.Ref{Name: "bare"}}},
			ast.Exists{
				Source: ast.Call{Fn: ast.Ref{Alias: "a4", Name: "list"}},
				Binder: "e",
				Cond:   ast.Call{Fn: ast.Ref{Alias: "a5", Name: "p"}},
			},
		}}},
	}
	used := map[string]bool{}
	collectRuleCallAliases(rules, used)
	for _, want := range []string{"a1", "a2", "a3", "a4", "a5"} {
		if !used[want] {
			t.Errorf("alias %q should be collected: %v", want, used)
		}
	}
	if len(used) != 5 {
		t.Errorf("bare calls and non-call operands must not add aliases: %v", used)
	}
}
