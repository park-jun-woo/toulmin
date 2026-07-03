//ff:func feature=tangl type=util control=selection
//ff:what ExecPropCode — map an ast.ExecKind (do/undo/run) to its Edge PropCode
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// ExecPropCode maps an ast.ExecKind to the PropCode of its execution Edge
// (§2.3): DoExec→PropDo, UndoExec→PropUndo, RunExec→PropRun. An unrecognized
// kind is a hard error — never a silent fallback.
func ExecPropCode(k ast.ExecKind) (PropCode, error) {
	switch k {
	case ast.DoExec:
		return PropDo, nil
	case ast.UndoExec:
		return PropUndo, nil
	case ast.RunExec:
		return PropRun, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown ast.ExecKind %d", int(k))
	}
}
