//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestIsRegion — tests IsRegion rule
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsRegion(t *testing.T) {
	tests := []struct {
		name   string
		region string
		back   string
		want   bool
	}{
		{"match", "KR", "KR", true},
		{"mismatch", "US", "KR", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("region", tt.region)
			got, _ := IsRegion(ctx, &RegionBacking{Region: tt.back})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
