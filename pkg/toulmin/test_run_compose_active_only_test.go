//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeActiveOnly — a non-Active node does NOT Run its sub-graph
package toulmin

import "testing"

func TestRunComposeActiveOnly(t *testing.T) {
	subFired := false
	subRule := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	sub := NewGraph("sub")
	sub.Rule(subRule).OnActive(func(ctx Context, ev NodeEvent, view RunView) error {
		subFired = true
		return nil
	})

	// Inactive node: rule fn returns false → never Active → sub-graph must not Run.
	inactive := func(ctx Context, specs Specs) (bool, any) { return false, nil }
	g := NewGraph("inactive-parent")
	g.Rule(inactive).Run(sub)
	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if subFired {
		t.Error("Inactive node must NOT Run its sub-graph")
	}

	// Defeated node: applies but is attacked → Defeated, not Active → sub must not Run.
	subFired = false
	warrant := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	attacker := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g2 := NewGraph("defeated-parent")
	w := g2.Rule(warrant).Run(sub)
	g2.Except(attacker).Attacks(w)
	if _, _, err := g2.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if subFired {
		t.Error("Defeated node must NOT Run its sub-graph")
	}
}
