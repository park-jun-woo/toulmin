//ff:func feature=tangl type=parser control=sequence
//ff:what renderField — render a struct definition field (`has `f` as <Type>`)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderField renders a struct Definition field back to "has `<field>` as <Type>",
// the inverse of the parser's parseFieldItem.
func renderField(f ast.Field) string {
	return "has `" + f.Name + "` as " + renderType(f.Type)
}
