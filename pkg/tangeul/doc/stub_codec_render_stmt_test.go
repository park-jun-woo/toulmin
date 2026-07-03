//ff:func feature=tangl type=util control=selection
//ff:what stubCodec.RenderStmt — canned single-line and multi-line statements per section
package doc

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// RenderStmt returns canned statement text: a Cases declaration or body node as a
// single line, and a Rules inline rule as a three-line subtree (with relative
// indentation) so Render's multi-line marker/shift path is exercised.
func (s stubCodec) RenderStmt(d *Doc, ref StmtRef) (string, error) {
	switch ref.Sec {
	case codebook.SectionCases:
		if len(ref.Path) == 1 {
			return "in case of `x`", nil
		}
		return "`a` is a general rule", nil
	case codebook.SectionRules:
		return "`r` when\n`a` is 5\nand `b`", nil
	default:
		return "", fmt.Errorf("stubCodec: no canned render for section %d", uint16(ref.Sec))
	}
}
