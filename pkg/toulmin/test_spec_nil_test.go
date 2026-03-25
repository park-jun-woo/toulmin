//ff:func feature=engine type=engine control=sequence
//ff:what TestSpecNil — tests that nil specs appears as nil in trace
package toulmin

import (
	"testing"
)

func TestSpecNil(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA)
	results, err := g.Evaluate(nil, EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Trace[0].Specs) != 0 {
		t.Errorf("expected specs nil, got %v", results[0].Trace[0].Specs)
	}
}
