//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphWarrantOnly — tests graph verdict for single warrant
package toulmin

import (
	"testing"
)

func TestGraphWarrantOnly(t *testing.T) {
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
