//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateNonDiscountEvidenceCase — verifies Pricer.Evaluate ignores evidence that is not a *DiscountSpec
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateNonDiscountEvidenceCase(t *testing.T) {
	// nonDiscountRule is an always-active rule whose evidence is not a *DiscountSpec,
	// exercising the `!ok` side of `db, ok := r.Evidence.(*DiscountSpec)`.
	nonDiscountRule := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
		return true, "not-a-discount-spec"
	}

	g := toulmin.NewGraph("test:nondiscount")
	g.Rule(nonDiscountRule)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{}

	result, err := p.Evaluate(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.AppliedDiscounts) != 0 {
		t.Errorf("expected no applied discounts, got %d", len(result.AppliedDiscounts))
	}
	if result.TotalDiscount != 0 {
		t.Errorf("expected 0 discount, got %f", result.TotalDiscount)
	}
}
