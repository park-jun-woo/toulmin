//ff:func feature=tangl type=validator control=sequence
//ff:what specCapable — whether a definition can be bound as a Spec (list, file, or const with a Spec constructor)
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// specCapable reports whether d materializes to a Spec: string-list and file
// definitions always do; a const does only when it names a Spec constructor.
func specCapable(d ast.Definition) bool {
	return d.Kind == ast.StringListDef || d.Kind == ast.FileSpecDef || (d.Kind == ast.ConstDef && d.SpecRef != nil)
}
