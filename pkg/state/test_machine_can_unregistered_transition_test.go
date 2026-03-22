//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_Can_UnregisteredTransition — tests error for unregistered transition
package state

import "testing"

func TestMachine_Can_UnregisteredTransition(t *testing.T) {
	m := NewMachine()

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	_, err := m.Can(req, ctx)
	if err == nil {
		t.Fatal("expected error for unregistered transition")
	}
}
