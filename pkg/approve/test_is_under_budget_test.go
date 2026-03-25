//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsUnderBudget — tests IsUnderBudget rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsUnderBudget(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		remaining float64
		want      bool
	}{
		{"under", 5000, 10000, true},
		{"equal", 10000, 10000, true},
		{"over", 15000, 10000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("amount", tt.amount)
			ctx.Set("budget", &Budget{Remaining: tt.remaining})
			got, _ := IsUnderBudget(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
