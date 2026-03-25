//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsCEOOverride — tests IsCEOOverride rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
			ctx := toulmin.NewContext()
			ctx.Set("approverRole", tt.role)
			got, _ := IsCEOOverride(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
