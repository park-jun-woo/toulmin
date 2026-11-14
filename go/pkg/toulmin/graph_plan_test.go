//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphPlan — plan: built on first use and kept, rebuilt after the graph changes, and carrying the cycle error of a cyclic graph
package toulmin

import "testing"

func TestGraphPlan(t *testing.T) {
	ok := func(ctx Context, specs Specs) (float64, any) { return 1.0, nil }
	off := func(ctx Context, specs Specs) (float64, any) { return -1.0, nil }
	g := NewGraph("plan")
	w := g.Rule(ok)
	p := g.plan()
	if p == nil || p.err != nil || p.static == nil || len(p.static.fnMap) != 1 {
		t.Fatalf("expected a plan of one rule, got %+v", p)
	}
	if g.plan() != p {
		t.Fatal("an unchanged graph keeps its plan")
	}
	c := g.Counter(off)
	q := g.plan()
	if q == p || len(q.static.fnMap) != 2 {
		t.Fatal("a changed graph gets a new plan")
	}
	c.Attacks(w)
	w.Attacks(c)
	if r := g.plan(); r == q || r.err == nil || r.static != nil {
		t.Fatalf("a cyclic graph's plan carries the error, got %+v", r)
	}
}
