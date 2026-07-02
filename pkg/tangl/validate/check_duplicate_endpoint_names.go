//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicateEndpointNames — reports repeated tangl:Provides endpoint names
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkDuplicateEndpointNames reports every `provides` endpoint name
// declared more than once.
func checkDuplicateEndpointNames(doc *ast.Document) []error {
	locs := make([]nameLoc, len(doc.Provides))
	for i, ep := range doc.Provides {
		locs[i] = nameLoc{Name: ep.Name, Line: ep.Line}
	}
	return checkDuplicates(doc.Path, "Endpoint name", locs)
}
