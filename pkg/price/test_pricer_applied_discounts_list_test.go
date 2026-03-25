//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_AppliedDiscountsList — tests applied discounts list has correct count
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_AppliedDiscountsList(t *testing.T) {
	g := toulmin.NewGraph("test:list")
	g.Rule(HasCoupon).With(&DiscountSpec{Name: "A", Rate: 0.1})
	g.Rule(IsMemberLevel).With(&MemberSpec{Level: "basic", Discount: &DiscountSpec{Name: "basic", Rate: 0.05}})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		User:       &testUser{Membership: "basic"},
		Membership: "basic",
		Coupons:    []Coupon{{Code: "A", MinPrice: 0}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.AppliedDiscounts) != 2 {
		t.Errorf("expected 2 applied discounts, got %d", len(result.AppliedDiscounts))
	}
}
