//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what collectCheckingAliases — scans the Document for every See alias used in a checking clause
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectCheckingAliases scans every case node's "checking" clause for an
// alias-qualified target and returns the set of aliases used. Unlike a
// "using" alias (a Go import path), a checking alias names a sibling TANGEUL
// document, so buildAliasImports resolves it via deriveGoImportPath rather
// than resolveSeePath.
func collectCheckingAliases(doc *ast.Document) map[string]bool {
	used := make(map[string]bool)
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking != nil && n.Checking.Alias != "" {
				used[n.Checking.Alias] = true
			}
		}
	}
	return used
}
