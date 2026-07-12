//ff:func feature=tangl type=util control=sequence
//ff:what emitCall — encode an ast.Call in predicate position as an ExprCall record [alias][name]
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitCall encodes a host-function call used as a predicate term (ast.Call as an
// Expr) as an ExprCall record (child count 0): the interned alias label (the
// empty string for a local call) and the interned name label (§design 3). The
// value-position Call inside a Compare is handled by emitOperand instead.
func emitCall(v ast.Call, labels *[]string) ([]word.Word, int, error) {
	header, err := packExprHeader(ExprCall, 0)
	if err != nil {
		return nil, 0, err
	}
	return []word.Word{
		header,
		word.Word(internLabel(labels, v.Fn.Alias)),
		word.Word(internLabel(labels, v.Fn.Name)),
	}, 1, nil
}
