//ff:func feature=tangl type=util control=selection
//ff:what LevelFromAST — map an ast.Level to its 2-bit LevelCode, erroring on unknown
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// LevelFromAST maps an ast.Level to its codebook LevelCode (§5.5). It preserves
// LevelUnset verbatim (no default injection — the unwritten level is itself
// information). An unrecognized Level is a hard error — never a silent fallback.
func LevelFromAST(l ast.Level) (LevelCode, error) {
	switch l {
	case ast.LevelUnset:
		return LevelUnset, nil
	case ast.LevelFail:
		return LevelFail, nil
	case ast.LevelReview:
		return LevelReview, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown ast.Level %d", int(l))
	}
}
