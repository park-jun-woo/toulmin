//ff:func feature=policy type=rule control=iteration dimension=1
//ff:what TestIsOwner — tests IsOwner rule
package policy

import "testing"

func TestIsOwner(t *testing.T) {
	ob := &OwnerBacking{
		UserIDFunc:     func(u any) string { return u.(*testUser).ID },
		ResourceIDFunc: func(ctx any) string { return ctx.(*RequestContext).ResourceOwnerID },
	}
	tests := []struct {
		name    string
		user    any
		ownerID string
		want    bool
	}{
		{"owner", &testUser{ID: "u1"}, "u1", true},
		{"not owner", &testUser{ID: "u1"}, "u2", false},
		{"nil user", nil, "u1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &RequestContext{User: tt.user, ResourceOwnerID: tt.ownerID}
			got, _ := IsOwner(nil, ctx, ob)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
