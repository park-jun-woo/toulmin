//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphRunDepth — runDepth covers depth guard, evaluate/handler/sub errors, and recursion
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

// TestGraphRunDepth exercises every branch of (*Graph).runDepth by calling it directly:
// the depth-cap backstop, the evaluate-error early return, the handler-error path, the
// sub-Run error path, the Active+RunGraph success recursion, and the plain no-handler /
// no-RunGraph leaf path.
func TestGraphRunDepth(t *testing.T) {
	// Distinct closures: two rules in one graph must be different function values.
	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	active2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	// (1) depth guard: depth above the cap returns an error before any work.
	t.Run("depthGuard", func(t *testing.T) {
		g := NewGraph("deep")
		g.Rule(active)
		results, view, err := g.runDepth(NewContext(), EvalOption{}, runMaxDepth+1)
		if err == nil || !strings.Contains(err.Error(), "depth exceeded") {
			t.Fatalf("expected depth exceeded error, got results=%v view=%v err=%v", results, view, err)
		}
		if results != nil || view != nil {
			t.Errorf("depth guard must return nil results/view, got %v %v", results, view)
		}
	})

	// (2) evaluate error: a circular defeat graph fails the full pass.
	t.Run("evaluateError", func(t *testing.T) {
		g := NewGraph("cycle")
		a := g.Rule(active)
		b := g.Counter(active2)
		a.Attacks(b)
		b.Attacks(a)
		results, view, err := g.runDepth(NewContext(), EvalOption{}, 0)
		if err == nil {
			t.Fatal("expected evaluate error from circular defeat graph")
		}
		if results != nil || view != nil {
			t.Errorf("evaluate error must return nil results/view, got %v %v", results, view)
		}
	})

	// (3) handler error: a node's handler returns an error; Run stops and wraps it.
	t.Run("handlerError", func(t *testing.T) {
		g := NewGraph("handler")
		g.Rule(active).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
			return fmt.Errorf("boom")
		})
		results, view, err := g.runDepth(NewContext(), EvalOption{}, 0)
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("expected wrapped handler error, got %v", err)
		}
		if results == nil || view == nil {
			t.Error("handler error must still return pre-dispatch results and view")
		}
	})

	// (4) sub-Run error: an Active node Runs a sub-graph whose handler errors.
	t.Run("subRunError", func(t *testing.T) {
		sub := NewGraph("sub")
		sub.Rule(active).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
			return fmt.Errorf("sub boom")
		})
		parent := NewGraph("parent")
		parent.Rule(active).Run(sub)
		results, view, err := parent.runDepth(NewContext(), EvalOption{}, 0)
		if err == nil || !strings.Contains(err.Error(), "run ") || !strings.Contains(err.Error(), "→") {
			t.Fatalf("expected wrapped sub-Run error, got %v", err)
		}
		if results == nil || view == nil {
			t.Error("sub-Run error must still return pre-dispatch results and view")
		}
	})

	// (5) success recursion: an Active node Runs its sub-graph; (6) a plain leaf node with
	// no handler and no RunGraph hits the fall-through branches.
	t.Run("recurseAndLeaf", func(t *testing.T) {
		subRuns := 0
		sub := NewGraph("sub")
		sub.Rule(active).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
			subRuns++
			return nil
		})
		parent := NewGraph("parent")
		parent.Rule(active).Run(sub) // Active + RunGraph -> recursion
		parent.Rule(active2)         // no handler, no RunGraph -> leaf fall-through
		results, view, err := parent.runDepth(NewContext(), EvalOption{}, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if results == nil || view == nil {
			t.Error("success path must return results and view")
		}
		if subRuns != 1 {
			t.Errorf("sub-graph should Run once, got %d", subRuns)
		}
	})
}
