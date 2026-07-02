//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestRenderInlineRule — tests renderInlineRule for successful render and error propagation from renderExpr
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderInlineRule(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		var w strings.Builder
		r := ast.InlineRule{Name: "myRule", Cond: ast.Compare{Field: "a", Op: "=", Value: "1"}}
		err := renderInlineRule(&w, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, "func myRule(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {") {
			t.Errorf("missing function signature: %q", out)
		}
		if !strings.Contains(out, `return tanglCompare(ctx, "a", "=", "1"), nil`) {
			t.Errorf("missing return statement: %q", out)
		}
	})

	t.Run("error from renderExpr propagates with rule name", func(t *testing.T) {
		var w strings.Builder
		r := ast.InlineRule{Name: "badRule", Cond: nil}
		err := renderInlineRule(&w, r)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "badRule") {
			t.Errorf("expected error to mention rule name, got %v", err)
		}
	})
}
