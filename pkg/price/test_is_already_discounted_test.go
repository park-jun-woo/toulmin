//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestIsAlreadyDiscounted — tests IsAlreadyDiscounted rule
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsAlreadyDiscounted(t *testing.T) {
	tests := []struct {
		name string
		meta map[string]any
		want bool
	}{
		{"discounted", map[string]any{"discounted": true}, true},
		{"not discounted", map[string]any{}, false},
		{"nil meta", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("requestMetadata", tt.meta)
			got, _ := IsAlreadyDiscounted(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
