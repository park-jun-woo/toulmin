//ff:func feature=tangl type=validator control=sequence
//ff:what Lint — non-fatal advisories (currently: tick-idempotency guard recommendations)
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// Lint returns advisory warnings that are not validation errors. Currently
// it recommends a `once` guard (or equivalent idempotency) for every `do`
// edge in a case reachable from tangl:Internal (spec §틱 멱등성 요구).
func Lint(doc *ast.Document) []string {
	reached := internalReachableCases(doc)
	return unguardedDoWarnings(doc, reached)
}
