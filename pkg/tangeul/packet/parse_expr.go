//ff:func feature=tangl type=util control=selection
//ff:what parseExpr — rebuild one ast.Expr subtree from prefix-order EXPR records, resolving labels
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseExpr rebuilds one ast.Expr subtree from the record words recs starting at
// *pos, advancing *pos past the subtree and resolving Text references through
// labels. It returns the node count of the subtree. The header word's 3-bit
// ExprKind and 13-bit child count (unpackExprHeader) drive the recursion, so the
// prefix stream reconstructs the tree uniquely (§design 1). Exists, Call and
// CompareX records dispatch to their own parsers.
func parseExpr(recs []word.Word, pos *int, labels []string) (ast.Expr, int, error) {
	if len(recs) < *pos+1 {
		return nil, 0, errors.New("tangeul/packet: EXPR record header missing")
	}
	kind, childCount := unpackExprHeader(recs[*pos])
	*pos++

	switch kind {
	case ExprCompare:
		if len(recs) < *pos+3 {
			return nil, 0, errors.New("tangeul/packet: EXPR Compare record truncated")
		}
		op, err := codebook.CompareOpToAST(codebook.CompareOp(recs[*pos]))
		if err != nil {
			return nil, 0, err
		}
		field := int(recs[*pos+1])
		vKind := ValueKind(recs[*pos+2])
		*pos += 3
		if field >= len(labels) {
			return nil, 0, errors.New("tangeul/packet: EXPR Compare field label out of range")
		}
		value, err := parseCompareValue(recs, pos, labels, vKind)
		if err != nil {
			return nil, 0, err
		}
		return ast.Compare{Field: labels[field], Op: op, Value: value}, 1, nil
	case ExprLogic:
		if len(recs) < *pos+1 {
			return nil, 0, errors.New("tangeul/packet: EXPR Logic op missing")
		}
		op, err := codebook.LogicOpToAST(codebook.LogicOp(recs[*pos]))
		if err != nil {
			return nil, 0, err
		}
		*pos++
		terms, count, err := parseChildren(recs, pos, labels, childCount)
		if err != nil {
			return nil, 0, err
		}
		return ast.Logic{Op: op, Terms: terms}, 1 + count, nil
	case ExprNot:
		terms, count, err := parseChildren(recs, pos, labels, childCount)
		if err != nil {
			return nil, 0, err
		}
		if len(terms) != 1 {
			return nil, 0, errors.New("tangeul/packet: EXPR Not must have exactly one child")
		}
		return ast.Not{Term: terms[0]}, 1 + count, nil
	case ExprEither:
		terms, count, err := parseChildren(recs, pos, labels, childCount)
		if err != nil {
			return nil, 0, err
		}
		return ast.Either{Terms: terms}, 1 + count, nil
	case ExprExists:
		return parseExists(recs, pos, labels, childCount)
	case ExprCall:
		return parseCall(recs, pos, labels)
	case ExprCompareX:
		return parseCompareX(recs, pos, labels)
	default:
		return nil, 0, errors.New("tangeul/packet: unknown ExprKind")
	}
}
