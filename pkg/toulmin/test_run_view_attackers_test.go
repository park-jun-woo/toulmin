//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewAttackers — view.Attackers returns the events of nodes that attacked a node
package toulmin

import "testing"

func TestRunViewAttackers(t *testing.T) {
	var attackers []NodeEvent
	var leafAttackers int
	capture := func(ctx Context, ev NodeEvent, view RunView) error {
		attackers = view.Attackers("authenticate")
		leafAttackers = len(view.Attackers("blockIP"))
		return nil
	}
	g := NewGraph("access")
	w := g.Rule(authenticate).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	c := g.Counter(blockIP)
	c.Attacks(w)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	ctx.Set("ip_blocked", true)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(attackers) != 1 {
		t.Fatalf("authenticate should have 1 attacker, got %d", len(attackers))
	}
	if attackers[0].Name != "blockIP" {
		t.Errorf("attacker want blockIP, got %q", attackers[0].Name)
	}
	if attackers[0].Type != Active {
		t.Errorf("blockIP attack should be Active, got %v", attackers[0].Type)
	}
	if leafAttackers != 0 {
		t.Errorf("blockIP has no attackers, want 0, got %d", leafAttackers)
	}
}
