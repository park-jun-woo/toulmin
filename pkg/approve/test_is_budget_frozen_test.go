//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsBudgetFrozen — tests IsBudgetFrozen rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsBudgetFrozen(t *testing.T) {
	tests := []struct {
		name   string
		frozen bool
		want   bool
	}{
		{"frozen", true, true},
		{"not frozen", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("budget", &Budget{Frozen: tt.frozen})
			got, _ := IsBudgetFrozen(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
