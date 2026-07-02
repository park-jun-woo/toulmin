//ff:func feature=tangl type=codegen control=sequence
//ff:what nodeFn — resolves the rule function a case node registers
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// nodeFn resolves the rule function a case node registers: an explicit
// "using" Ref, the dedup'd "checking" wrapper for its target case, or —
// with neither — the node's own name as a bare local identifier, which Go
// resolves to a package-level function of that name (hand-written, or a
// same-named tangl:Rules entry codegen'd elsewhere in this file).
func nodeFn(gc *genContext, n ast.Node) string {
	if n.Using != nil {
		return refSelector(n.Using)
	}
	if n.Checking != "" {
		if fn, ok := gc.CheckWrappers[n.Checking]; ok {
			return fn
		}
	}
	return goIdent(n.Name)
}
