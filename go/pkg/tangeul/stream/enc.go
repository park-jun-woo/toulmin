//ff:type feature=tangl type=model
//ff:what enc — the Writer's state: pass-1 TID table, SPAN width, next free Edge Id and the first encoding error
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// enc is the mutable state threaded through the Writer while it serializes a
// Doc. tt holds the pass-1 declaration-order Edge Ids of every primary entity
// (so forward references resolve); wire is the document's SPAN width; next
// issues the Edge Ids of every other record (SOURCE, TEXT, LABEL and relation
// records, Func nodes, framing controls) from primaryEnd upward, so every
// record of the document has a unique Edge Id (GEUL t002 TOULMIN). err holds
// the first record-encoding error: the record helpers record it instead of
// returning it, and emitBody reports it once the walk ends.
type enc struct {
	tt   *tidTable
	wire toulmin.Wire
	next toulmin.TID
	err  error
}
