//ff:func feature=tangl type=util control=sequence
//ff:what emitCompare — encode an ast.Compare as a legacy kind-0 record or, when it carries operands, an ExprCompareX record
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitCompare encodes a comparison leaf. When it carries no operands (Lhs == nil
// && Rhs == nil) it emits the legacy kind-0 record [ExprCompare][op][field]
// [value] byte-for-byte as before, so pre-existing documents keep an identical
// stream. When Lhs and/or Rhs is an operand it emits an ExprCompareX record
// [ExprCompareX][op][field][lhs present + operand?][value|operand] (§design 3):
// the field slot still carries Compare.Field (empty when Lhs is present) so the
// single-step-left case (Lhs nil, Rhs operand) round-trips, and the right side
// is a literal value or, tagged ValueOperand, an operand sub-record.
func emitCompare(v ast.Compare, labels *[]string) ([]word.Word, int, error) {
	op, err := codebook.CompareOpFromAST(v.Op)
	if err != nil {
		return nil, 0, err
	}
	if v.Lhs == nil && v.Rhs == nil {
		header, err := packExprHeader(ExprCompare, 0)
		if err != nil {
			return nil, 0, err
		}
		out := []word.Word{header, word.Word(op), word.Word(internLabel(labels, v.Field))}
		return append(out, emitCompareValue(v.Value, labels)...), 1, nil
	}
	header, err := packExprHeader(ExprCompareX, 0)
	if err != nil {
		return nil, 0, err
	}
	out := []word.Word{header, word.Word(op), word.Word(internLabel(labels, v.Field))}
	if v.Lhs != nil {
		lw, err := emitOperand(v.Lhs, labels)
		if err != nil {
			return nil, 0, err
		}
		out = append(append(out, word.Word(1)), lw...)
	} else {
		out = append(out, word.Word(0))
	}
	if v.Rhs != nil {
		rw, err := emitOperand(v.Rhs, labels)
		if err != nil {
			return nil, 0, err
		}
		return append(append(out, word.Word(ValueOperand)), rw...), 1, nil
	}
	return append(out, emitCompareValue(v.Value, labels)...), 1, nil
}
