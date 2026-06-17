//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestNewRunView — the run snapshot classifies active/defeated/inactive nodes and builds the attacker index over every edge
package toulmin

import "testing"

func TestNewRunView(t *testing.T) {
	var got RunView
	capture := func(ctx Context, ev NodeEvent, view RunView) error {
		got = view
		return nil
	}
	// Access-control graph with two attack edges plus an unrelated inactive
	// node, so newRunView's node loop (active/defeated/inactive), the edge
	// loop, and the inner attacker loop are all exercised in one pass.
	g := NewGraph("access")
	w := g.Rule(authenticate).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	c := g.Counter(blockIP).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	e := g.Except(exemptInternalIP).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	g.Rule(InactiveR).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	c.Attacks(w)
	e.Attacks(c)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	ctx.Set("ip_blocked", true)
	ctx.Set("ip_internal", false)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if got == nil {
		t.Fatal("handler never received the view")
	}

	// Node loop: every registered node is present, registration order kept.
	all := got.All()
	if len(all) != 4 {
		t.Fatalf("All() want 4 nodes, got %d", len(all))
	}
	wantOrder := []string{"authenticate", "blockIP", "exemptInternalIP", "InactiveR"}
	for i, name := range wantOrder {
		if all[i].Name != name {
			t.Errorf("order[%d] = %q, want %q", i, all[i].Name, name)
		}
	}

	// Classification: active warrant defeated by an active counter, the except
	// stays inactive, and the unrelated rule is inactive too.
	auth, _ := got.Get("authenticate")
	if auth.Type != Defeated {
		t.Errorf("authenticate want Defeated, got %v", auth.Type)
	}
	block, _ := got.Get("blockIP")
	if block.Type != Active {
		t.Errorf("blockIP want Active, got %v", block.Type)
	}
	exempt, _ := got.Get("exemptInternalIP")
	if exempt.Type != Inactive {
		t.Errorf("exemptInternalIP want Inactive, got %v", exempt.Type)
	}
	idle, _ := got.Get("InactiveR")
	if idle.Type != Inactive {
		t.Errorf("InactiveR want Inactive, got %v", idle.Type)
	}

	// Edge loop + inner attacker loop: authenticate has blockIP as its attacker.
	atk := got.Attackers("authenticate")
	if len(atk) != 1 || atk[0].Name != "blockIP" {
		t.Fatalf("authenticate attackers want [blockIP], got %+v", atk)
	}
	if atk[0].Type != Active {
		t.Errorf("blockIP attacker event want Active, got %v", atk[0].Type)
	}
	if len(got.Attackers("InactiveR")) != 0 {
		t.Errorf("InactiveR has no attackers, got %d", len(got.Attackers("InactiveR")))
	}
}
