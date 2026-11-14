//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestDemoHasCouponParity — the exists judgment expression matches price.HasCoupon verdicts
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/price"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// TestDemoHasCouponParity rewrites the C-group demo predicate
// price.HasCoupon (pkg/price/rule_has_coupon.go) as the judgment expression
// `any of coupons as c satisfies c.minPrice is at most basePrice` and checks
// the eval-compiled closure returns the same bool as the hand-written Go
// predicate across hit, miss, empty, and wrong-typed contexts.
func TestDemoHasCouponParity(t *testing.T) {
	cond := ast.Exists{
		Source: ast.Path{Steps: []string{"coupons"}},
		Binder: "c",
		Cond: ast.Compare{
			Lhs: ast.Path{Steps: []string{"c", "minPrice"}},
			Op:  "is at most",
			Rhs: ast.Path{Steps: []string{"basePrice"}},
		},
	}
	fn, err := Compile(ast.InlineRule{Name: "has coupon", Cond: cond}, nil)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	specs := toulmin.Specs{&price.DiscountSpec{Name: "basic"}}
	cases := []struct {
		name    string
		coupons any
		base    any
	}{
		{"second coupon applies", []price.Coupon{{Code: "A", MinPrice: 100}, {Code: "B", MinPrice: 40}}, 50.0},
		{"none applies", []price.Coupon{{Code: "A", MinPrice: 100}}, 50.0},
		{"empty slice", []price.Coupon{}, 50.0},
		{"wrong coupons type", "not-a-slice", 50.0},
		{"wrong base type", []price.Coupon{{Code: "A", MinPrice: 10}}, "cheap"},
	}
	for _, tc := range cases {
		ctx := toulmin.NewContext()
		ctx.Set("coupons", tc.coupons)
		ctx.Set("basePrice", tc.base)
		got, _ := fn(ctx, specs)
		want, _ := price.HasCoupon(ctx, specs)
		if got != want {
			t.Errorf("%s: eval=%t, hand-written=%t", tc.name, got, want)
		}
	}
}
