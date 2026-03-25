//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestDuration — tests that Duration option measures per-rule execution time
package toulmin

import (
	"testing"
)

func TestDuration(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 0.8)
	g.Defeat(r, w)
	results, err := g.Evaluate(nil, nil, EvalOption{Duration: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(results[0].Trace))
	}
	for _, te := range results[0].Trace {
		if te.Duration <= 0 {
			t.Errorf("expected Duration > 0 for %s, got %v", te.Name, te.Duration)
		}
	}
}
