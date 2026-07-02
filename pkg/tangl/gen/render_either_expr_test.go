//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderEitherExpr — tests renderEitherExpr for empty terms, multi-term join, and error propagation branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderEitherExpr(t *testing.T) {
	t.Run("empty terms", func(t *testing.T) {
		got, err := renderEitherExpr(ast.Either{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "()" {
			t.Errorf("got %q, want %q", got, "()")
		}
	})

	t.Run("multiple terms joined with or", func(t *testing.T) {
		e := ast.Either{Terms: []ast.Expr{
			ast.Compare{Field: "a", Op: "=", Value: "1"},
			ast.Compare{Field: "b", Op: "=", Value: "2"},
		}}
		got, err := renderEitherExpr(e)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `(tanglCompare(ctx, "a", "=", "1") || tanglCompare(ctx, "b", "=", "2"))`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("error from renderExpr propagates", func(t *testing.T) {
		e := ast.Either{Terms: []ast.Expr{nil}}
		_, err := renderEitherExpr(e)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
