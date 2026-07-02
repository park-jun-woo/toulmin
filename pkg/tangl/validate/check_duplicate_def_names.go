//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicateDefNames — reports repeated tangl:Definitions terms
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkDuplicateDefNames reports every Definitions term name declared more
// than once.
func checkDuplicateDefNames(doc *ast.Document) []error {
	locs := make([]nameLoc, len(doc.Defs))
	for i, d := range doc.Defs {
		locs[i] = nameLoc{Name: d.Name, Line: d.Line}
	}
	return checkDuplicates(doc.Path, "Definitions term", locs)
}
