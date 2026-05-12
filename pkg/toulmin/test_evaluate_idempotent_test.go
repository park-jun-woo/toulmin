//ff:func feature=engine type=engine control=sequence
//ff:what TestEvaluateIdempotent — tests that two Evaluate calls on same graph produce identical results
package toulmin

import (
	"testing"
)

func TestEvaluateIdempotent(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA)
	g.Counter(RebuttalB)

	r1, err1 := g.Evaluate(NewContext())
	if err1 != nil {
		t.Fatalf("first evaluate error: %v", err1)
	}
	r2, err2 := g.Evaluate(NewContext())
	if err2 != nil {
		t.Fatalf("second evaluate error: %v", err2)
	}
	if len(r1) != len(r2) {
		t.Fatalf("result count mismatch: %d vs %d", len(r1), len(r2))
	}
	for i := range r1 {
		if r1[i].Verdict != r2[i].Verdict {
			t.Errorf("verdict mismatch at %d: %f vs %f", i, r1[i].Verdict, r2[i].Verdict)
		}
		if r1[i].Name != r2[i].Name {
			t.Errorf("name mismatch at %d: %s vs %s", i, r1[i].Name, r2[i].Name)
		}
	}
}
