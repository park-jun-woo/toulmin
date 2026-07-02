//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildAliasImports — resolves every used See alias into a sorted importSpec list
package gen

import (
	"sort"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildAliasImports resolves every See alias the document's rules
// actually reference into an importSpec, sorted by alias for
// deterministic output. An alias with no matching tangl:See declaration
// is not a gen-time error — that check belongs to validate — so it falls
// back to a "tangl/<alias>" import path, the convention this module's own
// runtime companion packages (e.g. tangl/log) follow.
func buildAliasImports(doc *ast.Document) ([]importSpec, error) {
	used := collectUsedAliases(doc)
	aliases := make([]string, 0, len(used))
	for a := range used {
		aliases = append(aliases, a)
	}
	sort.Strings(aliases)
	specs := make([]importSpec, 0, len(aliases))
	for _, a := range aliases {
		path := "tangl/" + a
		if see, ok := findSee(doc, a); ok {
			path = see.Package
		}
		specs = append(specs, importSpec{Alias: a, Path: resolveSeePath(path)})
	}
	return specs, nil
}
