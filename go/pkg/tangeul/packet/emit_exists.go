//ff:func feature=tangl type=util control=sequence
//ff:what emitExists — encode an ast.Exists as an ExprExists record [Quant][source operand][binder] + Cond child
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitExists encodes an "any of <source> as <binder> satisfies <cond>"
// quantifier as an ExprExists record (child count 1): the QuantOp code, the
// source operand sub-record, the interned binder label, then the Cond subtree in
// prefix order via appendChildren (§design 3). The returned node count is 1 plus
// the Cond subtree's node count.
func emitExists(v ast.Exists, labels *[]string) ([]core.Word, int, error) {
	op, _ := codebook.QuantOpFromAST("any")
	header, _ := packExprHeader(ExprExists, 1)
	sw, err := emitOperand(v.Source, labels)
	if err != nil {
		return nil, 0, err
	}
	out := append([]core.Word{header, core.Word(op)}, sw...)
	out = append(out, core.Word(internLabel(labels, v.Binder)))
	return appendChildren(out, 1, []ast.Expr{v.Cond}, labels)
}
