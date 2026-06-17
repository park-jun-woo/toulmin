//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewGet — runView.Get returns (event, true) for a known name and (zero, false) for an unknown one
package toulmin

import "testing"

func TestRunViewGet(t *testing.T) {
	var hit NodeEvent
	var hitOK bool
	var miss NodeEvent
	var missOK bool
	probe := func(ctx Context, ev NodeEvent, view RunView) error {
		hit, hitOK = view.Get("WarrantA")
		miss, missOK = view.Get("nonexistent")
		return nil
	}
	g := NewGraph("get")
	g.Rule(WarrantA).OnActive(probe).OnDefeated(probe).OnInactive(probe)

	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if !hitOK {
		t.Errorf("present: ok = false, want true")
	}
	if hit.Name != "WarrantA" {
		t.Errorf("present: name = %q, want %q", hit.Name, "WarrantA")
	}
	if missOK {
		t.Errorf("absent: ok = true, want false")
	}
	if miss != (NodeEvent{}) {
		t.Errorf("absent: want zero NodeEvent, got %+v", miss)
	}
}
