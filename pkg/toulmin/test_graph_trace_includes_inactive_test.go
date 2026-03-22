//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphTraceIncludesInactive — tests that inactive rules appear in trace
package toulmin

import (
	"testing"
)

func TestGraphTraceIncludesInactive(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(InactiveR, nil, 1.0)
	g.Defeat(r, w)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	trace := results[0].Trace
	if len(trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(trace))
	}
	if trace[1].Activated {
		t.Errorf("expected InactiveR activated=false, got true")
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (rebuttal inactive), got %f", results[0].Verdict)
	}
}
