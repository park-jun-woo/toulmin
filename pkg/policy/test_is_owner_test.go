//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsOwner — tests IsOwner rule
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsOwner(t *testing.T) {
	ob := &OwnerSpec{}
	tests := []struct {
		name          string
		user          any
		userID        string
		resourceOwner string
		want          bool
	}{
		{"owner", &testUser{ID: "u1"}, "u1", "u1", true},
		{"not owner", &testUser{ID: "u1"}, "u1", "u2", false},
		{"nil user", nil, "", "u1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("user", tt.user)
			ctx.Set("userID", tt.userID)
			ctx.Set("resourceOwner", tt.resourceOwner)
			got, _ := IsOwner(ctx, toulmin.Specs{ob})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
