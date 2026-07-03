//ff:func feature=tangl type=util control=sequence
//ff:what exprCompareRecord — build one EXPR Compare record (equals, field 0, an inline NUM32 value of 0) for hand-built test fixtures
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// exprCompareRecord builds one EXPR Compare record (equals, field 0, an
// inline NUM32 value of 0) — the leaf term reused by parseExpr's Logic/Not/
// Either success fixtures.
func exprCompareRecord() []word.Word {
	num := word.EncodeNum32(0)
	return []word.Word{exprRecHeader(ExprCompare, 0), word.Word(codebook.CompareEquals), 0, word.Word(ValueNum), num[0], num[1]}
}
