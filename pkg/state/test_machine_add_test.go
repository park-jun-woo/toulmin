//ff:func feature=state type=engine control=sequence
//ff:what TestMachineAdd — verifies Machine.Add registers a transition keyed by "from:event" and records insertion order
package state

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMachineAdd(t *testing.T) {
	m := NewMachine()
	g := toulmin.NewGraph("test:add")

	m.Add("pending", "accept", "accepted", g)

	tr, ok := m.transitions["pending:accept"]
	if !ok {
		t.Fatal("expected transition registered under key \"pending:accept\"")
	}
	if tr.from != "pending" || tr.event != "accept" || tr.to != "accepted" || tr.graph != g {
		t.Errorf("unexpected transition: %+v", tr)
	}
	if len(m.order) != 1 || m.order[0] != "pending:accept" {
		t.Errorf("expected order [\"pending:accept\"], got %v", m.order)
	}

	g2 := toulmin.NewGraph("test:add2")
	m.Add("pending", "reject", "rejected", g2)
	if len(m.order) != 2 || m.order[1] != "pending:reject" {
		t.Errorf("expected order to append \"pending:reject\", got %v", m.order)
	}
}
