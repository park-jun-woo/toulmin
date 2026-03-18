//ff:func feature=analyzer type=analyzer control=iteration dimension=2
//ff:what findGraphCalls — collects CallExpr nodes from var declarations in AST
package analyzer

import "go/ast"

// findGraphCalls extracts all CallExpr nodes from top-level var declarations.
func findGraphCalls(f *ast.File) []*ast.CallExpr {
	specs := collectValueSpecs(f)
	var calls []*ast.CallExpr
	for _, vs := range specs {
		for _, val := range vs.Values {
			call, ok := val.(*ast.CallExpr)
			if !ok {
				continue
			}
			calls = append(calls, call)
		}
	}
	return calls
}
