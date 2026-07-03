//ff:func feature=tangl type=util control=selection
//ff:what emitExpr — flatten an ast.Expr subtree into prefix-order EXPR records, interning labels
package packet

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitExpr flattens e into prefix-order EXPR records (§design 4), appending
// field/value strings into *labels and returning the record words plus the node
// count of the subtree. A record header word packs the ExprKind in bits 15-14
// and the child count in bits 13-0. A Compare value that is an exact uint32
// (its decimal round-trip is identical) is stored as an inline NUM32; otherwise
// it is a Text label reference, so any value string round-trips exactly.
func emitExpr(e ast.Expr, labels *[]string) ([]word.Word, int, error) {
	switch v := e.(type) {
	case ast.Compare:
		op, err := codebook.CompareOpFromAST(v.Op)
		if err != nil {
			return nil, 0, err
		}
		out := []word.Word{
			word.Word(uint16(ExprCompare) << 14),
			word.Word(op),
			word.Word(internLabel(labels, v.Field)),
		}
		if n, err := strconv.ParseUint(v.Value, 10, 32); err == nil && strconv.FormatUint(n, 10) == v.Value {
			num := word.EncodeNum32(uint32(n))
			out = append(out, word.Word(ValueNum), num[0], num[1])
		} else {
			out = append(out, word.Word(ValueText), word.Word(internLabel(labels, v.Value)))
		}
		return out, 1, nil
	case ast.Logic:
		op, err := codebook.LogicOpFromAST(v.Op)
		if err != nil {
			return nil, 0, err
		}
		out := []word.Word{
			word.Word(uint16(ExprLogic)<<14 | uint16(len(v.Terms))),
			word.Word(op),
		}
		return appendChildren(out, 1, v.Terms, labels)
	case ast.Not:
		out := []word.Word{word.Word(uint16(ExprNot)<<14 | 1)}
		return appendChildren(out, 1, []ast.Expr{v.Term}, labels)
	case ast.Either:
		out := []word.Word{word.Word(uint16(ExprEither)<<14 | uint16(len(v.Terms)))}
		return appendChildren(out, 1, v.Terms, labels)
	default:
		return nil, 0, fmt.Errorf("tangeul/packet: unknown ast.Expr type %T", e)
	}
}
