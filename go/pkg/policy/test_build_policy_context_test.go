//ff:func feature=policy type=engine control=iteration dimension=1
//ff:what TestBuildPolicyContext — verifies buildPolicyContext maps all RequestContext fields into toulmin.Context
package policy

import "testing"

func TestBuildPolicyContext(t *testing.T) {
	rc := &RequestContext{
		User:            "alice",
		ClientIP:        "127.0.0.1",
		ResourceOwnerID: "owner-1",
		Headers:         map[string]string{"X-Test": "1"},
		RateLimiter:     nil,
		Metadata:        map[string]any{"k": "v"},
		Role:            "admin",
		UserID:          "u-1",
		ResourceOwner:   "u-1",
		IPBlocked:       true,
	}

	ctx := buildPolicyContext(rc)

	cases := []struct {
		key  string
		want any
	}{
		{"user", rc.User},
		{"clientIP", rc.ClientIP},
		{"resourceOwnerID", rc.ResourceOwnerID},
		{"headers", "X-Test"},
		{"rateLimiter", rc.RateLimiter},
		{"role", rc.Role},
		{"userID", rc.UserID},
		{"resourceOwner", rc.ResourceOwner},
		{"ipBlocked", rc.IPBlocked},
	}

	for _, c := range cases {
		v, ok := ctx.Get(c.key)
		if !ok {
			t.Fatalf("key %q not set", c.key)
		}
		if c.key == "headers" && !validHeaderValue(v) {
			t.Fatalf("headers not set correctly: %v", v)
		}
		if c.key == "headers" {
			continue
		}
		if v != c.want {
			t.Fatalf("key %q = %v, want %v", c.key, v, c.want)
		}
	}

	md, ok := ctx.Get("metadata")
	if !ok {
		t.Fatal("metadata not set")
	}
	m, ok := md.(map[string]any)
	if !ok || m["k"] != "v" {
		t.Fatalf("metadata not set correctly: %v", md)
	}
}
