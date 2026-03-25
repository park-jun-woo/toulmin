//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestHasApprovalRole — tests HasApprovalRole rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
			ab := &ApproverBacking{Role: tt.role}
			ctx := toulmin.NewContext()
			ctx.Set("approverRole", tt.have)
			got, _ := HasApprovalRole(ctx, ab)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
