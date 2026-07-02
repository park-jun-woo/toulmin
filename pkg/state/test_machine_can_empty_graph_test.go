//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_Can_EmptyGraph — tests Machine.Can returns -1 with no error when graph has no rules
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_Can_EmptyGraph(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	// no rules registered on g

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	verdict, err := m.Can(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verdict != -1 {
		t.Errorf("expected verdict -1 for empty graph, got %f", verdict)
	}
}
