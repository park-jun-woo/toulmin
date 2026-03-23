//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsDirectManager — tests IsDirectManager rule
package approve

import "testing"

func TestIsDirectManager(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverBacking{}
	tests := []struct {
		name     string
		approver *testApprover
		want     bool
	}{
		{"is manager", &testApprover{ID: "mgr-1"}, true},
		{"not manager", &testApprover{ID: "mgr-2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{RequesterID: "emp-1"}
			ctx := &ApprovalContext{Approver: tt.approver, ApproverID: tt.approver.ID, OrgTree: org}
			got, _ := IsDirectManager(req, ctx, ab)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
