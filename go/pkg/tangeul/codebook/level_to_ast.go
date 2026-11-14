//ff:func feature=tangl type=util control=selection
//ff:what LevelToAST — map a LevelCode back to its ast.Level, erroring on unknown
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// LevelToAST maps a codebook LevelCode back to its ast.Level (§5.5), the inverse
// of LevelFromAST. It preserves LevelUnset verbatim (no default injection). An
// unrecognized code is a hard error — never a silent fallback.
func LevelToAST(c toulmin.LevelCode) (ast.Level, error) {
	switch c {
	case toulmin.LevelUnset:
		return ast.LevelUnset, nil
	case toulmin.LevelFail:
		return ast.LevelFail, nil
	case toulmin.LevelReview:
		return ast.LevelReview, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown toulmin.LevelCode %d", uint16(c))
	}
}
