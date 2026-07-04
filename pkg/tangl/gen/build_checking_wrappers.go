//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildCheckingWrappers — renders one dedup'd wrapper per checking target case
package gen

import (
	"sort"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildCheckingWrappers renders one exported pure wrapper function per case
// (every case is externally addressable — see collectCheckingTargets),
// returning the case name -> generated function identifier map used by node
// registration. The wrapper name uses an ASCII "Checking" prefix so it is
// exported even for a Hangul case name (whose runes have no upper case), which
// is what lets another document call `alias.Checking<Case>` across packages.
func buildCheckingWrappers(w *strings.Builder, doc *ast.Document) map[string]string {
	targets := collectCheckingTargets(doc)
	names := make([]string, 0, len(targets))
	for t := range targets {
		names = append(names, t)
	}
	sort.Strings(names)
	wrappers := make(map[string]string, len(names))
	for _, t := range names {
		fn := "Checking" + goIdentExported(t)
		renderCheckingWrapper(w, fn, t)
		wrappers[t] = fn
	}
	return wrappers
}
