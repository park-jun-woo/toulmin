//ff:func feature=engine type=engine control=sequence
//ff:what TestRunVerdictZero — verdict == 0 classifies the node as Defeated
package toulmin

import "testing"

func TestRunVerdictZero(t *testing.T) {
	var gotType NodeEventType
	var gotVerdict float64
	rec := func(ctx Context, ev NodeEvent, view RunView) error {
		gotType = ev.Type
		gotVerdict = ev.Verdict
		return nil
	}
	g := NewGraph("zero")
	w := g.Rule(WarrantA).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	r := g.Counter(RebuttalB)
	r.Attacks(w)

	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if gotVerdict != 0.0 {
		t.Fatalf("expected verdict 0.0, got %f", gotVerdict)
	}
	if gotType != Defeated {
		t.Errorf("verdict 0 want Defeated, got %v", gotType)
	}
}
