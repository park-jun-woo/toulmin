//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_IsEnabled — tests feature flag enabled for beta user
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_IsEnabled(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	g.Rule(IsBetaUser)
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
