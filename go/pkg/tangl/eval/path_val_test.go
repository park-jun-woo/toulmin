//ff:func feature=tangl type=engine control=sequence
//ff:what TestPathVal — tests Path resolution over struct field, map key, binder, and misses
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestPathVal(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("author", &testAuthor{TrustScore: 0.8})
	ctx.Set("attributes", map[string]any{"internal": true})
	sc := newScope(ctx, nil)

	t.Run("StructFieldOnPointer", func(t *testing.T) {
		v, ok := pathVal(ast.Path{Steps: []string{"author", "trustScore"}}, sc)
		if !ok || v != 0.8 {
			t.Fatalf("want 0.8,true got %v,%v", v, ok)
		}
	})
	t.Run("MapKey", func(t *testing.T) {
		v, ok := pathVal(ast.Path{Steps: []string{"attributes", "internal"}}, sc)
		if !ok || v != true {
			t.Fatalf("want true,true got %v,%v", v, ok)
		}
	})
	t.Run("EmptySteps", func(t *testing.T) {
		if _, ok := pathVal(ast.Path{}, sc); ok {
			t.Fatal("expected ok=false for empty steps")
		}
	})
	t.Run("MissingCtxKey", func(t *testing.T) {
		if _, ok := pathVal(ast.Path{Steps: []string{"nope"}}, sc); ok {
			t.Fatal("expected ok=false for missing ctx key")
		}
	})
	t.Run("MissingField", func(t *testing.T) {
		if _, ok := pathVal(ast.Path{Steps: []string{"author", "ghost"}}, sc); ok {
			t.Fatal("expected ok=false for missing field")
		}
	})
	t.Run("BinderFirstStep", func(t *testing.T) {
		sc.binders["c"] = &binding{source: ast.Path{Steps: []string{"coupons"}}, index: 0, value: testCoupon{MinPrice: 42}}
		v, ok := pathVal(ast.Path{Steps: []string{"c", "minPrice"}}, sc)
		if !ok || v != float64(42) {
			t.Fatalf("want 42,true got %v,%v", v, ok)
		}
		delete(sc.binders, "c")
	})
}
