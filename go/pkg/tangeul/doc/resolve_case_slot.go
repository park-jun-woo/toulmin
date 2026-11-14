//ff:func feature=tangl type=util control=selection
//ff:what resolveCaseSlot — resolve a case body slot (0 Require·1 Node·2 Attack·3 Exec·4 Supersede)
package doc

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// resolveCaseSlot returns the slotIdx-th element of the case body slice selected
// by slotKind (0 Require, 1 Node, 2 Attack, 3 Exec, 4 Supersede), restoring the
// document order that ast.Case's split slices cannot (README contract §2).
func resolveCaseSlot(c *ast.Case, slotKind, slotIdx int) (any, error) {
	switch slotKind {
	case 0:
		if slotIdx < 0 || slotIdx >= len(c.Requires) {
			return nil, fmt.Errorf("tangeul/doc: Case Require index %d out of range (%d)", slotIdx, len(c.Requires))
		}
		return c.Requires[slotIdx], nil
	case 1:
		if slotIdx < 0 || slotIdx >= len(c.Nodes) {
			return nil, fmt.Errorf("tangeul/doc: Case Node index %d out of range (%d)", slotIdx, len(c.Nodes))
		}
		return c.Nodes[slotIdx], nil
	case 2:
		if slotIdx < 0 || slotIdx >= len(c.Attacks) {
			return nil, fmt.Errorf("tangeul/doc: Case Attack index %d out of range (%d)", slotIdx, len(c.Attacks))
		}
		return c.Attacks[slotIdx], nil
	case 3:
		if slotIdx < 0 || slotIdx >= len(c.Execs) {
			return nil, fmt.Errorf("tangeul/doc: Case Exec index %d out of range (%d)", slotIdx, len(c.Execs))
		}
		return c.Execs[slotIdx], nil
	case 4:
		if slotIdx < 0 || slotIdx >= len(c.Supersedes) {
			return nil, fmt.Errorf("tangeul/doc: Case Supersede index %d out of range (%d)", slotIdx, len(c.Supersedes))
		}
		return c.Supersedes[slotIdx], nil
	default:
		return nil, fmt.Errorf("tangeul/doc: unknown Case slotKind %d", slotKind)
	}
}
