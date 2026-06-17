//ff:func feature=engine type=engine control=selection
//ff:what selectHandler — picks the handler matching a node event type
package toulmin

// selectHandler returns the handler registered for the given event type (nil if none).
func selectHandler(r *RuleMeta, t NodeEventType) NodeHandler {
	switch t {
	case Active:
		return r.OnActive
	case Defeated:
		return r.OnDefeated
	default:
		return r.OnInactive
	}
}
