//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileExists — tests exists found/empty witnesses and binder-indexed Where
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCompileExists(t *testing.T) {
	empty := fakeResolver{}

	t.Run("Found", func(t *testing.T) {
		c := testCtx()
		c.Set("coupons", []testCoupon{{MinPrice: 100}, {MinPrice: 50}})
		c.Set("basePrice", 75)
		cond := ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "c",
			Cond: ast.Compare{
				Lhs: ast.Path{Steps: []string{"c", "minPrice"}},
				Op:  "is at most",
				Rhs: ast.Path{Steps: []string{"basePrice"}},
			},
		}
		b, ev := runRule(t, cond, empty, c)
		f := asFact(t, ev)
		if !b || f.Where != "coupons[1].minPrice" || f.Expected != "is at most basePrice" || f.Actual != float64(50) {
			t.Fatalf("exists-found witness mismatch: %v %+v", b, f)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		c := testCtx()
		c.Set("coupons", []testCoupon{})
		cond := ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "c",
			Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"c", "minPrice"}}, Op: "is at most", Value: "1"},
		}
		b, ev := runRule(t, cond, empty, c)
		f := asFact(t, ev)
		if b || f.Where != "coupons" || f.Expected != "any of coupons as c satisfies" || f.Actual != 0 {
			t.Fatalf("exists-empty witness mismatch: %v %+v", b, f)
		}
	})
	t.Run("CompileErrors", func(t *testing.T) {
		bad := ast.Exists{Source: ast.Call{Fn: ast.Ref{Name: "nope"}}, Binder: "c", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}}
		if _, err := Compile(ast.InlineRule{Name: "r", Cond: bad}, empty); err == nil {
			t.Fatal("want unresolved-source compile error")
		}
		bad.Source = ast.Path{Steps: []string{"xs"}}
		bad.Cond = nil
		if _, err := Compile(ast.InlineRule{Name: "r", Cond: bad}, empty); err == nil {
			t.Fatal("want nil-cond compile error")
		}
	})
	t.Run("MissingSource", func(t *testing.T) {
		cond := ast.Exists{
			Source: ast.Path{Steps: []string{"ghost"}},
			Binder: "c",
			Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"c"}}, Op: "is not empty"},
		}
		b, ev := runRule(t, cond, empty, testCtx())
		if b || asFact(t, ev).Actual != 0 {
			t.Fatalf("missing source should be zero elements, got %v %+v", b, ev)
		}
	})
	t.Run("NestedBinderRestored", func(t *testing.T) {
		c := testCtx()
		c.Set("xs", []int{1})
		inner := func(v string) ast.Exists {
			return ast.Exists{
				Source: ast.Path{Steps: []string{"xs"}},
				Binder: "e",
				Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"e"}}, Op: "equals", Value: v},
			}
		}
		if b, _ := runRule(t, ast.Exists{Source: ast.Path{Steps: []string{"xs"}}, Binder: "e", Cond: inner("1")}, empty, c); !b {
			t.Fatal("nested exists sharing a binder should find the element")
		}
		if b, _ := runRule(t, ast.Exists{Source: ast.Path{Steps: []string{"xs"}}, Binder: "e", Cond: inner("2")}, empty, c); b {
			t.Fatal("nested exists sharing a binder should miss")
		}
	})
	t.Run("NonSlice", func(t *testing.T) {
		c := testCtx()
		c.Set("coupons", "not a slice")
		cond := ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "c",
			Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"c"}}, Op: "is not empty", Value: ""},
		}
		b, ev := runRule(t, cond, empty, c)
		if b || asFact(t, ev).Actual != 0 {
			t.Fatalf("non-slice source should be zero elements, got %v %+v", b, ev)
		}
	})
}
