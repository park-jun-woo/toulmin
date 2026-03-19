package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_IsEnabled(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	g.Warrant(IsBetaUser, nil, 1.0)
	flags.Register("dark-mode", g)

	ctx := &UserContext{Attributes: map[string]any{"beta": true}}
	enabled, err := flags.IsEnabled("dark-mode", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enabled {
		t.Error("expected enabled for beta user")
	}
}

func TestFlags_IsEnabled_Disabled(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	g.Warrant(IsBetaUser, nil, 1.0)
	flags.Register("dark-mode", g)

	ctx := &UserContext{Attributes: map[string]any{"beta": false}}
	enabled, err := flags.IsEnabled("dark-mode", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enabled {
		t.Error("expected disabled for non-beta user")
	}
}

func TestFlags_IsEnabled_LegacyDefeat(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	beta := g.Warrant(IsBetaUser, nil, 1.0)
	legacy := g.Rebuttal(IsLegacyBrowser, nil, 1.0)
	internal := g.Defeater(IsInternalStaff, nil, 1.0)
	g.Defeat(legacy, beta)
	g.Defeat(internal, legacy)
	flags.Register("dark-mode", g)

	// beta + legacy + internal → internal defeats legacy → enabled
	ctx := &UserContext{Attributes: map[string]any{"beta": true, "legacy_browser": true, "internal": true}}
	enabled, err := flags.IsEnabled("dark-mode", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enabled {
		t.Error("expected enabled (internal defeats legacy)")
	}
}

func TestFlags_IsEnabled_Unregistered(t *testing.T) {
	flags := NewFlags()
	_, err := flags.IsEnabled("nonexistent", &UserContext{})
	if err == nil {
		t.Error("expected error for unregistered feature")
	}
}

func TestFlags_EvaluateTrace(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	g.Warrant(IsBetaUser, nil, 1.0)
	flags.Register("dark-mode", g)

	ctx := &UserContext{Attributes: map[string]any{"beta": true}}
	result, err := flags.EvaluateTrace("dark-mode", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Enabled {
		t.Error("expected enabled")
	}
	if len(result.Trace) == 0 {
		t.Error("expected non-empty trace")
	}
}

func TestFlags_List(t *testing.T) {
	flags := NewFlags()

	g1 := toulmin.NewGraph("feature:a")
	g1.Warrant(IsBetaUser, nil, 1.0)
	flags.Register("a", g1)

	g2 := toulmin.NewGraph("feature:b")
	g2.Warrant(IsRegion, "KR", 1.0)
	flags.Register("b", g2)

	ctx := &UserContext{
		Region:     "KR",
		Attributes: map[string]any{"beta": true},
	}
	active, err := flags.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(active) != 2 {
		t.Errorf("expected 2 active features, got %d", len(active))
	}
}

func TestFlags_Percentage_Deterministic(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:rollout")
	g.Warrant(IsUserInPercentage, 0.5, 1.0)
	flags.Register("rollout", g)

	ctx := &UserContext{ID: "user-42"}
	r1, _ := flags.IsEnabled("rollout", ctx)
	r2, _ := flags.IsEnabled("rollout", ctx)
	if r1 != r2 {
		t.Error("expected deterministic result for same user")
	}
}
