//ff:func feature=tangl type=codegen control=sequence
//ff:what internalFuncName — builds a unique unexported name for an on-event handler
package gen

import "fmt"

// internalFuncName builds a unique unexported name for an Internal "on"
// event handler: "<kind><PascalSeed><idx>", with the Internal's document
// index always appended so a blank or repeated seed never collides.
func internalFuncName(kind, seed string, idx int) string {
	ident := goIdentExported(seed)
	if ident == "" || ident == "_" {
		return fmt.Sprintf("%s%d", kind, idx)
	}
	return fmt.Sprintf("%s%s%d", kind, ident, idx)
}
