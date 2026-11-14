//ff:func feature=tangl type=util control=sequence
//ff:what aliasCheckNode — test helper building a case whose single node alias-checks a target
package registry

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// aliasCheckNode builds a case named caseName whose single node checks
// alias.target (a cross-document checking reference).
func aliasCheckNode(caseName, alias, target string) ast.Case {
	return ast.Case{
		Name:  caseName,
		Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: alias, Name: target}}},
	}
}
