//ff:func feature=tangl type=parser control=sequence
//ff:what renderCertainty — render a Certainty gate as `if <op> <N>% certain` (the inverse of parseCertaintyClause)
package en

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderCertainty renders a confidence gate back to "if <op> <N>% certain".
// Certainty.Op is already the canonical English operator ("at least"/"above"/
// "less than"/"at most"), so no reverse mapping is needed.
func renderCertainty(c ast.Certainty) string {
	return "if " + c.Op + " " + strconv.Itoa(c.Percent) + "% certain"
}
