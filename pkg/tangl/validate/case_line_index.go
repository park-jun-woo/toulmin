//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseLineIndex — maps each case name to its declaration line
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// caseLineIndex maps every case name in doc to the line where it was declared.
func caseLineIndex(doc *ast.Document) map[string]int {
	idx := make(map[string]int, len(doc.Cases))
	for _, c := range doc.Cases {
		idx[c.Name] = c.Line
	}
	return idx
}
