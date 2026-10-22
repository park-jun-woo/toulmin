//ff:func feature=engine type=engine control=sequence
//ff:what TestEvaluateTraceSupersedes — Trace mode wires Excluded per warrant, shared or not-shared supersede source
package toulmin

import "testing"

func TestEvaluateTraceSupersedes(t *testing.T) {
	ctx := NewContext()

	// Shared source: X supersedes both warrants Y1 and Y2. In Trace mode each warrant is
	// evaluated on reset state; each warrant's own trace must carry Excluded==true.
	t.Run("shared supersede source", func(t *testing.T) {
		g := NewGraph("trace-shared")
		y1 := g.Rule(WarrantA)
		y2 := g.Rule(DefeaterC)
		x := g.Counter(RebuttalB)
		x.Supersedes(y1)
		x.Supersedes(y2)
		results, err := g.Evaluate(ctx, EvalOption{Trace: true})
		if err != nil {
			t.Fatalf("Evaluate: %v", err)
		}
		if ex, ok := excludedIn(results, "WarrantA", "WarrantA"); !ok || !ex {
			t.Fatalf("WarrantA warrant should be Excluded in its own trace (ok=%v ex=%v)", ok, ex)
		}
		if ex, ok := excludedIn(results, "DefeaterC", "DefeaterC"); !ok || !ex {
			t.Fatalf("DefeaterC warrant should be Excluded in its own trace (ok=%v ex=%v)", ok, ex)
		}
	})

	// Not shared: X supersedes only Y1. Y1's warrant trace is Excluded; Y2's is not.
	t.Run("unshared supersede source", func(t *testing.T) {
		g := NewGraph("trace-unshared")
		y1 := g.Rule(WarrantA)
		g.Rule(DefeaterC)
		x := g.Counter(RebuttalB)
		x.Supersedes(y1)
		results, err := g.Evaluate(ctx, EvalOption{Trace: true})
		if err != nil {
			t.Fatalf("Evaluate: %v", err)
		}
		if ex, ok := excludedIn(results, "WarrantA", "WarrantA"); !ok || !ex {
			t.Fatalf("WarrantA warrant should be Excluded (ok=%v ex=%v)", ok, ex)
		}
		if ex, ok := excludedIn(results, "DefeaterC", "DefeaterC"); !ok || ex {
			t.Fatalf("DefeaterC warrant must NOT be Excluded (ok=%v ex=%v)", ok, ex)
		}
	})
}
