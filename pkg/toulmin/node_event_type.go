//ff:type feature=engine type=model
//ff:what NodeEventType — node event classification (inactive/active/defeated)
package toulmin

// NodeEventType classifies a node's post-evaluation event.
type NodeEventType int

const (
	// Inactive means the rule function returned false (rule did not apply).
	Inactive NodeEventType = iota
	// Active means the function returned true and verdict > 0 (applied and prevailed).
	Active
	// Defeated means the function returned true and verdict <= 0 (applied but defeated).
	Defeated
)
