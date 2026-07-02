//ff:func feature=tangl type=codegen control=sequence
//ff:what refSelector — renders a Ref exactly as the document names it
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// refSelector renders a Ref exactly as the document names it: an aliased
// reference becomes "alias.Name" (arm.placeCup for the source
// arm.placeCup pair), a local reference passes through as "Name"
// unchanged. Ref identifiers are never camelCased — they must match an
// existing Go identifier (a See-imported function, or a hand-written or
// generated local one).
func refSelector(ref *ast.Ref) string {
	if ref == nil {
		return ""
	}
	if ref.Alias != "" {
		return ref.Alias + "." + ref.Name
	}
	return ref.Name
}
