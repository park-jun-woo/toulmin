//ff:func feature=price type=engine control=sequence
//ff:what pricerEvaluateGraphErrorCase — verifies Pricer.Evaluate returns an error and nil result when the defeats graph evaluation fails (cyclic graph)
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func pricerEvaluateGraphErrorCase(t *testing.T) {
	cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
	cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }

	g := toulmin.NewGraph("test:cycle")
	ca := g.Rule(cycleA)
	cb := g.Counter(cycleB)
	cb.Attacks(ca)
	ca.Attacks(cb)

	p := NewPricer(g, nil)
	req := &PurchaseRequest{BasePrice: 100000}
	ctx := &PriceContext{}

	result, err := p.Evaluate(req, ctx)
	if err == nil {
		t.Fatal("expected error from cyclic graph evaluation")
	}
	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}
