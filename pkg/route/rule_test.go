package route

import "testing"

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		ctx  *RouteContext
		want bool
	}{
		{"authenticated", &RouteContext{User: &User{ID: "u1"}}, true},
		{"not authenticated", &RouteContext{User: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsAuthenticated(nil, tt.ctx, nil)
			if got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		ctx  *RouteContext
		want bool
	}{
		{"admin match", "admin", &RouteContext{User: &User{Role: "admin"}}, true},
		{"role mismatch", "admin", &RouteContext{User: &User{Role: "user"}}, false},
		{"nil user", "admin", &RouteContext{User: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasRole(nil, tt.ctx, tt.role)
			if got != tt.want {
				t.Errorf("HasRole(backing=%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestIsOwner(t *testing.T) {
	ownerFunc := func(ctx *RouteContext) string {
		id, _ := ctx.Metadata["owner_id"].(string)
		return id
	}
	tests := []struct {
		name string
		ctx  *RouteContext
		want bool
	}{
		{"is owner", &RouteContext{User: &User{ID: "u1"}, Metadata: map[string]any{"owner_id": "u1"}}, true},
		{"not owner", &RouteContext{User: &User{ID: "u1"}, Metadata: map[string]any{"owner_id": "u2"}}, false},
		{"nil user", &RouteContext{User: nil, Metadata: map[string]any{"owner_id": "u1"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsOwner(nil, tt.ctx, ownerFunc)
			if got != tt.want {
				t.Errorf("IsOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPBlocked(t *testing.T) {
	blocklist := func(ip string) bool { return ip == "1.2.3.4" }
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"blocked", "1.2.3.4", true},
		{"not blocked", "5.6.7.8", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsIPBlocked(nil, &RouteContext{ClientIP: tt.ip}, blocklist)
			if got != tt.want {
				t.Errorf("IsIPBlocked(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsWhitelisted(t *testing.T) {
	whitelist := func(ip string) bool { return ip == "10.0.0.1" }
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"whitelisted", "10.0.0.1", true},
		{"not whitelisted", "5.6.7.8", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsWhitelisted(nil, &RouteContext{ClientIP: tt.ip}, whitelist)
			if got != tt.want {
				t.Errorf("IsWhitelisted(%q) = %v, want %v", tt.ip, got, tt.want)
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
			got, _ := IsRateLimited(nil, &RouteContext{ClientIP: tt.ip}, limiter)
			if got != tt.want {
				t.Errorf("IsRateLimited(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsInternalService(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    bool
	}{
		{"has token", map[string]string{"X-Internal-Token": "secret"}, true},
		{"no token", map[string]string{}, false},
		{"empty token", map[string]string{"X-Internal-Token": ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsInternalService(nil, &RouteContext{Headers: tt.headers}, nil)
			if got != tt.want {
				t.Errorf("IsInternalService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAdminOverride(t *testing.T) {
	tests := []struct {
		name string
		ctx  *RouteContext
		want bool
	}{
		{"admin", &RouteContext{User: &User{Role: "admin"}}, true},
		{"not admin", &RouteContext{User: &User{Role: "user"}}, false},
		{"nil user", &RouteContext{User: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsAdminOverride(nil, tt.ctx, nil)
			if got != tt.want {
				t.Errorf("IsAdminOverride() = %v, want %v", got, tt.want)
			}
		})
	}
}
