//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCycleError — Run propagates the evaluation error from a circular defeat graph
package toulmin

import "testing"

func TestRunCycleError(t *testing.T) {
	cycleA := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	cycleB := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("cycle")
	a := g.Rule(cycleA)
	b := g.Counter(cycleB)
	b.Attacks(a)
	a.Attacks(b)

	results, view, err := g.Run(NewContext())
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
	if results != nil || view != nil {
		t.Errorf("on evaluate error Run must return nil results and nil view, got results=%v view=%v", results, view)
	}
}
