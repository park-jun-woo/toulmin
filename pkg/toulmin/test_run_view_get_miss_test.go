//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewGetMiss — view.Get for an unknown name returns the zero event and false
package toulmin

import "testing"

func TestRunViewGetMiss(t *testing.T) {
	var got NodeEvent
	var ok bool
	probe := func(ctx Context, ev NodeEvent, view RunView) error {
		got, ok = view.Get("nonexistent")
		return nil
	}
	g := NewGraph("miss")
	g.Rule(WarrantA).OnActive(probe).OnDefeated(probe).OnInactive(probe)

	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if ok {
		t.Error("Get for an unknown name must return false")
	}
	if got != (NodeEvent{}) {
		t.Errorf("Get miss must return the zero NodeEvent, got %+v", got)
	}
}
