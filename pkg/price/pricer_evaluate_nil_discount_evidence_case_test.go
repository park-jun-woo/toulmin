//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateNilDiscountEvidenceCase — verifies Pricer.Evaluate treats a matched rule with nil Discount evidence as no discount applied
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateNilDiscountEvidenceCase(t *testing.T) {
	g := toulmin.NewGraph("test:nildiscount")
	g.Rule(IsMemberLevel).With(&MemberSpec{Level: "basic", Discount: nil})

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{Membership: "basic"}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.AppliedDiscounts) != 0 {
		t.Errorf("expected no applied discounts for nil Discount evidence, got %d", len(result.AppliedDiscounts))
	}
	if result.TotalDiscount != 0 {
		t.Errorf("expected 0 discount, got %f", result.TotalDiscount)
	}
}
