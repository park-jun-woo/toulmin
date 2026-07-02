//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_EvaluateTrace_Branches — covers unregistered, error, empty, and disabled branches of EvaluateTrace
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_EvaluateTrace_Branches(t *testing.T) {
	t.Run("NotRegistered", func(t *testing.T) {
		flags := NewFlags()

		result, err := flags.EvaluateTrace("missing", &UserContext{})
		if err == nil {
			t.Fatal("expected error for unregistered feature")
		}
		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})

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

		result, err := flags.EvaluateTrace("cycle", &UserContext{})
		if err == nil {
			t.Fatal("expected error from cyclic defeat graph")
		}
		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})

	t.Run("EmptyGraph", func(t *testing.T) {
		flags := NewFlags()

		g := toulmin.NewGraph("feature:empty")
		flags.Register("empty", g)

		result, err := flags.EvaluateTrace("empty", &UserContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Enabled {
			t.Error("expected disabled for empty graph")
		}
		if result.Verdict != -1 {
			t.Errorf("expected verdict -1, got %v", result.Verdict)
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		flags := NewFlags()

		g := toulmin.NewGraph("feature:dark-mode-off")
		g.Rule(IsBetaUser)
		flags.Register("dark-mode-off", g)

		ctx := &UserContext{Attributes: map[string]any{"beta": false}}
		result, err := flags.EvaluateTrace("dark-mode-off", ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Enabled {
			t.Error("expected disabled")
		}
	})
}
