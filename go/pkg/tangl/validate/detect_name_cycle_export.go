//ff:func feature=tangl type=validator control=sequence
//ff:what DetectNameCycle — exported entry point for detectNameCycle, for cross-package reuse
package validate

// DetectNameCycle is detectNameCycle exported for reuse outside this
// package (e.g. pkg/tangeul/registry's cross-document checking graph).
// detectNameCycle is a pure algorithm over a generic map[string][]string
// edge graph — exporting it does not weaken the "validate consumes only
// *ast.Document" contract, since it takes no Document at all.
func DetectNameCycle(path, label string, edges map[string][]string, lineOf func(string) int) error {
	return detectNameCycle(path, label, edges, lineOf)
}
