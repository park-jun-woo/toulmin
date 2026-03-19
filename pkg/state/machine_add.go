//ff:func feature=state type=engine control=sequence
//ff:what Add: 전이 graph 등록
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Add registers a transition graph.
func (m *Machine) Add(from, event, to string, g *toulmin.Graph) {
	key := from + ":" + event
	m.transitions[key] = &transition{from: from, event: event, to: to, graph: g}
	m.order = append(m.order, key)
}
