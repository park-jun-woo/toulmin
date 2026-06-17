//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunOrder — handler firing order equals rule registration order
package toulmin

import "testing"

func TestRunOrder(t *testing.T) {
	var order []string
	rec := func(ctx Context, ev NodeEvent, view RunView) error {
		order = append(order, ev.Name)
		return nil
	}
	g := NewGraph("order")
	g.Rule(authenticate).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	g.Rule(WarrantA).OnActive(rec).OnDefeated(rec).OnInactive(rec)
	g.Rule(RebuttalB).OnActive(rec).OnDefeated(rec).OnInactive(rec)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	want := []string{"authenticate", "WarrantA", "RebuttalB"}
	if len(order) != len(want) {
		t.Fatalf("fired %d events, want %d: %v", len(order), len(want), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order[%d] = %q, want %q", i, order[i], want[i])
		}
	}
}
