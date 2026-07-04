//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what collectCheckingTargets — returns the set of every case name, each an externally-addressable checking target
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectCheckingTargets returns the set of every case name in the document.
// Every case is an externally-addressable checking target (TANGEUL.en.md §Cases
// — "Case names are referenced externally as namespace.`caseName`"), so gen
// emits one exported wrapper per case whether or not this document checks it
// itself; a case referenced only from another document must still be callable.
func collectCheckingTargets(doc *ast.Document) map[string]bool {
	targets := make(map[string]bool)
	for _, c := range doc.Cases {
		targets[c.Name] = true
	}
	return targets
}
