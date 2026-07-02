//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateTotalCapMaxNotPositiveCase — verifies Pricer.Evaluate does not clamp the total discount when the total cap's Max is not positive
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateTotalCapMaxNotPositiveCase(t *testing.T) {
	g := toulmin.NewGraph("test:capzero")
	g.Rule(HasCoupon).With(&DiscountSpec{Name: "A", Rate: 0.2})

	totalCap := &DiscountSpec{Max: 0}
	p := NewPricer(g, totalCap)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: []Coupon{{Code: "A", MinPrice: 0}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-20000) > 0.01 {
		t.Errorf("expected 20000 (cap not applied since Max<=0), got %f", result.TotalDiscount)
	}
}
