//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleWithSupersedes — With() rewrites supersede from/to regardless of declaration order
package toulmin

import "testing"

func TestRuleWithSupersedes(t *testing.T) {
	// Supersedes declared before .With() on either endpoint: both from and to must be
	// rewritten to the new ids so no stale id lingers in the supersede graph.
	g := NewGraph("with-sup")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.Supersedes(w) // from=r.id(old), to=w.id(old)

	oldFrom := r.id
	oldTo := w.id

	r.With(&testSpec{Value: "a"}) // rewrites from
	w.With(&testSpec{Value: "b"}) // rewrites to

	if r.id == oldFrom || w.id == oldTo {
		t.Fatalf("expected ids to change after With")
	}
	if len(g.supersedes) != 1 {
		t.Fatalf("expected 1 supersede edge, got %d", len(g.supersedes))
	}
	e := g.supersedes[0]
	if e.from != r.id {
		t.Errorf("expected supersede.from rewritten to %q, got %q", r.id, e.from)
	}
	if e.to != w.id {
		t.Errorf("expected supersede.to rewritten to %q, got %q", w.id, e.to)
	}
	if e.from == oldFrom || e.to == oldTo {
		t.Errorf("stale id lingering in supersede edge: %+v", e)
	}
}
