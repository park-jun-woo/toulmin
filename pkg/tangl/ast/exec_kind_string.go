//ff:func feature=tangl type=model control=selection
//ff:what String — human-readable name of an ExecKind
package ast

// String returns the human-readable name of the ExecKind.
func (k ExecKind) String() string {
	switch k {
	case DoExec:
		return "DoExec"
	case UndoExec:
		return "UndoExec"
	case RunExec:
		return "RunExec"
	default:
		return "ExecKind(?)"
	}
}
