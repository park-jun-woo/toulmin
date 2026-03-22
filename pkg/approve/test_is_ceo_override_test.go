//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsCEOOverride — tests IsCEOOverride rule
package approve

import "testing"

func TestIsCEOOverride(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"ceo", "ceo", true},
		{"not ceo", "manager", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := &ApproverBacking{RoleFunc: testAB.RoleFunc}
			ctx := &ApprovalContext{Approver: &testApprover{Role: tt.role}}
			got, _ := IsCEOOverride(nil, ctx, ab)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
