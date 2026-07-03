//ff:func feature=engine type=engine control=sequence
//ff:what TestAnySupersederActive — empty/inactive report false, active reports true (all branches)
package toulmin

import "testing"

func TestAnySupersederActive(t *testing.T) {
	ctx := NewContext()
	g := NewGraph("any-sup")
	w := g.Rule(WarrantA)     // active (returns true)
	r := g.Counter(InactiveR) // inactive (returns false)
	_, ec, err := g.evaluate(ctx, EvalOption{}, true)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}

	// Empty superseders: loop body never runs → not active.
	if ec.anySupersederActive(nil, ctx) {
		t.Fatalf("empty superseders must not be active")
	}
	// Active superseder → active (if ec.active[from] taken).
	if !ec.anySupersederActive([]string{w.id}, ctx) {
		t.Fatalf("active superseder must report active")
	}
	// Inactive superseder → not active (if ec.active[from] not taken).
	if ec.anySupersederActive([]string{r.id}, ctx) {
		t.Fatalf("inactive superseder must not report active")
	}
}
