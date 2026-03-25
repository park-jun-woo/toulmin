//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsSmallAmount — tests IsSmallAmount rule
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsSmallAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		threshold float64
		want      bool
	}{
		{"small", 5000, 10000, true},
		{"equal", 10000, 10000, true},
		{"large", 15000, 10000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("amount", tt.amount)
			got, _ := IsSmallAmount(ctx, &ThresholdBacking{Max: tt.threshold})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
