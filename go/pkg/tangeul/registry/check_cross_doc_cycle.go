//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what CheckCrossDocCycle — rejects checking cycles that span multiple documents
package registry

import (
	"sort"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

// CheckCrossDocCycle rejects a checking-composition cycle that spans the
// documents in reg. It builds one edge graph over "<canonical-path>#<case>"
// keys (addDocEdges): same-document checking targets keep their path prefix,
// an alias-qualified target resolves to its sibling document's path, so a
// mixed cycle (A's case -> a local case -> B's case -> back to A) is caught
// alongside a pure two-document one. The generic cycle detector is reused from
// validate over deterministically sorted document paths.
func CheckCrossDocCycle(reg *Registry) error {
	paths := make([]string, 0, len(reg.Docs))
	for p := range reg.Docs {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	edges := make(map[string][]string)
	lines := make(map[string]int)
	for _, p := range paths {
		if err := addDocEdges(edges, lines, p, reg.Docs[p]); err != nil {
			return err
		}
	}
	return validate.DetectNameCycle("<cross-document>", "cross-document checking", edges, func(name string) int { return lines[name] })
}
