//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what addDocEdges — records the cross-document checking edges contributed by one document
package registry

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// addDocEdges records into edges/lines the checking edges of every case in d,
// keyed by "<docPath>#<case>". Node-level resolution is delegated to
// addCaseEdges so this stays within the shallow nesting the codebook allows.
func addDocEdges(edges map[string][]string, lines map[string]int, docPath string, d *ast.Document) error {
	for _, c := range d.Cases {
		srcKey := docPath + "#" + c.Name
		lines[srcKey] = c.Line
		if err := addCaseEdges(edges, docPath, d, c, srcKey); err != nil {
			return err
		}
	}
	return nil
}
