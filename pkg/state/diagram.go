//ff:func feature=state type=util control=iteration dimension=1
//ff:what Mermaid: Machine에 등록된 전이 목록에서 Mermaid stateDiagram 생성
package state

import (
	"fmt"
	"strings"
)

// Mermaid returns a Mermaid stateDiagram-v2 string from registered transitions.
func (m *Machine) Mermaid() string {
	var b strings.Builder
	b.WriteString("stateDiagram-v2\n")
	for _, key := range m.order {
		t := m.transitions[key]
		fmt.Fprintf(&b, "    %s --> %s : %s\n", t.from, t.to, t.event)
	}
	return b.String()
}
