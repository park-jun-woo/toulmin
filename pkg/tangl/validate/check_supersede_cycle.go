//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkSupersedeCycle — rejects a circular supersede graph within a case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkSupersedeCycle rejects a cyclic "`x` supersedes `y`" graph within
// a single case (node A superseding node B superseding ... node A).
func checkSupersedeCycle(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		edges := caseSupersedeEdges(c)
		lines := caseNodeLineIndex(c)
		if err := detectNameCycle(doc.Path, "supersedes", edges, func(name string) int { return lines[name] }); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
