//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_IsEnabled_LegacyDefeat — tests internal staff defeats legacy browser rebuttal
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
