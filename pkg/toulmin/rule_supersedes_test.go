//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleSupersedes — tests Rule.Supersedes appends a supersede edge from r to target
package toulmin

import "testing"

func TestRuleSupersedes(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.Supersedes(w)
	if len(g.supersedes) != 1 {
		t.Fatalf("expected 1 supersede edge, got %d", len(g.supersedes))
	}
	if g.supersedes[0].from != r.id || g.supersedes[0].to != w.id {
		t.Fatalf("expected edge from %q to %q, got %+v", r.id, w.id, g.supersedes[0])
	}
	// Supersede must not leak into the Attacks (defeats) graph.
	if len(g.defeats) != 0 {
		t.Fatalf("supersede must not touch defeats graph, got %d defeats", len(g.defeats))
	}
}
