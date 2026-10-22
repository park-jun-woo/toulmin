//ff:type feature=tangl type=model
//ff:what caseState — the Writer's open-case tracker while walking Blocks
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// caseState tracks the case currently wrapped by a CASE_BEGIN as WriteDoc walks
// Blocks: whether one is open, its Blocks index and its TID. Leaving the case (a
// block outside it, or a different case) closes it with a CASE_END.
type caseState struct {
	open bool
	idx  int
	tid  toulmin.TID
}
