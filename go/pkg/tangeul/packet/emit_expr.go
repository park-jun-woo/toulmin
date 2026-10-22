//ff:func feature=tangl type=util control=selection
//ff:what emitExpr — flatten an ast.Expr subtree into prefix-order EXPR records, interning labels
package packet

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitExpr flattens e into prefix-order EXPR records (§design 1), appending
// field/value/operand strings into *labels and returning the record words plus
// the node count of the subtree. A record header word packs the 3-bit ExprKind
// (bits 15-14 + bit 13) and the 13-bit child count (packExprHeader). Compare,
// Exists and Call have their own emitters; Logic/Not/Either share appendChildren
// for their prefix-order child records.
func emitExpr(e ast.Expr, labels *[]string) ([]core.Word, int, error) {
	switch v := e.(type) {
	case ast.Compare:
		return emitCompare(v, labels)
	case ast.Logic:
		op, err := codebook.LogicOpFromAST(v.Op)
		if err != nil {
			return nil, 0, err
		}
		header, err := packExprHeader(ExprLogic, len(v.Terms))
		if err != nil {
			return nil, 0, err
		}
		return appendChildren([]core.Word{header, core.Word(op)}, 1, v.Terms, labels)
	case ast.Not:
		header, _ := packExprHeader(ExprNot, 1)
		return appendChildren([]core.Word{header}, 1, []ast.Expr{v.Term}, labels)
	case ast.Either:
		header, err := packExprHeader(ExprEither, len(v.Terms))
		if err != nil {
			return nil, 0, err
		}
		return appendChildren([]core.Word{header}, 1, v.Terms, labels)
	case ast.Exists:
		return emitExists(v, labels)
	case ast.Call:
		return emitCall(v, labels)
	default:
		return nil, 0, fmt.Errorf("tangeul/packet: unknown ast.Expr type %T", e)
	}
}
