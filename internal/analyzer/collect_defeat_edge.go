//ff:func feature=analyzer type=analyzer control=sequence
//ff:what collectDefeatEdge — extracts defeat edge from Defeat call arguments
package analyzer

import "go/ast"

// collectDefeatEdge extracts from/to identifiers from a Defeat(from, to) call.
func collectDefeatEdge(call *ast.CallExpr, dg *DefeatGraph) {
	if len(call.Args) < 2 {
		return
	}
	fromIdent, fromOk := call.Args[0].(*ast.Ident)
	toIdent, toOk := call.Args[1].(*ast.Ident)
	if !fromOk || !toOk {
		return
	}
	dg.Defeats[toIdent.Name] = append(dg.Defeats[toIdent.Name], fromIdent.Name)
}
