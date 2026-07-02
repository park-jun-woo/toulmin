//ff:func feature=state type=engine control=sequence
//ff:what TestMermaid_Empty — covers Machine.Mermaid loop-not-executed branch on an empty machine
package state

import "testing"

func TestMermaid_Empty(t *testing.T) {
	m := NewMachine()

	diagram := m.Mermaid()
	if diagram != "stateDiagram-v2\n" {
		t.Errorf("expected header-only diagram, got %q", diagram)
	}
}
