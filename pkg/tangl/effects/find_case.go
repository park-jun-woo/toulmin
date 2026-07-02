//ff:func feature=tangl type=analyzer control=iteration dimension=1
//ff:what findCase — looks up a tangl:Cases entry by name
package effects

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// findCase returns the case named name in doc, or nil if absent.
func findCase(doc *ast.Document, name string) *ast.Case {
	for i := range doc.Cases {
		if doc.Cases[i].Name == name {
			return &doc.Cases[i]
		}
	}
	return nil
}
