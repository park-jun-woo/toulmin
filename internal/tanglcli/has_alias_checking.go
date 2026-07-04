//ff:func feature=cli type=command control=iteration dimension=2
//ff:what hasAliasChecking — reports whether any node checks an alias-qualified (cross-document) case
package tanglcli

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// hasAliasChecking reports whether any case node in d uses an alias-qualified
// "checking" clause (n.Checking.Alias != ""). Only such documents pull in the
// multi-document registry (cross-document load + cycle check); an ordinary
// document without any alias-checking never touches that path, so its CLI
// behavior is unchanged.
func hasAliasChecking(d *ast.Document) bool {
	for _, c := range d.Cases {
		for _, n := range c.Nodes {
			if n.Checking != nil && n.Checking.Alias != "" {
				return true
			}
		}
	}
	return false
}
