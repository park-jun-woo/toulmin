//ff:func feature=tangl type=util control=sequence
//ff:what parseCompareX — rebuild an operand-carrying ast.Compare from an ExprCompareX record
package packet

import (
	"errors"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCompareX rebuilds an operand-carrying ast.Compare from an ExprCompareX
// record body at *pos (§design 3): the CompareOp, the field label (empty when a
// left operand is present), an lhs-present flag with an optional left operand,
// then a ValueKind whose ValueOperand tag selects a right operand and whose
// ValueNum/ValueText tags select a literal value via parseCompareValue. It
// returns node count 1. A bad op, out-of-range field, invalid flag, or bad value
// slot is a hard error.
func parseCompareX(recs []core.Word, pos *int, labels []string) (ast.Expr, int, error) {
	if len(recs) < *pos+3 {
		return nil, 0, errors.New("tangeul/packet: EXPR CompareX record truncated")
	}
	op, err := codebook.CompareOpToAST(toulmin.CompareOp(recs[*pos]))
	if err != nil {
		return nil, 0, err
	}
	field := int(recs[*pos+1])
	lhsPresent := recs[*pos+2]
	*pos += 3
	if field >= len(labels) {
		return nil, 0, errors.New("tangeul/packet: EXPR CompareX field label out of range")
	}
	cmp := ast.Compare{Field: labels[field], Op: op}
	if lhsPresent == 1 {
		lhs, err := parseOperand(recs, pos, labels)
		if err != nil {
			return nil, 0, err
		}
		cmp.Lhs = lhs
	} else if lhsPresent != 0 {
		return nil, 0, errors.New("tangeul/packet: EXPR CompareX invalid lhs presence flag")
	}
	if len(recs) < *pos+1 {
		return nil, 0, errors.New("tangeul/packet: EXPR CompareX value kind missing")
	}
	vKind := ValueKind(recs[*pos])
	*pos++
	if vKind == ValueOperand {
		rhs, err := parseOperand(recs, pos, labels)
		if err != nil {
			return nil, 0, err
		}
		cmp.Rhs = rhs
		return cmp, 1, nil
	}
	value, err := parseCompareValue(recs, pos, labels, vKind)
	if err != nil {
		return nil, 0, err
	}
	cmp.Value = value
	return cmp, 1, nil
}
