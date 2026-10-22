//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphCompensation — tests graph verdict with compensation chain
package toulmin

import (
	"math"
	"testing"
)

func TestGraphCompensation(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	d := g.Except(DefeaterC)
	r.Attacks(w)
	d.Attacks(r)
	results, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}
