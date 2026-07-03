//ff:func feature=engine type=engine control=sequence
//ff:what TestTracedExclusions — superseded trace entries are marked Excluded; empty trace stays empty
package toulmin

import "testing"

func TestTracedExclusions(t *testing.T) {
	ctx := NewContext()

	// Y (WarrantA, active) is superseded by active X (RebuttalB) and attacked by A (DefeaterC).
	// Evaluating Y in Trace mode records both Y and A; tracedExclusions must mark Y (superseded)
	// Excluded and leave A (not superseded) untouched — covering both if branches inside the loop.
	t.Run("marks superseded entries excluded", func(t *testing.T) {
		g := NewGraph("traced")
		y := g.Rule(WarrantA)
		a := g.Counter(DefeaterC)
		a.Attacks(y)
		x := g.Counter(RebuttalB)
		x.Supersedes(y)

		ec, err := newEvalContext(g.rules, g.defeats, g.supersedes, g.roles)
		if err != nil {
			t.Fatalf("newEvalContext: %v", err)
		}
		ec.reset()
		ec.evalRule(y.id, ctx, EvalOption{Trace: true})
		if ec.err != nil {
			t.Fatalf("evalRule: %v", ec.err)
		}

		out := ec.tracedExclusions(ctx)
		if ex, ok := traceExcluded(out, "WarrantA"); !ok || !ex {
			t.Fatalf("WarrantA should be Excluded (ok=%v ex=%v)", ok, ex)
		}
		if ex, ok := traceExcluded(out, "DefeaterC"); !ok || ex {
			t.Fatalf("DefeaterC must NOT be Excluded (ok=%v ex=%v)", ok, ex)
		}
	})

	// Empty trace: the loop body never runs and tracedExclusions returns an empty result.
	t.Run("empty trace returns empty result", func(t *testing.T) {
		g := NewGraph("traced-empty")
		g.Rule(WarrantA)
		ec, err := newEvalContext(g.rules, g.defeats, g.supersedes, g.roles)
		if err != nil {
			t.Fatalf("newEvalContext: %v", err)
		}
		if out := ec.tracedExclusions(ctx); len(out) != 0 {
			t.Fatalf("expected empty trace, got %v", out)
		}
	})
}
