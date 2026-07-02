//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_Can_CycleError — tests Machine.Can propagates Evaluate error on circular defeat graph
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_Can_CycleError(t *testing.T) {
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

	verdict, err := m.Can(req, ctx)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
	if verdict != -1 {
		t.Errorf("expected verdict -1 on error, got %f", verdict)
	}
}
