//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_AllPassed — tests approval flow when all steps pass
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_AllPassed(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverBacking{}

	g := toulmin.NewGraph("expense:manager")
	g.Rule(IsDirectManager).Backing(ab)
	g.Rule(IsUnderBudget)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{Amount: 5000, RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver:   &testApprover{ID: "mgr-1"},
			ApproverID: "mgr-1",
			Budget:     &Budget{Remaining: 10000},
			OrgTree:    org,
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Approved {
		t.Error("expected approved")
	}
}
