//ff:func feature=engine type=engine control=sequence
//ff:what TestWithThenAttacks — tests defeat edge after With then Attacks
package toulmin

import "testing"

func TestWithThenAttacks(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.With(&testSpec{Value: "y"})
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
