//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what findSee — looks up a tangl:See entry by alias
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// findSee looks up the tangl:See entry declaring alias, if any.
func findSee(doc *ast.Document, alias string) (ast.See, bool) {
	for _, s := range doc.Sees {
		if s.Alias == alias {
			return s, true
		}
	}
	return ast.See{}, false
}
