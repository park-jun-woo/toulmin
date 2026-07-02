//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateTotalCapNotExceededCase — verifies Pricer.Evaluate leaves the total discount unclamped when it stays below the total cap
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateTotalCapNotExceededCase(t *testing.T) {
	g := toulmin.NewGraph("test:capnotexceeded")
	g.Rule(HasCoupon).With(&DiscountSpec{Name: "A", Rate: 0.1})

	totalCap := &DiscountSpec{Max: 50000}
	p := NewPricer(g, totalCap)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: []Coupon{{Code: "A", MinPrice: 0}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-10000) > 0.01 {
		t.Errorf("expected 10000 (below cap, unclamped), got %f", result.TotalDiscount)
	}
}
