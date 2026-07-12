//ff:func feature=tangl type=util control=sequence
//ff:what parseExists — rebuild an ast.Exists from an ExprExists record [Quant][source][binder] + Cond child
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseExists rebuilds an ast.Exists from an ExprExists record body at *pos: the
// QuantOp code, the source operand sub-record, the binder label index, then the
// Cond subtree (childCount is 1) via parseChildren. It advances *pos and returns
// node count 1 plus the Cond subtree's node count. A bad QuantOp, out-of-range
// binder, or wrong child count is a hard error.
func parseExists(recs []word.Word, pos *int, labels []string, childCount int) (ast.Expr, int, error) {
	if len(recs) < *pos+1 {
		return nil, 0, errors.New("tangeul/packet: EXPR Exists quantifier missing")
	}
	if _, err := codebook.QuantOpToAST(codebook.QuantOp(recs[*pos])); err != nil {
		return nil, 0, err
	}
	*pos++
	source, err := parseOperand(recs, pos, labels)
	if err != nil {
		return nil, 0, err
	}
	if len(recs) < *pos+1 {
		return nil, 0, errors.New("tangeul/packet: EXPR Exists binder missing")
	}
	binder := int(recs[*pos])
	*pos++
	if binder >= len(labels) {
		return nil, 0, errors.New("tangeul/packet: EXPR Exists binder label out of range")
	}
	terms, count, err := parseChildren(recs, pos, labels, childCount)
	if err != nil {
		return nil, 0, err
	}
	if len(terms) != 1 {
		return nil, 0, errors.New("tangeul/packet: EXPR Exists must have exactly one Cond child")
	}
	return ast.Exists{Source: source, Binder: labels[binder], Cond: terms[0]}, 1 + count, nil
}
