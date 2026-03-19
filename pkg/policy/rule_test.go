package policy

import "testing"

type testUser struct {
	ID    string
	Role  string
	Email string
}

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

func TestIsInRole(t *testing.T) {
	roleFunc := func(u any) string { return u.(*testUser).Role }
	tests := []struct {
		name string
		user any
		role string
		want bool
	}{
		{"match", &testUser{Role: "admin"}, "admin", true},
		{"mismatch", &testUser{Role: "user"}, "admin", false},
		{"nil user", nil, "admin", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &RequestContext{User: tt.user}
			rb := &RoleBacking{Role: tt.role, RoleFunc: roleFunc}
			got, _ := IsInRole(nil, ctx, rb)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestIsIPInList(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}
	whitelist := &IPListBacking{Purpose: "whitelist", Check: func(ip string) bool { return ip == "10.0.0.1" }}
	tests := []struct {
		name string
		ip   string
		list *IPListBacking
		want bool
	}{
		{"blocked", "1.2.3.4", blocklist, true},
		{"not blocked", "5.6.7.8", blocklist, false},
		{"whitelisted", "10.0.0.1", whitelist, true},
		{"not whitelisted", "5.6.7.8", whitelist, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsIPInList(nil, &RequestContext{ClientIP: tt.ip}, tt.list)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

type mockLimiter struct {
	limited map[string]bool
}

func (m *mockLimiter) IsLimited(key string) bool { return m.limited[key] }

func TestIsRateLimited(t *testing.T) {
	limiter := &mockLimiter{limited: map[string]bool{"1.2.3.4": true}}
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"limited", "1.2.3.4", true},
		{"not limited", "5.6.7.8", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsRateLimited(nil, &RequestContext{ClientIP: tt.ip, RateLimiter: limiter}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		header  string
		want    bool
	}{
		{"has token", map[string]string{"X-Internal-Token": "secret"}, "X-Internal-Token", true},
		{"no token", map[string]string{}, "X-Internal-Token", false},
		{"empty token", map[string]string{"X-Internal-Token": ""}, "X-Internal-Token", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasHeader(nil, &RequestContext{Headers: tt.headers}, tt.header)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
