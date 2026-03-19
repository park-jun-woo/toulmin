package price

import (
	"math"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
	if math.Abs(result.FinalPrice-70000) > 0.01 {
		t.Errorf("expected 70000 final, got %f", result.FinalPrice)
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
	// 50% of 100000 = 50000, but max 30000
	if math.Abs(result.TotalDiscount-30000) > 0.01 {
		t.Errorf("expected 30000 (capped), got %f", result.TotalDiscount)
	}
}

func TestPricer_ComboWithMin(t *testing.T) {
	g := toulmin.NewGraph("test:combo")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "COMBO", Rate: 0.01, Fixed: 100, Min: 3000}, 1.0)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: []Coupon{{Code: "COMBO", MinPrice: 0}}}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1% of 100000 + 100 = 1100, but min 3000
	if math.Abs(result.TotalDiscount-3000) > 0.01 {
		t.Errorf("expected 3000 (min guarantee), got %f", result.TotalDiscount)
	}
}

func TestPricer_TotalCap(t *testing.T) {
	g := toulmin.NewGraph("test:totalcap")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "A", Rate: 0.2}, 1.0)
	g.Warrant(IsMemberLevel, &DiscountBacking{Name: "basic", Rate: 0.1}, 1.0)

	totalCap := &DiscountBacking{Max: 25000}
	p := NewPricer(g, totalCap)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		User:    &User{Membership: "basic"},
		Coupons: []Coupon{{Code: "A", MinPrice: 0}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 20% + 10% = 30000, but totalCap 25000
	if math.Abs(result.TotalDiscount-25000) > 0.01 {
		t.Errorf("expected 25000 (total cap), got %f", result.TotalDiscount)
	}
}

func TestPricer_NoDiscount(t *testing.T) {
	g := toulmin.NewGraph("test:none")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "X", Rate: 0.1}, 1.0)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Coupons: nil} // no coupons

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalDiscount != 0 {
		t.Errorf("expected 0 discount, got %f", result.TotalDiscount)
	}
	if result.FinalPrice != 100000 {
		t.Errorf("expected 100000 final, got %f", result.FinalPrice)
	}
}

func TestPricer_VIPDefeatsNoStack(t *testing.T) {
	g := toulmin.NewGraph("test:vip")
	coupon := g.Warrant(HasCoupon, &DiscountBacking{Name: "SAVE", Rate: 0.2}, 1.0)
	member := g.Warrant(IsMemberLevel, &DiscountBacking{Name: "basic", Rate: 0.1}, 1.0)
	noStack := g.Rebuttal(IsAlreadyDiscounted, nil, 1.0)
	vip := g.Defeater(IsMemberLevel, &DiscountBacking{Name: "vip"}, 1.0)
	g.Defeat(noStack, coupon)
	g.Defeat(noStack, member)
	g.Defeat(vip, noStack)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000, Metadata: map[string]any{"discounted": true}}
	ctx := &PriceContext{
		User:    &User{Membership: "vip"},
		Coupons: []Coupon{{Code: "SAVE", MinPrice: 0}},
	}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// VIP defeats noStack, so coupon still applies
	if result.TotalDiscount == 0 {
		t.Error("expected some discount (VIP defeats noStack)")
	}
}

func TestPricer_AppliedDiscountsList(t *testing.T) {
	g := toulmin.NewGraph("test:list")
	g.Warrant(HasCoupon, &DiscountBacking{Name: "A", Rate: 0.1}, 1.0)
	g.Warrant(IsMemberLevel, &DiscountBacking{Name: "basic", Rate: 0.05}, 1.0)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{
		User:    &User{Membership: "basic"},
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
