//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsInRole — tests IsInRole rule
package policy

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsInRole(t *testing.T) {
	tests := []struct {
		name    string
		user    any
		role    string
		ctxRole string
		want    bool
	}{
		{"match", &testUser{Role: "admin"}, "admin", "admin", true},
		{"mismatch", &testUser{Role: "user"}, "admin", "user", false},
		{"nil user", nil, "admin", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("user", tt.user)
			ctx.Set("role", tt.ctxRole)
			rb := &RoleBacking{Role: tt.role}
			got, _ := IsInRole(ctx, rb)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
