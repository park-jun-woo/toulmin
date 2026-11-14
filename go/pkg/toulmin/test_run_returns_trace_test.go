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
	// Trace: one entry per registered rule, in registration order.
	all := trace.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(all))
	}
	if all[0].Name != "WarrantA" || all[1].Name != "InactiveR" {
		t.Errorf("trace order mismatch: %q, %q", all[0].Name, all[1].Name)
	}
	// Get resolves an entry by short name; a miss returns false.
	warrant, ok := trace.Get("WarrantA")
	if !ok {
		t.Fatal("Get(\"WarrantA\") should hit")
	}
	inactive, ok := trace.Get("InactiveR")
	if !ok {
		t.Fatal("Get(\"InactiveR\") should hit")
	}
	if _, ok := trace.Get("Nope"); ok {
		t.Error("Get of an unknown name should miss")
	}
	// WarrantA active, InactiveR inactive (returns false).
	if !warrant.Activated {
		t.Error("WarrantA should be Activated")
	}
	if inactive.Activated {
		t.Error("InactiveR should be inactive")
	}
	// Verdict on the trace must match the evaluation result.
	if warrant.Verdict != results[0].Verdict {
		t.Errorf("trace verdict %f != result verdict %f", warrant.Verdict, results[0].Verdict)
	}
	// Inactive node has no cached verdict → 0.0.
	if inactive.Verdict != 0.0 {
		t.Errorf("inactive node verdict should be 0.0, got %f", inactive.Verdict)
	}
	// Ground is the ctx as-is (same reference for every entry).
	if warrant.Ground != Context(ctx) || inactive.Ground != Context(ctx) {
		t.Error("every entry's Ground must be the ctx passed to Run")
	}
	// Ctx exposes the same context passed to Run.
	if trace.Ctx() != Context(ctx) {
		t.Error("trace.Ctx() must be the ctx passed to Run")
	}
}
