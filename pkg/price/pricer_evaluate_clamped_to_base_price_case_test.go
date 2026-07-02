//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateClampedToBasePriceCase — verifies Pricer.Evaluate clamps the total discount to the request's base price when combined discounts exceed it
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateClampedToBasePriceCase(t *testing.T) {
	g := toulmin.NewGraph("test:overbase")
	g.Rule(HasCoupon).With(&DiscountSpec{Name: "A", Fixed: 80000})
	g.Rule(HasActivePromotion).With(&DiscountSpec{Name: "bf", Fixed: 80000})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		Coupons:    []Coupon{{Code: "A", MinPrice: 0}},
		Promotions: []Promotion{{Name: "bf", Active: true}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalDiscount != req.BasePrice {
		t.Errorf("expected discount clamped to base price %f, got %f", req.BasePrice, result.TotalDiscount)
	}
	if result.FinalPrice != 0 {
		t.Errorf("expected final price 0, got %f", result.FinalPrice)
	}
}
