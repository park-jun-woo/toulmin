//ff:func feature=state type=engine control=sequence
//ff:what TestNewMachine — tests NewMachine returns an initialized Machine
package state

import "testing"

func TestNewMachine(t *testing.T) {
	m := NewMachine()
	if m == nil {
		t.Fatal("expected non-nil Machine")
	}
	if m.transitions == nil {
		t.Error("expected initialized transitions map")
	}
	if len(m.transitions) != 0 {
		t.Errorf("expected empty transitions map, got %d entries", len(m.transitions))
	}
}
