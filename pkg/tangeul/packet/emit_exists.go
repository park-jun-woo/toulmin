//ff:func feature=tangl type=util control=sequence
//ff:what emitExists — encode an ast.Exists as an ExprExists record [Quant][source operand][binder] + Cond child
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitExists encodes an "any of <source> as <binder> satisfies <cond>"
// quantifier as an ExprExists record (child count 1): the QuantOp code, the
// source operand sub-record, the interned binder label, then the Cond subtree in
// prefix order via appendChildren (§design 3). The returned node count is 1 plus
// the Cond subtree's node count.
func emitExists(v ast.Exists, labels *[]string) ([]word.Word, int, error) {
	op, err := codebook.QuantOpFromAST("any")
	if err != nil {
		return nil, 0, err
	}
	header, err := packExprHeader(ExprExists, 1)
	if err != nil {
		return nil, 0, err
	}
	sw, err := emitOperand(v.Source, labels)
	if err != nil {
		return nil, 0, err
	}
	out := append([]word.Word{header, word.Word(op)}, sw...)
	out = append(out, word.Word(internLabel(labels, v.Binder)))
	return appendChildren(out, 1, []ast.Expr{v.Cond}, labels)
}
