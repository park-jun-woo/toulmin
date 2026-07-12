//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureLogic — tests and/or helper selection and a term error
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderClosureLogic(t *testing.T) {
	cmp := ast.Compare{Field: "x", Op: "equals", Value: "1"}

	got, err := renderClosureLogic(ast.Logic{Op: "and", Terms: []ast.Expr{cmp, cmp}}, nil, 0)
	if err != nil || !strings.Contains(got, "return tanglAnd(") || strings.Count(got, "tanglCompareVal(") != 2 {
		t.Errorf("and closure: %q (%v)", got, err)
	}
	got, err = renderClosureLogic(ast.Logic{Op: "or", Terms: []ast.Expr{cmp}}, nil, 0)
	if err != nil || !strings.Contains(got, "return tanglOr(") {
		t.Errorf("or closure: %q (%v)", got, err)
	}
	if _, err := renderClosureLogic(ast.Logic{Op: "and", Terms: []ast.Expr{nil}}, nil, 0); err == nil {
		t.Error("want term render error")
	}
}
