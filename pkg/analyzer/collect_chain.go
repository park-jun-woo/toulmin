//ff:func feature=analyzer type=analyzer control=selection
//ff:what collectChain — walks method chain collecting rule registrations and defeat edges
package analyzer

import "go/ast"

// collectChain walks up a method chain collecting Warrant/Rebuttal/Defeater/Defeat calls.
func collectChain(call *ast.CallExpr, dg *DefeatGraph) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	switch sel.Sel.Name {
	case "Warrant", "Rebuttal", "Defeater":
		collectRuleName(call, dg)
	case "Defeat":
		collectDefeatEdge(call, dg)
	}
	inner, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return
	}
	collectChain(inner, dg)
}
