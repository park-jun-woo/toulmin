//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_FixedDiscount — tests fixed amount promotion discount
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_FixedDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:fixed")
	g.Rule(HasActivePromotion).With(&DiscountSpec{Name: "bf", Fixed: 5000})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Promotions: []Promotion{{Name: "bf", Active: true}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-5000) > 0.01 {
		t.Errorf("expected 5000 discount, got %f", result.TotalDiscount)
	}
}
