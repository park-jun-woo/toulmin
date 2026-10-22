//ff:func feature=tangl type=model control=selection
//ff:what String — human-readable name of a Role
package ast

// String returns the human-readable name of the Role.
func (r Role) String() string {
	switch r {
	case GeneralRule:
		return "GeneralRule"
	case CounterRule:
		return "CounterRule"
	case ExceptRule:
		return "ExceptRule"
	default:
		return "Role(?)"
	}
}
