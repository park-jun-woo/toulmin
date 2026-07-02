//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_IsEnabled_Branches — covers graph error and empty-results branches of IsEnabled
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_IsEnabled_Branches(t *testing.T) {
	t.Run("GraphError", func(t *testing.T) {
		flags := NewFlags()

		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("feature:cycle")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)
		flags.Register("cycle", g)

		enabled, err := flags.IsEnabled("cycle", &UserContext{})
		if err == nil {
			t.Fatal("expected error from cyclic defeat graph")
		}
		if enabled {
			t.Error("expected disabled result on error")
		}
	})

	t.Run("EmptyGraph", func(t *testing.T) {
		flags := NewFlags()

		g := toulmin.NewGraph("feature:empty")
		flags.Register("empty", g)

		enabled, err := flags.IsEnabled("empty", &UserContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if enabled {
			t.Error("expected disabled for empty graph")
		}
	})
}
