//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphFuncReuse — tests same func used in different graphs
package toulmin

import (
	"testing"
)

func TestGraphFuncReuse(t *testing.T) {
	g1 := NewGraph("graph1")
	w1 := g1.Rule(WarrantA)
	r1 := g1.Counter(RebuttalB)
	r1.Attacks(w1)

	g2 := NewGraph("graph2")
	w2 := g2.Rule(WarrantA)
	d2 := g2.Except(DefeaterC)
	d2.Attacks(w2)

	res1, err := g1.Evaluate(NewContext())
	if err != nil {
		t.Fatalf("g1 error: %v", err)
	}
	res2, err := g2.Evaluate(NewContext())
	if err != nil {
		t.Fatalf("g2 error: %v", err)
	}
	if len(res1) != 1 || len(res2) != 1 {
		t.Fatalf("expected 1 result each, got %d and %d", len(res1), len(res2))
	}
	if res1[0].Verdict != 0.0 {
		t.Errorf("g1: expected 0.0, got %f", res1[0].Verdict)
	}
	if res2[0].Verdict != 0.0 {
		t.Errorf("g2: expected 0.0, got %f", res2[0].Verdict)
	}
}
