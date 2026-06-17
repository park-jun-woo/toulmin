//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewPostHoc — the RunView returned by Run is queryable after Run completes
package toulmin

import "testing"

func TestRunViewPostHoc(t *testing.T) {
	g := NewGraph("posthoc")
	w := g.Rule(authenticate)
	c := g.Counter(blockIP)
	e := g.Except(exemptInternalIP)
	c.Attacks(w)
	e.Attacks(c)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	ctx.Set("ip_blocked", true)
	ctx.Set("ip_internal", true)
	_, view, err := g.Run(ctx)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if view == nil {
		t.Fatal("Run must return a RunView")
	}
	if len(view.All()) != 3 {
		t.Fatalf("returned view should expose all 3 nodes, got %d", len(view.All()))
	}
	// Internal IP exempts the block, so block ip is defeated and authenticate prevails.
	if ev, ok := view.Get("exemptInternalIP"); !ok || ev.Type != Active {
		t.Errorf("exemptInternalIP want Active, got %v ok=%v", ev.Type, ok)
	}
	if ev, ok := view.Get("blockIP"); !ok || ev.Type != Defeated {
		t.Errorf("blockIP want Defeated, got %v ok=%v", ev.Type, ok)
	}
	if ev, ok := view.Get("authenticate"); !ok || ev.Type != Active {
		t.Errorf("authenticate want Active, got %v ok=%v", ev.Type, ok)
	}
	if got := view.Attackers("blockIP"); len(got) != 1 || got[0].Name != "exemptInternalIP" {
		t.Errorf("blockIP attacker want exemptInternalIP, got %+v", got)
	}
}
