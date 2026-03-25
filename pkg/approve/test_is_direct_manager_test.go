//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsDirectManager — tests IsDirectManager rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsDirectManager(t *testing.T) {
	org := &mockOrgTree{managers: map[string]string{"emp-1": "mgr-1"}}
	ab := &ApproverSpec{}
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
			ctx := toulmin.NewContext()
			ctx.Set("requesterID", "emp-1")
			ctx.Set("approverID", tt.approver.ID)
			ctx.Set("orgTree", org)
			got, _ := IsDirectManager(ctx, toulmin.Specs{ab})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
