//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what addCaseEdges — appends one case's checking-target edges under its source key
package registry

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// addCaseEdges appends, for every node in c, the "<path>#<case>" key of its
// checking target (resolved by checkingEdgeTarget) to edges[srcKey]. A node
// with no checking clause contributes no edge.
func addCaseEdges(edges map[string][]string, docPath string, d *ast.Document, c ast.Case, srcKey string) error {
	for _, n := range c.Nodes {
		dst, err := checkingEdgeTarget(d, docPath, n)
		if err != nil {
			return err
		}
		if dst != "" {
			edges[srcKey] = append(edges[srcKey], dst)
		}
	}
	return nil
}
