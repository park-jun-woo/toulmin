//ff:func feature=state type=engine control=sequence
//ff:what TestMachine_Can_ExpiredWithOverride — tests admin override defeats expired rebuttal
package state

import (
	"testing"
	"time"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachine_Can_ExpiredWithOverride(t *testing.T) {
	m := NewMachine()
	expiresAt := time.Now().Add(-1 * time.Hour) // expired

	g := toulmin.NewGraph("proposal:accept")
	current := g.Warrant(IsCurrentState, nil, 1.0)
	expired := g.Rebuttal(IsExpired, &ExpiryBacking{ExpiresAt: expiresAt}, 1.0)
	override := g.Defeater(isAdmin, nil, 1.0)
	g.Defeat(expired, current)
	g.Defeat(override, expired)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{
		CurrentState: "pending",
		Resource:     &testResource{ExpiresAt: expiresAt},
	}

	verdict, err := m.Can(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verdict <= 0 {
		t.Errorf("expected verdict > 0 (admin override defeats expired), got %f", verdict)
	}
}
