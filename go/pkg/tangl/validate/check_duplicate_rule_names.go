//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicateRuleNames — reports repeated tangl:Rules entry names
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkDuplicateRuleNames reports every tangl:Rules entry name declared more
// than once. The gen path already catches these as Go duplicate-identifier
// compile errors, but the eval consumer loads the document without codegen, so
// validate must catch them too.
func checkDuplicateRuleNames(doc *ast.Document) []error {
	locs := make([]nameLoc, len(doc.Rules))
	for i, r := range doc.Rules {
		locs[i] = nameLoc{Name: r.Name, Line: r.Line}
	}
	return checkDuplicates(doc.Path, "Rules entry", locs)
}
