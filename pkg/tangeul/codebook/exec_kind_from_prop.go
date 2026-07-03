//ff:func feature=tangl type=util control=selection
//ff:what ExecKindFromProp — map an execution PropCode back to its ast.ExecKind
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// ExecKindFromProp maps an execution-edge PropCode back to its ast.ExecKind
// (§2.3), the inverse of ExecPropCode. Only PropDo/PropUndo/PropRun are
// execution edges; any other PropCode is a hard error — never a silent
// fallback.
func ExecKindFromProp(p PropCode) (ast.ExecKind, error) {
	switch p {
	case PropDo:
		return ast.DoExec, nil
	case PropUndo:
		return ast.UndoExec, nil
	case PropRun:
		return ast.RunExec, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: PropCode 0x%04X is not an execution edge", uint16(p))
	}
}
