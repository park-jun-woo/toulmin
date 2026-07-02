//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicateSeeAlias — reports repeated tangl:See aliases
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkDuplicateSeeAlias reports every `see` alias declared more than once.
func checkDuplicateSeeAlias(doc *ast.Document) []error {
	locs := make([]nameLoc, len(doc.Sees))
	for i, s := range doc.Sees {
		locs[i] = nameLoc{Name: s.Alias, Line: s.Line}
	}
	return checkDuplicates(doc.Path, "See alias", locs)
}
