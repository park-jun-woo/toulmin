//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestHasApprovalRole — tests HasApprovalRole rule
package approve

import "testing"

func TestHasApprovalRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		have string
		want bool
	}{
		{"match", "finance", "finance", true},
		{"mismatch", "finance", "manager", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := &ApproverBacking{Role: tt.role, RoleFunc: testAB.RoleFunc}
			ctx := &ApprovalContext{Approver: &testApprover{Role: tt.have}}
			got, _ := HasApprovalRole(nil, ctx, ab)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
