//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRule — renderRule uses the one-line form for a lone compare and multi-line for a chain
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderRule checks that a single leaf Compare renders on one line, while an
// and/or Logic chain renders as a "`name` when" head plus indented list items.
func TestRenderRule(t *testing.T) {
	one := renderRule(ast.InlineRule{Name: "r", Cond: ast.Compare{Field: "f", Op: "equals", Value: `"x"`}})
	if one != "`r` when `f` equals \"x\"" {
		t.Errorf("one-line = %q", one)
	}
	chain := renderRule(ast.InlineRule{Name: "r", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
		ast.Compare{Field: "a", Op: "equals", Value: "1"},
		ast.Compare{Field: "b", Op: "equals", Value: "2"},
	}}})
	want := "`r` when\n  - `a` equals 1\n  - and `b` equals 2"
	if chain != want {
		t.Errorf("chain =\n%q\nwant\n%q", chain, want)
	}
}
