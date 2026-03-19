package feature

import "testing"

func TestIsBetaUser(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"beta", map[string]any{"beta": true}, true},
		{"not beta", map[string]any{"beta": false}, false},
		{"no attr", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsBetaUser(nil, &UserContext{Attributes: tt.attr}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInternalStaff(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"internal", map[string]any{"internal": true}, true},
		{"not internal", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsInternalStaff(nil, &UserContext{Attributes: tt.attr}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRegion(t *testing.T) {
	tests := []struct {
		name   string
		region string
		back   string
		want   bool
	}{
		{"match", "KR", "KR", true},
		{"mismatch", "US", "KR", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsRegion(nil, &UserContext{Region: tt.region}, tt.back)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		pair [2]any
		want bool
	}{
		{"match", map[string]any{"plan": "pro"}, [2]any{"plan", "pro"}, true},
		{"mismatch", map[string]any{"plan": "free"}, [2]any{"plan", "pro"}, false},
		{"missing", map[string]any{}, [2]any{"plan", "pro"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasAttribute(nil, &UserContext{Attributes: tt.attr}, tt.pair)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLegacyBrowser(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"legacy", map[string]any{"legacy_browser": true}, true},
		{"not legacy", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsLegacyBrowser(nil, &UserContext{Attributes: tt.attr}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUserInPercentage(t *testing.T) {
	ctx := &UserContext{ID: "user-123"}
	h := hashPercentage("user-123")

	got1, _ := IsUserInPercentage(nil, ctx, 1.0)
	if !got1 {
		t.Error("expected true for 100% rollout")
	}

	got2, _ := IsUserInPercentage(nil, ctx, 0.0)
	if got2 {
		t.Error("expected false for 0% rollout")
	}

	// deterministic
	got3, _ := IsUserInPercentage(nil, ctx, h+0.01)
	got4, _ := IsUserInPercentage(nil, ctx, h+0.01)
	if got3 != got4 {
		t.Error("expected deterministic result")
	}
}
