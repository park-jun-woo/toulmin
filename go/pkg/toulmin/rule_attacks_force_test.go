//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleAttacksForce — Attacks stores the default force 1, an explicit force, survives With(), and panics on invalid force
package toulmin

import "testing"

func TestRuleAttacksForce(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	d := g.Except(DefeaterC)
	r.Attacks(w)
	d.Attacks(r, 0.25)
	if g.defeats[0].force != 1.0 {
		t.Fatalf("expected default force 1.0, got %v", g.defeats[0].force)
	}
	if g.defeats[1].force != 0.25 {
		t.Fatalf("expected force 0.25, got %v", g.defeats[1].force)
	}
	d.With(&testSpec{Value: "x"})
	if g.defeats[1].from != d.id || g.defeats[1].force != 0.25 {
		t.Fatalf("expected With() to rename the edge but keep its force, got %+v", g.defeats[1])
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for force > 1")
		}
		if len(g.defeats) != 2 {
			t.Fatalf("invalid force must not append an edge, got %d edges", len(g.defeats))
		}
	}()
	r.Attacks(w, 2.0)
}
