//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureCompare — tests compare-leaf emission across lhs/rhs operand forms
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderClosureCompare(t *testing.T) {
	env := map[string]binderInfo{"c": {Elem: "e0", Idx: "i0", Src: reprCode{Expr: `"coupons"`, Const: true, Val: "coupons"}}}
	bad := struct{ ast.Path }{}

	t.Run("field vs literal", func(t *testing.T) {
		got, err := renderClosureCompare(ast.Compare{Field: "x", Op: "equals", Value: "1"}, nil)
		if err != nil || !strings.Contains(got, `tanglCtxVal(ctx, "x")`) ||
			!strings.Contains(got, `Where: "x"`) || !strings.Contains(got, `Expected: "equals 1"`) {
			t.Errorf("field-vs-literal closure: %q (%v)", got, err)
		}
	})
	t.Run("bare op without value", func(t *testing.T) {
		got, err := renderClosureCompare(ast.Compare{Field: "x", Op: "is empty"}, nil)
		if err != nil || !strings.Contains(got, `Expected: "is empty"`) {
			t.Errorf("bare-op closure: %q (%v)", got, err)
		}
	})
	t.Run("lhs and const rhs operands", func(t *testing.T) {
		c := ast.Compare{Lhs: ast.SpecRef{Steps: []string{"d", "max"}}, Op: "equals", Rhs: ast.Path{Steps: []string{"y"}}}
		got, err := renderClosureCompare(c, nil)
		if err != nil || !strings.Contains(got, `Where: "d.max"`) || !strings.Contains(got, `Expected: "equals y"`) {
			t.Errorf("operand closure: %q (%v)", got, err)
		}
	})
	t.Run("non-const rhs repr", func(t *testing.T) {
		c := ast.Compare{Field: "x", Op: "equals", Rhs: ast.Path{Steps: []string{"c", "minPrice"}}}
		got, err := renderClosureCompare(c, env)
		if err != nil || !strings.Contains(got, `Expected: "equals "+`) {
			t.Errorf("non-const rhs closure: %q (%v)", got, err)
		}
	})
	t.Run("operand errors", func(t *testing.T) {
		if _, err := renderClosureCompare(ast.Compare{Lhs: bad, Op: "is empty"}, nil); err == nil {
			t.Error("want lhs operand error")
		}
		if _, err := renderClosureCompare(ast.Compare{Field: "x", Op: "equals", Rhs: bad}, nil); err == nil {
			t.Error("want rhs operand error")
		}
	})
}
