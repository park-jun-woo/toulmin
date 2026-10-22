//ff:func feature=tangl type=util control=sequence
//ff:what exprCompareRecord — build one EXPR Compare record (equals, field 0, an inline NUM32 value of 0) for hand-built test fixtures
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// exprCompareRecord builds one EXPR Compare record (equals, field 0, an
// inline NUM32 value of 0) — the leaf term reused by parseExpr's Logic/Not/
// Either success fixtures.
func exprCompareRecord() []core.Word {
	num := uintVal(0)
	return []core.Word{exprRecHeader(ExprCompare, 0), core.Word(toulmin.CompareEquals), 0, core.Word(ValueNum), num[0], num[1]}
}
