//ff:func feature=tangl type=util control=sequence
//ff:what emitCase — dispatch a Cases statement to the case declaration or a body slot
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitCase routes a Cases coordinate: a one-element path is the `in case of x`
// declaration, a three-element path is a case-body slot (require/node/attack/exec).
func emitCase(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	if len(ref.Path) == 1 {
		return emitCaseDecl(e, ref, elem.(ast.Case), span)
	}
	return emitCaseSlot(e, ref, elem, span)
}
