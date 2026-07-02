//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what seeAliasSet — the set of all tangl:See package aliases
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// seeAliasSet returns the set of every package alias declared in doc's
// tangl:See section.
func seeAliasSet(doc *ast.Document) map[string]bool {
	set := make(map[string]bool, len(doc.Sees))
	for _, s := range doc.Sees {
		set[s.Alias] = true
	}
	return set
}
