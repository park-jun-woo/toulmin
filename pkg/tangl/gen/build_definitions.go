//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildDefinitions — renders every tangl:Definitions entry and indexes it by name
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildDefinitions renders every tangl:Definitions entry (consts, derived
// Spec vars, and struct types) into w, in document order, and returns the
// resulting name -> defInfo index for later "with" clause lookups.
func buildDefinitions(w *strings.Builder, defs []ast.Definition) map[string]defInfo {
	index := make(map[string]defInfo, len(defs))
	for _, d := range defs {
		if d.Kind == ast.StructDef {
			index[d.Name] = renderStructDef(w, d)
		} else {
			index[d.Name] = renderConstDef(w, d)
		}
	}
	return index
}
