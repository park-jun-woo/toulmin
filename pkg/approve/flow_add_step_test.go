//ff:func feature=approve type=engine control=sequence
//ff:what TestFlowAddStep — Flow.AddStep appends a Step and returns the receiver for chaining
package approve

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestFlowAddStep covers the single branch of Flow.AddStep: it always
// appends a new *Step wrapping the given name and graph, and returns the
// same *Flow so calls can be chained.
func TestFlowAddStep(t *testing.T) {
	f := NewFlow("expense")
	g1 := toulmin.NewGraph("step1")
	g2 := toulmin.NewGraph("step2")

	got := f.AddStep("manager", g1)
	if got != f {
		t.Errorf("AddStep must return the receiver for chaining, got %v want %v", got, f)
	}
	if len(f.steps) != 1 {
		t.Fatalf("len(steps) = %d, want 1", len(f.steps))
	}
	if f.steps[0].Name != "manager" || f.steps[0].Graph != g1 {
		t.Errorf("steps[0] = %+v, want Name=manager Graph=g1", f.steps[0])
	}

	// Chained call: a second AddStep appends onto the same slice.
	f.AddStep("finance", g2).AddStep("ceo", g1)
	if len(f.steps) != 3 {
		t.Fatalf("len(steps) = %d, want 3", len(f.steps))
	}
	if f.steps[1].Name != "finance" || f.steps[1].Graph != g2 {
		t.Errorf("steps[1] = %+v, want Name=finance Graph=g2", f.steps[1])
	}
	if f.steps[2].Name != "ceo" || f.steps[2].Graph != g1 {
		t.Errorf("steps[2] = %+v, want Name=ceo Graph=g1", f.steps[2])
	}
}
