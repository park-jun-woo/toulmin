//ff:func feature=tangl type=model control=selection
//ff:what String — human-readable name of a DefKind
package ast

// String returns the human-readable name of the DefKind.
func (k DefKind) String() string {
	switch k {
	case ConstDef:
		return "ConstDef"
	case StructDef:
		return "StructDef"
	default:
		return "DefKind(?)"
	}
}
