//ff:func feature=engine type=engine control=sequence
//ff:what Get — looks up a node's final event by short name
package toulmin

// Get returns the final event for the node with the given short name, or (zero, false).
func (v *runView) Get(name string) (NodeEvent, bool) {
	ev, ok := v.byName[name]
	return ev, ok
}
