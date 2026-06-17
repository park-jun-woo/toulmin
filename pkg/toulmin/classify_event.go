//ff:func feature=engine type=engine control=selection
//ff:what classifyEvent — classifies a node event from its active flag and verdict
package toulmin

// classifyEvent maps a node's active flag and verdict to a NodeEventType.
// verdict == 0 (balanced) counts as Defeated; Active requires verdict > 0.
func classifyEvent(active bool, verdict float64) NodeEventType {
	switch {
	case !active:
		return Inactive
	case verdict > 0:
		return Active
	default:
		return Defeated
	}
}
