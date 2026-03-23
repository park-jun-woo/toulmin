//ff:func feature=approve type=engine control=sequence
//ff:what TestFlow_StepRejected — tests approval flow when step is rejected
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlow_StepRejected(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverBacking{}

	g := toulmin.NewGraph("expense:manager")
	g.Warrant(IsDirectManager, ab, 1.0)

	f := NewFlow("expense").AddStep("manager", g)

	req := &ApprovalRequest{RequesterID: "emp-1"}
	result, err := f.Evaluate(req, func(step string) *ApprovalContext {
		return &ApprovalContext{
			Approver:   &testApprover{ID: "mgr-2"},
			ApproverID: "mgr-2",
			OrgTree:    org,
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Approved {
		t.Error("expected rejected (not direct manager)")
	}
}
