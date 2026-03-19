package policy

import "testing"

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name string
		ctx  *RequestContext
		want bool
	}{
		{"authenticated", &RequestContext{User: &User{ID: "u1"}}, true},
		{"not authenticated", &RequestContext{User: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsAuthenticated(nil, tt.ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		ctx  *RequestContext
		want bool
	}{
		{"match", "admin", &RequestContext{User: &User{Role: "admin"}}, true},
		{"mismatch", "admin", &RequestContext{User: &User{Role: "user"}}, false},
		{"nil user", "admin", &RequestContext{User: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsInRole(nil, tt.ctx, tt.role)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOwner(t *testing.T) {
	tests := []struct {
		name string
		ctx  *RequestContext
		want bool
	}{
		{"owner", &RequestContext{User: &User{ID: "u1"}, ResourceOwnerID: "u1"}, true},
		{"not owner", &RequestContext{User: &User{ID: "u1"}, ResourceOwnerID: "u2"}, false},
		{"nil user", &RequestContext{User: nil, ResourceOwnerID: "u1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsOwner(nil, tt.ctx, nil)
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
		{"custom header", map[string]string{"X-Api-Key": "abc"}, "X-Api-Key", true},
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
