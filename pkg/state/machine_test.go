package state

import (
	"testing"
	"time"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func isAuth(claim any, ground any, backing any) (bool, any) { return true, nil }
func isAdmin(claim any, ground any, backing any) (bool, any) { return true, nil }

func TestMachine_Can_Allowed(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	g.Warrant(IsCurrentState, nil, 1.0)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	verdict, err := m.Can(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verdict <= 0 {
		t.Errorf("expected verdict > 0, got %f", verdict)
	}
}

func TestMachine_Can_StateMismatch(t *testing.T) {
	m := NewMachine()

	g := toulmin.NewGraph("proposal:accept")
	g.Warrant(IsCurrentState, nil, 1.0)

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

func TestMachine_Can_ExpiredWithOverride(t *testing.T) {
	m := NewMachine()
	expiryFunc := func(r any) time.Time { return r.(*testResource).ExpiresAt }

	g := toulmin.NewGraph("proposal:accept")
	current := g.Warrant(IsCurrentState, nil, 1.0)
	expired := g.Rebuttal(IsExpired, expiryFunc, 1.0)
	override := g.Defeater(isAdmin, nil, 1.0)
	g.Defeat(expired, current)
	g.Defeat(override, expired)

	m.Add("pending", "accept", "accepted", g)

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{
		CurrentState: "pending",
		Resource:     &testResource{ExpiresAt: time.Now().Add(-1 * time.Hour)}, // expired
	}

	verdict, err := m.Can(req, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verdict <= 0 {
		t.Errorf("expected verdict > 0 (admin override defeats expired), got %f", verdict)
	}
}

func TestMachine_Can_UnregisteredTransition(t *testing.T) {
	m := NewMachine()

	req := &TransitionRequest{From: "pending", To: "accepted", Event: "accept"}
	ctx := &TransitionContext{CurrentState: "pending"}

	_, err := m.Can(req, ctx)
	if err == nil {
		t.Fatal("expected error for unregistered transition")
	}
}

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
