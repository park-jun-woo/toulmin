//ff:func feature=tangl type=util control=sequence
//ff:what emitCall — encode an ast.Call in predicate position as an ExprCall record [alias][name]
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitCall encodes a host-function call used as a predicate term (ast.Call as an
// Expr) as an ExprCall record (child count 0): the interned alias label (the
// empty string for a local call) and the interned name label (§design 3). The
// value-position Call inside a Compare is handled by emitOperand instead.
func emitCall(v ast.Call, labels *[]string) ([]core.Word, int, error) {
	header, _ := packExprHeader(ExprCall, 0)
	return []core.Word{
		header,
		core.Word(internLabel(labels, v.Fn.Alias)),
		core.Word(internLabel(labels, v.Fn.Name)),
	}, 1, nil
}
