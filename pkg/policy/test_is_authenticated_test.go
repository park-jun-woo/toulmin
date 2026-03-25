//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsAuthenticated — tests IsAuthenticated rule
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		user any
		want bool
	}{
		{"authenticated", &testUser{ID: "u1"}, true},
		{"not authenticated", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("user", tt.user)
			got, _ := IsAuthenticated(ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
