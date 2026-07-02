//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicateCaseNames — reports repeated `in case of` names
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkDuplicateCaseNames reports every case name declared more than once
// across the document.
func checkDuplicateCaseNames(doc *ast.Document) []error {
	locs := make([]nameLoc, len(doc.Cases))
	for i, c := range doc.Cases {
		locs[i] = nameLoc{Name: c.Name, Line: c.Line}
	}
	return checkDuplicates(doc.Path, "case name", locs)
}
