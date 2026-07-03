//ff:func feature=tangl type=util control=sequence
//ff:what exprRecHeader — pack an EXPR record header word from an ExprKind and a child count for hand-built test fixtures
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// exprRecHeader packs an EXPR record header word from an ExprKind and a child
// count, matching parseExpr's bit layout (kind bits 15-14, child count bits
// 13-0).
func exprRecHeader(kind ExprKind, childCount int) word.Word {
	return word.Word(uint16(kind)<<14 | uint16(childCount))
}
