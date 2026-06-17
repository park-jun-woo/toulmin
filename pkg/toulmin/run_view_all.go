//ff:func feature=engine type=engine control=sequence
//ff:what All — returns a copy of every node's final event in registration order
package toulmin

// All returns a copy of all node events in registration order (Inactive included).
func (v *runView) All() []NodeEvent {
	out := make([]NodeEvent, len(v.order))
	copy(out, v.order)
	return out
}
