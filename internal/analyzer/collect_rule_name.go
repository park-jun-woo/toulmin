//ff:func feature=analyzer type=analyzer control=sequence
//ff:what collectRuleName — extracts rule function name from Warrant/Rebuttal/Defeater call
package analyzer

import "go/ast"

// collectRuleName extracts the first argument's identifier name from a rule registration call.
func collectRuleName(call *ast.CallExpr, dg *DefeatGraph) {
	if len(call.Args) == 0 {
		return
	}
	ident, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return
	}
	dg.Rules = append(dg.Rules, ident.Name)
}
