//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildAliasImports — resolves every used See alias into a sorted importSpec list
package gen

import (
	"sort"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildAliasImports resolves every See alias the document's rules actually
// reference into an importSpec, sorted by alias for deterministic output. The
// resolution (resolveAliasImport) branches on how the alias is used: a
// "checking" alias names a sibling TANGEUL document (import path derived from
// its location), any other "using" alias is a Go import path. A checking alias
// takes precedence when an alias is (unusually) used both ways, since the
// derived sibling-package path is what the emitted cross-package call needs.
func buildAliasImports(doc *ast.Document) ([]importSpec, error) {
	using := collectUsedAliases(doc)
	checking := collectCheckingAliases(doc)
	names := make(map[string]bool, len(using)+len(checking))
	for a := range using {
		names[a] = true
	}
	for a := range checking {
		names[a] = true
	}
	aliases := make([]string, 0, len(names))
	for a := range names {
		aliases = append(aliases, a)
	}
	sort.Strings(aliases)
	specs := make([]importSpec, 0, len(aliases))
	for _, a := range aliases {
		spec, ok, err := resolveAliasImport(doc, a, checking[a])
		if err != nil {
			return nil, err
		}
		if ok {
			specs = append(specs, spec)
		}
	}
	return specs, nil
}
