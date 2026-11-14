//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleAttacks — tests Rule.Attacks appends a defeat edge from r to target
package toulmin

import "testing"

func TestRuleAttacks(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.Attacks(w)
	if len(g.defeats) != 1 {
		t.Fatalf("expected 1 defeat edge, got %d", len(g.defeats))
	}
	if g.defeats[0].from != r.id || g.defeats[0].to != w.id {
		t.Fatalf("expected edge from %q to %q, got %+v", r.id, w.id, g.defeats[0])
	}
}
