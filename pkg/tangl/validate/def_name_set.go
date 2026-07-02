//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what defNameSet — the set of all tangl:Definitions term names
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// defNameSet returns the set of every Definitions term name declared in doc.
func defNameSet(doc *ast.Document) map[string]bool {
	set := make(map[string]bool, len(doc.Defs))
	for _, d := range doc.Defs {
		set[d.Name] = true
	}
	return set
}
