//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphFuncReuse — tests same func used in different graphs
package toulmin

import (
	"testing"
)

func TestGraphFuncReuse(t *testing.T) {
	g1 := NewGraph("graph1")
	w1 := g1.Warrant(WarrantA, nil, 1.0)
	r1 := g1.Rebuttal(RebuttalB, nil, 1.0)
	g1.Defeat(r1, w1)

	g2 := NewGraph("graph2")
	w2 := g2.Warrant(WarrantA, nil, 1.0)
	d2 := g2.Defeater(DefeaterC, nil, 1.0)
	g2.Defeat(d2, w2)

	res1, err := g1.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("g1 error: %v", err)
	}
	res2, err := g2.Evaluate(nil, nil)
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
