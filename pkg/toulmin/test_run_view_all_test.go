//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunViewAll — a handler sees every node via view.All() during Run
package toulmin

import "testing"

func TestRunViewAll(t *testing.T) {
	var seen map[string]bool
	capture := func(ctx Context, ev NodeEvent, view RunView) error {
		seen = map[string]bool{}
		for _, n := range view.All() {
			seen[n.Name] = true
		}
		return nil
	}
	g := NewGraph("access")
	w := g.Rule(authenticate).OnActive(capture).OnDefeated(capture).OnInactive(capture)
	c := g.Counter(blockIP)
	e := g.Except(exemptInternalIP)
	c.Attacks(w)
	e.Attacks(c)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(seen) != 3 {
		t.Fatalf("view.All() should expose all 3 nodes, got %d: %v", len(seen), seen)
	}
	for _, name := range []string{"authenticate", "blockIP", "exemptInternalIP"} {
		if !seen[name] {
			t.Errorf("view.All() missing node %q", name)
		}
	}
}
