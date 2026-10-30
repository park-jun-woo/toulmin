//ff:func feature=engine type=engine control=sequence
//ff:what TestAnySupersederApplicable — empty/inactive/fully-gated report false, a_eff > 0 reports true
package toulmin

import "testing"

func TestAnySupersederApplicable(t *testing.T) {
	ctx := NewContext()
	g := NewGraph("any-sup")
	w := g.Rule(WarrantA)     // active (returns true), no attackers → a_eff 1
	r := g.Counter(InactiveR) // inactive (returns false) → a_eff 0
	x := g.Rule(RebuttalB)    // active, but fully gated by an exception → a_eff 0
	g.Except(DefeaterC).Attacks(x)
	_, ec, err := g.evaluate(ctx, EvalOption{}, true)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}

	// Empty superseders: loop body never runs → not applicable.
	if ec.anySupersederApplicable(nil, ctx) {
		t.Fatalf("empty superseders must not be applicable")
	}
	// a_eff > 0 → applicable.
	if !ec.anySupersederApplicable([]string{w.id}, ctx) {
		t.Fatalf("superseder with a_eff 1 must report applicable")
	}
	// Inactive superseder → not applicable.
	if ec.anySupersederApplicable([]string{r.id}, ctx) {
		t.Fatalf("inactive superseder must not report applicable")
	}
	// Active (v > -1) but gated to a_eff 0 by an exception → not applicable.
	if !ec.active[x.id] || ec.applicability[x.id] != 0 {
		t.Fatalf("precondition: x active with a_eff 0, got active=%v a_eff=%v", ec.active[x.id], ec.applicability[x.id])
	}
	if ec.anySupersederApplicable([]string{x.id}, ctx) {
		t.Fatalf("fully gated superseder must not report applicable")
	}
	// Any one applicable superseder among several decides the result.
	if !ec.anySupersederApplicable([]string{r.id, x.id, w.id}, ctx) {
		t.Fatalf("one applicable superseder among several must report applicable")
	}
}
