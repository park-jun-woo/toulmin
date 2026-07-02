//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseNameSet — the set of all declared case names in the document
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// caseNameSet returns the set of every case name declared in doc.
func caseNameSet(doc *ast.Document) map[string]bool {
	set := make(map[string]bool, len(doc.Cases))
	for _, c := range doc.Cases {
		set[c.Name] = true
	}
	return set
}
