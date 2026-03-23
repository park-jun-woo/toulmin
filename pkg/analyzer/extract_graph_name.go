//ff:func feature=analyzer type=analyzer control=sequence
//ff:what extractGraphName — checks if a call expression contains NewGraph and returns graph name
package analyzer

import (
	"go/ast"
	"go/token"
)

// extractGraphName checks if a call expression is or contains NewGraph("name")
// and returns the graph name. Returns "" if not found.
func extractGraphName(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	if sel.Sel.Name == "NewGraph" && len(call.Args) > 0 {
		lit, ok := call.Args[0].(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			return lit.Value[1 : len(lit.Value)-1]
		}
	}
	inner, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return ""
	}
	return extractGraphName(inner)
}
