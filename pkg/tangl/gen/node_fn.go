//ff:func feature=tangl type=codegen control=sequence
//ff:what nodeFn — resolves the rule function a case node registers
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// nodeFn resolves the rule function a case node registers: an explicit
// "using" Ref, an alias-qualified "checking" target rendered as a
// cross-package call to the referenced document's exported
// "alias.Checking<Case>" wrapper, the same-document "checking" wrapper for a
// local target case, or — with none of these — the node's own name as a bare
// local identifier, which Go resolves to a package-level function of that name
// (hand-written, or a same-named tangl:Rules entry codegen'd elsewhere).
func nodeFn(gc *genContext, n ast.Node) (string, error) {
	if n.Using != nil {
		return refSelector(n.Using), nil
	}
	if n.Checking != nil {
		if n.Checking.Alias != "" {
			return n.Checking.Alias + ".Checking" + goIdentExported(n.Checking.Name), nil
		}
		if fn, ok := gc.CheckWrappers[n.Checking.Name]; ok {
			return fn, nil
		}
	}
	if rule, ok := gc.Rules[n.Name]; ok && rule.Using != nil {
		return refSelector(rule.Using), nil
	}
	return goIdent(n.Name), nil
}
