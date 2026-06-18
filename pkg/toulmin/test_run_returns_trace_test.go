//ff:func feature=engine type=engine control=sequence
//ff:what TestRunReturnsTrace — Run's second return is a flat trace of every node with Verdict/Ground/Backing
package toulmin

import "testing"

func TestRunReturnsTrace(t *testing.T) {
	g := NewGraph("trace")
	w := g.Rule(WarrantA)
	r := g.Counter(InactiveR)
	r.Attacks(w)

	ctx := NewContext()
	ctx.Set("k", "v")
	results, trace, err := g.Run(ctx)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	// Flat trace: one entry per registered rule, in registration order.
	if len(trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(trace))
	}
	if trace[0].Name != "WarrantA" || trace[1].Name != "InactiveR" {
		t.Errorf("trace order mismatch: %q, %q", trace[0].Name, trace[1].Name)
	}
	// WarrantA active, InactiveR inactive (returns false).
	if !trace[0].Activated {
		t.Error("WarrantA should be Activated")
	}
	if trace[1].Activated {
		t.Error("InactiveR should be inactive")
	}
	// Verdict on the trace must match the evaluation result.
	if trace[0].Verdict != results[0].Verdict {
		t.Errorf("trace verdict %f != result verdict %f", trace[0].Verdict, results[0].Verdict)
	}
	// Inactive node has no cached verdict → 0.0.
	if trace[1].Verdict != 0.0 {
		t.Errorf("inactive node verdict should be 0.0, got %f", trace[1].Verdict)
	}
	// Ground is the ctx as-is (same reference for every entry).
	if trace[0].Ground != Context(ctx) || trace[1].Ground != Context(ctx) {
		t.Error("every entry's Ground must be the ctx passed to Run")
	}
}
