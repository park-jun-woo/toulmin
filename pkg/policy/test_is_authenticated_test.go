//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsAuthenticated — tests IsAuthenticated rule
package policy

import "testing"

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
			ctx := &RequestContext{User: tt.user}
			got, _ := IsAuthenticated(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
