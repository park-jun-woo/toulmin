package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_AllPassed(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}

	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsDirectManager, nil, 1.0)
	g.Warrant(IsUnderBudget, nil, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 5000, RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &Approver{ID: "mgr-1"},
			Budget:   &Budget{Remaining: 10000},
			OrgTree:  org,
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved")
	}
}

func TestFlow_StepRejected(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}

	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsDirectManager, nil, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 15000, RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &Approver{ID: "mgr-2"}, // not the manager
			Budget:   &Budget{Remaining: 10000},
			OrgTree:  org,
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Approved {
		t.Error("expected rejected (not direct manager)")
	}
}

func TestFlow_FrozenWithUrgentDefeat(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}

	g := toulmin.NewGraph("expense:manager")
	budget := g.Warrant(IsUnderBudget, nil, 1.0)
	frozen := g.Rebuttal(IsBudgetFrozen, nil, 1.0)
	urgent := g.Defeater(IsUrgent, nil, 1.0)
	g.Defeat(frozen, budget)
	g.Defeat(urgent, frozen)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{
		Amount:   5000,
		Metadata: map[string]any{"urgent": true},
	}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &Approver{ID: "mgr-1"},
			Budget:   &Budget{Remaining: 10000, Frozen: true},
			OrgTree:  org,
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved (urgent defeats frozen)")
	}
}

func TestFlow_MultiStep_SecondRejects(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}

	g1 := toulmin.NewGraph("expense:manager")
	g1.Warrant(IsDirectManager, nil, 1.0)

	g2 := toulmin.NewGraph("expense:finance")
	g2.Warrant(HasApprovalRole, "finance", 1.0)

	f := NewFlow("expense").
		AddStep("manager", g1).
		AddStep("finance", g2)

	req := &ApprovalRequest{RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		switch step {
		case "manager":
			return &ApprovalContext{Approver: &Approver{ID: "mgr-1"}, OrgTree: org}
		case "finance":
			return &ApprovalContext{Approver: &Approver{Role: "manager"}} // not finance
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Approved {
		t.Error("expected rejected at finance step")
	}
	if len(result.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(result.Steps))
	}
	if result.Steps[0].Verdict <= 0 {
		t.Error("manager step should have passed")
	}
}

func TestFlow_CEOOverride(t *testing.T) {
	g := toulmin.NewGraph("expense:finance")
	budget := g.Warrant(IsUnderBudget, nil, 1.0)
	frozen := g.Rebuttal(IsBudgetFrozen, nil, 1.0)
	ceo := g.Defeater(IsCEOOverride, nil, 1.0)
	g.Defeat(frozen, budget)
	g.Defeat(ceo, frozen)

	f := NewFlow("expense").AddStep("finance", g)

	req := &ApprovalRequest{Amount: 5000}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &Approver{Role: "ceo"},
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

func TestFlow_TraceIncluded(t *testing.T) {
	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsUnderBudget, nil, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 5000}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{Budget: &Budget{Remaining: 10000}}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Steps) == 0 {
		t.Fatal("expected steps")
	}
	if len(result.Steps[0].Trace) == 0 {
		t.Error("expected non-empty trace")
	}
}
