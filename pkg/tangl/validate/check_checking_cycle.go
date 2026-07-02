//ff:func feature=tangl type=validator control=sequence
//ff:what checkCheckingCycle — rejects circular case->case "checking" composition
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkCheckingCycle rejects a cyclic `checking <case>` composition graph
// (case A checking case B checking ... case A).
func checkCheckingCycle(doc *ast.Document) []error {
	edges := buildCheckingEdges(doc)
	lines := caseLineIndex(doc)
	if err := detectNameCycle(doc.Path, "checking", edges, func(name string) int { return lines[name] }); err != nil {
		return []error{err}
	}
	return nil
}
