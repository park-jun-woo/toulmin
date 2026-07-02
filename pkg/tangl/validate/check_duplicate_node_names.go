//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkDuplicateNodeNames — reports repeated node names within each case
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// checkDuplicateNodeNames reports every node name declared more than once
// inside the same case (node names are scoped per case, not globally).
func checkDuplicateNodeNames(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		locs := make([]nameLoc, len(c.Nodes))
		for i, n := range c.Nodes {
			locs[i] = nameLoc{Name: n.Name, Line: n.Line}
		}
		kind := fmt.Sprintf("node name in case %q", c.Name)
		errs = append(errs, checkDuplicates(doc.Path, kind, locs)...)
	}
	return errs
}
