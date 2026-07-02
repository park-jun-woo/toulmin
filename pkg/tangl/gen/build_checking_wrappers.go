//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildCheckingWrappers — renders one dedup'd wrapper per checking target case
package gen

import (
	"sort"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildCheckingWrappers scans every case node's "checking" clause, and
// for each unique target case name renders one pure wrapper function
// (dedup'd so two nodes checking the same case share it), returning the
// target case name -> generated function identifier map used by node
// registration.
func buildCheckingWrappers(w *strings.Builder, doc *ast.Document) map[string]string {
	targets := collectCheckingTargets(doc)
	names := make([]string, 0, len(targets))
	for t := range targets {
		names = append(names, t)
	}
	sort.Strings(names)
	wrappers := make(map[string]string, len(names))
	for _, t := range names {
		fn := "checking" + goIdentExported(t)
		renderCheckingWrapper(w, fn, t)
		wrappers[t] = fn
	}
	return wrappers
}
