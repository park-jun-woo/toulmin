//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestIsBulkOrder — tests IsBulkOrder rule
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsBulkOrder(t *testing.T) {
	tests := []struct {
		name string
		qty  int
		min  int
		want bool
	}{
		{"bulk", 100, 50, true},
		{"equal", 50, 50, true},
		{"not bulk", 10, 50, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("quantity", tt.qty)
			got, _ := IsBulkOrder(ctx, &BulkOrderBacking{MinQuantity: tt.min})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
