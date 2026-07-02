//ff:func feature=moderate type=engine control=sequence
//ff:what TestReview_Branches — covers error, empty-results, and flag branches of Review
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestReview_Branches(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("test:cycle")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)

		mod := NewModerator(g)
		content := &Content{Body: "hello"}
		ctx := &ContentContext{Author: &Author{}, Channel: &Channel{}}

		result, err := mod.Review(content, ctx)
		if err == nil {
			t.Fatal("expected error from cyclic defeat graph")
		}
		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})

	t.Run("EmptyGraph", func(t *testing.T) {
		g := toulmin.NewGraph("test:empty")

		mod := NewModerator(g)
		content := &Content{Body: "hello"}
		ctx := &ContentContext{Author: &Author{}, Channel: &Channel{}}

		result, err := mod.Review(content, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Action != ActionBlock {
			t.Errorf("expected block, got %s", result.Action)
		}
		if result.Allowed {
			t.Error("expected not allowed")
		}
		if result.Verdict != -1 {
			t.Errorf("expected verdict -1, got %v", result.Verdict)
		}
	})

	t.Run("Flag", func(t *testing.T) {
		g := toulmin.NewGraph("test:flag")
		g.Rule(IsVerifiedUser).Qualifier(0.6)

		mod := NewModerator(g)
		content := &Content{Body: "hello"}
		ctx := &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}

		result, err := mod.Review(content, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Action != ActionFlag {
			t.Errorf("expected flag, got %s", result.Action)
		}
		if result.Allowed {
			t.Error("expected not allowed")
		}
	})
}
