//ff:func feature=engine type=engine control=sequence
//ff:what TestQualifierZero — tests that qualifier=0 produces verdict -1.0
package toulmin

import (
	"testing"
)

func TestQualifierZero(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA).Qualifier(0.0)
	results, err := g.Evaluate(NewContext())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != -1.0 {
		t.Errorf("expected -1.0 (qualifier=0), got %f", results[0].Verdict)
	}
}
