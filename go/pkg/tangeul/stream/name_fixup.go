//ff:type feature=tangl type=model
//ff:what nameFixup — a deferred case-name reference (AST coordinate + the case TID that resolves it)
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// nameFixup records a forward case-name reference to patch once every case name
// is known: the AST coordinate (case index, slot index within Nodes or Execs)
// and the case TID whose recovered name fills the field. `run` targets and
// `checking` targets can point at a case declared later in the stream, so their
// names cannot be resolved during the single linear build pass.
type nameFixup struct {
	caseIdx int
	slotIdx int
	tid     toulmin.TID
}
