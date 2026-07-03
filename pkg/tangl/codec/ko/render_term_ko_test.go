//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what TestRenderTermKo — renderTermKo covers the prefix head, the one-line predicate (ok/err), the nested Either (ok/err), and the unsupported-term fallback
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderTermKo checks: an unprefixed and a prefixed one-line predicate (also
// covering the depth-driven indentation), a predicate render error propagated from
// renderCondPredKo, a nested Either that renders a "다음 중" head plus its inner
// condition list (both the success and the propagated-error shapes), and a term that
// is neither a one-line predicate nor an Either, which falls through to the
// unsupported-term error.
func TestRenderTermKo(t *testing.T) {
	okCompare := ast.Compare{Field: "직원", Op: "equals", Value: "true"}
	badCompare := ast.Compare{Field: "x", Op: "bogus", Value: "1"}

	t.Run("NoPrefix", func(t *testing.T) {
		got, err := renderTermKo("", okCompare, 0)
		if err != nil {
			t.Fatalf("renderTermKo: %v", err)
		}
		want := []string{"- `직원` 이면"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("renderTermKo = %+v, want %+v", got, want)
		}
	})

	t.Run("WithPrefixAndDepth", func(t *testing.T) {
		got, err := renderTermKo("그리고", okCompare, 1)
		if err != nil {
			t.Fatalf("renderTermKo: %v", err)
		}
		want := []string{"  - 그리고 `직원` 이면"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("renderTermKo = %+v, want %+v", got, want)
		}
	})

	t.Run("PredicateErr", func(t *testing.T) {
		if _, err := renderTermKo("", badCompare, 0); err == nil {
			t.Error("renderTermKo accepted a Compare with an unknown operator")
		}
	})

	t.Run("EitherOk", func(t *testing.T) {
		got, err := renderTermKo("또는", ast.Either{Terms: []ast.Expr{okCompare}}, 0)
		if err != nil {
			t.Fatalf("renderTermKo: %v", err)
		}
		want := []string{"- 또는 다음 중", "  - `직원` 이면"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("renderTermKo = %+v, want %+v", got, want)
		}
	})

	t.Run("EitherErr", func(t *testing.T) {
		inner := ast.Logic{Op: "and", Terms: []ast.Expr{okCompare, okCompare, okCompare}}
		if _, err := renderTermKo("", ast.Either{Terms: []ast.Expr{inner}}, 0); err == nil {
			t.Error("renderTermKo accepted an Either whose inner list cannot be rendered")
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		term := ast.Logic{Op: "and", Terms: []ast.Expr{okCompare, okCompare}}
		if _, err := renderTermKo("", term, 0); err == nil {
			t.Error("renderTermKo accepted a term that is neither a predicate nor an Either")
		}
	})
}
