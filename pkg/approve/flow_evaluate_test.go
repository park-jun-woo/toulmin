//ff:func feature=approve type=engine control=sequence
//ff:what TestFlowEvaluate — Flow.Evaluate propagates a step's graph error and defaults to a rejecting verdict when a step produces no results
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestFlowEvaluate covers the two branches of Flow.Evaluate not exercised
// by the existing scenario tests (all-passed, step-rejected, multi-step
// second-rejects, CEO override, frozen-with-urgent-defeat): a step whose
// graph evaluation errors (here, a circular defeat graph) propagates that
// error immediately, and a step whose graph produces zero results defaults
// its verdict to -1, rejecting the flow.
func TestFlowEvaluate(t *testing.T) {
	t.Run("step graph evaluation error is propagated", func(t *testing.T) {
		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("expense:cyclic")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)

		f := NewFlow("expense").AddStep("manager", g)
		req := &ApprovalRequest{Amount: 100}
		result, err := f.Evaluate(req, func(step string) *ApprovalContext {
			return &ApprovalContext{}
		})
		if err == nil {
			t.Fatal("expected a cycle error, got nil")
		}
		if result != nil {
			t.Errorf("result = %v, want nil", result)
		}
	})

	t.Run("no results defaults verdict to -1 and rejects", func(t *testing.T) {
		g := toulmin.NewGraph("expense:empty")
		f := NewFlow("expense").AddStep("manager", g)
		req := &ApprovalRequest{Amount: 100}
		result, err := f.Evaluate(req, func(step string) *ApprovalContext {
			return &ApprovalContext{}
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Approved {
			t.Error("expected rejected when the step graph produces no results")
		}
		if len(result.Steps) != 1 {
			t.Fatalf("len(Steps) = %d, want 1", len(result.Steps))
		}
		if result.Steps[0].Verdict != -1 {
			t.Errorf("Steps[0].Verdict = %v, want -1", result.Steps[0].Verdict)
		}
	})
}
