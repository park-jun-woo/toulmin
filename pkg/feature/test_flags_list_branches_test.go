//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_List_Branches — covers error, disabled, and empty-order branches of List
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_List_Branches(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		flags := NewFlags()

		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("feature:cycle")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)
		flags.Register("cycle", g)

		active, err := flags.List(&UserContext{})
		if err == nil {
			t.Fatal("expected error from cyclic defeat graph")
		}
		if active != nil {
			t.Errorf("expected nil active list, got %v", active)
		}
	})

	t.Run("SomeDisabled", func(t *testing.T) {
		flags := NewFlags()

		g1 := toulmin.NewGraph("feature:a")
		g1.Rule(IsBetaUser)
		flags.Register("a", g1)

		g2 := toulmin.NewGraph("feature:b")
		g2.Rule(IsRegion).With(&RegionSpec{Region: "KR"})
		flags.Register("b", g2)

		ctx := &UserContext{
			Region:     "US",
			Attributes: map[string]any{"beta": false},
		}
		active, err := flags.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(active) != 0 {
			t.Errorf("expected 0 active features, got %d: %v", len(active), active)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		flags := NewFlags()

		active, err := flags.List(&UserContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(active) != 0 {
			t.Errorf("expected 0 active features, got %d: %v", len(active), active)
		}
	})
}
