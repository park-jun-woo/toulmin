//ff:func feature=engine type=engine control=sequence
//ff:what TestRunVerdictZero — verdict == 0 is Defeated, so RunOn never fires
package toulmin

import "testing"

func TestRunVerdictZero(t *testing.T) {
	fired := false
	rec := func(t Trace) error {
		fired = true
		return nil
	}
	g := NewGraph("zero")
	w := g.Rule(WarrantA).RunOn(rec)
	r := g.Counter(RebuttalB)
	r.Attacks(w)

	_, trace, err := g.Run(NewContext())
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if fired {
		t.Error("verdict 0 (Defeated) must not fire RunOn")
	}
	// trace still records WarrantA with verdict 0.0 (balanced).
	got, ok := trace.Get("WarrantA")
	if !ok {
		t.Fatal("WarrantA missing from trace")
	}
	if got.Verdict != 0.0 {
		t.Errorf("expected verdict 0.0, got %f", got.Verdict)
	}
}
