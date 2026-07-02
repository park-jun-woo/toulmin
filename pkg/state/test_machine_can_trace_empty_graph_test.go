//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_CanTrace_EmptyGraph — tests CanTrace returns -1 verdict with no trace when graph has no rules
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_CanTrace_EmptyGraph(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	// no rules registered on g

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	result, err := m.CanTrace(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Verdict != -1 {
		t.Errorf("expected verdict -1 for empty graph, got %f", result.Verdict)
	}
	if result.From != "pending" || result.To != "accepted" || result.Event != "accept" {
		t.Errorf("unexpected result: %+v", result)
	}
	if len(result.Trace) != 0 {
		t.Errorf("expected empty trace, got %+v", result.Trace)
	}
}
