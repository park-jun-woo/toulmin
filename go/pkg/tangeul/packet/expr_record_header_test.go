//ff:func feature=tangl type=util control=sequence
//ff:what exprRecHeader — pack an EXPR record header word from an ExprKind and a child count for hand-built test fixtures
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
)

// exprRecHeader packs an EXPR record header word from an ExprKind and a child
// count via packExprHeader (kind split across bits 15-14 and bit 13, child count
// bits 12-0). It panics on an out-of-range count since fixtures are always
// well-formed.
func exprRecHeader(kind ExprKind, childCount int) core.Word {
	w, err := packExprHeader(kind, childCount)
	if err != nil {
		panic(err)
	}
	return w
}
