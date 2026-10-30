//ff:func feature=tangl type=validator control=sequence
//ff:what Lint — non-fatal advisories (tick-idempotency guards, unreachable/redundant certainty gates)
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// Lint returns advisory warnings that are not validation errors: it
// recommends a `once` guard (or equivalent idempotency) for every `do` edge
// in a case reachable from tangl:Internal (spec §틱 멱등성 요구), and flags
// certainty gates outside a firing handler's reachable range (50%, 100%].
func Lint(doc *ast.Document) []string {
	reached := internalReachableCases(doc)
	return append(unguardedDoWarnings(doc, reached), certaintyGateWarnings(doc)...)
}
