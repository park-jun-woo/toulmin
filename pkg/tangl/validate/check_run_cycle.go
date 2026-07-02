//ff:func feature=tangl type=validator control=sequence
//ff:what checkRunCycle — rejects circular case->case run cascade composition
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkRunCycle rejects a cyclic `run <case> when <node>` cascade graph
// (case A running case B running ... case A).
func checkRunCycle(doc *ast.Document) []error {
	edges := buildRunEdges(doc)
	lines := caseLineIndex(doc)
	if err := detectNameCycle(doc.Path, "run", edges, func(name string) int { return lines[name] }); err != nil {
		return []error{err}
	}
	return nil
}
