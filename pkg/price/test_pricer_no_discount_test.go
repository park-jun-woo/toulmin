//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_NoDiscount — tests pricer with no applicable discounts
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_NoDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:none")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "X", Rate: 0.1}, 1.0)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: nil}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalDiscount != 0 {
		t.Errorf("expected 0 discount, got %f", result.TotalDiscount)
	}
}
