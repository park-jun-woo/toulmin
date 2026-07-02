//ff:func feature=tangl type=codegen control=sequence
//ff:what goIdentExported — converts a tangl backtick name to an exported PascalCase Go identifier
package gen

import "strings"

// goIdentExported converts a tangl backtick name to an exported Go
// identifier: every word is capitalized ("make coffee" becomes
// "MakeCoffee").
func goIdentExported(name string) string {
	ident := goIdent(name)
	if ident == "" || ident == "_" {
		return ident
	}
	return strings.ToUpper(ident[:1]) + ident[1:]
}
