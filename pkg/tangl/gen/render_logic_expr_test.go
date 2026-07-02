//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderLogicExpr — tests renderLogicExpr for and/or separator branches and error propagation
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderLogicExpr(t *testing.T) {
	t.Run("default op uses and separator", func(t *testing.T) {
		l := ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Field: "a", Op: "=", Value: "1"},
			ast.Compare{Field: "b", Op: "=", Value: "2"},
		}}
		got, err := renderLogicExpr(l)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `(tanglCompare(ctx, "a", "=", "1") && tanglCompare(ctx, "b", "=", "2"))`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("or op uses or separator", func(t *testing.T) {
		l := ast.Logic{Op: "or", Terms: []ast.Expr{
			ast.Compare{Field: "a", Op: "=", Value: "1"},
			ast.Compare{Field: "b", Op: "=", Value: "2"},
		}}
		got, err := renderLogicExpr(l)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `(tanglCompare(ctx, "a", "=", "1") || tanglCompare(ctx, "b", "=", "2"))`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("error from renderExpr propagates", func(t *testing.T) {
		l := ast.Logic{Op: "and", Terms: []ast.Expr{nil}}
		_, err := renderLogicExpr(l)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
