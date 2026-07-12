//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCompareOperands — renderCompare covers the additive Lhs/Rhs operand branches
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderCompareOperands checks the additive Lhs/Rhs operand branches: a
// nil/nil Compare is byte-identical to the pre-existing path, a path Lhs, a spec
// Lhs, a value-Call Lhs, a single-step backtick Rhs, and both sides as operands.
func TestRenderCompareOperands(t *testing.T) {
	cases := []struct {
		c    ast.Compare
		want string
	}{
		// Both nil: unchanged flat-field path.
		{ast.Compare{Field: "credit score", Op: "is at least", Value: "650"}, "`credit score` is at least 650"},
		// Path Lhs.
		{ast.Compare{Op: "is at least", Value: "0.7", Lhs: ast.Path{Steps: []string{"author", "trustScore"}}}, "`author`.`trustScore` is at least 0.7"},
		// Spec Lhs, single-step backtick Rhs operand.
		{ast.Compare{Op: "is at most", Lhs: ast.SpecRef{Steps: []string{"discount", "min price"}}, Rhs: ast.Path{Steps: []string{"basePrice"}}}, "spec `discount`.`min price` is at most `basePrice`"},
		// Value-Call Lhs.
		{ast.Compare{Op: "is at least", Value: "0.7", Lhs: ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}}, "result of `nlp`.`score` is at least 0.7"},
		// Rhs operand only (bare-field Lhs via Field).
		{ast.Compare{Field: "ip", Op: "equals", Rhs: ast.Path{Steps: []string{"clientIp"}}}, "`ip` equals `clientIp`"},
	}
	for _, tc := range cases {
		if got := renderCompare(tc.c); got != tc.want {
			t.Errorf("renderCompare(%+v) = %q, want %q", tc.c, got, tc.want)
		}
	}
}
