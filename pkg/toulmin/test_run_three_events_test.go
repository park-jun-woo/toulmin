//ff:func feature=engine type=engine control=sequence
//ff:what TestRunThreeEvents — access-control graph fires Inactive/Active/Defeated correctly
package toulmin

import "testing"

func TestRunThreeEvents(t *testing.T) {
	recorded := map[string]NodeEventType{}
	rec := func(ctx Context, ev NodeEvent, view RunView) error {
		recorded[ev.Name] = ev.Type
		return nil
	}
	g := NewGraph("access")
	w := g.Rule(authenticate).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	c := g.Counter(blockIP).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	e := g.Except(exemptInternalIP).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	c.Attacks(w)
	e.Attacks(c)

	ctx1 := NewContext()
	ctx1.Set("authenticated", true)
	ctx1.Set("ip_blocked", true)
	ctx1.Set("ip_internal", false)
	if _, _, err := g.Run(ctx1); err != nil {
		t.Fatalf("run1 error: %v", err)
	}
	if recorded["blockIP"] != Active {
		t.Errorf("external: blockIP want Active, got %v", recorded["blockIP"])
	}
	if recorded["authenticate"] != Defeated {
		t.Errorf("external: authenticate want Defeated, got %v", recorded["authenticate"])
	}
	if recorded["exemptInternalIP"] != Inactive {
		t.Errorf("external: exemptInternalIP want Inactive, got %v", recorded["exemptInternalIP"])
	}

	recorded = map[string]NodeEventType{}
	ctx2 := NewContext()
	ctx2.Set("authenticated", true)
	ctx2.Set("ip_blocked", true)
	ctx2.Set("ip_internal", true)
	if _, _, err := g.Run(ctx2); err != nil {
		t.Fatalf("run2 error: %v", err)
	}
	if recorded["exemptInternalIP"] != Active {
		t.Errorf("internal: exemptInternalIP want Active, got %v", recorded["exemptInternalIP"])
	}
	if recorded["blockIP"] != Defeated {
		t.Errorf("internal: blockIP want Defeated, got %v", recorded["blockIP"])
	}
	if recorded["authenticate"] != Active {
		t.Errorf("internal: authenticate want Active, got %v", recorded["authenticate"])
	}
}
