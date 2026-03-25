//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_RateDiscount — tests rate-based coupon discount
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_RateDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:rate")
	g.Rule(HasCoupon).Backing(&DiscountBacking{Name: "SAVE30", Rate: 0.3})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: []Coupon{{Code: "SAVE30", MinPrice: 0}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-30000) > 0.01 {
		t.Errorf("expected 30000 discount, got %f", result.TotalDiscount)
	}
}
