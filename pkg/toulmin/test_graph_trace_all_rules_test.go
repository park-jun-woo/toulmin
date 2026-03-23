//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphTraceAllRules — tests that trace contains all executed rules
package toulmin

import (
	"testing"
)

func TestGraphTraceAllRules(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 0.8)
	g.Defeat(r, w)
	results, err := g.Evaluate(nil, nil, EvalOption{Trace: true})
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
	if trace[0].Name != "WarrantA" || trace[0].Role != "warrant" || !trace[0].Activated || trace[0].Qualifier != 1.0 {
		t.Errorf("trace[0] unexpected: %+v", trace[0])
	}
	if trace[1].Name != "RebuttalB" || trace[1].Role != "rebuttal" || !trace[1].Activated || trace[1].Qualifier != 0.8 {
		t.Errorf("trace[1] unexpected: %+v", trace[1])
	}
}
