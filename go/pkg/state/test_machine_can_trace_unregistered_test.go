//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_CanTrace_Unregistered — tests error for unregistered transition
package state

import "testing"

func TestMachine_CanTrace_Unregistered(t *testing.T) {
	m := NewMachine()

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	result, err := m.CanTrace(req, ctx)
	if err == nil {
		t.Fatal("expected error for unregistered transition")
	}
	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}
