package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

var flowAB = &ApproverBacking{
	IDFunc:    func(a any) string { return a.(*testApprover).ID },
	RoleFunc:  func(a any) string { return a.(*testApprover).Role },
	LevelFunc: func(a any) int { return a.(*testApprover).Level },
}

func TestFlow_AllPassed(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverBacking{IDFunc: flowAB.IDFunc}

	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsDirectManager, ab, 1.0)
	g.Warrant(IsUnderBudget, nil, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 5000, RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &testApprover{ID: "mgr-1"},
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
	ab := &ApproverBacking{IDFunc: flowAB.IDFunc}

	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsDirectManager, ab, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver: &testApprover{ID: "mgr-2"},
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

func TestFlow_MultiStep_SecondRejects(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverBacking{IDFunc: flowAB.IDFunc, RoleFunc: flowAB.RoleFunc}

	g1 := toulmin.NewGraph("expense:manager")
	g1.Warrant(IsDirectManager, ab, 1.0)

	g2 := toulmin.NewGraph("expense:finance")
	g2.Warrant(HasApprovalRole, &ApproverBacking{Role: "finance", RoleFunc: flowAB.RoleFunc}, 1.0)

	f := NewFlow("expense").AddStep("manager", g1).AddStep("finance", g2)

	req := &ApprovalRequest{RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		switch step {
		case "manager":
			return &ApprovalContext{Approver: &testApprover{ID: "mgr-1"}, OrgTree: org}
		case "finance":
			return &ApprovalContext{Approver: &testApprover{Role: "manager"}}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Approved {
		t.Error("expected rejected at finance step")
	}
}

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
