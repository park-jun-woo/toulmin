//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_RateWithMax — tests rate discount capped by max
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_RateWithMax(t *testing.T) {
	g := toulmin.NewGraph("test:max")
	g.Rule(HasCoupon).With(&DiscountSpec{Name: "BIG", Rate: 0.5, Max: 30000})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: []Coupon{{Code: "BIG", MinPrice: 0}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-30000) > 0.01 {
		t.Errorf("expected 30000 (capped), got %f", result.TotalDiscount)
	}
}
