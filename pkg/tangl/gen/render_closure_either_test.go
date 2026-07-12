//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureEither — tests either-group emission via tanglOr and a term error
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderClosureEither(t *testing.T) {
	cmp := ast.Compare{Field: "x", Op: "equals", Value: "1"}
	got, err := renderClosureEither(ast.Either{Terms: []ast.Expr{cmp, cmp}}, nil, 0)
	if err != nil || !strings.Contains(got, "return tanglOr(") || strings.Count(got, "tanglCompareVal(") != 2 {
		t.Errorf("either closure: %q (%v)", got, err)
	}
	if _, err := renderClosureEither(ast.Either{Terms: []ast.Expr{nil}}, nil, 0); err == nil {
		t.Error("want term render error")
	}
}
