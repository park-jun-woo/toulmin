//ff:func feature=tangl type=util control=selection
//ff:what resolveInternalSlot — resolve an internal body slot (1 Run·2 Check; 0 unused)
package doc

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// resolveInternalSlot returns the slotIdx-th element of the internal body slice
// selected by slotKind (1 Run, 2 Check). ast.Internal has no Requires, so
// slotKind 0 is unused (README contract §2).
func resolveInternalSlot(in *ast.Internal, slotKind, slotIdx int) (any, error) {
	switch slotKind {
	case 1:
		if slotIdx < 0 || slotIdx >= len(in.Runs) {
			return nil, fmt.Errorf("tangeul/doc: Internal Run index %d out of range (%d)", slotIdx, len(in.Runs))
		}
		return in.Runs[slotIdx], nil
	case 2:
		if slotIdx < 0 || slotIdx >= len(in.Checks) {
			return nil, fmt.Errorf("tangeul/doc: Internal Check index %d out of range (%d)", slotIdx, len(in.Checks))
		}
		return in.Checks[slotIdx], nil
	default:
		return nil, fmt.Errorf("tangeul/doc: unknown Internal slotKind %d (0 unused)", slotKind)
	}
}
