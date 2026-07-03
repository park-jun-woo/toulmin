//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderNode — renderNode covers role/article, using+with+qualified, and checking
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderNode checks the node variants: a bare general rule, a counter rule
// using a function, an except rule (article "an") with `with`+`qualified`, and a
// general rule composing another case via `checking`.
func TestRenderNode(t *testing.T) {
	q := 0.8
	cases := []struct {
		n    ast.Node
		want string
	}{
		{ast.Node{Name: "authenticate", Role: ast.GeneralRule},
			"`authenticate` is a general rule"},
		{ast.Node{Name: "block ip", Role: ast.CounterRule, Using: &ast.Ref{Alias: "policy", Name: "IsIPBlocked"}},
			"`block ip` is a counter rule using `policy`.`IsIPBlocked`"},
		{ast.Node{Name: "low temp", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "s", Name: "low"}, With: []string{"a", "b"}, Qualified: &q},
			"`low temp` is an except rule using `s`.`low` with `a` and `b` qualified 0.8"},
		{ast.Node{Name: "renewal allowed", Role: ast.GeneralRule, Checking: "can renew loan"},
			"`renewal allowed` is a general rule checking `can renew loan`"},
	}
	for _, tc := range cases {
		if got := renderNode(tc.n); got != tc.want {
			t.Errorf("renderNode(%+v) = %q, want %q", tc.n, got, tc.want)
		}
	}
}
