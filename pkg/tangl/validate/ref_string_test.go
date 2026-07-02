//ff:func feature=tangl type=util control=sequence
//ff:what TestRefString — tests refString for nil, no-alias, and alias-present branches
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRefString(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		if got := refString(nil); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("NoAlias", func(t *testing.T) {
		r := &ast.Ref{Alias: "", Name: "fn"}
		if got := refString(r); got != "fn" {
			t.Fatalf("expected %q, got %q", "fn", got)
		}
	})

	t.Run("WithAlias", func(t *testing.T) {
		r := &ast.Ref{Alias: "pkg", Name: "fn"}
		if got := refString(r); got != "pkg.fn" {
			t.Fatalf("expected %q, got %q", "pkg.fn", got)
		}
	})
}
