//ff:func feature=moderate type=rule control=iteration dimension=1
//ff:what TestHasMinPosts — tests HasMinPosts rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasMinPosts(t *testing.T) {
	tests := []struct {
		name  string
		count int
		min   int
		want  bool
	}{
		{"enough", 100, 10, true},
		{"not enough", 5, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("author", &Author{PostCount: tt.count})
			got, _ := HasMinPosts(ctx, &MinPostsBacking{MinPosts: tt.min})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
