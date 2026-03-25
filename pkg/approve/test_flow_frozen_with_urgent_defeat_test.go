//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_FrozenWithUrgentDefeat — tests urgent defeater overrides budget frozen
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_FrozenWithUrgentDefeat(t *testing.T) {
	g := toulmin.NewGraph("expense:manager")
	budget := g.Rule(IsUnderBudget)
	frozen := g.Counter(IsBudgetFrozen)
	urgent := g.Except(IsUrgent)
	frozen.Attacks(budget)
	urgent.Attacks(frozen)

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
