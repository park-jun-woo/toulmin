//ff:func feature=tangl type=analyzer control=iteration dimension=1
//ff:what findEndpoint — looks up a tangl:Provides endpoint by name
package effects

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// findEndpoint returns the endpoint named name in doc, or nil if absent.
func findEndpoint(doc *ast.Document, name string) *ast.Endpoint {
	for i := range doc.Provides {
		if doc.Provides[i].Name == name {
			return &doc.Provides[i]
		}
	}
	return nil
}
