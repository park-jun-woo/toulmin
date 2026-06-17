//ff:func feature=engine type=engine control=sequence
//ff:what TestRunInactiveAllNodes — full pass fires Inactive for nodes lazy eval never reaches
package toulmin

import "testing"

func TestRunInactiveAllNodes(t *testing.T) {
	recorded := map[string]NodeEventType{}
	rec := func(ctx Context, ev NodeEvent, view RunView) error {
		recorded[ev.Name] = ev.Type
		return nil
	}
	g := NewGraph("test")
	g.Rule(WarrantA).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	w2 := g.Rule(InactiveR).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	c := g.Counter(blockIP).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	c.Attacks(w2)

	ctx := NewContext()
	ctx.Set("ip_blocked", false)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if recorded["WarrantA"] != Active {
		t.Errorf("WarrantA want Active, got %v", recorded["WarrantA"])
	}
	if recorded["InactiveR"] != Inactive {
		t.Errorf("InactiveR want Inactive, got %v", recorded["InactiveR"])
	}
	if recorded["blockIP"] != Inactive {
		t.Errorf("blockIP want Inactive (full pass reaches it), got %v", recorded["blockIP"])
	}
}
