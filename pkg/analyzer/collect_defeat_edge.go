//ff:func feature=analyzer type=analyzer control=sequence
//ff:what collectDefeatEdge — extracts defeat edge from Defeat call arguments
package analyzer

import (
	"go/ast"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// collectDefeatEdge extracts from/to identifiers from a Defeat(from, to) call.
func collectDefeatEdge(call *ast.CallExpr, gc *graphCollector) {
	if len(call.Args) < 2 {
		return
	}
	fromIdent, fromOk := call.Args[0].(*ast.Ident)
	toIdent, toOk := call.Args[1].(*ast.Ident)
	if !fromOk || !toOk {
		return
	}
	gc.def.Defeats = append(gc.def.Defeats, toulmin.GraphEdgeDef{
		From: fromIdent.Name,
		To:   toIdent.Name,
	})
}
