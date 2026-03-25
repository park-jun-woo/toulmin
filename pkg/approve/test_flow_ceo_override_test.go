//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_CEOOverride — tests CEO override defeats budget frozen
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_CEOOverride(t *testing.T) {
	ab := &ApproverBacking{}

	g := toulmin.NewGraph("expense:finance")
	budget := g.Rule(IsUnderBudget)
	frozen := g.Counter(IsBudgetFrozen)
	ceo := g.Except(IsCEOOverride).Backing(ab)
	frozen.Attacks(budget)
	ceo.Attacks(frozen)

	f := NewFlow("expense").AddStep("finance", g)

	req := &ApprovalRequest{Amount: 5000}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver:     &testApprover{Role: "ceo"},
			ApproverRole: "ceo",
			Budget:       &Budget{Remaining: 10000, Frozen: true},
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved (CEO override defeats frozen)")
	}
}
