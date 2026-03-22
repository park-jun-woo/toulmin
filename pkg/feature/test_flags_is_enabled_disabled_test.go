//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_IsEnabled_Disabled — tests feature flag disabled for non-beta user
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

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
