//ff:func feature=state type=engine control=sequence
//ff:what TestMermaid — tests Mermaid diagram generation
package state

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestMermaid(t *testing.T) {
	m := NewMachine()

	g1 := toulmin.NewGraph("proposal:accept")
	g1.Rule(IsCurrentState)
	m.Add("pending", "accept", "accepted", g1)

	g2 := toulmin.NewGraph("proposal:reject")
	g2.Rule(IsCurrentState)
	m.Add("pending", "reject", "rejected", g2)

	diagram := m.Mermaid()

	if !strings.Contains(diagram, "stateDiagram-v2") {
		t.Error("expected stateDiagram-v2 header")
	}
	if !strings.Contains(diagram, "pending --> accepted : accept") {
		t.Error("expected pending --> accepted : accept")
	}
	if !strings.Contains(diagram, "pending --> rejected : reject") {
		t.Error("expected pending --> rejected : reject")
	}
}
