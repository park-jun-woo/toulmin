//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderOperand — renderOperand covers Path, SpecRef, aliased/bare Call, and the defensive default
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderOperand exercises each operand form: a single-step and multi-step
// Path, a SpecRef, and a value-position Call with and without an alias. The
// defensive default branch is reached with *ast.Path — the pointer still
// satisfies Operand (value-receiver marker) but matches no value-type case.
func TestRenderOperand(t *testing.T) {
	cases := []struct {
		o    ast.Operand
		want string
	}{
		{ast.Path{Steps: []string{"basePrice"}}, "`basePrice`"},
		{ast.Path{Steps: []string{"author", "trustScore"}}, "`author`.`trustScore`"},
		{ast.SpecRef{Steps: []string{"discount", "min price"}}, "spec `discount`.`min price`"},
		{ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, "result of `nlp`.`score`"},
		{ast.Call{Fn: ast.Ref{Name: "score"}}, "result of `score`"},
		{&ast.Path{Steps: []string{"x"}}, "`?`"},
	}
	for _, tc := range cases {
		if got := renderOperand(tc.o); got != tc.want {
			t.Errorf("renderOperand(%#v) = %q, want %q", tc.o, got, tc.want)
		}
	}
}
