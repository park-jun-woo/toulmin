//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestRenderCompareExpr — tests renderCompareExpr formats a quoted tanglCompare call
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderCompareExpr(t *testing.T) {
	got := renderCompareExpr(ast.Compare{Field: "amount", Op: ">=", Value: "100"})
	want := `tanglCompare(ctx, "amount", ">=", "100")`
	if got != want {
		t.Errorf("renderCompareExpr() = %q, want %q", got, want)
	}
}
