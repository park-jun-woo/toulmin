//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestDemoActivePromotionParity — the exists judgment expression matches price.HasActivePromotion verdicts
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/price"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// TestDemoActivePromotionParity rewrites the C-group demo predicate
// price.HasActivePromotion (pkg/price/rule_has_active_promotion.go) as
// `any of promotions as p satisfies p.name equals spec DiscountSpec.Name and
// p.active equals true` and checks the eval-compiled closure returns the same
// bool as the hand-written Go predicate.
func TestDemoActivePromotionParity(t *testing.T) {
	cond := ast.Exists{
		Source: ast.Path{Steps: []string{"promotions"}},
		Binder: "p",
		Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{
				Lhs: ast.Path{Steps: []string{"p", "name"}},
				Op:  "equals",
				Rhs: ast.SpecRef{Steps: []string{"DiscountSpec", "Name"}},
			},
			ast.Compare{Lhs: ast.Path{Steps: []string{"p", "active"}}, Op: "equals", Value: "true"},
		}},
	}
	fn, err := Compile(ast.InlineRule{Name: "active promotion", Cond: cond}, nil)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	specs := toulmin.Specs{&price.DiscountSpec{Name: "SAVE30"}}
	cases := []struct {
		name       string
		promotions any
	}{
		{"named active promotion", []price.Promotion{{Name: "other", Active: true}, {Name: "SAVE30", Active: true}}},
		{"named but inactive", []price.Promotion{{Name: "SAVE30", Active: false}}},
		{"active but other name", []price.Promotion{{Name: "other", Active: true}}},
		{"empty", []price.Promotion{}},
		{"wrong type", 42},
	}
	for _, tc := range cases {
		ctx := toulmin.NewContext()
		ctx.Set("promotions", tc.promotions)
		got, _ := fn(ctx, specs)
		want, _ := price.HasActivePromotion(ctx, specs)
		if got != want {
			t.Errorf("%s: eval=%t, hand-written=%t", tc.name, got, want)
		}
	}
}
