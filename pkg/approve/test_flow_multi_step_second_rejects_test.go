//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_MultiStep_SecondRejects — tests multi-step flow where second step rejects
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
