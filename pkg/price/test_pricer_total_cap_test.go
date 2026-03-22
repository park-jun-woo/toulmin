//ff:func feature=price type=engine control=sequence
//ff:what TestPricer_TotalCap — tests total discount cap across multiple discounts
package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPricer_TotalCap(t *testing.T) {
	g := toulmin.NewGraph("test:totalcap")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "A", Rate: 0.2}, 1.0)
	g.Warrant(IsMemberLevel, &MemberBacking{Level: "basic", MembershipFunc: memberFunc, Discount: &DiscountBacking{Name: "basic", Rate: 0.1}}, 1.0)

	totalCap := &DiscountBacking{Max: 25000}
	p := NewPricer(g, totalCap)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		User:    &testUser{Membership: "basic"},
		Coupons: []Coupon{{Code: "A", MinPrice: 0}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.TotalDiscount-25000) > 0.01 {
		t.Errorf("expected 25000 (total cap), got %f", result.TotalDiscount)
	}
}
