//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsOwner — tests IsOwner rule
package policy

import "testing"

func TestIsOwner(t *testing.T) {
	ob := &OwnerBacking{}
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
			ctx := &RequestContext{User: tt.user, UserID: tt.userID, ResourceOwner: tt.resourceOwner}
			got, _ := IsOwner(nil, ctx, ob)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
