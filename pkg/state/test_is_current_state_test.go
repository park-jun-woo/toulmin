//ff:func feature=state type=rule control=iteration dimension=1
//ff:what TestIsCurrentState — tests IsCurrentState rule
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsCurrentState(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		current string
		want    bool
	}{
		{"match", "pending", "pending", true},
		{"mismatch", "pending", "accepted", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("from", tt.from)
			ctx.Set("currentState", tt.current)
			got, _ := IsCurrentState(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
