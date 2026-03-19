package state

import (
	"testing"
	"time"
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
			req := &TransitionRequest{From: tt.from}
			ctx := &TransitionContext{CurrentState: tt.current}
			got, _ := IsCurrentState(req, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

type testUser struct{ ID string }
type testResource struct {
	OwnerID   string
	ExpiresAt time.Time
}

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

func TestIsExpired(t *testing.T) {
	expiryFunc := func(r any) time.Time { return r.(*testResource).ExpiresAt }
	tests := []struct {
		name string
		exp  time.Time
		want bool
	}{
		{"expired", time.Now().Add(-1 * time.Hour), true},
		{"not expired", time.Now().Add(1 * time.Hour), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &TransitionContext{Resource: &testResource{ExpiresAt: tt.exp}}
			got, _ := IsExpired(nil, ctx, expiryFunc)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
