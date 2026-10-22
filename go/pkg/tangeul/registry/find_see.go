//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what findSee — looks up a document's tangl:See entry by alias
package registry

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// findSee looks up the tangl:See entry declaring alias in d, if any.
func findSee(d *ast.Document, alias string) (ast.See, bool) {
	for _, s := range d.Sees {
		if s.Alias == alias {
			return s, true
		}
	}
	return ast.See{}, false
}
