//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestIsLegacyBrowser — tests IsLegacyBrowser rule
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsLegacyBrowser(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"legacy", map[string]any{"legacy_browser": true}, true},
		{"not legacy", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("attributes", tt.attr)
			got, _ := IsLegacyBrowser(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
