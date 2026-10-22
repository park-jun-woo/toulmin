//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunEvaluatePure — Evaluate stays idempotent and fires no handlers
package toulmin

import "testing"

func TestRunEvaluatePure(t *testing.T) {
	fired := false
	mark := func(self TraceEntry, t Trace) error {
		fired = true
		return nil
	}
	g := NewGraph("test")
	g.Rule(WarrantA).RunOn(mark)
	g.Counter(RebuttalB).RunOn(mark)

	r1, err1 := g.Evaluate(NewContext())
	if err1 != nil {
		t.Fatalf("first evaluate error: %v", err1)
	}
	r2, err2 := g.Evaluate(NewContext())
	if err2 != nil {
		t.Fatalf("second evaluate error: %v", err2)
	}
	if fired {
		t.Fatal("Evaluate must not fire node handlers")
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
