//ff:func feature=tangl type=parser control=sequence
//ff:what renderDefinition — render a Definitions entry head (const `means <v> [as <ref>]` or struct `means`)
package en

import (
	"encoding/json"
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderDefinition renders a tangl:Definitions entry back to its head form: a
// constant is "`x` means <value>" plus an optional "as <ref>" spec clause, while
// a struct is the bare "`x` means" head (its `has` fields render as separate
// StmtBlocks). It is the inverse of the parser's parseDefinitionItem.
func renderDefinition(d ast.Definition) string {
	if d.Kind == ast.StructDef {
		return "`" + d.Name + "` means"
	}
	if d.Kind == ast.StringListDef {
		b, _ := json.Marshal(d.Strings)
		return "`" + d.Name + "` means " + string(b)
	}
	if d.Kind == ast.FileSpecDef {
		return "`" + d.Name + "` reads `" + d.LoadType + "` from " + strconv.Quote(d.FilePath)
	}
	s := "`" + d.Name + "` means " + d.Value
	if d.SpecRef != nil {
		s += " as " + renderRef(*d.SpecRef)
	}
	return s
}
