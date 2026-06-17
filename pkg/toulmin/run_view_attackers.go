//ff:func feature=engine type=engine control=sequence
//ff:what Attackers — returns a copy of the final events of nodes that attacked name
package toulmin

// Attackers returns a copy of the final events of the nodes that attacked name
// (empty if none). name is the target node's short name.
func (v *runView) Attackers(name string) []NodeEvent {
	src := v.attackers[name]
	out := make([]NodeEvent, len(src))
	copy(out, src)
	return out
}
