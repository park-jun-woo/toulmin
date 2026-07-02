//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderExpr — tests renderExpr for Compare, Not (success/error), Logic, Either, and unsupported default branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderExpr(t *testing.T) {
	t.Run("compare leaf", func(t *testing.T) {
		got, err := renderExpr(ast.Compare{Field: "a", Op: "=", Value: "1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `tanglCompare(ctx, "a", "=", "1")`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("not with success inner", func(t *testing.T) {
		got, err := renderExpr(ast.Not{Term: ast.Compare{Field: "a", Op: "=", Value: "1"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `!(tanglCompare(ctx, "a", "=", "1"))`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("not with error inner propagates", func(t *testing.T) {
		_, err := renderExpr(ast.Not{Term: nil})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("logic node", func(t *testing.T) {
		got, err := renderExpr(ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Field: "a", Op: "=", Value: "1"},
			ast.Compare{Field: "b", Op: "=", Value: "2"},
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == "" {
			t.Error("expected non-empty logic expression")
		}
	})

	t.Run("either node", func(t *testing.T) {
		got, err := renderExpr(ast.Either{Terms: []ast.Expr{
			ast.Compare{Field: "a", Op: "=", Value: "1"},
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == "" {
			t.Error("expected non-empty either expression")
		}
	})

	t.Run("unsupported node type returns error", func(t *testing.T) {
		_, err := renderExpr(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
