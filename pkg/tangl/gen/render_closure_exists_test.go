//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderClosureExists — tests exists emission with const/binder sources and errors
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderClosureExists(t *testing.T) {
	inner := ast.Compare{Lhs: ast.Path{Steps: []string{"c"}}, Op: "is not empty"}

	t.Run("const source", func(t *testing.T) {
		e := ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: inner}
		got, err := renderClosureExists(e, nil, 0)
		if err != nil || !strings.Contains(got, "tanglAny(sv, sok, ") ||
			!strings.Contains(got, `"any of coupons as c satisfies"`) ||
			!strings.Contains(got, "func(e0 any, i0 int) (bool, any)") {
			t.Errorf("const-source exists: %q (%v)", got, err)
		}
	})
	t.Run("binder source non-const expected", func(t *testing.T) {
		env := map[string]binderInfo{"o": {Elem: "e0", Idx: "i0", Src: reprCode{Expr: `"orders"`, Const: true, Val: "orders"}}}
		e := ast.Exists{Source: ast.Path{Steps: []string{"o", "items"}}, Binder: "c", Cond: inner}
		got, err := renderClosureExists(e, env, 1)
		if err != nil || !strings.Contains(got, `"any of "+`) ||
			!strings.Contains(got, "func(e1 any, i1 int) (bool, any)") {
			t.Errorf("binder-source exists: %q (%v)", got, err)
		}
	})
	t.Run("errors", func(t *testing.T) {
		bad := struct{ ast.Path }{}
		e := ast.Exists{Source: bad, Binder: "c", Cond: inner}
		if _, err := renderClosureExists(e, nil, 0); err == nil {
			t.Error("want source render error")
		}
		e = ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: nil}
		if _, err := renderClosureExists(e, nil, 0); err == nil {
			t.Error("want cond render error")
		}
	})
}
