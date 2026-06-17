//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewImmutable — a handler mutating ctx never changes another handler's view
package toulmin

import "testing"

func TestRunViewImmutable(t *testing.T) {
	// First handler mutates ctx; the snapshot must stay fixed for the later handler.
	mutate := func(ctx Context, ev NodeEvent, view RunView) error {
		ctx.Set("authenticated", false)
		return nil
	}
	var laterType NodeEventType
	var laterFound bool
	inspect := func(ctx Context, ev NodeEvent, view RunView) error {
		got, ok := view.Get("authenticate")
		laterFound = ok
		laterType = got.Type
		return nil
	}
	g := NewGraph("access")
	g.Rule(authenticate).OnActive(mutate).OnDefeated(mutate).OnInactive(mutate)
	g.Counter(blockIP).OnActive(inspect).OnDefeated(inspect).OnInactive(inspect)

	ctx := NewContext()
	ctx.Set("authenticated", true)
	if _, _, err := g.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !laterFound {
		t.Fatal("later handler should find authenticate in the snapshot")
	}
	if laterType != Active {
		t.Errorf("snapshot must be immutable: authenticate want Active, got %v", laterType)
	}
}
