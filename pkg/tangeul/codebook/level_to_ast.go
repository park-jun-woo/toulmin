//ff:func feature=tangl type=util control=selection
//ff:what LevelToAST — map a LevelCode back to its ast.Level, erroring on unknown
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// LevelToAST maps a codebook LevelCode back to its ast.Level (§5.5), the inverse
// of LevelFromAST. It preserves LevelUnset verbatim (no default injection). An
// unrecognized code is a hard error — never a silent fallback.
func LevelToAST(c LevelCode) (ast.Level, error) {
	switch c {
	case LevelUnset:
		return ast.LevelUnset, nil
	case LevelFail:
		return ast.LevelFail, nil
	case LevelReview:
		return ast.LevelReview, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown LevelCode %d", uint16(c))
	}
}
