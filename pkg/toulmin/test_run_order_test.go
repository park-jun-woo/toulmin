//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunOrder — RunOn firing order equals rule registration order (Active nodes only)
package toulmin

import "testing"

func TestRunOrder(t *testing.T) {
	var order []string
	rec := func(name string) NodeHandler {
		return func(t Trace) error {
			// self = this handler's own node, looked up by its known name.
			self, _ := t.Get(name)
			order = append(order, self.Name)
			return nil
		}
	}
	g := NewGraph("order")
	g.Rule(authenticate).RunOn(rec("authenticate"))
	g.Rule(WarrantA).RunOn(rec("WarrantA"))
	g.Rule(RebuttalB).RunOn(rec("RebuttalB"))

	ctx := NewContext()
	ctx.Set("authenticated", true)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	// All three rules are Active (authenticated=true, WarrantA/RebuttalB always true,
	// none attacked) → all fire in registration order.
	want := []string{"authenticate", "WarrantA", "RebuttalB"}
	if len(order) != len(want) {
		t.Fatalf("fired %d handlers, want %d: %v", len(order), len(want), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order[%d] = %q, want %q", i, order[i], want[i])
		}
	}
}
