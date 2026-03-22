//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_CanTrace — tests CanTrace returns verdict and trace
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_CanTrace(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	g.Warrant(IsCurrentState, nil, 1.0)
	g.Warrant(isAuth, nil, 1.0)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	result, err := m.CanTrace(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Verdict <= 0 {
		t.Errorf("expected verdict > 0, got %f", result.Verdict)
	}
	if result.From != "pending" || result.To != "accepted" || result.Event != "accept" {
		t.Errorf("unexpected result: %+v", result)
	}
	if len(result.Trace) == 0 {
		t.Error("expected non-empty trace")
	}
}
