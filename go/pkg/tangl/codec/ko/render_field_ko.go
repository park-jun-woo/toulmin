//ff:func feature=tangl type=parser control=sequence
//ff:what renderFieldKo — render a struct definition field (`` `f` 를 가진다, 타입 <Type> ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderFieldKo renders a struct Definition field back to "`<field>` 를 가진다, 타입 <Type>",
// the inverse of parseFieldKo. The 를(을/를) particle is 받침-selected; 가진다 is followed
// by a comma before the 타입 tail.
func renderFieldKo(f ast.Field) string {
	return "`" + f.Name + "` " + josaFor(f.Name, "을", "를") + " 가진다, 타입 " + renderTypeKo(f.Type)
}
