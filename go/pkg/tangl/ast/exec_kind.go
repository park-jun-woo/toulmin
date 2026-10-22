//ff:type feature=tangl type=model
//ff:what ExecKind — kind of an execution edge (do/undo/run)
package ast

// ExecKind classifies an Exec edge.
type ExecKind int

const (
	// DoExec is a `do <func> [once] when <node>` leaf action edge.
	DoExec ExecKind = iota
	// UndoExec is an `undo <func> when <node>` compensation edge.
	UndoExec
	// RunExec is a `run <case> when <node>` composite cascade edge.
	RunExec
)
