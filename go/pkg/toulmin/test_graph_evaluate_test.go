//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphEvaluate — evaluate covers lazy/full passes, trace, and cycle/panic errors
package toulmin

import "testing"

func TestGraphEvaluate(t *testing.T) {
	ctx := NewContext()

	// Lazy pass: WarrantA active (result), InactiveR inactive (skipped),
	// counter attacks the inactive warrant so it is never reached lazily.
	g := NewGraph("eval")
	g.Rule(WarrantA)
	w := g.Rule(InactiveR)
	c := g.Counter(blockIP)
	c.Attacks(w)

	lazy, ecLazy, err := g.evaluate(ctx, EvalOption{}, false)
	if err != nil {
		t.Fatalf("lazy evaluate: %v", err)
	}
	if len(lazy) != 1 {
		t.Fatalf("lazy pass want 1 active result, got %d", len(lazy))
	}
	if ecLazy.ran[g.rules[2].Name] {
		t.Errorf("lazy pass should not run unreached node %q", g.rules[2].Name)
	}

	// Full pass: the unreached counter is calc'd to fill its state.
	_, ecFull, err := g.evaluate(ctx, EvalOption{}, true)
	if err != nil {
		t.Fatalf("full evaluate: %v", err)
	}
	if !ecFull.ran[g.rules[2].Name] {
		t.Errorf("full pass should run unreached node %q", g.rules[2].Name)
	}

	// Trace pass: exercises the reset() and result.Trace branches.
	traced, _, err := g.evaluate(ctx, EvalOption{Trace: true}, false)
	if err != nil {
		t.Fatalf("trace evaluate: %v", err)
	}
	if len(traced) != 1 || traced[0].Trace == nil {
		t.Errorf("trace pass should attach a trace, got %+v", traced)
	}

	// Cycle: newEvalContext returns an error.
	cycleA := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	cycleB := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	gc := NewGraph("cycle")
	ca := gc.Rule(cycleA)
	cb := gc.Counter(cycleB)
	cb.Attacks(ca)
	ca.Attacks(cb)
	if _, _, err := gc.evaluate(ctx, EvalOption{}, false); err == nil {
		t.Error("cyclic graph should return an error")
	}

	// Lazy panic: a panicking warrant sets ec.err during the lazy pass.
	panicRule := func(ctx Context, specs Specs) (bool, any) { panic("lazy boom") }
	gp := NewGraph("lazypanic")
	gp.Rule(panicRule)
	if _, _, err := gp.evaluate(ctx, EvalOption{}, false); err == nil {
		t.Error("panicking warrant should return an error on the lazy pass")
	}

	// Full panic: a node panicking only during the full pass sets ec.err there.
	panicCounter := func(ctx Context, specs Specs) (bool, any) { panic("full boom") }
	gf := NewGraph("fullpanic")
	wf := gf.Rule(InactiveR)
	cf := gf.Counter(panicCounter)
	cf.Attacks(wf)
	if _, _, err := gf.evaluate(ctx, EvalOption{}, true); err == nil {
		t.Error("panicking node should return an error on the full pass")
	}
}
