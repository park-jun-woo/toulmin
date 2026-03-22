//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_FrozenWithUrgentDefeat — tests urgent defeater overrides budget frozen
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_FrozenWithUrgentDefeat(t *testing.T) {
	g := toulmin.NewGraph("expense:manager")
	budget := g.Warrant(IsUnderBudget, nil, 1.0)
	frozen := g.Rebuttal(IsBudgetFrozen, nil, 1.0)
	urgent := g.Defeater(IsUrgent, nil, 1.0)
	g.Defeat(frozen, budget)
	g.Defeat(urgent, frozen)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 5000, Metadata: map[string]any{"urgent": true}}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{Budget: &Budget{Remaining: 10000, Frozen: true}}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved (urgent defeats frozen)")
	}
}
