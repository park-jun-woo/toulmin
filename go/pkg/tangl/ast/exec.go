//ff:type feature=tangl type=model
//ff:what Exec — an execution edge (do/undo/run) attached to a case node
package ast

// Exec is a single execution edge: do, undo, or run.
type Exec struct {
	Kind      ExecKind   `json:"kind"`
	Func      *Ref       `json:"func,omitempty"` // do/undo target
	Case      string     `json:"case,omitempty"` // run target case
	Node      string     `json:"node"`           // when <node> trigger
	Once      bool       `json:"once,omitempty"`
	Certainty *Certainty `json:"certainty,omitempty"`
	Line      int        `json:"line"`
}
