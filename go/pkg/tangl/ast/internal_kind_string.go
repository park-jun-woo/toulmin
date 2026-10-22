//ff:func feature=tangl type=model control=selection
//ff:what String — human-readable name of an InternalKind
package ast

// String returns the human-readable name of the InternalKind.
func (k InternalKind) String() string {
	switch k {
	case OnEvent:
		return "OnEvent"
	case EveryTick:
		return "EveryTick"
	default:
		return "InternalKind(?)"
	}
}
