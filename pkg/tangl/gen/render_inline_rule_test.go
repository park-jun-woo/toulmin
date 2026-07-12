//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderInlineRule — tests renderInlineRule closure emission and error propagation
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderInlineRule(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		var w strings.Builder
		r := ast.InlineRule{Name: "my rule", Cond: ast.Compare{Field: "a", Op: "equals", Value: "1"}}
		err := renderInlineRule(&w, r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, "func myRule(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {") {
			t.Errorf("missing function signature: %q", out)
		}
		if !strings.Contains(out, `tanglCtxVal(ctx, "a")`) {
			t.Errorf("missing ctx value resolution: %q", out)
		}
		if !strings.Contains(out, `tanglCompareVal(lv, lok, "equals", "1")`) {
			t.Errorf("missing value comparison: %q", out)
		}
		if !strings.Contains(out, `toulmin.Fact{Where: "a", Expected: "equals 1", Actual: lv}`) {
			t.Errorf("missing witness Fact: %q", out)
		}
	})

	t.Run("error from renderClosureExpr propagates with rule name", func(t *testing.T) {
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
