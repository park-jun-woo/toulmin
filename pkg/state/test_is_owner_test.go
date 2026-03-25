//ff:func feature=state type=rule control=iteration dimension=1
//ff:what TestIsOwner — tests IsOwner rule
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsOwner(t *testing.T) {
	b := &OwnerBacking{}
	tests := []struct {
		name            string
		userID          string
		resourceOwnerID string
		want            bool
	}{
		{"owner", "u1", "u1", true},
		{"not owner", "u1", "u2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("userID", tt.userID)
			ctx.Set("resourceOwnerID", tt.resourceOwnerID)
			got, _ := IsOwner(ctx, b)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
