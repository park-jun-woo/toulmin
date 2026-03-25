//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphTraceAllRules — tests that trace contains all executed rules
package toulmin

import (
	"testing"
)

func TestGraphTraceAllRules(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB).Qualifier(0.8)
	r.Attacks(w)
	results, err := g.Evaluate(nil, EvalOption{Trace: true})
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
	if trace[0].Name != "WarrantA" || trace[0].Role != "rule" || !trace[0].Activated || trace[0].Qualifier != 1.0 {
		t.Errorf("trace[0] unexpected: %+v", trace[0])
	}
	if trace[1].Name != "RebuttalB" || trace[1].Role != "counter" || !trace[1].Activated || trace[1].Qualifier != 0.8 {
		t.Errorf("trace[1] unexpected: %+v", trace[1])
	}
}
