//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_CanTrace_CycleError — tests CanTrace propagates Evaluate error on circular defeat graph
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_CanTrace_CycleError(t *testing.T) {
	m := NewMachine()

	cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
	cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }

	g := toulmin.NewGraph("proposal:accept")
	a := g.Rule(cycleA)
	b := g.Counter(cycleB)
	b.Attacks(a)
	a.Attacks(b)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	result, err := m.CanTrace(req, ctx)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}
