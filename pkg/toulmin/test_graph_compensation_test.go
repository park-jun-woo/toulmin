//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphCompensation — tests graph verdict with compensation chain
package toulmin

import (
	"math"
	"testing"
)

func TestGraphCompensation(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	d := g.Defeater(DefeaterC, nil, 1.0)
	g.Defeat(r, w)
	g.Defeat(d, r)
	results, err := g.Evaluate(nil, nil, EvalOption{Trace: true})
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
