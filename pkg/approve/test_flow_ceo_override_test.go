//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_CEOOverride — tests CEO override defeats budget frozen
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_CEOOverride(t *testing.T) {
	ab := &ApproverBacking{RoleFunc: flowAB.RoleFunc}

	g := toulmin.NewGraph("expense:finance")
	budget := g.Warrant(IsUnderBudget, nil, 1.0)
	frozen := g.Rebuttal(IsBudgetFrozen, nil, 1.0)
	ceo := g.Defeater(IsCEOOverride, ab, 1.0)
	g.Defeat(frozen, budget)
	g.Defeat(ceo, frozen)

	f := NewFlow("expense").AddStep("finance", g)

	req := &ApprovalRequest{Amount: 5000}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &testApprover{Role: "ceo"},
			Budget:   &Budget{Remaining: 10000, Frozen: true},
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved (CEO override defeats frozen)")
	}
}
