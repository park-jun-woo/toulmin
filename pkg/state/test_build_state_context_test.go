//ff:func feature=state type=engine control=iteration dimension=1
//ff:what TestBuildStateContext — verifies buildStateContext maps all TransitionRequest and TransitionContext fields into toulmin.Context
package state

import "testing"

func TestBuildStateContext(t *testing.T) {
	req := &TransitionRequest{
		From:  "draft",
		To:    "published",
		Event: "publish",
	}
	tc := &TransitionContext{
		CurrentState:    "draft",
		User:            "alice",
		Resource:        "doc-1",
		Metadata:        map[string]any{"k": "v"},
		UserID:          "u-1",
		ResourceOwnerID: "u-1",
	}

	ctx := buildStateContext(req, tc)

	cases := []struct {
		key  string
		want any
	}{
		{"from", req.From},
		{"to", req.To},
		{"event", req.Event},
		{"currentState", tc.CurrentState},
		{"user", tc.User},
		{"resource", tc.Resource},
		{"userID", tc.UserID},
		{"resourceOwnerID", tc.ResourceOwnerID},
	}
	for _, c := range cases {
		v, ok := ctx.Get(c.key)
		if !ok {
			t.Fatalf("key %q not set", c.key)
		}
		if v != c.want {
			t.Fatalf("key %q = %v, want %v", c.key, v, c.want)
		}
	}

	md, ok := ctx.Get("metadata")
	if !ok {
		t.Fatal("metadata not set")
	}
	if m, ok := md.(map[string]any); !ok || m["k"] != "v" {
		t.Fatalf("metadata not set correctly: %v", md)
	}
}
