//ff:func feature=approve type=model control=iteration dimension=1
//ff:what TestBuildApproveContext — buildApproveContext copies every request/context field into a toulmin.Context
package approve

import "testing"

// TestBuildApproveContext covers buildApproveContext's single straight-line
// branch by verifying every request and approval-context field is copied
// into the returned toulmin.Context under its documented key.
func TestBuildApproveContext(t *testing.T) {
	req := &ApprovalRequest{
		Amount:      1234.5,
		Category:    "travel",
		RequesterID: "u1",
		Description: "flight tickets",
		Metadata:    map[string]any{"trip": "seoul"},
	}
	tree := &mockOrgTree{managers: map[string]string{"u1": "m1"}}
	ac := &ApprovalContext{
		Approver:      "some-approver",
		ApproverID:    "m1",
		ApproverRole:  "manager",
		ApproverLevel: 2,
		Budget:        &Budget{Remaining: 500, Frozen: false},
		OrgTree:       tree,
		Metadata:      map[string]any{"region": "kr"},
	}

	ctx := buildApproveContext(req, ac)
	if ctx == nil {
		t.Fatal("buildApproveContext returned nil")
	}

	cases := []struct {
		key  string
		want any
	}{
		{"amount", req.Amount},
		{"category", req.Category},
		{"requesterID", req.RequesterID},
		{"description", req.Description},
		{"requestMetadata", req.Metadata},
		{"approver", ac.Approver},
		{"approverID", ac.ApproverID},
		{"approverRole", ac.ApproverRole},
		{"approverLevel", ac.ApproverLevel},
		{"budget", ac.Budget},
		{"orgTree", ac.OrgTree},
		{"metadata", ac.Metadata},
	}
	for _, c := range cases {
		got, ok := ctx.Get(c.key)
		if !ok {
			t.Errorf("Get(%q): not found", c.key)
			continue
		}
		if !equalAny(got, c.want) {
			t.Errorf("Get(%q) = %v, want %v", c.key, got, c.want)
		}
	}
}
