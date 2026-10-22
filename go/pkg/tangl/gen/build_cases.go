//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildCases — renders every tangl:Cases entry in document order
package gen

import "strings"

// buildCases renders every tangl:Cases entry as a graph-builder function
// plus its package-level *toulmin.Graph var, in document order.
func buildCases(w *strings.Builder, gc *genContext) error {
	for _, c := range gc.Doc.Cases {
		if err := renderCase(w, gc, c); err != nil {
			return err
		}
	}
	return nil
}
