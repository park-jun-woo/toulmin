//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_Can_StateMismatch — tests Machine rejects wrong current state
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_Can_StateMismatch(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	g.Rule(IsCurrentState)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "accepted"} // wrong state

	verdict, err := m.Can(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verdict > 0 {
		t.Errorf("expected verdict <= 0 for state mismatch, got %f", verdict)
	}
}
