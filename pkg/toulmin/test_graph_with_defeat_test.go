//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphWithDefeat — tests graph verdict with defeat edge
package toulmin

import (
	"testing"
)

func TestGraphWithDefeat(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	g.Defeat(r, w)
	results, err := g.EvaluateTrace(nil, nil)
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
