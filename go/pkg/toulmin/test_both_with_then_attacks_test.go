//ff:func feature=engine type=engine control=sequence
//ff:what TestBothWithThenAttacks — tests defeat edge when both rules use With then Attacks
package toulmin

import "testing"

func TestBothWithThenAttacks(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	w.With(&testSpec{Value: "a"})
	r := g.Counter(RebuttalB)
	r.With(&testSpec{Value: "b"})
	r.Attacks(w)
	results, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}
