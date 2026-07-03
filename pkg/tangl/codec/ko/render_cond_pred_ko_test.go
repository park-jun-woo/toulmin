//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCondPredKo — renderCondPredKo covers Compare, Not{Compare}, Not{non-Compare}, and the non-predicate default
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderCondPredKo checks all four renderCondPredKo branches: a leaf Compare
// (both its ok render and its propagated operator error), a Not wrapping a Compare
// (same two outcomes, negated), a Not wrapping something other than a Compare (the
// impossible shape — falls through to ok=false, nil), and any other Expr (Logic,
// Either) hitting the default case.
func TestRenderCondPredKo(t *testing.T) {
	okCompare := ast.Compare{Field: "직원", Op: "equals", Value: "true"}
	badCompare := ast.Compare{Field: "x", Op: "bogus", Value: "1"}

	t.Run("CompareOk", func(t *testing.T) {
		s, ok, err := renderCondPredKo(okCompare)
		if err != nil || !ok || s != "`직원` 이면" {
			t.Errorf("renderCondPredKo(Compare) = (%q, %v, %v), want (\"`직원` 이면\", true, nil)", s, ok, err)
		}
	})

	t.Run("CompareErr", func(t *testing.T) {
		s, ok, err := renderCondPredKo(badCompare)
		if err == nil || !ok || s != "" {
			t.Errorf("renderCondPredKo(bad Compare) = (%q, %v, %v), want (\"\", true, error)", s, ok, err)
		}
	})

	t.Run("NotCompareOk", func(t *testing.T) {
		s, ok, err := renderCondPredKo(ast.Not{Term: okCompare})
		if err != nil || !ok || s != "`직원` 이 아니고" {
			t.Errorf("renderCondPredKo(Not{Compare}) = (%q, %v, %v), want (\"`직원` 이 아니고\", true, nil)", s, ok, err)
		}
	})

	t.Run("NotCompareErr", func(t *testing.T) {
		s, ok, err := renderCondPredKo(ast.Not{Term: badCompare})
		if err == nil || !ok || s != "" {
			t.Errorf("renderCondPredKo(Not{bad Compare}) = (%q, %v, %v), want (\"\", true, error)", s, ok, err)
		}
	})

	t.Run("NotNonCompare", func(t *testing.T) {
		s, ok, err := renderCondPredKo(ast.Not{Term: ast.Logic{Op: "and", Terms: []ast.Expr{okCompare, okCompare}}})
		if err != nil || ok || s != "" {
			t.Errorf("renderCondPredKo(Not{Logic}) = (%q, %v, %v), want (\"\", false, nil)", s, ok, err)
		}
	})

	t.Run("Default", func(t *testing.T) {
		s, ok, err := renderCondPredKo(ast.Either{Terms: []ast.Expr{okCompare}})
		if err != nil || ok || s != "" {
			t.Errorf("renderCondPredKo(Either) = (%q, %v, %v), want (\"\", false, nil)", s, ok, err)
		}
	})
}
