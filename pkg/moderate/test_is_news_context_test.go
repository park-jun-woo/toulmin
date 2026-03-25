//ff:func feature=moderate type=rule control=iteration dimension=1
//ff:what TestIsNewsContext — tests IsNewsContext rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsNewsContext(t *testing.T) {
	tests := []struct {
		name   string
		chType string
		want   bool
	}{
		{"news", "news", true},
		{"general", "general", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("channel", &Channel{Type: tt.chType})
			got, _ := IsNewsContext(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
