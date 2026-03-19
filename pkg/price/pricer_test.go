package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

var memberFunc = func(u any) string { return u.(*testUser).Membership }

func TestPricer_RateDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:rate")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "SAVE30", Rate: 0.3}, 1.0)

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

func TestPricer_FixedDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:fixed")
	g.Warrant(HasActivePromotion, &DiscountBacking{Name: "bf", Fixed: 5000}, 1.0)

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

func TestPricer_RateWithMax(t *testing.T) {
	g := toulmin.NewGraph("test:max")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "BIG", Rate: 0.5, Max: 30000}, 1.0)

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

func TestPricer_AppliedDiscountsList(t *testing.T) {
	g := toulmin.NewGraph("test:list")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "A", Rate: 0.1}, 1.0)
	g.Warrant(IsMemberLevel, &MemberBacking{Level: "basic", MembershipFunc: memberFunc, Discount: &DiscountBacking{Name: "basic", Rate: 0.05}}, 1.0)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		User:    &testUser{Membership: "basic"},
		Coupons: []Coupon{{Code: "A", MinPrice: 0}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.AppliedDiscounts) != 2 {
		t.Errorf("expected 2 applied discounts, got %d", len(result.AppliedDiscounts))
	}
}
