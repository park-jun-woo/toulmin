//ff:func feature=state type=rule control=iteration dimension=1
//ff:what TestIsOwner — tests IsOwner rule
package state

import "testing"

func TestIsOwner(t *testing.T) {
	b := &OwnerBacking{
		OwnerIDFunc: func(r any) string { return r.(*testResource).OwnerID },
		UserIDFunc:  func(u any) string { return u.(*testUser).ID },
	}
	tests := []struct {
		name string
		user *testUser
		res  *testResource
		want bool
	}{
		{"owner", &testUser{ID: "u1"}, &testResource{OwnerID: "u1"}, true},
		{"not owner", &testUser{ID: "u1"}, &testResource{OwnerID: "u2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &TransitionContext{User: tt.user, Resource: tt.res}
			got, _ := IsOwner(nil, ctx, b)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
