//ff:func feature=engine type=util control=selection
//ff:what String — returns the name of a NodeEventType for logs and errors
package toulmin

// String returns the event type name (Active, Defeated, or Inactive).
func (t NodeEventType) String() string {
	switch t {
	case Active:
		return "Active"
	case Defeated:
		return "Defeated"
	default:
		return "Inactive"
	}
}
